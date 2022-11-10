package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/katasec/playground/utils"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

//"github.com/pulumi/pulumi/sdk/v3/go/pulumi"

var (
	Location string
)

type PlayGroundConfig struct {
	TailScaleRoutes string
	TailScaleKey    string
}

func main() {

	// resourceGroup := "myrg"
	// clusterName := "mycluster"
	// namespace := "default"

	// yamlFile := createTempK8sYaml()
	// mycommand := fmt.Sprintf("az aks command invoke --resource-group %s --name %s --command \"kubectl apply -f %s -n %s\" --file deployment.yaml", resourceGroup, clusterName, yamlFile, namespace)

	// fmt.Println(mycommand)

	// createTempK8sYaml()

	pulumi.Run(NewDC)

}

func createTempK8sYaml() string {
	file, err := ioutil.TempFile("", "")
	utils.ExitOnError(err)

	buf, err := ioutil.ReadFile("./deployment.yaml")
	utils.ExitOnError(err)

	yaml := string(buf)

	config := readConfig()
	yaml = strings.Replace(yaml, "$(TAILSCALE_KEY)", config.TailScaleKey, 1)
	yaml = strings.Replace(yaml, "$(TAILSCALE_ROUTES)", config.TailScaleRoutes, 1)

	file.Write([]byte(yaml))
	return file.Name()
}

func readConfig() *PlayGroundConfig {
	return &PlayGroundConfig{
		TailScaleKey:    getEnvVar("TAILSCALE_KEY"),
		TailScaleRoutes: getEnvVar("TAILSCALE_ROUTES"),
	}
}

func getEnvVar(envvar string) string {
	value := os.Getenv(envvar)
	if value == "" {
		fmt.Println("Environment variable '%s' cannot be null, exitting!", envvar)
		os.Exit(1)
	}

	return value
}

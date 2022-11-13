package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func ExitOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func createTempK8sYaml() string {
	file, err := ioutil.TempFile("", "")
	ExitOnError(err)

	buf, err := ioutil.ReadFile("./deployment.yaml")
	ExitOnError(err)

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
		fmt.Printf("Environment variable '%s' cannot be null, exitting!", envvar)
		os.Exit(1)
	}

	return value
}

type PlayGroundConfig struct {
	TailScaleRoutes string
	TailScaleKey    string
}

func TestK8sDeployment() {
	// resourceGroup := "myrg"
	// clusterName := "mycluster"
	// namespace := "default"

	// config := readConfig()
	// yamlFile := createTempK8sYaml()
	// mycommand := fmt.Sprintf("az aks command invoke --resource-group %s --name %s --command \"kubectl apply -f %s -n %s\" --file deployment.yaml", resourceGroup, clusterName, yamlFile, namespace)

	// fmt.Println(mycommand)

}

package main

import "fmt"

//"github.com/pulumi/pulumi/sdk/v3/go/pulumi"

var (
	Location string
)

func main() {

	resourceGroup := "myrg"
	clusterName := "mycluster"
	namespace := "default"
	yaml := "deployment.yaml"
	mycommand := fmt.Sprintf("az aks command invoke --resource-group %s --name %s--command \"kubectl apply -f %s -n %s\" --file deployment.yaml", resourceGroup, clusterName, yaml, namespace)

	fmt.Println(mycommand)

	//pulumi.Run(NewDC)

}

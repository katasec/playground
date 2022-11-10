package main

import (
	"fmt"
	"log"

	"github.com/katasec/playground/utils"
	"github.com/pulumi/pulumi-azure-native/sdk/go/azure/containerservice"
	network "github.com/pulumi/pulumi-azure-native/sdk/go/azure/network"
	"github.com/pulumi/pulumi-azure-native/sdk/go/azure/resources"
	"github.com/pulumi/pulumi-azuread/sdk/v5/go/azuread"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func launchK8s(ctx *pulumi.Context, rg *resources.ResourceGroup, vnet *network.VirtualNetwork) {

	k8sSubnet := network.LookupSubnetOutput(ctx, network.LookupSubnetOutputArgs{
		ResourceGroupName:  rg.Name,
		SubnetName:         pulumi.String("snet-tier2-aks"),
		VirtualNetworkName: vnet.Name,
	})

	myrg, err := resources.NewResourceGroup(ctx, "rg-ea-aks01", &resources.ResourceGroupArgs{})
	ctx.Export("aks_rg", myrg.Name)
	utils.ExitOnError(err)

	sp, spPwd, err := CreateAzureServicePrincipal(ctx, "k8s")
	utils.ExitOnError(err)

	managedCluster, err := containerservice.NewManagedCluster(ctx, "aks01", &containerservice.ManagedClusterArgs{
		ResourceGroupName: myrg.Name,
		Identity: &containerservice.ManagedClusterIdentityArgs{
			Type: containerservice.ResourceIdentityTypeSystemAssigned,
		},
		ServicePrincipalProfile: &containerservice.ManagedClusterServicePrincipalProfileArgs{
			ClientId: sp.ID(),
			Secret:   spPwd.Value,
		},
		AgentPoolProfiles: containerservice.ManagedClusterAgentPoolProfileArray{
			containerservice.ManagedClusterAgentPoolProfileArgs{
				AvailabilityZones: pulumi.StringArray{
					pulumi.String("1"),
					pulumi.String("2"),
					pulumi.String("3"),
				},
				Count:              pulumi.Int(1),
				EnableNodePublicIP: pulumi.Bool(false),
				Mode:               pulumi.String("System"),
				Name:               pulumi.String("agentpool"),
				OsType:             pulumi.String("Linux"),
				Type:               pulumi.String("VirtualMachineScaleSets"),
				VmSize:             pulumi.String("Standard_B4ms"),
				VnetSubnetID:       k8sSubnet.Id(),
			},
		},
		DnsPrefix: pulumi.String("ark"),
		NetworkProfile: &containerservice.ContainerServiceNetworkProfileArgs{
			NetworkPlugin:    pulumi.String("azure"),
			NetworkPolicy:    pulumi.String("calico"),
			DockerBridgeCidr: pulumi.String("10.17.0.1/16"),
		},
		ApiServerAccessProfile: &containerservice.ManagedClusterAPIServerAccessProfileArgs{
			EnablePrivateCluster: pulumi.BoolPtr(true),
		},
	})
	utils.ExitOnError(err)
	ctx.Export("aks_name", managedCluster.Name)
}

func CreateAzureServicePrincipal(ctx *pulumi.Context, name string) (*azuread.ServicePrincipal, *azuread.ServicePrincipalPassword, error) {

	var err error
	// Create a new registered Azure AD app
	app, err := azuread.NewApplication(ctx, name+"-app", &azuread.ApplicationArgs{
		DisplayName: pulumi.String(name + "-app"),
	})
	if err != nil {
		log.Fatal(err.Error())
	}

	// Create a service principal in the App
	sp, err := azuread.NewServicePrincipal(ctx, name+"-sp", &azuread.ServicePrincipalArgs{
		ApplicationId: app.ApplicationId,
	})
	if err != nil {
		log.Fatal(err.Error())
	}

	// Create password for the service principal
	spPassword, err := azuread.NewServicePrincipalPassword(ctx, name+"-secret", &azuread.ServicePrincipalPasswordArgs{
		ServicePrincipalId: sp.ID(),
		EndDateRelative:    pulumi.StringPtr("8760h"),
		DisplayName:        pulumi.StringPtr(fmt.Sprintf("%v-password", name+"-secret")),
	})

	ctx.Export("spPassword", spPassword)
	if err != nil {
		log.Fatal(err.Error())
	}

	return sp, spPassword, err
}

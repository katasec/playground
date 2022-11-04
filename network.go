package main

import (
	"fmt"

	"github.com/katasec/playground/azuredc"
	"github.com/katasec/playground/utils"
	containerinstance "github.com/pulumi/pulumi-azure-native/sdk/go/azure/containerinstance"
	network "github.com/pulumi/pulumi-azure-native/sdk/go/azure/network"
	"github.com/pulumi/pulumi-azure-native/sdk/go/azure/resources"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	// github.com/pulumi/pulumi-azure/sdk/v4/go/azure
	// github.com/pulumi/pulumi-azure/sdk/v4/go/azure
)

// NewDC creates a new data centre based on a reference azuredc
func NewDC(ctx *pulumi.Context) error {

	// Create hub resource group and VNET
	hubrg, err := resources.NewResourceGroup(ctx, "rg-play-hub-", &resources.ResourceGroupArgs{})
	utils.ExitOnError(err)
	hubVnet, firewall := CreateHub(ctx, hubrg, &azuredc.ReferenceHubVNET)

	AddSpoke(ctx, "nprod", hubrg, hubVnet, firewall, 0)
	AddSpoke(ctx, "prod", hubrg, hubVnet, firewall, 1)
	AddSpoke(ctx, "nprod2", hubrg, hubVnet, firewall, 2)

	return err
}

// Creates an Azure Virtual Network and subnets using the provided VNETInfo
func CreateHub(ctx *pulumi.Context, rg *resources.ResourceGroup, vnetInfo *azuredc.VNETInfo) (*network.VirtualNetwork, *network.AzureFirewall) {

	// Generate list of subnets to create
	subnets := network.SubnetTypeArray{}
	for _, subnet := range vnetInfo.SubnetsInfo {
		subnets = append(subnets, network.SubnetTypeArgs{
			AddressPrefix: pulumi.String(subnet.AddressPrefix),
			Name:          pulumi.String(subnet.Name),
		})
	}

	// Create VNET + subnets
	vnet, err := network.NewVirtualNetwork(ctx, vnetInfo.Name, &network.VirtualNetworkArgs{
		AddressSpace: &network.AddressSpaceArgs{
			AddressPrefixes: pulumi.StringArray{
				pulumi.String(vnetInfo.AddressPrefix),
			},
		},
		ResourceGroupName: rg.Name,
		Subnets:           &subnets,
	})

	// Create Firewall
	firewall := createFirewall(ctx, rg, vnet)
	utils.ExitOnError(err)

	return vnet, firewall
}

func AddSpoke(ctx *pulumi.Context, spokeName string, hubrg *resources.ResourceGroup, hubVnet *network.VirtualNetwork, firewall *network.AzureFirewall, offset int) {

	// Create a resource group
	rgName := fmt.Sprintf("%s-%s-", azuredc.RgPrefix, spokeName)
	spokeRg, err := resources.NewResourceGroup(ctx, rgName, &resources.ResourceGroupArgs{})
	utils.ExitOnError(err)

	// Create a route to firewall
	routeName := fmt.Sprintf("rt-%s", spokeName)
	spokeRoute := createFWRoute(ctx, spokeRg, routeName, firewall)

	// Generate CIDRs
	spokeCidrs := azuredc.NewSpokeVnetTemplate(spokeName, offset)

	// Create VNET using generated ResourceGroup, Route & CIDRs
	nprodVnet := CreateVNET(ctx, spokeRg, spokeCidrs, spokeRoute)

	// Peer hub with spoke
	pulumiUrn1 := fmt.Sprintf("hub-to-%s", spokeName)
	peerNetworks(ctx, pulumiUrn1, hubrg, hubVnet, nprodVnet)

	// Peer spoke with hub
	pulumiUrn2 := fmt.Sprintf("%s-to-hub", spokeName)
	peerNetworks(ctx, pulumiUrn2, spokeRg, nprodVnet, hubVnet)

}

func StartContainers(ctx *pulumi.Context, rg *resources.ResourceGroup, hub *network.VirtualNetwork) {
	// containerinstance.NewContainerGroup(ctx, "bash", &containerinstance.ContainerGroupArgs{
	// 	ContainerGroupName: pulumi.String("bash"),
	// 	ResourceGroupName:  rg.Name,
	// 	Containers: containerinstance.ContainerArray{
	// 		containerinstance.ContainerArgs{
	// 			Image: pulumi.String("registry.hub.docker.com/library/bash"),
	// 			Command: pulumi.StringArray{
	// 				pulumi.String("/usr/local/bin/bash"),
	// 				pulumi.String("-c"),
	// 				pulumi.String("echo hello; sleep 100000"),
	// 			},
	// 		},
	// 	},
	// })

	_, err := containerinstance.NewContainerGroup(ctx, "bash", &containerinstance.ContainerGroupArgs{
		OsType:             pulumi.String("Linux"),
		ContainerGroupName: pulumi.String("bash"),
		ResourceGroupName:  rg.Name,
		Containers: containerinstance.ContainerArray{
			containerinstance.ContainerArgs{
				Name: pulumi.String("hello-world"),
				Resources: containerinstance.ResourceRequirementsArgs{
					Limits: containerinstance.ResourceLimitsArgs{
						Cpu:        pulumi.Float64(0.5),
						MemoryInGB: pulumi.Float64(1.5),
					},
					Requests: containerinstance.ResourceRequestsArgs{
						Cpu:        pulumi.Float64(0.5),
						MemoryInGB: pulumi.Float64(1.5),
					},
				},
				Image: pulumi.String("registry.hub.docker.com/library/bash"),
				Command: pulumi.StringArray{
					pulumi.String("/usr/local/bin/bash"),
					pulumi.String("-c"),
					pulumi.String("echo hello; sleep 100000"),
				},
			},
		},
		NetworkProfile: containerinstance.ContainerGroupNetworkProfileArgs{
			Id: pulumi.String("/subscriptions/174c6cc1-faef-4e40-91f4-1bef3a703153/resourceGroups/rg-play-hub-cfe74535/providers/Microsoft.Network/virtualNetworks/hub60b0830a/subnets/snet-hub/aci"),
		},
	})
	utils.ExitOnError(err)

}

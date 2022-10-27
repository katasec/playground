package main

import (
	"github.com/katasec/playground/utils"
	network "github.com/pulumi/pulumi-azure-native/sdk/go/azure/network"
	"github.com/pulumi/pulumi-azure-native/sdk/go/azure/resources"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func Start(ctx *pulumi.Context) error {

	// Create Resource Group
	resourceGroup, err := resources.NewResourceGroup(ctx, "resourceGroup", &resources.ResourceGroupArgs{
		ResourceGroupName: pulumi.String("myResourceGroup"),
	})
	utils.ExitOnError(err)

	// Create VNET
	_, err = network.NewVirtualNetwork(ctx, "virtualNetwork", &network.VirtualNetworkArgs{
		AddressSpace: &network.AddressSpaceArgs{
			AddressPrefixes: pulumi.StringArray{
				pulumi.String("172.17.0.0/16"),
			},
		},
		Subnets: []network.SubnetTypeArgs{
			&network.SubnetTypeArgs{
				AddressPrefix: pulumi.String("10.0.0.0/24"),
				Name:          pulumi.String("test-1"),
			},
		},
		ResourceGroupName:  resourceGroup.Name,
		VirtualNetworkName: pulumi.String("test-vnet"),
	})
	utils.ExitOnError(err)

	return err
}

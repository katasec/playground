package main

import (
	"log"

	network "github.com/pulumi/pulumi-azure-native/sdk/go/azure/network"
	"github.com/pulumi/pulumi-azure-native/sdk/go/azure/resources"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func Start(ctx *pulumi.Context) error {

	resourceGroup, err := resources.NewResourceGroup(ctx, "resourceGroup", &resources.ResourceGroupArgs{
		ResourceGroupName: pulumi.String("myResourceGroup"),
	})
	if err != nil {
		log.Fatal(err.Error())
	}

	_, err = network.NewVirtualNetwork(ctx, "virtualNetwork", &network.VirtualNetworkArgs{
		AddressSpace: &network.AddressSpaceArgs{
			AddressPrefixes: pulumi.StringArray{
				pulumi.String("10.0.0.0/16"),
			},
		},
		ResourceGroupName:  resourceGroup.Name,
		VirtualNetworkName: pulumi.String("test-vnet"),
	})

	return err
}

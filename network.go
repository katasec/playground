package main

import (
	"github.com/katasec/playground/lookup"
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
	vnetTemplate := lookup.NewVnetTemplate()

	vnet, err := network.NewVirtualNetwork(ctx, "virtualNetwork", &network.VirtualNetworkArgs{
		AddressSpace: &network.AddressSpaceArgs{
			AddressPrefixes: pulumi.StringArray{
				pulumi.String(vnetTemplate.AddressPrefix),
			},
		},
		ResourceGroupName:  resourceGroup.Name,
		VirtualNetworkName: pulumi.String("test-vnet"),
	})
	utils.ExitOnError(err)

	// Create Subnets
	for _, subnet := range vnetTemplate.Subnets {
		_, err = network.NewSubnet(ctx, subnet.Name, &network.SubnetArgs{
			ResourceGroupName:  resourceGroup.Name,
			AddressPrefix:      pulumi.String(subnet.AddressPrefix),
			VirtualNetworkName: vnet.Name,
		})
		utils.ExitOnError(err)
	}

	return err
}

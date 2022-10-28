package main

import (
	"github.com/katasec/playground/azuredc"
	"github.com/katasec/playground/utils"
	network "github.com/pulumi/pulumi-azure-native/sdk/go/azure/network"
	"github.com/pulumi/pulumi-azure-native/sdk/go/azure/resources"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// NewDC creates a new data centre based on a reference azuredc
func NewDC(ctx *pulumi.Context) error {

	// Create Resource Group
	resourceGroup, err := resources.NewResourceGroup(ctx, "resourceGroup", &resources.ResourceGroupArgs{
		ResourceGroupName: pulumi.String("myResourceGroup"),
	})
	utils.ExitOnError(err)

	// Get a template for the VNET we want to create.
	nprodSpoke := azuredc.NewVnetTemplate("nprod")

	CreateVNET(ctx, resourceGroup, nprodSpoke)

	return err
}

// Creates an Azure Virtual Network and subnets using the provided VNETInfo
func CreateVNET(ctx *pulumi.Context, rg *resources.ResourceGroup, info *azuredc.VNETInfo) *network.VirtualNetwork {
	// Create VNET
	vnet, err := network.NewVirtualNetwork(ctx, "virtualNetwork", &network.VirtualNetworkArgs{
		AddressSpace: &network.AddressSpaceArgs{
			AddressPrefixes: pulumi.StringArray{
				pulumi.String(info.AddressPrefix),
			},
		},
		ResourceGroupName:  rg.Name,
		VirtualNetworkName: pulumi.String(info.Name),
	})
	utils.ExitOnError(err)

	// Create Subnets
	for _, subnet := range info.SubnetsInfo {
		_, err = network.NewSubnet(ctx, subnet.Name, &network.SubnetArgs{
			ResourceGroupName:  rg.Name,
			AddressPrefix:      pulumi.String(subnet.AddressPrefix),
			VirtualNetworkName: vnet.Name,
		})
		utils.ExitOnError(err)
	}

	return vnet
}

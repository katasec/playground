package main

import (
	"github.com/katasec/playground/utils"
	network "github.com/pulumi/pulumi-azure-native/sdk/go/azure/network"
	"github.com/pulumi/pulumi-azure-native/sdk/go/azure/resources"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type SubnetInfo struct {
	Name            string
	AddressPrefixes []string
	Description     string
}

func GetSubnets() *[]SubnetInfo {
	return &[]SubnetInfo{
		{
			Name:            "snet-tier1-agw",
			Description:     "Subnet for AGW",
			AddressPrefixes: []string{"172.17.1.0/24"},
		},
		{
			Name:            "snet-tier1-webin",
			Description:     "Subnet for other LBs",
			AddressPrefixes: []string{"172.17.1.0/24"},
		},
		{
			Name:            "snet-tier1-rsvd1",
			Description:     "Tier 1 reserved subnet",
			AddressPrefixes: []string{"172.17.1.0/24"},
		},
		{
			Name:            "snet-tier1-rsvd1",
			Description:     "Tier 1 reserved subnet",
			AddressPrefixes: []string{"172.17.1.0/24"},
		},
	}
}
func Start(ctx *pulumi.Context) error {

	// Create Resource Group
	resourceGroup, err := resources.NewResourceGroup(ctx, "resourceGroup", &resources.ResourceGroupArgs{
		ResourceGroupName: pulumi.String("myResourceGroup"),
	})
	utils.ExitOnError(err)

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

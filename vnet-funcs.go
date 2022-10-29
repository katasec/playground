package main

import (
	"github.com/katasec/playground/azuredc"
	"github.com/katasec/playground/utils"
	network "github.com/pulumi/pulumi-azure-native/sdk/go/azure/network"
	"github.com/pulumi/pulumi-azure-native/sdk/go/azure/resources"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Creates an Azure Virtual Network and subnets using the provided VNETInfo
func CreateVNET(ctx *pulumi.Context, rg *resources.ResourceGroup, vnetInfo *azuredc.VNETInfo) *network.VirtualNetwork {
	// Create VNET
	vnet, err := network.NewVirtualNetwork(ctx, vnetInfo.Name, &network.VirtualNetworkArgs{
		AddressSpace: &network.AddressSpaceArgs{
			AddressPrefixes: pulumi.StringArray{
				pulumi.String(vnetInfo.AddressPrefix),
			},
		},
		ResourceGroupName: rg.Name,
		// Subnets:           network.SubnetTypeArray{

		// },
	})
	utils.ExitOnError(err)

	// Create Subnets
	var previousSubnet *network.Subnet
	var dependsOn pulumi.ResourceOption

	for _, subnet := range vnetInfo.SubnetsInfo {

		// Add previously create subnet as a dependency for the next subnet
		// Avoids race conditions during create/destroy
		if previousSubnet != nil {
			dependsOn = pulumi.DependsOn([]pulumi.Resource{previousSubnet})
		} else {
			dependsOn = nil
		}

		// Create subnet
		current, err := network.NewSubnet(ctx, vnetInfo.Name+"-"+subnet.Name, &network.SubnetArgs{
			ResourceGroupName:  rg.Name,
			AddressPrefix:      pulumi.String(subnet.AddressPrefix),
			VirtualNetworkName: vnet.Name,
		}, dependsOn)
		utils.ExitOnError(err)

		previousSubnet = current

	}

	return vnet
}

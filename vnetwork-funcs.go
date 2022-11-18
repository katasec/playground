package main

import (
	"github.com/katasec/playground/azuredc"
	"github.com/katasec/playground/utils"
	network "github.com/pulumi/pulumi-azure-native/sdk/go/azure/network"
	"github.com/pulumi/pulumi-azure-native/sdk/go/azure/resources"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Creates an Azure Virtual Network and subnets using the provided VNETInfo
func CreateVNET(ctx *pulumi.Context, rg *resources.ResourceGroup, vnetInfo *azuredc.VNETInfo, routeTable *network.RouteTable) *network.VirtualNetwork {

	// Generate list of subnets to create
	subnets := network.SubnetTypeArray{}
	for _, subnet := range vnetInfo.SubnetsInfo {
		subnets = append(subnets, network.SubnetTypeArgs{
			AddressPrefix: pulumi.String(subnet.AddressPrefix),
			Name:          pulumi.String(subnet.Name),
			RouteTable: network.RouteTableTypeArgs{
				Id: routeTable.ID(),
			},
		})
	}

	// -- Create VNET
	vnet, err := network.NewVirtualNetwork(ctx, vnetInfo.Name, &network.VirtualNetworkArgs{
		AddressSpace: &network.AddressSpaceArgs{
			AddressPrefixes: pulumi.StringArray{
				pulumi.String(vnetInfo.AddressPrefix),
			},
		},
		ResourceGroupName: rg.Name,
		Subnets:           &subnets,
	})

	utils.ExitOnError(err)

	return vnet
}

func peerNetworks(ctx *pulumi.Context, pulumiUrn string, srcRg *resources.ResourceGroup, src *network.VirtualNetwork, dst *network.VirtualNetwork) {
	peeringName := pulumi.Sprintf("%s-to-%s", src.Name, dst.Name)
	_, err := network.NewVirtualNetworkPeering(ctx, pulumiUrn, &network.VirtualNetworkPeeringArgs{
		Name:                      peeringName,
		VirtualNetworkPeeringName: peeringName,
		ResourceGroupName:         srcRg.Name,
		VirtualNetworkName:        src.Name,

		AllowForwardedTraffic:     pulumi.Bool(true),
		AllowGatewayTransit:       pulumi.Bool(false),
		AllowVirtualNetworkAccess: pulumi.Bool(true),
		RemoteVirtualNetwork: &network.SubResourceArgs{
			Id: dst.ID(),
		},
		UseRemoteGateways: pulumi.Bool(false),
	})
	utils.ExitOnError(err)
}

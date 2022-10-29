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

	// Create hub resource group and VNET
	hubrg, err := resources.NewResourceGroup(ctx, "rg-play-hubrg-", &resources.ResourceGroupArgs{})
	utils.ExitOnError(err)
	hubVnet := CreateVNET(ctx, hubrg, &azuredc.ReferenceHubVNET)
	//CreateVNET(ctx, hubrg, &azuredc.ReferenceHubVNET)

	// Create nprod resource group and VNET
	nprodrg, err := resources.NewResourceGroup(ctx, "rg-play-nprodrg-", &resources.ResourceGroupArgs{})
	utils.ExitOnError(err)
	nprodCidrs := azuredc.NewSpokeVnetTemplate("nprod")
	nprodVnet := CreateVNET(ctx, nprodrg, nprodCidrs)
	//CreateVNET(ctx, nprodrg, nprodCidrs)

	// Create prod resource group and VNET
	prodResGroup, err := resources.NewResourceGroup(ctx, "rg-play-prodrg-", &resources.ResourceGroupArgs{})
	utils.ExitOnError(err)
	prodCidrs := azuredc.NewSpokeVnetTemplate("prod", 1)
	prodVnet := CreateVNET(ctx, prodResGroup, prodCidrs)

	peerNetworks(ctx, "hub-to-nprod", hubrg, hubVnet, nprodVnet)
	peerNetworks(ctx, "nprod-to-hub", nprodrg, nprodVnet, hubVnet)

	peerNetworks(ctx, "hub-to-prod", hubrg, hubVnet, prodVnet)
	peerNetworks(ctx, "prod-to-hub", prodResGroup, prodVnet, hubVnet)

	return err
}

func peerNetworks(ctx *pulumi.Context, urn string, srcRg *resources.ResourceGroup, src *network.VirtualNetwork, dst *network.VirtualNetwork) {
	name := pulumi.Sprintf("%s-to-%s", src.Name, dst.Name)
	_, err := network.NewVirtualNetworkPeering(ctx, urn, &network.VirtualNetworkPeeringArgs{
		Name:                      name,
		VirtualNetworkPeeringName: name,
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

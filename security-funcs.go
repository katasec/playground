package main

import (
	"github.com/katasec/playground/utils"
	network "github.com/pulumi/pulumi-azure-native/sdk/go/azure/network"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func CreateNSG(ctx *pulumi.Context, rg string, location string, nsgName string) *network.NetworkSecurityGroup {

	nsg, err := network.NewNetworkSecurityGroup(ctx, nsgName, &network.NetworkSecurityGroupArgs{
		ResourceGroupName: pulumi.String(rg),
	})

	utils.ExitOnError(err)

	return nsg
}

package azuredc

var (
	RgPrefix              = "rg-play"
	octetStart            = 17
	referenceSpokeSubnets = []SubnetInfo{
		{
			Name:          "snet-tier1-agw",
			Description:   "Subnet for AGW",
			AddressPrefix: "172.x.1.0/24",
			Tags: map[string]string{
				"snet:role": "tier1-agw",
			},
		},
		{
			Name:          "snet-tier1-webin",
			Description:   "Subnet for other LBs",
			AddressPrefix: "172.x.2.0/24",
			Tags: map[string]string{
				"snet:role": "tier1-webin",
			},
		},
		{
			Name:          "snet-tier1-rsvd1",
			Description:   "Tier 1 reserved subnet",
			AddressPrefix: "172.x.3.0/25",
			Tags: map[string]string{
				"snet:role": "tier1-rsvd1",
			},
		},
		{
			Name:          "snet-tier1-rsvd2",
			Description:   "Tier 1 reserved subnet",
			AddressPrefix: "172.x.3.128/25",
			Tags: map[string]string{
				"snet:role": "tier1-rsvd2",
			},
		},
		{
			Name:          "snet-tier2-wbapp",
			Description:   "Subnet for web apps",
			AddressPrefix: "172.x.4.0/23",
			Tags: map[string]string{
				"snet:role": "tier2-wbapp",
			},
		},
		{
			Name:          "snet-tier2-rsvd2",
			Description:   "Tier 2 reserved subnet",
			AddressPrefix: "172.x.6.0/24",
		},
		{
			Name:          "snet-tier2-pckr",
			Description:   "Subnet for packer",
			AddressPrefix: "172.x.7.0/24",
		},
		{
			Name:          "snet-tier2-vm",
			Description:   "Subnet for VMs",
			AddressPrefix: "172.x.8.0/21",
			Tags: map[string]string{
				"snet:role": "tier2-vm",
			},
		},
		{
			Name:          "snet-tier2-aks",
			Description:   "Subnet for AKS",
			AddressPrefix: "172.x.16.0/20",
		},
		{
			Name:          "snet-tier3-mi",
			Description:   "Subnet for managed instance",
			AddressPrefix: "172.x.32.0/26",
		},
		{
			Name:          "snet-tier3-dbaz",
			Description:   "Subnet for SQL Azure",
			AddressPrefix: "172.x.32.64/26",
		},
		{
			Name:          "snet-tier3-dblb",
			Description:   "Subnet for LB for SQL VM",
			AddressPrefix: "172.x.32.128/25",
		},
		{
			Name:          "snet-tier3-dbvm",
			Description:   "Subnet for SQL VM",
			AddressPrefix: "172.x.33.0/25",
		},
		{
			Name:          "snet-tier3-strg",
			Description:   "Subnet for storage account/fileshares",
			AddressPrefix: "172.x.33.128/25",
		},
		{
			Name:          "snet-tier3-redis",
			Description:   "Subnet for redis cache",
			AddressPrefix: "172.x.34.0/25",
		},
	}

	referenceHubSubnets = []SubnetInfo{
		{
			Name:          "AzureFirewallSubnet",
			Description:   "Subnet for Azure Firewall",
			AddressPrefix: "172.16.0.0/26",
		},
		{
			Name:          "AzureBastionSubnet",
			Description:   "Subnet for Bastion",
			AddressPrefix: "172.16.0.64/26",
		},
		{
			Name:          "AzureFirewallManagementSubnet",
			Description:   "Subnet for VPN Gateway",
			AddressPrefix: "172.16.0.128/26",
		},
		{
			Name:          "GatewaySubnet",
			Description:   "Subnet for VPN Gateway",
			AddressPrefix: "172.16.0.192/27",
		},
		{
			Name:          "snet-test",
			Description:   "Subnet for VPN Gateway",
			AddressPrefix: "172.16.0.224/27",
		},
	}

	// Template for creating spoke networks
	ReferenceSpokeVNET = VNETInfo{
		AddressPrefix: "172.x.0.0/16",
		SubnetsInfo:   referenceSpokeSubnets,
	}

	// Template for creating a hub network
	ReferenceHubVNET = VNETInfo{
		Name:          "hub",
		AddressPrefix: "172.16.0.0/24",
		SubnetsInfo:   referenceHubSubnets,
	}
)

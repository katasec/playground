package lookup

type SubnetInfo struct {
	Name          string
	AddressPrefix string
	Description   string
}

var (
	SpokeSubnets = []SubnetInfo{
		{
			Name:          "snet-tier1-agw",
			Description:   "Subnet for AGW",
			AddressPrefix: "172.17.1.0/24",
		},
		{
			Name:          "snet-tier1-webin",
			Description:   "Subnet for other LBs",
			AddressPrefix: "172.17.2.0/24",
		},
		{
			Name:          "snet-tier1-rsvd1",
			Description:   "Tier 1 reserved subnet",
			AddressPrefix: "172.17.3.0/25",
		},
		{
			Name:          "snet-tier1-rsvd2",
			Description:   "Tier 1 reserved subnet",
			AddressPrefix: "172.17.3.128/25",
		},
		{
			Name:          "snet-tier2-wbapp",
			Description:   "Subnet for web apps",
			AddressPrefix: "172.17.4.0/23",
		},
		{
			Name:          "snet-tier2-rsvd2",
			Description:   "Tier 2 reserved subnet",
			AddressPrefix: "172.17.6.0/24",
		},
		{
			Name:          "snet-tier2-pckr",
			Description:   "Subnet for packer",
			AddressPrefix: "172.17.7.0/24",
		},
		{
			Name:          "snet-tier2-vm",
			Description:   "Subnet for VMs",
			AddressPrefix: "172.17.8.0/21",
		},
		{
			Name:          "snet-tier2-aks",
			Description:   "Subnet for AKS",
			AddressPrefix: "172.17.16.0/20",
		},
		{
			Name:          "snet-tier3-mi",
			Description:   "Subnet for managed instance",
			AddressPrefix: "172.17.32.0/26",
		},
		{
			Name:          "snet-tier3-dbaz",
			Description:   "Subnet for SQL Azure",
			AddressPrefix: "172.17.32.64/26",
		},
		{
			Name:          "snet-tier3-dblb",
			Description:   "Subnet for LB for SQL VM",
			AddressPrefix: "172.17.32.128/25",
		},
		{
			Name:          "snet-tier3-dbvm",
			Description:   "Subnet for SQL VM",
			AddressPrefix: "172.17.33.0/25",
		},
		{
			Name:          "snet-tier3-strg",
			Description:   "Subnet for storage account/fileshares",
			AddressPrefix: "172.17.33.128/25",
		},
		{
			Name:          "snet-tier3-redis",
			Description:   "Subnet for redis cache",
			AddressPrefix: "172.17.34.0/25",
		},
	}
)

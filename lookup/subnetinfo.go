package lookup

type SubnetInfo struct {
	Name            string
	AddressPrefixes []string
	Description     string
}

var (
	SpokeSubnets = []SubnetInfo{
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
)

// func (s *SubnetInfo) SpokeList() *[]SubnetInfo {
// 	return &[]SubnetInfo{
// 		{
// 			Name:            "snet-tier1-agw",
// 			Description:     "Subnet for AGW",
// 			AddressPrefixes: []string{"172.17.1.0/24"},
// 		},
// 		{
// 			Name:            "snet-tier1-webin",
// 			Description:     "Subnet for other LBs",
// 			AddressPrefixes: []string{"172.17.1.0/24"},
// 		},
// 		{
// 			Name:            "snet-tier1-rsvd1",
// 			Description:     "Tier 1 reserved subnet",
// 			AddressPrefixes: []string{"172.17.1.0/24"},
// 		},
// 		{
// 			Name:            "snet-tier1-rsvd1",
// 			Description:     "Tier 1 reserved subnet",
// 			AddressPrefixes: []string{"172.17.1.0/24"},
// 		},
// 	}
// }

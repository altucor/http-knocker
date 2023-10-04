package firewallField

import (
	"errors"
	"strings"
)

type ProtocolType uint16

// Extracted from: /etc/protocols
const (
	PROTOCOL_INVALID         ProtocolType = 0xFFFF
	PROTOCOL_IP              ProtocolType = 0
	PROTOCOL_ICMP            ProtocolType = 1
	PROTOCOL_IGMP            ProtocolType = 2
	PROTOCOL_GGP             ProtocolType = 3
	PROTOCOL_IPENCAP         ProtocolType = 4
	PROTOCOL_ST2             ProtocolType = 5
	PROTOCOL_TCP             ProtocolType = 6
	PROTOCOL_CBT             ProtocolType = 7
	PROTOCOL_EGP             ProtocolType = 8
	PROTOCOL_IGP             ProtocolType = 9
	PROTOCOL_BBN_RCC         ProtocolType = 10
	PROTOCOL_NVP             ProtocolType = 11
	PROTOCOL_PUP             ProtocolType = 12
	PROTOCOL_ARGUS           ProtocolType = 13
	PROTOCOL_EMCON           ProtocolType = 14
	PROTOCOL_XNET            ProtocolType = 15
	PROTOCOL_CHAOS           ProtocolType = 16
	PROTOCOL_UDP             ProtocolType = 17
	PROTOCOL_MUX             ProtocolType = 18
	PROTOCOL_DCN             ProtocolType = 19
	PROTOCOL_HMP             ProtocolType = 20
	PROTOCOL_PRM             ProtocolType = 21
	PROTOCOL_XNS_IDP         ProtocolType = 22
	PROTOCOL_TRUNK_1         ProtocolType = 23
	PROTOCOL_TRUNK_2         ProtocolType = 24
	PROTOCOL_LEAF_1          ProtocolType = 25
	PROTOCOL_LEAF_2          ProtocolType = 26
	PROTOCOL_RDP             ProtocolType = 27
	PROTOCOL_IRTP            ProtocolType = 28
	PROTOCOL_ISO_TP4         ProtocolType = 29
	PROTOCOL_NETBLT          ProtocolType = 30
	PROTOCOL_MFE_NSP         ProtocolType = 31
	PROTOCOL_MERIT_INP       ProtocolType = 32
	PROTOCOL_DCCP            ProtocolType = 33
	PROTOCOL_THREE_PC        ProtocolType = 34
	PROTOCOL_IDPR            ProtocolType = 35
	PROTOCOL_XTP             ProtocolType = 36
	PROTOCOL_DDP             ProtocolType = 37
	PROTOCOL_IDPR_CMTP       ProtocolType = 38
	PROTOCOL_TPPLUSPLUS      ProtocolType = 39
	PROTOCOL_IL              ProtocolType = 40
	PROTOCOL_IPV6            ProtocolType = 41
	PROTOCOL_SDRP            ProtocolType = 42
	PROTOCOL_IPV6_ROUTE      ProtocolType = 43
	PROTOCOL_IPV6_FRAG       ProtocolType = 44
	PROTOCOL_IDRP            ProtocolType = 45
	PROTOCOL_RSVP            ProtocolType = 46
	PROTOCOL_GRE             ProtocolType = 47
	PROTOCOL_DSR             ProtocolType = 48
	PROTOCOL_BNA             ProtocolType = 49
	PROTOCOL_ESP             ProtocolType = 50
	PROTOCOL_AH              ProtocolType = 51
	PROTOCOL_I_NLSP          ProtocolType = 52
	PROTOCOL_SWIPE           ProtocolType = 53
	PROTOCOL_NARP            ProtocolType = 54
	PROTOCOL_MOBILE          ProtocolType = 55
	PROTOCOL_TLSP            ProtocolType = 56
	PROTOCOL_SKIP            ProtocolType = 57
	PROTOCOL_IPV6_ICMP       ProtocolType = 58
	PROTOCOL_IPV6_NONXT      ProtocolType = 59
	PROTOCOL_IPV6_OPTS       ProtocolType = 60
	PROTOCOL_CFTP            ProtocolType = 62
	PROTOCOL_SAT_EXPAK       ProtocolType = 64
	PROTOCOL_KRYPTOLAN       ProtocolType = 65
	PROTOCOL_RVD             ProtocolType = 66
	PROTOCOL_IPPC            ProtocolType = 67
	PROTOCOL_SAT_MON         ProtocolType = 69
	PROTOCOL_VISA            ProtocolType = 70
	PROTOCOL_IPCV            ProtocolType = 71
	PROTOCOL_CPNX            ProtocolType = 72
	PROTOCOL_CPHB            ProtocolType = 73
	PROTOCOL_WSN             ProtocolType = 74
	PROTOCOL_PVP             ProtocolType = 75
	PROTOCOL_BR_SAT_MON      ProtocolType = 76
	PROTOCOL_SUN_ND          ProtocolType = 77
	PROTOCOL_WB_MON          ProtocolType = 78
	PROTOCOL_WB_EXPAK        ProtocolType = 79
	PROTOCOL_ISO_IP          ProtocolType = 80
	PROTOCOL_VMTP            ProtocolType = 81
	PROTOCOL_SECURE_VMTP     ProtocolType = 82
	PROTOCOL_VINES           ProtocolType = 83
	PROTOCOL_TTP             ProtocolType = 84
	PROTOCOL_NSFNET_IGP      ProtocolType = 85
	PROTOCOL_DGP             ProtocolType = 86
	PROTOCOL_TCF             ProtocolType = 87
	PROTOCOL_EIGRP           ProtocolType = 88
	PROTOCOL_OSPF            ProtocolType = 89
	PROTOCOL_SPRITE_RPC      ProtocolType = 90
	PROTOCOL_LARP            ProtocolType = 91
	PROTOCOL_MTP             ProtocolType = 92
	PROTOCOL_AX_25           ProtocolType = 93
	PROTOCOL_IPIP            ProtocolType = 94
	PROTOCOL_MICP            ProtocolType = 95
	PROTOCOL_SCC_SP          ProtocolType = 96
	PROTOCOL_ETHERIP         ProtocolType = 97
	PROTOCOL_ENCAP           ProtocolType = 98
	PROTOCOL_GMTP            ProtocolType = 100
	PROTOCOL_IFMP            ProtocolType = 101
	PROTOCOL_PNNI            ProtocolType = 102
	PROTOCOL_PIM             ProtocolType = 103
	PROTOCOL_ARIS            ProtocolType = 104
	PROTOCOL_SCPS            ProtocolType = 105
	PROTOCOL_QNX             ProtocolType = 106
	PROTOCOL_A_N             ProtocolType = 107
	PROTOCOL_IPCOMP          ProtocolType = 108
	PROTOCOL_SNP             ProtocolType = 109
	PROTOCOL_COMPAQ_PEER     ProtocolType = 110
	PROTOCOL_IPX_IN_IP       ProtocolType = 111
	PROTOCOL_CARP            ProtocolType = 112
	PROTOCOL_PGM             ProtocolType = 113
	PROTOCOL_L2TP            ProtocolType = 115
	PROTOCOL_DDX             ProtocolType = 116
	PROTOCOL_IATP            ProtocolType = 117
	PROTOCOL_STP             ProtocolType = 118
	PROTOCOL_SRP             ProtocolType = 119
	PROTOCOL_UTI             ProtocolType = 120
	PROTOCOL_SMP             ProtocolType = 121
	PROTOCOL_SM              ProtocolType = 122
	PROTOCOL_PTP             ProtocolType = 123
	PROTOCOL_ISIS            ProtocolType = 124
	PROTOCOL_FIRE            ProtocolType = 125
	PROTOCOL_CRTP            ProtocolType = 126
	PROTOCOL_CRUDP           ProtocolType = 127
	PROTOCOL_SSCOPMCE        ProtocolType = 128
	PROTOCOL_IPLT            ProtocolType = 129
	PROTOCOL_SPS             ProtocolType = 130
	PROTOCOL_PIPE            ProtocolType = 131
	PROTOCOL_SCTP            ProtocolType = 132
	PROTOCOL_FC              ProtocolType = 133
	PROTOCOL_RSVP_E2E_IGNORE ProtocolType = 134
	PROTOCOL_MOBILITY_HEADER ProtocolType = 135
	PROTOCOL_UDPLITE         ProtocolType = 136
	PROTOCOL_MPLS_IN_IP      ProtocolType = 137
	PROTOCOL_MANET           ProtocolType = 138
	PROTOCOL_HIP             ProtocolType = 139
	PROTOCOL_SHIM6           ProtocolType = 140
	PROTOCOL_WESP            ProtocolType = 141
	PROTOCOL_ROHC            ProtocolType = 142
	PROTOCOL_PFSYNC          ProtocolType = 240
)

var (
	protocolsMap = map[ProtocolType]string{
		PROTOCOL_INVALID:         "<INVALID>",
		PROTOCOL_IP:              "ip",
		PROTOCOL_ICMP:            "icmp",
		PROTOCOL_IGMP:            "igmp",
		PROTOCOL_GGP:             "ggp",
		PROTOCOL_IPENCAP:         "ipencap",
		PROTOCOL_ST2:             "st2",
		PROTOCOL_TCP:             "tcp",
		PROTOCOL_CBT:             "cbt",
		PROTOCOL_EGP:             "egp",
		PROTOCOL_IGP:             "igp",
		PROTOCOL_BBN_RCC:         "bbn-rcc",
		PROTOCOL_NVP:             "nvp",
		PROTOCOL_PUP:             "pup",
		PROTOCOL_ARGUS:           "argus",
		PROTOCOL_EMCON:           "emcon",
		PROTOCOL_XNET:            "xnet",
		PROTOCOL_CHAOS:           "chaos",
		PROTOCOL_UDP:             "udp",
		PROTOCOL_MUX:             "mux",
		PROTOCOL_DCN:             "dcn",
		PROTOCOL_HMP:             "hmp",
		PROTOCOL_PRM:             "prm",
		PROTOCOL_XNS_IDP:         "xns-idp",
		PROTOCOL_TRUNK_1:         "trunk-1",
		PROTOCOL_TRUNK_2:         "trunk-2",
		PROTOCOL_LEAF_1:          "leaf-1",
		PROTOCOL_LEAF_2:          "leaf-2",
		PROTOCOL_RDP:             "rdp",
		PROTOCOL_IRTP:            "irtp",
		PROTOCOL_ISO_TP4:         "iso-tp4",
		PROTOCOL_NETBLT:          "netblt",
		PROTOCOL_MFE_NSP:         "mfe-nsp",
		PROTOCOL_MERIT_INP:       "merit-inp",
		PROTOCOL_DCCP:            "dccp",
		PROTOCOL_THREE_PC:        "3pc",
		PROTOCOL_IDPR:            "idpr",
		PROTOCOL_XTP:             "xtp",
		PROTOCOL_DDP:             "ddp",
		PROTOCOL_IDPR_CMTP:       "idpr-cmtp",
		PROTOCOL_TPPLUSPLUS:      "tp++",
		PROTOCOL_IL:              "il",
		PROTOCOL_IPV6:            "ipv6",
		PROTOCOL_SDRP:            "sdrp",
		PROTOCOL_IPV6_ROUTE:      "ipv6-route",
		PROTOCOL_IPV6_FRAG:       "ipv6-frag",
		PROTOCOL_IDRP:            "idrp",
		PROTOCOL_RSVP:            "rsvp",
		PROTOCOL_GRE:             "gre",
		PROTOCOL_DSR:             "dsr",
		PROTOCOL_BNA:             "bna",
		PROTOCOL_ESP:             "esp",
		PROTOCOL_AH:              "ah",
		PROTOCOL_I_NLSP:          "i-nlsp",
		PROTOCOL_SWIPE:           "swipe",
		PROTOCOL_NARP:            "narp",
		PROTOCOL_MOBILE:          "mobile",
		PROTOCOL_TLSP:            "tlsp",
		PROTOCOL_SKIP:            "skip",
		PROTOCOL_IPV6_ICMP:       "ipv6-icmp",
		PROTOCOL_IPV6_NONXT:      "ipv6-nonxt",
		PROTOCOL_IPV6_OPTS:       "ipv6-opts",
		PROTOCOL_CFTP:            "cftp",
		PROTOCOL_SAT_EXPAK:       "sat-expak",
		PROTOCOL_KRYPTOLAN:       "kryptolan",
		PROTOCOL_RVD:             "rvd",
		PROTOCOL_IPPC:            "ippc",
		PROTOCOL_SAT_MON:         "sat-mon",
		PROTOCOL_VISA:            "visa",
		PROTOCOL_IPCV:            "ipcv",
		PROTOCOL_CPNX:            "cpnx",
		PROTOCOL_CPHB:            "cphb",
		PROTOCOL_WSN:             "wsn",
		PROTOCOL_PVP:             "pvp",
		PROTOCOL_BR_SAT_MON:      "br-sat-mon",
		PROTOCOL_SUN_ND:          "sun-nd",
		PROTOCOL_WB_MON:          "wb-mon",
		PROTOCOL_WB_EXPAK:        "wb-expak",
		PROTOCOL_ISO_IP:          "iso-ip",
		PROTOCOL_VMTP:            "vmtp",
		PROTOCOL_SECURE_VMTP:     "secure-vmtp",
		PROTOCOL_VINES:           "vines",
		PROTOCOL_TTP:             "ttp",
		PROTOCOL_NSFNET_IGP:      "nsfnet-igp",
		PROTOCOL_DGP:             "dgp",
		PROTOCOL_TCF:             "tcf",
		PROTOCOL_EIGRP:           "eigrp",
		PROTOCOL_OSPF:            "ospf",
		PROTOCOL_SPRITE_RPC:      "sprite-rpc",
		PROTOCOL_LARP:            "larp",
		PROTOCOL_MTP:             "mtp",
		PROTOCOL_AX_25:           "ax.25",
		PROTOCOL_IPIP:            "ipip",
		PROTOCOL_MICP:            "micp",
		PROTOCOL_SCC_SP:          "scc-sp",
		PROTOCOL_ETHERIP:         "etherip",
		PROTOCOL_ENCAP:           "encap",
		PROTOCOL_GMTP:            "gmtp",
		PROTOCOL_IFMP:            "ifmp",
		PROTOCOL_PNNI:            "pnni",
		PROTOCOL_PIM:             "pim",
		PROTOCOL_ARIS:            "aris",
		PROTOCOL_SCPS:            "scps",
		PROTOCOL_QNX:             "qnx",
		PROTOCOL_A_N:             "a/n",
		PROTOCOL_IPCOMP:          "ipcomp",
		PROTOCOL_SNP:             "snp",
		PROTOCOL_COMPAQ_PEER:     "compaq-peer",
		PROTOCOL_IPX_IN_IP:       "ipx-in-ip",
		PROTOCOL_CARP:            "carp",
		PROTOCOL_PGM:             "pgm",
		PROTOCOL_L2TP:            "l2tp",
		PROTOCOL_DDX:             "ddx",
		PROTOCOL_IATP:            "iatp",
		PROTOCOL_STP:             "stp",
		PROTOCOL_SRP:             "srp",
		PROTOCOL_UTI:             "uti",
		PROTOCOL_SMP:             "smp",
		PROTOCOL_SM:              "sm",
		PROTOCOL_PTP:             "ptp",
		PROTOCOL_ISIS:            "isis",
		PROTOCOL_FIRE:            "fire",
		PROTOCOL_CRTP:            "crtp",
		PROTOCOL_CRUDP:           "crudp",
		PROTOCOL_SSCOPMCE:        "sscopmce",
		PROTOCOL_IPLT:            "iplt",
		PROTOCOL_SPS:             "sps",
		PROTOCOL_PIPE:            "pipe",
		PROTOCOL_SCTP:            "sctp",
		PROTOCOL_FC:              "fc",
		PROTOCOL_RSVP_E2E_IGNORE: "rsvp-e2e-ignore",
		PROTOCOL_MOBILITY_HEADER: "mobility-header",
		PROTOCOL_UDPLITE:         "udplite",
		PROTOCOL_MPLS_IN_IP:      "mpls-in-ip",
		PROTOCOL_MANET:           "manet",
		PROTOCOL_HIP:             "hip",
		PROTOCOL_SHIM6:           "shim6",
		PROTOCOL_WESP:            "wesp",
		PROTOCOL_ROHC:            "rohc",
		PROTOCOL_PFSYNC:          "pfsync",
	}
)

type Protocol struct {
	value ProtocolType
}

func ProtocolNew() *Protocol {
	p := Protocol{value: PROTOCOL_INVALID}
	return &p
}

func (ctx *Protocol) TryInitFromString(param string) error {
	if len(param) > 0 {
		param = strings.ToLower(param)
		for key, value := range protocolsMap {
			if value == param {
				ctx.value = key
				return nil
			}
		}
	}
	ctx.value = PROTOCOL_INVALID
	return errors.New("cannot init from string")
}

func (ctx *Protocol) SetValue(protocol ProtocolType) {
	ctx.value = protocol
}

func (ctx *Protocol) GetValue() ProtocolType {
	return ctx.value
}

func (ctx Protocol) GetString() string {
	return protocolsMap[ctx.value]
}

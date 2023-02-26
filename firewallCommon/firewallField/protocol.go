package firewallField

import (
	"errors"
	"strings"
)

type ProtocolType uint64

// Extracted from: /etc/protocols
const (
	PROTOCOL_INVALID ProtocolType = 0xFFFFFFFFFFFFFFFF
	IP               ProtocolType = 0
	ICMP             ProtocolType = 1
	IGMP             ProtocolType = 2
	GGP              ProtocolType = 3
	IPENCAP          ProtocolType = 4
	ST2              ProtocolType = 5
	TCP              ProtocolType = 6
	CBT              ProtocolType = 7
	EGP              ProtocolType = 8
	IGP              ProtocolType = 9
	BBN_RCC          ProtocolType = 10
	NVP              ProtocolType = 11
	PUP              ProtocolType = 12
	ARGUS            ProtocolType = 13
	EMCON            ProtocolType = 14
	XNET             ProtocolType = 15
	CHAOS            ProtocolType = 16
	UDP              ProtocolType = 17
	MUX              ProtocolType = 18
	DCN              ProtocolType = 19
	HMP              ProtocolType = 20
	PRM              ProtocolType = 21
	XNS_IDP          ProtocolType = 22
	TRUNK_1          ProtocolType = 23
	TRUNK_2          ProtocolType = 24
	LEAF_1           ProtocolType = 25
	LEAF_2           ProtocolType = 26
	RDP              ProtocolType = 27
	IRTP             ProtocolType = 28
	ISO_TP4          ProtocolType = 29
	NETBLT           ProtocolType = 30
	MFE_NSP          ProtocolType = 31
	MERIT_INP        ProtocolType = 32
	DCCP             ProtocolType = 33
	THREE_PC         ProtocolType = 34
	IDPR             ProtocolType = 35
	XTP              ProtocolType = 36
	DDP              ProtocolType = 37
	IDPR_CMTP        ProtocolType = 38
	TPPLUSPLUS       ProtocolType = 39
	IL               ProtocolType = 40
	IPV6             ProtocolType = 41
	SDRP             ProtocolType = 42
	IPV6_ROUTE       ProtocolType = 43
	IPV6_FRAG        ProtocolType = 44
	IDRP             ProtocolType = 45
	RSVP             ProtocolType = 46
	GRE              ProtocolType = 47
	DSR              ProtocolType = 48
	BNA              ProtocolType = 49
	ESP              ProtocolType = 50
	AH               ProtocolType = 51
	I_NLSP           ProtocolType = 52
	SWIPE            ProtocolType = 53
	NARP             ProtocolType = 54
	MOBILE           ProtocolType = 55
	TLSP             ProtocolType = 56
	SKIP             ProtocolType = 57
	IPV6_ICMP        ProtocolType = 58
	IPV6_NONXT       ProtocolType = 59
	IPV6_OPTS        ProtocolType = 60
	CFTP             ProtocolType = 62
	SAT_EXPAK        ProtocolType = 64
	KRYPTOLAN        ProtocolType = 65
	RVD              ProtocolType = 66
	IPPC             ProtocolType = 67
	SAT_MON          ProtocolType = 69
	VISA             ProtocolType = 70
	IPCV             ProtocolType = 71
	CPNX             ProtocolType = 72
	CPHB             ProtocolType = 73
	WSN              ProtocolType = 74
	PVP              ProtocolType = 75
	BR_SAT_MON       ProtocolType = 76
	SUN_ND           ProtocolType = 77
	WB_MON           ProtocolType = 78
	WB_EXPAK         ProtocolType = 79
	ISO_IP           ProtocolType = 80
	VMTP             ProtocolType = 81
	SECURE_VMTP      ProtocolType = 82
	VINES            ProtocolType = 83
	TTP              ProtocolType = 84
	NSFNET_IGP       ProtocolType = 85
	DGP              ProtocolType = 86
	TCF              ProtocolType = 87
	EIGRP            ProtocolType = 88
	OSPF             ProtocolType = 89
	SPRITE_RPC       ProtocolType = 90
	LARP             ProtocolType = 91
	MTP              ProtocolType = 92
	AX_25            ProtocolType = 93
	IPIP             ProtocolType = 94
	MICP             ProtocolType = 95
	SCC_SP           ProtocolType = 96
	ETHERIP          ProtocolType = 97
	ENCAP            ProtocolType = 98
	GMTP             ProtocolType = 100
	IFMP             ProtocolType = 101
	PNNI             ProtocolType = 102
	PIM              ProtocolType = 103
	ARIS             ProtocolType = 104
	SCPS             ProtocolType = 105
	QNX              ProtocolType = 106
	A_N              ProtocolType = 107
	IPCOMP           ProtocolType = 108
	SNP              ProtocolType = 109
	COMPAQ_PEER      ProtocolType = 110
	IPX_IN_IP        ProtocolType = 111
	CARP             ProtocolType = 112
	PGM              ProtocolType = 113
	L2TP             ProtocolType = 115
	DDX              ProtocolType = 116
	IATP             ProtocolType = 117
	STP              ProtocolType = 118
	SRP              ProtocolType = 119
	UTI              ProtocolType = 120
	SMP              ProtocolType = 121
	SM               ProtocolType = 122
	PTP              ProtocolType = 123
	ISIS             ProtocolType = 124
	FIRE             ProtocolType = 125
	CRTP             ProtocolType = 126
	CRUDP            ProtocolType = 127
	SSCOPMCE         ProtocolType = 128
	IPLT             ProtocolType = 129
	SPS              ProtocolType = 130
	PIPE             ProtocolType = 131
	SCTP             ProtocolType = 132
	FC               ProtocolType = 133
	RSVP_E2E_IGNORE  ProtocolType = 134
	MOBILITY_HEADER  ProtocolType = 135
	UDPLITE          ProtocolType = 136
	MPLS_IN_IP       ProtocolType = 137
	MANET            ProtocolType = 138
	HIP              ProtocolType = 139
	SHIM6            ProtocolType = 140
	WESP             ProtocolType = 141
	ROHC             ProtocolType = 142
	PFSYNC           ProtocolType = 240
	DIVERT           ProtocolType = 258
)

var (
	protocolsMap = map[ProtocolType]string{
		PROTOCOL_INVALID: "<INVALID>",
		IP:               "ip",
		ICMP:             "icmp",
		IGMP:             "igmp",
		GGP:              "ggp",
		IPENCAP:          "ipencap",
		ST2:              "st2",
		TCP:              "tcp",
		CBT:              "cbt",
		EGP:              "egp",
		IGP:              "igp",
		BBN_RCC:          "bbn-rcc",
		NVP:              "nvp",
		PUP:              "pup",
		ARGUS:            "argus",
		EMCON:            "emcon",
		XNET:             "xnet",
		CHAOS:            "chaos",
		UDP:              "udp",
		MUX:              "mux",
		DCN:              "dcn",
		HMP:              "hmp",
		PRM:              "prm",
		XNS_IDP:          "xns-idp",
		TRUNK_1:          "trunk-1",
		TRUNK_2:          "trunk-2",
		LEAF_1:           "leaf-1",
		LEAF_2:           "leaf-2",
		RDP:              "rdp",
		IRTP:             "irtp",
		ISO_TP4:          "iso-tp4",
		NETBLT:           "netblt",
		MFE_NSP:          "mfe-nsp",
		MERIT_INP:        "merit-inp",
		DCCP:             "dccp",
		THREE_PC:         "3pc",
		IDPR:             "idpr",
		XTP:              "xtp",
		DDP:              "ddp",
		IDPR_CMTP:        "idpr-cmtp",
		TPPLUSPLUS:       "tp++",
		IL:               "il",
		IPV6:             "ipv6",
		SDRP:             "sdrp",
		IPV6_ROUTE:       "ipv6-route",
		IPV6_FRAG:        "ipv6-frag",
		IDRP:             "idrp",
		RSVP:             "rsvp",
		GRE:              "gre",
		DSR:              "dsr",
		BNA:              "bna",
		ESP:              "esp",
		AH:               "ah",
		I_NLSP:           "i-nlsp",
		SWIPE:            "swipe",
		NARP:             "narp",
		MOBILE:           "mobile",
		TLSP:             "tlsp",
		SKIP:             "skip",
		IPV6_ICMP:        "ipv6-icmp",
		IPV6_NONXT:       "ipv6-nonxt",
		IPV6_OPTS:        "ipv6-opts",
		CFTP:             "cftp",
		SAT_EXPAK:        "sat-expak",
		KRYPTOLAN:        "kryptolan",
		RVD:              "rvd",
		IPPC:             "ippc",
		SAT_MON:          "sat-mon",
		VISA:             "visa",
		IPCV:             "ipcv",
		CPNX:             "cpnx",
		CPHB:             "cphb",
		WSN:              "wsn",
		PVP:              "pvp",
		BR_SAT_MON:       "br-sat-mon",
		SUN_ND:           "sun-nd",
		WB_MON:           "wb-mon",
		WB_EXPAK:         "wb-expak",
		ISO_IP:           "iso-ip",
		VMTP:             "vmtp",
		SECURE_VMTP:      "secure-vmtp",
		VINES:            "vines",
		TTP:              "ttp",
		NSFNET_IGP:       "nsfnet-igp",
		DGP:              "dgp",
		TCF:              "tcf",
		EIGRP:            "eigrp",
		OSPF:             "ospf",
		SPRITE_RPC:       "sprite-rpc",
		LARP:             "larp",
		MTP:              "mtp",
		AX_25:            "ax.25",
		IPIP:             "ipip",
		MICP:             "micp",
		SCC_SP:           "scc-sp",
		ETHERIP:          "etherip",
		ENCAP:            "encap",
		GMTP:             "gmtp",
		IFMP:             "ifmp",
		PNNI:             "pnni",
		PIM:              "pim",
		ARIS:             "aris",
		SCPS:             "scps",
		QNX:              "qnx",
		A_N:              "a/n",
		IPCOMP:           "ipcomp",
		SNP:              "snp",
		COMPAQ_PEER:      "compaq-peer",
		IPX_IN_IP:        "ipx-in-ip",
		CARP:             "carp",
		PGM:              "pgm",
		L2TP:             "l2tp",
		DDX:              "ddx",
		IATP:             "iatp",
		STP:              "stp",
		SRP:              "srp",
		UTI:              "uti",
		SMP:              "smp",
		SM:               "sm",
		PTP:              "ptp",
		ISIS:             "isis",
		FIRE:             "fire",
		CRTP:             "crtp",
		CRUDP:            "crudp",
		SSCOPMCE:         "sscopmce",
		IPLT:             "iplt",
		SPS:              "sps",
		PIPE:             "pipe",
		SCTP:             "sctp",
		FC:               "fc",
		RSVP_E2E_IGNORE:  "rsvp-e2e-ignore",
		MOBILITY_HEADER:  "mobility-header",
		UDPLITE:          "udplite",
		MPLS_IN_IP:       "mpls-in-ip",
		MANET:            "manet",
		HIP:              "hip",
		SHIM6:            "shim6",
		WESP:             "wesp",
		ROHC:             "rohc",
		PFSYNC:           "pfsync",
		DIVERT:           "divert",
	}
)

type Protocol struct {
	value ProtocolType
}

func (ctx *Protocol) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var protocolStr string = ""
	err := unmarshal(&protocolStr)
	if err != nil {
		return err
	}
	protocol, err := ProtocolTypeFromString(protocolStr)
	if err != nil {
		return err
	}
	ctx.SetValue(protocol.GetValue())
	return nil
}

func (ctx *Protocol) TryInitFromString(param string) error {
	param = strings.ToLower(param)
	for key, value := range protocolsMap {
		if value == param {
			ctx.value = key
			return nil
		}
	}
	ctx.value = PROTOCOL_INVALID
	return errors.New("Cannot init from string")
}

func (ctx *Protocol) TryInitFromRest(param string) error {
	return ctx.TryInitFromString(param)
}

func (ctx *Protocol) TryInitFromIpTables(param string) error {
	return ctx.TryInitFromString(param)
}

func ProtocolTypeFromString(protocolString string) (Protocol, error) {
	protocolString = strings.ToLower(protocolString)
	for key, value := range protocolsMap {
		if value == protocolString {
			return Protocol{value: key}, nil
		}
	}

	return Protocol{value: PROTOCOL_INVALID}, errors.New("Invalid protocol text name")
}

func ProtocolTypeFromValue(value ProtocolType) (Protocol, error) {
	_, ok := protocolsMap[value]
	if ok {
		return Protocol{value: value}, nil
	}
	return Protocol{value: PROTOCOL_INVALID}, errors.New("Invalid protocol type value")
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

func (ctx Protocol) MarshalRest() string {
	return protocolsMap[ctx.value]
}

func (ctx Protocol) MarshalIpTables() string {
	return protocolsMap[ctx.value]
}

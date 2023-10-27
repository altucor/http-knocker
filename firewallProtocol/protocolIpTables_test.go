package firewallProtocol

import (
	"net/netip"
	"testing"

	"github.com/altucor/http-knocker/firewallCommon"
	"github.com/altucor/http-knocker/firewallCommon/firewallField"
)

func Test_IpTablesRule_fromProtocol(t *testing.T) {
	{
		input := "iptables -I INPUT 2 -p tcp -s 10.1.1.2 --dport 22 -j ACCEPT -m comment --comment \"My comments here\""
		expected := firewallCommon.FirewallRule{
			Id:          firewallField.NumberNew(0),
			Action:      firewallField.ActionNew(firewallField.ACTION_ACCEPT),
			Chain:       firewallField.ChainNew(firewallField.CHAIN_INPUT),
			Disabled:    firewallField.BoolNew(false),
			Protocol:    firewallField.ProtocolNew(firewallField.PROTOCOL_TCP),
			SrcAddress:  firewallField.AddressNew(netip.AddrFrom4([4]byte{10, 1, 1, 2})),
			DstPort:     firewallField.PortNew(22),
			Comment:     firewallField.TextNew("\"My comments here\""),
			PlaceBefore: firewallField.NumberNew(0),
		}

		ipTablesRule := IpTablesRule{}
		result, err := ipTablesRule.fromProtocol(input)
		if err != nil {
			t.Errorf("\"Test_FirewallField_Action_TryInitFromString('%s')\" FAILED, Error -> %v", input, err)
		}

		if result == expected {
			t.Logf("\"Test_FirewallField_Action_TryInitFromString('%s')\" SUCCEDED, expected -> %v, got -> %v", input, expected, result)
		} else {
			t.Errorf("\"Test_FirewallField_Action_TryInitFromString('%s')\" FAILED\nexpected -> %v\ngot -> %v", input, expected, result)
		}
	}
}

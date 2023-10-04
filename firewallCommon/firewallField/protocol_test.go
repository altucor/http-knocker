package firewallField

import "testing"

func Test_FirewallField_Protocol_TryInitFromString(t *testing.T) {
	{
		input := "TCP"
		expected := PROTOCOL_TCP

		prot := Protocol{}
		prot.TryInitFromString(input)
		result := prot.GetValue()

		if result == expected {
			t.Logf("\"Test_FirewallField_Protocol_TryInitFromString('%s')\" SUCCEDED, expected -> %v, got -> %v", input, expected, result)
		} else {
			t.Errorf("\"Test_FirewallField_Protocol_TryInitFromString('%s')\" FAILED, expected -> %v, got -> %v", input, expected, result)
		}
	}
	{
		input := "UdP"
		expected := PROTOCOL_UDP

		prot := Protocol{}
		prot.TryInitFromString(input)
		result := prot.GetValue()

		if result == expected {
			t.Logf("\"Test_FirewallField_Protocol_TryInitFromString('%s')\" SUCCEDED, expected -> %v, got -> %v", input, expected, result)
		} else {
			t.Errorf("\"Test_FirewallField_Protocol_TryInitFromString('%s')\" FAILED, expected -> %v, got -> %v", input, expected, result)
		}
	}
	{
		input := "TWP"
		expected := PROTOCOL_INVALID

		prot := Protocol{}
		prot.TryInitFromString(input)
		result := prot.GetValue()

		if result == expected {
			t.Logf("\"Test_FirewallField_Protocol_TryInitFromString('%s')\" SUCCEDED, expected -> %v, got -> %v", input, expected, result)
		} else {
			t.Errorf("\"Test_FirewallField_Protocol_TryInitFromString('%s')\" FAILED, expected -> %v, got -> %v", input, expected, result)
		}
	}
	{
		input := ""
		expected := PROTOCOL_INVALID

		prot := Protocol{}
		prot.TryInitFromString(input)
		result := prot.GetValue()

		if result == expected {
			t.Logf("\"Test_FirewallField_Protocol_TryInitFromString('%s')\" SUCCEDED, expected -> %v, got -> %v", input, expected, result)
		} else {
			t.Errorf("\"Test_FirewallField_Protocol_TryInitFromString('%s')\" FAILED, expected -> %v, got -> %v", input, expected, result)
		}
	}
}

func Test_FirewallField_Protocol_SetGetValueString(t *testing.T) {
	{
		input := PROTOCOL_ICMP
		expected := PROTOCOL_ICMP

		prot := ProtocolNew()
		prot.SetValue(input)
		result := prot.GetValue()

		if result == expected {
			t.Logf("\"Test_FirewallField_Protocol_SetGetValueString('%v')\" SUCCEDED, expected -> %v, got -> %v", input, expected, result)
		} else {
			t.Errorf("\"Test_FirewallField_Protocol_SetGetValueString('%v')\" FAILED, expected -> %v, got -> %v", input, expected, result)
		}
	}
	{
		input := ""
		expected := PROTOCOL_INVALID

		prot := ProtocolNew()
		result := prot.GetValue()

		if result == expected {
			t.Logf("\"Test_FirewallField_Protocol_SetGetValueString('%v')\" SUCCEDED, expected -> %v, got -> %v", input, expected, result)
		} else {
			t.Errorf("\"Test_FirewallField_Protocol_SetGetValueString('%v')\" FAILED, expected -> %v, got -> %v", input, expected, result)
		}
	}
}

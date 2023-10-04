package firewallField

import "testing"

func Test_FirewallField_Chain_FromString(t *testing.T) {
	{
		input := "FORWARD"
		expected := CHAIN_FORWARD

		chain := Chain{}
		chain.TryInitFromString(input)
		result := chain.GetValue()

		if result == expected {
			t.Logf("\"Test_FirewallField_Chain_FromString('%v')\" SUCCEDED, expected -> %v, got -> %v", input, expected, result)
		} else {
			t.Errorf("\"Test_FirewallField_Chain_FromString('%v')\" FAILED, expected -> %v, got -> %v", input, expected, result)
		}
	}
	{
		input := "FOrWard"
		expected := CHAIN_FORWARD

		chain := Chain{}
		chain.TryInitFromString(input)
		result := chain.GetValue()

		if result == expected {
			t.Logf("\"Test_FirewallField_Chain_FromString('%v')\" SUCCEDED, expected -> %v, got -> %v", input, expected, result)
		} else {
			t.Errorf("\"Test_FirewallField_Chain_FromString('%v')\" FAILED, expected -> %v, got -> %v", input, expected, result)
		}
	}
	{
		input := "forward"
		expected := CHAIN_FORWARD

		chain := Chain{}
		chain.TryInitFromString(input)
		result := chain.GetValue()

		if result == expected {
			t.Logf("\"Test_FirewallField_Chain_FromString('%v')\" SUCCEDED, expected -> %v, got -> %v", input, expected, result)
		} else {
			t.Errorf("\"Test_FirewallField_Chain_FromString('%v')\" FAILED, expected -> %v, got -> %v", input, expected, result)
		}
	}
	{
		input := ""
		expected := CHAIN_INVALID

		chain := Chain{}
		chain.TryInitFromString(input)
		result := chain.GetValue()

		if result == expected {
			t.Logf("\"Test_FirewallField_Chain_FromString('%v')\" SUCCEDED, expected -> %v, got -> %v", input, expected, result)
		} else {
			t.Errorf("\"Test_FirewallField_Chain_FromString('%v')\" FAILED, expected -> %v, got -> %v", input, expected, result)
		}
	}
	{
		input := "f0rw4ard"
		expected := CHAIN_INVALID

		chain := Chain{}
		chain.TryInitFromString(input)
		result := chain.GetValue()

		if result == expected {
			t.Logf("\"Test_FirewallField_Chain_FromString('%v')\" SUCCEDED, expected -> %v, got -> %v", input, expected, result)
		} else {
			t.Errorf("\"Test_FirewallField_Chain_FromString('%v')\" FAILED, expected -> %v, got -> %v", input, expected, result)
		}
	}
}

func Test_FirewallField_Chain_SetGetValueString(t *testing.T) {
	{
		input := CHAIN_INPUT
		expected := CHAIN_INPUT

		chain := Chain{}
		chain.SetValue(input)
		result := chain.GetValue()

		if result == expected {
			t.Logf("\"Test_FirewallField_Chain_SetGetValueString('%v')\" SUCCEDED, expected -> %v, got -> %v", input, expected, result)
		} else {
			t.Errorf("\"Test_FirewallField_Chain_SetGetValueString('%v')\" FAILED, expected -> %v, got -> %v", input, expected, result)
		}
	}
	{
		input := CHAIN_INVALID
		expected := CHAIN_INVALID

		chain := Chain{}
		result := chain.GetValue()

		if result == expected {
			t.Logf("\"Test_FirewallField_Chain_SetGetValueString('%v')\" SUCCEDED, expected -> %v, got -> %v", input, expected, result)
		} else {
			t.Errorf("\"Test_FirewallField_Chain_SetGetValueString('%v')\" FAILED, expected -> %v, got -> %v", input, expected, result)
		}
	}
	{
		input := CHAIN_INPUT
		expected := chainMap[CHAIN_INPUT]

		chain := Chain{}
		chain.SetValue(input)
		result := chain.GetString()

		if result == expected {
			t.Logf("\"Test_FirewallField_Chain_SetGetValueString('%v')\" SUCCEDED, expected -> %v, got -> %v", input, expected, result)
		} else {
			t.Errorf("\"Test_FirewallField_Chain_SetGetValueString('%v')\" FAILED, expected -> %v, got -> %v", input, expected, result)
		}
	}
}

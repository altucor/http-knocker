package firewallField

import "testing"

func Test_FirewallField_Port_TryInitFromString(t *testing.T) {
	{
		input := "2354"
		var expected uint16 = 2354

		p := Port{}
		p.TryInitFromString(input)
		result := p.GetValue()

		if result == expected {
			t.Logf("\"Test_FirewallField_Port_TryInitFromString('%s')\" SUCCEDED, expected -> %v, got -> %v", input, expected, result)
		} else {
			t.Errorf("\"Test_FirewallField_Port_TryInitFromString('%s')\" FAILED, expected -> %v, got -> %v", input, expected, result)
		}
	}
	{
		input := "-2354"
		var expected uint16 = 0

		p := Port{}
		p.TryInitFromString(input)
		result := p.GetValue()

		if result == expected {
			t.Logf("\"Test_FirewallField_Port_TryInitFromString('%s')\" SUCCEDED, expected -> %v, got -> %v", input, expected, result)
		} else {
			t.Errorf("\"Test_FirewallField_Port_TryInitFromString('%s')\" FAILED, expected -> %v, got -> %v", input, expected, result)
		}
	}
	{
		input := "70000"
		var expected uint16 = 0

		p := Port{}
		p.TryInitFromString(input)
		result := p.GetValue()

		if result == expected {
			t.Logf("\"Test_FirewallField_Port_TryInitFromString('%s')\" SUCCEDED, expected -> %v, got -> %v", input, expected, result)
		} else {
			t.Errorf("\"Test_FirewallField_Port_TryInitFromString('%s')\" FAILED, expected -> %v, got -> %v", input, expected, result)
		}
	}
	{
		input := "0"
		var expected uint16 = 0

		p := Port{}
		p.TryInitFromString(input)
		result := p.GetValue()

		if result == expected {
			t.Logf("\"Test_FirewallField_Port_TryInitFromString('%s')\" SUCCEDED, expected -> %v, got -> %v", input, expected, result)
		} else {
			t.Errorf("\"Test_FirewallField_Port_TryInitFromString('%s')\" FAILED, expected -> %v, got -> %v", input, expected, result)
		}
	}
	{
		input := ""
		var expected uint16 = 0

		p := Port{}
		p.TryInitFromString(input)
		result := p.GetValue()

		if result == expected {
			t.Logf("\"Test_FirewallField_Port_TryInitFromString('%s')\" SUCCEDED, expected -> %v, got -> %v", input, expected, result)
		} else {
			t.Errorf("\"Test_FirewallField_Port_TryInitFromString('%s')\" FAILED, expected -> %v, got -> %v", input, expected, result)
		}
	}
	{
		input := "-10"
		var expected uint16 = 0

		p := Port{}
		p.TryInitFromString(input)
		result := p.GetValue()

		if result == expected {
			t.Logf("\"Test_FirewallField_Port_TryInitFromString('%s')\" SUCCEDED, expected -> %v, got -> %v", input, expected, result)
		} else {
			t.Errorf("\"Test_FirewallField_Port_TryInitFromString('%s')\" FAILED, expected -> %v, got -> %v", input, expected, result)
		}
	}
}

func Test_FirewallField_Port_SetGetValueString(t *testing.T) {
	{
		var input uint16 = 6643
		var expected uint16 = 6643

		p := Port{}
		p.SetValue(input)
		result := p.GetValue()

		if result == expected {
			t.Logf("\"Test_FirewallField_Port_SetGetValueString('%d')\" SUCCEDED, expected -> %v, got -> %v", input, expected, result)
		} else {
			t.Errorf("\"Test_FirewallField_Port_SetGetValueString('%d')\" FAILED, expected -> %v, got -> %v", input, expected, result)
		}
	}
	{
		var input uint16 = 6643
		expected := "6643"

		p := Port{}
		p.SetValue(input)
		result := p.GetString()

		if result == expected {
			t.Logf("\"Test_FirewallField_Port_SetGetValueString('%d')\" SUCCEDED, expected -> %v, got -> %v", input, expected, result)
		} else {
			t.Errorf("\"Test_FirewallField_Port_SetGetValueString('%d')\" FAILED, expected -> %v, got -> %v", input, expected, result)
		}
	}
	{
		var input uint16 = 0
		expected := "0"

		p := Port{}
		p.SetValue(input)
		result := p.GetString()

		if result == expected {
			t.Logf("\"Test_FirewallField_Port_SetGetValueString('%d')\" SUCCEDED, expected -> %v, got -> %v", input, expected, result)
		} else {
			t.Errorf("\"Test_FirewallField_Port_SetGetValueString('%d')\" FAILED, expected -> %v, got -> %v", input, expected, result)
		}
	}
}

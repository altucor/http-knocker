package firewallField

import "testing"

func Test_FirewallField_Number_TryInitFromString(t *testing.T) {
	{
		input := "123"
		var expected uint64 = 0x123

		n := Number{}
		n.TryInitFromString(input)
		result := n.GetValue()

		if result == expected {
			t.Logf("\"Test_FirewallField_Number_TryInitFromString('%s')\" SUCCEDED, expected -> %v, got -> %v", input, expected, result)
		} else {
			t.Errorf("\"Test_FirewallField_Number_TryInitFromString('%s')\" FAILED, expected -> %v, got -> %v", input, expected, result)
		}
	}
	{
		input := "0x123"
		var expected uint64 = 0x123

		n := Number{}
		n.TryInitFromString(input)
		result := n.GetValue()

		if result == expected {
			t.Logf("\"Test_FirewallField_Number_TryInitFromString('%s')\" SUCCEDED, expected -> %v, got -> %v", input, expected, result)
		} else {
			t.Errorf("\"Test_FirewallField_Number_TryInitFromString('%s')\" FAILED, expected -> %v, got -> %v", input, expected, result)
		}
	}
	{
		input := "-1"
		var expected uint64 = 0xFFFFFFFFFFFFFFFF

		n := Number{}
		n.TryInitFromString(input)
		result := n.GetValue()

		if result == expected {
			t.Logf("\"Test_FirewallField_Number_TryInitFromString('%s')\" SUCCEDED, expected -> %v, got -> %v", input, expected, result)
		} else {
			t.Errorf("\"Test_FirewallField_Number_TryInitFromString('%s')\" FAILED, expected -> %v, got -> %v", input, expected, result)
		}
	}
	{
		input := ""
		var expected uint64 = 0xFFFFFFFFFFFFFFFF

		n := Number{}
		n.TryInitFromString(input)
		result := n.GetValue()

		if result == expected {
			t.Logf("\"Test_FirewallField_Number_TryInitFromString('%s')\" SUCCEDED, expected -> %v, got -> %v", input, expected, result)
		} else {
			t.Errorf("\"Test_FirewallField_Number_TryInitFromString('%s')\" FAILED, expected -> %v, got -> %v", input, expected, result)
		}
	}
	{
		input := "0"
		var expected uint64 = 0x00

		n := Number{}
		n.TryInitFromString(input)
		result := n.GetValue()

		if result == expected {
			t.Logf("\"Test_FirewallField_Number_TryInitFromString('%s')\" SUCCEDED, expected -> %v, got -> %v", input, expected, result)
		} else {
			t.Errorf("\"Test_FirewallField_Number_TryInitFromString('%s')\" FAILED, expected -> %v, got -> %v", input, expected, result)
		}
	}
	{
		input := "0sssafsga"
		var expected uint64 = 0xFFFFFFFFFFFFFFFF

		n := Number{}
		n.TryInitFromString(input)
		result := n.GetValue()

		if result == expected {
			t.Logf("\"Test_FirewallField_Number_TryInitFromString('%s')\" SUCCEDED, expected -> %v, got -> %v", input, expected, result)
		} else {
			t.Errorf("\"Test_FirewallField_Number_TryInitFromString('%s')\" FAILED, expected -> %v, got -> %v", input, expected, result)
		}
	}
}

func Test_FirewallField_Number_SetGetValueString(t *testing.T) {
	{
		var input uint64 = 456
		var expected uint64 = 456

		n := Number{}
		n.SetValue(input)
		result := n.GetValue()

		if result == expected {
			t.Logf("\"Test_FirewallField_Number_SetGetValueString('%d')\" SUCCEDED, expected -> %v, got -> %v", input, expected, result)
		} else {
			t.Errorf("\"Test_FirewallField_Number_SetGetValueString('%d')\" FAILED, expected -> %v, got -> %v", input, expected, result)
		}
	}
	{
		var input uint64 = 0x456
		expected := "0x456"

		n := Number{}
		n.SetValue(input)
		result := n.GetString()

		if result == expected {
			t.Logf("\"Test_FirewallField_Number_SetGetValueString('%d')\" SUCCEDED, expected -> %v, got -> %v", input, expected, result)
		} else {
			t.Errorf("\"Test_FirewallField_Number_SetGetValueString('%d')\" FAILED, expected -> %v, got -> %v", input, expected, result)
		}
	}
	{
		var input uint64 = 0x00
		expected := "0x00"

		n := Number{}
		n.SetValue(input)
		result := n.GetString()

		if result == expected {
			t.Logf("\"Test_FirewallField_Number_SetGetValueString('%d')\" SUCCEDED, expected -> %v, got -> %v", input, expected, result)
		} else {
			t.Errorf("\"Test_FirewallField_Number_SetGetValueString('%d')\" FAILED, expected -> %v, got -> %v", input, expected, result)
		}
	}
	{
		var input uint64 = 18446603336221161335
		expected := "0xFFFF7FFFFFFF7777"

		n := Number{}
		n.SetValue(input)
		result := n.GetString()

		if result == expected {
			t.Logf("\"Test_FirewallField_Number_SetGetValueString('%d')\" SUCCEDED, expected -> %v, got -> %v", input, expected, result)
		} else {
			t.Errorf("\"Test_FirewallField_Number_SetGetValueString('%d')\" FAILED, expected -> %v, got -> %v", input, expected, result)
		}
	}
}

package firewallField

import "testing"

func Test_FirewallField_Bool_TryInitFromString(t *testing.T) {
	{
		input := "True"
		expected := true

		b := Bool{}
		b.TryInitFromString(input)
		result := b.GetValue()

		if result == expected {
			t.Logf("\"Test_FirewallField_Bool_TryInitFromString('%s')\" SUCCEDED, expected -> %v, got -> %v", input, expected, result)
		} else {
			t.Errorf("\"Test_FirewallField_Bool_TryInitFromString('%s')\" FAILED, expected -> %v, got -> %v", input, expected, result)
		}
	}
	{
		input := "fAlse"
		expected := false

		b := Bool{}
		b.TryInitFromString(input)
		result := b.GetValue()

		if result == expected {
			t.Logf("\"Test_FirewallField_Bool_TryInitFromString('%s')\" SUCCEDED, expected -> %v, got -> %v", input, expected, result)
		} else {
			t.Errorf("\"Test_FirewallField_Bool_TryInitFromString('%s')\" FAILED, expected -> %v, got -> %v", input, expected, result)
		}
	}
	{
		input := ""
		expected := false

		b := Bool{}
		b.TryInitFromString(input)
		result := b.GetValue()

		if result == expected {
			t.Logf("\"Test_FirewallField_Bool_TryInitFromString('%s')\" SUCCEDED, expected -> %v, got -> %v", input, expected, result)
		} else {
			t.Errorf("\"Test_FirewallField_Bool_TryInitFromString('%s')\" FAILED, expected -> %v, got -> %v", input, expected, result)
		}
	}
	{
		input := "TRUE"
		expected := true

		b := Bool{}
		b.TryInitFromString(input)
		result := b.GetValue()

		if result == expected {
			t.Logf("\"Test_FirewallField_Bool_TryInitFromString('%s')\" SUCCEDED, expected -> %v, got -> %v", input, expected, result)
		} else {
			t.Errorf("\"Test_FirewallField_Bool_TryInitFromString('%s')\" FAILED, expected -> %v, got -> %v", input, expected, result)
		}
	}
	{
		input := "false"
		expected := false

		b := Bool{}
		b.TryInitFromString(input)
		result := b.GetValue()

		if result == expected {
			t.Logf("\"Test_FirewallField_Bool_TryInitFromString('%s')\" SUCCEDED, expected -> %v, got -> %v", input, expected, result)
		} else {
			t.Errorf("\"Test_FirewallField_Bool_TryInitFromString('%s')\" FAILED, expected -> %v, got -> %v", input, expected, result)
		}
	}
	{
		input := "looool"
		expected := false

		b := Bool{}
		b.TryInitFromString(input)
		result := b.GetValue()

		if result == expected {
			t.Logf("\"Test_FirewallField_Bool_TryInitFromString('%s')\" SUCCEDED, expected -> %v, got -> %v", input, expected, result)
		} else {
			t.Errorf("\"Test_FirewallField_Bool_TryInitFromString('%s')\" FAILED, expected -> %v, got -> %v", input, expected, result)
		}
	}
}

func Test_FirewallField_Bool_SetGetValueString(t *testing.T) {
	{
		input := true
		expected := true

		b := Bool{}
		b.SetValue(input)
		result := b.GetValue()

		if result == expected {
			t.Logf("\"Test_FirewallField_Bool_TryInitFromString('%v')\" SUCCEDED, expected -> %v, got -> %v", input, expected, result)
		} else {
			t.Errorf("\"Test_FirewallField_Bool_TryInitFromString('%v')\" FAILED, expected -> %v, got -> %v", input, expected, result)
		}
	}
	{
		input := true
		expected := "true"

		b := Bool{}
		b.SetValue(input)
		result := b.GetString()

		if result == expected {
			t.Logf("\"Test_FirewallField_Bool_TryInitFromString('%v')\" SUCCEDED, expected -> %v, got -> %v", input, expected, result)
		} else {
			t.Errorf("\"Test_FirewallField_Bool_TryInitFromString('%v')\" FAILED, expected -> %v, got -> %v", input, expected, result)
		}
	}
}

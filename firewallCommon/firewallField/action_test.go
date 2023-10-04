package firewallField

import "testing"

func Test_FirewallField_Action_TryInitFromString(t *testing.T) {
	{
		input := "ACCEPT"
		expected := ACTION_ACCEPT

		act := Action{}
		act.TryInitFromString(input)
		result := act.GetValue()

		if result == expected {
			t.Logf("\"Test_FirewallField_Action_TryInitFromString('%s')\" SUCCEDED, expected -> %v, got -> %v", input, expected, result)
		} else {
			t.Errorf("\"Test_FirewallField_Action_TryInitFromString('%s')\" FAILED, expected -> %v, got -> %v", input, expected, result)
		}
	}
	{
		input := "dRoP"
		expected := ACTION_DROP

		act := Action{}
		act.TryInitFromString(input)
		result := act.GetValue()

		if result == expected {
			t.Logf("\"Test_FirewallField_Action_TryInitFromString('%s')\" SUCCEDED, expected -> %v, got -> %v", input, expected, result)
		} else {
			t.Errorf("\"Test_FirewallField_Action_TryInitFromString('%s')\" FAILED, expected -> %v, got -> %v", input, expected, result)
		}
	}
	{
		input := "jump"
		expected := ACTION_JUMP

		act := Action{}
		act.TryInitFromString(input)
		result := act.GetValue()

		if result == expected {
			t.Logf("\"Test_FirewallField_Action_TryInitFromString('%s')\" SUCCEDED, expected -> %v, got -> %v", input, expected, result)
		} else {
			t.Errorf("\"Test_FirewallField_Action_TryInitFromString('%s')\" FAILED, expected -> %v, got -> %v", input, expected, result)
		}
	}
	{
		input := "reMect"
		expected := ACTION_INVALID

		act := Action{}
		act.TryInitFromString(input)
		result := act.GetValue()

		if result == expected {
			t.Logf("\"Test_FirewallField_Action_TryInitFromString('%s')\" SUCCEDED, expected -> %v, got -> %v", input, expected, result)
		} else {
			t.Errorf("\"Test_FirewallField_Action_TryInitFromString('%s')\" FAILED, expected -> %v, got -> %v", input, expected, result)
		}
	}
	{
		input := ""
		expected := ACTION_INVALID

		act := Action{}
		act.TryInitFromString(input)
		result := act.GetValue()

		if result == expected {
			t.Logf("\"Test_FirewallField_Action_TryInitFromString('%s')\" SUCCEDED, expected -> %v, got -> %v", input, expected, result)
		} else {
			t.Errorf("\"Test_FirewallField_Action_TryInitFromString('%s')\" FAILED, expected -> %v, got -> %v", input, expected, result)
		}
	}
}

func Test_FirewallField_Action_SetGetValueString(t *testing.T) {
	{
		input := ACTION_DROP
		expected := ACTION_DROP

		act := Action{}
		act.SetValue(input)
		result := act.GetValue()

		if result == expected {
			t.Logf("\"Test_FirewallField_Action_SetGetValueString('%v')\" SUCCEDED, expected -> %v, got -> %v", input, expected, result)
		} else {
			t.Errorf("\"Test_FirewallField_Action_SetGetValueString('%v')\" FAILED, expected -> %v, got -> %v", input, expected, result)
		}
	}
	{
		input := ""
		expected := ACTION_INVALID

		act := Action{}
		result := act.GetValue()

		if result == expected {
			t.Logf("\"Test_FirewallField_Action_SetGetValueString('%v')\" SUCCEDED, expected -> %v, got -> %v", input, expected, result)
		} else {
			t.Errorf("\"Test_FirewallField_Action_SetGetValueString('%v')\" FAILED, expected -> %v, got -> %v", input, expected, result)
		}
	}
	{
		input := ACTION_REJECT
		expected := actionMap[ACTION_REJECT]

		act := Action{}
		act.SetValue(input)
		result := act.GetString()

		if result == expected {
			t.Logf("\"Test_FirewallField_Action_SetGetValueString('%v')\" SUCCEDED, expected -> %v, got -> %v", input, expected, result)
		} else {
			t.Errorf("\"Test_FirewallField_Action_SetGetValueString('%v')\" FAILED, expected -> %v, got -> %v", input, expected, result)
		}
	}

}

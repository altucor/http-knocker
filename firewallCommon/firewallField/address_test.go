package firewallField

import (
	"net/netip"
	"testing"
)

func Test_FirewallField_Address_TryInitFromString(t *testing.T) {
	{
		input := "1.1.1.1/32"
		expected := netip.AddrFrom4([4]byte{1, 1, 1, 1})

		addr := Address{}
		addr.TryInitFromString(input)
		result := addr.GetValue()

		if result == expected {
			t.Logf("\"Test_FirewallField_Address_TryInitFromString('%s')\" SUCCEDED, expected -> %v, got -> %v", input, expected, result)
		} else {
			t.Errorf("\"Test_FirewallField_Address_TryInitFromString('%s')\" FAILED, expected -> %v, got -> %v", input, expected, result)
		}
	}
	{
		input := "1.1.1.1/32/34"
		expected := netip.AddrFrom4([4]byte{1, 1, 1, 1})

		addr := Address{}
		addr.TryInitFromString(input)
		result := addr.GetValue()

		if result == expected {
			t.Logf("\"Test_FirewallField_Address_TryInitFromString('%s')\" SUCCEDED, expected -> %v, got -> %v", input, expected, result)
		} else {
			t.Errorf("\"Test_FirewallField_Address_TryInitFromString('%s')\" FAILED, expected -> %v, got -> %v", input, expected, result)
		}
	}
	{
		input := "1.1.1.1"
		expected := netip.AddrFrom4([4]byte{1, 1, 1, 1})

		addr := Address{}
		addr.TryInitFromString(input)
		result := addr.GetValue()

		if result == expected {
			t.Logf("\"Test_FirewallField_Address_TryInitFromString('%s')\" SUCCEDED, expected -> %v, got -> %v", input, expected, result)
		} else {
			t.Errorf("\"Test_FirewallField_Address_TryInitFromString('%s')\" FAILED, expected -> %v, got -> %v", input, expected, result)
		}
	}
	{
		input := "123.222.45.65"
		expected := netip.AddrFrom4([4]byte{123, 222, 45, 65})

		addr := Address{}
		addr.TryInitFromString(input)
		result := addr.GetValue()

		if result == expected {
			t.Logf("\"Test_FirewallField_Address_TryInitFromString('%s')\" SUCCEDED, expected -> %v, got -> %v", input, expected, result)
		} else {
			t.Errorf("\"Test_FirewallField_Address_TryInitFromString('%s')\" FAILED, expected -> %v, got -> %v", input, expected, result)
		}
	}
	{
		input := "8.8.8.8/8"
		expected := netip.AddrFrom4([4]byte{8, 8, 8, 8})

		addr := Address{}
		addr.TryInitFromString(input)
		result := addr.GetValue()

		if result == expected {
			t.Logf("\"Test_FirewallField_Address_TryInitFromString('%s')\" SUCCEDED, expected -> %v, got -> %v", input, expected, result)
		} else {
			t.Errorf("\"Test_FirewallField_Address_TryInitFromString('%s')\" FAILED, expected -> %v, got -> %v", input, expected, result)
		}
	}
	{
		input := "8.8.8.8/81241242352345"
		expected := netip.AddrFrom4([4]byte{8, 8, 8, 8})

		addr := Address{}
		addr.TryInitFromString(input)
		result := addr.GetValue()

		if result == expected {
			t.Logf("\"Test_FirewallField_Address_TryInitFromString('%s')\" SUCCEDED, expected -> %v, got -> %v", input, expected, result)
		} else {
			t.Errorf("\"Test_FirewallField_Address_TryInitFromString('%s')\" FAILED, expected -> %v, got -> %v", input, expected, result)
		}
	}
	{
		input := "286.32.41.5/8"
		expected := netip.Addr{}

		addr := Address{}
		addr.TryInitFromString(input)
		result := addr.GetValue()

		if result == expected {
			t.Logf("\"Test_FirewallField_Address_TryInitFromString('%s')\" SUCCEDED, expected -> %v, got -> %v", input, expected, result)
		} else {
			t.Errorf("\"Test_FirewallField_Address_TryInitFromString('%s')\" FAILED, expected -> %v, got -> %v", input, expected, result)
		}
	}
	{
		input := "8.2.8/8"
		expected := netip.Addr{}

		addr := Address{}
		addr.TryInitFromString(input)
		result := addr.GetValue()

		if result == expected {
			t.Logf("\"Test_FirewallField_Address_TryInitFromString('%s')\" SUCCEDED, expected -> %v, got -> %v", input, expected, result)
		} else {
			t.Errorf("\"Test_FirewallField_Address_TryInitFromString('%s')\" FAILED, expected -> %v, got -> %v", input, expected, result)
		}
	}
	{
		input := ""
		expected := netip.Addr{}

		addr := Address{}
		addr.TryInitFromString(input)
		result := addr.GetValue()

		if result == expected {
			t.Logf("\"Test_FirewallField_Address_TryInitFromString('%s')\" SUCCEDED, expected -> %v, got -> %v", input, expected, result)
		} else {
			t.Errorf("\"Test_FirewallField_Address_TryInitFromString('%s')\" FAILED, expected -> %v, got -> %v", input, expected, result)
		}
	}
	{
		input := "1/22.1.1.1"
		expected := netip.Addr{}

		addr := Address{}
		addr.TryInitFromString(input)
		result := addr.GetValue()

		if result == expected {
			t.Logf("\"Test_FirewallField_Address_TryInitFromString('%s')\" SUCCEDED, expected -> %v, got -> %v", input, expected, result)
		} else {
			t.Errorf("\"Test_FirewallField_Address_TryInitFromString('%s')\" FAILED, expected -> %v, got -> %v", input, expected, result)
		}
	}
}

func Test_FirewallField_Address_SetGetValueString(t *testing.T) {
	{
		input := "1.1.1.1/32"
		expected := netip.AddrFrom4([4]byte{1, 1, 1, 1})

		addr := Address{}
		addr.SetValue(expected)
		result := addr.GetValue()

		if result == expected {
			t.Logf("\"Test_FirewallField_Address_SetGetValuString('%s')\" SUCCEDED, expected -> %v, got -> %v", input, expected, result)
		} else {
			t.Errorf("\"Test_FirewallField_Address_SetGetValuString('%s')\" FAILED, expected -> %v, got -> %v", input, expected, result)
		}
	}
	{
		input := "1.11.1/32"
		expected := netip.Addr{}

		addr := Address{}
		addr.SetValue(expected)
		result := addr.GetValue()

		if result == expected {
			t.Logf("\"Test_FirewallField_Address_SetGetValuString('%s')\" SUCCEDED, expected -> %v, got -> %v", input, expected, result)
		} else {
			t.Errorf("\"Test_FirewallField_Address_SetGetValuString('%s')\" FAILED, expected -> %v, got -> %v", input, expected, result)
		}
	}
}

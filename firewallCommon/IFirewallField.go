package firewallCommon

// TODO: Resolve import loop
type IFirewallField interface {
	TryInitFromString(param string) error
	GetString() string
}

package firewallCommon

type IFirewallField interface {
	TryInitFromRest(param string) error
	TryInitFromIpTables(param string) error
	GetString() string
	MarshalRest() string
	MarshalIpTables() string
}

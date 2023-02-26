package firewallField

type Text struct {
	value string
}

func (ctx *Text) TryInitFromString(param string) error {
	ctx.value = param
	return nil
}

func (ctx *Text) TryInitFromRest(param string) error {
	return ctx.TryInitFromString(param)
}

func (ctx *Text) TryInitFromIpTables(param string) error {
	return ctx.TryInitFromString(param)
}

func TextTypeFromString(idString string) (Text, error) {
	id := Text{}
	return id, id.TryInitFromString(idString)
}

func TextTypeFromValue(value string) (Text, error) {
	return Text{value: value}, nil
}

func (ctx *Text) SetValue(value string) {
	ctx.value = value
}

func (ctx Text) GetValue() string {
	return ctx.value
}

func (ctx Text) GetString() string {
	return ctx.value
}

func (ctx Text) MarshalRest() string {
	return ctx.value
}

func (ctx Text) MarshalIpTables() string {
	return ctx.value
}

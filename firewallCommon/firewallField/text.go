package firewallField

type Text struct {
	value string
}

func TextNew(val string) Text {
	return Text{value: val}
}

func (ctx *Text) TryInitFromString(param string) error {
	ctx.value = param
	return nil
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

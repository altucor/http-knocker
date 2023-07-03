package firewallField

type Text struct {
	value string
}

func (ctx *Text) TryInitFromString(param string) error {
	ctx.value = param
	return nil
}

func TextTypeFromString(idString string) (Text, error) {
	id := Text{}
	return id, id.TryInitFromString(idString)
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

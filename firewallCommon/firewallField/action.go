package firewallField

import (
	"errors"
	"strings"
)

type ActionType uint8

const (
	ACTION_INVALID ActionType = 0xFF
	ACTION_ACCEPT  ActionType = 0
	ACTION_DROP    ActionType = 1
	ACTION_REJECT  ActionType = 3
	ACTION_JUMP    ActionType = 4
	ACTION_LOG     ActionType = 5
)

var (
	actionMap = map[ActionType]string{
		ACTION_INVALID: "<INVALID>",
		ACTION_ACCEPT:  "accept",
		ACTION_DROP:    "drop",
		ACTION_REJECT:  "reject",
		ACTION_JUMP:    "jump",
		ACTION_LOG:     "log",
	}
)

type Action struct {
	value ActionType
}

func (ctx *Action) TryInitFromString(param string) error {
	param = strings.ToLower(param)
	for key, value := range actionMap {
		if value == param {
			ctx.value = key
			return nil
		}
	}
	ctx.value = ACTION_INVALID
	return errors.New("cannot init from string")
}

func ActionTypeFromString(chainString string) (Action, error) {
	chainString = strings.ToLower(chainString)
	for key, value := range actionMap {
		if value == chainString {
			return Action{value: key}, nil
		}
	}

	return Action{value: ACTION_INVALID}, errors.New("invalid action text name")
}

func (ctx *Action) SetValue(value ActionType) {
	ctx.value = value
}

func (ctx Action) GetValue() ActionType {
	return ctx.value
}

func (ctx Action) GetString() string {
	return actionMap[ctx.value]
}

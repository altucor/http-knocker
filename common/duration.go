package common

import (
	"errors"
	"time"
)

type DurationSeconds struct {
	value time.Duration
}

func (ctx *DurationSeconds) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var durationStr string = ""
	err := unmarshal(&durationStr)
	if err != nil {
		return err
	}
	duration, err := DurationSecondsFromString(durationStr)
	if err != nil {
		return err
	}
	ctx.SetValue(duration.GetValue())
	return nil
}

func DurationSecondsFromString(input string) (DurationSeconds, error) {
	// here parsing of human readable values like 1h30m20s
	var err error = nil
	duration := DurationSeconds{
		value: 0,
	}
	duration.value, err = time.ParseDuration(input)
	if duration.value < 0 {
		return DurationSeconds{}, errors.New("duration cannot be lower than zero")
	}
	return duration, err
}

func (ctx *DurationSeconds) SetValue(seconds time.Duration) {
	ctx.value = seconds
}

func (ctx DurationSeconds) GetValue() time.Duration {
	return ctx.value
}

func (ctx DurationSeconds) GetSeconds() uint64 {
	return uint64(ctx.value.Seconds())
}

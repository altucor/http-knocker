package common

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/altucor/http-knocker/logging"
)

type DurationSeconds struct {
	// value uint64
	// TODO: Rewrite to type "time.Duration"
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

func getValueByRegex(input string, keyAnchor string) (uint64, error) {
	regexRule, err := regexp.Compile("(\\d+" + keyAnchor + ")")
	if err != nil {
		logging.CommonLog().Error("[DurationSeconds] Error compiling regexRule: %s\n", err)
		return 0, err
	}
	tempStr := strings.ReplaceAll(regexRule.FindString(input), keyAnchor, "")
	if tempStr == "" {
		return 0, nil
	}
	return strconv.ParseUint(tempStr, 10, 64)
}

func DurationSecondsFromString(input string) (DurationSeconds, error) {
	// here parsing of human readable values like 1h30m20s
	if !strings.ContainsAny(input, "hms") {
		logging.CommonLog().Error("[DurationSeconds] Error invalid duration format: cannot find \"h\" or \"m\" or \"s\" modifiers\n")
		return DurationSeconds{}, errors.New("cannot find \"h\" or \"m\" or \"s\" modifiers")
	}
	hours, err := getValueByRegex(input, "h")
	if err != nil {
		logging.CommonLog().Error("[DurationSeconds] Error getting hours: %s\n", err)
		return DurationSeconds{}, err
	}
	minutes, err := getValueByRegex(input, "m")
	if err != nil {
		logging.CommonLog().Error("[DurationSeconds] Error getting minutes: %s\n", err)
		return DurationSeconds{}, err
	}
	seconds, err := getValueByRegex(input, "s")
	if err != nil {
		logging.CommonLog().Error("[DurationSeconds] Error getting seconds: %s\n", err)
		return DurationSeconds{}, err
	}
	// duration := DurationSeconds{
	// 	value: (hours * 60 * 60) + (minutes * 60) + seconds,
	// }
	duration := DurationSeconds{
		value: time.Duration((hours * uint64(time.Hour)) + (minutes * uint64(time.Minute)) + (seconds * uint64(time.Second))),
	}
	return duration, nil
}

func (ctx *DurationSeconds) SetValue(seconds time.Duration) {
	ctx.value = seconds
}

func (ctx DurationSeconds) GetValue() time.Duration {
	return ctx.value
}

func (ctx DurationSeconds) GetSeconds() uint64 {
	return uint64(ctx.value / time.Second)
}

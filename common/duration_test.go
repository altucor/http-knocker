package common

import (
	"errors"
	"testing"
	"time"
)

func Test_DurationSeconds_FromString(t *testing.T) {
	input := "2h3m46s"
	expected := DurationSeconds{value: (2 * time.Hour) + (3 * time.Minute) + (46 * time.Second)}
	result, err := DurationSecondsFromString(input)

	if err != nil {
		t.Error("Error from DurationSecondsFromString:", err)
	}

	if result != expected {
		t.Errorf("\"DurationSecondsFromString('%s')\" FAILED, expected -> %v, got -> %v", input, expected, result)
	} else {
		t.Logf("\"DurationSecondsFromString('%s')\" SUCCEDED, expected -> %v, got -> %v", input, expected, result)
	}
}

func Test_DurationSeconds_GetSeconds(t *testing.T) {
	input := "57s"
	expected := DurationSeconds{value: (57 * time.Second)}.GetSeconds()
	result, err := DurationSecondsFromString(input)

	if err != nil {
		t.Error("Error from DurationSecondsFromString:", err)
	}

	if result.GetSeconds() != expected {
		t.Errorf("\"DurationSecondsFromString('%s')\" FAILED, expected -> %v, got -> %v", input, expected, result)
	} else {
		t.Logf("\"DurationSecondsFromString('%s')\" SUCCEDED, expected -> %v, got -> %v", input, expected, result)
	}
}

func Test_DurationSeconds_SetGetValue(t *testing.T) {
	duration := DurationSeconds{}
	expected := 2 * time.Minute
	duration.SetValue(expected)
	result := duration.GetValue()

	if result != expected {
		t.Errorf("\"Test_DurationSeconds_SetGetValue('%s')\" FAILED, expected -> %v, got -> %v", expected, expected, result)
	} else {
		t.Logf("\"Test_DurationSeconds_SetGetValue('%s')\" SUCCEDED, expected -> %v, got -> %v", expected, expected, result)
	}
}

func Test_DurationSeconds_NegativeValue(t *testing.T) {
	input := "-2h3m46s"
	_, err := DurationSecondsFromString(input)
	expected := errors.New("duration cannot be lower than zero")

	if err.Error() == expected.Error() {
		t.Logf("\"DurationSeconds_NegativeValue('%s')\" SUCCEDED, expected -> %v, got -> %v", input, expected, err.Error())
	} else {
		t.Errorf("\"DurationSeconds_NegativeValue('%s')\" FAILED, expected -> %v, got -> %v", input, expected, err.Error())
	}
}

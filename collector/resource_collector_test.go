package collector

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestParseDuration(t *testing.T) {
	var testCases = []struct {
		input    string
		output   float64
		hasError bool
	}{
		{
			"3d3h42m53s",
			272573,
			false,
		},
		{
			"15w3d3h42m53s",
			9344573,
			false,
		},
		{
			"42m53s",
			2573,
			false,
		},
		{
			"7w6d9h34m",
			4786440,
			false,
		},
		{
			"59",
			0,
			true,
		},
		{
			"s",
			0,
			false,
		},
		{
			"",
			0,
			false,
		},
	}

	for _, testCase := range testCases {
		f, err := parseDuration(testCase.input)

		switch testCase.hasError {
		case true:
			assert.Error(t, err)
		case false:
			assert.NoError(t, err)
		}

		assert.Equal(t, testCase.output, f)
	}
}

func TestParseWirelessRate(t *testing.T) {
	var testCases = []struct {
		input    string
		output   float64
		hasError bool
	}{
		{
			"1Mbps",
			1,
			false,
		},
		{
			"702Mbps-80MHz/2S",
			702,
			false,
		},
		{
			"433.3Mbps-80MHz/1S/SGI",
			433.3,
			false,
		},
		{
			"Mbps",
			0,
			true,
		},
		{
			"433..3Mbps",
			0,
			true,
		},
		{
			"",
			0,
			true,
		},
	}

	for _, testCase := range testCases {
		f, err := parseWirelessRate(testCase.input)

		switch testCase.hasError {
		case true:
			assert.Error(t, err)
		case false:
			assert.NoError(t, err)
		}

		assert.Equal(t, testCase.output, f)
	}
}

func TestParseDatetime(t *testing.T) {
	var testCases = []struct {
		input    string
		output   time.Time
		hasError bool
	}{
		{
			"sep/08/2019 18:09:55",
			time.Date(2019, 9, 8, 18, 9, 55, 0, time.UTC),
			false,
		},
		{
			"oct/05/2019 16:34:15",
			time.Date(2019, 10, 5, 16, 34, 15, 0, time.UTC),
			false,
		},
		{
			"oct-05-2019 16:34:15",
			time.Time{},
			true,
		},
		{
			"16:34:15",
			time.Time{},
			true,
		},
		{
			"25",
			time.Time{},
			true,
		},
	}

	for _, testCase := range testCases {
		tt, err := parseDatetime(testCase.input)

		switch testCase.hasError {
		case true:
			assert.Error(t, err)
		case false:
			assert.NoError(t, err)
		}

		assert.Equal(t, testCase.output, tt)
	}
}

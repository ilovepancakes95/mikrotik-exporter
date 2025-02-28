package collector

import (
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
	"time"
)

func TestSplitStringToFloats(t *testing.T) {
	var testCases = []struct {
		input    string
		expected struct {
			f1 float64
			f2 float64
		}
		isNaN    bool
		hasError bool
	}{
		{
			"1.2,2.1",
			struct {
				f1 float64
				f2 float64
			}{
				1.2,
				2.1,
			},
			false,
			false,
		},
		{
			input:    "1.2,",
			isNaN:    true,
			hasError: true,
		},
		{
			input:    ",2.1",
			isNaN:    true,
			hasError: true,
		},
		{
			"1.2,2.1,3.2",
			struct {
				f1 float64
				f2 float64
			}{
				1.2,
				2.1,
			},
			false,
			false,
		},
		{
			input:    "",
			isNaN:    true,
			hasError: true,
		},
	}

	for _, testCase := range testCases {
		f1, f2, err := splitStringToFloats(testCase.input)

		switch testCase.hasError {
		case true:
			assert.Error(t, err)
		case false:
			assert.NoError(t, err)
		}

		switch testCase.isNaN {
		case true:
			assert.True(t, math.IsNaN(f1))
			assert.True(t, math.IsNaN(f2))
		case false:
			assert.Equal(t, testCase.expected.f1, f1)
			assert.Equal(t, testCase.expected.f2, f2)
		}
	}
}

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

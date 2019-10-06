package collector

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

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
			"1",
			0,
			true,
		},
		{
			"433.3",
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

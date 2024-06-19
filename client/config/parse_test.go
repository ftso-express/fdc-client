package config_test

import (
	"fmt"
	"local/fdc/client/config"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStringToByte32(t *testing.T) {

	const a = "12!Ab( )"

	bytes, err := config.StringToByte32(a)

	require.NoError(t, err)

	result := [32]byte{49, 50, 33, 65, 98, 40, 32, 41}

	require.Equal(t, result, bytes, fmt.Sprintf("bytes %v do not match expectation, %v", bytes, result))

	c := strings.Repeat("A", 33)

	bytes, err = config.StringToByte32(c)

	require.Error(t, err)

	require.Equal(t, [32]byte{}, bytes, fmt.Sprintf("bytes %v do not match expectation, %v", bytes, result))

}

func TestTwoStringsToByte64(t *testing.T) {

	const a = "12!Ab( )"
	const b = "11"

	bytes, err := config.TwoStringsToByte64(a, b)

	require.NoError(t, err)

	result := [64]byte{49, 50, 33, 65, 98, 40, 32, 41}
	result[32] = 49
	result[33] = 49

	if bytes != result {
		t.Errorf("bytes %v do not match the expected result %v", bytes, result)

	}

}

func TestWhiteSpaceStrip(t *testing.T) {

	tests := []struct {
		input  string
		output string
	}{
		{
			input:  "a s\vd \t ad \f\n YY \n",
			output: "asdadYY",
		},
		{
			input:  "    ",
			output: "",
		},
		{
			input:  "  1  ",
			output: "1",
		},
		{
			input:  "  \n\f  ",
			output: "",
		},
	}

	for i, test := range tests {

		output := config.WhiteSpaceStrip(test.input)

		require.Equal(t, test.output, output, fmt.Sprintf("wrong output test %d", i))
	}

}

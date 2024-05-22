package config_test

import (
	"local/fdc/client/config"
	"strings"
	"testing"
)

func TestStringToByte32(t *testing.T) {

	const a = "12!Ab( )"

	bytes, err := config.StringToByte32(a)

	if err != nil {
		t.Errorf("unexpected error %s", err)
	}

	result := [32]byte{49, 50, 33, 65, 98, 40, 32, 41}

	if bytes != result {
		t.Errorf("bytes %v do not match the expected result %v", bytes, result)

	}

	c := strings.Repeat("A", 33)

	bytes, err = config.StringToByte32(c)

	if err == nil {
		t.Errorf("unexpected fail error %s", err)
	}

	if bytes != [32]byte{} {

		t.Errorf("bytes %v do not match the expected result %v", bytes, [32]byte{})

	}

}

func TestTwoStringsToByte64(t *testing.T) {

	const a = "12!Ab( )"
	const b = "11"

	bytes, err := config.TwoStringsToByte64(a, b)

	if err != nil {
		t.Errorf("unexpected error %s", err)
	}

	result := [64]byte{49, 50, 33, 65, 98, 40, 32, 41}
	result[32] = 49
	result[33] = 49

	if bytes != result {
		t.Errorf("bytes %v do not match the expected result %v", bytes, result)

	}

}

func TestWhiteSpaceStrip(t *testing.T) {

	const a = "a s\vd \t ad \f\n YY \n"

	aStriped := config.WhiteSpaceStrip(a)

	if aStriped != "asdadYY" {
		t.Errorf("expected %s, got %s", "asdadYY", aStriped)
	}

}

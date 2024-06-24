package restServer

import (
	"encoding/hex"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

const hexPrefix = "0x"

func ValidateEVMAddressString(address string) bool {
	address = strings.TrimPrefix(address, hexPrefix)
	dec, err := hex.DecodeString(address)
	if err != nil {
		return false
	}

	if len(dec) != 20 {
		return false
	}

	return err == nil
}

package restServer

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var (
	validate *validator.Validate
)

const (
	hexPrefix = "0x"
)

func init() {
	validate = validator.New(validator.WithRequiredStructEnabled())
	// if err := validate.RegisterValidation("submitAddress", validateEVMAddress); err != nil {
	// 	log.Fatal(err)
	// }
}

func validateEVMAddress(fl validator.FieldLevel) bool {
	fmt.Println("Validating EVM address")
	val := fl.Field().String()
	return ValidateEVMAddressString(val)
}

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

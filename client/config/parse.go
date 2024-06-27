package config

import (
	"errors"
	"flare-common/errorf"
	"fmt"
	"os"
	"strings"
	"unicode"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

// ParseAttestationTypesConfig parses AttestationTypesUnparsed as read from toml file into AttestationTypes.
func ParseAttestationTypesConfig(attTypesConfigUnparsed AttestationTypesUnparsed) (AttestationTypes, error) {

	attTypesConfig := make(AttestationTypes)

	for k := range attTypesConfigUnparsed {

		attType, err := StringToByte32(k)

		if err != nil {
			return nil, fmt.Errorf("reading type %w", err)
		}

		attTypeConfig, err := ParseAttestationTypeConfig(attTypesConfigUnparsed[k])

		if err != nil {
			return nil, fmt.Errorf("parsing type %s: %w", k, err)
		}

		attTypesConfig[attType] = attTypeConfig
	}

	return attTypesConfig, nil

}

// getAbi reads abi of a struct from a JSON file and converts it into abi.Arguments and string representation.
func getAbi(path string) (abi.Arguments, string, error) {

	file, err := os.ReadFile(path)

	if err != nil {
		return abi.Arguments{}, "", errorf.ReadingFile(path, err)
	}

	args, err := ArgumentsFromAbi(file)

	if err != nil {
		return abi.Arguments{}, "", fmt.Errorf("retrieving arguments from %s with %w", path, err)
	}

	abiString := WhiteSpaceStrip(string(file))

	return args, abiString, nil

}

// ArgumentsFromAbi convert byte encoded json abi into abu.Arguments.
func ArgumentsFromAbi(abiBytes []byte) (abi.Arguments, error) {

	var arg abi.Argument

	err := arg.UnmarshalJSON(abiBytes)

	if err != nil {
		return abi.Arguments{}, err
	}

	return abi.Arguments{arg}, nil

}

// parseSourceConfig takes sourceBig and converts LutLimit from big.int to uint64.
func parseSourceConfig(sourceConfigBig sourceBig) (Source, error) {

	if !sourceConfigBig.LutLimit.IsUint64() {
		return Source{
				Url:      sourceConfigBig.Url,
				ApiKey:   sourceConfigBig.ApiKey,
				LutLimit: 0,
			},
			errors.New("lutLimit does not fit in uint64")

	}

	return Source{
			Url:      sourceConfigBig.Url,
			ApiKey:   sourceConfigBig.ApiKey,
			LutLimit: sourceConfigBig.LutLimit.Uint64(),
		},
		nil

}

func ParseAttestationTypeConfig(attTypeConfigUnparsed AttestationTypeUnparsed) (AttestationType, error) {

	responseArguments, responseAbiString, err := getAbi(attTypeConfigUnparsed.Abi)

	if err != nil {
		return AttestationType{}, fmt.Errorf("getting abi %w", err)
	}

	sourcesConfig, err := parseSourcesConfig(attTypeConfigUnparsed.Sources)

	if err != nil {
		return AttestationType{}, fmt.Errorf("parsing: %w", err)

	}

	return AttestationType{
			ResponseArguments:  responseArguments,
			ResponseAbisString: responseAbiString,
			SourcesConfig:      sourcesConfig,
		},
		nil

}

func parseSourcesConfig(sourcesConfigUnparsed map[string]sourceBig) (map[[32]byte]Source, error) {

	sourcesConfig := make(map[[32]byte]Source)

	for k := range sourcesConfigUnparsed {

		source, err := StringToByte32(k)

		if err != nil {
			return nil, fmt.Errorf("reading source %w", err)
		}

		sourceConfig, err := parseSourceConfig(sourcesConfigUnparsed[k])

		if err != nil {
			return nil, fmt.Errorf("parsing source config %w", err)
		}

		sourcesConfig[source] = sourceConfig

	}

	return sourcesConfig, nil
}

// WhiteSpaceStrip removes any white space character from the string.
func WhiteSpaceStrip(str string) string {
	var b strings.Builder
	b.Grow(len(str))
	for _, ch := range str {
		if !unicode.IsSpace(ch) {
			b.WriteRune(ch)
		}
	}
	return b.String()
}

// StringToByte32 converts string str to utf-8 encoding and writes it to [32]byte.
// If str is longer than 32 it returns an error.
func StringToByte32(str string) ([32]byte, error) {

	var strBytes [32]byte
	if len(str) > 32 {
		return strBytes, fmt.Errorf("string %s to long", str)
	}

	copy(strBytes[:], []byte(str))

	return strBytes, nil

}

// TowStringToByte64 converts each of the two strings to utf-8 encoding and writes it to [32]byte and concatenates the result.
// If any of the string is longer than 32 it returns an error.
func TwoStringsToByte64(str1, str2 string) ([64]byte, error) {

	var strBytes [64]byte
	if len(str1) > 32 {
		return strBytes, fmt.Errorf("first string %s to long", str1)
	}
	if len(str2) > 32 {
		return strBytes, fmt.Errorf("second string %s to long", str2)
	}

	copy(strBytes[0:32], []byte(str1))

	copy(strBytes[32:64], []byte(str2))

	return strBytes, nil

}

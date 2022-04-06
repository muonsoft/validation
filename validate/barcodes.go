package validate

import (
	"errors"
	"fmt"
	"strconv"
)

var (
	ErrOnlyZeros        = errors.New("contains only zeros")
	ErrInvalidChecksum  = errors.New("invalid checksum")
	ErrUnexpectedLength = errors.New("unexpected length")
	ErrContainsNonDigit = errors.New("contains non-digit")
)

// EAN8 checks that string contains valid EAN-8 code.
//
// If the value is not valid then one of the errors will be returned:
//	• ErrOnlyZeros if code contains only zeros;
//	• ErrInvalidChecksum if check digit is not valid;
//	• ErrUnexpectedLength if value length has an unexpected size;
//	• ErrContainsNonDigit if string contains non-digit value.
//
// See https://en.wikipedia.org/wiki/EAN-8.
func EAN8(value string) error {
	return validateBarcode(value, 8)
}

// EAN13 checks that string contains valid EAN-13 code.
//
// If the value is not valid then one of the errors will be returned:
//	• ErrOnlyZeros if code contains only zeros;
//	• ErrInvalidChecksum if check digit is not valid;
//	• ErrUnexpectedLength if value length has an unexpected size;
//	• ErrContainsNonDigit if string contains non-digit value.
//
// See https://en.wikipedia.org/wiki/International_Article_Number.
func EAN13(value string) error {
	return validateBarcode(value, 13)
}

// UPCA checks that string contains valid UPC-A code.
//
// If the value is not valid then one of the errors will be returned:
//	• ErrOnlyZeros if code contains only zeros;
//	• ErrInvalidChecksum if check digit is not valid;
//	• ErrUnexpectedLength if value length has an unexpected size;
//	• ErrContainsNonDigit if string contains non-digit value.
//
// See https://en.wikipedia.org/wiki/Universal_Product_Code.
func UPCA(value string) error {
	return validateBarcode(value, 12)
}

// UPCE checks that string contains valid UPC-E code.
//
// If the value is not valid then one of the errors will be returned:
//	• ErrOnlyZeros if code contains only zeros;
//	• ErrInvalid if 8-digits code starts with number not equal to 0;
//	• ErrInvalidChecksum if check digit is not valid;
//	• ErrUnexpectedLength if value length has an unexpected size;
//	• ErrContainsNonDigit if string contains non-digit value.
//
// See https://en.wikipedia.org/wiki/Universal_Product_Code#UPC-E.
func UPCE(value string) error {
	upce, err := decodeUPCE(value)
	if err != nil {
		return err
	}
	// No check digit, it's always correct
	if len(upce) == 6 {
		return nil
	}

	upca := expandUPCEToUPCA(upce)
	if upcaChecksum(upca) != upca[11] {
		return ErrInvalidChecksum
	}

	return nil
}

func validateBarcode(value string, size int) error {
	sum, err := barcodeChecksum(value, size)
	if err != nil {
		return err
	}

	if strconv.Itoa(sum) != value[size-1:size] {
		return ErrInvalidChecksum
	}

	return nil
}

func barcodeChecksum(barcode string, size int) (int, error) {
	if len(barcode) != size {
		return -1, ErrUnexpectedLength
	}

	code := barcode[:size-1]
	multiplyWhenEven := size%2 == 0
	sum := 0

	for i, v := range code {
		digit, isDigit := digits[v]
		if !isDigit {
			return -1, fmt.Errorf("%w: %q", ErrContainsNonDigit, v)
		}

		if (i%2 == 0) == multiplyWhenEven {
			sum += 3 * digit
		} else {
			sum += digit
		}
	}
	if sum == 0 {
		return -1, ErrOnlyZeros
	}

	return (10 - sum%10) % 10, nil
}

func decodeUPCE(value string) ([]byte, error) {
	upce := make([]byte, 0, 8)
	sum := 0
	for i, v := range value {
		digit, isDigit := digits[v]
		if !isDigit {
			return nil, fmt.Errorf("%w: %q", ErrContainsNonDigit, v)
		}
		if i >= 8 {
			return nil, ErrUnexpectedLength
		}
		upce = append(upce, byte(digit))
		sum += digit
	}
	if len(upce) < 6 {
		return nil, ErrUnexpectedLength
	}
	if sum == 0 {
		return nil, ErrOnlyZeros
	}

	// Only 0 is allowed for 8-digit code
	if len(upce) == 8 && upce[0] != 0 {
		return nil, ErrInvalid
	}

	return upce, nil
}

// expandUPCEToUPCA converts UPC-E code into equivalent UPC-A code.
// See https://en.wikipedia.org/wiki/Universal_Product_Code#UPC-E for details.
func expandUPCEToUPCA(upce []byte) [12]byte {
	var checkdigit byte
	var code []byte
	if len(upce) == 7 {
		// 7th digit is check digit. No number system digit
		code = upce[0:6]
		checkdigit = upce[6]
	} else {
		// Both check and number system digits
		code = upce[1:7]
		checkdigit = upce[7]
	}

	var upca [12]byte
	upca[11] = checkdigit
	switch code[5] {
	case 0, 1, 2:
		upca[1] = code[0]
		upca[2] = code[1]
		upca[3] = code[5]
		upca[8] = code[2]
		upca[9] = code[3]
		upca[10] = code[4]
	case 3:
		upca[1] = code[0]
		upca[2] = code[1]
		upca[3] = code[2]
		upca[9] = code[3]
		upca[10] = code[4]
	case 4:
		upca[1] = code[0]
		upca[2] = code[1]
		upca[3] = code[2]
		upca[4] = code[3]
		upca[10] = code[4]
	default:
		upca[1] = code[0]
		upca[2] = code[1]
		upca[3] = code[2]
		upca[4] = code[3]
		upca[5] = code[4]
		upca[10] = code[5]
	}

	return upca
}

func upcaChecksum(upca [12]byte) byte {
	var sum byte

	for i := 0; i < 11; i++ {
		if i%2 == 0 {
			sum += 3 * upca[i]
		} else {
			sum += upca[i]
		}
	}

	return (10 - sum%10) % 10
}

var digits = map[rune]int{
	'0': 0,
	'1': 1,
	'2': 2,
	'3': 3,
	'4': 4,
	'5': 5,
	'6': 6,
	'7': 7,
	'8': 8,
	'9': 9,
}

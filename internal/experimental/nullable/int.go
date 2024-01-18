package nullable

import (
	"fmt"
	"math"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	TypeNullableInt = schema.TypeString
)

type Int string

func (i Int) IsNull() bool {
	return i == ""
}

func (i Int) Value() (int64, bool, error) {
	if i.IsNull() {
		return 0, true, nil
	}

	value, err := strconv.ParseInt(string(i), 10, 64)
	if err != nil {
		return 0, false, err
	}
	return value, false, nil
}

// ValidateTypeStringNullableInt provides custom error messaging for TypeString ints
// Some arguments require an int value or unspecified, empty field.
func ValidateTypeStringNullableInt(v interface{}, k string) (ws []string, es []error) {
	value, ok := v.(string)
	if !ok {
		es = append(es, fmt.Errorf("expected type of %s to be string", k))
		return
	}

	if value == "" {
		return
	}

	if _, err := strconv.ParseInt(value, 10, 64); err != nil {
		es = append(es, fmt.Errorf("%s: cannot parse '%s' as int: %w", k, value, err))
	}

	return
}

// ValidateTypeStringNullableIntAtLeast provides custom error messaging for TypeString ints
// Some arguments require an int value or unspecified, empty field.
func ValidateTypeStringNullableIntAtLeast(min int64) schema.SchemaValidateFunc {
	return func(i interface{}, k string) (ws []string, es []error) {
		value, ok := i.(string)
		if !ok {
			es = append(es, fmt.Errorf("expected type of %s to be string", k))
			return
		}

		if value == "" {
			return
		}

		v, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			es = append(es, fmt.Errorf("%s: cannot parse '%s' as int: %w", k, value, err))
			return
		}

		if v < min {
			es = append(es, fmt.Errorf("expected %s to be at least (%d), got %d", k, min, v))
		}

		return
	}
}

// ValidateTypeStringNullableIntBetween provides custom error messaging for TypeString ints
// Some arguments require an int value or unspecified, empty field.
func ValidateTypeStringNullableIntBetween(min int64, max int64) schema.SchemaValidateFunc {
	return func(i interface{}, k string) (ws []string, es []error) {
		value, ok := i.(string)
		if !ok {
			es = append(es, fmt.Errorf("expected type of %s to be string", k))
			return
		}

		if value == "" {
			return
		}

		v, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			es = append(es, fmt.Errorf("%s: cannot parse '%s' as int: %w", k, value, err))
			return
		}

		if v < min || v > max {
			es = append(es, fmt.Errorf("expected %s to be at between (%d) and (%d), got %d", k, min, max, v))
		}

		return
	}
}

// ValidateTypeStringNullableIntDivisibleBy returns a SchemaValidateFunc which tests if the provided value
// is of type string, can be converted to int and is divisible by a given number.
func ValidateTypeStringNullableIntDivisibleBy(divisor int64) schema.SchemaValidateFunc {
	return func(i interface{}, k string) (ws []string, es []error) {
		value, ok := i.(string)
		if !ok {
			es = append(es, fmt.Errorf("expected type of %s to be string", k))
			return
		}

		if value == "" {
			return
		}

		v, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			es = append(es, fmt.Errorf("%s: cannot parse '%s' as int: %w", k, value, err))
			return
		}

		if math.Mod(float64(v), float64(divisor)) != 0 {
			es = append(es, fmt.Errorf("expected %s to be divisible by %d, got: %v", k, divisor, i))
		}

		return
	}
}

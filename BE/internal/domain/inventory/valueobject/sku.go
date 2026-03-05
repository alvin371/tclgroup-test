package valueobject

import "errors"

// SKU is a stock-keeping unit identifier.
type SKU string

// ErrEmptySKU is returned when an empty SKU is provided.
var ErrEmptySKU = errors.New("SKU must not be empty")

// NewSKU creates a validated SKU value object.
func NewSKU(s string) (SKU, error) {
	if s == "" {
		return "", ErrEmptySKU
	}
	return SKU(s), nil
}

// String returns the string representation of the SKU.
func (s SKU) String() string {
	return string(s)
}

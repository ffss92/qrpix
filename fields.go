package qrpix

import (
	"errors"
	"fmt"
)

const (
	primitive = "primitive"
	template  = "template"
)

var (
	ErrFieldMetadataNotFound = errors.New("field metadata for provided id not found")

	IDMetadata = map[string]Metadata{
		"00": {
			Name:     "Payload Format Indicator",
			MaxSize:  2,
			MinSize:  2,
			Default:  "01",
			Type:     primitive,
			Required: true,
		},
		"26": {
			Name:     "Merchant Account Information",
			MinSize:  5,
			MaxSize:  99,
			Type:     template,
			Required: true,
		},
		"26-00": {
			Name:     "GUI",
			MaxSize:  14,
			Required: true,
			Type:     primitive,
			Default:  PIXGui,
		},
		"26-01": {
			Name:     "Chave",
			MaxSize:  77,
			Required: true,
			Type:     primitive,
		},
		"26-02": {
			Name:     "Info Adicional",
			MaxSize:  72,
			Required: false,
			Type:     primitive,
		},
		"26-03": {
			Name:     "FSS",
			MaxSize:  8,
			Required: false,
			Type:     primitive,
		},
		"52": {
			Name:     "Merchant Category Code",
			MinSize:  4,
			MaxSize:  4,
			Required: true,
			Type:     primitive,
			// N/I
			Default: "0000",
		},
		"53": {
			Name:     "Transaction Currency",
			MinSize:  3,
			MaxSize:  3,
			Required: true,
			Type:     primitive,
			// BRL
			Default: "986",
		},
		"54": {
			Name:     "Transaction Amount",
			MinSize:  1,
			MaxSize:  13,
			Required: true,
			Type:     primitive,
		},
		"58": {
			Name:     "Country Code",
			MaxSize:  2,
			MinSize:  2,
			Required: true,
			Type:     primitive,
			// ISO3166-1 alpha 2
			Default: "BR",
		},
		"59": {
			Name:     "Merchant Name",
			MinSize:  1,
			MaxSize:  25,
			Required: true,
			Type:     primitive,
		},
		"60": {
			Name:     "Merchant City",
			MinSize:  1,
			MaxSize:  25,
			Required: true,
			Type:     primitive,
		},
		"61": {
			Name:     "Postal Code",
			MinSize:  1,
			MaxSize:  99,
			Type:     primitive,
			Required: false,
		},
		"62": {
			Name:     "Addional Data Field Template",
			MinSize:  5,
			MaxSize:  29,
			Type:     template,
			Required: false,
		},
		"62-05": {
			Name:     "Reference Label",
			MinSize:  1,
			MaxSize:  25,
			Required: false,
			Type:     primitive,
			Default:  "***",
		},
		"63": {
			Name:     "CRC16",
			MaxSize:  4,
			MinSize:  4,
			Required: true,
			Type:     primitive,
		},
	}
)

type Metadata struct {
	Name     string
	MaxSize  int
	MinSize  int
	Default  string
	Type     string
	Required bool
}

// Validates a field value based on the provided id metadata
func ValidateField(id, value string) error {
	meta, ok := IDMetadata[id]
	if !ok {
		return ErrFieldMetadataNotFound
	}

	if !meta.Required && value == "" {
		return nil
	}
	if len(value) > meta.MaxSize {
		return fmt.Errorf("limit above max for field: %s", meta.Name)
	}
	if len(value) < meta.MinSize {
		return fmt.Errorf("limit below min for field: %s", meta.Name)
	}
	if meta.Required && value == "" {
		return fmt.Errorf("no value set for required field: %s", meta.Name)
	}
	return nil
}

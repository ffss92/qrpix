package qrpix

import (
	"errors"
	"fmt"
)

const (
	FieldPrimitive = "primitive"
	FieldTemplate  = "template"
)

var (
	ErrFieldIsRequired       = errors.New("field is required")
	ErrFieldMetadataNotFound = errors.New("field metadata for provided id not found")
)

var (
	IDMetadata = map[string]Metadata{
		"00": {
			Name:     "Payload Format Indicator",
			MaxSize:  2,
			MinSize:  2,
			Type:     FieldPrimitive,
			Required: true,
		},
		"26": {
			Name:     "Merchant Account Information",
			MinSize:  5,
			MaxSize:  99,
			Type:     FieldTemplate,
			Required: true,
		},
		"26-00": {
			Name:     "GUI",
			MaxSize:  14,
			Required: true,
			Type:     FieldPrimitive,
		},
		"26-01": {
			Name:     "Chave",
			MaxSize:  77,
			Required: true,
			Type:     FieldPrimitive,
		},
		"26-02": {
			Name:     "Info Adicional",
			MaxSize:  72,
			Required: false,
			Type:     FieldPrimitive,
		},
		"26-03": {
			Name:     "FSS",
			MaxSize:  8,
			Required: false,
			Type:     FieldPrimitive,
		},
		"52": {
			Name:     "Merchant Category Code",
			MinSize:  4,
			MaxSize:  4,
			Required: true,
			Type:     FieldPrimitive,
		},
		"53": {
			Name:     "Transaction Currency",
			MinSize:  3,
			MaxSize:  3,
			Required: true,
			Type:     FieldPrimitive,
		},
		"54": {
			Name:     "Transaction Amount",
			MinSize:  1,
			MaxSize:  13,
			Required: false,
			Type:     FieldPrimitive,
		},
		"58": {
			Name:     "Country Code",
			MaxSize:  2,
			MinSize:  2,
			Required: true,
			Type:     FieldPrimitive,
		},
		"59": {
			Name:     "Merchant Name",
			MinSize:  1,
			MaxSize:  25,
			Required: true,
			Type:     FieldPrimitive,
		},
		"60": {
			Name:     "Merchant City",
			MinSize:  1,
			MaxSize:  25,
			Required: true,
			Type:     FieldPrimitive,
		},
		"61": {
			Name:     "Postal Code",
			MinSize:  1,
			MaxSize:  99,
			Type:     FieldPrimitive,
			Required: false,
		},
		"62": {
			Name:     "Addional Data Field Template",
			MinSize:  5,
			MaxSize:  29,
			Type:     FieldTemplate,
			Required: false,
		},
		"62-05": {
			Name:     "Reference Label",
			MinSize:  1,
			MaxSize:  25,
			Required: false,
			Type:     FieldPrimitive,
		},
		"63": {
			Name:     "CRC16",
			MaxSize:  4,
			MinSize:  4,
			Required: true,
			Type:     FieldPrimitive,
		},
	}
)

type Metadata struct {
	Name     string
	MaxSize  int
	MinSize  int
	Type     string
	Required bool
}

func GetFieldMetadata(id string) (Metadata, error) {
	meta, ok := IDMetadata[id]
	if !ok {
		return Metadata{}, ErrFieldMetadataNotFound
	}
	return meta, nil
}

// Validates a field value based on the provided id metadata
func ValidateField(id, value string) error {
	meta, err := GetFieldMetadata(id)
	if err != nil {
		return err
	}

	if !meta.Required && value == "" {
		return nil
	}
	if meta.Required && value == "" {
		return ErrFieldIsRequired
	}
	if len(value) > meta.MaxSize {
		return fmt.Errorf("limit above max for field: %s", meta.Name)
	}
	if len(value) < meta.MinSize {
		return fmt.Errorf("limit below min for field: %s", meta.Name)
	}

	return nil
}

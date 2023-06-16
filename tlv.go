package qrpix

import (
	"strconv"
)

// Represents a Type, Value, Lenght (TLV) object
type TLV interface {
	// Returns the (id + length + value) of a given field
	TLV() (string, error)
	// Returns the field id
	FieldID() string
}

// Represents a primitive object, which only contains one value
type Primitive struct {
	ID    string
	Value string

	parentId string
}

func (p Primitive) FieldID() string {
	return p.ID
}

func (p Primitive) validationId() string {
	if p.parentId == "" {
		return p.ID
	}
	return p.parentId + "-" + p.ID
}

func (p Primitive) TLV() (string, error) {
	validationId := p.validationId()
	if err := ValidateField(validationId, p.Value); err != nil {
		return "", err
	}

	// If not value was set, ignore
	if p.Value == "" {
		return "", nil
	}
	length, err := convertLength(p.Value)
	if err != nil {
		return "", err
	}

	return p.ID + length + p.Value, nil
}

// Represets a template object, which contains multiple values.
// Template values also implement the TLV interface
type Template struct {
	ID     string
	values []Primitive
}

func (t Template) FieldID() string {
	return t.ID
}

func (t Template) TLV() (string, error) {

	var value string

	if len(t.values) == 0 || t.values == nil {
		value = ""
	} else {
		for _, p := range t.values {
			tlv, err := p.TLV()
			if err != nil {
				return "", err
			}

			value += tlv
		}
	}

	if err := ValidateField(t.ID, value); err != nil {
		return "", err
	}

	limit, err := convertLength(value)
	if err != nil {
		return "", err
	}

	return (t.ID + limit + value), nil
}

func (t *Template) AddValue(id, value string) {
	t.values = append(t.values, Primitive{ID: id, Value: value, parentId: t.ID})
}

func (t *Template) RemoveValue(id string) bool {
	for i, v := range t.values {
		if v.ID == id {
			t.values = append(t.values[:i], t.values[i+1:]...)
			return true
		}
	}
	return false
}

func (t *Template) UpdateValue(id, value string) bool {
	for _, v := range t.values {
		if v.ID == id {
			v.Value = value
			return true
		}
	}
	return false
}

func (t Template) GetValues() []Primitive {
	return t.values
}

func convertLength(value string) (string, error) {
	l := len(value)
	switch {
	case l < 10:
		return "0" + strconv.Itoa(l), nil
	default:
		return strconv.Itoa(l), nil
	}
}

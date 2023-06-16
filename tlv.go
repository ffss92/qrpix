package qrpix

import (
	"strconv"
	"strings"

	"golang.org/x/exp/slices"
)

// Represents a Type, Value, Lenght (TLV) object.
// TODO: Add a struct representing the data
type TLV interface {
	// Returns the (id, length, value) of a given field
	TLV() (id string, length string, value string, err error)
	// Returns the field id
	FieldID() string
	// Return id + length + value
	Code() (string, error)
	Unwrap() map[string]Primitive
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

func (p Primitive) TLV() (string, string, string, error) {
	validationId := p.validationId()
	if err := ValidateField(validationId, p.Value); err != nil {
		return "", "", "", err
	}
	// If not value was set, ignore
	if p.Value == "" {
		return "", "", "", nil
	}
	length, err := convertLength(p.Value)
	if err != nil {
		return "", "", "", err
	}

	return p.ID, length, p.Value, nil
}

func (p Primitive) Unwrap() map[string]Primitive {
	return nil
}

func (p Primitive) Code() (string, error) {
	id, length, value, err := p.TLV()
	if err != nil {
		return "", err
	}
	return id + length + value, nil
}

// Represets a template object, which contains multiple values.
// Template values also implement the TLV interface
type Template struct {
	ID     string
	values map[string]Primitive
}

func (t Template) FieldID() string {
	return t.ID
}

func (t Template) TLV() (string, string, string, error) {

	b := strings.Builder{}

	primitives := t.Sorted()

	if len(t.values) == 0 || t.values == nil {
		b.WriteString("")
	} else {
		for _, p := range primitives {
			id, length, value, err := p.TLV()
			if err != nil {
				return "", "", "", err
			}
			b.WriteString(id + length + value)
		}
	}
	value := b.String()
	if err := ValidateField(t.ID, value); err != nil {
		return "", "", "", err
	}

	limit, err := convertLength(value)
	if err != nil {
		return "", "", "", err
	}

	return t.ID, limit, value, nil
}

func (t Template) Code() (string, error) {
	id, length, value, err := t.TLV()
	if err != nil {
		return "", err
	}
	return id + length + value, nil
}

func (t *Template) AddValue(id, value string) {
	if t.values == nil {
		t.values = map[string]Primitive{}
	}
	t.values[id] = Primitive{ID: id, Value: value, parentId: t.ID}
}

func (t Template) Unwrap() map[string]Primitive {
	return t.values
}

// Ensures consistency
func (t Template) Sorted() []Primitive {
	primitives := []Primitive{}
	for _, v := range t.values {
		primitives = append(primitives, v)
	}
	slices.SortFunc[Primitive](primitives, func(a, b Primitive) bool {
		aid, _ := strconv.Atoi(a.ID)
		bid, _ := strconv.Atoi(b.ID)
		return aid < bid
	})
	return primitives
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

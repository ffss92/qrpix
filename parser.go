package qrpix

import (
	"errors"
	"fmt"
	"strconv"
)

var (
	ErrEmptyCode = errors.New("cannot parse empty code")
)

type Parser struct {
	Code string
	cur  int
}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) Parse(brCode string) (Builder, error) {
	p.cur = 0 // Reset cursor
	p.Code = brCode
	parts := Builder{} // Reset parts

	if p.Code == "" {
		return nil, ErrEmptyCode
	}

	for p.cur != len(p.Code) {
		id, err := p.readID()
		if err != nil {
			return nil, err
		}

		meta, err := GetFieldMetadata(id)
		if err != nil {
			return nil, fmt.Errorf("failed to get metadata for id %s: %w", id, err)
		}

		switch meta.Type {
		case FieldPrimitive:
			primitive := Primitive{
				ID: id,
			}
			value, err := p.parsePrimitive(id)
			if err != nil {
				return nil, fmt.Errorf("failed to parse primitive with id %s: %w", id, err)
			}
			primitive.Value = value
			parts.Add(primitive)
		case FieldTemplate:
			template := Template{
				ID: id,
			}
			values, err := p.parseTemplate(id)
			if err != nil {
				return nil, fmt.Errorf("failed to parse template with id %s: %w", id, err)
			}
			template.values = values
			parts.Add(template)
		}
	}

	if err := parts.CheckCRC(); err != nil {
		return nil, err
	}

	return parts, nil
}

func (p *Parser) parseTemplate(id string) ([]Primitive, error) {
	primitives := []Primitive{}
	n, err := p.readLength()
	if err != nil {
		return nil, err
	}
	currPos := p.cur
	for currPos+n != p.cur {
		pid, err := p.readID()
		if err != nil {
			return nil, fmt.Errorf("failed to read template value id: %w", err)
		}
		value, err := p.parsePrimitive(id + "-" + pid)
		if err != nil {
			return nil, fmt.Errorf("failed to parse template primitive with id %s: %w", pid, err)
		}
		primitives = append(primitives, Primitive{
			ID:       pid,
			Value:    value,
			parentId: id,
		})
	}
	return primitives, nil
}

func (p *Parser) parsePrimitive(id string) (string, error) {
	n, err := p.readLength()
	if err != nil {
		return "", fmt.Errorf("failed to read primitive length: %w", err)
	}

	value, err := p.readValue(n)
	if err != nil {
		return "", fmt.Errorf("failed to read primitive value: %w", err)
	}

	return value, nil
}

func (p *Parser) readID() (string, error) {
	if err := p.checkSize(2); err != nil {
		return "", fmt.Errorf("failed to read id: %w", err)
	}

	id := p.Code[p.cur : p.cur+2] // Move 2 chars
	p.moveCursor(2)
	return id, nil
}

func (p *Parser) readLength() (int, error) {
	if err := p.checkSize(2); err != nil {
		return 0, fmt.Errorf("failed to read length: %w", err)
	}
	s := p.Code[p.cur : p.cur+2]

	length, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("failed to convert length to int: %w", err)
	}

	p.moveCursor(2)
	return length, nil
}

func (p *Parser) readValue(n int) (string, error) {
	if err := p.checkSize(n); err != nil {
		return "", fmt.Errorf("failed to read value of size %v: %w", n, err)
	}
	value := p.Code[p.cur : p.cur+n]

	p.moveCursor(n)
	return value, nil
}

func (p *Parser) checkSize(n int) error {
	if len(p.Code[p.cur:]) < n {
		return errors.New("seek size exceeds rest of string")
	}
	return nil
}

func (p *Parser) moveCursor(n int) {
	p.cur += n
}

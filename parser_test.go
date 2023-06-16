package qrpix

import (
	"errors"
	"testing"
)

func TestBaseParserMethods(t *testing.T) {

	t.Run("readLength should fail for invalid string or length", func(t *testing.T) {
		cases := []struct {
			code string
		}{
			{code: "AA"},
			{code: "1A"},
			{code: "1"},
		}
		for _, c := range cases {
			p := NewParser()
			p.Code = c.code
			if _, err := p.readLength(); err == nil {
				t.Error("expected error but got nil")
			}
		}

	})

	t.Run("readValue should fail for invalid value lengths", func(t *testing.T) {
		cases := []struct {
			code string
		}{
			{code: "03AA"},
			{code: "041A"},
			{code: "021"},
		}
		for _, c := range cases {
			p := NewParser()
			p.Code = c.code
			length, err := p.readLength()
			if err != nil {
				t.Errorf("unexpected error reading length: %v", err)
			}
			if _, err := p.readValue(length); err == nil {
				t.Error("expected error but got nil")
			}
		}
	})

	t.Run("checkSize should return error if cur + n exceeds Code length", func(t *testing.T) {
		p := NewParser()
		p.Code = "000201"
		p.cur = 4
		if err := p.checkSize(3); err == nil {
			t.Error("expected error but got nil")
		}
	})

	t.Run("checkSize should not return error if cur + n dont exceed Code length", func(t *testing.T) {
		p := NewParser()
		p.Code = "000201"
		p.cur = 0
		if err := p.checkSize(3); err != nil {
			t.Errorf("expected nil but got err: %v", err)
		}
	})

	t.Run("move cursor should add n to the current cursor", func(t *testing.T) {
		p := NewParser()
		steps := []int{4, 5, 2, 3}
		expected := 0

		for _, step := range steps {
			expected += step
			p.moveCursor(step)
			if p.cur != expected {
				t.Errorf("expected cursor to be %v but is %v", expected, p.cur)
			}
		}
	})

	cases := []struct {
		code   string
		id     string
		length int
		value  string
	}{
		{code: "000201", id: "00", length: 2, value: "01"},
		{code: "52040000", id: "52", length: 4, value: "0000"},
		{code: "54041000", id: "54", length: 4, value: "1000"},
		{code: "5303986", id: "53", length: 3, value: "986"},
	}

	for _, c := range cases {
		p := Parser{}
		p.Code = c.code

		t.Run("readID should move cursor by 2 and get correct id for valid cases", func(t *testing.T) {
			id, err := p.readID()
			if err != nil {
				t.Errorf("unexpected error reading id: %v", err)
			}
			if p.cur != 2 {
				t.Errorf("expected %v for cur but got %v", 2, p.cur)
			}
			if id != c.id {
				t.Errorf("expect %s for id but got %s", c.id, id)
			}
		})

		t.Run("readLength should return the correct length and move cursor by 2", func(t *testing.T) {
			expectedCur := p.cur + 2
			lenght, err := p.readLength()
			if err != nil {
				t.Errorf("unexpected error reading length: %v", err)
			}
			if p.cur != expectedCur {
				t.Errorf("expected %v for cur but got %v", expectedCur, p.cur)
			}
			if lenght != c.length {
				t.Errorf("expect %v for length but got %v", c.length, lenght)
			}
		})

		t.Run("readValue should move cursor by the length size and return correct value", func(t *testing.T) {
			expectedCur := p.cur + c.length
			value, err := p.readValue(c.length)
			if err != nil {
				t.Errorf("unexpected error reading length: %v", err)
			}
			if p.cur != expectedCur {
				t.Errorf("expected %v for cur but got %v", expectedCur, p.cur)
			}
			if value != c.value {
				t.Errorf("expect %s for value but got %s", c.value, value)
			}
		})

	}
}

func TestPrimitiveParse(t *testing.T) {
	cases := []struct {
		code     string
		expected string
	}{
		{code: "000201", expected: "01"},
		{code: "52040000", expected: "0000"},
	}

	for _, c := range cases {
		p := NewParser()
		p.Code = c.code
		id, err := p.readID()
		if err != nil {
			t.Error(err)
		}
		value, err := p.parsePrimitive(id)
		if err != nil {
			t.Error(err)
		}
		if value != c.expected {
			t.Errorf("expected %s but got %s", c.expected, value)
		}
	}
}

func TestNewParser(t *testing.T) {
	parser := NewParser()
	if parser == nil {
		t.Error("NewParser should not return a nil pointer to Parser")
	}
}

func TestParser(t *testing.T) {

	t.Run("Parse should fail if there are extra characters", func(t *testing.T) {
		p := NewParser()
		if _, err := p.Parse("0002010"); err == nil {
			t.Error("expected error for cur != len(Code) but got nil")
		}
	})

	t.Run("Parse should fail for empty code or not enough characters", func(t *testing.T) {
		cases := []string{"", "1", "2"}
		for _, c := range cases {
			p := NewParser()
			if _, err := p.Parse(c); err == nil {
				t.Error("expected error but got nil")
			}
		}
	})

	t.Run("valid code", func(t *testing.T) {
		p := Parser{}
		code := "00020126580014br.gov.bcb.pix0136123e4567-e12b-12d1-a456-4266554400005204000053039865802BR5913Fulano de Tal6008BRASILIA62070503***63041D3D"
		builder, err := p.Parse(code)
		if err != nil {
			t.Errorf("unexpected error while parsing valid code: %v", err)
		}
		brCode, err := builder.Build()
		if err != nil {
			t.Error(err)
		}
		if brCode != code {
			t.Errorf("expected %s but got %s", code, brCode)
		}
	})

	t.Run("invalid crc", func(t *testing.T) {
		p := Parser{}
		code := "00020126580014br.gov.bcb.pix0136123e4567-e12b-12d1-a456-4266554400005204000053039865802BR5913Fulano de Tal6008BRASILIA62070503***63041D4D"
		_, err := p.Parse(code)
		if !errors.Is(err, ErrInvalidCRC) {
			t.Errorf("expected ErrInvalidCRC but got: %v", err)
		}
	})

	t.Run("crc not present", func(t *testing.T) {
		p := Parser{}
		code := "00020126580014br.gov.bcb.pix0136123e4567-e12b-12d1-a456-4266554400005204000053039865802BR5913Fulano de Tal6008BRASILIA62070503***"
		_, err := p.Parse(code)
		if !errors.Is(err, ErrCRCNotPresent) {
			t.Errorf("expected ErrCRCNotPresent but got: %v", err)
		}
	})

	t.Run("invalid code", func(t *testing.T) {
		p := Parser{}
		code := "0002026580014br.gov.bcb.pix0136123e4567-e12b-12d1-a456-4266554400005204000053039865802BR5913Fulano de Tal6008BRASILIA62070503***63041D3D"
		_, err := p.Parse(code)
		if err == nil {
			t.Errorf("expected parser to return error but got nil")
		}

		p = Parser{}
		code = "00020"
		_, err = p.Parse(code)
		if err == nil {
			t.Errorf("expected parser to return error but got nil")
		}
	})

}

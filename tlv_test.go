package qrpix

import (
	"testing"
)

func TestTemplateTLV(t *testing.T) {
	expected := "26580014br.gov.bcb.pix0136123e4567-e12b-12d1-a456-426655440000"
	template := Template{
		ID: "26",
	}
	template.AddValue("00", PIXGui)
	template.AddValue("01", "123e4567-e12b-12d1-a456-426655440000")

	code, err := template.Code()
	if err != nil {
		t.Errorf("unexpected error getting tlv from template: %v", err)
	}

	if code != expected {
		t.Errorf("expected %s but got %s", expected, code)
	}

}

func TestPrimitiveTLV(t *testing.T) {
	expected := "000201"
	primitive := Primitive{
		ID:    "00",
		Value: "01",
	}
	code, err := primitive.Code()
	if err != nil {
		t.Errorf("unexpected error getting tlv from template: %v", err)
	}

	if code != expected {
		t.Errorf("expected %s but got %s", expected, code)
	}
}

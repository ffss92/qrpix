package qrpix

import (
	"testing"
)

var (
	exampleCode = "00020126580014br.gov.bcb.pix0136123e4567-e12b-12d1-a456-4266554400005204000053039865802BR5913Fulano de Tal6008BRASILIA62070503***63041D3D"
)

func TestBuilder(t *testing.T) {

	t.Run("valid input should return correct code", func(t *testing.T) {
		builder := Builder{}

		builder.AddPayloadFormatIndicator(PayloadFormatIndicator)
		builder.AddMerchantAccountInformation(PIXGui, "123e4567-e12b-12d1-a456-426655440000")
		builder.AddMerchantCategoryCode("0000")
		builder.AddTransactionCurrency("986")
		builder.AddCountryCode("BR")
		builder.AddMerchantName("Fulano de Tal")
		builder.AddMerchantCity("BRASILIA")
		builder.AddAdditionalDataField("***")

		code, err := builder.Build()
		if err != nil {
			t.Error(err)
		}
		if code != exampleCode {
			t.Errorf("expected %s but got %s", exampleCode, code)
		}
	})

	t.Run("unordered calls should be sorted by id before building", func(t *testing.T) {
		builder := Builder{}

		builder.AddMerchantAccountInformation(PIXGui, "123e4567-e12b-12d1-a456-426655440000")
		builder.AddCountryCode("BR")
		builder.AddMerchantName("Fulano de Tal")
		builder.AddMerchantCity("BRASILIA")
		builder.AddAdditionalDataField("***")
		builder.AddPayloadFormatIndicator(PayloadFormatIndicator)
		builder.AddMerchantCategoryCode("0000")
		builder.AddTransactionCurrency("986")

		code, err := builder.Build()
		if err != nil {
			t.Error(err)
		}
		if code != exampleCode {
			t.Errorf("expected %s but got %s", exampleCode, code)
		}
	})

	t.Run("invalid input (min) should return error", func(t *testing.T) {
		builder := Builder{}
		builder.AddMerchantCategoryCode("000")

		_, err := builder.Build()
		if err == nil {
			t.Error("Build should return an error")
		}
		if err.Error() != "limit below min for field: Merchant Category Code" {
			t.Error("invalid error message")
		}
	})

	t.Run("invalid input (max) should return error", func(t *testing.T) {
		builder := Builder{}
		builder.AddMerchantCategoryCode("00000")

		_, err := builder.Build()
		if err == nil {
			t.Error("Build should return an error")
		}
		if err.Error() != "limit above max for field: Merchant Category Code" {
			t.Error("invalid error message")
		}
	})

	t.Run("optional input should not return error when not set", func(t *testing.T) {
		builder := Builder{}
		builder.AddPostalCode("")
		builder.AddTransactionAmount(0)

		_, err := builder.Build()
		if err != nil {
			t.Error("Should not return error", err)
		}
	})

	t.Run("clear should erase all entries", func(t *testing.T) {
		builder := Builder{}
		builder.AddPayloadFormatIndicator("01")
		builder.Clear()

		if len(builder) != 0 {
			t.Errorf("expected builder lenght of 0 after clear but got %v", len(builder))
		}
	})

	t.Run("add transaction amount should convert int to float", func(t *testing.T) {
		cases := []struct {
			expected string
			value    int
		}{
			{expected: "10.50", value: 1050},
			{expected: "10.00", value: 1000},
			{expected: "12.12", value: 1212},
			{expected: "100.42", value: 10042},
		}

		for _, c := range cases {
			builder := Builder{}
			builder.AddTransactionAmount(c.value)
			ta, ok := builder["54"]
			if !ok {
				t.Error("expected entry with id 54 (Transaction Amount) to be set")
			}

			_, _, value, err := ta.TLV()
			if err != nil {
				t.Errorf("unexpected error while getting transaction amount tlv: %v", err)
			}
			if value != c.expected {
				t.Errorf("expected value to be %s but got %s", c.expected, value)
			}
		}
	})

	t.Run("add transaction of value 0 should be ignore", func(t *testing.T) {
		builder := Builder{}
		builder.AddTransactionAmount(0)
		_, ok := builder["54"]
		if ok {
			t.Error("transaction value should be ignored for value 0 but was set")
		}
	})

}

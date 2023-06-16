package qrpix

import (
	"net/http"

	qrcode "github.com/skip2/go-qrcode"
)

const (
	imageSize = 256

	PIXGui                 = "br.gov.bcb.pix"
	PayloadFormatIndicator = "01"
)

type Static struct {
	Chave                string
	MerchantCategoryCode string
	TransactionCurrency  string
	CountryCode          string
	MerchantName         string
	MerchantCity         string
	PostalCode           string
	TransactionId        string
	// Transaction amount in cents
	TransactionAmount int

	builder Builder
}

type StaticOptFn func(*Static)

func NewStatic(chave, merchantName, merchantCity, txId string, fns ...StaticOptFn) *Static {
	qr := &Static{
		Chave:                chave,
		MerchantCategoryCode: "0000",
		TransactionCurrency:  "986",
		CountryCode:          "BR",
		MerchantName:         merchantName,
		MerchantCity:         merchantCity,
		TransactionId:        txId,
		builder:              Builder{},
	}
	for _, fn := range fns {
		fn(qr)
	}
	return qr
}

func WithTransactionAmount(value int) StaticOptFn {
	return func(s *Static) {
		s.TransactionAmount = value
	}
}

func WithPostalCode(postalCode string) StaticOptFn {
	return func(s *Static) {
		s.PostalCode = postalCode
	}
}

func WithCountryCode(code string) StaticOptFn {
	return func(s *Static) {
		s.CountryCode = code
	}
}

func WithMerchantCategoryCode(code string) StaticOptFn {
	return func(s *Static) {
		s.MerchantCategoryCode = code
	}
}

func (s *Static) BRCode() (string, error) {
	defer s.builder.Clear()

	s.builder.AddPayloadFormatIndicator(PayloadFormatIndicator)
	s.builder.AddMerchantAccountInformation(PIXGui, s.Chave)
	s.builder.AddMerchantCategoryCode(s.MerchantCategoryCode)
	s.builder.AddTransactionCurrency(s.TransactionCurrency)
	s.builder.AddTransactionAmount(s.TransactionAmount)
	s.builder.AddCountryCode(s.CountryCode)
	s.builder.AddMerchantName(s.MerchantName)
	s.builder.AddMerchantCity(s.MerchantCity)
	s.builder.AddPostalCode(s.PostalCode)
	s.builder.AddAdditionalDataField(s.TransactionId)

	return s.builder.Build()
}

func (s Static) SaveFile(path string) error {
	brCode, err := s.BRCode()
	if err != nil {
		return err
	}

	if err := qrcode.WriteFile(brCode, qrcode.Medium, imageSize, path); err != nil {
		return err
	}

	return nil
}

func (s Static) Encode() ([]byte, error) {
	brCode, err := s.BRCode()
	if err != nil {
		return nil, err
	}

	png, err := qrcode.Encode(brCode, qrcode.Medium, imageSize)
	if err != nil {
		return nil, err
	}

	return png, nil
}

func (s Static) Serve(w http.ResponseWriter) error {
	w.Header().Add("Content-Type", "image/png")
	png, err := s.Encode()
	if err != nil {
		return err
	}
	w.Write(png)
	return nil
}

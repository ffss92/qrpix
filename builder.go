package qrpix

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/snksoft/crc"
	"golang.org/x/exp/slices"
)

var (
	ErrInvalidCRC    = errors.New("crc is not valid")
	ErrCRCNotPresent = errors.New("crc is not present")
)

// Used to build the BRCode
type Builder map[string]TLV

func (b Builder) Add(tlv TLV) Builder {
	b[tlv.FieldID()] = tlv
	return b
}

func (b Builder) toSlice() []TLV {
	tlvs := []TLV{}
	for _, tlv := range b {
		tlvs = append(tlvs, tlv)
	}
	return tlvs
}

func (b Builder) Sorted() []TLV {
	tlvs := b.toSlice()
	slices.SortFunc[TLV](tlvs, func(a, b TLV) bool {
		aid, _ := strconv.Atoi(a.FieldID())
		bid, _ := strconv.Atoi(b.FieldID())
		return aid < bid
	})
	return tlvs
}

// Build and validates the BRCode. CRC16 is added automatically.
func (b Builder) Build() (string, error) {
	var (
		res  string
		tlvs = b.Sorted()
	)

	for _, tlv := range tlvs {
		id, length, value, err := tlv.TLV()
		if err != nil {
			return "", err
		}
		res += id + length + value
	}

	res = b.addCRC16(res)
	return res, nil
}

func (b Builder) buildRaw() (string, error) {
	var (
		res  string
		tlvs = b.Sorted()
	)

	for _, tlv := range tlvs {
		id, length, value, err := tlv.TLV()
		if err != nil {
			return "", err
		}
		res += id + length + value
	}

	return res, nil
}

// Clears the builder TLV items
func (b *Builder) Clear() {
	(*b) = Builder{}
}

func (b Builder) AddPayloadFormatIndicator(value string) {
	b.Add(Primitive{
		ID:    "00",
		Value: value,
	})
}

func (b Builder) AddMerchantAccountInformation(gui, chave string) {
	t := Template{
		ID: "26",
	}
	t.AddValue("00", gui)
	t.AddValue("01", chave)
	b.Add(t)
}

func (b Builder) AddMerchantCategoryCode(code string) {
	b.Add(Primitive{
		ID:    "52",
		Value: code,
	})
}

func (b Builder) AddTransactionCurrency(code string) {
	b.Add(Primitive{
		ID:    "53",
		Value: code,
	})
}

// Adds the transaction amount in cents. Ex: 100 == 1 real
func (b Builder) AddTransactionAmount(amount int) {
	if amount == 0 {
		return
	}

	b.Add(Primitive{
		ID:    "54",
		Value: fmt.Sprintf("%.2f", float64(amount)/float64(100)),
	})
}

func (b Builder) AddCountryCode(code string) {
	b.Add(Primitive{
		ID:    "58",
		Value: code,
	})
}

func (b Builder) AddMerchantName(name string) {
	b.Add(Primitive{
		ID:    "59",
		Value: name,
	})
}

func (b Builder) AddMerchantCity(city string) {
	b.Add(Primitive{
		ID:    "60",
		Value: city,
	})
}

func (b Builder) AddPostalCode(postalCode string) {
	b.Add(Primitive{
		ID:    "61",
		Value: postalCode,
	})
}

func (b Builder) AddAdditionalDataField(transactionId string) {
	t := Template{
		ID: "62",
	}
	t.AddValue("05", transactionId)
	b.Add(t)
}

func (b Builder) addCRC16(data string) string {
	appended := data + "6304"
	ccittCrc := crc.CalculateCRC(crc.CCITT, []byte(appended))
	return fmt.Sprintf("%s%04X", appended, ccittCrc)
}

// Removes crc and then builds without adding it and checks if hashes matches
func (b Builder) CheckCRC() error {
	crcField, ok := b["63"]
	if !ok {
		return ErrCRCNotPresent
	}
	_, _, value, err := crcField.TLV()
	if err != nil {
		return err
	}

	delete(b, "63") // Remove CRC from map

	r, err := b.buildRaw()
	if err != nil {
		return err
	}

	appended := r + "6304"
	ccittCrc := crc.CalculateCRC(crc.CCITT, []byte(appended))
	if fmt.Sprintf("%04X", ccittCrc) != value {
		return ErrInvalidCRC
	}

	return nil
}

package qrpix

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/snksoft/crc"
	"golang.org/x/exp/slices"
)

var (
	ErrInvalidCRC              = errors.New("crc is not valid")
	ErrCRCNotPresent           = errors.New("crc is not present")
	ErrRequiredFieldNotPresent = errors.New("required field not present")
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
	b.Add(&Primitive{
		ID:    "00",
		Value: value,
	})
}

func (b Builder) GetPayloadFormatIndicator() (string, error) {
	return b.GetPrimitiveField("00")
}

func (b Builder) AddMerchantAccountInformation(gui, chave string) {
	t := &Template{
		ID: "26",
	}
	t.AddValue("00", gui)
	t.AddValue("01", chave)
	b.Add(t)
}

func (b Builder) GetMerchantAccountInformationGui() (string, error) {
	return b.GetTemplateField("26", "00")
}

func (b Builder) GetMerchantAccountInformationChave() (string, error) {
	return b.GetTemplateField("26", "01")
}

func (b Builder) AddMerchantCategoryCode(code string) {
	b.Add(&Primitive{
		ID:    "52",
		Value: code,
	})
}

func (b Builder) GetMerchantCategoryCode() (string, error) {
	return b.GetPrimitiveField("52")
}

func (b Builder) AddTransactionCurrency(code string) {
	b.Add(&Primitive{
		ID:    "53",
		Value: code,
	})
}

func (b Builder) GetTransactionCurrency() (string, error) {
	return b.GetPrimitiveField("53")
}

// Adds the transaction amount in cents. Ex: 100 == 1 real
func (b Builder) AddTransactionAmount(amount int) {
	if amount == 0 {
		return
	}

	b.Add(&Primitive{
		ID:    "54",
		Value: fmt.Sprintf("%.2f", float64(amount)/float64(100)),
	})
}

func (b Builder) GetTransactionAmount() (int, error) {
	val, err := b.GetPrimitiveField("54")
	if err != nil {
		return 0, err
	}
	if val == "" {
		return 0, nil
	}
	f, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return 0, err
	}

	return int(f * float64(100)), nil
}

func (b Builder) AddCountryCode(code string) {
	b.Add(&Primitive{
		ID:    "58",
		Value: code,
	})
}

func (b Builder) GetCountryCode() (string, error) {
	return b.GetPrimitiveField("58")
}

func (b Builder) AddMerchantName(name string) {
	b.Add(&Primitive{
		ID:    "59",
		Value: name,
	})
}

func (b Builder) GetMerchanName() (string, error) {
	return b.GetPrimitiveField("59")
}

func (b Builder) AddMerchantCity(city string) {
	b.Add(&Primitive{
		ID:    "60",
		Value: city,
	})
}

func (b Builder) GetMerchantCity() (string, error) {
	return b.GetPrimitiveField("60")
}

func (b Builder) AddPostalCode(postalCode string) {
	b.Add(&Primitive{
		ID:    "61",
		Value: postalCode,
	})
}

func (b Builder) GetPostalCode() (string, error) {
	return b.GetPrimitiveField("61")
}

func (b Builder) AddAdditionalDataField(transactionId string) {
	t := &Template{
		ID: "62",
	}
	t.AddValue("05", transactionId)
	b.Add(t)
}

func (b Builder) GetTransactionId() (string, error) {
	return b.GetTemplateField("62", "05")
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

func (b Builder) GetPrimitiveField(id string) (string, error) {
	meta, err := GetFieldMetadata(id)
	if err != nil {
		return "", err
	}

	tlv, ok := b[id]
	if meta.Required && !ok {
		return "", fmt.Errorf("required field not present: %s", meta.Name)
	}

	if !ok {
		return "", nil
	}

	// Get value from TLV for further validation
	_, _, value, err := tlv.TLV()
	return value, err
}

func (b Builder) GetTemplateField(id string, fieldId string) (string, error) {
	tempMeta, err := GetFieldMetadata(id)
	if err != nil {
		return "", err
	}

	fieldMeta, err := GetFieldMetadata(id + "-" + fieldId)
	if err != nil {
		return "", err
	}

	template, ok := b[id]
	if !ok && tempMeta.Required {
		return "", fmt.Errorf("required field not present: %s", tempMeta.Name)
	}

	vals := template.Unwrap()
	if vals == nil && fieldMeta.Required {
		return "", fmt.Errorf("required field not present: %s", fieldMeta.Name)
	}

	field, ok := vals[fieldId]
	if !ok && fieldMeta.Required {
		return "", fmt.Errorf("required field not present: %s", fieldMeta.Name)
	}
	if !ok {
		return "", nil
	}

	// Get value from TLV for further validation
	_, _, value, err := field.TLV()
	if err != nil {
		return "", err
	}
	return value, nil
}

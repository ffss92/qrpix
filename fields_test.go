package qrpix

import (
	"errors"
	"testing"
)

func TestFields(t *testing.T) {
	t.Run("get field metadata should return error for unknown ids", func(t *testing.T) {
		_, err := GetFieldMetadata("invalid")
		if err == nil {
			t.Error("expected error for invalid metadata id but got nil")
		}
	})

	t.Run("validate field should return error for unknown ids", func(t *testing.T) {
		err := ValidateField("invalid", "")
		if !errors.Is(err, ErrFieldMetadataNotFound) {
			t.Errorf("expected ErrFieldMetadataNotFound but got: %v", err)
		}
	})

	t.Run("validate field should return error for required empty fields", func(t *testing.T) {
		err := ValidateField("00", "")
		if !errors.Is(err, ErrFieldIsRequired) {
			t.Errorf("expected ErrFieldIsRequired but got: %v", err)
		}
	})
}

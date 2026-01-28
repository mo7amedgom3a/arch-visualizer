package errors

import (
	"testing"
)

func TestAppError_Error(t *testing.T) {
	err := New("TEST_CODE", KindValidation, "test message")
	if err.Error() == "" {
		t.Error("Error() should not return empty string")
	}
}

func TestAppError_WithOp(t *testing.T) {
	err := New("TEST_CODE", KindValidation, "test message")
	err = err.WithOp("TestOperation")
	if err.Op != "TestOperation" {
		t.Errorf("Expected Op to be 'TestOperation', got %s", err.Op)
	}
}

func TestAppError_WithMeta(t *testing.T) {
	err := New("TEST_CODE", KindValidation, "test message")
	err = err.WithMeta("key", "value")
	if err.Metadata["key"] != "value" {
		t.Errorf("Expected metadata key 'key' to be 'value', got %v", err.Metadata["key"])
	}
}

func TestWrap(t *testing.T) {
	originalErr := New("ORIGINAL_CODE", KindInternal, "original error")
	wrapped := Wrap(originalErr, "WRAPPED_CODE", KindValidation, "wrapped message")
	
	if wrapped.Cause != originalErr {
		t.Error("Wrapped error should have original error as cause")
	}
	if wrapped.Code != "WRAPPED_CODE" {
		t.Errorf("Expected code 'WRAPPED_CODE', got %s", wrapped.Code)
	}
}

func TestIsKind(t *testing.T) {
	err := New("TEST_CODE", KindValidation, "test message")
	if !IsKind(err, KindValidation) {
		t.Error("IsKind should return true for matching kind")
	}
	if IsKind(err, KindNotFound) {
		t.Error("IsKind should return false for non-matching kind")
	}
}

func TestIsCode(t *testing.T) {
	err := New("TEST_CODE", KindValidation, "test message")
	if !IsCode(err, "TEST_CODE") {
		t.Error("IsCode should return true for matching code")
	}
	if IsCode(err, "OTHER_CODE") {
		t.Error("IsCode should return false for non-matching code")
	}
}

func TestAsAppError(t *testing.T) {
	err := New("TEST_CODE", KindValidation, "test message")
	appErr := AsAppError(err)
	if appErr == nil {
		t.Error("AsAppError should return non-nil for AppError")
	}
	if appErr.Code != "TEST_CODE" {
		t.Errorf("Expected code 'TEST_CODE', got %s", appErr.Code)
	}
}

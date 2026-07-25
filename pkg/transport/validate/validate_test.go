package validate

import "testing"

type sample struct {
	Email string `json:"email" validate:"required,email"`
	Name  string `json:"name" validate:"required,min=2"`
}

func TestValidatorStruct(t *testing.T) {
	v := New()
	if err := v.Struct(sample{Email: "a@b.co", Name: "Ал"}); err != nil {
		t.Fatalf("valid payload rejected: %v", err)
	}
	err := v.Struct(sample{Email: "not-an-email", Name: ""})
	if err == nil {
		t.Fatal("invalid payload should error")
	}
	verr, ok := err.(ValidationError)
	if !ok {
		t.Fatalf("want ValidationError, got %T", err)
	}
	if verr.Fields["email"] == "" || verr.Fields["name"] == "" {
		t.Errorf("both fields should report errors, got %v", verr.Fields)
	}
	if verr.Error() == "" {
		t.Error("ValidationError.Error should be non-empty")
	}
}

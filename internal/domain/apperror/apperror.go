package apperror

import "encoding/json"

type ErrorType string

var (
	INVALID    ErrorType = "invalid"
	MIN_LENGTH ErrorType = "minlength"
	MAX_LENGTH ErrorType = "maxlength"
	REQUIRED   ErrorType = "required"
	INTERNAL   ErrorType = "internal"
	EXISTS     ErrorType = "exists"
	MISMATCH   ErrorType = "mismatch"
	EXPIRED    ErrorType = "expired"
)

type ErrorTarget string

var (
	AUTH          ErrorTarget = "auth"
	EMAIL_CONFIRM ErrorTarget = "emailconfirm"
)

type AppError struct {
	Value  any         `json:"value"`
	Target ErrorTarget `json:"target"`
	Type   ErrorType   `json:"type"`
	Detail error       `json:"detail"`
}

func (ae *AppError) MarshalJSON() ([]byte, error) {
	type Alias AppError
	return json.Marshal(&struct {
		*Alias
		Detail string `json:"Detail"`
	}{
		Alias:  (*Alias)(ae),
		Detail: ae.Detail.Error(), // Convert error to string
	})
}

func NewAppError(value any, target ErrorTarget, t ErrorType, detail error) *AppError {
	return &AppError{
		Value:  value,
		Target: target,
		Type:   t,
		Detail: detail,
	}
}

type AppErrors []AppError

func (aes *AppErrors) Append(ae *AppError) {
	*aes = append(*aes, *ae)
}

func (aes *AppErrors) HasErrors() bool {
	return len(*aes) > 0
}

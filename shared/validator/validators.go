package validator

import validation "github.com/go-ozzo/ozzo-validation/v4"

type Validator struct {
	Rule Rule
}

func (v *Validator) ValidateResponse(error error) map[string]any {
	errors := make(map[string]any)

	validationErrors, ok := error.(validation.Errors)
	if ok {
		for field, fieldError := range validationErrors {
			customError, ok := fieldError.(validation.Error)
			if ok {
				errors[field] = customError.Error()
			} else {
				errors[field] = fieldError.Error()
			}
		}
	}

	return errors
}

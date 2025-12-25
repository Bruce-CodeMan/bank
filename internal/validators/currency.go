package validators

import (
	"github.com/go-playground/validator/v10"

	"github.com/BruceCompiler/bank/utils"
)

var ValidCurrency validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if currency, ok := fieldLevel.Field().Interface().(string); ok {
		// check currency is supported
		return utils.IsSupportedCurrency(currency)
	}
	return false
}

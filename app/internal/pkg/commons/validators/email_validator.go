package validators

import (
	"github.com/asaskevich/govalidator"
)

func setEmailValidator() {
	govalidator.CustomTypeTagMap.Set("email", func(i any, _ any) bool {
		email, ok := i.(string)
		if !ok {
			return false
		}
		return govalidator.IsEmail(email)
	})
}

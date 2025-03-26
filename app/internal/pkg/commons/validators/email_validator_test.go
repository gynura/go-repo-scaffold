package validators

import (
	"testing"

	"github.com/asaskevich/govalidator"
	"github.com/stretchr/testify/assert"
)

type EmailTestType struct {
	Email string `valid:"required,email"`
}

func TestEmailValidator(t *testing.T) {
	setEmailValidator()

	tests := []struct {
		name string
		have EmailTestType
		want bool
	}{
		{
			name: "empty",
			have: EmailTestType{},
			want: false,
		},
		{
			name: "missing prefix",
			have: EmailTestType{Email: "@gmail.com"},
			want: false,
		},
		{
			name: "missing suffix",
			have: EmailTestType{Email: "foo"},
			want: false,
		},
		{
			name: "correct email",
			have: EmailTestType{Email: "foo@gmail.com"},
			want: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ok, _ := govalidator.ValidateStruct(tt.have)
			assert.Equal(t, tt.want, ok)
		})
	}
}

package entities

import "golang.org/x/crypto/bcrypt"

// password type wraps a string to enforce hashing on the pass
type password string

// User defines the user model in the service.
type User struct {
	MetaInformation
	UserSpecificData
}

// UserSpecificData defines the attributes that are not part of the meta information of a User resource.
type UserSpecificData struct {
	FirstName string   `valid:"required"`
	LastName  string   `valid:"required"`
	Email     string   `valid:"required,email"`
	password  password `valid:"required"` // private field to prevent direct assignment from outside the package
	Country   string   `valid:"required"`
	Nickname  string   `valid:"required"`
}

// EncryptAndSetPassword encapsulates the password assignment so it is not possible to set a password that is not hashed
// from outside this package where the User resource logic ends up becoming more diffuse
func (u *User) EncryptAndSetPassword(rawPassword string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.password = password(hashed)
	return nil
}

// SetPassword sets the password without encryption.
// Should just be used when retrieving users from storage
func (u *User) SetPassword(pass string) error {
	u.password = password(pass)
	return nil
}

func (u *User) Password() string {
	return string(u.password)
}

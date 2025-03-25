package user

// UseCase defines the functionality to manage User resources
type UseCase interface {
	// TODO: complete
}

// UserUseCase defines the default implementation of the user.UseCase interface
type UserUseCase struct {
}

// ProvideUserUseCaseOptions defines the options to configure a UserUseCase implementation
type ProvideUserUseCaseOptions struct {
}

// ProvideUserUseCase creates a UserUseCase with the given options.
func ProvideUserUseCase(options ProvideUserUseCaseOptions) (*UserUseCase, error) {
	return &UserUseCase{}, nil
}

var _ UseCase = new(UserUseCase)

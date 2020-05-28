package transfermodels

type CreateUserRequest struct {
	Email string `json:"email" validate:"required,email,max=64"`
	// The password should be between 8 and 64 characters, containing lowercase letter, uppercase letter, special character, and a digit
	Password        string `json:"password" validate:"required,min=8,max=64,containsany=abcdefghijklmnopqrstuvwxyz,containsany=ABCDEFGHIJKLMNOPQRSTUVWXYZ,containsany=!@#$%^&*()_-+<>?,containsany=0123456789,eqfield=ConfirmPassword"`
	ConfirmPassword string `json:"confirm_password"`
	IsOwner         bool   `json:"is_owner"`
}

type LoginRequest struct {
	Email    string
	Password string
}

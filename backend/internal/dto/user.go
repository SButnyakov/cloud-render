package dto

type CreateUserDTO struct {
	Login    string
	Email    string
	Password string
}

type AuthUserDTO struct {
	LoginOrEmail string
	Password     string
	AccessToken  string
	RefreshToken string
}

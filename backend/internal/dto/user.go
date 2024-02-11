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

type ReAuthUserDTO struct {
	AccessToken  string
	RefreshToken string
}

type EditUserDTO struct {
	Id       int64
	Login    string
	Email    string
	Password string
}

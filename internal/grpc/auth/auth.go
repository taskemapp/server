package auth

type Auth interface {
	Login(LoginOpts) (resp LoginResponse, err error)
}

type LoginOpts struct {
	email    string
	password string
}

type LoginResponse struct {
	Token        string
	RefreshToken string
}

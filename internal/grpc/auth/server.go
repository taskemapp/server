package auth

import (
	"context"
	"errors"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"math"
	"regexp"
	"taskem-server/internal/repositories/user"
	"taskem-server/internal/service/auth"
	"taskem-server/tools/gen/grpc/v1"
	"unicode"
)

type Opts struct {
	fx.In
	Auth   auth.Service
	Logger *zap.Logger
}

type Server struct {
	v1.UnimplementedAuthServer
	auth   auth.Service
	logger *zap.Logger
}

func New(opts Opts) *Server {
	return &Server{
		auth:   opts.Auth,
		logger: opts.Logger,
	}
}

func (s *Server) Login(
	ctx context.Context,
	req *v1.LoginRequest,
) (*v1.LoginResponse, error) {
	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "Missing argument: email")
	}

	if req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "Missing argument: password")
	}

	resp, err := s.auth.Login(
		ctx,
		auth.LoginOpts{
			Email:    req.Email,
			Password: req.Password,
		})

	if err != nil {
		s.logger.Sugar().Error(err)
		switch {
		case errors.Is(err, user.ErrNotFound):
			return nil, status.Error(codes.NotFound, "Not found")
		case errors.Is(err, auth.ErrTokenGen):
			return nil, status.Error(codes.Internal, "Token generate failed")
		case errors.Is(err, auth.ErrPwdMatch):
			return nil, status.Error(codes.InvalidArgument, "Wrong password")
		}
		return nil, err
	}

	return &v1.LoginResponse{
		Token:        resp.Token,
		RefreshToken: resp.RefreshToken,
		TokenType:    resp.TokenType,
	}, nil
}

const (
	emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
)

func (s *Server) SignUp(
	ctx context.Context,
	req *v1.SignupRequest,
) (*emptypb.Empty, error) {
	if req.UserName == "" {
		return nil, status.Error(codes.InvalidArgument, "Missing argument: user_name")
	}

	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "Missing argument: email")
	}

	if req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "Missing argument: password")
	}

	if !regexp.MustCompile(emailRegex).MatchString(req.Email) {
		return nil, status.Error(codes.InvalidArgument, "Invalid email: use format example@example.com")
	}

	err := PasswordComplexity(req.Password)
	if err != nil {
		return nil, err
	}

	err = s.auth.Registration(
		ctx,
		auth.RegistrationOpts{
			Email:    req.Email,
			Name:     req.UserName,
			Password: req.Password,
		})
	if err != nil {
		s.logger.Sugar().Error(err)
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func PasswordComplexity(password string) error {
	isDigit := false
	isUpper := false
	isLower := false
	isSpecial := false

	for _, c := range password {
		if unicode.IsDigit(c) {
			isDigit = true
		}
		if unicode.IsUpper(c) {
			isUpper = true
		}
		if unicode.IsLower(c) {
			isLower = true
		}
		if unicode.IsPunct(c) || unicode.IsSymbol(c) {
			isSpecial = true
		}
	}

	var symbolPool int
	if isLower && isUpper && isDigit && isSpecial {
		symbolPool = 95 // contains (a-z, A-Z, ASCII, space)
	} else if isLower && isUpper && isDigit {
		symbolPool = 62 // contains (a-z, A-Z, 0-9)
	} else if isLower && isDigit {
		symbolPool = 36 // contains (a-z, 0-9)
	} else {
		symbolPool = 26 // contains (a-z)
	}

	passwordComplexity := math.Log2(float64(symbolPool)) * float64(len(password))

	const minComplexity = 40.0
	if passwordComplexity < minComplexity {
		return status.Error(codes.InvalidArgument, "Password is too weak")
	}

	return nil
}

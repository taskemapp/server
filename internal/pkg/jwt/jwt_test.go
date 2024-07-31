package jwt

import (
	"github.com/brianvoe/gofakeit/v7"
	jwt2 "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestNewToken(t *testing.T) {
	tests := []struct {
		name        string
		opts        Opts
		failed      bool
		expectedErr error
	}{
		{
			name: "Very short life duration",
			opts: Opts{
				Email:    gofakeit.Email(),
				Secret:   "secret",
				Duration: time.Second * 10,
				ID:       uuid.New(),
			},
			failed:      true,
			expectedErr: ErrDuration,
		},
		{
			name: "Short life duration",
			opts: Opts{
				Email:    gofakeit.Email(),
				Secret:   "secret",
				Duration: time.Minute,
				ID:       uuid.New(),
			},
			failed: false,
		},
		{
			name: "Long life duration",
			opts: Opts{
				Email:    gofakeit.Email(),
				Secret:   "secret",
				Duration: time.Hour,
			},
			failed: false,
		},
		{
			name: "Missing secret",
			opts: Opts{
				Email:    gofakeit.Email(),
				Secret:   "",
				Duration: time.Hour,
			},
			failed:      true,
			expectedErr: ErrMissingSecret,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			token, err := NewToken(test.opts)

			if !test.failed {
				require.NoError(t, err)
				parsed, err := jwt2.Parse(token, func(token *jwt2.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt2.SigningMethodHMAC); !ok {
						t.Fatalf("unexpected signing method: %v", token.Header["alg"])
					}

					return []byte(test.opts.Secret), nil
				})
				if claims, ok := parsed.Claims.(jwt2.MapClaims); ok {
					require.Equal(t, test.opts.ID.String(), claims["uid"])
					require.Equal(t, test.opts.Email, claims["email"])
				} else {
					t.Fatalf("failed get claims from token: %s", token)
				}
				require.NoError(t, err)
			} else {
				if test.expectedErr != nil {
					require.ErrorIs(t, err, test.expectedErr)
				}
			}
		})
	}
}
func TestGetPayload(t *testing.T) {
	tests := []struct {
		name        string
		token       string
		secret      string
		expectedErr error
	}{
		{
			name:        "InvalidToken",
			token:       "invalid.token.string",
			secret:      "secret",
			expectedErr: ErrTokenParse,
		},
		{
			name:        "ValidToken",
			token:       "",
			secret:      "secret",
			expectedErr: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.name == "ValidToken" {
				opts := Opts{
					Email:    gofakeit.Email(),
					Secret:   test.secret,
					Duration: time.Minute,
					ID:       uuid.New(),
				}
				token, err := NewToken(opts)
				require.NoError(t, err)
				test.token = token
			}

			payload, err := GetPayload(test.token, test.secret)
			if test.expectedErr != nil {
				require.ErrorIs(t, err, test.expectedErr)
				require.Nil(t, payload)
			} else {
				require.NoError(t, err)
				require.NotNil(t, payload)
			}
		})
	}
}

func TestValidate(t *testing.T) {
	token, err := NewToken(Opts{
		Email:    gofakeit.Email(),
		Secret:   "secret",
		Duration: time.Minute,
		ID:       uuid.New(),
	})
	require.NoError(t, err)
	if err != nil {
		return
	}
	tests := []struct {
		name        string
		token       string
		secret      string
		expectedErr error
	}{
		{
			name:        "InvalidToken",
			token:       "invalid token string",
			secret:      "secret",
			expectedErr: ErrTokenParse,
		},
		{
			name:        "ValidToken",
			token:       token,
			secret:      "secret",
			expectedErr: nil,
		},
		{
			name:        "InvalidSecret",
			token:       token,
			secret:      "invalid_secret",
			expectedErr: ErrTokenValidation,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			valid, err := Validate(test.token, test.secret)
			if test.expectedErr != nil {
				require.ErrorIs(t, err, test.expectedErr)
				require.Nil(t, valid)
			} else {
				require.NoError(t, err)
				require.NotNil(t, valid)
			}
		})
	}
}

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
			name: "",
			opts: Opts{
				Email:    gofakeit.Email(),
				Secret:   "secret",
				Duration: time.Second * 10,
				Id:       uuid.New(),
			},
			failed:      true,
			expectedErr: ErrDuration,
		},
		{
			name: "",
			opts: Opts{
				Email:    gofakeit.Email(),
				Secret:   "secret",
				Duration: time.Minute,
				Id:       uuid.New(),
			},
			failed: false,
		},
		{
			name: "",
			opts: Opts{
				Email:    gofakeit.Email(),
				Secret:   "secret",
				Duration: time.Hour,
			},
			failed: false,
		},
		{
			name: "",
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
				require.Nil(t, err)
				parsed, err := jwt2.Parse(token, func(token *jwt2.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt2.SigningMethodHMAC); !ok {
						t.Fatalf("unexpected signing method: %v", token.Header["alg"])
					}

					return []byte(test.opts.Secret), nil
				})
				if claims, ok := parsed.Claims.(jwt2.MapClaims); ok {
					require.Equal(t, test.opts.Id.String(), claims["uid"])
					require.Equal(t, test.opts.Email, claims["email"])
				} else {
					t.Fatalf("failed get claims from token: %s", token)
				}
				require.Nil(t, err)
			} else {
				if test.expectedErr != nil {
					require.ErrorIs(t, err, test.expectedErr)
				}
			}
		})
	}
}

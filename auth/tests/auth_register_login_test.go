package tests

import (
	"fmt"
	"github.com/074yara/AuthGrpc/auth/tests/suite"
	"github.com/074yara/AuthGrpc/protos/gen/authGrpc"
	"github.com/brianvoe/gofakeit"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

const (
	emptyAppId          = 0
	appId          uint = 1
	appSecret           = "test-secret"
	passDefaultLen      = 10
)

func TestRegisterLogin_HappyPath(t *testing.T) {
	ctx, s := suite.New(t)
	email, password := gofakeit.Email(), randomPassword()
	respRegister, err := s.AuthClient.Register(ctx, &authGrpc.RegisterRequest{
		Email:    email,
		Password: password},
	)
	require.NoError(t, err)
	require.NotEmpty(t, t, respRegister.GetUserId())

	respLogin, err := s.AuthClient.Login(ctx, &authGrpc.LoginRequest{
		Email:    email,
		Password: password,
		AppId:    uint64(appId),
	})
	require.NoError(t, err)

	loginTime := time.Now()

	tokenParsed, err := jwt.Parse(respLogin.Token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(appSecret), nil
	})
	require.NoError(t, err)

	claims, ok := tokenParsed.Claims.(jwt.MapClaims)
	assert.True(t, ok)

	assert.Equal(t, respRegister.GetUserId(), uint64(claims["uid"].(float64)))
	assert.Equal(t, email, claims["email"].(string))
	assert.Equal(t, appId, uint(claims["app_id"].(float64)))

	const deltaSeconds = 1
	assert.InDelta(t, loginTime.Add(s.Cfg.TokenTTl).Unix(), claims["exp"].(float64), deltaSeconds)

}

func TestRegisterWrongEmailPassword(t *testing.T) {
	ctx, s := suite.New(t)
	email, password := gofakeit.Email(), randomPassword()
	_, err := s.AuthClient.Register(ctx, &authGrpc.RegisterRequest{
		Email:    email,
		Password: password,
	})
	require.NoError(t, err)
	_, err = s.AuthClient.Login(ctx, &authGrpc.LoginRequest{
		Email:    email,
		Password: "",
		AppId:    1,
	})
	require.Error(t, err)
	_, err = s.AuthClient.Login(ctx, &authGrpc.LoginRequest{
		Email:    "",
		Password: password,
		AppId:    1,
	})
	require.Error(t, err)
	_, err = s.AuthClient.Login(ctx, &authGrpc.LoginRequest{
		Email:    "wrong-email",
		Password: password,
		AppId:    1,
	})
	require.Error(t, err)
	_, err = s.AuthClient.Login(ctx, &authGrpc.LoginRequest{
		Email:    email,
		Password: "wrong-password",
		AppId:    1,
	})
	require.Error(t, err)
}

func TestRegistrationTwice(t *testing.T) {
	ctx, s := suite.New(t)
	email, password := gofakeit.Email(), randomPassword()
	_, err := s.AuthClient.Register(ctx, &authGrpc.RegisterRequest{
		Email:    email,
		Password: password,
	})
	require.NoError(t, err)
	_, err = s.AuthClient.Register(ctx, &authGrpc.RegisterRequest{
		Email:    email,
		Password: password,
	})
	require.Error(t, err)
}

func TestRegistrationNoPassNoEmail(t *testing.T) {
	ctx, s := suite.New(t)
	email, password := gofakeit.Email(), randomPassword()
	_, err := s.AuthClient.Register(ctx, &authGrpc.RegisterRequest{
		Email:    email,
		Password: "",
	})
	require.Error(t, err)
	_, err = s.AuthClient.Register(ctx, &authGrpc.RegisterRequest{
		Email:    "",
		Password: password,
	})
	require.Error(t, err)
}

func randomPassword() string {
	return gofakeit.Password(true, true, true, true, false, passDefaultLen)
}

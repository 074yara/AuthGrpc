package suite

import (
	"context"
	"fmt"
	"github.com/074yara/AuthGrpc/auth/internal/config"
	"github.com/074yara/AuthGrpc/protos/gen/authGrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"testing"
)

type Suite struct {
	*testing.T
	Cfg        *config.Config
	AuthClient authGrpc.AuthClient
}

func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()
	cfg, err := config.LoadByPath("../config/local-tests.yml")
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), cfg.GRPC.Timeout)

	t.Cleanup(func() {
		t.Helper()
		cancel()
	})

	conn, err := grpc.Dial(fmt.Sprintf("localhost:%v", cfg.GRPC.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	client := authGrpc.NewAuthClient(conn)

	return ctx, &Suite{T: t, Cfg: cfg, AuthClient: client}

}

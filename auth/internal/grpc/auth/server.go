package auth

import "github.com/074yara/AuthGrpc/protos/gen/authGrpc"

type ServerAPI struct {
	authGrpc.UnimplementedAuthServer
}

package jwtauth

import (
	"context"

	"github.com/filecoin-project/go-jsonrpc/auth"
	"github.com/gbrlsnchs/jwt/v3"
	logging "github.com/ipfs/go-log/v2"
	"golang.org/x/xerrors"
)

var log = logging.Logger("jwt-auth")

type jwtPayload struct {
	Allow []auth.Permission
}

type IAuthWay interface {
	AuthVerify(context.Context, string) ([]auth.Permission, error)
	AuthNew(context.Context, []auth.Permission) ([]byte, error)
}

type JwtAuth struct {
	alg *jwt.HMACSHA
}

func NewJwtAuth(s string) *JwtAuth {
	return &JwtAuth{
		alg: jwt.NewHS256([]byte(s)),
	}
}

// AuthVerify resolve allow permission
func (j *JwtAuth) AuthVerify(ctx context.Context, info string) ([]auth.Permission, error) {
	var p jwtPayload
	log.Info("222")
	if _, err := jwt.Verify([]byte(info), j.alg, &p); err != nil {
		return nil, xerrors.Errorf("JWT Verification failed: %w", err)
	}
	log.Info("hhh")
	log.Info(p.Allow)
	return p.Allow, nil
}

func (j *JwtAuth) AuthNew(ctx context.Context, perms []auth.Permission) ([]byte, error) {
	p := jwtPayload{
		Allow: perms,
	}

	return jwt.Sign(&p, j.alg)
}

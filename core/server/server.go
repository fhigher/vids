package server

import (
	"context"

	"github.com/fhigher/vids/config"
	logging "github.com/ipfs/go-log/v2"
)

var log = logging.Logger("server-impl")

type VidsServer struct {
	Ver string
}

func NewVidsServer(c *config.Config) VidsServerAPI {
	return &VidsServer{
		Ver: "0.0.1",
	}
}

func (vids *VidsServer) Version(ctx context.Context) string {
	return vids.Ver
}

var _ VidsServerAPI = new(VidsServer)

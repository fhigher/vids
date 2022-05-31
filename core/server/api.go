package server

import "context"

type VidsServerAPI interface {
	Version(context.Context) string
}

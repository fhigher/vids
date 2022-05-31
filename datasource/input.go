package input

import "context"



type Input interface {
	Start(context.Context)
}

package api

import (
	"context"
)

type TestService interface {
	Hello(context.Context, string) (string, error)
}

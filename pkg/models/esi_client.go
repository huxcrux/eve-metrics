package models

import (
	"context"

	"github.com/antihax/goesi"
)

type ESIClient struct {
	Ctx    context.Context
	Client *goesi.APIClient
}

package services

import (
	"context"

	"github.com/hblab-ngocnd/get-started/models"
)

type Visitor interface {
	CreateVisitor(ctx context.Context, v models.Visitor) error
}

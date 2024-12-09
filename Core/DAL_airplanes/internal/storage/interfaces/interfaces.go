package interfaces

import "context"

type IAirplanesRepository interface {
	GetAll(context.Context) (any, error)
	Insert(context.Context, any) (any, error)
	Update(context.Context, any) (any, error)
	Delete(context.Context, any) (any, error)
}
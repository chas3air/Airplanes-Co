package interfaces

import "context"

type IPlanesRepository interface {
	GetAll(context.Context) (any, error)
	GetById(context.Context, any) (any, error)
	Insert(context.Context, any) (any, error)
	Update(context.Context, any) (any, error)
	Delete(context.Context, any) (any, error)
}

package intefaces

import "context"

type IFlightsRepository interface {
	GetAll(context.Context) (any, error)
	GetById(context.Context, int) (any, error)
	Insert(context.Context, any) (any, error)
	Update(context.Context, any) (any, error)
	Delete(context.Context, int) (any, error)
}

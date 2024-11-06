package intefaces

import (
	"context"
)

type ICustomersRepository interface {
	GetAll(context.Context) (any, error)
	GetById(context.Context, any) (any, error)
	GetByLoginAndPassword(context.Context, string, string) (any, error)
	Insert(context.Context, any) (any, error)
	Update(context.Context, any) (any, error)
	Delete(context.Context, any) (any, error)
}

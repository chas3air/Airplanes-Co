package storage

import "context"

type DBRepository interface {
	GetAll(context.Context) (any, error)
	GetById(context.Context, any) (any, error)
	Insert(context.Context, any) error
	Update(context.Context, any) error
	Delete(context.Context, any) error
}

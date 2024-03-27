package conference

import (
	"context"
	paginator "github.com/dmitryburov/gorm-paginator"
	"gorm.io/gorm"
)

type Paged[T any] struct {
	Page    *paginator.Pagination
	Content []T
}

type Repository interface {
	FindByID(ctx context.Context, id int) (*Conference, error)
	FindPaged(ctx context.Context, req *paginator.Paging) (*Paged[Conference], error)
	Create(ctx context.Context, entity *Conference) error
	Update(ctx context.Context, entity *Conference) error
	Delete(ctx context.Context, id int) error
}

func NewRepository(conn *gorm.DB) Repository {
	return repository{conn: conn}
}

type repository struct {
	conn *gorm.DB
}

func (r repository) Delete(ctx context.Context, id int) error {
	return r.conn.WithContext(ctx).Delete(&Conference{}, id).Error
}

func (r repository) FindPaged(ctx context.Context, req *paginator.Paging) (*Paged[Conference], error) {
	var out []Conference
	p, err := paginator.Pages(&paginator.Param{
		DB:     r.conn.WithContext(ctx),
		Paging: req,
	}, &out)
	if err != nil {
		return nil, err
	}
	return &Paged[Conference]{
		Page:    p,
		Content: out,
	}, nil
}

func (r repository) FindByID(ctx context.Context, id int) (*Conference, error) {
	var conference Conference
	if err := r.conn.WithContext(ctx).First(&conference, id).Error; err != nil {
		return nil, err
	}
	return &conference, nil
}

func (r repository) Create(ctx context.Context, entity *Conference) error {
	return r.conn.WithContext(ctx).Create(entity).Error
}

func (r repository) Update(ctx context.Context, entity *Conference) error {
	return r.conn.WithContext(ctx).Save(entity).Error
}

func (r repository) Count(ctx context.Context) (int64, error) {
	var count int64
	return count, r.conn.WithContext(ctx).Table("conferences").Count(&count).Error
}

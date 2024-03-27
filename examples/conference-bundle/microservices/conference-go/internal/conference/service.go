package conference

import (
	"context"
	paginator "github.com/dmitryburov/gorm-paginator"
)

type Service interface {
	FindByID(ctx context.Context, id int) (*Conference, error)
	FindPaged(ctx context.Context, req *paginator.Paging) (*Paged[Conference], error)
	Create(ctx context.Context, req CreateRequest) (*Conference, error)
	Update(ctx context.Context, req Conference) (*Conference, error)
	Delete(ctx context.Context, id int) error
}

func NewService(repository Repository) Service {
	return service{
		repository: repository,
	}
}

type service struct {
	repository Repository
}

func (s service) FindByID(ctx context.Context, id int) (*Conference, error) {
	return s.repository.FindByID(ctx, id)
}

func (s service) FindPaged(ctx context.Context, req *paginator.Paging) (*Paged[Conference], error) {
	return s.repository.FindPaged(ctx, req)
}

func (s service) Create(ctx context.Context, req CreateRequest) (*Conference, error) {
	conference := &Conference{
		Name:     req.Name,
		Location: req.Location,
	}
	if err := s.repository.Create(ctx, conference); err != nil {
		return nil, err
	}
	return conference, nil
}

func (s service) Update(ctx context.Context, req Conference) (*Conference, error) {
	if _, err := s.FindByID(ctx, req.ID); err != nil {
		return nil, err
	}
	if err := s.repository.Update(ctx, &req); err != nil {
		return nil, err
	}
	return &req, nil
}

func (s service) Delete(ctx context.Context, id int) error {
	return s.repository.Delete(ctx, id)
}

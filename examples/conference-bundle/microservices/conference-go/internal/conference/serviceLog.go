package conference

import (
	"context"
	paginator "github.com/dmitryburov/gorm-paginator"
	"log/slog"
	"time"
)

func NewServiceLog(logger *slog.Logger, next Service) Service {
	return serviceLog{logger: logger.With("service", "conference-service"), next: next}
}

type serviceLog struct {
	logger *slog.Logger
	next   Service
}

func (s serviceLog) FindByID(ctx context.Context, id int) (out *Conference, err error) {
	defer func(begin time.Time) {
		s.logger.Debug(
			"FindByID",
			"input", id,
			"err", err,
			"took", time.Since(begin),
			"out", out,
		)
	}(time.Now())
	out, err = s.next.FindByID(ctx, id)
	return
}

func (s serviceLog) FindPaged(ctx context.Context, req *paginator.Paging) (out *Paged[Conference], err error) {
	defer func(begin time.Time) {
		s.logger.Debug(
			"FindPaged",
			"input", req,
			"err", err,
			"took", time.Since(begin),
			"out", out,
		)
	}(time.Now())
	out, err = s.next.FindPaged(ctx, req)
	return
}

func (s serviceLog) Create(ctx context.Context, req CreateRequest) (out *Conference, err error) {
	defer func(begin time.Time) {
		s.logger.Debug(
			"FindPaged",
			"input", req,
			"err", err,
			"took", time.Since(begin),
			"out", out,
		)
	}(time.Now())
	out, err = s.next.Create(ctx, req)
	return
}

func (s serviceLog) Update(ctx context.Context, req Conference) (conference *Conference, err error) {
	defer func(begin time.Time) {
		s.logger.Debug(
			"Update",
			"input", req,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	conference, err = s.next.Update(ctx, req)
	return
}

func (s serviceLog) Delete(ctx context.Context, id int) (err error) {
	defer func(begin time.Time) {
		s.logger.Debug(
			"Delete",
			"input", id,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	err = s.next.Delete(ctx, id)
	return
}

package conference

import (
	"context"
	paginator "github.com/dmitryburov/gorm-paginator"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

func NewServiceTracer(next Service) Service {
	tr := otel.GetTracerProvider().Tracer("conference-service")

	return serviceTracer{
		next: next,
		tr:   tr,
	}
}

type serviceTracer struct {
	next Service
	tr   trace.Tracer
}

func (s serviceTracer) FindByID(ctx context.Context, id int) (out *Conference, err error) {
	_, span := s.tr.Start(ctx, "FindByID")
	defer span.End()
	out, err = s.next.FindByID(ctx, id)
	return
}

func (s serviceTracer) FindPaged(ctx context.Context, req *paginator.Paging) (out *Paged[Conference], err error) {
	_, span := s.tr.Start(ctx, "FindPaged")
	defer span.End()
	out, err = s.next.FindPaged(ctx, req)
	return
}

func (s serviceTracer) Create(ctx context.Context, req CreateRequest) (out *Conference, err error) {
	_, span := s.tr.Start(ctx, "Create")
	defer span.End()
	out, err = s.next.Create(ctx, req)
	return
}

func (s serviceTracer) Update(ctx context.Context, req Conference) (conference *Conference, err error) {
	_, span := s.tr.Start(ctx, "Update")
	defer span.End()
	conference, err = s.next.Update(ctx, req)
	return
}

func (s serviceTracer) Delete(ctx context.Context, id int) (err error) {
	_, span := s.tr.Start(ctx, "Delete")
	defer span.End()
	err = s.next.Delete(ctx, id)
	return
}

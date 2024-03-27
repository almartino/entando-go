package conference

import (
	"encoding/json"
	"fmt"
	"github.com/debyten/apierr"
	"github.com/debyten/httplayer"
	paginator "github.com/dmitryburov/gorm-paginator"
	"net/http"
	"strconv"
)

func NewApi(service Service) httplayer.Routing {
	return api{
		service: service,
	}
}

type api struct {
	service Service
}

func (a api) Routes(with *httplayer.RoutingDefinition) []httplayer.Route {
	return with.
		Add(http.MethodGet, "/api/conferences", a.listConferences).
		Add(http.MethodPost, "/api/conferences", a.createConferences).
		Add(http.MethodGet, "/api/conferences/{id}", a.findConference).
		Add(http.MethodPut, "/api/conferences/{id}", a.updateConference).
		Add(http.MethodPatch, "/api/conferences/{id}", a.partialUpdateConference).
		Add(http.MethodDelete, "/api/conferences/{id}", a.deleteConference).
		Done()
}

func (a api) listConferences(w http.ResponseWriter, r *http.Request) {
	page, err := parsePage(r)
	if err != nil {
		apierr.HandleISE(err, w, r)
		return
	}
	found, err := a.service.FindPaged(r.Context(), page)
	if err != nil {
		apierr.HandleISE(err, w, r)
		return
	}
	writePageHeaders(w, found)
	_ = json.NewEncoder(w).Encode(&found.Content)
}

func (a api) createConferences(w http.ResponseWriter, r *http.Request) {
	var req CreateRequest
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apierr.HandleISE(err, w, r)
		return
	}
	conference, err := a.service.Create(r.Context(), req)
	if err != nil {
		apierr.HandleISE(err, w, r)
		return
	}
	_ = json.NewEncoder(w).Encode(conference)
}

func (a api) findConference(w http.ResponseWriter, r *http.Request) {
	id, err := parseId(r)
	if err != nil {
		apierr.HandleISE(err, w, r)
		return
	}
	conference, err := a.service.FindByID(r.Context(), id)
	if err != nil {
		apierr.HandleISE(err, w, r)
		return
	}
	_ = json.NewEncoder(w).Encode(conference)
}

func (a api) updateConference(w http.ResponseWriter, r *http.Request) {
	id, err := parseId(r)
	if err != nil {
		apierr.HandleISE(err, w, r)
		return
	}
	var req Conference
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apierr.HandleISE(err, w, r)
		return
	}
	req.ID = id
	out, err := a.service.Update(r.Context(), req)
	if err != nil {
		apierr.HandleISE(err, w, r)
		return
	}
	_ = json.NewEncoder(w).Encode(out)
}

func (a api) partialUpdateConference(w http.ResponseWriter, r *http.Request) {
	a.updateConference(w, r)
}

func (a api) deleteConference(w http.ResponseWriter, r *http.Request) {
	id, err := parseId(r)
	if err != nil {
		apierr.HandleISE(err, w, r)
		return
	}
	if err := a.service.Delete(r.Context(), id); err != nil {
		apierr.HandleISE(err, w, r)
		return
	}
}

func parseId(r *http.Request) (int, error) {
	id := r.PathValue("id")
	v, err := strconv.Atoi(id)
	if err != nil {
		return -1, apierr.BadRequest.Err(err)
	}
	return v, nil
}

func parsePage(r *http.Request) (*paginator.Paging, error) {
	req := &PageRequest{}
	q := r.URL.Query()
	page := q.Get("page")
	size := q.Get("size")
	sort := q.Get("sort")
	if err := req.ParsePage(page); err != nil {
		return nil, err
	}
	if err := req.ParseSort(sort); err != nil {
		return nil, err
	}
	if err := req.ParseSize(size); err != nil {
		return nil, err
	}
	return &paginator.Paging{
		Page:    req.Page,
		OrderBy: req.Sort,
		Limit:   req.Size,
		ShowSQL: false,
	}, nil
}

func writePageHeaders[T any](w http.ResponseWriter, resp *Paged[T]) {
	pg := resp.Page
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Total-Count", fmt.Sprint(pg.TotalRecords))
}

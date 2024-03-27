package conference

import (
	"github.com/debyten/apierr"
	"strconv"
	"strings"
)

type CreateRequest struct {
	Name     string `json:"name"`
	Location string `json:"location"`
}

type PageRequest struct {
	Page int      `json:"page,omitempty"`
	Size int      `json:"size,omitempty"`
	Sort []string `json:"sort,omitempty"`
}

func (r *PageRequest) ParsePage(page string) error {
	if page == "" {
		return nil
	}
	p, err := strconv.Atoi(page)
	if err != nil {
		return apierr.BadRequest.Err(err)
	}
	r.Page = p
	return nil
}

func (r *PageRequest) ParseSize(size string) error {
	if size == "" {
		r.Size = 10
		return nil
	}
	s, err := strconv.Atoi(size)
	if err != nil {
		return apierr.BadRequest.Err(err)
	}
	r.Size = s
	return nil
}

func (r *PageRequest) ParseSort(sort string) error {
	if sort == "" {
		r.Sort = []string{"id"}
		return nil
	}
	r.Sort = strings.Split(sort, ",")
	return nil
}

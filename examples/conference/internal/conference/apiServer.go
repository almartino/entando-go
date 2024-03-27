package conference

import (
	"github.com/debyten/httplayer"
	"gorm.io/gorm"
	"log/slog"
)

func NewServer(logger *slog.Logger, db *gorm.DB) httplayer.Routing {
	repository := NewRepository(db)
	service := NewService(repository)
	service = NewServiceLog(logger, service)
	return NewApi(NewServiceTracer(service))
}

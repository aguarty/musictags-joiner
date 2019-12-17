package artists

import (
	"musictags-joiner/internal/utils"
	"musictags-joiner/pkgs/logger"
	"musictags-joiner/pkgs/storage"
	"net/http"

	jsoniter "github.com/json-iterator/go"
)

var (
	json = jsoniter.ConfigFastest
)

// Service -
type Service struct {
	stor         *storage.Storage
	logger       *logger.Logger
	stortagsPath string
}

// NewService -
func NewService(stor *storage.Storage, logger *logger.Logger, stortagsPath string) *Service {
	return &Service{stor, logger, stortagsPath}
}

func (s *Service) artistsList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		utils.SendResponse(s.logger, w, http.StatusOK, "LIST")
	}
}

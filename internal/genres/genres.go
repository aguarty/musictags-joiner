package genres

import (
	"musictags-joiner/internal/utils"
	"musictags-joiner/pkgs/logger"
	"musictags-joiner/pkgs/storage"
	"net/http"
	"os"

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

func (s *Service) joiningtags() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := make([]string, 3)
		err := json.NewDecoder(r.Body).Decode(&params)
		if err != nil {
			s.logger.Error(err.Error())
			utils.SendResponse(s.logger, w, http.StatusBadRequest, "Invalid request parameters")
			return
		}

		result := make([]storage.ResGetTopArtists, 0)

		if result, err = s.stor.Joining(params); err != nil {
			s.logger.Error(err.Error())
			utils.SendResponse(s.logger, w, http.StatusInternalServerError, "Couldn't done this operation")
			return
		}

		utils.SendResponse(s.logger, w, http.StatusOK, result)
	}
}

func (s *Service) genresList(stortagsPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		result := make([]string, 0)

		f, err := os.Open(stortagsPath)
		if err != nil {
			s.logger.Error(err.Error())
			utils.SendResponse(s.logger, w, http.StatusInternalServerError, "Couldn't done this operation")
			return
		}
		files, err := f.Readdir(-1)
		f.Close()
		if err != nil {
			s.logger.Error(err.Error())
			utils.SendResponse(s.logger, w, http.StatusInternalServerError, "Couldn't done this operation")
			return
		}

		for _, file := range files {
			result = append(result, file.Name())
		}

		utils.SendResponse(s.logger, w, http.StatusOK, result)
	}
}

package main

import (
	"net/http"
	"os"
)

func (a *application) joiningtags(w http.ResponseWriter, r *http.Request) {

	params := make([]string, 3)
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		a.logger.Error(err.Error())
		a.sendResponse(w, http.StatusBadRequest, "Invalid request parameters")
		return
	}

	result := make([]resGetTopArtists, 0)

	if err := a.joining(
		params,
		&result,
	); err != nil {
		a.logger.Error(err.Error())
		a.sendResponse(w, http.StatusInternalServerError, "Couldn't done this operation")
		return
	}

	a.sendResponse(w, http.StatusOK, result)

}

func (a *application) genresList(w http.ResponseWriter, r *http.Request) {

	result := make([]string, 0)

	f, err := os.Open(StortagsPath)
	if err != nil {
		a.logger.Error(err.Error())
		a.sendResponse(w, http.StatusInternalServerError, "Couldn't done this operation")
		return
	}
	files, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		a.logger.Error(err.Error())
		a.sendResponse(w, http.StatusInternalServerError, "Couldn't done this operation")
		return
	}

	for _, file := range files {
		result = append(result, file.Name())
	}

	a.sendResponse(w, http.StatusOK, result)
}

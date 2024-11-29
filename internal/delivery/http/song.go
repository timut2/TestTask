package http

import (
	"errors"
	"log"
	"net/http"

	"github.com/timut2/music-library-api/internal/delivery"
	"github.com/timut2/music-library-api/internal/model"
	"github.com/timut2/music-library-api/internal/repository/postgresql"
	errorresp "github.com/timut2/music-library-api/pkg/errors"
	"github.com/timut2/music-library-api/pkg/jsonutil"
	"github.com/timut2/music-library-api/pkg/validator"
)

type MusicLibraryService interface {
	Get(id int64) (*model.Song, error)
	GetVerse(id int64, filter model.VerseFilter) ([]*model.SongVerse, error)
	GetAll(model.SongFilter) ([]*model.Song, error)
	Update(song *model.Song) error
	Delete(id int64) error
	InsertSong(*model.NewSong) error
	InsertMusicInfo(string, string) error
}

// @Summary Get list of all songs
// @Tags song
// @Description Receive list of all songs
// @Accept json
// @Produce json
// @Param songName query string  false  "name search by songName"
// @Param group query string  false  "name search by group"
// @Param page query int false "Page number for pagination"
// @Param page_size query int false "Number of verses per page"
// @Success 200 {object} model.Songs
// @Failure 500 {object} model.ErrResponse
// @Router /song [get]
func (h *Handler) listLibraryHandler(w http.ResponseWriter, r *http.Request) {
	var filter model.SongFilter
	qs := r.URL.Query()
	v := validator.New()

	filter.Name = readString(qs, "songName", "")
	filter.Group = readString(qs, "group", "")

	filter.Page = readInt(qs, "page", 1, v)
	filter.PageSize = readInt(qs, "page_size", 10, v)

	if delivery.ValidateFilters(v, filter); !v.Valid() {

	}

	songs, err := h.service.GetAll(filter)
	if err != nil {
		log.Fatal(err)
		return
	}

	err = jsonutil.WriteJSON(w, http.StatusOK, jsonutil.Wrap{"songs": songs}, nil)
	if err != nil {
		log.Fatal(err)
	}
}

// @Summary Get all verses of song
// @Tags song
// @Description Receive list of all verses from a song
// @Accept json
// @Produce json
// @Param verseNumber query int false "Filter by verse number"
// @Param page query int false "Page number for pagination"
// @Param page_size query int false "Number of verses per page"
// @Success 200 {object} model.Verses
// @Failure 500 {object} model.ErrResponse
// @Router /song/{id} [get]
func (h *Handler) getSongTextHandler(w http.ResponseWriter, r *http.Request) {
	id, err := readIDParam(r)
	if err != nil {
		log.Fatal(err)
	}

	var filter model.VerseFilter
	qs := r.URL.Query()
	v := validator.New()

	filter.VerseNumber = readInt(qs, "verseNumber", 0, v)
	filter.Page = readInt(qs, "page", 1, v)
	filter.PageSize = readInt(qs, "page_size", 10, v)

	if delivery.ValidateVerseFilters(v, filter); !v.Valid() {

	}

	verses, err := h.service.GetVerse(id, filter)
	if err != nil {
		log.Fatal(err)
	}

	err = jsonutil.WriteJSON(w, http.StatusOK, jsonutil.Wrap{"verses": verses}, nil)
	if err != nil {
		log.Fatal(err)
	}
}

// @Summary Add a new song
// @Tags song
// @Description Add a new song to the music library
// @Accept json
// @Produce json
// @Param newSong body model.NewSong true "New song information"
// @Success 201 {object} model.Song
// @Failure 500 {object} model.ErrResponse
// @Router /song [post]
func (h *Handler) addSongHandler(w http.ResponseWriter, r *http.Request) {
	var input *model.NewSong

	err := jsonutil.ReadJSON(w, r, &input)
	if err != nil {
		log.Fatal(err)
	}

	err = h.service.InsertSong(input)
	if err != nil {
		log.Fatal(err)
	}

	err = jsonutil.WriteJSON(w, http.StatusOK, jsonutil.Wrap{"message": "songs inserted successfully"}, nil)
	if err != nil {
		log.Fatal(err)
	}
}

// @Summary Update an existing song
// @Tags song
// @Description Update the details of an existing song in the music library.
// @Accept json
// @Produce json
// @Param id path int true "Song ID"
// @Param songHolder body model.SongHolder true "Updated song information"
// @Success 200 {object} model.Song "Successfully updated song"
// @Success 201 {object} model.Song
// @Failure 500 {object} model.ErrResponse
// @Router /song/{id} [patch]
func (h *Handler) updateSongHandler(w http.ResponseWriter, r *http.Request) {
	id, err := readIDParam(r)
	if err != nil {
		errorresp.NotFoundResponse(w, r)
		return
	}

	song, err := h.service.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, postgresql.ErrRecordNotFound):
			errorresp.NotFoundResponse(w, r)
		default:
			errorresp.ServerErrorResponse(w, r, err)
		}
		return
	}

	var songHolder model.SongHolder

	err = jsonutil.ReadJSON(w, r, &songHolder)
	if err != nil {
		errorresp.BadRequestResponse(w, r, err)
		return
	}

	if songHolder.Name != nil {
		song.Name = *songHolder.Name
	}

	if songHolder.Name != nil {
		song.Group = *songHolder.Group
	}
	err = h.service.Update(song)
	if err != nil {
		errorresp.ServerErrorResponse(w, r, err)
	}

	err = jsonutil.WriteJSON(w, http.StatusOK, jsonutil.Wrap{"song": song}, nil)
	if err != nil {
		errorresp.ServerErrorResponse(w, r, err)
	}
}

// @Summary Delete an existing song
// @Tags song
// @Description Delete the details of an existing song from the music library.
// @Accept json
// @Produce json
// @Param id path int true "Song ID"
// @Success 200 {object} model.ErrResponse "Successfully deleted song"
// @Failure 400 {object} model.ErrResponse "Invalid ID supplied"
// @Failure 404 {object} model.ErrResponse "Song not found"
// @Failure 500 {object} model.ErrResponse "Internal server error"
// @Router /song/{id} [delete]
func (h *Handler) deleteSongHandler(w http.ResponseWriter, r *http.Request) {
	id, err := readIDParam(r)
	if err != nil {
		errorresp.NotFoundResponse(w, r)
		return
	}

	err = h.service.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, postgresql.ErrRecordNotFound):
			errorresp.NotFoundResponse(w, r)
		default:
			errorresp.ServerErrorResponse(w, r, err)
		}
	}

	err = jsonutil.WriteJSON(w, http.StatusOK, jsonutil.Wrap{"message": "song deleted successfully"}, nil)
	if err != nil {
		errorresp.ServerErrorResponse(w, r, err)
	}
}

// @Summary Add a new song
// @Tags song
// @Description Insert a new song into the music library.
// @Accept json
// @Produce json
// @Success 200 {object} model.ErrResponse "Successfully inserted song"
// @Failure 400 {object} model.ErrResponse "Invalid input"
// @Failure 500 {object} model.ErrResponse "Internal server error"
// @Router /verse [post]
func (h *Handler) addVerseHandler(w http.ResponseWriter, r *http.Request) {
	var input model.NewSong

	err := jsonutil.ReadJSON(w, r, &input)
	if err != nil {
		errorresp.BadRequestResponse(w, r, err)
		return
	}

	err = h.service.InsertMusicInfo(input.Group, input.Name)
	if err != nil {
		errorresp.ServerErrorResponse(w, r, err)
		return
	}

	err = jsonutil.WriteJSON(w, http.StatusOK, jsonutil.Wrap{"message": "inserted successfully"}, nil)
	if err != nil {
		errorresp.ServerErrorResponse(w, r, err)
	}

}

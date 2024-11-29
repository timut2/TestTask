package http

import (
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/julienschmidt/httprouter"
)

func (h *Handler) Routes() http.Handler {

	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/song", h.listLibraryHandler)
	router.HandlerFunc(http.MethodGet, "/song/:id", h.getSongTextHandler)
	router.HandlerFunc(http.MethodDelete, "/song/:id", h.deleteSongHandler)
	router.HandlerFunc(http.MethodPatch, "/song/:id", h.updateSongHandler)
	router.HandlerFunc(http.MethodPost, "/song", h.addSongHandler)
	router.HandlerFunc(http.MethodPost, "/verse", h.addVerseHandler)

	router.HandlerFunc(http.MethodGet, "/swagger/*any", httpSwagger.WrapHandler)

	return router
}

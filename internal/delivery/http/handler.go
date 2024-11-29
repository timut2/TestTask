package http

type Handler struct {
	service MusicLibraryService
}

func NewHandler(service MusicLibraryService) *Handler {
	return &Handler{service: service}
}

package delivery

import (
	"github.com/timut2/music-library-api/internal/model"
	"github.com/timut2/music-library-api/pkg/validator"
)

func ValidateFilters(v *validator.Validator, f model.SongFilter) {
	v.Check(f.Page > 0, "page", "must be greater than zero")
	v.Check(f.Page <= 10000, "page", "must be a maximum of 10 thousand")
	v.Check(f.PageSize > 0, "page_size", "must be greater than zero")
	v.Check(f.PageSize <= 100, "page_size", "must be a maximum of 100")
}

func ValidateVerseFilters(v *validator.Validator, f model.VerseFilter) {
	v.Check(f.Page > 0, "page", "must be greater than zero")
	v.Check(f.Page <= 10000, "page", "must be a maximum of 10 thousand")
	v.Check(f.PageSize > 0, "page_size", "must be greater than zero")
	v.Check(f.PageSize <= 100, "page_size", "must be a maximum of 100")
}

func ValidateSong(v *validator.Validator, song *model.Song) {

	v.Check(song.Name != "", "song_name", "must be provided")
	v.Check(len(song.Name) <= 60, "song_name", "must not be more than 60 bytes long")

	v.Check(song.Group != "", "group_name", "must be provided")
	v.Check(len(song.Group) <= 60, "group_name", "must not be more than 500 bytes long")

}

package postgresql

import (
	"context"
	"database/sql"
	"time"

	"github.com/timut2/music-library-api/internal/model"
)

type MusicInfoRepository struct {
	db *sql.DB
}

func NewMusicInfoRepository(db *sql.DB) *MusicInfoRepository {
	return &MusicInfoRepository{db: db}
}

func (gr *MusicInfoRepository) InsertMusicInfo(musicInfo *model.MusicInfo) error {
	query := `
	INSERT INTO music_info (release_date, text, link)
	VALUES ($1, $2, $3)
	RETURNING music_info_id
	`

	args := []any{musicInfo.ReleaseDate, musicInfo.Text, musicInfo.Link}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return gr.db.QueryRowContext(ctx, query, args...).Scan(&musicInfo.Id)
}

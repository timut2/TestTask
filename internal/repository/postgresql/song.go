package postgresql

import (
	"context"
	"database/sql"
	"time"

	"github.com/timut2/music-library-api/internal/model"
)

type SongsRepository struct {
	db *sql.DB
}

// InsertMusicInfo implements service.SongStorage.
func (sr *SongsRepository) InsertMusicInfo(*model.MusicInfo) error {
	panic("unimplemented")
}

func NewSongsRepository(db *sql.DB) *SongsRepository {
	return &SongsRepository{db: db}
}

func (sr *SongsRepository) GetAll(filter model.SongFilter) ([]*model.Song, error) {
	query := `
	SELECT song_id, name, music_group 
	FROM song
	WHERE (LOWER(name)=LOWER($1) OR $1 = '')
	AND (LOWER(music_group)=LOWER($2) or $2 = '') 
	ORDER BY song_id	
	LIMIT $3 OFFSET $4
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{
		filter.Name,
		filter.Group,
		filter.PageSize,
		(filter.Page - 1) * filter.PageSize,
	}

	rows, err := sr.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	songs := []*model.Song{}

	for rows.Next() {

		var song model.Song
		err = rows.Scan(
			&song.ID,
			&song.Name,
			&song.Group,
		)
		if err != nil {
			return nil, err
		}
		songs = append(songs, &song)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	for _, song := range songs {
		query2 := `
		SELECT verse_id, verse_number, text
		FROM verse
		WHERE song_id = $1 ORDER BY verse_number
		`
		arg := []any{
			song.ID,
		}
		rows2, err := sr.db.QueryContext(ctx, query2, arg...)
		if err != nil {
			return nil, err
		}
		defer rows2.Close()
		var verses []model.SongVerse
		for rows2.Next() {
			var verse model.SongVerse
			err = rows2.Scan(
				&verse.ID,
				&verse.VerseNumber,
				&verse.Text,
			)
			if err != nil {
				return nil, err
			}
			verses = append(verses, verse)
		}
		if err = rows.Err(); err != nil {
			return nil, err
		}
		song.Verse = verses
	}

	return songs, nil
}

func (sr *SongsRepository) GetVerse(id int64, filter model.VerseFilter) ([]*model.SongVerse, error) {
	query := `
		SELECT verse_id, verse_number, text
		FROM verse
		WHERE song_id = $1 
		AND (verse_number=$2 or $2=0)
		ORDER BY verse_id	
		LIMIT $3 OFFSET $4
		`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{
		id,
		filter.VerseNumber,
		filter.PageSize,
		(filter.Page - 1) * filter.PageSize,
	}

	rows, err := sr.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	verses := []*model.SongVerse{}

	for rows.Next() {
		var verse model.SongVerse
		err = rows.Scan(
			&verse.ID,
			&verse.VerseNumber,
			&verse.Text,
		)
		if err != nil {
			return nil, err
		}
		verses = append(verses, &verse)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return verses, nil
}

func (sr *SongsRepository) Get(id int64) (*model.Song, error) {
	query := `
	SELECT song.song_id, song.name, song.music_group
	FROM song  
	WHERE song_id=$1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var song model.Song

	err := sr.db.QueryRowContext(ctx, query, id).Scan(
		&song.ID, &song.Name, &song.Group,
	)
	if err != nil {
		return nil, err
	}

	query2 := `
	SELECT verse_id, verse_number, text
	FROM verse
	WHERE song_id = $1 ORDER BY verse_number
	`
	rows, err := sr.db.QueryContext(ctx, query2, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var verses []model.SongVerse

	for rows.Next() {
		var verse model.SongVerse
		err = rows.Scan(
			&verse.ID,
			&verse.VerseNumber,
			&verse.Text,
		)
		if err != nil {
			return nil, err
		}
		verses = append(verses, verse)
	}
	song.Verse = verses
	return &song, nil
}

func (sr *SongsRepository) Insert(newSong *model.NewSong) error {
	query := `
	INSERT INTO song (name, music_group)
	VALUES ($1, $2)
	RETURNING song_id
	`
	song := &model.Song{
		Name:  newSong.Name,
		Group: newSong.Group,
	}
	if newSong.Group != "" && newSong.Name != "" {
		args := []any{song.Name, song.Group}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return sr.db.QueryRowContext(ctx, query, args...).Scan(&song.ID)
	}

	return nil
}

func (sr *SongsRepository) Update(song *model.Song) error {
	query := `
		UPDATE song
		SET name = $1, music_group = $2
		WHERE song_id = $3`

	args := []any{
		song.Name,
		song.Group,
		song.ID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := sr.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (sr *SongsRepository) Delete(id int64) error {
	query := `
	DELETE FROM	song 
	WHERE song_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := sr.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return err
	}
	return nil
}

package service

import (
	"errors"

	"github.com/timut2/music-library-api/internal/model"
	"github.com/timut2/music-library-api/internal/repository/api"
	"github.com/timut2/music-library-api/pkg/logger"
)

type SongStorage interface {
	Get(id int64) (*model.Song, error)
	GetAll(model.SongFilter) ([]*model.Song, error)
	Delete(id int64) error
	Update(song *model.Song) error
	Insert(song *model.NewSong) error
	GetVerse(id int64, filter model.VerseFilter) ([]*model.SongVerse, error)
}

type MusicInfoStorage interface {
	InsertMusicInfo(*model.MusicInfo) error
}

type ApiClient interface {
	GetMusicInfo(group string, song string) (*model.MusicInfo, error)
}

type MusicLibrary struct {
	songStorage      SongStorage
	musicInfoStorage MusicInfoStorage
	apiClient        ApiClient
}

func (ml *MusicLibrary) Delete(id int64) error {
	return ml.songStorage.Delete(id)
}

func (ml *MusicLibrary) GetVerse(id int64, filter model.VerseFilter) ([]*model.SongVerse, error) {
	return ml.songStorage.GetVerse(id, filter)
}

func (ml *MusicLibrary) GetAll(filter model.SongFilter) ([]*model.Song, error) {
	return ml.songStorage.GetAll(filter)
}

func (ml *MusicLibrary) InsertSong(newSong *model.NewSong) error {
	return ml.songStorage.Insert(newSong)
}

func (ml *MusicLibrary) Update(song *model.Song) error {
	return ml.songStorage.Update(song)
}

func NewMusicLibrary(songStorage SongStorage, musicInfoStorage MusicInfoStorage, apiClient ApiClient) *MusicLibrary {
	return &MusicLibrary{
		songStorage:      songStorage,
		musicInfoStorage: musicInfoStorage,
		apiClient:        apiClient,
	}
}

func (ml *MusicLibrary) Get(id int64) (*model.Song, error) {
	return ml.songStorage.Get(id)
}

func (ml *MusicLibrary) InsertMusicInfo(group, name string) error {
	musicInfo, err := ml.apiClient.GetMusicInfo(group, name)
	if err != nil {
		if errors.Is(err, api.ErrBadRequest) {
			logger.PrintDebug("bad request ", map[string]any{
				"group": group,
				"name":  name,
				"error": err,
			})
		}
		if errors.Is(err, api.ErrInternalServerError) {
			logger.PrintDebug("error on server side", map[string]any{
				"error": err,
			})
		}
	}

	logger.PrintDebug("info from external API", map[string]any{
		"musicInfo": musicInfo,
	})

	if musicInfo != nil {
		err := ml.musicInfoStorage.InsertMusicInfo(musicInfo)
		if err != nil {
			return err
		}
	}
	return nil
}

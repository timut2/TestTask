package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/timut2/music-library-api/config"
	"github.com/timut2/music-library-api/internal/model"
)

type ApiClient struct {
	config *config.Config
}

func NewApiClient(config *config.Config) *ApiClient {
	return &ApiClient{config: config}
}

func (ac *ApiClient) GetMusicInfo(group string, song string) (*model.MusicInfo, error) {
	url := fmt.Sprintf("%s/info?group=%s&song=%s", ac.config.ExternalApiUrl, group, song)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received status code %d", resp.StatusCode)
	}

	var musicInfo model.MusicInfo
	if err := json.NewDecoder(resp.Body).Decode(&musicInfo); err != nil {
		return nil, err
	}

	return &musicInfo, nil
}

package model

type Song struct {
	ID    int64       `json:"id"`
	Name  string      `json:"song"`
	Group string      `json:"group"`
	Verse []SongVerse `json:"verses"`
}

type SongVerse struct {
	ID          int64  `json:"id"`
	VerseNumber int64  `json:"verse_number"`
	Text        string `json:"name"`
}

type SongHolder struct {
	Name  *string `json:"song"`
	Group *string `json:"group"`
}

type NewSong struct {
	Group string `json:"group"`
	Name  string `json:"song"`
}

type SongFilter struct {
	Name     string
	Group    string
	Page     int
	PageSize int
}

type VerseFilter struct {
	VerseNumber int
	Page        int
	PageSize    int
}

type Songs struct {
	Songs []Song
}

type Verses struct {
	Verses []SongVerse
}

type ErrResponse struct {
	Error any `json:"error"`
}

type SuccessResponse struct {
	Error any `json:"error"`
}

type MusicInfo struct {
	Id          string `json:"id"`
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

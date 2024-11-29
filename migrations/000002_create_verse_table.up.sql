CREATE TABLE IF NOT EXISTS verse(
    verse_id bigserial PRIMARY KEY,
    text TEXT NOT NULL,
    verse_number  integer,
    song_id bigint,
    FOREIGN KEY (song_id) REFERENCES song (song_id) ON DELETE SET NULL
);
CREATE TABLE IF NOT EXISTS song(
    song_id bigserial PRIMARY KEY,
    name text NOT NULL,
    music_group text NOT NULL
);
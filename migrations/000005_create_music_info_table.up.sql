CREATE TABLE IF NOT EXISTS music_info(
    music_info_id bigserial PRIMARY KEY,
    release_date date,
    text text NOT NULL,
    link text NOT NULL
);
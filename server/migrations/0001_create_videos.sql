-- +goose Up
CREATE TABLE IF NOT EXISTS videos (
	id 					text primary key,
	title 				text not null default '',
	description 		text not null default '',
	original_filename 	text not null,
	content_type 		text not null,
	size 				bigint not null,
	storage_key 		text not null,
	status 				text not null,
	created_at 			timestamptz not null default now(),
	updated_at 			timestamptz not null default now()
);

-- +goose Down
DROP TABLE videos;
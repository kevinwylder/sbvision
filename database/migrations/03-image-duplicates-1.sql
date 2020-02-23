
ALTER TABLE `videos` ADD COLUMN `session_id` INT(11) NOT NULL DEFAULT 0 AFTER `id`;
UPDATE videos
INNER JOIN images ON videos.thumbnail_id = images.id
SET videos.session_id = images.session_id;
ALTER TABLE videos ADD FOREIGN KEY videos_ibfk_5 (`session_id`) REFERENCES sessions (`id`);

ALTER TABLE `frames` ADD COLUMN `session_id` INT(11) NOT NULL DEFAULT 0 AFTER `id`;

ALTER TABLE `frames` ADD COLUMN `image_hash` BIGINT NOT NULL DEFAULT 0 AFTER `time`;

UPDATE frames 
INNER JOIN images ON frames.image_id = images.id
SET frames.session_id = images.session_id;

ALTER TABLE frames ADD FOREIGN KEY frames_ibfk_3 (`session_id`) REFERENCES sessions (id);


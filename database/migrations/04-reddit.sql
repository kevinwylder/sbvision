
ALTER TABLE `videos` DROP COLUMN `fps` ;

ALTER TABLE `videos` CHANGE COLUMN `url` `url` VARCHAR(1024) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '0'  COMMENT '' AFTER `duration`;

ALTER TABLE `videos` ADD COLUMN `source_url` VARCHAR(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' AFTER `url`;

ALTER TABLE `videos` ADD COLUMN `link_expires` TIMESTAMP NULL DEFAULT NULL AFTER `source_url`;

ALTER TABLE `videos` CHANGE COLUMN `discovery_time` `discovery_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP() COMMENT '' AFTER `link_expires`;

UPDATE videos 
INNER JOIN youtube_videos ON youtube_videos.video_id = videos.id
SET videos.url = youtube_videos.mirror_url,
    videos.link_expires = youtube_videos.mirror_expire,
    videos.source_url = CONCAT("https://www.youtube.com/watch?v=", youtube_videos.youtube_id);

DROP TABLE youtube_videos;

ALTER TABLE `videos` DROP FOREIGN KEY `videos_ibfk_1`;

DROP TABLE video_types;
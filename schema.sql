#
# SQL Export
# Created by Querious (201067)
# Created: February 5, 2020 at 11:25:01 AM PST
# Encoding: Unicode (UTF-8)
#

CREATE DATABASE IF NOT EXISTS `sbvision` DEFAULT CHARACTER SET utf8mb4 DEFAULT COLLATE utf8mb4_general_ci;
USE `sbvision`;

CREATE TABLE `video_types` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `description` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

REPLACE INTO `video_types` (`id`, `description`) VALUES 
	(1,'Youtube'),
	(2,'Reddit Gif');

CREATE TABLE `sessions` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `start` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `source_ip` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE `images` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `key` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `session_id` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_key` (`key`) USING BTREE,
  KEY `session_id` (`session_id`),
  CONSTRAINT `images_ibfk_1` FOREIGN KEY (`session_id`) REFERENCES `sessions` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE `videos` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `title` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `type` int(11) NOT NULL DEFAULT '0',
  `duration` int(11) NOT NULL DEFAULT '0',
  `fps` double NOT NULL DEFAULT '0',
  `thumbnail_id` int(11) NOT NULL DEFAULT '0',
  `discovery_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `type` (`type`),
  KEY `thumbnail_id` (`thumbnail_id`),
  CONSTRAINT `videos_ibfk_1` FOREIGN KEY (`type`) REFERENCES `video_types` (`id`),
  CONSTRAINT `videos_ibfk_4` FOREIGN KEY (`thumbnail_id`) REFERENCES `images` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE `youtube_videos` (
  `youtube_id` varchar(48) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '0',
  `video_id` int(11) NOT NULL DEFAULT '0',
  `mirror_url` varchar(4096) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `mirror_expire` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`youtube_id`),
  KEY `video_id` (`video_id`),
  CONSTRAINT `youtube_videos_ibfk_1` FOREIGN KEY (`video_id`) REFERENCES `videos` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE `frames` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `image_key` int(11) NOT NULL DEFAULT '0',
  `video_id` int(11) NOT NULL DEFAULT '0',
  `frame` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `video_id` (`video_id`),
  KEY `image_key` (`image_key`),
  CONSTRAINT `frames_ibfk_1` FOREIGN KEY (`video_id`) REFERENCES `videos` (`id`),
  CONSTRAINT `frames_ibfk_2` FOREIGN KEY (`image_key`) REFERENCES `images` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE `clips` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `session_id` int(11) NOT NULL DEFAULT '0',
  `frame_id` int(11) NOT NULL DEFAULT '0',
  `r` double NOT NULL DEFAULT '0',
  `i` double NOT NULL DEFAULT '0',
  `j` double NOT NULL DEFAULT '0',
  `k` double NOT NULL DEFAULT '0',
  `x` int(11) NOT NULL DEFAULT '0',
  `y` int(11) NOT NULL DEFAULT '0',
  `width` int(11) NOT NULL DEFAULT '0',
  `height` int(11) NOT NULL DEFAULT '0',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `session_id` (`session_id`),
  KEY `frame_id` (`frame_id`),
  CONSTRAINT `clips_ibfk_1` FOREIGN KEY (`session_id`) REFERENCES `sessions` (`id`),
  CONSTRAINT `clips_ibfk_2` FOREIGN KEY (`frame_id`) REFERENCES `frames` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;


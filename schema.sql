#
# SQL Export
# Created by Querious (201067)
# Created: February 23, 2020 at 3:50:16 PM PST
# Encoding: Unicode (UTF-8)
#


CREATE DATABASE IF NOT EXISTS `sbvision` DEFAULT CHARACTER SET utf8mb4 DEFAULT COLLATE utf8mb4_0900_ai_ci;
USE `sbvision`;




CREATE TABLE `video_types` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `description` varchar(128) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;


CREATE TABLE `videos` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `title` varchar(128) NOT NULL DEFAULT '',
  `type` int(11) NOT NULL DEFAULT '0',
  `format` varchar(16) NOT NULL DEFAULT 'video/mp4',
  `duration` int(11) NOT NULL DEFAULT '0',
  `fps` double NOT NULL DEFAULT '0',
  `discovery_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `type` (`type`),
  CONSTRAINT `videos_ibfk_1` FOREIGN KEY (`type`) REFERENCES `video_types` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=13 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;


CREATE TABLE `youtube_videos` (
  `youtube_id` varchar(48) NOT NULL DEFAULT '0',
  `video_id` int(11) NOT NULL DEFAULT '0',
  `mirror_url` varchar(4096) NOT NULL DEFAULT '',
  `mirror_expire` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`youtube_id`),
  KEY `video_id` (`video_id`),
  CONSTRAINT `youtube_videos_ibfk_1` FOREIGN KEY (`video_id`) REFERENCES `videos` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;




CREATE TABLE `sessions` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `start` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `source_ip` varchar(128) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=611 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;


CREATE TABLE `frames` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `session_id` int(11) NOT NULL DEFAULT '0',
  `video_id` int(11) NOT NULL DEFAULT '0',
  `time` int(11) NOT NULL DEFAULT '0',
  `image_hash` bigint(20) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `video_id` (`video_id`),
  KEY `frames_ibfk_3` (`session_id`),
  CONSTRAINT `frames_ibfk_1` FOREIGN KEY (`video_id`) REFERENCES `videos` (`id`),
  CONSTRAINT `frames_ibfk_3` FOREIGN KEY (`session_id`) REFERENCES `sessions` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2217 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;


CREATE TABLE `bounds` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `session_id` int(11) NOT NULL DEFAULT '0',
  `frame_id` int(11) NOT NULL DEFAULT '0',
  `x` int(11) NOT NULL DEFAULT '0',
  `y` int(11) NOT NULL DEFAULT '0',
  `width` int(11) NOT NULL DEFAULT '0',
  `height` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `frame_id` (`frame_id`),
  KEY `session_id` (`session_id`),
  CONSTRAINT `bounds_ibfk_1` FOREIGN KEY (`frame_id`) REFERENCES `frames` (`id`),
  CONSTRAINT `bounds_ibfk_2` FOREIGN KEY (`session_id`) REFERENCES `sessions` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1841 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;


CREATE TABLE `rotations` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `session_id` int(11) NOT NULL DEFAULT '0',
  `bounds_id` int(11) NOT NULL DEFAULT '0',
  `r` double NOT NULL DEFAULT '0',
  `i` double NOT NULL DEFAULT '0',
  `j` double NOT NULL DEFAULT '0',
  `k` double NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `session_id` (`session_id`),
  KEY `bounds_id` (`bounds_id`),
  CONSTRAINT `rotations_ibfk_1` FOREIGN KEY (`session_id`) REFERENCES `sessions` (`id`),
  CONSTRAINT `rotations_ibfk_2` FOREIGN KEY (`bounds_id`) REFERENCES `bounds` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=58 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;



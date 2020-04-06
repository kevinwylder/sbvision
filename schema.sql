#
# SQL Export
# Created by Querious (201069)
# Created: April 6, 2020 at 1:44:55 PM PDT
# Encoding: Unicode (UTF-8)
#

CREATE DATABASE IF NOT EXISTS sbvision;
USE sbvision;

SET @PREVIOUS_FOREIGN_KEY_CHECKS = @@FOREIGN_KEY_CHECKS;
SET FOREIGN_KEY_CHECKS = 0;


CREATE TABLE `sessions` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `start` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `source_ip` varchar(128) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=723 DEFAULT CHARSET=utf8mb4;


CREATE TABLE `frames` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `session_id` int(11) NOT NULL DEFAULT '0',
  `video_id` int(11) NOT NULL DEFAULT '0',
  `time` int(11) NOT NULL DEFAULT '0',
  `image_hash` bigint(20) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `video_id` (`video_id`),
  KEY `frames_ibfk_3` (`session_id`),
  CONSTRAINT `frames_ibfk_1` FOREIGN KEY (`video_id`) REFERENCES `videos` (`id`) ON DELETE CASCADE,
  CONSTRAINT `frames_ibfk_3` FOREIGN KEY (`session_id`) REFERENCES `sessions` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=3467 DEFAULT CHARSET=utf8mb4;


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
  CONSTRAINT `bounds_ibfk_1` FOREIGN KEY (`frame_id`) REFERENCES `frames` (`id`) ON DELETE CASCADE,
  CONSTRAINT `bounds_ibfk_2` FOREIGN KEY (`session_id`) REFERENCES `sessions` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=2601 DEFAULT CHARSET=utf8mb4;


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
  CONSTRAINT `rotations_ibfk_1` FOREIGN KEY (`session_id`) REFERENCES `sessions` (`id`) ON DELETE CASCADE,
  CONSTRAINT `rotations_ibfk_2` FOREIGN KEY (`bounds_id`) REFERENCES `bounds` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=7958 DEFAULT CHARSET=utf8mb4;


CREATE TABLE `videos` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `title` varchar(128) NOT NULL DEFAULT '',
  `type` int(11) NOT NULL DEFAULT '0',
  `format` varchar(16) NOT NULL DEFAULT 'video/mp4',
  `duration` int(11) NOT NULL DEFAULT '0',
  `url` varchar(1024) NOT NULL DEFAULT '',
  `source_url` varchar(128) NOT NULL DEFAULT '',
  `link_expires` timestamp NOT NULL DEFAULT '1970-01-01 00:00:01',
  `discovery_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_source_url` (`source_url`) USING BTREE,
  KEY `type` (`type`)
) ENGINE=InnoDB AUTO_INCREMENT=32 DEFAULT CHARSET=utf8mb4;




SET FOREIGN_KEY_CHECKS = @PREVIOUS_FOREIGN_KEY_CHECKS;



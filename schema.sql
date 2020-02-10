#
# SQL Export
# Created by Querious (201067)
# Created: February 10, 2020 at 11:07:40 AM PST
# Encoding: Unicode (UTF-8)
#


CREATE DATABASE IF NOT EXISTS `sbvision` DEFAULT CHARACTER SET utf8mb4 DEFAULT COLLATE utf8mb4_general_ci;
USE `sbvision`;




CREATE TABLE `sessions` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `start` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `source_ip` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=318 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;


CREATE TABLE `images` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `key` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `session_id` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_key` (`key`) USING BTREE,
  KEY `session_id` (`session_id`),
  CONSTRAINT `images_ibfk_1` FOREIGN KEY (`session_id`) REFERENCES `sessions` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=74 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;


CREATE TABLE `video_types` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `description` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;


CREATE TABLE `videos` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `title` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `type` int(11) NOT NULL DEFAULT '0',
  `format` varchar(16) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'video/mp4',
  `duration` int(11) NOT NULL DEFAULT '0',
  `fps` double NOT NULL DEFAULT '0',
  `thumbnail_id` int(11) NOT NULL DEFAULT '0',
  `discovery_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `type` (`type`),
  KEY `thumbnail_id` (`thumbnail_id`),
  CONSTRAINT `videos_ibfk_1` FOREIGN KEY (`type`) REFERENCES `video_types` (`id`),
  CONSTRAINT `videos_ibfk_4` FOREIGN KEY (`thumbnail_id`) REFERENCES `images` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;


CREATE TABLE `frames` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `image_id` int(11) NOT NULL DEFAULT '0',
  `video_id` int(11) NOT NULL DEFAULT '0',
  `frame` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `video_id` (`video_id`),
  KEY `image_key` (`image_id`),
  CONSTRAINT `frames_ibfk_1` FOREIGN KEY (`video_id`) REFERENCES `videos` (`id`),
  CONSTRAINT `frames_ibfk_2` FOREIGN KEY (`image_id`) REFERENCES `images` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=66 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;


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
) ENGINE=InnoDB AUTO_INCREMENT=27 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;


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



INSERT INTO `sessions` (`id`, `start`, `source_ip`) VALUES (1,'2020-02-06 22:02:57','172.31.0.1:36622');


INSERT INTO `images` (`id`, `key`, `session_id`) VALUES (1,'thumbnail/HGzalft_B_M.jpg',1);








INSERT INTO `video_types` (`id`, `description`) VALUES (1,'Youtube');
INSERT INTO `video_types` (`id`, `description`) VALUES (2,'Reddit Gif');
INSERT INTO `video_types` (`id`, `description`) VALUES (3,'Local');


INSERT INTO `videos` (`id`, `title`, `type`, `format`, `duration`, `fps`, `thumbnail_id`, `discovery_time`) VALUES (1,'Local Video for Testing',3,'video/mp4',95,30,1,'2020-02-09 02:40:08');







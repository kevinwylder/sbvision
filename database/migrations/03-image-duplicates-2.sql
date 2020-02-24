ALTER TABLE frames DROP FOREIGN KEY frames_ibfk_2;
ALTER TABLE `frames` DROP COLUMN `image_id` ;

ALTER TABLE videos DROP FOREIGN KEY videos_ibfk_4;
ALTER TABLE `videos` DROP COLUMN `thumbnail_id`;

DROP TABLE images;
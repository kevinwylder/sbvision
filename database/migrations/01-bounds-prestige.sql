ALTER TABLE `frames` ADD COLUMN `time` INT(11) NOT NULL DEFAULT 0 AFTER `frame`;

UPDATE `frames` SET `time` = 1000 * frame / 24.0;

ALTER TABLE `sbvision`.`frames` DROP COLUMN `frame` ;


CREATE DATABASE IF NOT EXISTS `test`;
USE `test`;

CREATE TABLE IF NOT EXISTS `user` (
   `id` INT NOT NULL AUTO_INCREMENT,
   `name` VARCHAR(20) NOT NULL,
   `sex` CHAR(1) NOT NULL COMMENT 'f/m',
   `age` TINYINT NOT NULL,
   `birth` DATE NOT NULL,
   PRIMARY KEY (`id`),
   UNIQUE KEY `idx_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
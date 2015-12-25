create database if not exists daylove;
use daylove;
CREATE TABLE `article` (
  `aid` int(11) NOT NULL AUTO_INCREMENT,
  `content` longtext CHARACTER SET utf8,
  `publish_time` datetime DEFAULT CURRENT_TIMESTAMP,
  `publish_status` tinyint(1) DEFAULT '1',
  PRIMARY KEY (`aid`),
  FULLTEXT KEY `content` (`content`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

create database if not exists gosense;
use gosense;
CREATE TABLE `top_article` (
  `aid` int(11) NOT NULL AUTO_INCREMENT,
  `title` varchar(255) CHARACTER SET utf8 DEFAULT '',
  `content` longtext CHARACTER SET utf8,
  `publish_time` datetime DEFAULT CURRENT_TIMESTAMP,
  `publish_status` tinyint(1) DEFAULT '1',
  `views` int(11) DEFAULT '1',
  PRIMARY KEY (`aid`),
  FULLTEXT KEY `content` (`title`,`content`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE `top_vistors` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `aid` int(11) DEFAULT '0',
  `user_agent` varchar(512) DEFAULT NULL,
  `ip` varchar(218) DEFAULT NULL,
  `access_time` varchar(45) DEFAULT NULL,
  `status` tinyint(1) DEFAULT '1',
  PRIMARY KEY (`id`),
  KEY `aid` (`aid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

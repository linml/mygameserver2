CREATE DATABASE `gameserver` /*!40100 DEFAULT CHARACTER SET utf8 */;
USE `gameserver`;


DROP TABLE IF EXISTS `bet_record`;
CREATE TABLE `bet_record` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned NOT NULL COMMENT '用戶編號',
  `aid` bigint unsigned NOT NULL DEFAULT 1,
  `period_number` varchar(30) NOT NULL COMMENT '期數',
  `rule_id` int unsigned NOT NULL COMMENT '玩法編號',
  `bet_content` varchar(30) COLLATE utf8_unicode_ci NOT NULL COMMENT '投注內容',
  `amount` decimal(20,4) unsigned NOT NULL COMMENT '投注額',
  `odds` decimal(7,4) unsigned NOT NULL COMMENT '賠率',
  `rebate` decimal(8,4) unsigned NOT NULL DEFAULT 0.0000 COMMENT '退水',
  `status` int unsigned NOT NULL DEFAULT 1 COMMENT '狀態 [0:要退單] [1:要下注]',
  `finished` bool  NOT NULL DEFAULT false COMMENT '結算狀態 [false:未結算] [true:已結算]',
  `room_id` int  NOT NULL DEFAULT 1 COMMENT '房間編號',
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;


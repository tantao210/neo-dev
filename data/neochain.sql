/*
SQLyog Ultimate v10.42 
MySQL - 5.7.21 : Database - neochain
*********************************************************************
*/

/*!40101 SET NAMES utf8 */;

/*!40101 SET SQL_MODE=''*/;

/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;
CREATE DATABASE /*!32312 IF NOT EXISTS*/`neochain` /*!40100 DEFAULT CHARACTER SET utf8 */;

USE `neochain`;

/*Table structure for table `account` */

DROP TABLE IF EXISTS `account`;

CREATE TABLE `account` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `address` varchar(50) NOT NULL COMMENT '地址',
  `private` varchar(100) NOT NULL COMMENT '私钥',
  `ismain` tinyint(4) NOT NULL DEFAULT '0' COMMENT '是否为主账号',
  `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `marks` tinyint(4) DEFAULT '0' COMMENT '删除标识',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=128 DEFAULT CHARSET=utf8;

/*Table structure for table `tx_order` */

DROP TABLE IF EXISTS `tx_order`;

CREATE TABLE `tx_order` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `txid` varchar(100) NOT NULL COMMENT '交易hash',
  `type` tinyint(4) NOT NULL COMMENT '交易类型1NEO 2GAS',
  `from` varchar(50) DEFAULT NULL COMMENT '输出账号',
  `to` varchar(50) DEFAULT NULL COMMENT '输入账号',
  `amount` decimal(20,8) DEFAULT '0.00000000' COMMENT '输出金额',
  `gas_price` decimal(20,8) DEFAULT '0.00000000' COMMENT '交易费',
  `sys_fee` decimal(20,8) DEFAULT '0.00000000' COMMENT '系统费',
  `net_fee` decimal(20,8) DEFAULT '0.00000000' COMMENT '网络费',
  `marks` tinyint(4) NOT NULL DEFAULT '0' COMMENT '删除标识',
  `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=303 DEFAULT CHARSET=utf8;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;


DROP TABLE IF EXISTS `gjc_account_key_tb`;
CREATE TABLE `gjc_account_key_tb`  (
      `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'ID',
      `account_name` varchar(40) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT NULL COMMENT '名称',
      `coin_type` varchar(20) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT NULL,
      `walletid` bigint(20) DEFAULT NULL,
      `priv_key` varchar(200) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '账户公钥',
      `pub_key` varchar(200) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '账户私钥',
      `address_id` varchar(200) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '账户公钥hash地址',
      `utxoid` varchar(500) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT NULL COMMENT 'UTXO交易id',
      `created_time` bigint(20) NOT NULL COMMENT '创建时间',
      `status` int(11) DEFAULT 1 COMMENT '使用状态',
      `updated_time` bigint(20) DEFAULT NULL COMMENT '更新时间',
      PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 82 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

##1020 sgj add

DROP TABLE IF EXISTS `wdc_account_key`;
CREATE TABLE `wdc_account_key`  (
	`id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'ID',
	`wallettype` varchar(40) DEFAULT NULL COMMENT 'walletid密码',
	`walletid` varchar(40) DEFAULT NULL COMMENT 'walletid',
	`coin_type` varchar(20)  DEFAULT NULL COMMENT '币种类型',
	`address` varchar(200)  NOT NULL COMMENT '账户公钥hash地址',
	`priv_key` varchar(200)  NOT NULL COMMENT '账户公钥',
	`pub_key` varchar(200) NOT NULL COMMENT '账户私钥',
	`pub_key_hash` varchar(200) NOT NULL COMMENT '账户私钥',
	`time_create` bigint(20) NOT NULL COMMENT '创建时间',
	`time_update` bigint(20) NOT NULL COMMENT '创建时间',
	PRIMARY KEY (`id`) USING BTREE
  )
  ENGINE = InnoDB
  DEFAULT CHARACTER SET = utf8;



--1030 add:
---  use ggexaddressall;
---  ALTER TABLE `wdc_account_key` ADD COLUMN `pub_key_hash` varchar(200)  DEFAULT NULL AFTER `pub_key`;
---
##1103add
DROP TABLE IF EXISTS `ggex_tran_state`;
CREATE TABLE `ggex_tran_state`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `settleid` bigint(20) DEFAULT NULL,
  `txhash` varchar(255) CHARACTER SET latin1 COLLATE latin1_swedish_ci DEFAULT NULL,
  `from` varchar(255) CHARACTER SET latin1 COLLATE latin1_swedish_ci DEFAULT NULL,
  `to` varchar(255) CHARACTER SET latin1 COLLATE latin1_swedish_ci DEFAULT NULL,
  `amount` double DEFAULT NULL,
  `amountfee` double DEFAULT NULL,
  `coincode` varchar(255) CHARACTER SET latin1 COLLATE latin1_swedish_ci DEFAULT NULL,
  `status` varchar(255) CHARACTER SET latin1 COLLATE latin1_swedish_ci DEFAULT NULL,
  `verifystatus` bigint(4) NULL DEFAULT '0',
  `time_create` varchar(30) CHARACTER SET latin1 COLLATE latin1_swedish_ci DEFAULT NULL,
  `time_update` varchar(30) CHARACTER SET latin1 COLLATE latin1_swedish_ci DEFAULT NULL,
  `errcode` bigint(20) DEFAULT NULL,
  `desc` text CHARACTER SET latin1 COLLATE latin1_swedish_ci,
  `raw` text CHARACTER SET latin1 COLLATE latin1_swedish_ci,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `index_settleid`(`settleid`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 186 CHARACTER SET = latin1 COLLATE = latin1_swedish_ci ROW_FORMAT = Dynamic;

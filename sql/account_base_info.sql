-- MySQL dump 10.13  Distrib 8.0.32, for Linux (x86_64)
--
-- Host: localhost    Database: account_base_info
-- ------------------------------------------------------
-- Server version	8.0.32-0ubuntu0.20.04.2

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

create database account_base_info default character set utf8mb4 collate utf8mb4_0900_ai_ci;
use account_base_info;

--
-- Table structure for table `holiday`
--

DROP TABLE IF EXISTS `holiday`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `holiday` (
  `ymd` date NOT NULL,
  PRIMARY KEY (`ymd`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='法定节假日日期';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `holiday`
--

LOCK TABLES `holiday` WRITE;
/*!40000 ALTER TABLE `holiday` DISABLE KEYS */;
INSERT INTO `holiday` VALUES ('2023-04-05'),('2023-05-01'),('2023-06-22'),('2023-09-29'),('2023-10-01'),('2023-10-02'),('2023-10-03');
/*!40000 ALTER TABLE `holiday` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `mail_tpl`
--

DROP TABLE IF EXISTS `mail_tpl`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `mail_tpl` (
  `id` mediumint NOT NULL AUTO_INCREMENT,
  `type` tinyint(1) NOT NULL COMMENT '1: 注册, 2: 忘记密码, 3: 账号绑定, 4: 账号解绑, 5:登录',
  `lang_id` varchar(32) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT '' COMMENT '语言标识, 简体:zh-CN, 繁体: zh-TW, 英文: en',
  `title` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT '' COMMENT '标题',
  `content` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci COMMENT '内容',
  `updated_time` int DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `type` (`type`,`lang_id`)
) ENGINE=InnoDB AUTO_INCREMENT=16 DEFAULT CHARSET=utf8mb3 COMMENT='用户忘记密码/绑定/解绑邮件模板表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `mail_tpl`
--

LOCK TABLES `mail_tpl` WRITE;
/*!40000 ALTER TABLE `mail_tpl` DISABLE KEYS */;
INSERT INTO `mail_tpl` VALUES (1,2,'zh-TW','XGame 會員更改密碼通知','<h3 style=\"font-weight: 900;font-weight: 700\r\n            margin-bottom: 20px;font-size:20px\">親愛的XGame會員您好:</h3>\r\n\r\n    <div style=\"font-size:16px;margin-bottom: 10px;\">\r\n        <p style=\"line-height: 30px;\r\n            margin-bottom: 10px;\">&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;您已通過信箱重設新密碼，本次請求的認證碼為：\r\n            <span style=\"color:red\">{CODE}</span>\r\n            請在認證碼輸入框中輸入此認證>碼以完成認證。（認證碼有效期為<span style=\"color:red\">120</span>分鐘）\r\n        </p>\r\n\r\n        <div>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;如您有疑問，請聯繫我們客服：\r\n            <a style=\"color: red;\" href=\"Mailto:service@xgame.com\">客服>信箱</a>\r\n        </div>\r\n        <div style=\"margin: 10px 0;\">&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-XGame-</div>\r\n        <div>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;此為系統自動發送信件，請勿回覆。</div>\r\n    </div>',NULL),(2,2,'zh-CN','XGame 会员更改密码通知','<h3 style=\"font-weight: 900;font-weight: 700\r\n            margin-bottom: 20px;font-size:20px\">亲爱的XGame用户您好:</h3>\r\n\r\n    <div style=\"font-size:16px;margin-bottom: 10px;\">\r\n        <p style=\"line-height: 30px;\r\n            margin-bottom: 10px;\">&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;您已通过信箱重设新密码，本次请求的认证码为：\r\n            <span style=\"color:red\">{CODE}</span>\r\n            请在认证码输入框中输入此认证码以完成认证。（认证码有效期为<span style=\"color:red\">10</span>分钟）\r\n        </p>\r\n\r\n        <div style=\"margin: 10px 0;\">&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-XGame-</div>\r\n        <div>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;此为系统自动发送信件，请勿回复。</div>\r\n    </div>\r\n',NULL),(3,2,'en','XGame Reminder: Reset Password','<h3 style=\"font-weight: 900;font-weight: 700\r\n            margin-bottom: 20px;font-size:20px\">Dear XGame user,</h3>\r\n\r\n    <div style=\"font-size:16px;margin-bottom: 10px;\">\r\n        <p style=\"line-height: 30px;\r\n            margin-bottom: 10px;\">&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;We received a request to reset your password via E-mail. Here is the verifcation code for this request:\r\n            <span style=\"color:red\">{CODE}</span>\r\n            Please enter it in the verfication code box to continue your reset.（The code will expire in <span style=\"color:red\">120</span>minutes）\r\n        </p>\r\n\r\n        <div>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;For more questions, please contact our customer service:\r\n            <a style=\"color: red;\" href=\"Mailto:service@xgame.com\">Customer Service Mail</a>\r\n        </div>\r\n        <div style=\"margin: 10px 0;\">&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-XGame-</div>\r\n        <div>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;This is an automatically generated email. Please do not reply.</div>\r\n    </div>\r\n',NULL),(4,3,'zh-TW','XGame 會員綁定郵件通知','<h3 style=\"font-weight: 900;font-weight: 700\r\n            margin-bottom: 20px;font-size:20px\">親愛的XGame會員您好:</h3>\r\n\r\n    <div style=\"font-size:16px;margin-bottom: 10px;\">\r\n        <p style=\"line-height: 30px;\r\n            margin-bottom: 10px;\">&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;您正在進行帳號綁定，本次請求的認證碼\n為：\r\n            <span style=\"color:red\">{CODE}</span>\r\n            請在認證碼輸入框中輸入此認證碼以完成認證。（認證碼有效期為<span style=\"color:red\">120</span>分鐘）\r\n        </p>\r\n\r\n        <div>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;如您有疑>問，請聯繫我們客服：\r\n            <a style=\"color: red;\" href=\"Mailto:service@xgame.com\">客服信箱</a>\r\n        </div>\r\n        <div style=\"margin: 10px 0;\">&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-XGame-</div>\r\n        <div>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;此為系統自動發送信件，請勿回覆。</div>\r\n    </div>\r\n',NULL),(5,3,'zh-CN','XGame 会员绑定邮件通知','<h3 style=\"font-weight: 900;font-weight: 700\r\n            margin-bottom: 20px;font-size:20px\">亲爱的XGame用户您好:</h3>\r\n\r\n    <div style=\"font-size:16px;margin-bottom: 10px;\">\r\n        <p style=\"line-height: 30px;\r\n            margin-bottom: 10px;\">&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;您正在进行账号绑定，本次请求的认证码为：\r\n            <span style=\"color:red\">{CODE}</span>\r\n            请在认证码输入框中输入此认证码以完成认证。（认证码有效期为<span style=\"color:red\">10</span>分钟）\r\n        </p>\r\n\r\n        <div style=\"margin: 10px 0;\">&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-XGame-</div>\r\n        <div>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;此为系统自动发送信件，请勿回复。</div>\r\n    </div>',NULL),(6,3,'en','XGame Reminder: Bind Account','<h3 style=\"font-weight: 900;font-weight: 700\r\n            margin-bottom: 20px;font-size:20px\">Dear XGame user,</h3>\r\n\r\n    <div style=\"font-size:16px;margin-bottom: 10px;\">\r\n        <p style=\"line-height: 30px;\r\n            margin-bottom: 10px;\">&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;You are trying to bind your account. Here is the verifcation code for this request:\r\n            <span style=\"color:red\">{CODE}</span>\r\n            Please enter it in the verfication code box and continue your reset.（The code will expire in <span style=\"color:red\">120</span>minutes）\r\n        </p>\r\n\r\n        <div>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;For more inquiry, please contact our customer service:\r\n            <a style=\"color: red;\" href=\"Mailto:service@xgame.com\">Customer Service Mail</a>\r\n        </div>\r\n        <div style=\"margin: 10px 0;\">&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-XGame-</div>\r\n        <div>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;This is an automatically generated email. Please do not reply.</div>\r\n    </div>',NULL),(7,4,'zh-TW','XGame 會員解除綁定通知','<h3 style=\"font-weight: 900;font-weight: 700\r\n            margin-bottom: 20px;font-size:20px\">親愛的XGame會員您好:</h3>\r\n\r\n    <div style=\"font-size:16px;margin-bottom: 10px;\">\r\n        <p style=\"line-height: 30px;\r\n            margin-bottom: 10px;\">&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;您正在進行帳號解除綁定，本次請求的認證碼為：\r\n            <span style=\"color:red\">{CODE}</span>\r\n            請在認證碼輸入框中輸入此認證碼以完成認證。（認證碼有效期為<span style=\"color:red\">120</span>分鐘）\r\n        </p>\r\n\r\n        <div>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;如您有疑問，請聯繫我們客服：\r\n            <a style=\"color: red;\" href=\"Mailto:service@xgame.com\">客服信箱</a>\r\n        </div>\r\n        <div style=\"margin: 10px 0;\">&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-XGame-</div>\r\n        <div>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;此為系統自動發送信件，請勿回覆。</div>\r\n    </div>',NULL),(8,4,'zh-CN','XGame 会员解除绑定通知','<h3 style=\"font-weight: 900;font-weight: 700\r\n            margin-bottom: 20px;font-size:20px\">亲爱的XGame会员您好:</h3>\r\n\r\n    <div style=\"font-size:16px;margin-bottom: 10px;\">\r\n        <p style=\"line-height: 30px;\r\n            margin-bottom: 10px;\">&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;您正在进行账号解除绑\n定，本次请求的认证码为：\r\n            <span style=\"color:red\">{CODE}</span>\r\n            请在认证码输入框中输入此认证码以完成认证。（认证码有效期为<span style=\"color:red\">10</span>分钟）\r\n        </p>\r\n        <div style=\"margin: 10px 0;\">&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-XGame-</div>\r\n        <div>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;此为系统自动发送信件，请勿回复。</div>\r\n    </div>',NULL),(9,4,'en','XGame Reminder: Unbind Account','<h3 style=\"font-weight: 900;font-weight: 700\r\n            margin-bottom: 20px;font-size:20px\">Dear XGame user,</h3>\r\n\r\n    <div style=\"font-size:16px;margin-bottom: 10px;\">\r\n        <p style=\"line-height: 30px;\r\n            margin-bottom: 10px;\">&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;You are trying to unbind your account. Here is the verifcation code for this request:\r\n            <span style=\"color:red\">{CODE}</span>\r\n            Please enter it in the verfication code box and continue your reset.（The code will expire in <span style=\"color:red\">120</span>minutes）\r\n        </p>\r\n\r\n        <div>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;For more inquiry, please contact our customer service:\r\n            <a style=\"color: red;\" href=\"Mailto:service@xgame.com\">Customer Service Mail</a>\r\n        </div>\r\n        <div style=\"margin: 10px 0;\">&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-XGame-</div>\r\n        <div>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;This is an automatically generated email. Please do not reply.</div>\r\n    </div>',NULL),(14,1,'zh-CN','XGame 会员注册验证码通知','<h3 style=\"font-weight: 900;font-weight: 700\r\n            margin-bottom: 20px;font-size:20px\">亲爱的XGame用户您好:</h3>\r\n\r\n    <div style=\"font-size:16px;margin-bottom: 10px;\">\r\n        <p style=\"line-height: 30px;\r\n            margin-bottom: 10px;\">&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;您本次请求的验证码为：\r\n            <span style=\"color:red\">{CODE}</span>\r\n            请在验证码输入框中输入此验证码完成检验。（认证码有效期为<span style=\"color:red\">10</span>分钟）\r\n        </p>\r\n\r\n        <div style=\"margin: 10px 0;\">&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-XGame-</div>\r\n        <div>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;此为系统自动发送信件，请勿回复。</div>\r\n    </div>\r\n',NULL);
/*!40000 ALTER TABLE `mail_tpl` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sms_tpl`
--

DROP TABLE IF EXISTS `sms_tpl`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sms_tpl` (
  `id` int NOT NULL AUTO_INCREMENT,
  `type` tinyint(1) NOT NULL COMMENT '类型, 1: 注册, 2: 忘记密码, 3: 账号绑定, 4: 账号解绑, 5:登录',
  `lang_id` varchar(16) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL DEFAULT '' COMMENT '语言id, i18n约束id, zh-CN: 简体中文, zh-TW: 繁体中文, en-US: 英文 ',
  `sms_id` varchar(32) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL DEFAULT '' COMMENT '模板id(目前是阿里模板id)',
  `title` varchar(64) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT '' COMMENT '短信标题',
  PRIMARY KEY (`id`),
  UNIQUE KEY `lang_id` (`type`,`lang_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb3 COMMENT='用户短信模板';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sms_tpl`
--

LOCK TABLES `sms_tpl` WRITE;
/*!40000 ALTER TABLE `sms_tpl` DISABLE KEYS */;
INSERT INTO `sms_tpl` VALUES (1,1,'zh-CN','SMS_123456789','注册');
/*!40000 ALTER TABLE `sms_tpl` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `user_delete_config`
--

DROP TABLE IF EXISTS `user_delete_config`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `user_delete_config` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '自增id',
  `game_id` int NOT NULL COMMENT '游戏id',
  `platform_id` int NOT NULL COMMENT '大区id',
  `user_delete_wait_duration` tinyint NOT NULL DEFAULT '15' COMMENT '注销冷静期时长，单位天',
  `apple_client_id` varchar(50) NOT NULL DEFAULT '' COMMENT '苹果客户端id',
  `apple_client_secret` varchar(512) NOT NULL DEFAULT '' COMMENT '苹果客户端秘钥',
  `apple_config` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '定时生成secret需要的数据,json格式',
  `last_apple_update_time` int unsigned NOT NULL DEFAULT '0' COMMENT '最后一次更新secret时间',
  `ext` varchar(1024) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT '' COMMENT '扩展备用信息',
  PRIMARY KEY (`id`),
  UNIQUE KEY `game_platform` (`game_id`,`platform_id`)
) ENGINE=InnoDB AUTO_INCREMENT=17 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='项目配置表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `user_delete_config`
--

LOCK TABLES `user_delete_config` WRITE;
/*!40000 ALTER TABLE `user_delete_config` DISABLE KEYS */;
INSERT INTO `user_delete_config` VALUES (1,16,1,15,'x','x','{}',1678416829,'');
/*!40000 ALTER TABLE `user_delete_config` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `white_user_list`
--

DROP TABLE IF EXISTS `white_user_list`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `white_user_list` (
  `id` int NOT NULL AUTO_INCREMENT,
  `game_id` int NOT NULL DEFAULT '0' COMMENT '游戏编号',
  `platform_id` int NOT NULL DEFAULT '0' COMMENT '大区id',
  `uid` bigint DEFAULT NULL COMMENT '用户uid',
  `ip` varchar(64) DEFAULT '' COMMENT '用户ip',
  `env` varchar(255) NOT NULL DEFAULT '' COMMENT '环境，竖线分隔',
  `desc` varchar(255) NOT NULL DEFAULT '' COMMENT '备注信息',
  `updated_time` int NOT NULL DEFAULT '0' COMMENT '最后更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `game_platform_uid` (`game_id`,`platform_id`,`uid`) USING BTREE,
  UNIQUE KEY `game_platform_ip` (`game_id`,`platform_id`,`ip`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb3 ROW_FORMAT=DYNAMIC COMMENT='用户白名单';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `white_user_list`
--

LOCK TABLES `white_user_list` WRITE;
/*!40000 ALTER TABLE `white_user_list` DISABLE KEYS */;
INSERT INTO `white_user_list` VALUES (1,16,1,1610000100014,'','dev|newversion|prerelease|produce|review','',0),(2,16,1,NULL,'192.0.2.1','dev|newversion|review','',0);
/*!40000 ALTER TABLE `white_user_list` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2023-04-05 14:23:05

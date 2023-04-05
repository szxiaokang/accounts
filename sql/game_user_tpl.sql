CREATE TABLE `user_%d`
(
    `account`      varchar(128) NOT NULL COMMENT '账号account,确保数据唯一性',
    `uid`          bigint       NOT NULL COMMENT '项目uid',
    `main_uid`     bigint       NOT NULL COMMENT '主账号uid',
    `status`       tinyint(1) NOT NULL DEFAULT '1' COMMENT '1：正常状态, 2：被禁用，3：注销中，4：已注销',
    `type`         tinyint(4) NOT NULL DEFAULT '0' COMMENT '注册类型 1: email, 2:手机号, 3: 用户名, 4: 游客, 5: 第三方',
    `created_time` int          NOT NULL COMMENT '创建时间',
    `ext`          varchar(255) DEFAULT '' COMMENT '扩展备用信息',
    PRIMARY KEY (`account`),
    KEY            `main_uid` (`main_uid`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='项目用户表';


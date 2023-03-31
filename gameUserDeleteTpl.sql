CREATE TABLE `user_delete_apply_%d`
(
    `uid`                 bigint       NOT NULL COMMENT '项目uid',
    `main_uid`            bigint       NOT NULL COMMENT '账号uid',
    `account`             varchar(128) NOT NULL COMMENT '账号account',
    `type`                tinyint(4) NOT NULL DEFAULT '0' COMMENT '注册类型 1: email, 2:手机号, 3: 用户名, 4: 游客, 5: 第三方',
    `status`              tinyint(1) NOT NULL DEFAULT '1' COMMENT '1：申请注销冷静期中或立即注销，2：系统处理中，3：删除成功, 4: 申请恢复，5：恢复成功',
    `apply_time`          int          NOT NULL COMMENT '申请时间',
    `execute_delete_time` int          NOT NULL COMMENT '执行注销时间',
    `ext_info`            text         NOT NULL COMMENT '其他信息，json格式，例如包含调用苹果revoke请求的refresh_token',
    `ext`                 varchar(255) DEFAULT '' COMMENT '扩展备用信息, 信息来自于项目用户表ext字段',
    PRIMARY KEY (`account`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='注销申请表';


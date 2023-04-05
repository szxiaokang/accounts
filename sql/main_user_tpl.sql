CREATE TABLE `account_%d`
(
    `uid`             bigint NOT NULL,
    `username`        varchar(64)           DEFAULT NULL COMMENT '用户名',
    `email`           varchar(128)          DEFAULT NULL COMMENT '用户的email',
    `guest`           varchar(128)          DEFAULT NULL COMMENT '游客',
    `third`           varchar(128)          DEFAULT NULL COMMENT '第三方账号，第三方名称id加uid,如fb账号：1001_112257954430192',
    `mobile`          varchar(32)           DEFAULT NULL COMMENT '手机号',
    `password`        char(32)     NOT NULL DEFAULT '' COMMENT '密码',
    `created_time`    int(11) NOT NULL DEFAULT '0' COMMENT '注册时间',
    `created_ip`      varchar(16)  NOT NULL DEFAULT '' COMMENT '注册时ip',
    `updated_time`    int(11) NOT NULL DEFAULT '0' COMMENT '最后更新时间',
    `last_login_time` int(11) NOT NULL DEFAULT '0' COMMENT '最后登录时间',
    `name`            varchar(32)  NOT NULL DEFAULT '' COMMENT '实名认证',
    `card_id`         varchar(32)  NOT NULL DEFAULT '' COMMENT '身份证id',
    `type`            tinyint(1) NOT NULL DEFAULT '0' COMMENT '注册类型 1: email, 2:手机号, 3: 用户名, 4: 游客, 5: 第三方',
    `salt`            char(16)     NOT NULL DEFAULT '' COMMENT '密码盐',
    `device_type`     tinyint(1) NOT NULL DEFAULT '0' COMMENT '1 ios, 2 android, 3 pc/h5, 4 其他',
    `lang`            varchar(32)  NOT NULL DEFAULT '' COMMENT '用户语言',
    `extra`           varchar(255) NOT NULL DEFAULT '' COMMENT '其他扩展部分',
    PRIMARY KEY (`uid`) USING BTREE,
    UNIQUE KEY `email` (`email`) USING BTREE,
    UNIQUE KEY `username` (`username`) USING BTREE,
    UNIQUE KEY `guest` (`guest`) USING BTREE,
    UNIQUE KEY `third` (`third`) USING BTREE,
    UNIQUE KEY `mobile` (`mobile`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='账号表';

CREATE TABLE `account_hash_%d`
(
    `account` varchar(128) NOT NULL COMMENT 'Hash的账号',
    `uid`     bigint DEFAULT NULL COMMENT '账号id',
    PRIMARY KEY (`account`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='账号Hash表';


CREATE TABLE `t_build`
(
    `id`                  varchar(64) NOT NULL,
    `pipeline_id`         varchar(64) NULL DEFAULT NULL,
    `pipeline_version_id` varchar(64) NULL DEFAULT NULL,
    `status`              varchar(100) NULL DEFAULT NULL COMMENT '构建状态',
    `error`               varchar(500) NULL DEFAULT NULL COMMENT '错误信息',
    `event`               varchar(100) NULL DEFAULT NULL COMMENT '事件',
    `time_stamp`          datetime(0) NULL DEFAULT NULL COMMENT '执行时长',
    `title`               varchar(255) NULL DEFAULT NULL COMMENT '标题',
    `message`             varchar(255) NULL DEFAULT NULL COMMENT '构建信息',
    `started`             datetime(0) NULL DEFAULT NULL COMMENT '开始时间',
    `finished`            datetime(0) NULL DEFAULT NULL COMMENT '结束时间',
    `created`             datetime(0) NULL DEFAULT NULL COMMENT '创建时间',
    `updated`             datetime(0) NULL DEFAULT NULL COMMENT '更新时间',
    `version`             varchar(255) NULL DEFAULT NULL COMMENT '版本',
    PRIMARY KEY (`id`) USING BTREE
);
CREATE TABLE `t_cmd_line`
(
    `id`       varchar(64) NOT NULL,
    `group_id` varchar(64) NULL DEFAULT NULL,
    `build_id` varchar(64) NULL DEFAULT NULL,
    `step_id`  varchar(64) NULL DEFAULT NULL,
    `status`   varchar(50) NULL DEFAULT NULL,
    `num`      int(11) NULL DEFAULT NULL,
    `content`  text NULL,
    `created`  datetime(0) NULL DEFAULT NULL,
    `started`  datetime(0) NULL DEFAULT NULL,
    `finished` datetime(0) NULL DEFAULT NULL,
    PRIMARY KEY (`id`) USING BTREE
);
CREATE TABLE `t_stage`
(
    `id`                  varchar(64) NOT NULL,
    `pipeline_version_id` varchar(64) NULL DEFAULT NULL COMMENT '流水线id',
    `build_id`            varchar(64) NULL DEFAULT NULL,
    `status`              varchar(100) NULL DEFAULT NULL COMMENT '构建状态',
    `exit_code`           bigint(20) NULL DEFAULT NULL COMMENT '退出码',
    `error`               varchar(500) NULL DEFAULT NULL COMMENT '错误信息',
    `name`                varchar(255) NULL DEFAULT NULL COMMENT '名字',
    `display_name`        varchar(255) NULL DEFAULT NULL,
    `started`             datetime(0) NULL DEFAULT NULL COMMENT '开始时间',
    `finished`            datetime(0) NULL DEFAULT NULL COMMENT '结束时间',
    `created`             datetime(0) NULL DEFAULT NULL COMMENT '创建时间',
    `updated`             datetime(0) NULL DEFAULT NULL COMMENT '更新时间',
    `version`             varchar(255) NULL DEFAULT NULL COMMENT '版本',
    `on_success`          varchar(5) NULL DEFAULT NULL,
    `on_failure`          varchar(5) NULL DEFAULT NULL,
    `sort`                int(11) NULL DEFAULT NULL,
    `stage`               varchar(255) NULL DEFAULT NULL,
    PRIMARY KEY (`id`) USING BTREE
);
CREATE TABLE `t_step`
(
    `id`                  varchar(64) NOT NULL,
    `build_id`            varchar(64) NULL DEFAULT NULL,
    `stage_id`            varchar(100) NULL DEFAULT NULL COMMENT '流水线id',
    `display_name`        varchar(255) NULL DEFAULT NULL,
    `pipeline_version_id` varchar(64) NULL DEFAULT NULL COMMENT '流水线id',
    `step`                varchar(255) NULL DEFAULT NULL,
    `status`              varchar(100) NULL DEFAULT NULL COMMENT '构建状态',
    `event`               varchar(100) NULL DEFAULT NULL COMMENT '事件',
    `exit_code`           int(11) NULL DEFAULT NULL COMMENT '退出码',
    `error`               varchar(500) NULL DEFAULT NULL COMMENT '错误信息',
    `name`                varchar(100) NULL DEFAULT NULL COMMENT '名字',
    `started`             datetime(0) NULL DEFAULT NULL COMMENT '开始时间',
    `finished`            datetime(0) NULL DEFAULT NULL COMMENT '结束时间',
    `created`             datetime(0) NULL DEFAULT NULL COMMENT '创建时间',
    `updated`             datetime(0) NULL DEFAULT NULL COMMENT '更新时间',
    `version`             varchar(255) NULL DEFAULT NULL COMMENT '版本',
    `errignore`           varchar(5) NULL DEFAULT NULL,
    `commands`            text NULL,
    `depends_on`          json NULL,
    `image`               varchar(255) NULL DEFAULT NULL,
    `environments`        json NULL,
    `sort`                int(11) NULL DEFAULT NULL,
    PRIMARY KEY (`id`) USING BTREE
);
CREATE TABLE `t_message`
(
    `id`      varchar(64) NOT NULL,
    `aid`     BIGINT      NOT NULL AUTO_INCREMENT,
    `uid`     varchar(64) NULL DEFAULT NULL COMMENT '发送者（可空）',
    `title`   varchar(255) NULL DEFAULT NULL,
    `content` longtext NULL,
    `types`   varchar(50) NULL DEFAULT NULL,
    `created` datetime(0) NULL DEFAULT NULL,
    `infos`   text NULL,
    `url`     varchar(500) NULL DEFAULT NULL,
    PRIMARY KEY (`aid`, `id`) USING BTREE
);
CREATE TABLE `t_org`
(
    `id`           varchar(64) NOT NULL,
    `aid`          BIGINT      NOT NULL AUTO_INCREMENT,
    `uid`          varchar(64) NULL DEFAULT NULL,
    `name`         varchar(200) NULL DEFAULT NULL,
    `desc`         TEXT NULL DEFAULT NULL,
    `public`       INT(1) NULL DEFAULT 0 COMMENT '公开',
    `created`      datetime(0) NULL DEFAULT NULL COMMENT '创建时间',
    `updated`      datetime(0) NULL DEFAULT NULL COMMENT '更新时间',
    `deleted`      int(1) NULL DEFAULT 0,
    `deleted_time` datetime(0) NULL DEFAULT NULL,
    PRIMARY KEY (`aid`, `id`) USING BTREE,
    INDEX          `uid`(`uid`) USING BTREE
);
CREATE TABLE `t_org_pipe`
(
    `aid`     bigint(20) NOT NULL AUTO_INCREMENT,
    `org_id`  varchar(64) NULL DEFAULT NULL,
    `pipe_id` varchar(64) NULL DEFAULT NULL COMMENT '收件人',
    `created` datetime(0) NULL DEFAULT NULL,
    `public`  INT(1) NULL DEFAULT 0 COMMENT '公开',
    PRIMARY KEY (`aid`) USING BTREE,
    INDEX     `org_id`(`org_id`) USING BTREE
);
CREATE TABLE `t_pipeline`
(
    `id`            varchar(64) NOT NULL,
    `uid`           varchar(64)  DEFAULT NULL,
    `name`          varchar(255) DEFAULT NULL,
    `display_name`  varchar(255) DEFAULT NULL,
    `pipeline_type` varchar(255) DEFAULT NULL,
    `json_content`  longtext,
    `yml_content`   longtext,
    `access_token`  varchar(255) DEFAULT NULL,
    `url`           varchar(255) DEFAULT NULL,
    `username`      varchar(255) DEFAULT NULL,
    PRIMARY KEY (`id`) USING BTREE
)
CREATE TABLE `t_pipeline_version`
(
    `id`                    varchar(64) NOT NULL,
    `number`                bigint(20) DEFAULT NULL COMMENT '构建次数',
    `events`                varchar(100) DEFAULT NULL COMMENT '事件push、pr、note',
    `branch`                varchar(255) DEFAULT NULL,
    `sha`                   varchar(255) DEFAULT NULL,
    `pipeline_name`         varchar(255) DEFAULT NULL,
    `pipeline_display_name` varchar(255) DEFAULT NULL,
    `pipeline_id`           varchar(64)  DEFAULT NULL,
    `version`               varchar(255) DEFAULT NULL,
    `content`               longtext,
    `created`               datetime     DEFAULT NULL,
    `deleted`               tinyint(1) DEFAULT '0',
    `pr_number`             bigint(20) DEFAULT NULL,
    `repo_clone_url`        varchar(255) DEFAULT NULL,
    PRIMARY KEY (`id`) USING BTREE
);
CREATE TABLE `t_param`
(
    `aid`   bigint(20) NOT NULL AUTO_INCREMENT,
    `name`  varchar(100) NULL DEFAULT NULL,
    `title` varchar(255) NULL DEFAULT NULL,
    `data`  text NULL,
    `times` datetime(0) NULL DEFAULT NULL,
    PRIMARY KEY (`aid`) USING BTREE
);
CREATE TABLE `t_user`
(
    `id`         varchar(64) NOT NULL,
    `aid`        BIGINT      NOT NULL AUTO_INCREMENT,
    `name`       varchar(100) NULL DEFAULT NULL,
    `pass`       varchar(255) NULL DEFAULT NULL,
    `nick`       varchar(100) NULL DEFAULT NULL,
    `avatar`     varchar(500) NULL DEFAULT NULL,
    `created`    datetime(0) NULL DEFAULT NULL,
    `login_time` datetime(0) NULL DEFAULT NULL,
    PRIMARY KEY (`aid`, `id`) USING BTREE
);
-- ----------------------------
-- Records of t_user
-- ----------------------------
INSERT INTO `t_user`
VALUES ("admin",
        1,
        'admin',
        'e10adc3949ba59abbe56e057f20f883e',
        '管理员',
        NULL,
        NOW(),
        NULL);
CREATE TABLE `t_user_org`
(
    `aid`       bigint(20) NOT NULL AUTO_INCREMENT,
    `uid`       varchar(64) NULL DEFAULT NULL,
    `org_id`    varchar(64) NULL DEFAULT NULL,
    `created`   datetime(0) NULL DEFAULT NULL,
    `perm_adm`  INT(1) NULL DEFAULT 0 COMMENT '管理员',
    `perm_rw`   INT(1) NULL DEFAULT 0 COMMENT '编辑权限',
    `perm_exec` INT(1) NULL DEFAULT 0 COMMENT '执行权限',
    PRIMARY KEY (`aid`) USING BTREE,
    INDEX       `uid`(`uid`) USING BTREE,
    INDEX       `oid`(`org_id`) USING BTREE,
    INDEX       `uoid`(`uid`, `org_id`) USING BTREE
);
CREATE TABLE `t_user_msg`
(
    `aid`          BIGINT NOT NULL AUTO_INCREMENT,
    `uid`          varchar(64) NULL DEFAULT NULL COMMENT '收件人',
    `msg_id`       varchar(64) NULL DEFAULT NULL,
    `created`      datetime(0) NULL DEFAULT NULL,
    `readtm`       datetime(0) NULL DEFAULT NULL,
    `status`       int(11) NULL DEFAULT 0,
    `deleted`      int(1) NULL DEFAULT 0,
    `deleted_time` datetime(0) NULL DEFAULT NULL,
    PRIMARY KEY (`aid`) USING BTREE,
    INDEX          `uid`(`uid`) USING BTREE
);
CREATE TABLE `t_user_token`
(
    `aid`           bigint(20) NOT NULL AUTO_INCREMENT,
    `uid`           bigint(20) NULL DEFAULT NULL,
    `type`          varchar(50) NULL DEFAULT NULL,
    `openid`        varchar(100) NULL DEFAULT NULL,
    `name`          varchar(255) NULL DEFAULT NULL,
    `nick`          varchar(255) NULL DEFAULT NULL,
    `avatar`        varchar(500) NULL DEFAULT NULL,
    `access_token`  text NULL DEFAULT NULL,
    `refresh_token` text NULL DEFAULT NULL,
    `expires_in`    bigint(20) NULL DEFAULT 0,
    `expires_time`  datetime(0) NULL DEFAULT NULL,
    `refresh_time`  datetime(0) NULL DEFAULT NULL,
    `created`       datetime(0) NULL DEFAULT NULL,
    `tokens`        text NULL,
    `uinfos`        text NULL,
    PRIMARY KEY (`aid`) USING BTREE,
    INDEX           `uid`(`uid`) USING BTREE,
    INDEX           `openid`(`openid`) USING BTREE
);
CREATE TABLE `t_repo`
(
    `id`   varchar(64) NOT NULL,
    `name` varchar(255) DEFAULT NULL,
    `url`  varchar(255) DEFAULT NULL,
    PRIMARY KEY (`id`)
);
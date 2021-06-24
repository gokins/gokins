
-- ----------------------------
-- Table structure for t_build
-- ----------------------------
CREATE TABLE `t_build`  (
                            `id` varchar(64) NOT NULL,
                            `pipeline_id` varchar(64) NULL DEFAULT NULL,
                            `pipeline_version_id` varchar(64) NULL DEFAULT NULL,
                            `status` varchar(100) NULL DEFAULT NULL COMMENT '构建状态',
                            `error` varchar(500) NULL DEFAULT NULL COMMENT '错误信息',
                            `event` varchar(100) NULL DEFAULT NULL COMMENT '事件',
                            `time_stamp` datetime(0) NULL DEFAULT NULL COMMENT '执行时长',
                            `title` varchar(255) NULL DEFAULT NULL COMMENT '标题',
                            `message` varchar(255) NULL DEFAULT NULL COMMENT '构建信息',
                            `started` datetime(0) NULL DEFAULT NULL COMMENT '开始时间',
                            `finished` datetime(0) NULL DEFAULT NULL COMMENT '结束时间',
                            `created` datetime(0) NULL DEFAULT NULL COMMENT '创建时间',
                            `updated` datetime(0) NULL DEFAULT NULL COMMENT '更新时间',
                            `version` varchar(255) NULL DEFAULT NULL COMMENT '版本',
                            PRIMARY KEY (`id`) USING BTREE
);

CREATE TABLE `t_cmd_line`  (
  `id` varchar(64) NOT NULL,
  `group_id` varchar(64) NULL DEFAULT NULL,
  `build_id` varchar(64) NULL DEFAULT NULL,
  `job_id` varchar(64) NULL DEFAULT NULL,
  `status` varchar(50) NULL DEFAULT NULL,
  `num` int(11) NULL DEFAULT NULL,
  `content` text NULL,
  `created` datetime(0) NULL DEFAULT NULL,
  `started` datetime(0) NULL DEFAULT NULL,
  `finished` datetime(0) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
);

CREATE TABLE `t_hook`  (
                           `id` varchar(64) NOT NULL,
                           `type` varchar(255) NULL DEFAULT NULL,
                           `snapshot` longtext NULL,
                           `status` varchar(255) NULL DEFAULT NULL,
                           `msg` varchar(255) NULL DEFAULT NULL,
                           `hook_type` varchar(255) NULL DEFAULT NULL,
                           PRIMARY KEY (`id`) USING BTREE
);

-- ----------------------------
-- Records of t_hook
-- ----------------------------
-- ----------------------------
-- Table structure for t_job
-- ----------------------------
CREATE TABLE `t_job`  (
                          `id` varchar(64) NOT NULL,
                          `build_id` varchar(64) NULL DEFAULT NULL,
                          `stage_id` varchar(100) NULL DEFAULT NULL COMMENT '流水线id',
                          `display_name` varchar(255) NULL DEFAULT NULL,
                          `pipeline_version_id` varchar(64) NULL DEFAULT NULL COMMENT '流水线id',
                          `job` varchar(255) NULL DEFAULT NULL,
                          `status` varchar(100) NULL DEFAULT NULL COMMENT '构建状态',
                          `exit_code` bigint(20) NULL DEFAULT NULL COMMENT '退出码',
                          `error` varchar(500) NULL DEFAULT NULL COMMENT '错误信息',
                          `name` varchar(100) NULL DEFAULT NULL COMMENT '名字',
                          `started` datetime(0) NULL DEFAULT NULL COMMENT '开始时间',
                          `finished` datetime(0) NULL DEFAULT NULL COMMENT '结束时间',
                          `created` datetime(0) NULL DEFAULT NULL COMMENT '创建时间',
                          `updated` datetime(0) NULL DEFAULT NULL COMMENT '更新时间',
                          `version` varchar(255) NULL DEFAULT NULL COMMENT '版本',
                          `errignore` varchar(5) NULL DEFAULT NULL,
                          `number` bigint(20) NULL DEFAULT NULL,
                          `commands` text NULL,
                          `depends_on` json NULL,
                          `image` varchar(255) NULL DEFAULT NULL,
                          `environments` json NULL,
                          `sort` bigint(10) NULL DEFAULT NULL,
                          PRIMARY KEY (`id`) USING BTREE
);

-- ----------------------------
-- Records of t_job
-- ----------------------------
-- ----------------------------
-- Table structure for t_message
-- ----------------------------
CREATE TABLE `t_message`  (
                              `id` bigint(20) NOT NULL AUTO_INCREMENT,
                              `xid` varchar(64) NOT NULL,
                              `uid` varchar(64) NULL DEFAULT NULL COMMENT '发送者（可空）',
                              `title` varchar(255) NULL DEFAULT NULL,
                              `content` longtext NULL,
                              `types` varchar(50) NULL DEFAULT NULL,
                              `created` datetime(0) NULL DEFAULT NULL,
                              `infos` text NULL,
                              `url` varchar(500) NULL DEFAULT NULL,
                              PRIMARY KEY (`id`, `xid`) USING BTREE
);

-- ----------------------------
-- Records of t_message
-- ----------------------------
-- ----------------------------
-- Table structure for t_param
-- ----------------------------
CREATE TABLE `t_param`  (
                            `id` bigint(20) NOT NULL AUTO_INCREMENT,
                            `name` varchar(100) NULL DEFAULT NULL,
                            `title` varchar(255) NULL DEFAULT NULL,
                            `data` text NULL,
                            `times` datetime(0) NULL DEFAULT NULL,
                            PRIMARY KEY (`id`) USING BTREE
);

-- ----------------------------
-- Records of t_param
-- ----------------------------
-- ----------------------------
-- Table structure for t_permssion
-- ----------------------------
CREATE TABLE `t_permssion`  (
                                `id` bigint(20) NOT NULL AUTO_INCREMENT,
                                `xid` varchar(64) NOT NULL,
                                `parent` varchar(64) NULL DEFAULT NULL,
                                `title` varchar(100) NULL DEFAULT NULL,
                                `value` varchar(100) NULL DEFAULT NULL,
                                `times` datetime(0) NULL DEFAULT CURRENT_TIMESTAMP(0),
                                `sort` int(11) NULL DEFAULT 10,
                                PRIMARY KEY (`id`, `xid`) USING BTREE,
                                INDEX `IDX_sys_permssion_parent`(`parent`) USING BTREE
);

-- ----------------------------
-- Records of t_permssion
-- ----------------------------
INSERT INTO `t_permssion` VALUES (1, 'common', NULL, '通用', 'common', '2019-10-23 11:29:32', 0);
INSERT INTO `t_permssion` VALUES (2, 'login', 'common', '登录', 'login', '2019-10-23 11:29:40', 11);
INSERT INTO `t_permssion` VALUES (3, 'uppass', 'common', '修改密码', 'comm:uppass', '2020-03-01 21:47:42', 12);
INSERT INTO `t_permssion` VALUES (4, 'admin', NULL, '后台界面', 'admin', '2019-10-27 00:11:30', 2);
INSERT INTO `t_permssion` VALUES (5, 'roles', NULL, '角色管理', 'role:list', '2019-11-03 13:48:58', 3);
INSERT INTO `t_permssion` VALUES (6, 'role1', 'roles', '角色编辑', 'role:edit', '2019-11-03 13:49:21', 11);
INSERT INTO `t_permssion` VALUES (7, 'role2', 'roles', '角色删除', 'role:del', '2019-11-03 13:49:22', 12);
INSERT INTO `t_permssion` VALUES (8, 'grant', 'roles', '用户授权', 'user:grant', '2019-11-04 09:52:47', 13);
INSERT INTO `t_permssion` VALUES (9, 'users', NULL, '用户管理', 'user:list', '2019-11-04 09:51:14', 4);
INSERT INTO `t_permssion` VALUES (10, 'userxg', 'users', '用户编辑', 'user:edit', '2019-11-04 20:27:02', 10);

-- ----------------------------
-- Table structure for t_pipeline
-- ----------------------------
CREATE TABLE `t_pipeline`  (
                               `id` varchar(64) NOT NULL,
                               `name` varchar(255) NULL DEFAULT NULL,
                               `repo_id` varchar(64) NULL DEFAULT NULL,
                               `display_name` varchar(255) NULL DEFAULT NULL,
                               `pipeline_type` varchar(255) NULL DEFAULT NULL,
                               `json_content` longtext NULL,
                               PRIMARY KEY (`id`) USING BTREE
);

-- ----------------------------
-- Table structure for t_role
-- ----------------------------
CREATE TABLE `t_role`  (
                           `id` bigint(20) NOT NULL AUTO_INCREMENT,
                           `xid` varchar(64) NOT NULL,
                           `title` varchar(100) NULL DEFAULT NULL,
                           `perms` text NULL,
                           `times` datetime(0) NULL DEFAULT CURRENT_TIMESTAMP(0),
                           PRIMARY KEY (`id`, `xid`) USING BTREE
);

-- ----------------------------
-- Records of t_role
-- ----------------------------
INSERT INTO `t_role` VALUES (1, 'common', '通用权限', 'common,login,uppass', '2019-10-23 11:31:34');
INSERT INTO `t_role` VALUES (2, 'admin', '权限管理员', 'admin,roles,role1,role2,users,grant', '2019-11-03 13:49:35');

-- ----------------------------
-- Table structure for t_stage
-- ----------------------------
CREATE TABLE `t_stage`  (
                            `id` varchar(64) NOT NULL,
                            `pipeline_version_id` varchar(64) NULL DEFAULT NULL COMMENT '流水线id',
                            `build_id` varchar(64) NULL DEFAULT NULL,
                            `status` varchar(100) NULL DEFAULT NULL COMMENT '构建状态',
                            `exit_code` bigint(20) NULL DEFAULT NULL COMMENT '退出码',
                            `error` varchar(500) NULL DEFAULT NULL COMMENT '错误信息',
                            `name` varchar(255) NULL DEFAULT NULL COMMENT '名字',
                            `display_name` varchar(255) NULL DEFAULT NULL,
                            `started` datetime(0) NULL DEFAULT NULL COMMENT '开始时间',
                            `finished` datetime(0) NULL DEFAULT NULL COMMENT '结束时间',
                            `created` datetime(0) NULL DEFAULT NULL COMMENT '创建时间',
                            `updated` datetime(0) NULL DEFAULT NULL COMMENT '更新时间',
                            `version` varchar(255) NULL DEFAULT NULL COMMENT '版本',
                            `on_success` varchar(5) NULL DEFAULT NULL,
                            `on_failure` varchar(5) NULL DEFAULT NULL,
                            `sort` bigint(10) NULL DEFAULT NULL,
                            `stage` varchar(255) NULL DEFAULT NULL,
                            PRIMARY KEY (`id`) USING BTREE
);

CREATE TABLE `t_variable`  (
                                   `id` varchar(64) NOT NULL,
                                   `name` varchar(255) NULL DEFAULT NULL,
                                   `value` varchar(255) NULL DEFAULT NULL,
                                   PRIMARY KEY (`id`) USING BTREE
);

CREATE TABLE `t_user`  (
                           `id` bigint(20) NOT NULL AUTO_INCREMENT,
                           `name` varchar(100) NULL DEFAULT NULL,
                           `pass` varchar(255) NULL DEFAULT NULL,
                           `nick` varchar(100) NULL DEFAULT NULL,
                           `avatar` varchar(500) NULL DEFAULT NULL,
                           `create_time` datetime(0) NULL DEFAULT NULL,
                           `login_time` datetime(0) NULL DEFAULT NULL,
                           PRIMARY KEY (`id`) USING BTREE
);

-- ----------------------------
-- Records of t_user
-- ----------------------------
INSERT INTO `t_user` VALUES (1, 'admin', 'e10adc3949ba59abbe56e057f20f883e', '管理员', NULL, NOW(), NULL);
-- ----------------------------
-- Table structure for t_user_msg
-- ----------------------------
CREATE TABLE `t_user_msg`  (
                               `id` bigint(20) NOT NULL AUTO_INCREMENT,
                               `mid` varchar(64) NULL DEFAULT NULL,
                               `uid` varchar(64) NULL DEFAULT NULL COMMENT '收件人',
                               `created` datetime(0) NULL DEFAULT NULL,
                               `readtm` datetime(0) NULL DEFAULT NULL,
                               `status` int(11) NULL DEFAULT 0,
                               `deleted` int(1) NULL DEFAULT 0,
                               `deleted_time` datetime(0) NULL DEFAULT NULL,
                               PRIMARY KEY (`id`) USING BTREE
);

CREATE TABLE `t_user_repo`  (
                                `id` varchar(64) NOT NULL,
                                `repo_id` varchar(64) NOT NULL,
                                `user_id` bigint(20) NOT NULL,
                                PRIMARY KEY (`id`) USING BTREE
);

CREATE TABLE `t_user_role`  (
                                `user_id` bigint(20) NOT NULL,
                                `role_codes` text NULL,
                                `limits` text NULL,
                                PRIMARY KEY (`user_id`) USING BTREE
);

CREATE TABLE `t_user_token`  (
                                 `id` bigint(20) NOT NULL AUTO_INCREMENT,
                                 `uid` bigint(20) NULL DEFAULT NULL,
                                 `type` varchar(50) NULL DEFAULT NULL,
                                 `openid` varchar(100) NULL DEFAULT NULL,
                                 `name` varchar(255) NULL DEFAULT NULL,
                                 `nick` varchar(255) NULL DEFAULT NULL,
                                 `avatar` varchar(500) NULL DEFAULT NULL,
                                 `access_token` text NULL DEFAULT NULL,
                                 `refresh_token` text NULL DEFAULT NULL,
                                 `expires_in` bigint(20) NULL DEFAULT 0,
                                 `expires_time` datetime(0) NULL DEFAULT NULL,
                                 `refresh_time` datetime(0) NULL DEFAULT NULL,
                                 `create_time` datetime(0) NULL DEFAULT NULL,
                                 `tokens` text NULL,
                                 `uinfos` text NULL,
                                 PRIMARY KEY (`id`) USING BTREE
);

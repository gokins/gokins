-- ----------------------------
-- Table structure for t_build
-- ----------------------------
CREATE TABLE `t_build` (
  `id` TEXT NOT NULL PRIMARY KEY,
  `pipeline_id` TEXT DEFAULT NULL,
  `pipeline_version_id` TEXT DEFAULT NULL,
  `status` TEXT DEFAULT NULL COMMENT '构建状态',
  `error` TEXT DEFAULT NULL COMMENT '错误信息',
  `event` TEXT DEFAULT NULL COMMENT '事件',
  `time_stamp` DATETIME DEFAULT NULL COMMENT '执行时长',
  `title` TEXT NULL DEFAULT NULL COMMENT '标题',
  `message` TEXT NULL DEFAULT NULL COMMENT '构建信息',
  `started` DATETIME NULL DEFAULT NULL COMMENT '开始时间',
  `finished` DATETIME NULL DEFAULT NULL COMMENT '结束时间',
  `created` DATETIME NULL DEFAULT NULL COMMENT '创建时间',
  `updated` DATETIME NULL DEFAULT NULL COMMENT '更新时间',
  `version` TEXT NULL DEFAULT NULL COMMENT '版本'
);
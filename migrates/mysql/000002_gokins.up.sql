
CREATE TABLE `t_org_var` (
    `aid` bigint(20) NOT NULL AUTO_INCREMENT,
    `uid` varchar(64) DEFAULT NULL,
    `org_id` varchar(64) DEFAULT NULL,
    `name` varchar(255) DEFAULT NULL,
    `value` text DEFAULT NULL,
    `remarks` varchar(255) DEFAULT NULL,
    `public` int(1) DEFAULT '0' COMMENT '公开',
    PRIMARY KEY (`aid`) USING BTREE
  ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;

ALTER TABLE `t_pipeline_var`
MODIFY COLUMN `value` text DEFAULT NULL;

PRAGMA foreign_keys = false;

-- ----------------------------
-- Table structure for schema_migrations
-- ----------------------------
DROP TABLE IF EXISTS "schema_migrations";
CREATE TABLE "schema_migrations" (
    "version" uint64,
    "dirty" bool
);

-- ----------------------------
-- Records of schema_migrations
-- ----------------------------
INSERT INTO "schema_migrations" VALUES (1, 0);

-- ----------------------------
-- Table structure for t_artifact_package
-- ----------------------------
DROP TABLE IF EXISTS "t_artifact_package";
CREATE TABLE "t_artifact_package" (
    "id" varchar(64) NOT NULL,
    "aid" integer NOT NULL,
    "repo_id" varchar(64) DEFAULT NULL,
    "name" varchar(100) DEFAULT NULL,
    "display_name" varchar(255) DEFAULT NULL,
    "desc" varchar(500) DEFAULT NULL,
    "created" datetime DEFAULT NULL,
    "updated" datetime DEFAULT NULL,
    "deleted" int(1) DEFAULT NULL,
    "deleted_time" datetime DEFAULT NULL,
    PRIMARY KEY ("aid")
);

-- ----------------------------
-- Table structure for t_artifact_version
-- ----------------------------
DROP TABLE IF EXISTS "t_artifact_version";
CREATE TABLE "t_artifact_version" (
    "id" varchar(64) NOT NULL,
    "aid" integer NOT NULL,
    "repo_id" varchar(64) DEFAULT NULL,
    "package_id" varchar(64) DEFAULT NULL,
    "name" varchar(100) DEFAULT NULL,
    "version" varchar(100) DEFAULT NULL,
    "sha" varchar(100) DEFAULT NULL,
    "desc" varchar(500) DEFAULT NULL,
    "preview" int(1) DEFAULT NULL,
    "created" datetime DEFAULT NULL,
    "updated" datetime DEFAULT NULL,
    PRIMARY KEY ("aid")
);

-- ----------------------------
-- Table structure for t_artifactory
-- ----------------------------
DROP TABLE IF EXISTS "t_artifactory";
CREATE TABLE "t_artifactory" (
    "id" varchar(64) NOT NULL,
    "aid" integer NOT NULL,
    "uid" varchar(64) DEFAULT NULL,
    "org_id" varchar(64) DEFAULT NULL,
    "identifier" varchar(50) DEFAULT NULL,
    "name" varchar(200) DEFAULT NULL,
    "disabled" int(1) DEFAULT '0',
    "source" varchar(50) DEFAULT NULL,
    "desc" varchar(500) DEFAULT NULL,
    "logo" varchar(255) DEFAULT NULL,
    "created" datetime DEFAULT NULL,
    "updated" datetime DEFAULT NULL,
    "deleted" int(1) DEFAULT NULL,
    "deleted_time" datetime DEFAULT NULL,
    PRIMARY KEY ("aid")
);

-- ----------------------------
-- Table structure for t_build
-- ----------------------------
DROP TABLE IF EXISTS "t_build";
CREATE TABLE "t_build" (
    "id" TEXT NOT NULL,
    "pipeline_id" TEXT DEFAULT NULL,
    "pipeline_version_id" TEXT DEFAULT NULL,
    "status" TEXT DEFAULT NULL,
    "error" TEXT DEFAULT NULL,
    "event" TEXT DEFAULT NULL,
    "time_stamp" DATETIME DEFAULT NULL,
    "title" TEXT DEFAULT NULL,
    "message" TEXT DEFAULT NULL,
    "started" DATETIME DEFAULT NULL,
    "finished" DATETIME DEFAULT NULL,
    "created" DATETIME DEFAULT NULL,
    "updated" DATETIME DEFAULT NULL,
    "version" TEXT DEFAULT NULL,
    PRIMARY KEY ("id")
);

-- ----------------------------
-- Records of t_build
-- ----------------------------
INSERT INTO "t_build" VALUES ('611a0805ff7b6750a0000003', '611a07efff7b6750a0000001', '611a0805ff7b6750a0000002', 'error', '', '', NULL, NULL, NULL, '2021-08-16 06:39:01', '2021-08-16 06:39:02', '2021-08-16 06:39:01', '2021-08-16 06:39:02', '');
INSERT INTO "t_build" VALUES ('611a0942ff7b6750a000000c', '611a07efff7b6750a0000001', '611a0942ff7b6750a000000b', 'error', '', '', NULL, NULL, NULL, '2021-08-16 06:44:18', '2021-08-16 06:44:19', '2021-08-16 06:44:18', '2021-08-16 06:44:19', '');
INSERT INTO "t_build" VALUES ('611a0a8dff7b6750a0000015', '611a07efff7b6750a0000001', '611a0a8dff7b6750a0000014', 'error', '', '', NULL, NULL, NULL, '2021-08-16 06:49:50', '2021-08-16 06:49:51', '2021-08-16 06:49:49', '2021-08-16 06:49:51', '');
INSERT INTO "t_build" VALUES ('611a0a9cff7b6750a000001e', '611a07efff7b6750a0000001', '611a0a9cff7b6750a000001d', 'error', '', '', NULL, NULL, NULL, '2021-08-16 06:50:04', '2021-08-16 06:50:05', '2021-08-16 06:50:04', '2021-08-16 06:50:05', '');
INSERT INTO "t_build" VALUES ('611a0acaff7b6750a0000027', '611a07efff7b6750a0000001', '611a0acaff7b6750a0000026', 'error', 'repo err:http: unexpected EOF reading trailer', 'err_get_repo', NULL, NULL, NULL, '2021-08-16 06:50:50', '2021-08-16 06:52:15', '2021-08-16 06:50:50', '2021-08-16 06:52:15', '');
INSERT INTO "t_build" VALUES ('611a0b36ff7b6750a0000030', '611a07efff7b6750a0000001', '611a0b36ff7b6750a000002f', 'error', '', '', NULL, NULL, NULL, '2021-08-16 06:52:38', '2021-08-16 06:52:39', '2021-08-16 06:52:38', '2021-08-16 06:52:39', '');

-- ----------------------------
-- Table structure for t_cmd_line
-- ----------------------------
DROP TABLE IF EXISTS "t_cmd_line";
CREATE TABLE "t_cmd_line" (
    "id" varchar(64) NOT NULL,
    "group_id" varchar(64) DEFAULT NULL,
    "build_id" varchar(64) DEFAULT NULL,
    "step_id" varchar(64) DEFAULT NULL,
    "status" varchar(50) DEFAULT NULL,
    "num" int(11) DEFAULT NULL,
    "code" int(11) DEFAULT NULL,
    "content" text,
    "created" datetime DEFAULT NULL,
    "started" datetime DEFAULT NULL,
    "finished" datetime DEFAULT NULL,
    PRIMARY KEY ("id")
);

-- ----------------------------
-- Table structure for t_message
-- ----------------------------
DROP TABLE IF EXISTS "t_message";
CREATE TABLE "t_message" (
    "id" varchar(64) NOT NULL,
    "aid" integer NOT NULL,
    "uid" varchar(64) DEFAULT NULL,
    "title" varchar(255) DEFAULT NULL,
    "content" longtext,
    "types" varchar(50) DEFAULT NULL,
    "created" datetime DEFAULT NULL,
    "infos" text,
    "url" varchar(500) DEFAULT NULL,
    PRIMARY KEY ("aid")
);

-- ----------------------------
-- Table structure for t_org
-- ----------------------------
DROP TABLE IF EXISTS "t_org";
CREATE TABLE "t_org" (
    "id" varchar(64) NOT NULL,
    "aid" integer NOT NULL,
    "uid" varchar(64) DEFAULT NULL,
    "name" varchar(200) DEFAULT NULL,
    "desc" TEXT DEFAULT NULL,
    "public" INT(1) DEFAULT 0,
    "created" datetime DEFAULT NULL,
    "updated" datetime DEFAULT NULL,
    "deleted" int(1) DEFAULT 0,
    "deleted_time" datetime DEFAULT NULL,
    PRIMARY KEY ("aid")
);

-- ----------------------------
-- Table structure for t_org_pipe
-- ----------------------------
DROP TABLE IF EXISTS "t_org_pipe";
CREATE TABLE "t_org_pipe" (
    "aid" integer NOT NULL,
    "org_id" varchar(64) DEFAULT NULL,
    "pipe_id" varchar(64) DEFAULT NULL,
    "created" datetime DEFAULT NULL,
    "public" INT(1) DEFAULT 0,
    PRIMARY KEY ("aid")
);

-- ----------------------------
-- Table structure for t_param
-- ----------------------------
DROP TABLE IF EXISTS "t_param";
CREATE TABLE "t_param" (
    "aid" integer NOT NULL,
    "name" varchar(100) DEFAULT NULL,
    "title" varchar(255) DEFAULT NULL,
    "data" text,
    "times" datetime DEFAULT NULL,
    PRIMARY KEY ("aid")
);

-- ----------------------------
-- Table structure for t_pipeline
-- ----------------------------
DROP TABLE IF EXISTS "t_pipeline";
CREATE TABLE "t_pipeline" (
    "id" varchar(64) NOT NULL,
    "uid" varchar(64) DEFAULT NULL,
    "name" varchar(255) DEFAULT NULL,
    "display_name" varchar(255) DEFAULT NULL,
    "pipeline_type" varchar(255) DEFAULT NULL,
    "created" datetime DEFAULT NULL,
    "deleted" int(1) DEFAULT '0',
    "deleted_time" datetime DEFAULT NULL,
    PRIMARY KEY ("id")
);

-- ----------------------------
-- Table structure for t_pipeline_conf
-- ----------------------------
DROP TABLE IF EXISTS "t_pipeline_conf";
CREATE TABLE "t_pipeline_conf" (
    "aid" integer NOT NULL,
    "pipeline_id" varchar(64) NOT NULL,
    "url" varchar(255) DEFAULT NULL,
    "access_token" varchar(255) DEFAULT NULL,
    "yml_content" longtext,
    "username" varchar(255) DEFAULT NULL,
    PRIMARY KEY ("aid")
);

-- ----------------------------
-- Table structure for t_pipeline_var
-- ----------------------------
DROP TABLE IF EXISTS "t_pipeline_var";
CREATE TABLE "t_pipeline_var" (
    "aid" integer NOT NULL,
    "uid" varchar(64) DEFAULT NULL,
    "pipeline_id" varchar(64) DEFAULT NULL,
    "name" varchar(255) DEFAULT NULL,
    "value" varchar(255) DEFAULT NULL,
    "remarks" varchar(255) DEFAULT NULL,
    "public" int(1) DEFAULT '0',
    PRIMARY KEY ("aid")
);

-- ----------------------------
-- Table structure for t_pipeline_version
-- ----------------------------
DROP TABLE IF EXISTS "t_pipeline_version";
CREATE TABLE "t_pipeline_version" (
    "id" varchar(64) NOT NULL,
    "uid" varchar(64) DEFAULT NULL,
    "number" bigint(20) DEFAULT NULL,
    "events" varchar(100) DEFAULT NULL,
    "sha" varchar(255) DEFAULT NULL,
    "pipeline_name" varchar(255) DEFAULT NULL,
    "pipeline_display_name" varchar(255) DEFAULT NULL,
    "pipeline_id" varchar(64) DEFAULT NULL,
    "version" varchar(255) DEFAULT NULL,
    "content" longtext,
    "created" datetime DEFAULT NULL,
    "deleted" tinyint(1) DEFAULT '0',
    "pr_number" bigint(20) DEFAULT NULL,
    "repo_clone_url" varchar(255) DEFAULT NULL,
    PRIMARY KEY ("id")
);

-- ----------------------------
-- Table structure for t_stage
-- ----------------------------
DROP TABLE IF EXISTS "t_stage";
CREATE TABLE "t_stage" (
    "id" varchar(64) NOT NULL,
    "pipeline_version_id" varchar(64) DEFAULT NULL,
    "build_id" varchar(64) DEFAULT NULL,
    "status" varchar(100) DEFAULT NULL,
    "error" varchar(500) DEFAULT NULL,
    "name" varchar(255) DEFAULT NULL,
    "display_name" varchar(255) DEFAULT NULL,
    "started" datetime DEFAULT NULL,
    "finished" datetime DEFAULT NULL,
    "created" datetime DEFAULT NULL,
    "updated" datetime DEFAULT NULL,
    "sort" int(11) DEFAULT NULL,
    "stage" varchar(255) DEFAULT NULL,
    PRIMARY KEY ("id")
);

-- ----------------------------
-- Table structure for t_step
-- ----------------------------
DROP TABLE IF EXISTS "t_step";
CREATE TABLE "t_step" (
    "id" varchar(64) NOT NULL,
    "build_id" varchar(64) DEFAULT NULL,
    "stage_id" varchar(100) DEFAULT NULL,
    "display_name" varchar(255) DEFAULT NULL,
    "pipeline_version_id" varchar(64) DEFAULT NULL,
    "step" varchar(255) DEFAULT NULL,
    "status" varchar(100) DEFAULT NULL,
    "event" varchar(100) DEFAULT NULL,
    "exit_code" int(11) DEFAULT NULL,
    "error" varchar(500) DEFAULT NULL,
    "name" varchar(100) DEFAULT NULL,
    "started" datetime DEFAULT NULL,
    "finished" datetime DEFAULT NULL,
    "created" datetime DEFAULT NULL,
    "updated" datetime DEFAULT NULL,
    "version" varchar(255) DEFAULT NULL,
    "errignore" int(11) DEFAULT NULL,
    "commands" text,
    "waits" json,
    "sort" int(11) DEFAULT NULL,
    PRIMARY KEY ("id")
);

-- ----------------------------
-- Table structure for t_trigger
-- ----------------------------
DROP TABLE IF EXISTS "t_trigger";
CREATE TABLE "t_trigger" (
    "id" varchar(64) NOT NULL,
    "aid" integer NOT NULL,
    "uid" varchar(64) DEFAULT NULL,
    "pipeline_id" varchar(64) NOT NULL,
    "types" varchar(50) DEFAULT NULL,
    "name" varchar(100) DEFAULT NULL,
    "desc" varchar(255) DEFAULT NULL,
    "params" json DEFAULT NULL,
    "enabled" int(1) DEFAULT NULL,
    "created" datetime DEFAULT NULL,
    "updated" datetime DEFAULT NULL,
    PRIMARY KEY ("aid")
);

-- ----------------------------
-- Table structure for t_trigger_run
-- ----------------------------
DROP TABLE IF EXISTS "t_trigger_run";
CREATE TABLE "t_trigger_run" (
    "id" varchar(64) NOT NULL,
    "aid" integer NOT NULL,
    "tid" varchar(64) DEFAULT NULL,
    "pipe_version_id" varchar(64) DEFAULT NULL,
    "infos" json DEFAULT NULL,
    "error" varchar(255) DEFAULT NULL,
    "created" datetime DEFAULT NULL,
    PRIMARY KEY ("aid")
);

-- ----------------------------
-- Table structure for t_user
-- ----------------------------
DROP TABLE IF EXISTS "t_user";
CREATE TABLE "t_user" (
    "id" varchar(64) NOT NULL,
    "aid" integer NOT NULL,
    "name" varchar(100) DEFAULT NULL,
    "pass" varchar(255) DEFAULT NULL,
    "nick" varchar(100) DEFAULT NULL,
    "avatar" varchar(500) DEFAULT NULL,
    "created" datetime DEFAULT NULL,
    "login_time" datetime DEFAULT NULL,
    "active" int(1) DEFAULT '0',
    PRIMARY KEY ("aid")
);

-- ----------------------------
-- Records of t_user
-- ----------------------------
INSERT INTO "t_user" VALUES ('admin', 1, 'gokins', 'e10adc3949ba59abbe56e057f20f883e', '管理员', NULL, 1628584680, '2021-08-16 06:09:57', 0);

-- ----------------------------
-- Table structure for t_user_info
-- ----------------------------
DROP TABLE IF EXISTS "t_user_info";
CREATE TABLE "t_user_info" (
    "id" varchar(64) NOT NULL,
    "phone" varchar(100) DEFAULT NULL,
    "email" varchar(200) DEFAULT NULL,
    "birthday" datetime DEFAULT NULL,
    "remark" text,
    "perm_user" int(1) DEFAULT NULL,
    "perm_org" int(1) DEFAULT NULL,
    "perm_pipe" int(1) DEFAULT NULL,
    PRIMARY KEY ("id")
);

-- ----------------------------
-- Records of t_user_info
-- ----------------------------

-- ----------------------------
-- Table structure for t_user_msg
-- ----------------------------
DROP TABLE IF EXISTS "t_user_msg";
CREATE TABLE "t_user_msg" (
    "aid" integer NOT NULL,
    "uid" varchar(64) DEFAULT NULL,
    "msg_id" varchar(64) DEFAULT NULL,
    "created" datetime DEFAULT NULL,
    "readtm" datetime DEFAULT NULL,
    "status" int(11) DEFAULT 0,
    "deleted" int(1) DEFAULT 0,
    "deleted_time" datetime DEFAULT NULL,
    PRIMARY KEY ("aid")
);

-- ----------------------------
-- Table structure for t_user_org
-- ----------------------------
DROP TABLE IF EXISTS "t_user_org";
CREATE TABLE "t_user_org" (
    "aid" integer NOT NULL,
    "uid" varchar(64) DEFAULT NULL,
    "org_id" varchar(64) DEFAULT NULL,
    "created" datetime DEFAULT NULL,
    "perm_adm" INT(1) DEFAULT 0,
    "perm_rw" INT(1) DEFAULT 0,
    "perm_exec" INT(1) DEFAULT 0,
    "perm_down" int(1) DEFAULT NULL,
    PRIMARY KEY ("aid")
);

-- ----------------------------
-- Table structure for t_user_token
-- ----------------------------
DROP TABLE IF EXISTS "t_user_token";
CREATE TABLE "t_user_token" (
    "aid" integer NOT NULL,
    "uid" bigint(20) DEFAULT NULL,
    "type" varchar(50) DEFAULT NULL,
    "openid" varchar(100) DEFAULT NULL,
    "name" varchar(255) DEFAULT NULL,
    "nick" varchar(255) DEFAULT NULL,
    "avatar" varchar(500) DEFAULT NULL,
    "access_token" text DEFAULT NULL,
    "refresh_token" text DEFAULT NULL,
    "expires_in" bigint(20) DEFAULT 0,
    "expires_time" datetime DEFAULT NULL,
    "refresh_time" datetime DEFAULT NULL,
    "created" datetime DEFAULT NULL,
    "tokens" text,
    "uinfos" text,
    PRIMARY KEY ("aid")
);

-- ----------------------------
-- Table structure for t_yml_plugin
-- ----------------------------
DROP TABLE IF EXISTS "t_yml_plugin";
CREATE TABLE "t_yml_plugin" (
    "aid" integer NOT NULL,
    "name" varchar(64) DEFAULT NULL,
    "yml_content" longtext,
    "deleted" int(1) DEFAULT '0',
    "deleted_time" datetime DEFAULT NULL,
    PRIMARY KEY ("aid")
);

-- ----------------------------
-- Records of t_yml_plugin
-- ----------------------------
INSERT INTO "t_yml_plugin" VALUES (1, 'sh', '      - step: shell@sh\n        displayName: sh\n        name: sh\n        commands:\n          - echo hello world', 0, NULL);
INSERT INTO "t_yml_plugin" VALUES (2, 'bash', '      - step: shell@bash\n        displayName: bash\n        name: bash\n        commands:\n          - echo hello world', 0, NULL);
INSERT INTO "t_yml_plugin" VALUES (3, 'powershell', '      - step: shell@powershell\n        displayName: powershell\n        name: powershell\n        commands:\n          - echo hello world', 1, NULL);
INSERT INTO "t_yml_plugin" VALUES (4, 'ssh', '      - step: shell@ssh\r\n        displayName: ssh\r\n        name: ssh\r\n        input:\r\n          host: localhost:22  #端口必填\r\n          user: root\r\n          pass: 123456\r\n          workspace: /root/test #为空就是 $HOME 用户目录\r\n        commands:\r\n          - echo hello world', 0, NULL);

-- ----------------------------
-- Table structure for t_yml_template
-- ----------------------------
DROP TABLE IF EXISTS "t_yml_template";
CREATE TABLE "t_yml_template" (
    "aid" integer NOT NULL,
    "name" varchar(64) DEFAULT NULL,
    "yml_content" longtext,
    "deleted" int(1) DEFAULT '0',
    "deleted_time" datetime DEFAULT NULL,
    PRIMARY KEY ("aid")
);

-- ----------------------------
-- Records of t_yml_template
-- ----------------------------
INSERT INTO "t_yml_template" VALUES (1, 'Golang', 'version: 1.0\nvars:\nstages:\n  - stage:\n    displayName: build\n    name: build\n    steps:\n      - step: shell@sh\n        displayName: go-build-1\n        name: build\n        env:\n        commands:\n          - go build main.go\n      - step: shell@sh\n        displayName: go-build-2\n        name: test\n        env:\n        commands:\n          - go test -v\n', 0, NULL);
INSERT INTO "t_yml_template" VALUES (2, 'Maven', 'version: 1.0\nvars:\nstages:\n  - stage:\n    displayName: build\n    name: build\n    steps:\n      - step: shell@sh\n        displayName: java-build-1\n        name: build\n        env:\n        commands:\n          - mvn clean\n          - mvn install\n      - step: shell@sh\n        displayName: java-build-2\n        name: test\n        env:\n        commands:\n          - mvn test -v', 0, NULL);
INSERT INTO "t_yml_template" VALUES (3, 'Npm', 'version: 1.0\nvars:\nstages:\n  - stage:\n    displayName: build\n    name: build\n    steps:\n      - step: shell@sh\n        displayName: npm-build-1\n        name: build\n        env:\n        commands:\n          - npm build\n      - step: shell@sh\n        displayName: npm-build-2\n        name: publish\n        env:\n        commands:\n          - npm publish ', 0, NULL);

-- ----------------------------
-- Indexes structure for table schema_migrations
-- ----------------------------
CREATE UNIQUE INDEX "version_unique" ON "schema_migrations" ("version" ASC);

PRAGMA foreign_keys = true;

CREATE TABLE "t_artifact_package" (
  "id" varchar(64) NOT NULL,
  "aid" serial8 NOT NULL,
  "repo_id" varchar(64) DEFAULT NULL,
  "name" varchar(100) DEFAULT NULL,
  "display_name" varchar(255) DEFAULT NULL,
  "desc" varchar(500) DEFAULT NULL,
  "created" timestamp DEFAULT NULL,
  "updated" timestamp DEFAULT NULL,
  "deleted" int2 DEFAULT NULL,
  "deleted_time" timestamp DEFAULT NULL,
  PRIMARY KEY ("aid", "id")
);
CREATE INDEX t_artifact_package_pid on t_artifact_package("repo_id");
CREATE INDEX t_artifact_package_rpnm on t_artifact_package("repo_id", "name");

CREATE TABLE "t_artifact_version" (
  "id" varchar(64) NOT NULL,
  "aid" serial8 NOT NULL,
  "repo_id" varchar(64) DEFAULT NULL,
  "package_id" varchar(64) DEFAULT NULL,
  "name" varchar(100) DEFAULT NULL,
  "version" varchar(100) DEFAULT NULL,
  "sha" varchar(100) DEFAULT NULL,
  "desc" varchar(500) DEFAULT NULL,
  "preview" int2 DEFAULT NULL,
  "created" timestamp DEFAULT NULL,
  "updated" timestamp DEFAULT NULL,
  PRIMARY KEY ("aid", "id")
);

CREATE INDEX t_artifact_version_rpnm on t_artifact_version("repo_id", "name");


CREATE TABLE "t_artifactory" (
  "id" varchar(64) NOT NULL,
  "aid" serial8 NOT NULL,
  "uid" varchar(64) DEFAULT NULL,
  "org_id" varchar(64) DEFAULT NULL,
  "identifier" varchar(50) DEFAULT NULL,
  "name" varchar(200) DEFAULT NULL,
  "disabled" int2 DEFAULT 0,
  "source" varchar(50) DEFAULT NULL,
  "desc" varchar(500) DEFAULT NULL,
  "logo" varchar(255) DEFAULT NULL,
  "created" timestamp DEFAULT NULL,
  "updated" timestamp DEFAULT NULL,
  "deleted" int2 DEFAULT NULL,
  "deleted_time" timestamp DEFAULT NULL,
  PRIMARY KEY ("aid", "id")
);
COMMENT on column "t_artifactory"."disabled" IS '是否归档(1归档|0正常)';

CREATE TABLE "t_build" (
  "id" varchar(64) NOT NULL,
  "pipeline_id" varchar(64) NULL DEFAULT NULL,
  "pipeline_version_id" varchar(64) NULL DEFAULT NULL,
  "status" varchar(100) NULL DEFAULT NULL,
  "error" varchar(500) NULL DEFAULT NULL,
  "event" varchar(100) NULL DEFAULT NULL,
  "started" timestamp NULL DEFAULT NULL,
  "finished" timestamp NULL DEFAULT NULL,
  "created" timestamp NULL DEFAULT NULL,
  "updated" timestamp NULL DEFAULT NULL,
  "version" varchar(255) NULL DEFAULT NULL,
  PRIMARY KEY ("id")
) ;

COMMENT on column "t_build"."status" IS '构建状态';
COMMENT on column "t_build"."error" IS '错误信息';
COMMENT on column "t_build"."event" IS '事件';
COMMENT on column "t_build"."started" IS '开始时间';
COMMENT on column "t_build"."finished" IS '结束时间';
COMMENT on column "t_build"."created" IS '创建时间';
COMMENT on column "t_build"."updated" IS '更新时间';
COMMENT on column "t_build"."version" IS '版本';


CREATE TABLE "t_cmd_line" (
  "id" varchar(64) NOT NULL,
  "group_id" varchar(64) NULL DEFAULT NULL,
  "build_id" varchar(64) NULL DEFAULT NULL,
  "step_id" varchar(64) NULL DEFAULT NULL,
  "status" varchar(50) NULL DEFAULT NULL,
  "num" int4 NULL DEFAULT NULL,
  "code" int4 NULL DEFAULT NULL,
  "content" text NULL,
  "created" timestamp NULL DEFAULT NULL,
  "started" timestamp NULL DEFAULT NULL,
  "finished" timestamp NULL DEFAULT NULL,
  PRIMARY KEY ("id")
);


CREATE TABLE "t_stage" (
  "id" varchar(64) NOT NULL,
  "pipeline_version_id" varchar(64) NULL DEFAULT NULL,
  "build_id" varchar(64) NULL DEFAULT NULL,
  "status" varchar(100) NULL DEFAULT NULL,
  "error" varchar(500) NULL DEFAULT NULL,
  "name" varchar(255) NULL DEFAULT NULL,
  "display_name" varchar(255) NULL DEFAULT NULL,
  "started" timestamp NULL DEFAULT NULL,
  "finished" timestamp NULL DEFAULT NULL,
  "created" timestamp NULL DEFAULT NULL,
  "updated" timestamp NULL DEFAULT NULL,
  "sort" int2 NULL DEFAULT NULL,
  "stage" varchar(255) NULL DEFAULT NULL,
  PRIMARY KEY ("id")
);
COMMENT on column "t_stage"."pipeline_version_id" IS '流水线id';
COMMENT on column "t_stage"."status" IS '构建状态';
COMMENT on column "t_stage"."error" IS '错误信息';
COMMENT on column "t_stage"."name" IS '名字';
COMMENT on column "t_stage"."started" IS '开始时间';
COMMENT on column "t_stage"."finished" IS '结束时间';
COMMENT on column "t_stage"."created" IS '创建时间';
COMMENT on column "t_stage"."updated" IS '更新时间';


CREATE TABLE "t_step" (
  "id" varchar(64) NOT NULL,
  "build_id" varchar(64) NULL DEFAULT NULL,
  "stage_id" varchar(100) NULL DEFAULT NULL,
  "display_name" varchar(255) NULL DEFAULT NULL,
  "pipeline_version_id" varchar(64) NULL DEFAULT NULL,
  "step" varchar(255) NULL DEFAULT NULL,
  "status" varchar(100) NULL DEFAULT NULL,
  "event" varchar(100) NULL DEFAULT NULL,
  "exit_code" int8 NULL DEFAULT NULL,
  "error" varchar(500) NULL DEFAULT NULL,
  "name" varchar(100) NULL DEFAULT NULL,
  "started" timestamp NULL DEFAULT NULL,
  "finished" timestamp NULL DEFAULT NULL,
  "created" timestamp NULL DEFAULT NULL,
  "updated" timestamp NULL DEFAULT NULL,
  "version" varchar(255) NULL DEFAULT NULL,
  "errignore" int4 NULL DEFAULT NULL,
  "commands" text NULL,
  "waits" jsonb NULL,
  "sort" int2 NULL DEFAULT NULL,
  PRIMARY KEY ("id")
);
COMMENT on column "t_step"."pipeline_version_id" IS '流水线id';
COMMENT on column "t_step"."stage_id" IS '流水线id';
COMMENT on column "t_step"."status" IS '构建状态';
COMMENT on column "t_step"."event" IS '事件';
COMMENT on column "t_step"."exit_code" IS '退出码';
COMMENT on column "t_step"."error" IS '错误信息';
COMMENT on column "t_step"."name" IS '名字';
COMMENT on column "t_step"."started" IS '开始时间';
COMMENT on column "t_step"."finished" IS '结束时间';
COMMENT on column "t_step"."created" IS '创建时间';
COMMENT on column "t_step"."updated" IS '更新时间';



CREATE TABLE "t_trigger" (
  "id" varchar(64) NOT NULL,
  "aid" serial8 NOT NULL,
  "uid" varchar(64) DEFAULT NULL,
  "pipeline_id" varchar(64) NOT NULL,
  "types" varchar(50) DEFAULT NULL,
  "name" varchar(100) DEFAULT NULL,
  "desc" varchar(255) DEFAULT NULL,
  "params" jsonb DEFAULT NULL,
  "enabled" int2 DEFAULT NULL,
  "created" timestamp DEFAULT NULL,
  "updated" timestamp DEFAULT NULL,
  PRIMARY KEY ("aid", "id")
);


CREATE INDEX t_trigger_uid on t_trigger("uid");


CREATE TABLE "t_trigger_run" (
  "id" varchar(64) NOT NULL,
  "aid" serial8 NOT NULL,
  "tid" varchar(64) DEFAULT NULL,
  "pipe_version_id" varchar(64) DEFAULT NULL,
  "infos" jsonb DEFAULT NULL,
  "error" varchar(255) DEFAULT NULL,
  "created" timestamp DEFAULT NULL,
  PRIMARY KEY ("aid", "id")
);
CREATE INDEX t_trigger_run_tid on t_trigger_run("tid");
COMMENT on column "t_trigger_run"."tid" IS '触发器ID';


CREATE TABLE "t_message" (
  "id" varchar(64) NOT NULL,
  "aid" serial8 NOT NULL,
  "uid" varchar(64) NULL DEFAULT NULL,
  "title" varchar(255) NULL DEFAULT NULL,
  "content" text NULL,
  "types" varchar(50) NULL DEFAULT NULL,
  "created" timestamp NULL DEFAULT NULL,
  "infos" text NULL,
  "url" varchar(500) NULL DEFAULT NULL,
  PRIMARY KEY ("aid", "id")
);
COMMENT on column "t_message"."uid" IS '发送者（可空）';


CREATE TABLE "t_org" (
  "id" varchar(64) NOT NULL,
  "aid" serial8 NOT NULL,
  "uid" varchar(64) NULL DEFAULT NULL,
  "name" varchar(200) NULL DEFAULT NULL,
  "desc" TEXT NULL DEFAULT NULL,
  "public" smallint NULL DEFAULT 0,
  "created" timestamp NULL DEFAULT NULL ,
  "updated" timestamp NULL DEFAULT NULL,
  "deleted" smallint NULL DEFAULT 0,
  "deleted_time" timestamp NULL DEFAULT NULL,
  PRIMARY KEY ("aid", "id")
);
CREATE INDEX uid on t_org("uid");
COMMENT on column "t_org"."public" IS '公开';
COMMENT on column "t_org"."created" IS '创建时间';
COMMENT on column "t_org"."updated" IS '更新时间';

CREATE TABLE "t_org_pipe" (
  "aid" serial8 NOT NULL,
  "org_id" varchar(64) NULL DEFAULT NULL,
  "pipe_id" varchar(64) NULL DEFAULT NULL,
  "created" timestamp NULL DEFAULT NULL,
  "public" smallint NULL DEFAULT 0,
  PRIMARY KEY ("aid")
);
CREATE INDEX t_org_pipe_org_id on t_org_pipe("org_id");
COMMENT on column "t_org_pipe"."pipe_id" IS '收件人';
COMMENT on column "t_org_pipe"."public" IS '公开';

CREATE TABLE "t_pipeline" (
  "id" varchar(64) NOT NULL,
  "uid" varchar(64) DEFAULT NULL,
  "name" varchar(255) DEFAULT NULL,
  "display_name" varchar(255) DEFAULT NULL,
  "pipeline_type" varchar(255) DEFAULT NULL,
  "created" timestamp DEFAULT NULL,
  "deleted" smallint DEFAULT 0,
  "deleted_time" timestamp DEFAULT NULL,
  PRIMARY KEY ("id")
);
CREATE TABLE "t_pipeline_conf" (
  "aid" serial8 NOT NULL,
  "pipeline_id" varchar(64) NOT NULL,
  "url" varchar(255) DEFAULT NULL,
  "access_token" varchar(255) DEFAULT NULL,
  "yml_content" text,
  "username" varchar(255) DEFAULT NULL,
  PRIMARY KEY ("aid")
) ;

CREATE TABLE "t_pipeline_version" (
  "id" varchar(64) NOT NULL,
  "uid" varchar(64) DEFAULT NULL,
  "number" int8 DEFAULT NULL ,
  "events" varchar(100) DEFAULT NULL,
  "sha" varchar(255) DEFAULT NULL,
  "pipeline_name" varchar(255) DEFAULT NULL,
  "pipeline_display_name" varchar(255) DEFAULT NULL,
  "pipeline_id" varchar(64) DEFAULT NULL,
  "version" varchar(255) DEFAULT NULL,
  "content" text,
  "created" timestamp DEFAULT NULL,
  "deleted" smallint DEFAULT 0,
  "pr_number" int8 DEFAULT NULL,
  "repo_clone_url" varchar(255) DEFAULT NULL,
  PRIMARY KEY ("id")
);
COMMENT on column "t_pipeline_version"."number" IS '构建次数';
COMMENT on column "t_pipeline_version"."events" IS '事件push、pr、note';

CREATE TABLE "t_param" (
  "aid" serial8 NOT NULL,
  "name" varchar(100) NULL DEFAULT NULL,
  "title" varchar(255) NULL DEFAULT NULL,
  "data" text NULL,
  "times" timestamp NULL DEFAULT NULL,
  PRIMARY KEY ("aid")
);
CREATE TABLE "t_user" (
  "id" varchar(64) NOT NULL,
  "aid" serial8 NOT NULL,
  "name" varchar(100) NULL DEFAULT NULL,
  "pass" varchar(255) NULL DEFAULT NULL,
  "nick" varchar(100) NULL DEFAULT NULL,
  "avatar" varchar(500) NULL DEFAULT NULL,
  "created" timestamp NULL DEFAULT NULL,
  "login_time" timestamp NULL DEFAULT NULL,
  "active" smallint DEFAULT 0,
  PRIMARY KEY ("aid", "id")
);
-- ----------------------------
INSERT INTO
  "t_user"
VALUES
  (
    'admin',
    1,
    'gokins',
    'e10adc3949ba59abbe56e057f20f883e',
    '管理员',
    NULL,
    NOW(),
    NULL,
    1
  );
-- ----------------------------
  CREATE TABLE "t_user_info" (
    "id" varchar(64) NOT NULL,
    "phone" varchar(100) DEFAULT NULL,
    "email" varchar(200) DEFAULT NULL,
    "birthday" timestamp DEFAULT NULL,
    "remark" text,
    "perm_user" int2 DEFAULT NULL,
    "perm_org" int2 DEFAULT NULL,
    "perm_pipe" int2 DEFAULT NULL,
    PRIMARY KEY ("id")
  );
CREATE TABLE "t_user_org" (
    "aid" serial8 NOT NULL,
    "uid" varchar(64) NULL DEFAULT NULL,
    "org_id" varchar(64) NULL DEFAULT NULL,
    "created" timestamp NULL DEFAULT NULL,
    "perm_adm" int2 NULL DEFAULT 0,
    "perm_rw" int2 NULL DEFAULT 0,
    "perm_exec" int2 NULL DEFAULT 0,
    "perm_down" int2 DEFAULT NULL,
    PRIMARY KEY ("aid")
  );
CREATE INDEX t_user_org_uid on t_user_org("uid");
CREATE INDEX t_user_org_oid on t_user_org("org_id");
CREATE INDEX t_user_org_uoid on t_user_org("uid","org_id");
COMMENT on column "t_user_org"."perm_adm" IS '管理员';
COMMENT on column "t_user_org"."perm_rw" IS '编辑权限';
COMMENT on column "t_user_org"."perm_exec" IS '执行权限';
COMMENT on column "t_user_org"."perm_down" IS '下载制品权限';


CREATE TABLE "t_user_msg" (
    "aid" serial8 NOT NULL,
    "uid" varchar(64) NULL DEFAULT NULL,
    "msg_id" varchar(64) NULL DEFAULT NULL,
    "created" timestamp NULL DEFAULT NULL,
    "readtm" timestamp NULL DEFAULT NULL,
    "status" smallint NULL DEFAULT 0,
    "deleted" smallint NULL DEFAULT 0,
    "deleted_time" timestamp NULL DEFAULT NULL,
    PRIMARY KEY ("aid")
  );
CREATE INDEX t_user_msg_uid on t_user_msg("uid");
COMMENT on column "t_user_msg"."uid" IS '收件人';

CREATE TABLE "t_user_token" (
    "aid" serial8 NOT NULL,
    "uid" int8 NULL DEFAULT NULL,
    "type" varchar(50) NULL DEFAULT NULL,
    "openid" varchar(100) NULL DEFAULT NULL,
    "name" varchar(255) NULL DEFAULT NULL,
    "nick" varchar(255) NULL DEFAULT NULL,
    "avatar" varchar(500) NULL DEFAULT NULL,
    "access_token" text NULL DEFAULT NULL,
    "refresh_token" text NULL DEFAULT NULL,
    "expires_in" int8 NULL DEFAULT 0,
    "expires_time" timestamp NULL DEFAULT NULL,
    "refresh_time" timestamp NULL DEFAULT NULL,
    "created" timestamp NULL DEFAULT NULL,
    "tokens" text NULL,
    "uinfos" text NULL,
    PRIMARY KEY ("aid")
  );
CREATE INDEX t_user_token_uid on t_user_token("uid");
CREATE INDEX t_user_token_openid on t_user_token("openid");
CREATE TABLE "t_pipeline_var" (
    "aid" serial8 NOT NULL,
    "uid" varchar(64) DEFAULT NULL,
    "pipeline_id" varchar(64) DEFAULT NULL,
    "name" varchar(255) DEFAULT NULL,
    "value" varchar(255) DEFAULT NULL,
    "remarks" varchar(255) DEFAULT NULL,
    "public" smallint DEFAULT 0,
    PRIMARY KEY ("aid")
  ) ;
COMMENT on column "t_pipeline_var"."public" IS '公开';

CREATE TABLE "t_yml_plugin" (
    "aid" serial8 NOT NULL,
    "name" varchar(64) DEFAULT NULL,
    "yml_content" text,
    "deleted" smallint DEFAULT 0,
    "deleted_time" timestamp DEFAULT NULL,
    PRIMARY KEY ("aid")
  );
INSERT INTO
  "t_yml_plugin"
VALUES
  (
    1,
    'sh',
    '      - step: shell@sh\n        displayName: sh\n        name: sh\n        commands:\n          - echo hello world',
    0,
    NULL
  );
INSERT INTO
  "t_yml_plugin"
VALUES
  (
    2,
    'bash',
    '      - step: shell@bash\n        displayName: bash\n        name: bash\n        commands:\n          - echo hello world',
    0,
    NULL
  );
INSERT INTO
  "t_yml_plugin"
VALUES
  (
    3,
    'powershell',
    '      - step: shell@powershell\n        displayName: powershell\n        name: powershell\n        commands:\n          - echo hello world',
    1,
    NULL
  );
INSERT INTO
  "t_yml_plugin"
VALUES
  (
    4,
    'ssh',
    '      - step: shell@ssh\r\n        displayName: ssh\r\n        name: ssh\r\n        input:\r\n          host: localhost:22  #端口必填\r\n          user: root\r\n          pass: 123456\r\n          workspace: /root/test #为空就是 $HOME 用户目录\r\n        commands:\r\n          - echo hello world',
    0,
    NULL
  );
CREATE TABLE "t_yml_template" (
    "aid" serial8 NOT NULL,
    "name" varchar(64) DEFAULT NULL,
    "yml_content" text,
    "deleted" smallint DEFAULT 0,
    "deleted_time" timestamp DEFAULT NULL,
    PRIMARY KEY ("aid")
  );
INSERT INTO
  "t_yml_template"(
    "aid",
    "name",
    "yml_content",
    "deleted",
    "deleted_time"
  )
VALUES
  (
    1,
    'Golang',
    'version: 1.0\nvars:\nstages:\n  - stage:\n    displayName: build\n    name: build\n    steps:\n      - step: shell@sh\n        displayName: go-build-1\n        name: build\n        env:\n        commands:\n          - go build main.go\n      - step: shell@sh\n        displayName: go-build-2\n        name: test\n        env:\n        commands:\n          - go test -v\n',
    0,
    NULL
  );
INSERT INTO
  "t_yml_template"(
    "aid",
    "name",
    "yml_content",
    "deleted",
    "deleted_time"
  )
VALUES
  (
    2,
    'Maven',
    'version: 1.0\nvars:\nstages:\n  - stage:\n    displayName: build\n    name: build\n    steps:\n      - step: shell@sh\n        displayName: java-build-1\n        name: build\n        env:\n        commands:\n          - mvn clean\n          - mvn install\n      - step: shell@sh\n        displayName: java-build-2\n        name: test\n        env:\n        commands:\n          - mvn test -v',
    0,
    NULL
  );
INSERT INTO
  "t_yml_template"(
    "aid",
    "name",
    "yml_content",
    "deleted",
    "deleted_time"
  )
VALUES
  (
    3,
    'Npm',
    'version: 1.0\nvars:\nstages:\n  - stage:\n    displayName: build\n    name: build\n    steps:\n      - step: shell@sh\n        displayName: npm-build-1\n        name: build\n        env:\n        commands:\n          - npm build\n      - step: shell@sh\n        displayName: npm-build-2\n        name: publish\n        env:\n        commands:\n          - npm publish ',
    0,
    NULL
  );
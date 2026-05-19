-- 002_add_user_password.sql
-- 用户表新增密码字段（开发阶段手机号+密码登录）

ALTER TABLE users ADD COLUMN password VARCHAR(100) NOT NULL DEFAULT '';

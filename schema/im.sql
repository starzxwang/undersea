-- 本文件为im数据库建表语句，可重复执行
CREATE DATABASE IF NOT EXISTS im;

USE im;

-- 创建用户表
CREATE TABLE IF NOT EXISTS im_user (
    id int primary key auto_increment,
    pwd char(32) not null default '' comment '密码(md5)',
    avatar varchar(1000) not null default '' comment '头像链接',
    `name` varchar(30) not null default '' comment '用户名',
    deleted tinyint not null default 0 comment '是否删除',
    created_at datetime not null default current_timestamp(),
    updated_at datetime not null default current_timestamp(),
    index idx_name on(`name`)
) engine=innodb charset=utf8mb4;

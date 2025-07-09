-- Project: gorchat 
-- Author: wens
-- Date: 2025/07/07

set NAMES 'utf8mb4';
set FOREIGN_KEY_CHECKS = 0;

-- User table
drop table if exists `user`;
create table `user` (
    `user_id` int primary key auto_increment comment '用户id',
    `user_name` varchar(128) not null comment '用户名',
    `user_password` varchar(256) not null comment '用户密码',
    `user_email` varchar(64) unique not null comment '用户邮箱',
    `create_time` datetime not null default current_timestamp comment '创建时间',
    `update_time` datetime not null default current_timestamp comment '更新时间',
    `deleted` tinyint(1) not null default 0 check(deleted in (0,1)) comment '逻辑删除'
) engine=InnoDB default charset=utf8mb4 collate=utf8mb4_general_ci;

set FOREIGN_KEY_CHECKS = 1;
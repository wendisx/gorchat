-- Porject: gorchat
-- Author: wens
-- Date: 2025/07/14

set NAMES 'utf8mb4';
set FOREIGN_KEY_CHECKS = 0;

-- 用户详细信息表
DROP TABLE IF EXISTS `im_users_detail`;
CREATE TABLE `im_users_detail` (
  `user_id` bigint PRIMARY KEY COMMENT '用户账号',
  `email` varchar(64) COMMENT '邮箱',
  `phone` varchar(20) COMMENT '电话',
  `gender` varchar(4) COMMENT '性别',
  `age` int default 0 COMMENT '年龄',
  `address` varchar(64) COMMENT '地址',
  `location` varchar(64) COMMENT '位置',
  `avatar` text COMMENT '头像',
  `created_time` timestamp default current_timestamp COMMENT '用户详细创建时间',
  `updated_time` timestamp default current_timestamp on update current_timestamp COMMENT '用户详细修改时间',
  constraint `gender_check` check ((`gender` in ('男','女'))),
  constraint `age_check` check ((`age` >=0 and `age` <= 150)),
  constraint `fk_ud_to_user` FOREIGN KEY (`user_id`) REFERENCES `im_users` (`user_id`) on delete cascade
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- 用户表
DROP TABLE IF EXISTS `im_users`;
CREATE TABLE `im_users` (
  `user_id` bigint PRIMARY KEY auto_increment COMMENT '用户账号',
  `user_name` varchar(64) not null COMMENT '用户昵称',
  `user_password` varchar(64) not null COMMENT '用户密码',
  `created_time` timestamp default current_timestamp COMMENT '用户创建时间',
  `updated_time` timestamp default current_timestamp on update current_timestamp COMMENT '用户修改时间',
  `deleted` int COMMENT '逻辑删除',
  index i_user_name(user_name)
)ENGINE=InnoDB auto_increment=100000 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- 用户职责表
DROP TABLE IF EXISTS `im_users_role`;
CREATE TABLE `im_users_role` (
  `role_id` int PRIMARY KEY auto_increment COMMENT '用户职责标识',
  `role_name` varchar(32) COMMENT '职责名称',
  `created_time` timestamp COMMENT '职责创建时间',
  `updated_time` timestamp COMMENT '职责更新时间',
  `deleted` int COMMENT '逻辑删除'
)ENGINE=InnoDB auto_increment=1 default CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- 群信息表
DROP TABLE IF EXISTS `im_groups`;
CREATE TABLE `im_groups` (
  `group_id` bigint auto_increment PRIMARY KEY COMMENT '群号',
  `group_name` varchar(32) not null COMMENT '群名称',
  `group_password` varchar(20) COMMENT '群密码',
  `created_time` timestamp COMMENT '群创建时间',
  `updated_time` timestamp COMMENT '群更新时间',
  `deleted` int COMMENT '逻辑删除',
  index i_group_name(group_name)
)ENGINE=InnoDB auto_increment=1000000 default CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- 群详细信息表
DROP TABLE IF EXISTS `im_groups_detail`;
CREATE TABLE `im_groups_detail` (
  `group_id` bigint PRIMARY KEY COMMENT '群号',
  `group_avatar` text COMMENT '群头像',
  `max_size` int not null COMMENT '群容量',
  `current_size` int COMMENT '群当前人数',
  `online_size` int COMMENT '在线人数',
  `created_time` timestamp COMMENT '群创建时间',
  `updated_time` timestamp COMMENT '群更新时间',
  constraint `fk_gd_to_group` FOREIGN KEY (`group_id`) REFERENCES `im_groups` (`group_id`) on delete cascade
)ENGINE=InnoDB default CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- 群用户关系表
DROP TABLE IF EXISTS `im_groups_users`;
CREATE TABLE `im_groups_users` (
  `group_id` bigint COMMENT '群号',
  `user_id` bigint COMMENT '用户账号',
  `group_nickname` varchar(32) COMMENT '群别称',
  `user_nickname` varchar(32) COMMENT '用户别称',
  `user_role` int COMMENT '用户职责',
  `user_role_nickname` varchar(32) COMMENT '用户职责别称',
  `disturb` int default 1 COMMENT '群打扰模式',
  `created_time` timestamp COMMENT '用户入群时间',
  `updated_time` timestamp COMMENT '用户退出群时间',
  `deleted` int COMMENT '逻辑删除',
  PRIMARY KEY (`group_id`, `user_id`),
  constraint `fk_gu_to_urole` FOREIGN KEY (`user_role`) REFERENCES `im_users_role` (`role_id`),
  constraint `fk_gu_to_groups` FOREIGN KEY (`group_id`) REFERENCES `im_groups` (`group_id`),
  constraint `fk_gu_to_users` FOREIGN KEY (`user_id`) REFERENCES `im_users` (`user_id`)
)ENGINE=InnoDB default CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- 单聊表
DROP TABLE IF EXISTS `im_single_chat`;
CREATE TABLE `im_single_chat` (
  `single_id` bigint COMMENT '单聊标识',
  `inviter_id` bigint COMMENT '邀请人账号',
  `invitee_id` bigint COMMENT '受邀人账号',
  `inviter_nickname` varchar(32) COMMENT '邀请人别称',
  `invitee_nickname` varchar(32) COMMENT '受邀人别称',
  `inviter_disturb` int default 1 COMMENT '邀请人打扰模式',
  `invitee_disturb` int default 1 COMMENT '受邀请人打扰模式',
  `created_time` timestamp COMMENT '单聊创建时间',
  `updated_time` timestamp COMMENT '单聊修改时间',
  `deleted` int COMMENT '逻辑删除',
  PRIMARY KEY (`single_id`, `inviter_id`, `invitee_id`),
  constraint `fk_sc_to_users_1` FOREIGN KEY (`inviter_id`) REFERENCES `im_users` (`user_id`) on delete cascade,
  constraint `fk_sc_to_users_2` FOREIGN KEY (`invitee_id`) REFERENCES `im_users` (`user_id`) on delete cascade
)ENGINE=InnoDB default CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- 对话分类表
DROP TABLE IF EXISTS `im_dialog`;
CREATE TABLE `im_dialog` (
  `dialog_id` int auto_increment PRIMARY KEY COMMENT '对话标识',
  `dialog_name` varchar(32) not null unique COMMENT '对话名称',
  `created_time` timestamp COMMENT '创建时间',
  `updated_time` timestamp COMMENT '更新时间',
  `deleted` int COMMENT '逻辑删除'
)ENGINE=InnoDB default CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- 历史消息
DROP TABLE IF EXISTS `im_timeline`;
CREATE TABLE `im_timeline` (
  `timeline_id` bigint COMMENT '时间线标识',
  `sequence_id` timestamp default current_timestamp COMMENT '序列标识',
  `sender` bigint COMMENT '发送者',
  `dialog_type` int COMMENT '对话类型',
  `message_id` bigint COMMENT '消息标识',
  `deleted` int COMMENT '逻辑删除',
  PRIMARY KEY (`timeline_id`, `sequence_id`),
  constraint `fk_timeline_to_dialog` FOREIGN KEY (`dialog_type`) REFERENCES `im_dialog` (`dialog_id`),
  constraint `fk_timeline_to_groups` FOREIGN KEY (`timeline_id`) REFERENCES `im_groups` (`group_id`),
  constraint `fk_timeline_to_sc` FOREIGN KEY (`timeline_id`) REFERENCES `im_single_chat` (`single_id`),
  constraint `fk_timeline_to_message` FOREIGN KEY (`message_id`) REFERENCES `im_message` (`message_id`),
  index i_sender_dtype(sender,dialog_type)
)ENGINE=InnoDB default CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- 消息类型表
DROP TABLE IF EXISTS `im_message_type`;
CREATE TABLE `im_message_type` (
  `type_id` int auto_increment PRIMARY KEY COMMENT '类型标识',
  `type_name` varchar(32) not null unique COMMENT '类型名称',
  `created_time` timestamp COMMENT '类型创建时间',
  `updated_time` timestamp COMMENT '类型更新时间',
  `deleted` int COMMENT '逻辑删除'
)ENGINE=InnoDB default CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- 消息状态表
DROP TABLE IF EXISTS `im_message_status`;
CREATE TABLE `im_message_status` (
  `status_id` int auto_increment PRIMARY KEY COMMENT '状态标识',
  `status_name` varchar(32) not null unique COMMENT '状态名称',
  `created_time` timestamp COMMENT '状态创建时间',
  `updated_time` timestamp COMMENT '状态更新时间',
  `deleted` int COMMENT '逻辑删除'
)ENGINE=InnoDB default CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- 消息表
DROP TABLE IF EXISTS `im_message`;
CREATE TABLE `im_message` (
  `message_id` bigint PRIMARY KEY auto_increment COMMENT '消息标识',
  `sender` bigint COMMENT '发送者',
  `receiver` bigint COMMENT '接收者',
  `type` int COMMENT '消息类型',
  `content` text COMMENT '消息内容',
  `status` int COMMENT '消息状态',
  `send_time` timestamp default current_timestamp COMMENT '发送时间',
  `deleted` int COMMENT '逻辑删除',
  constraint `fk_message_to_type`FOREIGN KEY (`type`) REFERENCES `im_message_type` (`type_id`),
  constraint `fk_message_to_status` FOREIGN KEY (`status`) REFERENCES `im_message_status` (`status_id`)
)ENGINE=InnoDB default CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

set FOREIGN_KEY_CHECKS = 1;
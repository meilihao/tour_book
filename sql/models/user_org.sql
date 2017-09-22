-- 先有user,再org
-- 参考: 组织管理相关的表 - http://www.zentao.net/book/zentaopmshelp/157.html

-- 用户表
CREATE TABLE `user` (
  `id` int NOT NULL AUTO_INCREMENT,
  `mobile` char(11) NOT NULL DEFAULT '',
  `password` varchar(128) NOT NULL DEFAULT '',
  `realname` varchar(30) NOT NULL DEFAULT '',
  `avatar` varchar(128) NOT NULL DEFAULT '',
  `birthday` date NOT NULL DEFAULT '0000-00-00',
  `gender` tinyint NOT NULL DEFAULT '0', -- 性别:0,未知;1,女;2,男
  `address` varchar(120) DEFAULT '', -- 现住地
  `zipcode` varchar(10) DEFAULT '', -- 现住地邮编
  `verify_status` int NOT NULL DEFAULT '0', -- 验证状态:1,实名认证
  `visit_count` int NOT NULL DEFAULT '0', -- 登录次数
  `ip` varchar(32) NOT NULL DEFAULT '', -- 最近登录IP
  `last_at` datetime NOT NULL DEFAULT '0000-00-00 00:00:00', -- 最近登录IP时间
  `status` tinyint NOT NULL DEFAULT '0', -- 状态: 0,禁用;1,正常
  `created_at` datetime NOT NULL DEFAULT '0000-00-00 00:00:00',
  `updated_at` datetime NOT NULL DEFAULT '0000-00-00 00:00:00',
  PRIMARY KEY (`id`),
  UNIQUE KEY `mobile` (`mobile`)
) ENGINE=InnoDB AUTO_INCREMENT=10000 DEFAULT CHARSET=utf8;

-- 组织表(通常是公司)
CREATE TABLE `org` (
  `id` int NOT NULL AUTO_INCREMENT,
  `code` char(18) DEFAULT NULL, -- 公司用统一社会信用代码(长度18)
  `name` varchar(120) NOT NULL DEFAULT '', -- 组织名称
  `short_name` varchar(32) NOT NUll DEFAULT '',
  `manager_userid` int NOT NULL DEFAULT '0', -- 管理员的userid
  `address` varchar(120) DEFAULT '',
  `zipcode` varchar(10)  DEFAULT '',
  `website` varchar(120) DEFAULT '', -- 官网
  `created_at` datetime NOT NULL DEFAULT '0000-00-00 00:00:00',
  `updated_at` datetime NOT NULL DEFAULT '0000-00-00 00:00:00',
  PRIMARY KEY (`id`),
  UNIQUE (`code`)
) ENGINE=InnoDB AUTO_INCREMENT=10000 DEFAULT CHARSET=utf8

-- 部门结构
CREATE TABLE `dept` (
  `id` int NOT NULL AUTO_INCREMENT,
  `org_id` int NOT NULL DEFAULT '0',
  `name` varchar(60) NOT NULL DEFAULT '',
  `parent` int NOT NULL DEFAULT '0',
  `path` varchar(255) NOT NULL DEFAULT '', -- 具体层级,用逗号分隔
  `grade` tinyint NOT NULL DEFAULT '0', -- 层级
  `order` tinyint NOT NULL DEFAULT '0', -- 排序,小的靠前
  `leader_userid` int NOT NULL DEFAULT '', -- 负责人
  `created_at` datetime NOT NULL DEFAULT '0000-00-00 00:00:00',
  `updated_at` datetime NOT NULL DEFAULT '0000-00-00 00:00:00',
  PRIMARY KEY (`id`),
  KEY `dept` (`org_id`,`parent`,`path`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- 用户和组织之间的对应关系
CREATE TABLE `user_org` (
  `org_id` int NOT NULL DEFAULT '0',
  `user_id` int NOT NULL DEFAULT '0',
  `dept_id` int NOT NULL DEFAULT '0', -- 暂不使用
  `nickname` varchar(60) NOT NULL DEFAULT '', -- 花名
  `joinby_userid` int NOT NULL DEFAULT '0', -- 引入人
  `created_at` datetime NOT NULL DEFAULT '0000-00-00 00:00:00',
  `updated_at` datetime NOT NULL DEFAULT '0000-00-00 00:00:00',
  UNIQUE KEY `userorg` (`org_id`,`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8

-- 用户和分组之间的对应关系
CREATE TABLE `user_group` (
  `org_id` int NOT NULL DEFAULT '0',
  `user_id` int NOT NULL DEFAULT '0',
  `group_id` int NOT NULL DEFAULT '0',
  UNIQUE KEY `usergroup` (`org_id`,`user_id`,`group_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8

-- 分组表
CREATE TABLE `group` (
  `id` int NOT NULL DEFAULT '0', -- 有初始分组admin(0),不使用自增
  `org_id`int NOT NULL DEFAULT '0', 
  `name` varchar(30) NOT NULL,
  `desc` varchar(255) NOT NULL DEFAULT '',
  `created_at` datetime NOT NULL DEFAULT '0000-00-00 00:00:00',
  `updated_at` datetime NOT NULL DEFAULT '0000-00-00 00:00:00',
  PRIMARY KEY (`org_id`,`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8

-- 分组的权限
CREATE TABLE `privilege` (
  `org_id`int NOT NULL DEFAULT '0', 
  `group` int NOT NULL DEFAULT '0',
  `module` varchar(30) NOT NULL DEFAULT '',
  `method` varchar(30) NOT NULL DEFAULT '',
  UNIQUE KEY `group` (`org_id`,`group`,`module`,`method`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8
CREATE DATABASE IF NOT EXISTS `douyin` /*!40100 DEFAULT CHARACTER SET utf8 COLLATE utf8_unicode_ci */;
USE `douyin`;
DROP TABLE IF EXISTS `userinfo`;
CREATE TABLE `userinfo`
(
    `id`            bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
    `name`          varchar(128)        NOT NULL DEFAULT '' COMMENT '用户昵称',
    `token`         varchar(128)        NOT NULL COMMENT '用户鉴权',
    `create_time`   timestamp           NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '用户创建时间',
    PRIMARY KEY(`id`)
) ENGINE = InnoDB
    DEFAULT CHARSET = utf8mb4 COMMENT = '用户信息表';

DROP TABLE IF EXISTS `userrelation`;
CREATE TABLE `userrelation`
(
    `id`            bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
    `from_user_id`  bigint(20) unsigned NOT NULL COMMENT '关注用户id',
    `to_user_id`    bigint(20) unsigned NOT NULL COMMENT '被关注用户id',
    `is_follow`     int                 NOT NULL DEFAULT 1 COMMENT '是否关注',
    PRIMARY KEY (`id`),
    FOREIGN KEY(`from_user_id`) REFERENCES `userinfo`(`id`),
    FOREIGN KEY(`to_user_id`) REFERENCES `userinfo`(`id`)
) ENGINE = InnoDB
    DEFAULT CHARSET = utf8mb4 COMMENT = '用户关系表';

DROP TABLE IF EXISTS `message`;
CREATE TABLE `message`
(
    `id`            bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
    `from_user_id`  bigint(20) unsigned NOT NULL COMMENT '发送消息用户id',
    `to_user_id`    bigint(20) unsigned NOT NULL COMMENT '接收消息用户id',
    `content`       text                NOT NULL COMMENT '消息内容',
    `create_time`   timestamp           NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '发送消息时间',
    PRIMARY KEY (`id`),
    FOREIGN KEY(`from_user_id`) REFERENCES `userinfo`(`id`),
    FOREIGN KEY(`to_user_id`) REFERENCES `userinfo`(`id`)
) ENGINE = InnoDB
    DEFAULT CHARSET = utf8mb4 COMMENT = '消息表';

DROP TABLE IF EXISTS `video`;
CREATE TABLE `video`
(
    `id`            bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
    `user_id`       bigint(20) unsigned NOT NULL COMMENT '发视频用户id',
    `play_url`      varchar(128)        NOT NULL COMMENT '视频播放地址',
    `cover_url`     varchar(128)        NOT NULL COMMENT '视频封面地址',
    `title`         varchar(128)        NOT NULL COMMENT '视频标题',
    `create_time`   timestamp           NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '视频发布时间',
    PRIMARY KEY (`id`),
    FOREIGN KEY(`user_id`) REFERENCES `userinfo`(`id`)
) ENGINE = InnoDB
    DEFAULT CHARSET = utf8mb4 COMMENT = '视频表';

DROP TABLE IF EXISTS `videorelation`;
CREATE TABLE `videorelation`
(
    `id`            bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
    `from_user_id`  bigint(20) unsigned NOT NULL COMMENT '点赞用户id',
    `to_video_id`   bigint(20) unsigned NOT NULL COMMENT '点赞视频id',
    `is_like`       int                 NOT NULL DEFAULT 1  COMMENT '是否点赞',  
    PRIMARY KEY (`id`),
    FOREIGN KEY(`from_user_id`) REFERENCES `userinfo`(`id`),
    FOREIGN KEY(`to_video_id`) REFERENCES `video`(`id`)
) ENGINE = InnoDB
    DEFAULT CHARSET = utf8mb4 COMMENT = '用户视频关系表';

DROP TABLE IF EXISTS `comment`;
CREATE TABLE `comment`
(
    `id`            bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `from_user_id`  bigint(20) unsigned NOT NULL COMMENT '评论用户id',    
    `to_video_id`   bigint(20) unsigned NOT NULL COMMENT '评论视频id',
    `content`       text                NOT NULL COMMENT '评论内容',
    `status`        int                 NOT NULL DEFAULT 1 COMMENT '存在、修改、删除',
    `create_time`   timestamp           NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '评论创建时间',
    `update_time`   timestamp           COMMENT '评论修改或删除时间',
    PRIMARY KEY (`id`),
    FOREIGN KEY (`from_user_id`) REFERENCES `userinfo`(`id`),
    FOREIGN KEY (`to_video_id`) REFERENCES `video`(`id`)
) ENGINE = InnoDB
    DEFAULT CHARSET = utf8mb4 COMMENT = '评论表';


INSERT INTO `userinfo`
VALUES(1, 'lulingyi', 'lulingyiroot123','2023-1-21 10:00:00');
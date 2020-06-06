CREATE DATABASE IF NOT EXISTS `test_tcc` CHARACTER SET utf8mb4 COLLATE utf8mb4_bin;
CREATE TABLE IF NOT EXISTS `test_tcc`.`tcc_async_task` (
    `id` bigint NOT NULL AUTO_INCREMENT,
    `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `uid` varchar(36) NOT NULL,
    `name` varchar(64) NOT NULL,
    `status` varchar(64) NOT NULL,
    `value` TEXT NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_uid` (`uid`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_bin;
CREATE TABLE IF NOT EXISTS `test_tcc`.`tcc_lock` (
    `id` bigint NOT NULL AUTO_INCREMENT,
    `key` varchar(127) NOT NULL,
    `expire_at` datetime(3) NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_key` (`key`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_bin;
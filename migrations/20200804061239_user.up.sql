CREATE TABLE `users`(
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `username` VARCHAR(32) NOT NULL,
    `email` VARCHAR(64) NOT NULL,
    `email_verified` TINYINT NOT NULL DEFAULT 0,
    `verification_token` CHAR(43) NULL,
    `hashed_password` VARCHAR(128) NOT NULL,
    `password_reset_token` CHAR(43) NULL,
    `created_at` DATETIME NOT NULL,
    `updated_at` DATETIME NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_users_username` (`username`),
    UNIQUE KEY `idx_users_email` (`email`),
    KEY `idx_users_email_verified` (`email_verified`),
    UNIQUE KEY `idx_users_verification_token` (`verification_token`),
    UNIQUE KEY `idx_users_password_reset_token` (`password_reset_token`)
);

ALTER TABLE `domains` DROP INDEX IF EXISTS `idx_domains_user_id_name`;
ALTER TABLE `domains` DROP COLUMN `user_id`;
ALTER TABLE `domains` ADD COLUMN `user_id` BIGINT NOT NULL DEFAULT 0 AFTER `id`;
CREATE INDEX `idx_domains_user_id` ON `domains` (`user_id`);
UPDATE `domains` SET `user_id`=1;

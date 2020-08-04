CREATE TABLE `users`(
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `username` VARCHAR(32) NOT NULL,
    `email` VARCHAR(64) NOT NULL,
    `email_verified` TINYINT NOT NULL DEFAULT 0,
    `verification_token` VARCHAR(32) NULL,
    `hashed_password` VARCHAR(128) NOT NULL,
    `created_at` DATETIME NOT NULL,
    `updated_at` DATETIME NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_users_username` (`username`),
    UNIQUE KEY `idx_users_email` (`email`),
    KEY `idx_users_email_verified` (`email_verified`),
    UNIQUE KEY `idx_users_verification_token` (`verification_token`)
);

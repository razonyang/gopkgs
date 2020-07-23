CREATE TABLE `domains` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `user_id` VARCHAR(64) NOT NULL,
    `name` VARCHAR(64) NOT NULL,
    `verified` TINYINT NOT NULL DEFAULT 0,
    `challenge_txt` VARCHAR(32) NOT NULL,
    `created_at` DATETIME NOT NULL,
    `updated_at` DATETIME NOT NULL,
    PRIMARY KEY(`id`),
    UNIQUE KEY `idx_domains_user_id_name` (`user_id`, `name`),
    KEY `idx_domains_verified` (`verified`)
);

CREATE TABLE `packages` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `domain_id` BIGINT NOT NULL,
    `path` VARCHAR(64) NOT NULL,
    `vcs` VARCHAR(6) NOT NULL,
    `root` VARCHAR(256) NOT NULL,
    `docs` VARCHAR(256) NOT NULL,
    `created_at` DATETIME NOT NULL,
    `updated_at` DATETIME NOT NULL,
    PRIMARY KEY(`id`),
    UNIQUE KEY `idx_packages_domain_id_path` (`domain_id`, `path`),
    KEY`idx_packages_vcs` (`vcs`)
);

CREATE TABLE `actions` (
    `id` VARCHAR(32) NOT NULL,
    `package_id` BIGINT NOT NULL,
    `kind` VARCHAR(16) NOT NULL,
    `created_at` DATETIME NOT NULL,
    PRIMARY KEY (`id`),
    KEY `idx_actions_package_id` (`package_id`),
    KEY `idx_actions_kind` (`kind`),
    KEY `idx_actions_created_at` (`created_at`)
);

CREATE TABLE `calendars` (
    id DATE NOT NULL,
    PRIMARY KEY(`id`) 
);
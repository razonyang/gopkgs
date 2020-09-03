ALTER TABLE `users` ADD COLUMN `timezone` VARCHAR(64) NOT NULL DEFAULT 'Asia/Shanghai';
ALTER TABLE `users` ADD COLUMN `auth_key` CHAR(32);
CREATE UNIQUE INDEX idx_users_auth_key ON `users` (`auth_key`);

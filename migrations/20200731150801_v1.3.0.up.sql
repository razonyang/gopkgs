ALTER TABLE `packages` ADD COLUMN `description` VARCHAR(256) NOT NULL DEFAULT '' AFTER `docs`;
ALTER TABLE `packages` ADD COLUMN `private` TINYINT NOT NULL DEFAULT 0 AFTER `domain_id`;
ALTER TABLE `packages` ADD INDEX `idx_packages_private` (`private`);
ALTER TABLE `packages` MODIFY COLUMN `docs` VARCHAR(256) NOT NULL DEFAULT '';
ALTER TABLE `packages` ADD COLUMN `homepage` VARCHAR(256) NOT NULL DEFAULT '' AFTER `description`;
ALTER TABLE `packages` ADD COLUMN `license` VARCHAR(256) NOT NULL DEFAULT '' AFTER `homepage`;

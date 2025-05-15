-- Create "availabilities" table
CREATE TABLE `availabilities` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL,
  `updated_at` datetime(3) NULL,
  `deleted_at` datetime(3) NULL,
  `uuid` binary(16) NOT NULL,
  `resource_id` binary(16) NOT NULL,
  `start_time` datetime(3) NOT NULL,
  `end_time` datetime(3) NOT NULL,
  `is_recurring` bool NULL DEFAULT 0,
  `recur_rule` varchar(255) NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_availabilities_deleted_at` (`deleted_at`),
  INDEX `idx_availabilities_end_time` (`end_time`),
  INDEX `idx_availabilities_resource_id` (`resource_id`),
  INDEX `idx_availabilities_start_time` (`start_time`),
  INDEX `idx_availabilities_uuid` (`uuid`),
  UNIQUE INDEX `uni_availabilities_uuid` (`uuid`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "bookings" table
CREATE TABLE `bookings` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL,
  `updated_at` datetime(3) NULL,
  `deleted_at` datetime(3) NULL,
  `uuid` binary(16) NOT NULL,
  `resource_id` binary(16) NOT NULL,
  `user_id` binary(16) NOT NULL,
  `start_time` datetime(3) NOT NULL,
  `end_time` datetime(3) NOT NULL,
  `status` varchar(50) NULL DEFAULT "pending",
  `notes` text NULL,
  `reference` varchar(100) NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_bookings_deleted_at` (`deleted_at`),
  INDEX `idx_bookings_end_time` (`end_time`),
  INDEX `idx_bookings_resource_id` (`resource_id`),
  INDEX `idx_bookings_start_time` (`start_time`),
  INDEX `idx_bookings_user_id` (`user_id`),
  INDEX `idx_bookings_uuid` (`uuid`),
  UNIQUE INDEX `uni_bookings_uuid` (`uuid`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "organizations" table
CREATE TABLE `organizations` (
  `id` binary(16) NOT NULL,
  `name` longtext NOT NULL,
  `location` longtext NULL,
  `established_at` datetime(3) NULL,
  `created_at` datetime(3) NULL,
  `updated_at` datetime(3) NULL,
  PRIMARY KEY (`id`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "resources" table
CREATE TABLE `resources` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL,
  `updated_at` datetime(3) NULL,
  `deleted_at` datetime(3) NULL,
  `uuid` binary(16) NOT NULL,
  `name` varchar(255) NOT NULL,
  `description` text NULL,
  `type` varchar(50) NOT NULL,
  `capacity` bigint NULL DEFAULT 1,
  `location` varchar(255) NULL,
  `attributes` json NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_resources_deleted_at` (`deleted_at`),
  INDEX `idx_resources_uuid` (`uuid`),
  UNIQUE INDEX `uni_resources_uuid` (`uuid`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "todos" table
CREATE TABLE `todos` (
  `id` binary(16) NOT NULL,
  `title` longtext NOT NULL,
  `description` longtext NULL,
  `created_at` datetime(3) NULL,
  `updated_at` datetime(3) NULL,
  PRIMARY KEY (`id`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;

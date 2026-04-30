SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

CREATE DATABASE IF NOT EXISTS `go_base_project`
  DEFAULT CHARACTER SET utf8mb4
  COLLATE utf8mb4_general_ci;

USE `go_base_project`;

DROP TABLE IF EXISTS `gbp_operation_audit_logs`;
DROP TABLE IF EXISTS `gbp_login_audit_logs`;
DROP TABLE IF EXISTS `gbp_system_params`;
DROP TABLE IF EXISTS `gbp_dictionary_items`;
DROP TABLE IF EXISTS `gbp_dictionaries`;
DROP TABLE IF EXISTS `gbp_role_data_scopes`;
DROP TABLE IF EXISTS `gbp_permission_policies`;
DROP TABLE IF EXISTS `gbp_api_skip_rules`;
DROP TABLE IF EXISTS `gbp_api_resources`;
DROP TABLE IF EXISTS `gbp_role_actions`;
DROP TABLE IF EXISTS `gbp_menu_actions`;
DROP TABLE IF EXISTS `gbp_role_menus`;
DROP TABLE IF EXISTS `gbp_menu_route_params`;
DROP TABLE IF EXISTS `gbp_menus`;
DROP TABLE IF EXISTS `gbp_jwt_blocklist`;
DROP TABLE IF EXISTS `gbp_user_roles`;
DROP TABLE IF EXISTS `gbp_users`;
DROP TABLE IF EXISTS `gbp_roles`;

CREATE TABLE `gbp_roles` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `deleted_at` DATETIME(3) NULL,
  `role_code` VARCHAR(64) NOT NULL,
  `role_name` VARCHAR(128) NOT NULL,
  `parent_role_id` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `default_route` VARCHAR(128) NOT NULL DEFAULT 'dashboard',
  `sort_no` INT NOT NULL DEFAULT 0,
  `role_status` TINYINT NOT NULL DEFAULT 1,
  `remark` VARCHAR(255) NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_gbp_roles_role_code` (`role_code`),
  KEY `idx_gbp_roles_parent_role_id` (`parent_role_id`),
  KEY `idx_gbp_roles_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='roles';

CREATE TABLE `gbp_users` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `deleted_at` DATETIME(3) NULL,
  `user_uuid` VARCHAR(64) NOT NULL,
  `login_name` VARCHAR(128) NOT NULL,
  `password_hash` VARCHAR(255) NOT NULL,
  `display_name` VARCHAR(128) NOT NULL DEFAULT '系统用户',
  `avatar_url` VARCHAR(512) NULL,
  `primary_role_id` BIGINT UNSIGNED NULL,
  `phone_number` VARCHAR(32) NULL,
  `email_address` VARCHAR(128) NULL,
  `user_status` TINYINT NOT NULL DEFAULT 1,
  `must_change_password` TINYINT(1) NOT NULL DEFAULT 0,
  `last_login_at` DATETIME(3) NULL,
  `profile_config` JSON NULL,
  `remark` VARCHAR(255) NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_gbp_users_user_uuid` (`user_uuid`),
  UNIQUE KEY `uk_gbp_users_login_name` (`login_name`),
  KEY `idx_gbp_users_primary_role_id` (`primary_role_id`),
  KEY `idx_gbp_users_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='users';

CREATE TABLE `gbp_user_roles` (
  `user_id` BIGINT UNSIGNED NOT NULL,
  `role_id` BIGINT UNSIGNED NOT NULL,
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`user_id`, `role_id`),
  KEY `idx_gbp_user_roles_role_id` (`role_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='user roles';

CREATE TABLE `gbp_jwt_blocklist` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `deleted_at` DATETIME(3) NULL,
  `token_hash` CHAR(64) NOT NULL,
  `expires_at` DATETIME(3) NOT NULL,
  `logout_reason` VARCHAR(255) NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_gbp_jwt_blocklist_token_hash` (`token_hash`),
  KEY `idx_gbp_jwt_blocklist_expires_at` (`expires_at`),
  KEY `idx_gbp_jwt_blocklist_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='jwt blocklist';

CREATE TABLE `gbp_menus` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `deleted_at` DATETIME(3) NULL,
  `parent_id` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `menu_type` TINYINT NOT NULL DEFAULT 2,
  `route_path` VARCHAR(255) NULL,
  `route_name` VARCHAR(128) NULL,
  `component_path` VARCHAR(255) NULL,
  `redirect_path` VARCHAR(255) NULL,
  `menu_title` VARCHAR(128) NOT NULL,
  `menu_icon` VARCHAR(128) NULL,
  `sort_no` INT NOT NULL DEFAULT 0,
  `is_hidden` TINYINT(1) NOT NULL DEFAULT 0,
  `is_keep_alive` TINYINT(1) NOT NULL DEFAULT 0,
  `is_affix` TINYINT(1) NOT NULL DEFAULT 0,
  `active_route` VARCHAR(128) NULL,
  `transition_name` VARCHAR(64) NULL,
  `external_url` VARCHAR(512) NULL,
  `menu_status` TINYINT NOT NULL DEFAULT 1,
  PRIMARY KEY (`id`),
  KEY `idx_gbp_menus_parent_id` (`parent_id`),
  KEY `idx_gbp_menus_route_name` (`route_name`),
  KEY `idx_gbp_menus_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='menus';

CREATE TABLE `gbp_menu_route_params` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `deleted_at` DATETIME(3) NULL,
  `menu_id` BIGINT UNSIGNED NOT NULL,
  `param_mode` VARCHAR(32) NOT NULL DEFAULT 'query',
  `param_key` VARCHAR(128) NOT NULL,
  `param_value` VARCHAR(255) NULL,
  PRIMARY KEY (`id`),
  KEY `idx_gbp_menu_route_params_menu_id` (`menu_id`),
  KEY `idx_gbp_menu_route_params_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='menu route params';

CREATE TABLE `gbp_role_menus` (
  `role_id` BIGINT UNSIGNED NOT NULL,
  `menu_id` BIGINT UNSIGNED NOT NULL,
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`role_id`, `menu_id`),
  KEY `idx_gbp_role_menus_menu_id` (`menu_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='role menus';

CREATE TABLE `gbp_menu_actions` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `deleted_at` DATETIME(3) NULL,
  `menu_id` BIGINT UNSIGNED NOT NULL,
  `action_code` VARCHAR(64) NOT NULL,
  `action_name` VARCHAR(128) NOT NULL,
  `action_desc` VARCHAR(255) NULL,
  `sort_no` INT NOT NULL DEFAULT 0,
  `action_status` TINYINT NOT NULL DEFAULT 1,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_gbp_menu_actions_menu_code` (`menu_id`, `action_code`),
  KEY `idx_gbp_menu_actions_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='menu actions';

CREATE TABLE `gbp_role_actions` (
  `role_id` BIGINT UNSIGNED NOT NULL,
  `menu_id` BIGINT UNSIGNED NOT NULL,
  `action_id` BIGINT UNSIGNED NOT NULL,
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`role_id`, `menu_id`, `action_id`),
  KEY `idx_gbp_role_actions_action_id` (`action_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='role actions';

CREATE TABLE `gbp_api_resources` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `deleted_at` DATETIME(3) NULL,
  `api_path` VARCHAR(255) NOT NULL,
  `api_method` VARCHAR(16) NOT NULL DEFAULT 'POST',
  `api_group` VARCHAR(128) NULL,
  `api_desc` VARCHAR(255) NULL,
  `api_status` TINYINT NOT NULL DEFAULT 1,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_gbp_api_resources_path_method` (`api_path`, `api_method`),
  KEY `idx_gbp_api_resources_group` (`api_group`),
  KEY `idx_gbp_api_resources_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='api resources';

CREATE TABLE `gbp_api_skip_rules` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `deleted_at` DATETIME(3) NULL,
  `api_path` VARCHAR(255) NOT NULL,
  `api_method` VARCHAR(16) NOT NULL DEFAULT 'POST',
  `skip_reason` VARCHAR(255) NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_gbp_api_skip_rules_path_method` (`api_path`, `api_method`),
  KEY `idx_gbp_api_skip_rules_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='api skip rules';

CREATE TABLE `gbp_permission_policies` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `deleted_at` DATETIME(3) NULL,
  `subject_type` VARCHAR(32) NOT NULL DEFAULT 'role',
  `subject_id` BIGINT UNSIGNED NOT NULL,
  `resource_type` VARCHAR(32) NOT NULL DEFAULT 'api',
  `resource_id` BIGINT UNSIGNED NULL,
  `resource_key` VARCHAR(255) NOT NULL,
  `action` VARCHAR(64) NOT NULL DEFAULT '*',
  `effect` VARCHAR(16) NOT NULL DEFAULT 'allow',
  `condition_expr` VARCHAR(512) NULL,
  `policy_status` TINYINT NOT NULL DEFAULT 1,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_gbp_permission_policies` (`subject_type`, `subject_id`, `resource_type`, `resource_key`, `action`),
  KEY `idx_gbp_permission_policies_subject` (`subject_type`, `subject_id`),
  KEY `idx_gbp_permission_policies_resource` (`resource_type`, `resource_id`),
  KEY `idx_gbp_permission_policies_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='permission policies';

CREATE TABLE `gbp_role_data_scopes` (
  `role_id` BIGINT UNSIGNED NOT NULL,
  `visible_role_id` BIGINT UNSIGNED NOT NULL,
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`role_id`, `visible_role_id`),
  KEY `idx_gbp_role_data_scopes_visible_role_id` (`visible_role_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='role data scopes';

CREATE TABLE `gbp_dictionaries` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `deleted_at` DATETIME(3) NULL,
  `dict_name` VARCHAR(128) NOT NULL,
  `dict_code` VARCHAR(128) NOT NULL,
  `dict_status` TINYINT NOT NULL DEFAULT 1,
  `parent_id` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `remark` VARCHAR(255) NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_gbp_dictionaries_dict_code` (`dict_code`),
  KEY `idx_gbp_dictionaries_parent_id` (`parent_id`),
  KEY `idx_gbp_dictionaries_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='dictionaries';

CREATE TABLE `gbp_dictionary_items` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `deleted_at` DATETIME(3) NULL,
  `dict_id` BIGINT UNSIGNED NOT NULL,
  `item_label` VARCHAR(128) NOT NULL,
  `item_value` VARCHAR(128) NOT NULL,
  `item_extra` VARCHAR(255) NULL,
  `item_status` TINYINT NOT NULL DEFAULT 1,
  `sort_no` INT NOT NULL DEFAULT 0,
  `parent_id` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `tree_level` INT NOT NULL DEFAULT 1,
  `tree_path` VARCHAR(512) NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_gbp_dictionary_items_dict_value` (`dict_id`, `item_value`),
  KEY `idx_gbp_dictionary_items_parent_id` (`parent_id`),
  KEY `idx_gbp_dictionary_items_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='dictionary items';

CREATE TABLE `gbp_system_params` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `deleted_at` DATETIME(3) NULL,
  `param_name` VARCHAR(128) NOT NULL,
  `param_key` VARCHAR(128) NOT NULL,
  `param_value` TEXT NULL,
  `param_type` VARCHAR(32) NOT NULL DEFAULT 'string',
  `is_encrypted` TINYINT(1) NOT NULL DEFAULT 0,
  `param_status` TINYINT NOT NULL DEFAULT 1,
  `remark` VARCHAR(255) NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_gbp_system_params_param_key` (`param_key`),
  KEY `idx_gbp_system_params_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='system params';

CREATE TABLE `gbp_login_audit_logs` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `deleted_at` DATETIME(3) NULL,
  `user_id` BIGINT UNSIGNED NULL,
  `login_name` VARCHAR(128) NULL,
  `source_ip` VARCHAR(64) NULL,
  `login_success` TINYINT(1) NOT NULL DEFAULT 0,
  `fail_reason` VARCHAR(255) NULL,
  `user_agent` VARCHAR(512) NULL,
  PRIMARY KEY (`id`),
  KEY `idx_gbp_login_audit_logs_user_id` (`user_id`),
  KEY `idx_gbp_login_audit_logs_login_name` (`login_name`),
  KEY `idx_gbp_login_audit_logs_source_ip` (`source_ip`),
  KEY `idx_gbp_login_audit_logs_created_at` (`created_at`),
  KEY `idx_gbp_login_audit_logs_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='login audit logs';

CREATE TABLE `gbp_operation_audit_logs` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `deleted_at` DATETIME(3) NULL,
  `user_id` BIGINT UNSIGNED NULL,
  `source_ip` VARCHAR(64) NULL,
  `request_method` VARCHAR(16) NULL,
  `request_path` VARCHAR(255) NULL,
  `status_code` INT NULL,
  `cost_ms` BIGINT NULL,
  `user_agent` TEXT NULL,
  `error_message` VARCHAR(512) NULL,
  `request_body` LONGTEXT NULL,
  `response_body` LONGTEXT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_gbp_operation_audit_logs_user_id` (`user_id`),
  KEY `idx_gbp_operation_audit_logs_request_path` (`request_path`),
  KEY `idx_gbp_operation_audit_logs_status_code` (`status_code`),
  KEY `idx_gbp_operation_audit_logs_created_at` (`created_at`),
  KEY `idx_gbp_operation_audit_logs_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='operation audit logs';

SET FOREIGN_KEY_CHECKS = 1;

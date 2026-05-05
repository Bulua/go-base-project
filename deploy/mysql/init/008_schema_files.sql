-- 008_schema_files.sql
-- 文件管理模块：表结构 + 菜单 + API 资源 + super_admin 权限

USE `go_base_project`;

-- ── 文件记录表 ─────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS `gbp_files` (
  `id`            BIGINT UNSIGNED  NOT NULL AUTO_INCREMENT,
  `original_name` VARCHAR(255)     NOT NULL COMMENT '原始文件名',
  `storage_key`   VARCHAR(512)     NOT NULL COMMENT '磁盘存储相对路径，格式 YYYY/MM/DD/<uuid>.<ext>',
  `file_size`     BIGINT           NOT NULL COMMENT '文件大小（字节）',
  `mime_type`     VARCHAR(128)     NOT NULL DEFAULT '' COMMENT 'MIME 类型',
  `uploader_id`   BIGINT UNSIGNED  NULL COMMENT '上传用户 ID',
  `created_at`    DATETIME(3)      NOT NULL DEFAULT NOW(3),
  `deleted_at`    DATETIME(3)      NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_storage_key` (`storage_key`),
  KEY `idx_uploader`  (`uploader_id`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='文件记录';

-- ── 菜单 ──────────────────────────────────────────────────────────────────
INSERT INTO `gbp_menus`
  (`id`, `parent_id`, `menu_type`, `route_path`, `route_name`, `component_path`,
   `redirect_path`, `menu_title`, `menu_icon`, `sort_no`, `is_hidden`, `is_keep_alive`, `menu_status`)
VALUES
  (30, 0, 1, '/file', 'file', 'layouts/default',
   '/file/list', '文件管理', 'FolderOpened', 30, 0, 0, 1),
  (31, 30, 2, 'list', 'fileList', 'views/system/file/index.vue',
   NULL, '文件列表', 'Files', 1, 0, 0, 1)
ON DUPLICATE KEY UPDATE `menu_title` = VALUES(`menu_title`);

INSERT IGNORE INTO `gbp_role_menus` (`role_id`, `menu_id`) VALUES (1, 30), (1, 31);

-- ── API 资源 ───────────────────────────────────────────────────────────────
INSERT INTO `gbp_api_resources` (`api_path`, `api_method`, `api_group`, `api_desc`, `api_status`) VALUES
  ('/api/v1/files',          'POST',   'file', '上传文件',   1),
  ('/api/v1/files',          'GET',    'file', '文件列表',   1),
  ('/api/v1/files/{id}',     'DELETE', 'file', '删除文件',   1),
  ('/api/v1/files/{id}/raw', 'GET',    'file', '下载/预览文件', 1)
ON DUPLICATE KEY UPDATE `api_group` = VALUES(`api_group`), `api_desc` = VALUES(`api_desc`),
  `api_status` = 1, `deleted_at` = NULL;

-- super_admin 权限策略
INSERT INTO `gbp_permission_policies`
  (`subject_type`, `subject_id`, `resource_type`, `resource_id`, `resource_key`, `action`, `effect`, `policy_status`)
SELECT 'role', 1, 'api', id, CONCAT(api_method, ':', api_path), api_method, 'allow', 1
FROM `gbp_api_resources`
WHERE `api_group` = 'file' AND `deleted_at` IS NULL
ON DUPLICATE KEY UPDATE `deleted_at` = NULL, `policy_status` = 1;

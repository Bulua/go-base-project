-- 004_seed_dict_apis.sql
-- Adds dictionary module API resources and grants super_admin allow policies.

INSERT INTO `gbp_api_resources` (`api_path`, `api_method`, `api_group`, `api_desc`, `api_status`) VALUES
  ('/api/v1/dictionaries',                     'GET',    'dict', '字典列表',         1),
  ('/api/v1/dictionaries',                     'POST',   'dict', '创建字典',         1),
  ('/api/v1/dictionaries/{id}',                'PUT',    'dict', '修改字典',         1),
  ('/api/v1/dictionaries/{id}',                'DELETE', 'dict', '删除字典',         1),
  ('/api/v1/dictionaries/{id}/items',          'POST',   'dict', '新增字典项',       1),
  ('/api/v1/dictionaries/{dictCode}/items',    'GET',    'dict', '查询字典项',       1),
  ('/api/v1/dictionary-items/{id}',            'PUT',    'dict', '修改字典项',       1),
  ('/api/v1/dictionary-items/{id}',            'DELETE', 'dict', '删除字典项',       1)
ON DUPLICATE KEY UPDATE api_group = VALUES(api_group), api_desc = VALUES(api_desc), api_status = 1, deleted_at = NULL;

INSERT INTO `gbp_permission_policies`
  (`subject_type`, `subject_id`, `resource_type`, `resource_id`, `resource_key`, `action`, `effect`, `policy_status`)
SELECT 'role', 1, 'api', id, CONCAT(api_method, ':', api_path), api_method, 'allow', 1
FROM `gbp_api_resources`
WHERE `api_group` = 'dict' AND `deleted_at` IS NULL
ON DUPLICATE KEY UPDATE `deleted_at` = NULL, `policy_status` = 1;

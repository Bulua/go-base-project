-- 005_seed_dict_items_skip.sql
-- 字典项查询接口供所有登录用户使用（用于页面下拉数据），加入 skip rules
-- 绕过权限检查，handler 层仍然验证 Token

INSERT INTO `gbp_api_skip_rules` (`api_path`, `api_method`, `skip_reason`) VALUES
  ('/api/v1/dictionaries/{dictCode}/items', 'GET', '字典项查询，所有登录用户均可访问')
ON DUPLICATE KEY UPDATE `skip_reason` = VALUES(`skip_reason`), `deleted_at` = NULL;

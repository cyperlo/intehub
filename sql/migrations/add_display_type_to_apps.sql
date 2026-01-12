-- 添加 display_type 字段到 apps 表
ALTER TABLE apps ADD COLUMN IF NOT EXISTS display_type VARCHAR(20) DEFAULT 'dialog';

COMMENT ON COLUMN apps.display_type IS '展示类型: none-无展示, page-独立页面, dialog-弹窗';

-- 为定时任务表添加工作流支持
ALTER TABLE schedule_tasks ADD COLUMN IF NOT EXISTS workflow_id INTEGER;
ALTER TABLE schedule_tasks ADD CONSTRAINT fk_schedule_workflow FOREIGN KEY (workflow_id) REFERENCES workflows(id) ON DELETE CASCADE;

-- 更新注释
COMMENT ON COLUMN schedule_tasks.workflow_id IS '关联的工作流ID';
COMMENT ON COLUMN schedule_tasks.task_type IS '任务类型: push, app, workflow';

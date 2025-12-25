-- InteHub 数据库表结构
-- PostgreSQL

-- 1. 用户表
CREATE TABLE IF NOT EXISTS "user" (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    nickname VARCHAR(255),
    password VARCHAR(255) NOT NULL,
    role VARCHAR(50) DEFAULT 'user',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_user_deleted_at ON "user"(deleted_at);
COMMENT ON TABLE "user" IS '用户表';
COMMENT ON COLUMN "user".role IS '用户角色: admin, user';

-- 2. 应用表
CREATE TABLE IF NOT EXISTS apps (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    code TEXT,
    language VARCHAR(50) DEFAULT 'go',
    enabled BOOLEAN DEFAULT true,
    user_id INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_apps_deleted_at ON apps(deleted_at);
CREATE INDEX IF NOT EXISTS idx_apps_user_id ON apps(user_id);
COMMENT ON TABLE apps IS '应用表';
COMMENT ON COLUMN apps.language IS '编程语言: go, python, javascript';

-- 3. 应用日志表
CREATE TABLE IF NOT EXISTS app_logs (
    id SERIAL PRIMARY KEY,
    app_id INTEGER NOT NULL,
    app_name VARCHAR(255),
    status VARCHAR(50),
    output TEXT,
    error TEXT,
    duration BIGINT,
    started_at TIMESTAMP WITH TIME ZONE,
    finished_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_app_logs_app_id ON app_logs(app_id);
CREATE INDEX IF NOT EXISTS idx_app_logs_created_at ON app_logs(created_at);
COMMENT ON TABLE app_logs IS '应用执行日志表';
COMMENT ON COLUMN app_logs.status IS '执行状态: success, error';
COMMENT ON COLUMN app_logs.duration IS '执行时长(毫秒)';

-- 4. 应用商店模板表
CREATE TABLE IF NOT EXISTS app_templates (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    display_name VARCHAR(255) NOT NULL,
    description TEXT,
    icon VARCHAR(255),
    code TEXT NOT NULL,
    language VARCHAR(50) DEFAULT 'go',
    category VARCHAR(100),
    version VARCHAR(50),
    author VARCHAR(255),
    tags VARCHAR(500),
    config_schema TEXT,
    enabled BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_app_templates_deleted_at ON app_templates(deleted_at);
CREATE INDEX IF NOT EXISTS idx_app_templates_category ON app_templates(category);
COMMENT ON TABLE app_templates IS '应用商店模板表';
COMMENT ON COLUMN app_templates.config_schema IS 'JSON Schema 配置模板';

-- 5. 应用配置表
CREATE TABLE IF NOT EXISTS app_configs (
    id SERIAL PRIMARY KEY,
    app_id INTEGER NOT NULL,
    key VARCHAR(255) NOT NULL,
    value TEXT,
    type VARCHAR(50) DEFAULT 'string',
    encrypted BOOLEAN DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_app_configs_app_id ON app_configs(app_id);
CREATE INDEX IF NOT EXISTS idx_app_configs_deleted_at ON app_configs(deleted_at);
COMMENT ON TABLE app_configs IS '应用配置表';
COMMENT ON COLUMN app_configs.type IS '配置类型: string, number, boolean, json';

-- 6. 字段定义表
CREATE TABLE IF NOT EXISTS field_schemas (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    key VARCHAR(255) NOT NULL UNIQUE,
    type VARCHAR(50) NOT NULL,
    description TEXT,
    required BOOLEAN DEFAULT false,
    default_value VARCHAR(500),
    options TEXT,
    placeholder VARCHAR(255),
    validation VARCHAR(500),
    user_id INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_field_schemas_deleted_at ON field_schemas(deleted_at);
COMMENT ON TABLE field_schemas IS '字段定义表';
COMMENT ON COLUMN field_schemas.type IS '字段类型: text, textarea, number, date, select, etc.';

-- 7. 推送配置表
CREATE TABLE IF NOT EXISTS push_configs (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    url VARCHAR(500) NOT NULL,
    method VARCHAR(20) DEFAULT 'POST',
    headers TEXT,
    template TEXT,
    enabled BOOLEAN DEFAULT true,
    user_id INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_push_configs_deleted_at ON push_configs(deleted_at);
COMMENT ON TABLE push_configs IS '推送配置表';
COMMENT ON COLUMN push_configs.method IS 'HTTP方法: GET, POST, PUT, DELETE';

-- 8. 推送历史表
CREATE TABLE IF NOT EXISTS push_histories (
    id SERIAL PRIMARY KEY,
    config_id INTEGER NOT NULL,
    config_name VARCHAR(255),
    url VARCHAR(500),
    method VARCHAR(20),
    content TEXT,
    response TEXT,
    success BOOLEAN,
    status_code INTEGER,
    duration BIGINT,
    user_id INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_push_histories_config_id ON push_histories(config_id);
CREATE INDEX IF NOT EXISTS idx_push_histories_created_at ON push_histories(created_at);
COMMENT ON TABLE push_histories IS '推送历史表';
COMMENT ON COLUMN push_histories.duration IS '请求时长(毫秒)';

-- 9. 配置字段关联表
CREATE TABLE IF NOT EXISTS config_field_relations (
    id SERIAL PRIMARY KEY,
    config_id INTEGER NOT NULL,
    field_id INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_config_field_relations_config_id ON config_field_relations(config_id);
CREATE INDEX IF NOT EXISTS idx_config_field_relations_field_id ON config_field_relations(field_id);
COMMENT ON TABLE config_field_relations IS '推送配置与字段关联表';

-- 10. 工作流表
CREATE TABLE IF NOT EXISTS workflows (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    nodes TEXT,
    enabled BOOLEAN DEFAULT true,
    user_id INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_workflows_deleted_at ON workflows(deleted_at);
CREATE INDEX IF NOT EXISTS idx_workflows_user_id ON workflows(user_id);
COMMENT ON TABLE workflows IS '工作流表';
COMMENT ON COLUMN workflows.nodes IS 'JSON格式的节点配置';

-- 11. 工作流日志表
CREATE TABLE IF NOT EXISTS workflow_logs (
    id SERIAL PRIMARY KEY,
    workflow_id INTEGER NOT NULL,
    name VARCHAR(255),
    status VARCHAR(50),
    input TEXT,
    output TEXT,
    error TEXT,
    node_logs TEXT,
    duration BIGINT,
    started_at TIMESTAMP WITH TIME ZONE,
    finished_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_workflow_logs_workflow_id ON workflow_logs(workflow_id);
CREATE INDEX IF NOT EXISTS idx_workflow_logs_created_at ON workflow_logs(created_at);
COMMENT ON TABLE workflow_logs IS '工作流执行日志表';
COMMENT ON COLUMN workflow_logs.status IS '执行状态: running, success, error';
COMMENT ON COLUMN workflow_logs.node_logs IS 'JSON格式的节点执行日志';

-- 12. 定时任务表
CREATE TABLE IF NOT EXISTS schedule_tasks (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    cron_expr VARCHAR(100) NOT NULL,
    task_type VARCHAR(50) NOT NULL,
    config_id INTEGER,
    app_id INTEGER,
    workflow_id INTEGER,
    field_data TEXT,
    enabled BOOLEAN DEFAULT false,
    last_run_at TIMESTAMP WITH TIME ZONE,
    next_run_at TIMESTAMP WITH TIME ZONE,
    user_id INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_schedule_tasks_deleted_at ON schedule_tasks(deleted_at);
CREATE INDEX IF NOT EXISTS idx_schedule_tasks_enabled ON schedule_tasks(enabled);
COMMENT ON TABLE schedule_tasks IS '定时任务表';
COMMENT ON COLUMN schedule_tasks.task_type IS '任务类型: push, app, workflow';
COMMENT ON COLUMN schedule_tasks.cron_expr IS 'Cron表达式';

-- 13. 定时任务日志表
CREATE TABLE IF NOT EXISTS schedule_logs (
    id SERIAL PRIMARY KEY,
    task_id INTEGER NOT NULL,
    task_name VARCHAR(255),
    status VARCHAR(50),
    message TEXT,
    duration BIGINT,
    started_at TIMESTAMP WITH TIME ZONE,
    finished_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_schedule_logs_task_id ON schedule_logs(task_id);
CREATE INDEX IF NOT EXISTS idx_schedule_logs_created_at ON schedule_logs(created_at);
COMMENT ON TABLE schedule_logs IS '定时任务执行日志表';
COMMENT ON COLUMN schedule_logs.status IS '执行状态: success, error';

-- 14. 系统日志表
CREATE TABLE IF NOT EXISTS system_logs (
    id SERIAL PRIMARY KEY,
    module VARCHAR(100),
    action VARCHAR(100),
    content TEXT,
    user_id INTEGER,
    username VARCHAR(255),
    ip VARCHAR(50),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_system_logs_created_at ON system_logs(created_at);
CREATE INDEX IF NOT EXISTS idx_system_logs_user_id ON system_logs(user_id);
COMMENT ON TABLE system_logs IS '系统操作日志表';

-- 15. 菜单表
CREATE TABLE IF NOT EXISTS menus (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    path VARCHAR(255),
    icon VARCHAR(100),
    parent_id INTEGER,
    sort INTEGER DEFAULT 0,
    visible BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_menus_parent_id ON menus(parent_id);
CREATE INDEX IF NOT EXISTS idx_menus_sort ON menus(sort);
COMMENT ON TABLE menus IS '系统菜单表';

-- 外键约束
ALTER TABLE apps ADD CONSTRAINT fk_apps_user FOREIGN KEY (user_id) REFERENCES "user"(id) ON DELETE CASCADE;
ALTER TABLE app_logs ADD CONSTRAINT fk_app_logs_app FOREIGN KEY (app_id) REFERENCES apps(id) ON DELETE CASCADE;
ALTER TABLE app_configs ADD CONSTRAINT fk_app_configs_app FOREIGN KEY (app_id) REFERENCES apps(id) ON DELETE CASCADE;
ALTER TABLE push_configs ADD CONSTRAINT fk_push_configs_user FOREIGN KEY (user_id) REFERENCES "user"(id) ON DELETE CASCADE;
ALTER TABLE push_histories ADD CONSTRAINT fk_push_histories_config FOREIGN KEY (config_id) REFERENCES push_configs(id) ON DELETE CASCADE;
ALTER TABLE config_field_relations ADD CONSTRAINT fk_config_field_config FOREIGN KEY (config_id) REFERENCES push_configs(id) ON DELETE CASCADE;
ALTER TABLE config_field_relations ADD CONSTRAINT fk_config_field_field FOREIGN KEY (field_id) REFERENCES field_schemas(id) ON DELETE CASCADE;
ALTER TABLE workflows ADD CONSTRAINT fk_workflows_user FOREIGN KEY (user_id) REFERENCES "user"(id) ON DELETE CASCADE;
ALTER TABLE workflow_logs ADD CONSTRAINT fk_workflow_logs_workflow FOREIGN KEY (workflow_id) REFERENCES workflows(id) ON DELETE CASCADE;
ALTER TABLE schedule_tasks ADD CONSTRAINT fk_schedule_tasks_config FOREIGN KEY (config_id) REFERENCES push_configs(id) ON DELETE SET NULL;
ALTER TABLE schedule_tasks ADD CONSTRAINT fk_schedule_tasks_app FOREIGN KEY (app_id) REFERENCES apps(id) ON DELETE SET NULL;
ALTER TABLE schedule_tasks ADD CONSTRAINT fk_schedule_tasks_workflow FOREIGN KEY (workflow_id) REFERENCES workflows(id) ON DELETE SET NULL;
ALTER TABLE schedule_logs ADD CONSTRAINT fk_schedule_logs_task FOREIGN KEY (task_id) REFERENCES schedule_tasks(id) ON DELETE CASCADE;

-- 初始化管理员账号 (密码: admin123)
INSERT INTO "user" (username, nickname, password, role) 
VALUES ('admin', '管理员', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'admin')
ON CONFLICT (username) DO NOTHING;

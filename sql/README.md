# 数据库脚本说明

## 文件说明

- `schema.sql` - 完整的数据库表结构定义
- `drop_all.sql` - 删除所有表的脚本（谨慎使用）
- `add_workflow_to_schedule.sql` - 为定时任务添加工作流支持的迁移脚本

## 使用方法

### 初始化数据库

```bash
# 连接到PostgreSQL数据库
psql -U your_username -d your_database -f sql/schema.sql
```

### 清空数据库（谨慎！）

```bash
psql -U your_username -d your_database -f sql/drop_all.sql
```

### 应用迁移

```bash
psql -U your_username -d your_database -f sql/add_workflow_to_schedule.sql
```

## 表结构说明

### 核心表

1. **user** - 用户表
   - 存储用户账号信息
   - 默认管理员账号: admin / admin123

2. **apps** - 应用表
   - 存储用户创建的应用代码
   - 支持多种编程语言

3. **app_logs** - 应用日志表
   - 记录应用执行历史

### 应用商店

4. **app_templates** - 应用模板表
   - 应用商店中的可安装模板

5. **app_configs** - 应用配置表
   - 应用的配置参数

### 字段管理

6. **field_schemas** - 字段定义表
   - 定义可复用的字段模板

### 推送集成

7. **push_configs** - 推送配置表
   - HTTP推送的配置信息

8. **push_histories** - 推送历史表
   - 推送执行记录

9. **config_field_relations** - 配置字段关联表
   - 推送配置与字段的关联

### 工作流

10. **workflows** - 工作流表
    - 工作流定义和节点配置

11. **workflow_logs** - 工作流日志表
    - 工作流执行记录

### 定时任务

12. **schedule_tasks** - 定时任务表
    - 定时任务配置
    - 支持执行应用、推送、工作流

13. **schedule_logs** - 定时任务日志表
    - 定时任务执行记录

### 系统管理

14. **system_logs** - 系统日志表
    - 系统操作审计日志

15. **menus** - 菜单表
    - 系统菜单配置

## 注意事项

1. 项目使用GORM自动迁移，通常不需要手动执行SQL脚本
2. 这些脚本主要用于：
   - 新环境快速初始化
   - 数据库备份恢复
   - 了解完整的表结构
3. 生产环境请做好数据备份后再执行任何DDL操作

package workflow

import "gorm.io/gorm"

type model struct {
	db *gorm.DB
}

func New(db *gorm.DB) Model {
	return &model{db: db}
}

func (m *model) Create(workflow *Workflow) error {
	return m.db.Create(workflow).Error
}

func (m *model) GetByID(id uint) (*Workflow, error) {
	var workflow Workflow
	err := m.db.First(&workflow, id).Error
	if err != nil {
		return nil, err
	}
	return &workflow, nil
}

func (m *model) List(userID uint) ([]*Workflow, error) {
	var workflows []*Workflow
	query := m.db.Order("created_at DESC")
	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	}
	err := query.Find(&workflows).Error
	return workflows, err
}

func (m *model) Update(workflow *Workflow) error {
	return m.db.Save(workflow).Error
}

func (m *model) Delete(id uint) error {
	return m.db.Delete(&Workflow{}, id).Error
}

func (m *model) GetLogs(workflowID *uint, page, pageSize int) ([]*WorkflowLog, int64, error) {
	var logs []*WorkflowLog
	var total int64

	query := m.db.Model(&WorkflowLog{})
	if workflowID != nil && *workflowID > 0 {
		query = query.Where("workflow_id = ?", *workflowID)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&logs).Error

	return logs, total, err
}

func (m *model) CreateLog(log *WorkflowLog) error {
	return m.db.Create(log).Error
}

func (m *model) UpdateLog(log *WorkflowLog) error {
	return m.db.Save(log).Error
}

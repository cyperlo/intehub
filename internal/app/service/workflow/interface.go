package workflow

import (
	workflowModel "intehub/internal/app/models/workflow"
)

type Service interface {
	Create(workflow *workflowModel.Workflow) error
	GetByID(id uint) (*workflowModel.Workflow, error)
	List(userID uint) ([]*workflowModel.Workflow, error)
	Update(workflow *workflowModel.Workflow) error
	Delete(id uint) error
	Run(id uint, input map[string]interface{}) (*workflowModel.WorkflowLog, error)
	GetLogs(workflowID *uint, page, pageSize int) ([]*workflowModel.WorkflowLog, int64, error)
}

package workflow

type Model interface {
	Create(workflow *Workflow) error
	GetByID(id uint) (*Workflow, error)
	List(userID uint) ([]*Workflow, error)
	Update(workflow *Workflow) error
	Delete(id uint) error
	GetLogs(workflowID *uint, page, pageSize int) ([]*WorkflowLog, int64, error)
	CreateLog(log *WorkflowLog) error
	UpdateLog(log *WorkflowLog) error
}

package system

import systemModel "intehub/internal/app/models/system"

type service struct {
	model systemModel.Model
}

func New(model systemModel.Model) Service {
	return &service{model: model}
}

// Menu methods
func (s *service) GetMenus() ([]*systemModel.Menu, error) {
	return s.model.GetMenus()
}

func (s *service) CreateMenu(menu *systemModel.Menu) error {
	return s.model.CreateMenu(menu)
}

func (s *service) UpdateMenu(menu *systemModel.Menu) error {
	return s.model.UpdateMenu(menu)
}

func (s *service) DeleteMenu(id uint) error {
	return s.model.DeleteMenu(id)
}

// Log methods
func (s *service) GetLogs(page, pageSize int) ([]*systemModel.SystemLog, int64, error) {
	return s.model.GetLogs(page, pageSize)
}

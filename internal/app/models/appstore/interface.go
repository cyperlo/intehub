package appstore

type Model interface {
	// AppTemplate methods
	CreateTemplate(template *AppTemplate) error
	GetTemplateByID(id uint) (*AppTemplate, error)
	GetTemplateByName(name string) (*AppTemplate, error)
	ListTemplates(category string) ([]*AppTemplate, error)
	UpdateTemplate(template *AppTemplate) error
	DeleteTemplate(id uint) error

	// AppConfig methods
	CreateConfig(config *AppConfig) error
	GetConfig(appID uint, key string) (*AppConfig, error)
	GetConfigsByAppID(appID uint) ([]*AppConfig, error)
	UpdateConfig(config *AppConfig) error
	DeleteConfig(id uint) error
	DeleteConfigsByAppID(appID uint) error
}

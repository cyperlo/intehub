package user

type Model interface {
	// 查询方法
	GetUserByUsername(username string) (*DataObject, error)
	GetUserByID(id uint) (*DataObject, error)
	GetUsers() ([]*DataObject, error)

	// 创建和更新方法
	CreateUser(user *DataObject) error
	UpdateUser(user *DataObject) error
	UpdateUserFields(id uint, fields map[string]interface{}) error

	// 删除方法
	DeleteUser(id uint) error
}

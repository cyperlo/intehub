package user

type Model interface {
	GetUserByUsername(username string) (*DataObject, error)
}

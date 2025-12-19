package models

import "intehub/internal/app/models/user"

type Model interface {
	UserModel() user.Model
}

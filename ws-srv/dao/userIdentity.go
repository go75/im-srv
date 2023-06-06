package dao

import (
	"im-srv/ws-srv/global"
	"im-srv/ws-srv/model"

	"gorm.io/gorm"
)

func QueryUserIdentity(identity *model.UserIdentity) *gorm.DB {
	return global.DB.Where("name=?", identity.Name).First(identity)
}

func CreateUserIdentity(identity *model.UserIdentity) *gorm.DB {
	return global.DB.Create(identity)
}

func DeleteUserIdentity(identity *model.UserIdentity) *gorm.DB {
	return global.DB.Where("name = ?", identity.Name).Delete(identity)
}
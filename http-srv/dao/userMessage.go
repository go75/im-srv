package dao

import (
	"im-srv/http-srv/global"
	"im-srv/http-srv/model"

	"gorm.io/gorm"
)

func CreateUserMessage(msg *model.UserMessage) *gorm.DB {
	return global.DB.Create(msg)
}

func QueryUserMessageByIds(id1, id2 uint32) ([]*model.UserMessage, *gorm.DB) {
	var list []*model.UserMessage
	db := global.DB.Where("(sender_id=? and receiver_id=?) or (receiver_id=? and sender_id=?)", id1, id2, id1, id2).Order("id").Find(&list)
	return list, db
}
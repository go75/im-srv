package dao

import (
	"im-srv/ws-srv/global"
	"im-srv/ws-srv/model"

	"gorm.io/gorm"
)

func CreateGroupSession(session *model.GroupSession) *gorm.DB {
	return global.DB.Create(session)
}

func DeleteGroupSession(GroupSession *model.GroupSession) *gorm.DB {
	return global.DB.Delete(GroupSession)
}

func QueryGroupSessionsByUserId(userId uint) ([]*model.GroupSession, *gorm.DB){
	res := make([]*model.GroupSession, 0)
	db := global.DB.Where("user_id", userId).Find(&res)
	return res, db
}

func QueryGroupIdsByUserId(userId uint32) ([]uint32, *gorm.DB) {
	res := make([]uint32, 0)
	db := global.DB.Model(&model.GroupSession{}).Select("group_id").Where("user_id=?", userId).Find(&res)
	return res, db
}

func QueryGroupSessionsByGroupId(session *model.GroupSession) ([]*model.GroupSession, *gorm.DB) {
	res := make([]*model.GroupSession, 0)
	db := global.DB.Where("group_id", session.GroupId).Find(&res)
	return res, db	
}

func QueryUserIdsByGroupId(id uint32) ([]uint32, *gorm.DB) {
	ids := make([]uint32, 0)
	db := global.DB.Model(&model.GroupSession{}).Select("user_id").Where("group_id=?", id).Find(&ids)
	return ids, db
}

func QueryGroupSessionUserIdsByGroupId(session *model.GroupSession) ([]*model.GroupSession, *gorm.DB){
	res := make([]*model.GroupSession, 0)
	db := global.DB.Select("user_id").Where("group_id=?", session.GroupId).Find(&res)
	return res, db
}
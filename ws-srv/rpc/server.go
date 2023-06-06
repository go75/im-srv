package rpc

import (
	"context"
	"errors"
	"im-srv/ws-srv/conversion"
	"im-srv/ws-srv/dao"
	"im-srv/ws-srv/global"
	"im-srv/ws-srv/model"
	"im-srv/ws-srv/pb"
	"log"
)

type WebSocketServer struct {
}

var _ pb.WebSocketServer = WebSocketServer{}

// ================================================好友模块================================================
// 获取好友会话
func (s WebSocketServer) GetFriendSession(ctx context.Context, id2 *pb.Id2) (*pb.UserMessages, error) {
	msgs, _ := dao.QueryUserMessageByIds(id2.SenderId, id2.ProcessId)
	res := conversion.UserMessages(msgs)

	return res, nil
}

// 获取新好友信息
func (s WebSocketServer) GetNewFriend(ctx context.Context, id *pb.ProcessId) (*pb.Users, error) {
	ids, _ := dao.GetAddUserSenderIdsByReveiverId(id.Id)
	if len(ids) == 0 {
		return nil, errors.New("操作失败")
	}
	users, _ := dao.GetUsersByIds(ids)
	res := conversion.Users(users)

	return res, nil
}

// 用户消息持久化
func (s WebSocketServer) SaveUserMessage(ctx context.Context, msg *pb.UserMessage) (*pb.Ok, error) {
	res := conversion.UserMessage2(msg)
	db := dao.CreateUserMessage(res)
	if db.RowsAffected == 1 {
		return &pb.Ok{
			Ok: true,
		}, nil
	}
	return &pb.Ok{
		Ok: false,
	}, nil
}

// 获取好友列表
func (s WebSocketServer) GetFriendList(ctx context.Context, id *pb.ProcessId) (*pb.Users, error) {
	// 查询用户会话
	users := make([]*model.User, 0)

	db := global.DB.Raw("select user.id, user.name, user.introduce from user, user_session where user_session.small_id = ? and user.id = user_session.big_id", id.Id).Scan(&users)
	log.Println(users)

	tmp := make([]*model.User, 0)
	db.Raw("select user.id, user.name, user.introduce from user, user_session where user_session.big_id = ? and user.id = user_session.small_id", id.Id).Scan(&tmp)
	log.Println(tmp)

	users = append(users, tmp...)

	us := conversion.Users(users)

	return us, nil
}

// 添加好友
func (s WebSocketServer) AddFriend(ctx context.Context, id2 *pb.Id2) (*pb.User, error) {
	// 新好友消息持久化
	msg := &model.AddUserMessage{
		SenderId:   id2.SenderId,
		ReceiverId: id2.ProcessId,
	}
	dao.CreateAddUserMessage(msg)
	user := &model.User{
		ID: msg.SenderId,
	}
	dao.GetUserById(user)

	res := conversion.User(user)
	return res, nil
}

// 通过用户名称模糊查询用户
func (s WebSocketServer) GetFuzzyUserByUserName(ctx context.Context, name *pb.Name) (*pb.Users, error) {
	ls, _ := dao.GetFuzzyUserByUserName(name.Name)

	us := conversion.Users(ls)
	return us, nil
}

// 同意新好友请求
func (s WebSocketServer) AgreeNewFriend(ctx context.Context, id2 *pb.Id2) (*pb.User, error) {
	var session *model.UserSession
	if id2.ProcessId < id2.SenderId {
		session = &model.UserSession{
			SmallId: id2.ProcessId,
			BigId:   id2.SenderId,
		}
	} else if id2.ProcessId > id2.SenderId {
		session = &model.UserSession{
			SmallId: id2.SenderId,
			BigId:   id2.ProcessId,
		}
	} else {
		return nil, errors.New("操作失败")
	}

	// 开启事务
	transaction := global.DB.Begin()
	db := dao.DeleteAddUserMessage(id2.ProcessId, id2.SenderId)
	if db.RowsAffected != 1 {
		// 事务回滚
		transaction.Rollback()
		return nil, errors.New("操作失败")
	}
	db = dao.CreateUserSession(session)
	if db.RowsAffected != 1 {
		// 事务回滚
		transaction.Rollback()
		return nil, errors.New("操作失败")
	}

	// 提交事务
	transaction.Commit()

	user := &model.User{
		ID: id2.SenderId,
	}
	dao.GetUserById(user)

	res := conversion.User(user)
	return res, nil
}

// 拒绝新好友请求
func (s WebSocketServer) RefuseNewFriend(ctx context.Context, id2 *pb.Id2) (*pb.Ok, error) {
	if dao.DeleteAddUserMessage(id2.ProcessId, id2.SenderId).RowsAffected != 1 {
		return &pb.Ok{
			Ok: false,
		}, nil
	}
	return &pb.Ok{
		Ok: true,
	}, nil
}

// ================================================群聊模块================================================
// 获取群聊会话
func (s WebSocketServer) GetGroupSession(ctx context.Context, id *pb.ProcessId) (*pb.GroupMessages, error) {
	msgs, _ := dao.QueryGroupMessage(id.Id)
	if len(msgs) == 0 {
		return nil, nil
	}

	res := conversion.GroupMessages(msgs)

	return res, nil
}

// 获取新群聊信息
func (s WebSocketServer) GetNewGroup(ctx context.Context, id *pb.ProcessId) (*pb.NewGroupMessages, error) {
	// 获取当前用户创建过的所有群聊id
	groupIds, _ := dao.QueryGroupIdsByMasterId(id.Id)

	if groupIds == nil {
		return nil, errors.New("操作失败")
	}

	msgs, _ := dao.GetAddGroupMessagesByGroupIds(groupIds)
	log.Println(msgs)
	if len(msgs) == 0 {
		return nil, errors.New("操作失败")
	}

	res := conversion.NewGroupMessages(msgs)

	return res, nil
}

// 发送群聊消息
func (s WebSocketServer) SendGroupMsg(ctx context.Context, message *pb.GroupMessage) (*pb.SendGroupMsgRes, error) {
	user := &model.User{ID: message.SenderId}
	dao.GetUserNameById(user)
	if user.Name == "" {
		// 该用户不存在
		return nil, errors.New("操作失败")
	}

	// 群聊消息持久化
	msg := conversion.GroupMessage2(message)

	// 消息持久化
	dao.CreateGroupMessage(msg)

	// 发布消息
	ids, _ := dao.QueryUserIdsByGroupId(message.GroupId)

	return &pb.SendGroupMsgRes{
		UserIds: &pb.Ids{
			Ids: ids,
		},
		Message: conversion.GroupMessage(msg),
	}, nil
}

// 获取群聊列表
func (s WebSocketServer) GetGroupList(ctx context.Context, id *pb.ProcessId) (*pb.Groups, error) {
	ids, _ := dao.QueryGroupIdsByUserId(id.Id)
	if len(ids) == 0 {
		return nil, errors.New("操作失败")
	}
	groups, db := dao.QueryGroupsByGroupIds(ids)
	if db.Error != nil {
		return nil, errors.New("操作失败")
	}

	return conversion.Groups(groups), nil
}

// 添加群聊
func (s WebSocketServer) AddGroup(ctx context.Context, req *pb.AddGroupReq) (*pb.AddGroupRes, error) {
	// 新群友消息持久化
	msg := &model.AddGroupMessage{
		SenderId: req.Id2.SenderId,
		GroupId:  req.Id2.ProcessId,
	}

	if dao.CreateAddGroupMessage(msg).RowsAffected != 1 {
		// 请求失败
		log.Println("创建添加群聊消息失败")
		return nil, errors.New("操作失败")
	}

	group := &model.Group{
		ID: req.Id2.ProcessId,
	}

	dao.QueryGroupByGroupId(group)
	if group.Name == "" {
		log.Println("群聊查询失败")
		return nil, errors.New("操作失败")
	}

	return conversion.AddGroup(group), nil
}

// 通过群聊名称模糊查询群聊
func (s WebSocketServer) GetFuzzyGroupByGroupName(ctx context.Context, name *pb.Name) (*pb.Groups, error) {
	ls, _ := dao.GetFuzzyGroupByGroupName(name.Name)

	return conversion.Groups(ls), nil
}

// 同意新群友请求
func (s WebSocketServer) AgreeNewGroup(ctx context.Context, req *pb.AgreeNewGroupReq) (*pb.AgreeNewGroupRes, error) {
	user := &model.User{
		Name: req.SenderName,
	}

	// 通过 加好友的请求方名称 获取 加好友的请求方id
	dao.GetUserByName(user)

	if user.ID == 0 {
		return nil, errors.New("操作失败")
	}

	var session = &model.GroupSession{
		GroupId: req.GroupId,
		UserId:  user.ID,
	}

	msg := &model.AddGroupMessage{
		SenderId: user.ID,
		GroupId:  req.GroupId,
	}
	// 开启事务
	transaction := global.DB.Begin()

	db := dao.DeleteAddGroupMessageBySenderIdAndGroupId(msg)
	if db.RowsAffected != 1 {
		// 事务回滚
		transaction.Rollback()
		return nil, errors.New("操作失败")
	}
	log.Println("create session", session)
	db = dao.CreateGroupSession(session)
	if db.RowsAffected != 1 {
		// 事务回滚
		transaction.Rollback()
		return nil, errors.New("操作失败")
	}

	group := &model.Group{
		ID: req.GroupId,
	}

	dao.QueryGroupByGroupId(group)
	if group.Name == "" {
		log.Printf("id为%d的群聊信息获取失败\n", group.ID)
		transaction.Rollback()
		return nil, errors.New("操作失败")
	}

	// 提交事务
	transaction.Commit()
	// 返回新群聊的信息
	return &pb.AgreeNewGroupRes{
		Group:      conversion.Group(group),
		ReceiverId: user.ID,
	}, nil
}

// 拒绝新群友请求
func (s WebSocketServer) RefuseNewGroup(ctx context.Context, id2 *pb.Id2) (*pb.Ok, error) {
	msg := &model.AddGroupMessage{
		SenderId: id2.SenderId,
		GroupId:  id2.ProcessId,
	}
	db := dao.DeleteAddGroupMessageBySenderIdAndGroupId(msg)
	if db.RowsAffected != 1 {
		return &pb.Ok{
			Ok: false,
		}, nil
	}

	return &pb.Ok{
		Ok: true,
	}, nil
}

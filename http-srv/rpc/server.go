package rpc

import (
	"context"
	"fmt"
	"im-srv/http-srv/dao"
	"im-srv/http-srv/global"
	"im-srv/http-srv/model"
	"im-srv/http-srv/pb"
	"im-srv/http-srv/utils"
	"log"
	"math/rand"
)

type HTTPServer struct {
}

var _ pb.HTTPServer = HTTPServer{}

// 创建用户
func (s HTTPServer) UserRegist(ctx context.Context, info *pb.UserRegistInfo) (*pb.UserRegistRes, error) {
	log.Println("user regist")
	salt := fmt.Sprintf("%10d", rand.Int31())
	identity := &model.UserIdentity{
		Name: info.Name,
		Pwd: utils.MakePwd(info.Pwd, salt),
		Salt: salt,
	}
	user := &model.User {
		Name: identity.Name,
	}

	// 先将用户信息存入数据库
	transaction := global.DB.Begin()

	db := dao.CreateUserIdentity(identity)
	if db.RowsAffected != 1 {
		transaction.Rollback()
		log.Println("用户标记信息新增失败")
		return &pb.UserRegistRes{
			Ok: false,
			Msg: "用户标记信息新增失败",
		}, nil
	}

	db = dao.CreateUser(user)
	if db.RowsAffected != 1 {
		transaction.Rollback()
		log.Println("用户新增失败")
		return &pb.UserRegistRes{
			Ok: false,
			Msg: "用户新增失败",
		}, nil
	}
	transaction.Commit()
	return &pb.UserRegistRes{
		Ok: true,
		Msg: "注册成功",
	}, nil
}

// 用户登录
func (s HTTPServer) UserLogin(ctx context.Context, info *pb.UserLoginInfo) (*pb.UserLoginRes, error) {
	log.Println("user login")
	var identity model.UserIdentity

	global.DB.Raw("select * from user_identity where name=?", info.Name).Scan(&identity)
	//global.DB.Where("name=?", identity.Name).Find(identity)
	//dao.QueryUserIdentity(identity)
	if ok := utils.CheckPwd(info.Pwd, identity.Salt, identity.Pwd); ok {
		
		user := &model.User{
			Name: info.Name,
		}
		
		dao.GetUserByName(user)
		if user.ID == 0 {
			// 未查询到用户
			return nil, nil
		}

		token, err := utils.GenerateToken(user.ID, user.Name)
		if err != nil {
			log.Println("generate token err: ", err)
			return nil, nil
		}

		log.Println("generate token: ", token)

		return &pb.UserLoginRes{
			ID: user.ID,
			Name: user.Name,
			Token: token,
		}, nil

	} else {
		return nil, nil
	}
}

// 创建群聊
func (s HTTPServer) GroupRegist(ctx context.Context, info *pb.GroupRegistInfo) (*pb.GroupRegistRes, error) {

	log.Println(info.Name, info.Introduce, info.MasterId)
	group := &model.Group{
		Name: info.Name,
		MasterId: info.MasterId,
		Introduce: info.Introduce,
	}
	
	transaction := global.DB.Begin()

	if dao.CreateGroup(group).RowsAffected != 1 {
		transaction.Rollback()
		return &pb.GroupRegistRes{
			Ok: false,
			Msg: "创建失败",
		}, nil
	}

	if dao.GetGroupByName(group).Error != nil {
		transaction.Rollback()
		return &pb.GroupRegistRes{
			Ok: false,
			Msg: "创建失败",
		}, nil
	}

	session := model.GroupSession{
		GroupId: group.ID,
		UserId: info.MasterId,
	}

	if dao.CreateGroupSession(&session).RowsAffected != 1 {
		transaction.Rollback()
		return &pb.GroupRegistRes{
			Ok: false,
			Msg: "创建失败",
		}, nil
	}

	return &pb.GroupRegistRes{
		Ok: true,
		Msg: "创建成功",
		GroupId: group.ID,
	}, nil
}
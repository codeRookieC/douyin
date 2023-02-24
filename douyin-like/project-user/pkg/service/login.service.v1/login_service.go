package login_service_v1

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/jinzhu/copier"
	"go.uber.org/zap"
	"log"
	"strconv"
	"strings"
	common "test.com/project-common"
	"test.com/project-common/encrypts"
	"test.com/project-common/errs"
	"test.com/project-common/jwts"
	"test.com/project-common/logs"
	"test.com/project-grpc/user/login"
	"test.com/project-user/config"
	"test.com/project-user/internal/dao"
	"test.com/project-user/internal/data/member"
	"test.com/project-user/internal/data/organization"
	"test.com/project-user/internal/database"
	"test.com/project-user/internal/database/tran"
	"test.com/project-user/internal/repo"
	"test.com/project-user/pkg/model"
	"time"
)

type LoginService struct {
	login.UnimplementedLoginServiceServer
	cache            repo.Cache
	memberRepo       repo.MemberRepo
	organizationRepo repo.OrganizationRepo
	transaction      tran.Transaction
}

func New() *LoginService {
	return &LoginService{
		cache:            dao.Rc,
		memberRepo:       dao.NewMemberDao(),
		organizationRepo: dao.NewOrganizationDao(),
		transaction:      dao.NewTransaction(),
	}
}

func (ls *LoginService) GetCaptcha(ctx context.Context, msg *login.CaptchaMessage) (*login.CaptchaResponse, error) {
	//1.获取手机号
	mobile := msg.Mobile
	//2.验证是否合法
	if !common.VerifyMobile(mobile) {
		return nil, errs.GrpcError(model.NoLegalMobile)
	}
	//3.生成验证码
	code := "123456"
	//4.发送验证码
	go func() {
		time.Sleep(2 * time.Second)
		//log.Println("短信平台调用成功 发送短信")
		zap.L().Info("短信平台调用成功 发送短信 INFO")
		logs.LG.Debug("短信平台调用成功 发送短信 DEBUG")
		zap.L().Error("短信平台调用成功 发送短信 ERROR")
		//存储验证码到redis当中 过期时间是15分钟
		//redis  假设之后我们要改成存储在mysql中 mongodb中 或者其他数据库中的话 我们要直接在这里改代码 耦合性太强  考虑使用接口的方式
		//面向接口编程 高耦合 低内聚
		c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		err := ls.cache.Put(c, model.RegisterRedisKey+mobile, code, 15*time.Minute)
		if err != nil {
			log.Printf("验证码存入redis出错, cause by : %v \n", err)
		}
		log.Printf("将手机号和验证码存入redis成功: REGISTER_%s : %s", mobile, code)
	}()
	return &login.CaptchaResponse{Code: code}, nil
}

func (ls *LoginService) Register(ctx context.Context, msg *login.RegisterMessage) (*login.RegisterResponse, error) {
	c := context.Background()
	//1.可以校验参数
	//2.校验验证码
	redisCode, err := ls.cache.Get(c, model.RegisterRedisKey+msg.Mobile)
	if err == redis.Nil { //验证码不存在或者已经过期
		return nil, errs.GrpcError(model.CaptchaNotExist)
	}
	if err != nil {
		zap.L().Error("Register redis get error", zap.Error(err))
		return nil, errs.GrpcError(model.RedisError)
	}
	if redisCode != msg.Captcha {
		return nil, errs.GrpcError(model.CaptchaError)
	}
	//3.校验业务逻辑(邮箱是否被注册, 用户名是否被注册, 手机是否被注册)
	exist, err := ls.memberRepo.GetMemberByEmail(c, msg.Email)
	if err != nil {
		zap.L().Error("Register db get error", zap.Error(err))
		return nil, errs.GrpcError(model.DbError)
	}
	if exist {
		return nil, errs.GrpcError(model.EmailExist)
	}
	exist, err = ls.memberRepo.GetMemberByAccount(c, msg.Name)
	if err != nil {
		zap.L().Error("Register db get error", zap.Error(err))
		return nil, errs.GrpcError(model.DbError)
	}
	if exist {
		return nil, errs.GrpcError(model.AccountExist)
	}
	exist, err = ls.memberRepo.GetMemberByMobile(c, msg.Mobile)
	if err != nil {
		zap.L().Error("Register db get error", zap.Error(err))
		return nil, errs.GrpcError(model.DbError)
	}
	if exist {
		return nil, errs.GrpcError(model.MobileExist)
	}
	//4.执行业务 将数据存入member表中 然后将数据存入organization表中
	pwd := encrypts.Md5(msg.Password)
	mem := &member.Member{
		Account:       msg.Name,
		Password:      pwd,
		Name:          msg.Name,
		Mobile:        msg.Mobile,
		Email:         msg.Email,
		CreateTime:    time.Now().UnixMilli(),
		LastLoginTime: time.Now().UnixMilli(),
		Status:        model.Normal,
	}
	//存储到两个表中应该是事务操作 要么都成功 要么都失败
	err = ls.transaction.Action(func(conn database.DbConn) error {
		err = ls.memberRepo.SaveMember(conn, c, mem)
		if err != nil {
			zap.L().Error("Register db SaveMember error", zap.Error(err))
			return errs.GrpcError(model.DbError)
		}
		//存入organization表
		org := &organization.Organization{
			Name:       mem.Name + "个人组织",
			MemberId:   mem.Id,
			CreateTime: time.Now().UnixMilli(),
			Personal:   model.Personal,
			Avatar:     "https://gimg2.baidu.com/image_search/src=http%3A%2F%2Fc-ssl.dtstatic.com%2Fuploads%2Fblog%2F202103%2F31%2F20210331160001_9a852.thumb.1000_0.jpg&refer=http%3A%2F%2Fc-ssl.dtstatic.com&app=2002&size=f9999,10000&q=a80&n=0&g=0n&fmt=auto?sec=1673017724&t=ced22fc74624e6940fd6a89a21d30cc5",
		}
		err = ls.organizationRepo.SaveOrganization(conn, c, org)
		if err != nil {
			zap.L().Error("register SaveOrganization db err", zap.Error(err))
			return errs.GrpcError(model.DbError)
		}
		return nil
	})
	//5.返回
	return &login.RegisterResponse{}, err
}

func (ls *LoginService) Login(ctx context.Context, msg *login.LoginMessage) (*login.LoginResponse, error) {
	c := context.Background()
	//1.去数据库中查询账号密码是否正确
	pwd := encrypts.Md5(msg.Password)
	mem, err := ls.memberRepo.FindMember(c, msg.Account, pwd)
	if err != nil {
		zap.L().Error("Login db get error", zap.Error(err))
		return nil, errs.GrpcError(model.DbError)
	}
	if mem == nil {
		return nil, errs.GrpcError(model.AccountOrPwdError)
	}
	memMsg := &login.MemberMessage{}
	err = copier.Copy(memMsg, mem)
	//如果向外暴露id 会有安全隐患 所以这里采取aes加密
	memMsg.Code, _ = encrypts.EncryptInt64(mem.Id, model.AESKey)
	//2.根据用户id查组织
	orgs, err := ls.organizationRepo.FindOrganizationByMemId(c, mem.Id)
	if err != nil {
		zap.L().Error("Login db get error", zap.Error(err))
		return nil, errs.GrpcError(model.DbError)
	}
	var orgsMessage []*login.OrganizationMessage
	err = copier.Copy(&orgsMessage, orgs)
	for _, v := range orgsMessage {
		v.Code, _ = encrypts.EncryptInt64(v.Id, model.AESKey)
	}
	//3.用jwt生成token 使用用户id(主键)来创建token
	memIdStr := strconv.FormatInt(mem.Id, 10) //将十进制的mem.Id 转换成string类型
	exp := time.Duration(config.C.JwtConfig.AccessExp * 3600 * 24) * time.Second
	rExp := time.Duration(config.C.JwtConfig.RefreshExp * 3600 * 24) * time.Second
	token := jwts.CreateToken(memIdStr, exp, config.C.JwtConfig.AccessSecret, rExp, config.C.JwtConfig.RefreshSecret)
	tokenList := &login.TokenMessage{
		AccessToken:    token.AccessToken,
		AccessTokenExp: token.AccessExp,
		RefreshToken:   token.RefreshToken,
		TokenType:      "bearer",
	}
	return &login.LoginResponse{
		Member:           memMsg,
		OrganizationList: orgsMessage,
		TokenList:        tokenList,
	}, nil
}
func (ls *LoginService) TokenVerify(ctx context.Context, msg *login.LoginMessage) (*login.LoginResponse, error) {
	token := msg.Token
	if strings.Contains(token, "bearer") {
		token = strings.ReplaceAll(token, "bearer", "")
	}
	parseToken, err := jwts.ParseToken(token, config.C.JwtConfig.AccessSecret)
	if err != nil {
		zap.L().Error("Login TokenVerify error", zap.Error(err))
		return nil, errs.GrpcError(model.NoLogin)
	}
	//数据库查询 优化点：TODO 登录之后将用户信息缓存 redis
	id, _ := strconv.ParseInt(parseToken, 10, 64)
	memberById, err := ls.memberRepo.FindMemberById(context.Background(), id)
	if err != nil {
		zap.L().Error("TokenVerify db FindMemberById error", zap.Error(err))
		return nil, errs.GrpcError(model.DbError)
	}
	memMsg := &login.MemberMessage{}
	copier.Copy(memMsg, memberById)
	//如果向外暴露id 会有安全隐患 所以这里采取aes加密
	memMsg.Code, _ = encrypts.EncryptInt64(memberById.Id, model.AESKey)
	return &login.LoginResponse{Member: memMsg}, nil
}

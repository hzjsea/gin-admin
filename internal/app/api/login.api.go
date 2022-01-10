package api

import (
	"fmt"
	"github.com/LyricTian/captcha"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"

	"github.com/LyricTian/gin-admin/v8/internal/app/config"
	"github.com/LyricTian/gin-admin/v8/internal/app/contextx"
	"github.com/LyricTian/gin-admin/v8/internal/app/ginx"
	"github.com/LyricTian/gin-admin/v8/internal/app/schema"
	"github.com/LyricTian/gin-admin/v8/internal/app/service"
	"github.com/LyricTian/gin-admin/v8/pkg/errors"
	"github.com/LyricTian/gin-admin/v8/pkg/logger"
)

var LoginSet = wire.NewSet(wire.Struct(new(LoginAPI), "*"))

type LoginAPI struct {
	LoginSrv *service.LoginSrv
}

func (a *LoginAPI) GetCaptcha(c *gin.Context) {
	ctx := c.Request.Context()
	item, err := a.LoginSrv.GetCaptcha(ctx, config.C.Captcha.Length)
	logger.Infof("验证码ID是%s", item)
	if err != nil {
		ginx.ResError(c, err)
		return
	}

	ginx.ResSuccess(c, item)
}

// ResCaptcha 图片验证码
func (a *LoginAPI) ResCaptcha(c *gin.Context) {
	ctx := c.Request.Context()
	captchaID := c.Query("id")
	if captchaID == "" {
		ginx.ResError(c, errors.New400Response("captcha id not empty"))
		return
	}

	if c.Query("reload") != "" {
		if !captcha.Reload(captchaID) {
			ginx.ResError(c, errors.New400Response("not found captcha id"))
			return
		}
	}

	cfg := config.C.Captcha
	logger.Infof("配置文件 %s", cfg)
	err := a.LoginSrv.ResCaptcha(ctx, c.Writer, captchaID, cfg.Width, cfg.Height)
	if err != nil {
		ginx.ResError(c, err)
	}
}


// Login 登录流程
// 需要检验两个内容
// - 验证码
// - 根据验证码获取的图形验证码, 这里先给他设置为true pass: 因为是底层写入的
// - 用户密码
func (a *LoginAPI) Login(c *gin.Context) {
	ctx := c.Request.Context()
	var item schema.LoginParam
	if err := ginx.ParseJSON(c, &item); err != nil {
		ginx.ResError(c, err)
		return
	}

	//logger.Infof("解析后的内容%v", item)

	// 检验验证码
	//logger.Infof("验证码ID %v 验证码Code %v", item.CaptchaID, item.CaptchaCode)
	if captcha.VerifyString(item.CaptchaID, item.CaptchaCode){
	//if !captcha.VerifyString(item.CaptchaID, item.CaptchaCode) {
		ginx.ResError(c, errors.New400Response("无效的验证码"))
		return
	}

	//logger.Infof("加密后的密码应该是 %s", hash.SHA1String(item.Password))
	// 检验用户准确性
	user, err := a.LoginSrv.Verify(ctx, item.UserName, item.Password)
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	//logger.Infof("%s",json.MarshalToString(user))

	// 检验token
	tokenInfo, err := a.LoginSrv.GenerateToken(ctx, a.formatTokenUserID(user.ID, user.UserName))
	if err != nil {
		ginx.ResError(c, err)
		return
	}

	ctx = logger.NewUserIDContext(ctx, user.ID)
	ctx = logger.NewUserNameContext(ctx, user.UserName)
	ctx = logger.NewTagContext(ctx, "__login__")
	logger.WithContext(ctx).Infof("login")

	ginx.ResSuccess(c, tokenInfo)
}

func (a *LoginAPI) formatTokenUserID(userID uint64, userName string) string {
	return fmt.Sprintf("%d-%s", userID, userName)
}

func (a *LoginAPI) Logout(c *gin.Context) {
	ctx := c.Request.Context()

	userID := contextx.FromUserID(ctx)
	if userID != 0 {
		ctx = logger.NewTagContext(ctx, "__logout__")
		err := a.LoginSrv.DestroyToken(ctx, ginx.GetToken(c))
		if err != nil {
			logger.WithContext(ctx).Errorf(err.Error())
		}
		logger.WithContext(ctx).Infof("logout")
	}
	ginx.ResOK(c)
}

func (a *LoginAPI) RefreshToken(c *gin.Context) {
	ctx := c.Request.Context()
	tokenInfo, err := a.LoginSrv.GenerateToken(ctx, a.formatTokenUserID(contextx.FromUserID(ctx), contextx.FromUserName(ctx)))
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResSuccess(c, tokenInfo)
}

func (a *LoginAPI) GetUserInfo(c *gin.Context) {
	ctx := c.Request.Context()
	info, err := a.LoginSrv.GetLoginInfo(ctx, contextx.FromUserID(ctx))
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResSuccess(c, info)
}

func (a *LoginAPI) QueryUserMenuTree(c *gin.Context) {
	ctx := c.Request.Context()
	menus, err := a.LoginSrv.QueryUserMenuTree(ctx, contextx.FromUserID(ctx))
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResList(c, menus)
}

func (a *LoginAPI) UpdatePassword(c *gin.Context) {
	ctx := c.Request.Context()
	var item schema.UpdatePasswordParam
	if err := ginx.ParseJSON(c, &item); err != nil {
		ginx.ResError(c, err)
		return
	}

	err := a.LoginSrv.UpdatePassword(ctx, contextx.FromUserID(ctx), item)
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResOK(c)
}

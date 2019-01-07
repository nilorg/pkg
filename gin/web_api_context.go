package gin

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/nilorg/sdk/errors"
)

// WebAPIContext Web上下文
type WebAPIContext struct {
	*gin.Context
}

// SetCurrentAccount 设置当前账户
func (ctx *WebAPIContext) SetCurrentAccount(data interface{}) error {
	session := sessions.Default(ctx.Context)
	session.Set(CurrentAccount, data)
	return session.Save()
}

// GetCurrentAccount 设置当前账户
func (ctx *WebAPIContext) GetCurrentAccount() interface{} {
	session := sessions.Default(ctx.Context)
	return session.Get(CurrentAccount)
}

// DelCurrentAccount 删除当前账户
func (ctx *WebAPIContext) DelCurrentAccount() error {
	session := sessions.Default(ctx.Context)
	session.Delete(CurrentAccount)
	return session.Save()
}

// ResultError 返回错误
func (ctx *WebAPIContext) ResultError(err error) {
	if berr, ok := err.(errors.BusinessError); ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": berr,
		})
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": errors.New(0, err.Error()),
		})
	}
}

// WebAPIControllerFunc WebAPI控制器函数
func WebAPIControllerFunc(ctlFunc func(ctx *WebAPIContext)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tmplCtx := &WebAPIContext{
			Context: ctx,
		}
		ctlFunc(tmplCtx)
	}
}

package gin

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// WebAPIContext Web上下文
type WebAPIContext struct {
	*gin.Context
}

// SetCurrentAccount 设置当前账户
func (ctx *WebAPIContext) SetCurrentAccount(data interface{}) {
	session := sessions.Default(ctx.Context)
	session.Set(currentAccount, data)
	session.Save()
}

// GetCurrentAccount 设置当前账户
func (ctx *WebAPIContext) GetCurrentAccount() interface{} {
	session := sessions.Default(ctx.Context)
	return session.Get(currentAccount)
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

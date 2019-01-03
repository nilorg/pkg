package gin

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// CurrentAccount ...
const CurrentAccount = "current_account"

// WebContext Web上下文
type WebContext struct {
	*gin.Context
	pageName string
}

// RenderPage 渲染页面
func (ctx *WebContext) RenderPage(data gin.H) {
	layout := "layout.tmpl"
	if ctx.GetHeader("X-PJAX") == "true" {
		layout = "pjax_layout.tmpl"
	}
	tmplName := fmt.Sprintf("%s_pages_%s", layout, ctx.pageName)
	if data == nil {
		data = gin.H{
			CurrentAccount: ctx.GetCurrentAccount(),
		}
	} else {
		data[CurrentAccount] = ctx.GetCurrentAccount()
	}
	ctx.HTML(http.StatusOK, tmplName, data)
}

// RenderSinglePage 渲染单页面
func (ctx *WebContext) RenderSinglePage(data gin.H) {
	tmplName := fmt.Sprintf("singles_%s.tmpl", ctx.pageName)
	if data == nil {
		data = gin.H{
			CurrentAccount: ctx.GetCurrentAccount(),
		}
	} else {
		data[CurrentAccount] = ctx.GetCurrentAccount()
	}
	ctx.HTML(http.StatusOK, tmplName, data)
}

// SetCurrentAccount 设置当前账户
func (ctx *WebContext) SetCurrentAccount(data interface{}) error {
	session := sessions.Default(ctx.Context)
	session.Set(CurrentAccount, data)
	return session.Save()
}

// GetCurrentAccount 设置当前账户
func (ctx *WebContext) GetCurrentAccount() interface{} {
	session := sessions.Default(ctx.Context)
	return session.Get(CurrentAccount)
}

// DelCurrentAccount 删除当前账户
func (ctx *WebContext) DelCurrentAccount() error {
	session := sessions.Default(ctx.Context)
	session.Delete(CurrentAccount)
	return session.Save()
}

// WebControllerFunc Web控制器函数
func WebControllerFunc(ctlFunc func(ctx *WebContext), pageName string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tmplCtx := &WebContext{
			Context:  ctx,
			pageName: pageName,
		}
		ctlFunc(tmplCtx)
	}
}

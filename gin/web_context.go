package gin

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

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
	ctx.HTML(http.StatusOK, tmplName, data)
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

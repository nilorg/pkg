
# 模板结构
```bash
./templates
└── default
    ├── errors
    ├── layouts
    │   ├── layout.tmpl
    │   └── pjax_layout.tmpl
    └── pages
        ├── posts
        │   ├── posts.tmpl
        └── posts_detail
            └── posts_detail.tmpl
```
# 安装插件
```bash
go get -u github.com/gin-contrib/multitemplate

go get -u github.com/nilorg/pkg/gin
```

# 加载模板
```go
// 主题名称
const themeName string = "default"

func loadTemplates(templatesDir string) multitemplate.Render {
	r := multitemplate.New()
	// 加载布局
	layouts, err := filepath.Glob(filepath.Join(templatesDir, themeName, "layouts/*.tmpl"))
	if err != nil {
		panic(err.Error())
	}
	// 加载错误页面
	errors, err := filepath.Glob(filepath.Join(templatesDir, themeName, "errors/*.tmpl"))
	if err != nil {
		panic(err.Error())
	}
	for _, errPage := range errors {
		tmplName := fmt.Sprintf("error_%s", filepath.Base(errPage))
		r.AddFromFiles(tmplName, errPage)
	}
	// 页面文件夹
	pages, err := ioutil.ReadDir(filepath.Join(templatesDir, themeName, "pages"))
	if err != nil {
		panic(err.Error())
	}
	for _, page := range pages {
		if !page.IsDir() {
			continue
		}
		for _, layout := range layouts {
			pageItems, err := filepath.Glob(filepath.Join(templatesDir, themeName, fmt.Sprintf("pages/%s/*.tmpl", page.Name())))
			if err != nil {
				panic(err.Error())
			}
			files := []string{
				layout,
			}
			files = append(files,pageItems...)
			tmplName := fmt.Sprintf("%s_pages_%s", filepath.Base(layout), page.Name())
			fmt.Println(tmplName, files)
			r.AddFromFiles(tmplName, files...)
		}
	}
	return r
}
```
# 使用模板
```go
engine := gin.Default()
engine.HTMLRender = loadTemplates("./templates")
```
# 给路由配置模板
```go 
import ngin "github.com/nilorg/pkg/gin"

engine.GET("/detail", ngin.WebControllerFunc(func(ctx *ngin.WebContext) {
		ctx.RenderPage(gin.H{
			"title":"标题",
		})
	}, "posts_detail"))
```
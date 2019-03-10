package gin

import (
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
)

// LoadTemplateFunc 加载模板函数类
type LoadTemplateFunc func(templatesDir string) multitemplate.Render

// EngineTemplate gin引擎模板
type EngineTemplate struct {
	templatesDir     string
	engine           *gin.Engine
	watcher          *fsnotify.Watcher
	Errors           <-chan error
	loadTemplateFunc LoadTemplateFunc
}

// NewEngineTemplate 创建一个gin引擎模板
func NewEngineTemplate(templateDir string, engine *gin.Engine, tmplFunc LoadTemplateFunc) (*EngineTemplate, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	return &EngineTemplate{
		templatesDir:     templateDir,
		engine:           engine,
		loadTemplateFunc: tmplFunc,
		watcher:          watcher,
		Errors:           make(<-chan error),
	}, nil
}

// Watching 监听模板文件夹中是否有变动
func (tmpl *EngineTemplate) Watching() error {
	tmpl.Errors = tmpl.watcher.Errors

	go func() {
		for {
			event := <-tmpl.watcher.Events
			loadFlag := true
			switch event.Op {
			case fsnotify.Create:
				fileInfo, err := os.Stat(event.Name)
				if err == nil && fileInfo.IsDir() {
					tmpl.watcher.Add(event.Name)
				}
			case fsnotify.Remove, fsnotify.Rename:
				fileInfo, err := os.Stat(event.Name)
				if err == nil && fileInfo.IsDir() {
					tmpl.watcher.Remove(event.Name)
				}
			case fsnotify.Chmod:
				loadFlag = false
			}
			if loadFlag {
				tmpl.engine.HTMLRender = tmpl.LoadTemplate()
			}
		}
	}()

	//遍历目录下的所有子目录
	err := filepath.Walk(tmpl.templatesDir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			err := tmpl.watcher.Add(path)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

// Close 关闭
func (tmpl *EngineTemplate) Close() error {
	return tmpl.watcher.Close()
}

// LoadTemplate 加载模板
func (tmpl *EngineTemplate) LoadTemplate() multitemplate.Render {
	return tmpl.loadTemplateFunc(tmpl.templatesDir)
}

package server

import (
	"embed"
	"io/fs"
	"moeCounter/server/controller"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 注册基础路由（静态文件和首页）
func registerBaseRoutes(router *gin.Engine, publicFS embed.FS) {
	// 首页路由
	router.GET("/", func(c *gin.Context) {
		data, err := fs.ReadFile(publicFS, "index.html")
		if err != nil {
			c.String(http.StatusInternalServerError, "Internal Server Error")
			return
		}
		c.Data(http.StatusOK, "text/html; charset=utf-8", data)
	})

	// favicon图标路由
	router.GET("/favicon.ico", func(c *gin.Context) {
		data, err := fs.ReadFile(publicFS, "favicon.ico")
		if err != nil {
			c.String(http.StatusInternalServerError, "Internal Server Error")
			return
		}
		c.Data(http.StatusOK, "image/x-icon", data)
	})

	// 静态文件服务（挂载到/public/assets路径）
	fsAssets, _ := fs.Sub(publicFS, "assets")
	router.StaticFS("/assets", http.FS(fsAssets))
}

// 注册API路由
func registerAPIRoutes(apiGroup *gin.RouterGroup) {
	// 计数器接口
	apiGroup.GET("/counter", controller.CounterHandler)

	// 主题列表接口
	apiGroup.GET("/themes", controller.ThemeListHandler)
}

// 启动服务器
func RunServer(router *gin.Engine, port int) {
	if err := router.Run(":" + strconv.Itoa(port)); err != nil {
		router.Run() // 使用随机端口
	}
}

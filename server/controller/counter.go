package controller

import (
	"log"
	"math/rand"
	"moeCounter/internal/database"
	"moeCounter/internal/utils"
	"moeCounter/public"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// 处理计数器请求
func CounterHandler(c *gin.Context) {
	var req CounterRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 简单校验name，防止注入
	validName := regexp.MustCompile(`^[a-zA-Z0-9_-]{1,64}$`)
	if !validName.MatchString(req.Name) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "非法的name参数"})
		return
	}

	publicFS := public.Public

	// 随机主题（如果未提供）
	if req.Theme == "" {
		themes, err := ListThemes(publicFS)
		if err != nil || len(themes) == 0 {
			req.Theme = "original-new"
		} else {
			r := rand.New(rand.NewSource(time.Now().UnixNano()))
			req.Theme = themes[r.Intn(len(themes))]
		}
	}

	// 设置默认值
	if req.Length == 0 {
		req.Length = 7
	}
	if req.Scale == 0 {
		req.Scale = 1
	}
	if req.Offset == 0 {
		req.Offset = 0
	}
	if req.Align == "" {
		req.Align = "center"
	}
	if req.Pixelate == "" {
		req.Pixelate = "off"
	}

	var count uint
	var err error
	if req.Num != "" {
		var numVal uint64
		numVal, err = strconv.ParseUint(req.Num, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的num参数"})
			return
		}
		count = uint(numVal)
	} else {
		count, err = database.IncrementCounter(req.Name)
		if err != nil {
			log.Printf("[ERROR] 数据库错误: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器内部错误"})
			return
		}
		if req.Base > 0 {
			count += uint(req.Base)
		}
	}

	svg, err := utils.CombineImages(count, publicFS, req.Theme, req.Length, req.Scale, req.Offset, req.Align, req.Pixelate, req.Darkmode)
	if err != nil {
		log.Printf("[ERROR] 图片生成失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "图片生成失败"})
		return
	}

	// 设置安全响应头
	c.Header("Cache-Control", "no-store")
	c.Header("X-Content-Type-Options", "nosniff")
	c.Header("X-Frame-Options", "DENY")

	c.Header("Content-Type", "image/svg+xml")
	c.String(http.StatusOK, svg)
}

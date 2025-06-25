package controller

import (
	"math/rand"
	"moeCounter/internal/database"
	"moeCounter/internal/utils"
	"moeCounter/public"
	"net/http"
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

	// 使用全局public文件系统
	publicFS := public.Public

	// 设置默认值
	if req.Length == 0 {
		req.Length = 7
	}
	if req.Scale == 0 {
		req.Scale = 1
	}
	if req.Offset == 0 {
		req.Offset = 0 // 默认偏移量为0
	}

	// 随机主题（如果未提供）
	if req.Theme == "" {
		// 获取所有可用主题
		themes, err := ListThemes(publicFS)
		if err != nil || len(themes) == 0 {
			req.Theme = "original-new" // 默认值
		} else {
			// 随机选择一个主题
			r := rand.New(rand.NewSource(time.Now().UnixNano()))
			req.Theme = themes[r.Intn(len(themes))]
		}
	}

	// 设置新参数的默认值
	if req.Align == "" {
		req.Align = "center" // 默认居中
	}
	if req.Pixelate == "" {
		req.Pixelate = "off" // 默认关闭像素化
	}

	var count uint
	var err error
	// 如果提供了num参数，直接使用num值
	if req.Num != "" {
		// 尝试将num转换为整数
		var numVal uint64
		numVal, err = strconv.ParseUint(req.Num, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的num参数"})
			return
		}
		count = uint(numVal)
	} else {
		// 否则增加计数器
		count, err = database.IncrementCounter(req.Name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "数据库错误"})
			return
		}

		// 如果提供了base参数，将计数值与base相加
		if req.Base > 0 {
			// 确保不会出现负数
			count += uint(req.Base)
		}
	}

	// 生成SVG图片（使用文件系统）
	svg, err := utils.CombineImages(count, publicFS, req.Theme, req.Length, req.Scale, req.Offset, req.Align, req.Pixelate, req.Darkmode)
	if err != nil {
		// 添加详细错误日志
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "图片生成失败",
			"details": err.Error(),
		})
		return
	}

	// 直接返回SVG
	c.Header("Content-Type", "image/svg+xml")
	c.String(http.StatusOK, svg)
}

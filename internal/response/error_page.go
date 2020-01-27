package response

import "github.com/gin-gonic/gin"

// ErrTemplateFile is
var ErrTemplateFile = "error.tpl"

// RenderErrorPage is
func RenderErrorPage(c *gin.Context, errCode int, errMsg string) {
	c.AbortWithStatus(errCode)
	c.HTML(errCode, ErrTemplateFile, gin.H{
		"code":    errCode,
		"message": errMsg,
	})
}

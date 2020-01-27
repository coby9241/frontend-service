package response

import "github.com/gin-gonic/gin"

// ErrTemplateFile is
var ErrTemplateFile = "error.tpl"

// RenderErrorPage is
func RenderErrorPage(c *gin.Context, errCode int, err error) {
	c.AbortWithStatus(errCode)
	c.HTML(errCode, ErrTemplateFile, gin.H{
		"code":    errCode,
		"message": err.Error(),
	})
}

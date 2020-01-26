package response

import "github.com/gin-gonic/gin"

// RenderErrorPage is
func RenderErrorPage(c *gin.Context, errCode int, errMsg string) {
	c.AbortWithStatus(errCode)
	c.HTML(errCode, "error.tpl", gin.H{
		"code":    errCode,
		"message": errMsg,
	})
}

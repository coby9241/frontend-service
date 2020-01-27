package response_test

import (
	"errors"
	"html/template"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/coby9241/frontend-service/internal/response"
	"github.com/gin-gonic/gin"
)

func TestRenderErrorPage(t *testing.T) {
	// setup test context
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, router := gin.CreateTestContext(resp)
	// setup variable for test template

	tpl := template.Must(template.New(ErrTemplateFile).Parse(`Error: {{.code}} {{.message}}`))
	router.SetHTMLTemplate(tpl)

	// render page
	RenderErrorPage(c, http.StatusBadGateway, errors.New("server error"))
}

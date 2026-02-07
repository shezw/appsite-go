package redirect

import (
	"github.com/gin-gonic/gin"
	"appsite-go/internal/apis/response"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) WechatCallback(c *gin.Context) {
	// TODO: Implement WeChat OAuth callback
	response.Success(c, gin.H{"status": "mock_wechat_ok"})
}

func (h *Handler) OSSCallback(c *gin.Context) {
	// TODO: Implement OSS Callback
	response.Success(c, gin.H{"status": "mock_oss_ok"})
}

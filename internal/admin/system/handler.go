package system

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"appsite-go/internal/apis/response"
	"appsite-go/internal/core/setting"
)

type Handler struct {
	cfg *setting.Config
}

func NewHandler(cfg *setting.Config) *Handler {
	return &Handler{
		cfg: cfg,
	}
}

type MenuItem struct {
	Key      string     `json:"key"`
	Label    string     `json:"label"`
	Icon     string     `json:"icon"`
	Path     string     `json:"path,omitempty"`
	Roles    []string   `json:"roles"`
	Children []MenuItem `json:"children,omitempty"`
}

func (h *Handler) GetMenu(c *gin.Context) {
	// In a real app, we would filter by the current user's role.
	// We read from the String configuration which is a JSON blob.
	
	if h.cfg.AdminMenu == "" {
		response.Success(c, []MenuItem{})
		return
	}

	var menu []MenuItem
	if err := json.Unmarshal([]byte(h.cfg.AdminMenu), &menu); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, menu)
}

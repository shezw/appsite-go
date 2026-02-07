package content

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"appsite-go/internal/apis/response"
	"appsite-go/internal/services/contents"
	"appsite-go/internal/services/contents/entity"
)

type Handler struct {
	articleSvc *contents.ArticleService
	bannerSvc  *contents.BannerService
}

func NewHandler(articleSvc *contents.ArticleService, bannerSvc *contents.BannerService) *Handler {
	return &Handler{
		articleSvc: articleSvc,
		bannerSvc:  bannerSvc,
	}
}

// --- Requests ---
type CreateArticleReq struct {
	Title       string `json:"title" binding:"required"`
	Type        string `json:"type"`
	Mode        string `json:"mode"`
	Content     string `json:"content"` // Maps to Introduce
	Cover       string `json:"cover"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

type UpdateArticleReq struct {
	Title       string `json:"title"`
	Content     string `json:"content"`
	Cover       string `json:"cover"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

// --- Article Handlers ---

func (h *Handler) CreateArticle(c *gin.Context) {
	var req CreateArticleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}

	article := &entity.Article{
		Title:       req.Title,
		Type:        req.Type,
		Mode:        req.Mode,
		Introduce:   req.Content,
		Cover:       req.Cover,
		Description: req.Description,
		Status:      req.Status,
	}
	// TODO: Assign AuthorID from context

	if err := h.articleSvc.Create(article); err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, gin.H{"id": article.ID})
}

func (h *Handler) GetArticle(c *gin.Context) {
	id := c.Param("id")
	article, err := h.articleSvc.Get(id)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, article)
}

func (h *Handler) ListArticles(c *gin.Context) {
	page := 1
	pageSize := 20
	if p := c.Query("page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil && v > 0 {
			page = v
		}
	}
	if ps := c.Query("page_size"); ps != "" {
		if v, err := strconv.Atoi(ps); err == nil && v > 0 {
			pageSize = v
		}
	}

	filters := make(map[string]interface{})
	if t := c.Query("type"); t != "" {
		filters["type"] = t
	}

	list, count, err := h.articleSvc.List(page, pageSize, filters)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{
		"list":  list,
		"total": count,
	})
}

func (h *Handler) UpdateArticle(c *gin.Context) {
	id := c.Param("id")
	var req UpdateArticleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}

	updates := make(map[string]interface{})
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Content != "" {
		updates["introduce"] = req.Content
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}

	if err := h.articleSvc.Update(id, updates); err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, nil)
}

func (h *Handler) DeleteArticle(c *gin.Context) {
	id := c.Param("id")
	if err := h.articleSvc.Delete(id); err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, nil)
}

// --- Banner Handlers (Simplified) ---

func (h *Handler) CreateBanner(c *gin.Context) {
	var banner entity.Banner
	if err := c.ShouldBindJSON(&banner); err != nil {
		response.Error(c, err)
		return
	}
	if err := h.bannerSvc.Create(&banner); err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, gin.H{"id": banner.ID})
}

func (h *Handler) ListBanners(c *gin.Context) {
	// Simplified list
	list, _, err := h.bannerSvc.List(1, 100, nil)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, list)
}

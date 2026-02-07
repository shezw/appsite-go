package contents

import (
"appsite-go/internal/apis/response"
apperr "appsite-go/internal/core/error"
"appsite-go/internal/services/contents"
"appsite-go/internal/services/contents/entity"
"github.com/gin-gonic/gin"
"strconv"
)

type Handler struct {
articleService *contents.ArticleService
bannerService  *contents.BannerService
}

func NewHandler(as *contents.ArticleService, bs *contents.BannerService) *Handler {
return &Handler{
articleService: as,
bannerService:  bs,
}
}

// ---- Articles ----

// ListArticles lists articles with pagination and filtering
func (h *Handler) ListArticles(c *gin.Context) {
page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

filters := make(map[string]interface{})
if status := c.Query("status"); status != "" {
filters["status"] = status
}
if categoryId := c.Query("category_id"); categoryId != "" {
filters["category_id"] = categoryId
}

list, total, err := h.articleService.List(page, size, filters)
if err != nil {
response.Error(c, apperr.NewWithMessage(apperr.ServerError, err.Error()))
return
}

response.Success(c, gin.H{
"list":  list,
"total": total,
"page":  page,
"size":  size,
})
}

// GetArticle gets a single article
func (h *Handler) GetArticle(c *gin.Context) {
id := c.Param("id")
article, err := h.articleService.Get(id)
if err != nil {
response.Error(c, apperr.NewWithMessage(apperr.NotFound, "Article not found"))
return
}
response.Success(c, article)
}

// CreateArticle creates a new article
func (h *Handler) CreateArticle(c *gin.Context) {
var article entity.Article
if err := c.ShouldBindJSON(&article); err != nil {
response.Error(c, apperr.NewWithMessage(apperr.InvalidParams, "Invalid request"))
return
}

if err := h.articleService.Create(&article); err != nil {
response.Error(c, apperr.NewWithMessage(apperr.ServerError, err.Error()))
return
}
response.Success(c, article)
}

// UpdateArticle updates an article
func (h *Handler) UpdateArticle(c *gin.Context) {
id := c.Param("id")
var updates map[string]interface{}
if err := c.ShouldBindJSON(&updates); err != nil {
response.Error(c, apperr.NewWithMessage(apperr.InvalidParams, "Invalid request"))
return
}

if err := h.articleService.Update(id, updates); err != nil {
response.Error(c, apperr.NewWithMessage(apperr.ServerError, err.Error()))
return
}
response.Success(c, nil)
}

// DeleteArticle deletes an article
func (h *Handler) DeleteArticle(c *gin.Context) {
id := c.Param("id")
if err := h.articleService.Delete(id); err != nil {
response.Error(c, apperr.NewWithMessage(apperr.ServerError, err.Error()))
return
}
response.Success(c, nil)
}

// ---- Banners ----

// ListBanners lists banners
func (h *Handler) ListBanners(c *gin.Context) {
page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

filters := make(map[string]interface{})
if status := c.Query("status"); status != "" {
filters["status"] = status
}
if position := c.Query("position"); position != "" {
filters["position"] = position
}

list, total, err := h.bannerService.List(page, size, filters)
if err != nil {
response.Error(c, apperr.NewWithMessage(apperr.ServerError, err.Error()))
return
}
response.Success(c, gin.H{
"list":  list,
"total": total,
"page":  page,
"size":  size,
})
}

// GetBanner gets a single banner
func (h *Handler) GetBanner(c *gin.Context) {
id := c.Param("id")
banner, err := h.bannerService.Get(id)
if err != nil {
response.Error(c, apperr.NewWithMessage(apperr.NotFound, "Banner not found"))
return
}
response.Success(c, banner)
}

// CreateBanner creates a new banner
func (h *Handler) CreateBanner(c *gin.Context) {
var banner entity.Banner
if err := c.ShouldBindJSON(&banner); err != nil {
response.Error(c, apperr.NewWithMessage(apperr.InvalidParams, "Invalid request"))
return
}

if err := h.bannerService.Create(&banner); err != nil {
response.Error(c, apperr.NewWithMessage(apperr.ServerError, err.Error()))
return
}
response.Success(c, banner)
}

// UpdateBanner updates a banner
func (h *Handler) UpdateBanner(c *gin.Context) {
id := c.Param("id")
var updates map[string]interface{}
if err := c.ShouldBindJSON(&updates); err != nil {
response.Error(c, apperr.NewWithMessage(apperr.InvalidParams, "Invalid request"))
return
}

if err := h.bannerService.Update(id, updates); err != nil {
response.Error(c, apperr.NewWithMessage(apperr.ServerError, err.Error()))
return
}
response.Success(c, nil)
}

// DeleteBanner deletes a banner
func (h *Handler) DeleteBanner(c *gin.Context) {
id := c.Param("id")
if err := h.bannerService.Delete(id); err != nil {
response.Error(c, apperr.NewWithMessage(apperr.ServerError, err.Error()))
return
}
response.Success(c, nil)
}

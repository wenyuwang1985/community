package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wenyuwang1985/community/internal/middleware"
	"github.com/wenyuwang1985/community/internal/service"
)

type MarketHandler struct {
	marketService *service.MarketService
}

func NewMarketHandler(marketService *service.MarketService) *MarketHandler {
	return &MarketHandler{marketService: marketService}
}

type createItemRequest struct {
	CommunityID int64    `json:"community_id" binding:"required"`
	Title       string   `json:"title" binding:"required,max=100"`
	Price       int      `json:"price" binding:"required,min=0"`
	Condition   string   `json:"condition" binding:"required,oneof=new like_new lightly_used heavily_used"`
	Category    string   `json:"category" binding:"required,oneof=appliance furniture book baby sports other"`
	Images      []string `json:"images" binding:"omitempty,max=9,dive,max=500"`
}

func (h *MarketHandler) CreateItem(c *gin.Context) {
	var req createItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.Error(c, http.StatusBadRequest, 400, "参数错误: "+err.Error())
		return
	}

	userID := c.GetInt64("userID")
	item, err := h.marketService.CreateItem(c.Request.Context(), userID, req.CommunityID, req.Title, req.Price, req.Condition, req.Category, req.Images)
	if err != nil {
		middleware.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	middleware.Success(c, gin.H{
		"id": item.ID,
	})
}

func (h *MarketHandler) GetItem(c *gin.Context) {
	idStr := c.Param("id")
	itemID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		middleware.Error(c, http.StatusBadRequest, 400, "商品 ID 格式错误")
		return
	}

	ctx := c.Request.Context()
	item, err := h.marketService.GetItemByID(ctx, itemID)
	if err != nil {
		middleware.Error(c, http.StatusNotFound, 404, err.Error())
		return
	}

	// 组装卖家信息
	seller, _ := h.marketService.GetUserByID(ctx, item.SellerID)
	sellerInfo := gin.H{}
	if seller != nil {
		sellerInfo = gin.H{
			"id":         seller.ID,
			"nickname":   seller.Nickname,
			"avatar_url": seller.AvatarURL,
		}
	}

	middleware.Success(c, gin.H{
		"id":           item.ID,
		"seller_id":    item.SellerID,
		"community_id": item.CommunityID,
		"title":        item.Title,
		"price":        item.Price,
		"condition":    item.Condition,
		"category":     item.Category,
		"images":       item.Images,
		"status":       item.Status,
		"created_at":   item.CreatedAt,
		"seller":       sellerInfo,
	})
}

type updateItemRequest struct {
	Title     string   `json:"title" binding:"required,max=100"`
	Price     int      `json:"price" binding:"required,min=0"`
	Condition string   `json:"condition" binding:"required,oneof=new like_new lightly_used heavily_used"`
	Category  string   `json:"category" binding:"required,oneof=appliance furniture book baby sports other"`
	Images    []string `json:"images" binding:"omitempty,max=9,dive,max=500"`
}

func (h *MarketHandler) UpdateItem(c *gin.Context) {
	idStr := c.Param("id")
	itemID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		middleware.Error(c, http.StatusBadRequest, 400, "商品 ID 格式错误")
		return
	}

	var req updateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.Error(c, http.StatusBadRequest, 400, "参数错误: "+err.Error())
		return
	}

	userID := c.GetInt64("userID")
	item, err := h.marketService.UpdateItem(c.Request.Context(), userID, itemID, req.Title, req.Price, req.Condition, req.Category, req.Images)
	if err != nil {
		middleware.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	middleware.Success(c, gin.H{
		"id": item.ID,
	})
}

func (h *MarketHandler) ListItems(c *gin.Context) {
	communityIDStr := c.Query("community_id")
	communityID, err := strconv.ParseInt(communityIDStr, 10, 64)
	if err != nil {
		middleware.Error(c, http.StatusBadRequest, 400, "community_id 格式错误")
		return
	}

	category := c.Query("category")
	lastIDStr := c.Query("last_id")
	lastID, _ := strconv.ParseInt(lastIDStr, 10, 64)
	limitStr := c.Query("limit")
	limit, _ := strconv.Atoi(limitStr)

	ctx := c.Request.Context()
	items, err := h.marketService.ListItems(ctx, communityID, category, lastID, limit)
	if err != nil {
		middleware.Error(c, http.StatusInternalServerError, 500, err.Error())
		return
	}

	var result []gin.H
	for _, item := range items {
		seller, _ := h.marketService.GetUserByID(ctx, item.SellerID)
		sellerInfo := gin.H{}
		if seller != nil {
			sellerInfo = gin.H{
				"id":         seller.ID,
				"nickname":   seller.Nickname,
				"avatar_url": seller.AvatarURL,
			}
		}
		result = append(result, gin.H{
			"id":           item.ID,
			"seller_id":    item.SellerID,
			"community_id": item.CommunityID,
			"title":        item.Title,
			"price":        item.Price,
			"condition":    item.Condition,
			"category":     item.Category,
			"images":       item.Images,
			"status":       item.Status,
			"created_at":   item.CreatedAt,
			"seller":       sellerInfo,
		})
	}

	middleware.Success(c, result)
}

func (h *MarketHandler) MarkSold(c *gin.Context) {
	idStr := c.Param("id")
	itemID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		middleware.Error(c, http.StatusBadRequest, 400, "商品 ID 格式错误")
		return
	}

	userID := c.GetInt64("userID")
	if err := h.marketService.MarkSold(c.Request.Context(), userID, itemID); err != nil {
		middleware.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	middleware.Success(c, nil)
}

func (h *MarketHandler) MarkOff(c *gin.Context) {
	idStr := c.Param("id")
	itemID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		middleware.Error(c, http.StatusBadRequest, 400, "商品 ID 格式错误")
		return
	}

	userID := c.GetInt64("userID")
	if err := h.marketService.MarkOff(c.Request.Context(), userID, itemID); err != nil {
		middleware.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	middleware.Success(c, nil)
}

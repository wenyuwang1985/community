package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wenyuwang1985/community/internal/middleware"
	"github.com/wenyuwang1985/community/internal/service"
)

type PostHandler struct {
	postService *service.PostService
}

func NewPostHandler(postService *service.PostService) *PostHandler {
	return &PostHandler{postService: postService}
}

type createPostRequest struct {
	CommunityID int64    `json:"community_id" binding:"required"`
	Tag         string   `json:"tag" binding:"omitempty,oneof=help share notice qa"`
	Content     string   `json:"content" binding:"required,max=2000"`
	Images      []string `json:"images" binding:"omitempty,max=9,dive,max=500"`
}

func (h *PostHandler) CreatePost(c *gin.Context) {
	var req createPostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.Error(c, http.StatusBadRequest, 400, "参数错误: "+err.Error())
		return
	}

	userID := c.GetInt64("userID")
	post, err := h.postService.CreatePost(c.Request.Context(), userID, req.CommunityID, req.Tag, req.Content, req.Images)
	if err != nil {
		middleware.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	middleware.Success(c, gin.H{
		"id": post.ID,
	})
}

func (h *PostHandler) GetPost(c *gin.Context) {
	idStr := c.Param("id")
	postID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		middleware.Error(c, http.StatusBadRequest, 400, "帖子 ID 格式错误")
		return
	}

	ctx := c.Request.Context()
	post, err := h.postService.GetPostByID(ctx, postID)
	if err != nil {
		middleware.Error(c, http.StatusNotFound, 404, err.Error())
		return
	}

	// 组装作者信息
	author, _ := h.postService.GetUserByID(ctx, post.UserID)
	authorInfo := gin.H{}
	if author != nil {
		authorInfo = gin.H{
			"id":         author.ID,
			"nickname":   author.Nickname,
			"avatar_url": author.AvatarURL,
		}
	}

	// 查询当前用户是否点赞
	var isLiked bool
	if userID, exists := c.Get("userID"); exists {
		isLiked, _ = h.postService.IsLiked(ctx, postID, userID.(int64))
	}

	middleware.Success(c, gin.H{
		"id":            post.ID,
		"user_id":       post.UserID,
		"community_id":  post.CommunityID,
		"tag":           post.Tag,
		"content":       post.Content,
		"images":        post.Images,
		"like_count":    post.LikeCount,
		"comment_count": post.CommentCount,
		"created_at":    post.CreatedAt,
		"author":        authorInfo,
		"is_liked":      isLiked,
	})
}

func (h *PostHandler) DeletePost(c *gin.Context) {
	idStr := c.Param("id")
	postID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		middleware.Error(c, http.StatusBadRequest, 400, "帖子 ID 格式错误")
		return
	}

	userID := c.GetInt64("userID")
	if err := h.postService.DeletePost(c.Request.Context(), userID, postID); err != nil {
		middleware.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	middleware.Success(c, nil)
}

func (h *PostHandler) ListPosts(c *gin.Context) {
	communityIDStr := c.Query("community_id")
	communityID, err := strconv.ParseInt(communityIDStr, 10, 64)
	if err != nil {
		middleware.Error(c, http.StatusBadRequest, 400, "community_id 格式错误")
		return
	}

	lastIDStr := c.Query("last_id")
	lastID, _ := strconv.ParseInt(lastIDStr, 10, 64)

	limitStr := c.Query("limit")
	limit, _ := strconv.Atoi(limitStr)

	ctx := c.Request.Context()
	posts, err := h.postService.ListPosts(ctx, communityID, lastID, limit)
	if err != nil {
		middleware.Error(c, http.StatusInternalServerError, 500, err.Error())
		return
	}

	userID := c.GetInt64("userID")
	var result []gin.H
	for _, post := range posts {
		author, _ := h.postService.GetUserByID(ctx, post.UserID)
		authorInfo := gin.H{}
		if author != nil {
			authorInfo = gin.H{
				"id":         author.ID,
				"nickname":   author.Nickname,
				"avatar_url": author.AvatarURL,
			}
		}
		isLiked, _ := h.postService.IsLiked(ctx, post.ID, userID)
		result = append(result, gin.H{
			"id":            post.ID,
			"user_id":       post.UserID,
			"community_id":  post.CommunityID,
			"tag":           post.Tag,
			"content":       post.Content,
			"images":        post.Images,
			"like_count":    post.LikeCount,
			"comment_count": post.CommentCount,
			"created_at":    post.CreatedAt,
			"author":        authorInfo,
			"is_liked":      isLiked,
		})
	}

	middleware.Success(c, result)
}

func (h *PostHandler) Like(c *gin.Context) {
	idStr := c.Param("id")
	postID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		middleware.Error(c, http.StatusBadRequest, 400, "帖子 ID 格式错误")
		return
	}

	userID := c.GetInt64("userID")
	if err := h.postService.Like(c.Request.Context(), postID, userID); err != nil {
		middleware.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	middleware.Success(c, nil)
}

func (h *PostHandler) Unlike(c *gin.Context) {
	idStr := c.Param("id")
	postID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		middleware.Error(c, http.StatusBadRequest, 400, "帖子 ID 格式错误")
		return
	}

	userID := c.GetInt64("userID")
	if err := h.postService.Unlike(c.Request.Context(), postID, userID); err != nil {
		middleware.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	middleware.Success(c, nil)
}

type createCommentRequest struct {
	Content string `json:"content" binding:"required,max=500"`
}

func (h *PostHandler) CreateComment(c *gin.Context) {
	idStr := c.Param("id")
	postID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		middleware.Error(c, http.StatusBadRequest, 400, "帖子 ID 格式错误")
		return
	}

	var req createCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.Error(c, http.StatusBadRequest, 400, "参数错误: "+err.Error())
		return
	}

	userID := c.GetInt64("userID")
	comment, err := h.postService.CreateComment(c.Request.Context(), postID, userID, req.Content)
	if err != nil {
		middleware.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	middleware.Success(c, gin.H{
		"id": comment.ID,
	})
}

func (h *PostHandler) ListComments(c *gin.Context) {
	idStr := c.Param("id")
	postID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		middleware.Error(c, http.StatusBadRequest, 400, "帖子 ID 格式错误")
		return
	}

	ctx := c.Request.Context()
	comments, err := h.postService.ListComments(ctx, postID)
	if err != nil {
		middleware.Error(c, http.StatusInternalServerError, 500, err.Error())
		return
	}

	var result []gin.H
	for _, comment := range comments {
		author, _ := h.postService.GetUserByID(ctx, comment.UserID)
		authorInfo := gin.H{}
		if author != nil {
			authorInfo = gin.H{
				"id":         author.ID,
				"nickname":   author.Nickname,
				"avatar_url": author.AvatarURL,
			}
		}
		result = append(result, gin.H{
			"id":         comment.ID,
			"post_id":    comment.PostID,
			"user_id":    comment.UserID,
			"content":    comment.Content,
			"created_at": comment.CreatedAt,
			"author":     authorInfo,
		})
	}

	middleware.Success(c, result)
}

package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wenyuwang1985/community/internal/middleware"
	"github.com/wenyuwang1985/community/internal/service"
)

type UploadHandler struct {
	uploadService *service.UploadService
}

func NewUploadHandler(uploadService *service.UploadService) *UploadHandler {
	return &UploadHandler{uploadService: uploadService}
}

func (h *UploadHandler) Upload(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		middleware.Error(c, http.StatusBadRequest, 400, "解析表单失败: "+err.Error())
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		middleware.Error(c, http.StatusBadRequest, 400, "未选择文件")
		return
	}

	urls, err := h.uploadService.SaveFiles(files)
	if err != nil {
		middleware.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	middleware.Success(c, gin.H{"urls": urls})
}

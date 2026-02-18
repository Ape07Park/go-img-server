package handler

import (
	"go-img-server/internal/storage"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ListHandler - 이미지 목록 조회 담당
type ListHandler struct {
	Storage storage.Storage
}

// List - GET /api/v1/projects/:project/images
//
// Spring 비유:
//   @GetMapping("/projects/{project}/images")
//   public ResponseEntity<List<FileInfo>> list(@PathVariable String project)
func (h *ListHandler) List(c *gin.Context) {
	project := c.Param("project")

	files, err := h.Storage.List(project)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "목록 조회 실패: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, files)
}

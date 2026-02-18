package handler

import (
	"go-img-server/internal/storage"
	"net/http"

	"github.com/gin-gonic/gin"
)

// DeleteHandler - 이미지 삭제 담당
type DeleteHandler struct {
	Storage storage.Storage
}

// Delete - DELETE /api/v1/projects/:project/images/:filename
//
// Spring 비유:
//   @DeleteMapping("/projects/{project}/images/{filename}")
//   public ResponseEntity<Void> delete(
//       @PathVariable String project,
//       @PathVariable String filename)
func (h *DeleteHandler) Delete(c *gin.Context) {
	project := c.Param("project")
	filename := c.Param("filename")

	if err := h.Storage.Delete(project, filename); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "파일을 찾을 수 없습니다"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "삭제되었습니다"})
}

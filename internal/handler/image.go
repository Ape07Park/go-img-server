package handler

import (
	"go-img-server/internal/storage"
	"io"
	"mime"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// ImageHandler - 이미지 서빙(보기/다운로드) 담당
type ImageHandler struct {
	Storage storage.Storage
}

// Serve - GET /i/:project/:filename
//
// 브라우저에서 직접 열거나, <img src="..."> 태그에서 호출하는 엔드포인트.
// 인증 없이 공개 접근 가능 (URL을 아는 사람만 접근).
//
// Spring 비유:
//   @GetMapping("/i/{project}/{filename}")
//   public ResponseEntity<StreamingResponseBody> serve(...)
func (h *ImageHandler) Serve(c *gin.Context) {
	project := c.Param("project")
	filename := c.Param("filename")

	file, err := h.Storage.Get(project, filename)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "파일을 찾을 수 없습니다"})
		return
	}
	defer file.Close()

	// 파일 확장자로 Content-Type 결정
	// 예: .jpg → image/jpeg, .png → image/png
	contentType := mime.TypeByExtension(filepath.Ext(filename))
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	c.Header("Content-Type", contentType)
	c.Header("Cache-Control", "public, max-age=2592000") // 30일 캐시

	// 파일 스트림을 응답으로 그대로 복사
	// Spring의 StreamingResponseBody 또는 Resource 반환과 동일
	io.Copy(c.Writer, file)
}

// Download - GET /api/v1/projects/:project/images/:filename/download
//
// 브라우저에서 파일 다운로드 대화상자를 띄우는 엔드포인트.
// Content-Disposition: attachment 헤더를 추가하는 것 외에는 Serve와 동일.
func (h *ImageHandler) Download(c *gin.Context) {
	project := c.Param("project")
	filename := c.Param("filename")

	file, err := h.Storage.Get(project, filename)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "파일을 찾을 수 없습니다"})
		return
	}
	defer file.Close()

	contentType := mime.TypeByExtension(filepath.Ext(filename))
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	c.Header("Content-Type", contentType)
	// attachment → 브라우저가 파일 저장 대화상자를 띄움
	c.Header("Content-Disposition", "attachment; filename=\""+filename+"\"")

	io.Copy(c.Writer, file)
}

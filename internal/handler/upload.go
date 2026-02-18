package handler

import (
	"go-img-server/internal/storage"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UploadHandler - 이미지 업로드 담당
// Spring의 @RestController + @PostMapping 과 동일
type UploadHandler struct {
	Storage storage.Storage // Spring의 @Autowired 필드 주입과 동일
}

// Upload - POST /api/v1/projects/:project/images
//
// Spring 비유:
//   @PostMapping("/projects/{project}/images")
//   public ResponseEntity<FileInfo> upload(
//       @PathVariable String project,
//       @RequestParam("file") MultipartFile file)
func (h *UploadHandler) Upload(c *gin.Context) {
	project := c.Param("project") // Spring의 @PathVariable

	// multipart/form-data 에서 "file" 필드 추출
	// Spring의 @RequestParam("file") MultipartFile 과 동일
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "파일을 찾을 수 없습니다. 필드명 'file'로 전송해주세요",
		})
		return
	}

	// 10MB 크기 제한 체크
	// Spring의 spring.servlet.multipart.max-file-size=10MB 와 동일
	const maxSize = 10 << 20 // 10MB
	if fileHeader.Size > maxSize {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "파일 크기는 10MB를 초과할 수 없습니다",
		})
		return
	}

	// 이미지 파일 타입만 허용
	if !isAllowedType(fileHeader.Header.Get("Content-Type")) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "이미지 파일만 업로드 가능합니다 (jpg, png, gif, webp)",
		})
		return
	}

	// 파일 열기
	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "파일 열기 실패"})
		return
	}
	defer file.Close()

	// 스토리지에 저장
	info, err := h.Storage.Save(project, fileHeader.Filename, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "파일 저장 실패: " + err.Error(),
		})
		return
	}

	// 201 Created + 저장된 파일 정보 반환
	// Spring의 ResponseEntity.created(...).body(info) 와 동일
	c.JSON(http.StatusCreated, info)
}

// isAllowedType - 허용된 이미지 MIME 타입 검사
func isAllowedType(contentType string) bool {
	allowed := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/gif":  true,
		"image/webp": true,
	}
	return allowed[contentType]
}

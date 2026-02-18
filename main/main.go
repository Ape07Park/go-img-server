package main

import (
	"fmt"
	"go-img-server/internal/config"
	"go-img-server/internal/handler"
	"go-img-server/internal/middleware"
	"go-img-server/internal/storage"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. 설정 로드 (환경변수 → 기본값)
	// Spring의 application.properties 로딩과 동일
	cfg := config.Load()

	// 2. 의존성 생성 및 주입
	// Spring의 @Bean / ApplicationContext와 동일한 역할을 수동으로 수행
	//
	// 나중에 MinIO로 전환할 때 이 줄만 바꾸면 됨:
	//   storage.NewLocalStorage(...)  →  storage.NewMinioStorage(...)
	store := storage.NewLocalStorage(cfg.UploadDir, cfg.BaseURL)

	uploadHandler := &handler.UploadHandler{Storage: store}
	listHandler := &handler.ListHandler{Storage: store}
	imageHandler := &handler.ImageHandler{Storage: store}
	deleteHandler := &handler.DeleteHandler{Storage: store}

	// 3. Gin 라우터 설정
	// Spring의 DispatcherServlet + @RequestMapping 설정과 동일
	r := gin.Default() // Logger + Recover 미들웨어 자동 포함

	// ── 공개 엔드포인트 (인증 불필요) ──────────────────────────────────
	// 이미지 보기: 브라우저 주소창, <img src=""> 태그에서 사용
	r.GET("/i/:project/:filename", imageHandler.Serve)

	// ── 인증 필요 API (/api/v1) ────────────────────────────────────────
	// Spring의 SecurityConfig에서 특정 경로에만 인증 적용하는 것과 동일
	api := r.Group("/api/v1")
	api.Use(middleware.APIKeyAuth(cfg.APIKey)) // X-API-Key 헤더 검사
	{
		projects := api.Group("/projects/:project")
		{
			// 업로드:   POST   /api/v1/projects/{project}/images
			projects.POST("/images", uploadHandler.Upload)

			// 목록:     GET    /api/v1/projects/{project}/images
			projects.GET("/images", listHandler.List)

			// 삭제:     DELETE /api/v1/projects/{project}/images/{filename}
			projects.DELETE("/images/:filename", deleteHandler.Delete)

			// 다운로드: GET    /api/v1/projects/{project}/images/{filename}/download
			projects.GET("/images/:filename/download", imageHandler.Download)
		}
	}

	// 4. 서버 시작
	addr := fmt.Sprintf(":%s", cfg.Port)
	fmt.Printf("이미지 서버 시작: http://localhost%s\n", addr)
	if err := r.Run(addr); err != nil {
		panic(err)
	}
}

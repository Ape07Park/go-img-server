package main

import (
	"fmt"
	"image-server/handler"
	"image-server/storage"
	"net/http"
)

func main() {
	// 1. 의존성 주입 (Dependency Injection)
	// Spring Bean 등록 과정을 수동으로 하는 것과 같습니다.
	// LocalStorage 생성 -> UploadHandler에 주입
	store := storage.NewLocalStorage("./uploads")
	uploadHandler := &handler.UploadHandler{Storage: store}

	// 2. 라우터 설정 (ServeMux)
	// Spring의 DispatcherServlet 설정과 유사
	mux := http.NewServeMux()

	// 3. 엔드포인트 등록
	// /upload 경로로 들어오면 uploadHandler의 HandleUpload 함수 실행
	mux.HandleFunc("/upload", uploadHandler.HandleUpload)

	// 4. 서버 시작
	fmt.Println("Server started on :8080")
	// Tomcat 시작과 동일
	if err := http.ListenAndServe(":8080", mux); err != nil {
		panic(err)
	}
}
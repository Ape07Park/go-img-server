package config

import "os"

// Config - 서버 설정 값 모음
// Spring의 application.properties / @ConfigurationProperties 와 동일한 역할
type Config struct {
	Port      string // 서버 포트 (기본: 8080)
	UploadDir string // 파일 저장 경로 (기본: ./uploads)
	BaseURL   string // 이미지 URL 앞에 붙는 도메인 (기본: http://localhost:8080)
	APIKey    string // API 인증 키 (기본: dev-secret-key)
}

// Load - 환경변수에서 설정 로드, 없으면 기본값 사용
// 배포 시: 환경변수로 값을 주입 (Docker -e, systemd Environment= 등)
func Load() *Config {
	return &Config{
		Port:      getEnv("PORT", "8080"),
		UploadDir: getEnv("UPLOAD_DIR", "./uploads"),
		BaseURL:   getEnv("BASE_URL", "http://localhost:8080"),
		APIKey:    getEnv("API_KEY", "dev-secret-key"),
	}
}

func getEnv(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}

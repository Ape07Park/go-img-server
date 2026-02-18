package storage

import "io"

// FileInfo - 파일 메타데이터
// Spring의 DTO(Data Transfer Object)와 동일한 역할
type FileInfo struct {
	Name    string `json:"name"`
	URL     string `json:"url"`
	Size    int64  `json:"size"`
	Project string `json:"project"`
}

// Storage - 스토리지 인터페이스
//
// Spring의 Repository 인터페이스와 동일한 패턴:
//   type UserRepository interface { ... }
//
// 이 인터페이스 덕분에 나중에 MinIO로 전환할 때
// main.go에서 구현체만 교체하면 나머지 코드는 변경 없음.
type Storage interface {
	// Save - 파일 저장 후 메타데이터 반환
	Save(project, originalName string, data io.Reader) (FileInfo, error)

	// Get - 파일 스트림 반환
	Get(project, filename string) (io.ReadCloser, error)

	// Delete - 파일 삭제
	Delete(project, filename string) error

	// List - 프로젝트 내 파일 목록 반환
	List(project string) ([]FileInfo, error)
}

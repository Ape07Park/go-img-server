package storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// LocalStorage - 로컬 파일시스템 기반 Storage 구현체
//
// Spring 비유:
//   interface Storage         → UserRepository (인터페이스)
//   struct LocalStorage       → UserRepositoryImpl (구현체)
//   func NewLocalStorage(...) → @Bean / @Component 생성자
type LocalStorage struct {
	baseDir string // 파일 저장 루트 경로 (예: ./uploads)
	baseURL string // 이미지 URL 생성용 도메인 (예: https://img.mycompany.com)
}

// NewLocalStorage - LocalStorage 생성자
// Spring의 @Bean 메서드 또는 생성자 주입과 동일
func NewLocalStorage(baseDir, baseURL string) *LocalStorage {
	return &LocalStorage{
		baseDir: baseDir,
		baseURL: baseURL,
	}
}

// Save - 파일을 로컬 디스크에 저장
// 저장 경로: {baseDir}/{project}/{타임스탬프}.{확장자}
func (s *LocalStorage) Save(project, originalName string, data io.Reader) (FileInfo, error) {
	// 보안: 프로젝트명에서 경로 탐색 문자(../) 제거
	project = sanitize(project)

	// 디렉터리 생성 (없으면 자동 생성)
	dir := filepath.Join(s.baseDir, project)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return FileInfo{}, fmt.Errorf("디렉터리 생성 실패: %w", err)
	}

	// 파일명: 타임스탬프 + 원본 확장자 (예: 1700000000000000000.jpg)
	ext := filepath.Ext(originalName)
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	path := filepath.Join(dir, filename)

	// 파일 생성 및 데이터 쓰기
	f, err := os.Create(path)
	if err != nil {
		return FileInfo{}, fmt.Errorf("파일 생성 실패: %w", err)
	}
	defer f.Close()

	size, err := io.Copy(f, data)
	if err != nil {
		return FileInfo{}, fmt.Errorf("파일 쓰기 실패: %w", err)
	}

	return FileInfo{
		Name:    filename,
		URL:     fmt.Sprintf("%s/i/%s/%s", s.baseURL, project, filename),
		Size:    size,
		Project: project,
	}, nil
}

// Get - 파일 스트림 반환 (호출자가 Close() 책임)
func (s *LocalStorage) Get(project, filename string) (io.ReadCloser, error) {
	project = sanitize(project)
	filename = sanitize(filename)
	path := filepath.Join(s.baseDir, project, filename)
	return os.Open(path)
}

// Delete - 파일 삭제
func (s *LocalStorage) Delete(project, filename string) error {
	project = sanitize(project)
	filename = sanitize(filename)
	path := filepath.Join(s.baseDir, project, filename)
	return os.Remove(path)
}

// List - 프로젝트 폴더 내 파일 목록 반환
func (s *LocalStorage) List(project string) ([]FileInfo, error) {
	project = sanitize(project)
	dir := filepath.Join(s.baseDir, project)

	entries, err := os.ReadDir(dir)
	if os.IsNotExist(err) {
		// 프로젝트 폴더가 아직 없으면 빈 목록 반환 (에러 아님)
		return []FileInfo{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("목록 조회 실패: %w", err)
	}

	var files []FileInfo
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			continue
		}
		files = append(files, FileInfo{
			Name:    entry.Name(),
			URL:     fmt.Sprintf("%s/i/%s/%s", s.baseURL, project, entry.Name()),
			Size:    info.Size(),
			Project: project,
		})
	}
	return files, nil
}

// sanitize - 경로 탐색 공격(path traversal) 방지
// 예: "../../etc/passwd" → "etcpasswd"
func sanitize(name string) string {
	name = filepath.Clean(name)
	name = strings.ReplaceAll(name, "..", "")
	name = strings.ReplaceAll(name, "/", "")
	name = strings.ReplaceAll(name, "\\", "")
	return name
}

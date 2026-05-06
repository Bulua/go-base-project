package fileservice

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	filemodel "gobaseproject/server/internal/model/file"
)

const (
	defaultPageSize = 10
	maxPageSize     = 200
	DefaultMaxBytes = 10 << 20 // 10 MB
)

type repository interface {
	Insert(ctx context.Context, f filemodel.FileRecord) (uint64, error)
	Count(ctx context.Context, q filemodel.FileListQuery) (int64, error)
	List(ctx context.Context, q filemodel.FileListQuery) ([]filemodel.FileRecord, error)
	GetByID(ctx context.Context, id uint64) (*filemodel.FileRecord, error)
	Delete(ctx context.Context, id uint64) error
}

type Service struct {
	repo      repository
	uploadDir string
	maxBytes  int64
}

func NewService(repo repository, uploadDir string, maxBytes int64) *Service {
	if maxBytes <= 0 {
		maxBytes = DefaultMaxBytes
	}
	if uploadDir == "" {
		uploadDir = "./uploads"
	}
	return &Service{repo: repo, uploadDir: uploadDir, maxBytes: maxBytes}
}

func (s *Service) Upload(ctx context.Context, fh *multipart.FileHeader, uploaderID *uint64) (*filemodel.FileRecord, error) {
	if fh.Size > s.maxBytes {
		return nil, filemodel.ErrFileTooLarge
	}

	src, err := fh.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	ext := strings.ToLower(filepath.Ext(fh.Filename))
	key := fmt.Sprintf("%s/%s%s", time.Now().Format("2006/01/02"), randomHex(), ext)

	fullPath := filepath.Join(s.uploadDir, filepath.FromSlash(key))
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return nil, err
	}

	dst, err := os.Create(fullPath)
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		os.Remove(fullPath)
		return nil, err
	}

	mimeType := fh.Header.Get("Content-Type")
	if mimeType == "" || mimeType == "application/octet-stream" {
		if t := mime.TypeByExtension(ext); t != "" {
			mimeType = t
		}
	}
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	rec := filemodel.FileRecord{
		OriginalName: fh.Filename,
		StorageKey:   key,
		FileSize:     fh.Size,
		MimeType:     mimeType,
		UploaderID:   uploaderID,
		CreatedAt:    time.Now(),
	}

	id, err := s.repo.Insert(ctx, rec)
	if err != nil {
		os.Remove(fullPath)
		return nil, err
	}
	rec.ID = id
	return &rec, nil
}

func (s *Service) List(ctx context.Context, q filemodel.FileListQuery) (*filemodel.FileListResult, error) {
	if q.Page < 1 {
		q.Page = 1
	}
	if q.PageSize < 1 || q.PageSize > maxPageSize {
		q.PageSize = defaultPageSize
	}
	total, err := s.repo.Count(ctx, q)
	if err != nil {
		return nil, err
	}
	items, err := s.repo.List(ctx, q)
	if err != nil {
		return nil, err
	}
	if items == nil {
		items = []filemodel.FileRecord{}
	}
	return &filemodel.FileListResult{Total: total, Items: items, Page: q.Page, PageSize: q.PageSize}, nil
}

func (s *Service) GetByID(ctx context.Context, id uint64) (*filemodel.FileRecord, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) Delete(ctx context.Context, id uint64) error {
	rec, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if rec == nil {
		return filemodel.ErrFileNotFound
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}
	os.Remove(s.StoragePath(rec.StorageKey)) //nolint:errcheck
	return nil
}

// StoragePath resolves a storage key to the absolute disk path.
func (s *Service) StoragePath(key string) string {
	return filepath.Join(s.uploadDir, filepath.FromSlash(key))
}

// UploadDir returns the configured upload directory.
func (s *Service) UploadDir() string { return s.uploadDir }

func randomHex() string {
	var b [16]byte
	rand.Read(b[:]) //nolint:errcheck
	return hex.EncodeToString(b[:])
}

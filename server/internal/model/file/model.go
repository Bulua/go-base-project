package filemodel

import (
	"errors"
	"time"
)

type FileRecord struct {
	ID           uint64    `json:"id"`
	OriginalName string    `json:"original_name"`
	StorageKey   string    `json:"storage_key"`
	FileSize     int64     `json:"file_size"`
	MimeType     string    `json:"mime_type"`
	UploaderID   *uint64   `json:"uploader_id,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

type FileListQuery struct {
	Page         int
	PageSize     int
	Keyword      string
	MimeCategory string // image | video | audio | text | pdf | other
	StartDate    string // YYYY-MM-DD
	EndDate      string // YYYY-MM-DD
}

type FileListResult struct {
	Total    int64        `json:"total"`
	Items    []FileRecord `json:"items"`
	Page     int          `json:"page"`
	PageSize int          `json:"page_size"`
}

var (
	ErrFileNotFound = errors.New("文件不存在")
	ErrFileTooLarge = errors.New("文件大小超出限制")
)

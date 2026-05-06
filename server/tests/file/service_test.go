package file_test

import (
	"context"
	"errors"
	"mime/multipart"
	"testing"

	filemodel "gobaseproject/server/internal/model/file"
	fileservice "gobaseproject/server/internal/service/file"
)

type stubRepository struct{}

func (stubRepository) Insert(context.Context, filemodel.FileRecord) (uint64, error) {
	return 0, errors.New("insert should not be called")
}

func (stubRepository) Count(context.Context, filemodel.FileListQuery) (int64, error) {
	return 0, nil
}

func (stubRepository) List(context.Context, filemodel.FileListQuery) ([]filemodel.FileRecord, error) {
	return nil, nil
}

func (stubRepository) GetByID(context.Context, uint64) (*filemodel.FileRecord, error) {
	return nil, nil
}

func (stubRepository) Delete(context.Context, uint64) error {
	return nil
}

func TestUploadUsesTenMegabytesAsDefaultLimit(t *testing.T) {
	service := fileservice.NewService(stubRepository{}, t.TempDir(), 0)
	header := &multipart.FileHeader{
		Filename: "too-large.bin",
		Size:     (10 << 20) + 1,
	}

	_, err := service.Upload(context.Background(), header, nil)
	if !errors.Is(err, filemodel.ErrFileTooLarge) {
		t.Fatalf("expected ErrFileTooLarge for file larger than 10MB, got %v", err)
	}
}

package auditservice

import (
	"context"

	auditmodel "gobaseproject/server/internal/model/audit"
)

type repository interface {
	CountLoginLogs(ctx context.Context, q auditmodel.LoginLogQuery) (int64, error)
	ListLoginLogs(ctx context.Context, q auditmodel.LoginLogQuery) ([]auditmodel.LoginLogRecord, error)
	CleanupLoginLogs(ctx context.Context, days int) (int64, error)

	CountOperationLogs(ctx context.Context, q auditmodel.OperationLogQuery) (int64, error)
	ListOperationLogs(ctx context.Context, q auditmodel.OperationLogQuery) ([]auditmodel.OperationLogRecord, error)
	GetOperationLogByID(ctx context.Context, id uint64) (*auditmodel.OperationLogRecord, error)
	CleanupOperationLogs(ctx context.Context, days int) (int64, error)
}

type Service struct {
	repo repository
}

func NewService(repo repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) ListLoginLogs(ctx context.Context, q auditmodel.LoginLogQuery) (*auditmodel.LoginLogResult, error) {
	if q.Page < 1 {
		q.Page = 1
	}
	if q.PageSize < 1 || q.PageSize > 200 {
		q.PageSize = 20
	}
	total, err := s.repo.CountLoginLogs(ctx, q)
	if err != nil {
		return nil, err
	}
	items, err := s.repo.ListLoginLogs(ctx, q)
	if err != nil {
		return nil, err
	}
	if items == nil {
		items = []auditmodel.LoginLogRecord{}
	}
	return &auditmodel.LoginLogResult{
		Total:    total,
		Items:    items,
		Page:     q.Page,
		PageSize: q.PageSize,
	}, nil
}

func (s *Service) CleanupLoginLogs(ctx context.Context, days int) (*auditmodel.CleanupResult, error) {
	if days < 1 {
		days = 90
	}
	deleted, err := s.repo.CleanupLoginLogs(ctx, days)
	if err != nil {
		return nil, err
	}
	return &auditmodel.CleanupResult{Deleted: deleted}, nil
}

func (s *Service) ListOperationLogs(ctx context.Context, q auditmodel.OperationLogQuery) (*auditmodel.OperationLogResult, error) {
	if q.Page < 1 {
		q.Page = 1
	}
	if q.PageSize < 1 || q.PageSize > 200 {
		q.PageSize = 20
	}
	total, err := s.repo.CountOperationLogs(ctx, q)
	if err != nil {
		return nil, err
	}
	items, err := s.repo.ListOperationLogs(ctx, q)
	if err != nil {
		return nil, err
	}
	if items == nil {
		items = []auditmodel.OperationLogRecord{}
	}
	return &auditmodel.OperationLogResult{
		Total:    total,
		Items:    items,
		Page:     q.Page,
		PageSize: q.PageSize,
	}, nil
}

func (s *Service) GetOperationLogByID(ctx context.Context, id uint64) (*auditmodel.OperationLogRecord, error) {
	return s.repo.GetOperationLogByID(ctx, id)
}

func (s *Service) CleanupOperationLogs(ctx context.Context, days int) (*auditmodel.CleanupResult, error) {
	if days < 1 {
		days = 90
	}
	deleted, err := s.repo.CleanupOperationLogs(ctx, days)
	if err != nil {
		return nil, err
	}
	return &auditmodel.CleanupResult{Deleted: deleted}, nil
}

package dictservice

import (
	"context"
	"strings"

	dictmodel "gobaseproject/server/internal/model/dict"
)

type Repository interface {
	Count(ctx context.Context, q dictmodel.ListQuery) (int64, error)
	List(ctx context.Context, q dictmodel.ListQuery) ([]dictmodel.Dictionary, error)
	GetByID(ctx context.Context, id uint64) (*dictmodel.Dictionary, error)
	GetByCode(ctx context.Context, code string) (*dictmodel.Dictionary, error)
	CountByCode(ctx context.Context, code string, excludeID uint64) (int64, error)
	Create(ctx context.Context, p dictmodel.SaveDictPayload) (uint64, error)
	Update(ctx context.Context, id uint64, p dictmodel.SaveDictPayload) error
	SoftDelete(ctx context.Context, id uint64) error

	ListItems(ctx context.Context, dictID uint64) ([]dictmodel.DictItem, error)
	GetItemByID(ctx context.Context, id uint64) (*dictmodel.DictItem, error)
	CountItemByValue(ctx context.Context, dictID uint64, value string, excludeID uint64) (int64, error)
	CreateItem(ctx context.Context, dictID uint64, p dictmodel.SaveItemPayload) (uint64, error)
	UpdateItem(ctx context.Context, id uint64, p dictmodel.SaveItemPayload) error
	SoftDeleteItem(ctx context.Context, id uint64) error
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) List(ctx context.Context, q dictmodel.ListQuery) (dictmodel.ListResult, error) {
	total, err := s.repo.Count(ctx, q)
	if err != nil {
		return dictmodel.ListResult{}, err
	}
	items, err := s.repo.List(ctx, q)
	if err != nil {
		return dictmodel.ListResult{}, err
	}
	return dictmodel.ListResult{
		Total:    total,
		Items:    items,
		Page:     q.Page,
		PageSize: q.PageSize,
	}, nil
}

func (s *Service) Create(ctx context.Context, p dictmodel.SaveDictPayload) (*dictmodel.Dictionary, error) {
	if strings.TrimSpace(p.DictName) == "" {
		return nil, dictmodel.ErrDictNameEmpty
	}
	if strings.TrimSpace(p.DictCode) == "" {
		return nil, dictmodel.ErrDictCodeEmpty
	}
	if p.DictStatus != dictmodel.StatusActive && p.DictStatus != dictmodel.StatusDisabled {
		p.DictStatus = dictmodel.StatusActive
	}
	if n, err := s.repo.CountByCode(ctx, p.DictCode, 0); err != nil {
		return nil, err
	} else if n > 0 {
		return nil, dictmodel.ErrDictCodeTaken
	}
	id, err := s.repo.Create(ctx, p)
	if err != nil {
		return nil, err
	}
	return s.repo.GetByID(ctx, id)
}

func (s *Service) Update(ctx context.Context, id uint64, p dictmodel.SaveDictPayload) (*dictmodel.Dictionary, error) {
	if strings.TrimSpace(p.DictName) == "" {
		return nil, dictmodel.ErrDictNameEmpty
	}
	if strings.TrimSpace(p.DictCode) == "" {
		return nil, dictmodel.ErrDictCodeEmpty
	}
	if p.DictStatus != dictmodel.StatusActive && p.DictStatus != dictmodel.StatusDisabled {
		p.DictStatus = dictmodel.StatusActive
	}
	if n, err := s.repo.CountByCode(ctx, p.DictCode, id); err != nil {
		return nil, err
	} else if n > 0 {
		return nil, dictmodel.ErrDictCodeTaken
	}
	if err := s.repo.Update(ctx, id, p); err != nil {
		return nil, err
	}
	return s.repo.GetByID(ctx, id)
}

func (s *Service) Delete(ctx context.Context, id uint64) error {
	return s.repo.SoftDelete(ctx, id)
}

func (s *Service) ListItemsByCode(ctx context.Context, code string) ([]dictmodel.DictItem, error) {
	d, err := s.repo.GetByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	return s.repo.ListItems(ctx, d.ID)
}

func (s *Service) ListItems(ctx context.Context, dictID uint64) ([]dictmodel.DictItem, error) {
	if _, err := s.repo.GetByID(ctx, dictID); err != nil {
		return nil, err
	}
	return s.repo.ListItems(ctx, dictID)
}

func (s *Service) CreateItem(ctx context.Context, dictID uint64, p dictmodel.SaveItemPayload) (*dictmodel.DictItem, error) {
	if _, err := s.repo.GetByID(ctx, dictID); err != nil {
		return nil, err
	}
	if strings.TrimSpace(p.ItemLabel) == "" {
		return nil, dictmodel.ErrItemLabelEmpty
	}
	if strings.TrimSpace(p.ItemValue) == "" {
		return nil, dictmodel.ErrItemValueEmpty
	}
	if p.ItemStatus != dictmodel.StatusActive && p.ItemStatus != dictmodel.StatusDisabled {
		p.ItemStatus = dictmodel.StatusActive
	}
	if n, err := s.repo.CountItemByValue(ctx, dictID, p.ItemValue, 0); err != nil {
		return nil, err
	} else if n > 0 {
		return nil, dictmodel.ErrItemValueTaken
	}
	id, err := s.repo.CreateItem(ctx, dictID, p)
	if err != nil {
		return nil, err
	}
	return s.repo.GetItemByID(ctx, id)
}

func (s *Service) UpdateItem(ctx context.Context, id uint64, p dictmodel.SaveItemPayload) (*dictmodel.DictItem, error) {
	existing, err := s.repo.GetItemByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(p.ItemLabel) == "" {
		return nil, dictmodel.ErrItemLabelEmpty
	}
	if strings.TrimSpace(p.ItemValue) == "" {
		return nil, dictmodel.ErrItemValueEmpty
	}
	if p.ItemStatus != dictmodel.StatusActive && p.ItemStatus != dictmodel.StatusDisabled {
		p.ItemStatus = dictmodel.StatusActive
	}
	if n, err := s.repo.CountItemByValue(ctx, existing.DictID, p.ItemValue, id); err != nil {
		return nil, err
	} else if n > 0 {
		return nil, dictmodel.ErrItemValueTaken
	}
	if err := s.repo.UpdateItem(ctx, id, p); err != nil {
		return nil, err
	}
	return s.repo.GetItemByID(ctx, id)
}

func (s *Service) DeleteItem(ctx context.Context, id uint64) error {
	return s.repo.SoftDeleteItem(ctx, id)
}

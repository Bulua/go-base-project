package menuservice

import (
	"context"
	"strings"

	menumodel "gobaseproject/server/internal/model/menu"
)

const (
	defaultPage     = 1
	defaultPageSize = 20
	maxPageSize     = 200
)

const superAdminRoleID uint64 = 1

type Repository interface {
	Count(ctx context.Context, q menumodel.ListQuery) (int64, error)
	List(ctx context.Context, q menumodel.ListQuery) ([]menumodel.Menu, error)
	ListAll(ctx context.Context) ([]menumodel.Menu, error)
	GetByID(ctx context.Context, id uint64) (*menumodel.Menu, error)
	CountChildren(ctx context.Context, parentID uint64) (int64, error)
	Create(ctx context.Context, req menumodel.SaveRequest) (uint64, error)
	Update(ctx context.Context, id uint64, req menumodel.SaveRequest) error
	SoftDelete(ctx context.Context, id uint64) error
	AssignToRole(ctx context.Context, roleID, menuID uint64) error
	ListParams(ctx context.Context, menuID uint64) ([]menumodel.RouteParam, error)
	CreateParam(ctx context.Context, menuID uint64, req menumodel.CreateParamRequest) (uint64, error)
	DeleteParam(ctx context.Context, paramID uint64) error
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) List(ctx context.Context, q menumodel.ListQuery) (menumodel.ListResult, error) {
	q = normalizeQuery(q)
	total, err := s.repo.Count(ctx, q)
	if err != nil {
		return menumodel.ListResult{}, err
	}
	if total == 0 {
		return menumodel.ListResult{Items: []menumodel.Menu{}, Page: q.Page, Size: q.PageSize}, nil
	}
	items, err := s.repo.List(ctx, q)
	if err != nil {
		return menumodel.ListResult{}, err
	}
	return menumodel.ListResult{Total: total, Items: items, Page: q.Page, Size: q.PageSize}, nil
}

func (s *Service) Tree(ctx context.Context) ([]menumodel.Menu, error) {
	menus, err := s.repo.ListAll(ctx)
	if err != nil {
		return nil, err
	}
	return buildTree(menus), nil
}

func (s *Service) Get(ctx context.Context, id uint64) (*menumodel.Menu, error) {
	menu, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	params, err := s.repo.ListParams(ctx, id)
	if err != nil {
		return nil, err
	}
	menu.Params = params
	return menu, nil
}

func (s *Service) Create(ctx context.Context, req menumodel.SaveRequest) (*menumodel.Menu, error) {
	if err := validate(req); err != nil {
		return nil, err
	}
	if req.ParentID > 0 {
		if _, err := s.repo.GetByID(ctx, req.ParentID); err != nil {
			return nil, menumodel.ErrParentNotFound
		}
	}
	id, err := s.repo.Create(ctx, req)
	if err != nil {
		return nil, err
	}
	// 自动关联超级管理员角色，确保新菜单立即可见
	_ = s.repo.AssignToRole(ctx, superAdminRoleID, id)
	return s.Get(ctx, id)
}

func (s *Service) Update(ctx context.Context, id uint64, req menumodel.SaveRequest) (*menumodel.Menu, error) {
	if err := validate(req); err != nil {
		return nil, err
	}
	if req.ParentID > 0 {
		if req.ParentID == id {
			return nil, menumodel.ErrParentNotFound
		}
		if _, err := s.repo.GetByID(ctx, req.ParentID); err != nil {
			return nil, menumodel.ErrParentNotFound
		}
	}
	if err := s.repo.Update(ctx, id, req); err != nil {
		return nil, err
	}
	return s.Get(ctx, id)
}

func (s *Service) Delete(ctx context.Context, id uint64) error {
	if _, err := s.repo.GetByID(ctx, id); err != nil {
		return err
	}
	count, err := s.repo.CountChildren(ctx, id)
	if err != nil {
		return err
	}
	if count > 0 {
		return menumodel.ErrHasChildren
	}
	return s.repo.SoftDelete(ctx, id)
}

func (s *Service) ListParams(ctx context.Context, menuID uint64) ([]menumodel.RouteParam, error) {
	if _, err := s.repo.GetByID(ctx, menuID); err != nil {
		return nil, err
	}
	return s.repo.ListParams(ctx, menuID)
}

func (s *Service) CreateParam(ctx context.Context, menuID uint64, req menumodel.CreateParamRequest) (*menumodel.RouteParam, error) {
	req.ParamKey = strings.TrimSpace(req.ParamKey)
	req.ParamMode = strings.TrimSpace(req.ParamMode)
	if req.ParamKey == "" {
		return nil, menumodel.ErrParamKeyEmpty
	}
	if req.ParamMode != "query" && req.ParamMode != "path" {
		return nil, menumodel.ErrParamModeBad
	}
	if _, err := s.repo.GetByID(ctx, menuID); err != nil {
		return nil, err
	}
	id, err := s.repo.CreateParam(ctx, menuID, req)
	if err != nil {
		return nil, err
	}
	params, err := s.repo.ListParams(ctx, menuID)
	if err != nil {
		return nil, err
	}
	for i := range params {
		if params[i].ID == id {
			return &params[i], nil
		}
	}
	return nil, menumodel.ErrParamNotFound
}

func (s *Service) DeleteParam(ctx context.Context, paramID uint64) error {
	return s.repo.DeleteParam(ctx, paramID)
}

// ── helpers ──────────────────────────────────────────────────────────────────

func validate(req menumodel.SaveRequest) error {
	if strings.TrimSpace(req.MenuTitle) == "" {
		return menumodel.ErrMenuTitleEmpty
	}
	if req.MenuType < menumodel.TypeDirectory || req.MenuType > menumodel.TypeExternal {
		return menumodel.ErrInvalidType
	}
	if req.MenuStatus != menumodel.StatusActive && req.MenuStatus != menumodel.StatusDisabled {
		req.MenuStatus = menumodel.StatusActive
	}
	return nil
}

func normalizeQuery(q menumodel.ListQuery) menumodel.ListQuery {
	if q.Page <= 0 {
		q.Page = defaultPage
	}
	if q.PageSize <= 0 {
		q.PageSize = defaultPageSize
	}
	if q.PageSize > maxPageSize {
		q.PageSize = maxPageSize
	}
	q.Keyword = strings.TrimSpace(q.Keyword)
	return q
}

func buildTree(menus []menumodel.Menu) []menumodel.Menu {
	byParent := map[uint64][]menumodel.Menu{}
	for _, m := range menus {
		byParent[m.ParentID] = append(byParent[m.ParentID], m)
	}
	var attach func(parentID uint64) []menumodel.Menu
	attach = func(parentID uint64) []menumodel.Menu {
		items := byParent[parentID]
		for i := range items {
			items[i].Children = attach(items[i].ID)
		}
		return items
	}
	return attach(0)
}

package dictmodel

import (
	"errors"
	"time"
)

type Dictionary struct {
	ID         uint64     `json:"id"`
	DictName   string     `json:"dict_name"`
	DictCode   string     `json:"dict_code"`
	DictStatus int        `json:"dict_status"`
	ParentID   uint64     `json:"parent_id"`
	Remark     *string    `json:"remark"`
	ItemCount  int64      `json:"item_count"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

type DictItem struct {
	ID         uint64     `json:"id"`
	DictID     uint64     `json:"dict_id"`
	ItemLabel  string     `json:"item_label"`
	ItemValue  string     `json:"item_value"`
	ItemExtra  *string    `json:"item_extra"`
	ItemStatus int        `json:"item_status"`
	SortNo     int        `json:"sort_no"`
	ParentID   uint64     `json:"parent_id"`
	TreeLevel  int        `json:"tree_level"`
	TreePath   *string    `json:"tree_path"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

type ListQuery struct {
	Page       int
	PageSize   int
	Keyword    string
	DictStatus int
}

type ListResult struct {
	Total    int64        `json:"total"`
	Items    []Dictionary `json:"items"`
	Page     int          `json:"page"`
	PageSize int          `json:"page_size"`
}

type SaveDictPayload struct {
	DictName   string
	DictCode   string
	DictStatus int
	ParentID   uint64
	Remark     *string
}

type SaveItemPayload struct {
	ItemLabel  string
	ItemValue  string
	ItemExtra  *string
	ItemStatus int
	SortNo     int
	ParentID   uint64
}

const (
	StatusActive   = 1
	StatusDisabled = 2
)

var (
	ErrDictNotFound   = errors.New("dict: not found")
	ErrDictCodeTaken  = errors.New("dict: code already exists")
	ErrDictCodeEmpty  = errors.New("dict: code is empty")
	ErrDictNameEmpty  = errors.New("dict: name is empty")
	ErrItemNotFound   = errors.New("dict item: not found")
	ErrItemValueTaken = errors.New("dict item: value already exists in this dict")
	ErrItemLabelEmpty = errors.New("dict item: label is empty")
	ErrItemValueEmpty = errors.New("dict item: value is empty")
)

package apperror

import (
	"errors"
	"net/http"

	authmodel "gobaseproject/server/internal/model/auth"
	dictmodel "gobaseproject/server/internal/model/dict"
	menumodel "gobaseproject/server/internal/model/menu"
	rolemodel "gobaseproject/server/internal/model/role"
	usermodel "gobaseproject/server/internal/model/user"
	"gobaseproject/server/pkg/response"
)

type Definition struct {
	Code    int
	Status  int
	Message string
}

var (
	InvalidRequestBody  = Definition{Code: 400001, Status: http.StatusBadRequest, Message: "请求参数格式不正确"}
	InvalidParams       = Definition{Code: 400002, Status: http.StatusBadRequest, Message: "请求参数不正确"}
	InvalidUserID       = Definition{Code: 400201, Status: http.StatusBadRequest, Message: "用户ID不正确"}
	InvalidRoleID       = Definition{Code: 400301, Status: http.StatusBadRequest, Message: "角色ID不正确"}
	InvalidMenuID       = Definition{Code: 400401, Status: http.StatusBadRequest, Message: "菜单ID不正确"}
	NotFound            = Definition{Code: 404001, Status: http.StatusNotFound, Message: "资源不存在"}
	MissingAuthToken    = Definition{Code: 401001, Status: http.StatusUnauthorized, Message: "缺少登录凭证"}
	InvalidAuthToken    = Definition{Code: 401002, Status: http.StatusUnauthorized, Message: "登录状态已失效，请重新登录"}
	MissingRefreshToken = Definition{Code: 401003, Status: http.StatusUnauthorized, Message: "缺少刷新凭证"}
	Forbidden           = Definition{Code: 403001, Status: http.StatusForbidden, Message: "无权执行该操作"}
	MethodNotAllowed    = Definition{Code: 405001, Status: http.StatusMethodNotAllowed, Message: "请求方法不允许"}
	Internal            = Definition{Code: 500001, Status: http.StatusInternalServerError, Message: "系统繁忙，请稍后重试"}
)

var registry = []struct {
	err error
	def Definition
}{
	{authmodel.ErrInvalidCredentials, Definition{Code: 401101, Status: http.StatusUnauthorized, Message: "账号或密码错误"}},
	{authmodel.ErrUserDisabled, Definition{Code: 403101, Status: http.StatusForbidden, Message: "用户已被禁用"}},
	{authmodel.ErrInvalidToken, InvalidAuthToken},
	{authmodel.ErrTokenBlocked, InvalidAuthToken},
	{usermodel.ErrUserNotFound, Definition{Code: 404201, Status: http.StatusNotFound, Message: "用户不存在"}},
	{usermodel.ErrLoginNameTaken, Definition{Code: 409201, Status: http.StatusConflict, Message: "登录账号已存在"}},
	{usermodel.ErrLoginNameInvalid, Definition{Code: 400202, Status: http.StatusBadRequest, Message: "账号必须以字母开头，长度 3-32 位，仅支持字母、数字、下划线、点和横线"}},
	{usermodel.ErrPasswordWeak, Definition{Code: 400203, Status: http.StatusBadRequest, Message: "密码至少 6 位"}},
	{usermodel.ErrInvalidStatus, Definition{Code: 400204, Status: http.StatusBadRequest, Message: "用户状态不正确"}},
	{usermodel.ErrAdminProtected, Definition{Code: 403201, Status: http.StatusForbidden, Message: "内置 admin 用户不允许执行该操作"}},
	{usermodel.ErrRoleAssignmentBad, Definition{Code: 400205, Status: http.StatusBadRequest, Message: "角色分配不正确"}},
	{rolemodel.ErrRoleNotFound, Definition{Code: 404301, Status: http.StatusNotFound, Message: "角色不存在"}},
	{rolemodel.ErrRoleCodeTaken, Definition{Code: 409301, Status: http.StatusConflict, Message: "角色编码已存在"}},
	{rolemodel.ErrRoleCodeInvalid, Definition{Code: 400302, Status: http.StatusBadRequest, Message: "角色编码必须以小写字母开头，长度 3-64 位，仅支持小写字母、数字和下划线"}},
	{rolemodel.ErrInvalidStatus, Definition{Code: 400303, Status: http.StatusBadRequest, Message: "角色状态不正确"}},
	{rolemodel.ErrBuiltinProtect, Definition{Code: 403301, Status: http.StatusForbidden, Message: "内置角色不允许执行该操作"}},
	{rolemodel.ErrRoleHasUsers, Definition{Code: 409302, Status: http.StatusConflict, Message: "角色下仍有关联用户，请先解除关联"}},
	{rolemodel.ErrRoleHasChildren, Definition{Code: 409303, Status: http.StatusConflict, Message: "角色下仍有子角色，请先删除子角色"}},
	{rolemodel.ErrParentLoop, Definition{Code: 400304, Status: http.StatusBadRequest, Message: "上级角色不能形成循环关系"}},
	{menumodel.ErrMenuNotFound, Definition{Code: 404401, Status: http.StatusNotFound, Message: "菜单不存在"}},
	{menumodel.ErrMenuTitleEmpty, Definition{Code: 400402, Status: http.StatusBadRequest, Message: "菜单名称不能为空"}},
	{menumodel.ErrInvalidType, Definition{Code: 400403, Status: http.StatusBadRequest, Message: "菜单类型不正确"}},
	{menumodel.ErrInvalidStatus, Definition{Code: 400404, Status: http.StatusBadRequest, Message: "菜单状态不正确"}},
	{menumodel.ErrParentNotFound, Definition{Code: 404402, Status: http.StatusNotFound, Message: "上级菜单不存在"}},
	{menumodel.ErrHasChildren, Definition{Code: 409401, Status: http.StatusConflict, Message: "菜单下仍有子菜单，请先删除子菜单"}},
	{menumodel.ErrParamNotFound, Definition{Code: 404403, Status: http.StatusNotFound, Message: "路由参数不存在"}},
	{menumodel.ErrParamKeyEmpty, Definition{Code: 400405, Status: http.StatusBadRequest, Message: "参数键名不能为空"}},
	{menumodel.ErrParamModeBad, Definition{Code: 400406, Status: http.StatusBadRequest, Message: "参数模式必须为 query 或 path"}},
	{menumodel.ErrActionNotFound, Definition{Code: 404404, Status: http.StatusNotFound, Message: "按钮权限不存在"}},
	{menumodel.ErrActionCodeEmpty, Definition{Code: 400407, Status: http.StatusBadRequest, Message: "按钮编码不能为空"}},
	{menumodel.ErrActionNameEmpty, Definition{Code: 400408, Status: http.StatusBadRequest, Message: "按钮名称不能为空"}},
	{menumodel.ErrActionCodeTaken, Definition{Code: 409402, Status: http.StatusConflict, Message: "该菜单下按钮编码已存在"}},
	{dictmodel.ErrDictNotFound, Definition{Code: 404501, Status: http.StatusNotFound, Message: "字典不存在"}},
	{dictmodel.ErrDictCodeTaken, Definition{Code: 409501, Status: http.StatusConflict, Message: "字典编码已存在"}},
	{dictmodel.ErrDictCodeEmpty, Definition{Code: 400501, Status: http.StatusBadRequest, Message: "字典编码不能为空"}},
	{dictmodel.ErrDictNameEmpty, Definition{Code: 400502, Status: http.StatusBadRequest, Message: "字典名称不能为空"}},
	{dictmodel.ErrItemNotFound, Definition{Code: 404502, Status: http.StatusNotFound, Message: "字典项不存在"}},
	{dictmodel.ErrItemValueTaken, Definition{Code: 409502, Status: http.StatusConflict, Message: "该字典下字典项值已存在"}},
	{dictmodel.ErrItemLabelEmpty, Definition{Code: 400503, Status: http.StatusBadRequest, Message: "字典项标签不能为空"}},
	{dictmodel.ErrItemValueEmpty, Definition{Code: 400504, Status: http.StatusBadRequest, Message: "字典项值不能为空"}},
}

func Write(w http.ResponseWriter, r *http.Request, err error) {
	WriteDefinition(w, r, FromError(err))
}

func WriteDefinition(w http.ResponseWriter, r *http.Request, def Definition) {
	response.WriteJSON(w, def.Status, response.Body{
		Code:    def.Code,
		Message: def.Message,
		TraceID: response.TraceID(r),
	})
}

func AllowMethod(w http.ResponseWriter, r *http.Request, method string) bool {
	if r.Method == method {
		return true
	}
	w.Header().Set("Allow", method)
	WriteDefinition(w, r, MethodNotAllowed)
	return false
}

func FromError(err error) Definition {
	for _, item := range registry {
		if errors.Is(err, item.err) {
			return item.def
		}
	}
	return Internal
}

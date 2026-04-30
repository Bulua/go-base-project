package apperror

import (
	"errors"
	"net/http"

	authmodel "gobaseproject/server/internal/model/auth"
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
	InvalidUserID       = Definition{Code: 400201, Status: http.StatusBadRequest, Message: "用户ID不正确"}
	InvalidRoleID       = Definition{Code: 400301, Status: http.StatusBadRequest, Message: "角色ID不正确"}
	MissingAuthToken    = Definition{Code: 401001, Status: http.StatusUnauthorized, Message: "缺少登录凭证"}
	InvalidAuthToken    = Definition{Code: 401002, Status: http.StatusUnauthorized, Message: "登录状态已失效，请重新登录"}
	MissingRefreshToken = Definition{Code: 401003, Status: http.StatusUnauthorized, Message: "缺少刷新凭证"}
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

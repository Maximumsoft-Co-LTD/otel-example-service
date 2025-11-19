package dto

type LoginRequest struct {
	Username   string `json:"username" validate:"required"`
	Password   string `json:"password" validate:"required"`
	RememberMe bool   `json:"remember_me" validate:"omitempty"`
}

type UserCreateRequest struct {
	Name     string              `json:"name" validate:"required"`
	Username string              `json:"username" validate:"required"`
	Password string              `json:"password" validate:"required"`
	Email    string              `json:"email" validate:"required,email"`
	Phone    string              `json:"phone" validate:"required,numeric"`
	Roles    []RoleCreateRequest `json:"roles" validate:"required"`
}

type UserUpdateRequest struct {
	Name     *string              `json:"name" validate:"required"`
	Username *string              `json:"username" validate:"required"`
	Password *string              `json:"password" validate:"required"`
	Email    *string              `json:"email" validate:"required,email"`
	Phone    *string              `json:"phone" validate:"required,numeric"`
	Roles    *[]RoleCreateRequest `json:"roles" validate:"required"`
}

type RoleCreateRequest struct {
	Name        string   `json:"name" validate:"required"`
	Permissions []string `json:"permissions" validate:"required"`
}

type RoleUpdateRequest struct {
	Name        *string   `json:"name" validate:"required"`
	Permissions *[]string `json:"permissions" validate:"required"`
}

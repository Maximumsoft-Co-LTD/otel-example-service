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

type MangaCheckRequest struct {
	Name string `json:"name"`
}

type MangaFindRequest struct {
	Keyword string `json:"keyword" validate:"required"`
}

type MangaCreateRequest struct {
	Name        string   `json:"name"`
	NameTH      string   `json:"name_th,omitempty"`
	NameEN      string   `json:"name_en,omitempty"`
	Pseudonym   string   `json:"pseudonym,omitempty"`
	Description string   `json:"description,omitempty"`
	Group       string   `json:"group"`
	Type        string   `json:"type"`
	Tags        []string `json:"tags,omitempty"`
	Categories  []string `json:"categories,omitempty"`
	CoverImage  string   `json:"cover_image_url,omitempty"`
}

type MangaUpdateRequest struct {
	Name        *string   `json:"name"`
	NameTH      *string   `json:"name_th,omitempty"`
	NameEN      *string   `json:"name_en,omitempty"`
	Pseudonym   *string   `json:"pseudonym,omitempty"`
	Description *string   `json:"description,omitempty"`
	Group       *string   `json:"group"`
	Type        *string   `json:"type"`
	Tags        *[]string `json:"tags,omitempty"`
	Categories  *[]string `json:"categories,omitempty"`
	CoverImage  *string   `json:"cover_image_url,omitempty"`
	Active      *bool     `json:"active,omitempty"`
}

type EpisodeCheckRequest struct {
	MangaID string `json:"manga_id" validate:"required"`
	Name    string `json:"name" validate:"required"`
}

type EpisodeFindRequest struct {
	Keyword string `json:"keyword" validate:"required"`
}

type EpisodeCreateRequest struct {
	MangaID     string `json:"manga_id"`
	OrderNo     int    `json:"order_no"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	CoverImage  string `json:"cover_image_url,omitempty"`
}

type EpisodeUpdateRequest struct {
	OrderNo     *int    `json:"order_no"`
	Name        *string `json:"name"`
	Description *string `json:"description,omitempty"`
	CoverImage  *string `json:"cover_image_url,omitempty"`
	Active      *bool   `json:"active,omitempty"`
}

type CategoryCheckRequest struct {
	Name string `json:"name" validate:"required"`
}

type CategoryFindRequest struct {
	Keyword string `json:"keyword" validate:"required"`
}

type CategoryCreateRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description,omitempty" validate:"omitempty"`
}

type CategoryUpdateRequest struct {
	Name        *string `json:"name" validate:"required"`
	Description *string `json:"description,omitempty" validate:"omitempty"`
	Active      *bool   `json:"active,omitempty" validate:"required"`
}

type TagCheckRequest struct {
	Name string `json:"name" validate:"required"`
}

type TagFindRequest struct {
	Keyword string `json:"keyword" validate:"required"`
}

type TagCreateRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description,omitempty" validate:"omitempty"`
}

type TagUpdateRequest struct {
	Name        *string `json:"name" validate:"required"`
	Description *string `json:"description,omitempty" validate:"required"`
	Active      *bool   `json:"active,omitempty" validate:"required"`
}

type ImageCheckRequest struct {
	Name string `json:"name" validate:"required"`
}

type ImageCreateRequest struct {
	Name   string `json:"name" validate:"required"`
	Active bool   `json:"active,omitempty" validate:"required"`
}

type ImageUpdateRequest struct {
	Name   *string `json:"name" validate:"required"`
	Active *bool   `json:"active,omitempty" validate:"required"`
}

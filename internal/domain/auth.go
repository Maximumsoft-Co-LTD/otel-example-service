package domain

type UserProfile struct {
	Name     string   `json:"name"`
	Username string   `json:"username"`
	Email    string   `json:"email"`
	Roles    []string `json:"roles,omitempty"`
}

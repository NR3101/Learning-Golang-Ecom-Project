package utils

type ContextKey string

const (
	UserID     ContextKey = "user_id"
	UserEmail  ContextKey = "user_email"
	UserRole   ContextKey = "user_role"
	GinContext ContextKey = "gin_context"
)

package context

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// contextKey for avoid key collisions
type contextKey string

type Context struct{}

// Set is used to store a new key/value pair exclusively for context
func (c *Context) Set(r *http.Request, key contextKey, value any) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), key, value))
}

// Get returns the value as any for the given ctx and key,
// if the value does not exist it returns nil
func (c *Context) Get(ctx context.Context, key contextKey) any {
	return ctx.Value(key)
}

// GetString returns the value as string
func (c *Context) GetString(ctx context.Context, key contextKey) string {
	value, ok := c.Get(ctx, key).(string)
	if !ok {
		return ""
	}

	return value
}

// GetTime returns the value as time
func (c *Context) GetTime(ctx context.Context, key contextKey) time.Time {
	value, ok := c.Get(ctx, key).(time.Time)
	if !ok {
		return time.Time{}
	}

	return value
}

// GetUUID returns the value as uuid
func (c *Context) GetUUID(ctx context.Context, key contextKey) uuid.UUID {
	var value uuid.UUID
	var ok bool

	value, ok = c.Get(ctx, key).(uuid.UUID)
	if !ok {
		return value
	}

	return value
}

// GetSqlx returns the value as *sqlx.DB
func (c *Context) GetSqlx(ctx context.Context, key contextKey) *sqlx.DB {
	value, ok := c.Get(ctx, key).(*sqlx.DB)
	if !ok {
		return nil
	}

	return value
}

// GetDB is shorthand to get database, the value as *sqlx.DB
func (c *Context) GetDB(ctx context.Context) *sqlx.DB {
	value := c.GetSqlx(ctx, "db")
	if value == nil {
		return nil
	}

	return value
}

// GetUserID is shorthand to get user_id, the value as uuid.
func (c *Context) GetUserID(ctx context.Context) uuid.UUID {
	return c.GetUUID(ctx, "user_id")
}

// GetUserUsername is shorthand to get user_username, the value as string.
func (c *Context) GetUserUsername(ctx context.Context) string {
	value := c.GetString(ctx, "user_username")
	if value == "" {
		return ""
	}

	return value
}

// GetUserRole is shorthand to get user_role, the value as string.
func (c *Context) GetUserRole(ctx context.Context) string {
	value := c.GetString(ctx, "user_role")
	if value == "" {
		return ""
	}

	return value
}

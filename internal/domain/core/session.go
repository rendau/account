package core

import (
	"context"
	"strconv"

	"github.com/rendau/account/internal/domain/entities"
	"github.com/rendau/dop/adapters/jwt"
)

const sessionContextKey = "user_session"
const sessionDur = int64(600) // seconds

type Session struct {
	r *St
}

func NewSession(r *St) *Session {
	return &Session{r: r}
}

func (c *Session) GetFromToken(token string) *entities.Session {
	result := &entities.Session{}

	defer func() {
		if result.Roles == nil {
			result.Roles = make([]string, 0)
		}

		if result.Perms == nil {
			result.Perms = make([]string, 0)
		}
	}()

	if token == "" {
		return result
	}

	err := jwt.ParsePayload(token, result)
	if err != nil {
		return &entities.Session{}
	}

	result.Id, _ = strconv.ParseInt(result.Sub, 10, 64)

	return result
}

func (c *Session) SetToContext(ctx context.Context, ses *entities.Session) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	return context.WithValue(ctx, sessionContextKey, ses)
}

func (c *Session) GetFromContext(ctx context.Context) *entities.Session {
	contextV := ctx.Value(sessionContextKey)
	if contextV == nil {
		return c.GetFromToken("")
	}

	switch ses := contextV.(type) {
	case *entities.Session:
		return ses
	default:
		c.r.lg.Errorw("wrong type of session in context", nil)
		return c.GetFromToken("")
	}
}

func (c *Session) CreateToken(ses *entities.Session, durSeconds int64) (string, error) {

	token, _ := c.r.jwts.Create(
		strconv.FormatInt(ses.Id, 10),
		durSeconds,
		map[string]any{
			"roles": ses.Roles,
			"perms": ses.Perms,
		},
	)

	return token, nil
}

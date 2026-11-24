package core

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/rendau/account/internal/domain/entities"
)

const sessionContextKey = "user_session"
const sessionDur = int64(600) // seconds

// tokenValidateTimeout — сколько ждать ответа jwts при проверке токена.
const tokenValidateTimeout = 3 * time.Second

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

	// Подпись и срок проверяет jwts: без этого account верил бы любому склеенному токену
	// (сам он полезную нагрузку только раскодирует), а на него ходят и мимо ruto.
	ctx, cancel := context.WithTimeout(context.Background(), tokenValidateTimeout)
	defer cancel()

	valid, claims, err := c.r.jwts.Validate(ctx, token)
	if err != nil {
		c.r.lg.Errorw("jwts.Validate", err)
		return result
	}
	if !valid {
		return result
	}

	var ses entities.Session
	if err = json.Unmarshal(claims, &ses); err != nil {
		return result
	}

	ses.Id, _ = strconv.ParseInt(ses.Sub, 10, 64)
	if ses.Id <= 0 {
		return result
	}

	result = &ses

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

func (c *Session) CreateToken(ses *entities.Session, dur time.Duration) (string, error) {

	token, _ := c.r.jwts.Create(
		context.Background(),
		strconv.FormatInt(ses.Id, 10),
		dur,
		map[string]any{
			"roles": ses.Roles,
			"perms": ses.Perms,
		},
	)

	return token, nil
}

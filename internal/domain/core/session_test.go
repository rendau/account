package core

import (
	"context"
	"encoding/base64"
	"errors"
	"testing"
	"time"

	"github.com/rendau/dop/adapters/logger/zap"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// fakeJwts отвечает как jwts: годны только токены из valid, остальные — нет.
type fakeJwts struct {
	valid map[string][]byte
	err   error
	calls int
}

func (f *fakeJwts) Create(context.Context, string, time.Duration, any) (string, error) {
	return "", nil
}

func (f *fakeJwts) Validate(_ context.Context, token string) (bool, []byte, error) {
	f.calls++
	if f.err != nil {
		return false, nil, f.err
	}
	claims, ok := f.valid[token]
	return ok, claims, nil
}

func newTestSession(jwts JwtI) *Session {
	return NewSession(&St{lg: zap.New("error", false), jwts: jwts})
}

// forgedToken — «токен», склеенный без ключа: правильный вид, подпись мусорная.
func forgedToken(payload string) string {
	enc := base64.RawURLEncoding.EncodeToString
	return enc([]byte(`{"alg":"RS256"}`)) + "." + enc([]byte(payload)) + ".x"
}

func TestSessionGetFromToken(t *testing.T) {
	signed := "signed.token.ok"
	jwts := &fakeJwts{valid: map[string][]byte{
		signed: []byte(`{"sub":"43","roles":["acc:super_admin"],"perms":["acc:*"]}`),
	}}
	s := newTestSession(jwts)

	t.Run("signed token is accepted", func(t *testing.T) {
		ses := s.GetFromToken(signed)
		assert.Equal(t, int64(43), ses.Id)
		assert.Equal(t, []string{"acc:super_admin"}, ses.Roles)
		assert.Equal(t, []string{"acc:*"}, ses.Perms)
	})

	t.Run("forged token is rejected", func(t *testing.T) {
		ses := s.GetFromToken(forgedToken(`{"sub":"43","roles":["acc:super_admin"],"perms":["acc:*"]}`))
		assert.Zero(t, ses.Id)
		assert.Empty(t, ses.Roles)
		assert.Empty(t, ses.Perms)
	})

	t.Run("empty token does not call jwts", func(t *testing.T) {
		before := jwts.calls
		ses := s.GetFromToken("")
		assert.Zero(t, ses.Id)
		assert.NotNil(t, ses.Perms)
		assert.Equal(t, before, jwts.calls)
	})
}

func TestSessionGetFromTokenJwtsDown(t *testing.T) {
	s := newTestSession(&fakeJwts{err: errors.New("connection refused")})

	ses := s.GetFromToken(forgedToken(`{"sub":"43","perms":["acc:*"]}`))
	require.NotNil(t, ses)
	assert.Zero(t, ses.Id, "jwts unavailable must not let the token through")
	assert.Empty(t, ses.Perms)
}

func TestSessionGetFromTokenBadSub(t *testing.T) {
	signed := "signed.token.nosub"
	s := newTestSession(&fakeJwts{valid: map[string][]byte{
		signed: []byte(`{"sub":"","perms":["acc:*"]}`),
	}})

	ses := s.GetFromToken(signed)
	assert.Zero(t, ses.Id)
	assert.Empty(t, ses.Perms)
}

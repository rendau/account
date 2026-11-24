package repo

import (
	"fmt"
	"time"

	"github.com/goccy/go-json"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/rendau/account/pkg/proto/jwts_v1"
)

type Repo struct {
	client jwts_v1.JwtClient
}

func New(uri string) (*Repo, error) {
	conn, err := grpc.Dial(uri, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("grpc.Dial: %w", err)
	}

	return &Repo{
		client: jwts_v1.NewJwtClient(conn),
	}, nil
}

func (r *Repo) Create(ctx context.Context, sub string, ttl time.Duration, payload any) (_ string, finalError error) {
	payloadJson, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("json.Marshal payload: %w", err)
	}

	repObj, err := r.client.Create(ctx, &jwts_v1.JwtCreateReq{
		Sub:        sub,
		ExpSeconds: int64(ttl.Seconds()),
		Payload:    payloadJson,
	})
	if err != nil {
		return "", fmt.Errorf("client.Create: %w", err)
	}

	return repObj.Token, nil
}

// Validate проверяет подпись и срок токена в jwts. Возвращает признак годности и claims (json).
func (r *Repo) Validate(ctx context.Context, token string) (bool, []byte, error) {
	repObj, err := r.client.Validate(ctx, &jwts_v1.JwtValidateReq{Token: token})
	if err != nil {
		return false, nil, fmt.Errorf("client.Validate: %w", err)
	}

	return repObj.Valid, repObj.Claims, nil
}

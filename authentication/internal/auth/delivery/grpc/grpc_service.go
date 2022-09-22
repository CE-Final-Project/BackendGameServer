package grpc

import (
	"context"
	"github.com/ce-final-project/backend_game_server/authentication/config"
	"github.com/ce-final-project/backend_game_server/authentication/internal/auth/service"
	authService "github.com/ce-final-project/backend_game_server/authentication/proto"
	"github.com/ce-final-project/backend_game_server/pkg/logger"
	"github.com/go-playground/validator"
	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type grpcService struct {
	log logger.Logger
	cfg *config.Config
	v   *validator.Validate
	as  *service.AuthService
}

func NewGrpcService(log logger.Logger, cfg *config.Config, v *validator.Validate, as *service.AuthService) *grpcService {
	return &grpcService{
		log: log,
		cfg: cfg,
		v:   v,
		as:  as,
	}
}

func (g *grpcService) Login(ctx context.Context, req *authService.LoginReq) (*authService.LoginRes, error) {
	//TODO implement me
	panic("implement me")
}

func (g *grpcService) Register(ctx context.Context, req *authService.RegisterReq) (*authService.RegisterRes, error) {
	//TODO implement me
	panic("implement me")
}

func (g *grpcService) VerifyToken(ctx context.Context, req *authService.VerifyTokenReq) (*authService.VerifyTokenRes, error) {
	result, err := g.validateToken(req.GetToken())
	if err != nil {
		g.log.Error(err)
		return &authService.VerifyTokenRes{
			Valid:     false,
			AccountID: "",
			PlayerID:  "",
		}, status.Error(codes.Internal, "jwt token pars error")
	}
	return result, nil
}

func (g *grpcService) generateJwtToken(accountID string, playerID string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS512)
	claim := token.Claims.(jwt.MapClaims)
	claim["exp"] = time.Now().Add(10 * time.Minute).Unix()
	claim["id"] = accountID
	claim["player_id"] = playerID
	tokenStr, err := token.SignedString([]byte(g.cfg.JWT.Secret))
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}

func (g *grpcService) validateToken(tokenString string) (*authService.VerifyTokenRes, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("jwt.Parse.Token.Method")
		}
		return []byte(g.cfg.JWT.Secret), nil
	})
	if err != nil {
		return &authService.VerifyTokenRes{
			Valid:     false,
			AccountID: "",
			PlayerID:  "",
		}, nil
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return &authService.VerifyTokenRes{
			Valid:     true,
			AccountID: claims["id"].(string),
			PlayerID:  claims["player_id"].(string),
		}, nil
	} else {
		return &authService.VerifyTokenRes{
			Valid:     false,
			AccountID: "",
			PlayerID:  "",
		}, nil
	}
}

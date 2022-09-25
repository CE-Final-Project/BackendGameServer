package grpc

import (
	"context"
	"github.com/ce-final-project/backend_game_server/authentication/config"
	"github.com/ce-final-project/backend_game_server/authentication/internal/auth/commands"
	"github.com/ce-final-project/backend_game_server/authentication/internal/auth/queries"
	"github.com/ce-final-project/backend_game_server/authentication/internal/auth/service"
	authService "github.com/ce-final-project/backend_game_server/authentication/proto"
	"github.com/ce-final-project/backend_game_server/pkg/logger"
	"github.com/ce-final-project/backend_game_server/pkg/tracing"
	"github.com/ce-final-project/backend_game_server/pkg/utils"
	"github.com/go-playground/validator"
	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/speps/go-hashids/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"math/rand"
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
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "grpcService.LoginAccount")
	defer span.Finish()

	query := queries.NewGetAccountByUsernameQuery(req.GetUsername())
	if err := g.v.StructCtx(ctx, query); err != nil {
		g.log.WarnMsg("validate", err)
		return nil, g.errResponse(codes.InvalidArgument, err)
	}

	account, err := g.as.Queries.GetAccountByUsername.Handle(ctx, query)
	if err != nil {
		g.log.WarnMsg("Register.Handle", err)
		return nil, g.errResponse(codes.Internal, err)
	}

	if !utils.CheckPasswordHash(req.GetPassword(), account.PasswordHashed) {
		err = errors.New("Invalid Username or Password")
		g.log.WarnMsg("Unauthenticated", err)
		return nil, g.errResponse(codes.Unauthenticated, err)
	}

	token, err := g.generateJwtToken(account.GetAccountID(), account.GetPlayerID())
	if err != nil {
		g.log.WarnMsg("GenerateJWTToken", err)
		return nil, g.errResponse(codes.Internal, err)
	}

	return &authService.LoginRes{
		Token:     token,
		AccountID: account.GetAccountID(),
		PlayerID:  account.GetPlayerID(),
	}, nil
}

func (g *grpcService) Register(ctx context.Context, req *authService.RegisterReq) (*authService.RegisterRes, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "grpcService.RegisterAccount")
	defer span.Finish()

	// Account UUID
	accountUUID := uuid.NewV4()

	// Player Short ID
	hd := hashids.NewData()
	hd.Salt = req.GetEmail() + req.GetUsername()
	hd.MinLength = 5
	h, err := hashids.NewWithData(hd)
	if err != nil {
		g.log.WarnMsg("hashIds.NewWithData", err)
		return nil, err
	}
	var playerID string
	playerID, err = h.Encode([]int{rand.Intn(5000)})
	if err != nil {
		g.log.WarnMsg("h.Encode", err)
		return nil, err
	}

	command := commands.NewRegisterCommand(accountUUID, playerID, req.GetUsername(), req.GetPassword(), req.GetPassword())
	if err := g.v.StructCtx(ctx, command); err != nil {
		g.log.WarnMsg("validate", err)
		return nil, g.errResponse(codes.InvalidArgument, err)
	}

	err = g.as.Commands.Register.Handle(ctx, command)
	if err != nil {
		g.log.WarnMsg("Register.Handle", err)
		return nil, g.errResponse(codes.Internal, err)
	}
	token, err := g.generateJwtToken(accountUUID.String(), playerID)
	if err != nil {
		g.log.WarnMsg("GenerateJWTToken", err)
		return nil, g.errResponse(codes.Internal, err)
	}
	return &authService.RegisterRes{
		Token:     token,
		AccountID: accountUUID.String(),
		PlayerID:  playerID,
	}, nil
}

func (g *grpcService) VerifyToken(_ context.Context, req *authService.VerifyTokenReq) (*authService.VerifyTokenRes, error) {
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

func (g *grpcService) errResponse(c codes.Code, err error) error {
	return status.Error(c, err.Error())
}

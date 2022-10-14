package grpc

import (
	"context"
	"database/sql"
	"github.com/ce-final-project/backend_game_server/authentication/config"
	accountCommands "github.com/ce-final-project/backend_game_server/authentication/internal/account/commands"
	accountQueries "github.com/ce-final-project/backend_game_server/authentication/internal/account/queries"
	accountService "github.com/ce-final-project/backend_game_server/authentication/internal/account/service"
	"github.com/ce-final-project/backend_game_server/authentication/internal/models"
	authGRPCService "github.com/ce-final-project/backend_game_server/authentication/proto"
	"github.com/ce-final-project/backend_game_server/pkg/logger"
	"github.com/ce-final-project/backend_game_server/pkg/tracing"
	"github.com/ce-final-project/backend_game_server/pkg/utils"
	"github.com/go-playground/validator"
	"github.com/golang-jwt/jwt/v4"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/speps/go-hashids/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"math/rand"
	"strconv"
	"time"
)

type grpcService struct {
	log logger.Logger
	cfg *config.Config
	v   *validator.Validate
	acc *accountService.AccountService
}

func NewAuthGRPCService(log logger.Logger, cfg *config.Config, v *validator.Validate, acc *accountService.AccountService) authGRPCService.AuthServiceServer {
	return &grpcService{
		log: log,
		cfg: cfg,
		v:   v,
		acc: acc,
	}
}

func (g *grpcService) Login(ctx context.Context, req *authGRPCService.LoginReq) (*authGRPCService.LoginRes, error) {
	var span opentracing.Span
	ctx, span = tracing.StartGrpcServerTracerSpan(ctx, "grpcService.LoginAccount")
	defer span.Finish()

	isEmail := g.isEmail(ctx, isEmail{Email: req.GetUsername()})
	var account *models.Account

	if isEmail {
		query := accountQueries.NewGetAccountByEmailQuery(req.GetUsername())
		if err := g.v.StructCtx(ctx, query); err != nil {
			g.log.WarnMsg("validate", err)
			return nil, g.errResponse(codes.InvalidArgument, err)
		}
		var err error
		account, err = g.acc.Queries.GetAccountByEmail.Handle(ctx, query)
		if err != nil {
			g.log.WarnMsg("Register.Handle", err)
			if errors.Cause(err) == sql.ErrNoRows {
				return nil, g.errResponse(codes.NotFound, err)
			}
			return nil, g.errResponse(codes.Internal, err)
		}
	} else {
		query := accountQueries.NewGetAccountByUsernameQuery(req.GetUsername())
		if err := g.v.StructCtx(ctx, query); err != nil {
			g.log.WarnMsg("validate", err)
			return nil, g.errResponse(codes.InvalidArgument, err)
		}
		var err error
		account, err = g.acc.Queries.GetAccountByUsername.Handle(ctx, query)
		if err != nil {
			g.log.WarnMsg("Register.Handle", err)
			if errors.Cause(err) == sql.ErrNoRows {
				return nil, g.errResponse(codes.NotFound, err)
			}
			return nil, g.errResponse(codes.Internal, err)
		}
	}

	if !utils.CheckPasswordHash(req.GetPassword(), account.PasswordHashed) {
		err := errors.New("Invalid Username or Password")
		g.log.WarnMsg("Unauthenticated", err)
		return nil, g.errResponse(codes.Unauthenticated, err)
	}

	token, err := g.generateJwtToken(strconv.FormatUint(account.ID, 10), account.PlayerID)
	if err != nil {
		g.log.WarnMsg("GenerateJWTToken", err)
		return nil, g.errResponse(codes.Internal, err)
	}

	return &authGRPCService.LoginRes{
		Token:     token,
		AccountID: account.ID,
		PlayerID:  account.PlayerID,
	}, nil
}

func (g *grpcService) Register(ctx context.Context, req *authGRPCService.RegisterReq) (*authGRPCService.RegisterRes, error) {
	var span opentracing.Span
	ctx, span = tracing.StartGrpcServerTracerSpan(ctx, "grpcService.RegisterAccount")
	defer span.Finish()

	// Player Short ID
	hd := hashids.NewData()
	hd.Salt = req.GetEmail() + req.GetUsername()
	hd.MinLength = 5
	h, err := hashids.NewWithData(hd)
	if err != nil {
		g.log.WarnMsg("hashIds.NewWithData", err)
		return nil, g.errResponse(codes.Internal, err)
	}
	var playerID string
	playerID, err = h.Encode([]int{rand.Intn(5000)})
	if err != nil {
		g.log.WarnMsg("h.Encode", err)
		return nil, g.errResponse(codes.Internal, err)
	}

	command := accountCommands.NewCreateAccountCommand(playerID, req.GetUsername(), req.GetEmail(), req.GetPassword())
	if err := g.v.StructCtx(ctx, command); err != nil {
		g.log.WarnMsg("validate", err)
		return nil, g.errResponse(codes.InvalidArgument, err)
	}
	var account *models.Account
	account, err = g.acc.Commands.CreateAccount.Handle(ctx, command)
	if err != nil {
		g.log.WarnMsg("Register.Handle", err)
		return nil, g.errResponse(codes.Internal, err)
	}
	var token string
	token, err = g.generateJwtToken(strconv.FormatUint(account.ID, 10), account.PlayerID)
	if err != nil {
		g.log.WarnMsg("GenerateJWTToken", err)
		return nil, g.errResponse(codes.Internal, err)
	}
	return &authGRPCService.RegisterRes{
		Token:     token,
		AccountID: account.ID,
		PlayerID:  account.PlayerID,
	}, nil
}

func (g *grpcService) VerifyToken(ctx context.Context, req *authGRPCService.VerifyTokenReq) (*authGRPCService.VerifyTokenRes, error) {
	var span opentracing.Span
	ctx, span = tracing.StartGrpcServerTracerSpan(ctx, "grpcService.VerifyToken")
	defer span.Finish()

	result, err := g.validateToken(req.GetToken())
	if err != nil {
		g.log.Error(err)
		return nil, g.errResponse(codes.Unauthenticated, err)
	}
	var account *models.Account
	query := accountQueries.NewGetAccountByIdQuery(result.GetAccountID())
	account, err = g.acc.Queries.GetAccountById.Handle(ctx, query)
	if err != nil {
		g.log.WarnMsg("GetAccount_VerifyToken.Handle", err)
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, g.errResponse(codes.Unauthenticated, err)
		}
		return nil, g.errResponse(codes.Internal, err)
	}
	if account == nil {
		return nil, g.errResponse(codes.Unauthenticated, err)
	}

	return result, nil
}

func (g *grpcService) generateJwtToken(accountID string, playerID string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS512)
	claim := token.Claims.(jwt.MapClaims)
	m, err := time.ParseDuration(g.cfg.JWT.ExpireTime)
	if err != nil {
		return "", err
	}
	claim["exp"] = time.Now().Add(m).Unix()
	claim["id"] = accountID
	claim["player_id"] = playerID
	tokenStr, err := token.SignedString([]byte(g.cfg.JWT.Secret))
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}

func (g *grpcService) validateToken(tokenString string) (*authGRPCService.VerifyTokenRes, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("jwt.Parse.Token.Method")
		}
		return []byte(g.cfg.JWT.Secret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return &authGRPCService.VerifyTokenRes{
			Valid:     true,
			AccountID: claims["id"].(uint64),
			PlayerID:  claims["player_id"].(string),
		}, nil
	} else {
		return nil, errors.New("Invalid Token")
	}
}

func (g *grpcService) errResponse(c codes.Code, err error) error {
	return status.Error(c, err.Error())
}

type isEmail struct {
	Email string `validate:"required,email,max=320"`
}

func (g *grpcService) isEmail(ctx context.Context, isEmail isEmail) bool {
	if err := g.v.StructCtx(ctx, isEmail); err != nil {
		return false
	}
	return true
}

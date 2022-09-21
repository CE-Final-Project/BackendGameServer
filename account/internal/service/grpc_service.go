package service

import (
	"context"
	"database/sql"
	"github.com/ce-final-project/backend_game_server/authentication/config"
	"github.com/ce-final-project/backend_game_server/authentication/internal/dto"
	authService "github.com/ce-final-project/backend_game_server/authentication/proto"
	"github.com/ce-final-project/backend_game_server/pkg/logger"
	"github.com/ce-final-project/backend_game_server/pkg/utils"
	"github.com/go-playground/validator"
	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type grpcService struct {
	log logger.Logger
	cfg *config.Config
	v   *validator.Validate
	as  *AccountService
}

func NewGrpcService(log logger.Logger, cfg *config.Config, v *validator.Validate, as *AccountService) *grpcService {
	return &grpcService{
		log: log,
		cfg: cfg,
		v:   v,
		as:  as,
	}
}

func (g *grpcService) Login(ctx context.Context, payload *authService.LoginReq) (*authService.LoginRes, error) {
	account, err := g.as.GetAccountByUsernameOrEmail(ctx, payload.Username, "")
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, status.Error(codes.Unauthenticated, "invalid Username or Password")
		}
		return nil, status.Errorf(codes.Internal, "GetAccountByUsernameOrEmail Error: %v", err)
	}
	if correct := utils.CheckPasswordHash(payload.GetPassword(), account.PasswordHashed); !correct {
		return nil, status.Error(codes.Unauthenticated, "invalid Username or Password")
	}

	var token string
	token, err = g.GenerateJwtToken(account.AccountID.String(), account.PlayerID)
	if err != nil {
		g.log.Error(err)
		return nil, status.Error(codes.Internal, "can not generate token")
	}
	return &authService.LoginRes{
		Token:     token,
		AccountID: account.AccountID.String(),
		PlayerID:  account.PlayerID,
	}, nil
}

func (g *grpcService) Register(ctx context.Context, payload *authService.RegisterReq) (*authService.RegisterRes, error) {

	result, err := g.as.CreateNewAccount(ctx, &dto.CreateAccountDto{
		Username: payload.GetUsername(),
		Email:    payload.GetEmail(),
		Password: payload.GetPassword(),
	})
	if err != nil {
		g.log.Error(err)
		return nil, status.Error(codes.InvalidArgument, "can not create account user")
	}
	var token string
	token, err = g.GenerateJwtToken(result.AccountID.String(), result.PlayerID)
	if err != nil {
		g.log.Error(err)
		return nil, status.Error(codes.Internal, "can not generate token")
	}

	return &authService.RegisterRes{
		Token:     token,
		AccountID: result.AccountID.String(),
		PlayerID:  result.PlayerID,
	}, nil
}

func (g *grpcService) UpdateAccount(ctx context.Context, payload *authService.UpdateAccountReq) (*authService.UpdateAccountRes, error) {
	accountUUID, err := uuid.FromString(payload.GetAccountID())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "account uuid invalid format")
	}
	result, err := g.as.UpdateAccount(ctx, &dto.UpdateAccountDto{
		AccountID: accountUUID,
		Username:  payload.Username,
		Email:     payload.Email,
	})
	if err != nil {
		g.log.Error(err)
		return nil, status.Error(codes.Internal, "can not update account")
	}
	return &authService.UpdateAccountRes{
		AccountID: result.AccountID.String(),
		UpdatedAt: timestamppb.New(result.UpdatedAt),
	}, nil
}

func (g *grpcService) ChangePassword(ctx context.Context, payload *authService.ChangePasswordReq) (*authService.UpdateAccountRes, error) {
	accountUUID, err := uuid.FromString(payload.GetAccountID())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "account uuid invalid format")
	}
	account, err := g.as.GetAccountById(ctx, accountUUID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "GetAccountById Error: %v", err)
	}
	if correct := utils.CheckPasswordHash(payload.GetOldPassword(), account.PasswordHashed); !correct {
		return nil, status.Error(codes.Unauthenticated, "invalid Username or Password")
	}

	result, err := g.as.ChangePasswordAccount(ctx, &dto.ChangePasswordDto{
		AccountID:   accountUUID,
		OldPassword: payload.GetOldPassword(),
		NewPassword: payload.GetNewPassword(),
	})

	if err != nil {
		g.log.Error(err)
		return nil, status.Error(codes.Internal, "can not change password")
	}
	return &authService.UpdateAccountRes{
		AccountID: result.AccountID.String(),
		UpdatedAt: timestamppb.New(result.UpdatedAt),
	}, nil
}

func (g *grpcService) GetAccountById(ctx context.Context, payload *authService.GetAccountByIdReq) (*authService.Account, error) {

	accountUUID, err := uuid.FromString(payload.GetAccountID())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "account uuid invalid format")
	}
	account, err := g.as.GetAccountById(ctx, accountUUID)
	if err != nil {
		return nil, status.Error(codes.Internal, "can not get account by id")
	}
	return &authService.Account{
		AccountID: account.AccountID.String(),
		PlayerID:  account.PlayerID,
		Username:  account.Username,
		Email:     account.Email,
		IsBan:     account.IsBan,
		CreatedAt: timestamppb.New(account.CreatedAt),
		UpdatedAt: timestamppb.New(account.UpdatedAt),
	}, nil
}

func (g *grpcService) VerifyToken(ctx context.Context, payload *authService.VerifyTokenReq) (*authService.VerifyTokenRes, error) {
	result, err := g.ValidateToken(payload.GetToken())
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

func (g *grpcService) GenerateJwtToken(accountID string, playerID string) (string, error) {
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

func (g *grpcService) ValidateToken(tokenString string) (*authService.VerifyTokenRes, error) {
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

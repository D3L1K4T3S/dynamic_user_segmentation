package service

import (
	"context"
	"dynamic-user-segmentation/internal/entity"
	"dynamic-user-segmentation/internal/repository"
	"dynamic-user-segmentation/internal/repository/respository_errors"
	"dynamic-user-segmentation/internal/service/dto"
	"dynamic-user-segmentation/pkg/hash"
	e "dynamic-user-segmentation/pkg/util/errors"
	"errors"
	"github.com/golang-jwt/jwt"
	"time"
)

type AuthService struct {
	signKey  string
	tokenTTL time.Duration

	usersRepository repository.Users
	hash            hash.PasswordHash
}

type Token struct {
	UserId int
	jwt.StandardClaims
}

func NewAuthService(usersRepository repository.Users, ph hash.PasswordHash, signKey string, tokenTTL time.Duration) *AuthService {
	return &AuthService{
		signKey:         signKey,
		tokenTTL:        tokenTTL,
		usersRepository: usersRepository,
		hash:            ph,
	}
}

func (as *AuthService) CreateUser(ctx context.Context, authUser dto.AuthUser) (int, error) {
	var err error
	defer func() {
		err = e.WrapIfErr("AuthService: ", err)
	}()

	user := entity.Users{
		Username: authUser.Username,
		Password: as.hash.Hash(authUser.Password),
	}

	userId, err := as.usersRepository.CreateUser(ctx, user)
	if err != nil {
		if errors.Is(err, respository_errors.ErrAlreadyExists) {
			return 0, respository_errors.ErrAlreadyExists
		}
		return 0, e.Wrap("can't create user: ", err)
	}
	return userId, nil
}
func (as *AuthService) DeleteUser(ctx context.Context, authUser dto.AuthUser) error {
	var err error
	defer func() {
		err = e.WrapIfErr("AuthService: ", err)
	}()

	user := entity.Users{
		Username: authUser.Username,
		Password: as.hash.Hash(authUser.Password),
	}

	err = as.usersRepository.DeleteUser(ctx, user)
	if err != nil {
		if errors.Is(err, respository_errors.ErrNotFound) {
			return respository_errors.ErrNotFound
		}
		return e.Wrap("can't delete user: ", err)
	}
	return nil
}
func (as *AuthService) GenerateToken(ctx context.Context, authUser dto.AuthUser) (string, error) {
	user, err := as.usersRepository.GetUserByUsername(ctx, authUser.Username)
	if err != nil {
		if errors.Is(err, respository_errors.ErrNotFound) {
			return "", respository_errors.ErrNotFound
		}
		return "", e.Wrap("can't get user: ", err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &Token{
		UserId: user.Id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(as.tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	})

	tokenString, err := token.SignedString([]byte(as.signKey))
	if err != nil {
		return "", ErrCannotCreateToken
	}
	return tokenString, nil
}
func (as *AuthService) ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &Token{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, ErrUnexpectedSigningMethod
			}
			return []byte(as.signKey), nil
		})

	if err != nil {
		return 0, ErrCannotParseToken
	}
	claims, ok := token.Claims.(*Token)
	if !ok {
		return 0, ErrCannotParseToken
	}
	return claims.UserId, nil
}

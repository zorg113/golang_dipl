package grpcapi

import (
	"context"
	"errors"

	"github.com/rs/zerolog"
	"github.com/zorg113/golang_dipl/atibruteforce/internal/controller/grpcapi/authorizationpb"
	"github.com/zorg113/golang_dipl/atibruteforce/internal/controller/httpapi/handlers"
	"github.com/zorg113/golang_dipl/atibruteforce/internal/controller/proto/authorizationpb"
	"github.com/zorg113/golang_dipl/atibruteforce/model/entity"
	"github.com/zorg113/golang_dipl/atibruteforce/model/service"
)

type AuthorizationServer struct {
	authorizationpb.UnimplementedAuthorizationSrver
	service *service.Authorization
	log     *zerolog.Logger
}

func NewAuthosization(service *service.Authorization, log *zerolog.Logger) *AuthorizationServer {
	return &AuthorizationServer{service: service, log: log}
}

func (s *AuthorizationServer) Authorization(ctx context.Context, in *authorizationpb.AuthorizationRequest) (*authorizationpb.AuthorizationResponse, error) {
	s.log.Info().Msgf("Authorization request GRPC")
	req := entity.Request{
		Login:    in.GetRequest().GetLogin(),
		Password: in.GetRequest().GetPassword(),
		Ip:       in.GetRequest().GetIp(),
	}
	if !handlers.ValidateRequest(req) {
		return nil, errors.New("Invalid authorization request")
	}
	isAlllowed, err := s.service
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to authorization request")
		return nil, err
	}
	return &authorizationpb.AuthorizationResponse{IsAllowed: isAllowed}, nil
}

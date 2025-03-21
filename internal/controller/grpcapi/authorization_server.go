package grpcapi

import (
	"context"
	"errors"

	"github.com/rs/zerolog"
	"github.com/zorg113/golang_dipl/atibruteforce/internal/common"
	"github.com/zorg113/golang_dipl/atibruteforce/internal/controller/grpcapi/authorizationpb"
	"github.com/zorg113/golang_dipl/atibruteforce/model/entity"
	"github.com/zorg113/golang_dipl/atibruteforce/model/service"
)

type AuthorizationServer struct {
	authorizationpb.UnimplementedAuthorizationServer
	service *service.Authorization
	log     *zerolog.Logger
}

func NewAuthorization(service *service.Authorization, log *zerolog.Logger) *AuthorizationServer {
	return &AuthorizationServer{service: service, log: log}
}

func (s *AuthorizationServer) Authorization(_ context.Context, in *authorizationpb.AuthorizationRequest) (*authorizationpb.AuthorizationResponse, error) { //nolint:lll
	s.log.Info().Msgf("Authorization request GRPC")
	req := entity.Request{
		Login:    in.GetRequest().GetLogin(),
		Password: in.GetRequest().GetPassword(),
		IP:       in.GetRequest().GetIp(),
	}
	if !common.ValidateRequest(req) {
		return nil, errors.New("invalid authorization request")
	}
	isAllowed, err := s.service.Authorization(req)
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to authorization request")
		return nil, err
	}
	return &authorizationpb.AuthorizationResponse{IsAllow: isAllowed}, nil
}

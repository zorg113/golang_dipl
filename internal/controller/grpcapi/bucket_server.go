package grpcapi

import (
	"context"
	"errors"

	"github.com/rs/zerolog"
	"github.com/zorg113/golang_dipl/atibruteforce/internal/common"
	"github.com/zorg113/golang_dipl/atibruteforce/internal/controller/grpcapi/bucketpb"
	"github.com/zorg113/golang_dipl/atibruteforce/model/entity"
	"github.com/zorg113/golang_dipl/atibruteforce/model/service"
)

type BucketServer struct {
	bucketpb.UnimplementedBucketServiceServer
	service *service.Authorization
	log     *zerolog.Logger
}

func NewBucketServer(service *service.Authorization, log *zerolog.Logger) *BucketServer {
	return &BucketServer{service: service, log: log}
}

func (s *BucketServer) ResetBucket(_ context.Context, req *bucketpb.ResetBucketRequest) (*bucketpb.ResetBucketResponse, error) { //nolint:lll
	s.log.Info().Msg("resetting bucket GRPC")
	request := entity.Request{
		Login:    req.GetRequest().GetLogin(),
		Password: req.GetRequest().GetPassword(),
		IP:       req.GetRequest().GetIp(),
	}
	request.Password = "empty"
	if !common.ValidateRequest(request) {
		return nil, errors.New("invalid input request received")
	}
	resp := &bucketpb.ResetBucketResponse{ResetLogin: true, ResetIp: true}
	isLoginReset := s.service.ResetLoginInBucket(request.Login)
	if !isLoginReset {
		resp.ResetLogin = false
	}
	isIPReset := s.service.ResetIPBucket(request.IP)
	if !isIPReset {
		resp.ResetIp = false
	}
	return resp, nil
}

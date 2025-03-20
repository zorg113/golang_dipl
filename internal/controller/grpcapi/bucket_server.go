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

func (s *BucketServer) ResetBucket(ctx context.Context, req *bucketpb.ResetBucketRequest) (*bucketpb.ResetBucketResponse, error) {
	s.log.Info().Msg("resetting bucket GRPC")
	request := entity.Request{
		Login:    req.GetRequest().GetLogin(),
		Password: req.GetRequest().GetPassword(),
		Ip:       req.GetRequest().GetIp(),
	}
	request.Password = "empty"
	if !common.ValidateRequest(request) {
		return nil, errors.New("invalid input request recived")
	}
	resp := &bucketpb.ResetBucketResponse{ResetLogin: true, ResetIp: true}
	isLoginReset := s.service.ResetLoginInBucket(request.Login)
	if !isLoginReset {
		resp.ResetLogin = false
	}
	isIpReset := s.service.ResetIpBucket(request.Ip)
	if !isIpReset {
		resp.ResetIp = false
	}
	return resp, nil
}

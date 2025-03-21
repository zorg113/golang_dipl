package grpcapi

import (
	"context"
	"errors"

	"github.com/rs/zerolog"
	"github.com/zorg113/golang_dipl/atibruteforce/internal/common"
	"github.com/zorg113/golang_dipl/atibruteforce/internal/controller/grpcapi/blacklistpb"
	"github.com/zorg113/golang_dipl/atibruteforce/model/entity"
	"github.com/zorg113/golang_dipl/atibruteforce/model/service"
)

var errInvalidInputIP = errors.New("invalid input IP address from client")

type BlackListServer struct {
	blacklistpb.UnimplementedBlackListServiceServer
	service *service.BlackList
	log     *zerolog.Logger
}

func NewBlackListServer(service *service.BlackList, log *zerolog.Logger) *BlackListServer {
	return &BlackListServer{service: service, log: log}
}

func (s *BlackListServer) AddIP(_ context.Context, in *blacklistpb.AddIpRequest) (*blacklistpb.AddIpResponse, error) {
	s.log.Info().Msg("add IP to blacklist by GRPC")
	ipNetwork := entity.IPNetwork{
		IP:   in.GetIpNetwork().GetIp(),
		Mask: in.GetIpNetwork().GetMask(),
	}
	idValideted := common.ValidateIP(ipNetwork)
	if !idValideted {
		s.log.Info().Msg("invalid IP address")
		return nil, errInvalidInputIP
	}
	err := s.service.AddIP(ipNetwork)
	if err != nil {
		s.log.Error().Err(err).Msg("failed to add IP to blacklist")
		return nil, err
	}
	return &blacklistpb.AddIpResponse{IsAddIp: true}, nil
}

func (s *BlackListServer) RemoveIP(_ context.Context, req *blacklistpb.RemoveIPRequest) (*blacklistpb.RemoveIPResponse, error) { //nolint:lll
	s.log.Info().Msg("removing IP address from blacklist GRPC")
	ipNetwork := entity.IPNetwork{
		IP:   req.GetIpNetwork().GetIp(),
		Mask: req.GetIpNetwork().GetMask(),
	}
	idValideted := common.ValidateIP(ipNetwork)
	if !idValideted {
		s.log.Info().Msg("invalid IP address")
		return nil, errInvalidInputIP
	}
	err := s.service.DeleteIP(ipNetwork)
	if err != nil {
		s.log.Error().Err(err).Msg("failed to remove IP from blacklist")
		return nil, err
	}
	return &blacklistpb.RemoveIPResponse{IsRemoveIp: true}, nil
}

func (s *BlackListServer) GetIPs(_ context.Context, stream blacklistpb.BlackListService_GetIpListServer) error {
	s.log.Info().Msg("getting IP addresses from blacklist GRPC")
	ips, err := s.service.GetIPs()
	if err != nil {
		s.log.Error().Err(err).Msg("failed to get IP addresses from blacklist")
		return err
	}
	for _, net := range ips {
		err := stream.Send(&blacklistpb.GetIpListResponse{IpNetwork: &blacklistpb.IpNetwork{
			Ip:   net.IP,
			Mask: net.Mask,
		}})
		if err != nil {
			s.log.Error().Err(err).Msg("failed to send IP address to client")
			return err
		}
	}
	return nil
}

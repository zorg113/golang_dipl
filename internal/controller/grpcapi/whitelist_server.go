package grpcapi

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/zorg113/golang_dipl/atibruteforce/internal/common"
	"github.com/zorg113/golang_dipl/atibruteforce/internal/controller/grpcapi/whitelistpb"
	"github.com/zorg113/golang_dipl/atibruteforce/model/entity"
	"github.com/zorg113/golang_dipl/atibruteforce/model/service"
)

type WhiteListServer struct {
	whitelistpb.UnimplementedWhiteListServiceServer
	service *service.WhiteList
	log     *zerolog.Logger
}

func NewWhiteListServer(service *service.WhiteList, log *zerolog.Logger) *WhiteListServer {
	return &WhiteListServer{service: service, log: log}
}
func (s *WhiteListServer) AddIp(ctx context.Context, req *whitelistpb.AddIpRequest) (*whitelistpb.AddIpResponse, error) {
	s.log.Info().Msg("add ip to whitelist GRPC")
	ipNetwork := entity.IpNetwork{
		Ip:   req.GetIpNetwork().GetIp(),
		Mask: req.GetIpNetwork().GetMask(),
	}
	if !common.ValidateIP(ipNetwork) {
		s.log.Error().Msg("invalid IP address")
		return nil, errInvalidInputIP
	}
	err := s.service.AddIP(ipNetwork)
	if err != nil {
		s.log.Error().Err(err).Msg("failed to add IP to whitelist")
		return nil, err
	}
	return &whitelistpb.AddIpResponse{IsAddIp: true}, nil
}

func (s *WhiteListServer) RemoveIp(ctx context.Context, req *whitelistpb.RemoveIPRequest) (*whitelistpb.RemoveIPResponse, error) {
	s.log.Info().Msg("removing IP from whitelist GRPC")
	ipNetwork := entity.IpNetwork{
		Ip:   req.GetIpNetwork().GetIp(),
		Mask: req.GetIpNetwork().GetMask(),
	}
	if !common.ValidateIP(ipNetwork) {
		s.log.Error().Msg("invalid IP address")
		return nil, errInvalidInputIP
	}
	err := s.service.DeleteIP(ipNetwork)
	if err != nil {
		s.log.Error().Err(err).Msg("failed to remove IP from whitelist")
		return nil, err
	}
	return &whitelistpb.RemoveIPResponse{IsRemoveIp: true}, nil
}

func (s *WhiteListServer) GetIPs(ctx context.Context, stream whitelistpb.WhiteListService_GetIpListServer) error {
	s.log.Info().Msg("getting IP addresses from whitelist GRPC")
	ips, err := s.service.GetIPs()
	if err != nil {
		s.log.Error().Err(err).Msg("failed to get IP addresses from whitelist")
		return err
	}
	for _, net := range ips {
		ip := &whitelistpb.IpNetwork{
			Ip:   net.Ip,
			Mask: net.Mask,
		}
		err := stream.Send(&whitelistpb.GetIpListResponse{IpNetwork: ip})
		if err != nil {
			return err
		}
	}
	return nil
}

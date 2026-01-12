package main

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/zorg113/golang_dipl/atibruteforce/internal/config"
	"github.com/zorg113/golang_dipl/atibruteforce/internal/controller/grpcapi/authorizationpb"
	"github.com/zorg113/golang_dipl/atibruteforce/internal/controller/grpcapi/blacklistpb"
	"github.com/zorg113/golang_dipl/atibruteforce/internal/controller/grpcapi/bucketpb"
	"github.com/zorg113/golang_dipl/atibruteforce/internal/controller/grpcapi/whitelistpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conf, err := config.NewConfig("./../../config/conf.yaml")
	if err != nil {
		fmt.Println(err)
	}

	trn := grpc.WithTransportCredentials(insecure.NewCredentials())
	dial, err := grpc.NewClient(conf.Listen.BindIP+":"+conf.Listen.Port, trn)
	if err != nil {
		fmt.Printf("Failed to dial: %v", err)
		return
	}
	clientBlackList := blacklistpb.NewBlackListServiceClient(dial)
	clientWhiteList := whitelistpb.NewWhiteListServiceClient(dial)
	clientBucket := bucketpb.NewBucketServiceClient(dial)
	clientAuthorization := authorizationpb.NewAuthorizationClient(dial)
	getIPListInBlackList(clientBlackList)
	fmt.Println()
	getIPListInWhiteList(clientWhiteList)
	fmt.Println()
	resetBucket(clientBucket)
	fmt.Println()
	execAuthorization(clientAuthorization)
}

func execAuthorization(client authorizationpb.AuthorizationClient) {
	resp, err := client.TryAuthorization(context.Background(),
		&authorizationpb.AuthorizationRequest{
			Request: &authorizationpb.Request{
				Login:    "admin",
				Password: "admin",
				Ip:       "192.168.0.1",
			},
		})
	if err != nil {
		fmt.Printf("Error: %v", err)
	} else {
		fmt.Printf("Authorization status: %v\n", resp.GetIsAllow())
	}
}

func getIPListInBlackList(client blacklistpb.BlackListServiceClient) {
	stream, err := client.GetIpList(context.Background(), &blacklistpb.GetIpListRequest{})
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}
	for {
		ip, err := stream.Recv()
		if errors.Is(err, context.Canceled) {
			break
		}
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			fmt.Printf("Error: %v", err)
			return
		}
		fmt.Printf("IP: %s\n", ip.GetIpNetwork())
	}
}

func getIPListInWhiteList(client whitelistpb.WhiteListServiceClient) {
	stream, err := client.GetIpList(context.Background(), &whitelistpb.GetIpListRequest{})
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}
	for {
		ip, err := stream.Recv()
		if errors.Is(err, context.Canceled) {
			break
		}
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			fmt.Printf("Error: %v", err)
			return
		}
		fmt.Printf("IP: %s\n", ip.GetIpNetwork())
	}
}

func resetBucket(client bucketpb.BucketServiceClient) {
	rest, err := client.ResetBucket(context.Background(),
		&bucketpb.ResetBucketRequest{
			Request: &bucketpb.Request{
				Login:    "admin",
				Password: "admin",
				Ip:       "192.168.0.1",
			},
		})
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
	fmt.Printf("Reset status: Login: %v, IP: %v\n", rest.GetResetLogin(), rest.GetResetIp())
}

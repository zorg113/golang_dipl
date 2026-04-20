package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/zorg113/golang_dipl/atibruteforce/internal/config"
	"github.com/zorg113/golang_dipl/atibruteforce/internal/controller/grpcapi/authorizationpb"
	"github.com/zorg113/golang_dipl/atibruteforce/internal/controller/grpcapi/blacklistpb"
	"github.com/zorg113/golang_dipl/atibruteforce/internal/controller/grpcapi/bucketpb"
	"github.com/zorg113/golang_dipl/atibruteforce/internal/controller/grpcapi/commonpb"
	"github.com/zorg113/golang_dipl/atibruteforce/internal/controller/grpcapi/whitelistpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func main() {
	conf, err := config.NewConfig("./../../config/conf.yaml")
	if err != nil {
		fmt.Println(err)
	}
	apiKey := os.Getenv("ADMIN_API_KEY")
	if apiKey == "" {
		fmt.Println("ADMIN_API_KEY is not set")
		return
	}

	adminCtx := metadata.AppendToOutgoingContext(
		context.Background(),
		"x-admin-key", apiKey,
	)

	trn := grpc.WithTransportCredentials(insecure.NewCredentials())
	dial, err := grpc.NewClient(conf.Listen.BindIP+":"+conf.Listen.Port, trn)
	if err != nil {
		fmt.Printf("Failed to dial: %v", err)
		return
	}
	clientBlackList := blacklistpb.NewBlackListClient(dial)
	clientWhiteList := whitelistpb.NewWhiteListClient(dial)
	clientBucket := bucketpb.NewBucketClient(dial)
	clientAuthorization := authorizationpb.NewAuthorizationClient(dial)
	getIPListInBlackList(adminCtx, clientBlackList)
	fmt.Println()
	getIPListInWhiteList(adminCtx, clientWhiteList)
	fmt.Println()
	resetBucket(adminCtx, clientBucket)
	fmt.Println()
	execAuthorization(context.Background(), clientAuthorization)
}

func execAuthorization(ctx context.Context, client authorizationpb.AuthorizationClient) {
	resp, err := client.TryAuthorization(ctx,
		&authorizationpb.AuthorizationRequest{
			Request: &commonpb.Request{
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

func getIPListInBlackList(ctx context.Context, client blacklistpb.BlackListClient) {
	stream, err := client.GetIpList(ctx, &blacklistpb.GetIpListRequest{})
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

func getIPListInWhiteList(ctx context.Context, client whitelistpb.WhiteListClient) {
	stream, err := client.GetIpList(ctx, &whitelistpb.GetIpListRequest{})
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

func resetBucket(ctx context.Context, client bucketpb.BucketClient) {
	rest, err := client.ResetBucket(ctx,
		&bucketpb.ResetBucketRequest{
			Request: &commonpb.Request{
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

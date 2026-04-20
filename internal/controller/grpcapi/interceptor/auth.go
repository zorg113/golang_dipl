package interceptor

import (
	"context"
	"crypto/subtle"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// adminMethods — методы, требующие аутентификации.
// Authorization публичный, всё остальное — административное.

var adminMethods = map[string]struct{}{
	"/blacklist.BlackListService/AddIP":     {},
	"/blacklist.BlackListService/RemoveIP":  {},
	"/blacklist.BlackListService/GetIpList": {},
	"/whitelist.WhiteListService/AddIP":     {},
	"/whitelist.WhiteListService/RemoveIP":  {},
	"/whitelist.WhiteListService/GetIpList": {},
	"/bucket.BucketService/ResetBucket":     {},
}

func AdminAuth(apiKey string, log *zerolog.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		// Проверяем только административные методы
		if _, isAdmin := adminMethods[info.FullMethod]; !isAdmin {
			return handler(ctx, req)
		}

		// Читаем метаданные — аналог HTTP-заголовков в gRPC
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			log.Warn().Str("method", info.FullMethod).Msg("no metadata in request")
			return nil, status.Error(codes.Unauthenticated, "missing credentials")
		}

		keys := md.Get("x-admin-key") // metadata ключи в нижнем регистре
		if len(keys) == 0 {
			log.Warn().Str("method", info.FullMethod).Msg("missing x-admin-key")
			return nil, status.Error(codes.Unauthenticated, "missing credentials")
		}

		// ConstantTimeCompare — защита от timing-атак, как и в HTTP
		if subtle.ConstantTimeCompare([]byte(keys[0]), []byte(apiKey)) != 1 {
			log.Warn().Str("method", info.FullMethod).Msg("invalid admin key")
			return nil, status.Error(codes.Unauthenticated, "invalid credentials")
		}

		return handler(ctx, req)
	}
}

// StreamAdminAuth — то же самое для стриминговых методов (GetIpList)

func StreamAdminAuth(apiKey string, log *zerolog.Logger) grpc.StreamServerInterceptor {
	return func(
		srv any,
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		if _, isAdmin := adminMethods[info.FullMethod]; !isAdmin {
			return handler(srv, ss)
		}

		md, ok := metadata.FromIncomingContext(ss.Context())
		if !ok {
			return status.Error(codes.Unauthenticated, "missing credentials")
		}

		keys := md.Get("x-admin-key")
		if len(keys) == 0 || subtle.ConstantTimeCompare([]byte(keys[0]), []byte(apiKey)) != 1 {
			log.Warn().Str("method", info.FullMethod).Msg("unauthorized stream request")
			return status.Error(codes.Unauthenticated, "invalid credentials")
		}

		return handler(srv, ss)
	}
}

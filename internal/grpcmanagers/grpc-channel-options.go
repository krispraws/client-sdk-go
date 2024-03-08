package grpcmanagers

import (
	"crypto/tls"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"

	"github.com/momentohq/client-sdk-go/config"
)

func GrpcChannelOptionsFromGrpcConfig(grpcConfig config.GrpcConfiguration) []grpc.DialOption {
	// Default to 5mb message sizes and keepalives turned on (defaults are set in NewStaticGrpcConfiguration)
	return []grpc.DialOption{
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(grpcConfig.GetMaxReceiveMessageLength()),
			grpc.MaxCallSendMsgSize(grpcConfig.GetMaxSendMessageLength()),
		),
		grpc.WithKeepaliveParams(
			keepalive.ClientParameters{
				PermitWithoutStream: grpcConfig.GetKeepAlivePermitWithoutCalls(),
				Time:                grpcConfig.GetKeepAliveTime(),
				Timeout:             grpcConfig.GetKeepAliveTimeout(),
			},
		),
	}
}

func TransportCredentialsChannelOption() grpc.DialOption {
	config := &tls.Config{
		InsecureSkipVerify: false,
	}
	return grpc.WithTransportCredentials(credentials.NewTLS(config))
}

func AllDialOptions(grpcConfig config.GrpcConfiguration, options ...grpc.DialOption) []grpc.DialOption {
	options = append(options, TransportCredentialsChannelOption())
	options = append(options, GrpcChannelOptionsFromGrpcConfig(grpcConfig)...)
	return options
}

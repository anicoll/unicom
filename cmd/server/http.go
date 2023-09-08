package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/anicoll/unicom/gen/pb/go/unicom/api/v1"
)

func runHTTPGateway(ctx context.Context, httpPort, grpcPort int) error {
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := pb.RegisterUnicomHandlerFromEndpoint(ctx, mux, fmt.Sprintf("localhost:%d", grpcPort), opts)
	if err != nil {
		return err
	}

	// Start HTTP server (and proxy calls to gRPC server endpoint)
	return http.ListenAndServe(fmt.Sprintf(":%d", httpPort), mux)
}

package main

import (
	"context"
	"fmt"
	config "github.com/gynshu-one/go-metric-collector/internal/config/agent"
	ag "github.com/gynshu-one/go-metric-collector/internal/controller/http/agent"
	"github.com/gynshu-one/go-metric-collector/internal/domain/service"
	"github.com/gynshu-one/go-metric-collector/proto"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"syscall"
	"time"
)

var (
	agent        ag.Handler
	buildVersion string
	buildDate    string
	buildCommit  string
)

func main() {
	if buildVersion == "" {
		buildVersion = "N/A"
	}
	if buildDate == "" {
		buildDate = "N/A"
	}
	if buildCommit == "" {
		buildCommit = "N/A"
	}
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// grpc
	conn, err := grpc.DialContext(ctx,
		config.GetConfig().Server.GRPCAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*20)),
	)

	if err != nil {
		log.Fatal().Err(err).Msg("failed to dial server")
	}

	defer func(conn *grpc.ClientConn) {
		err = conn.Close()
		if err != nil {
			log.Warn().Err(err).Msg("failed to close connection")
		}
	}(conn)

	waitForConnection(ctx, conn)
	grpcClient := proto.NewMetricServiceClient(conn)
	agent = ag.NewAgent(service.NewMemService(), grpcClient)

	live, err := grpcClient.Live(ctx, &emptypb.Empty{})
	if err != nil {
		return
	}
	if live.Message != "OK" {
		log.Fatal().Msgf("Server live: %s", live.Message)
	}

	log.Info().Msg("Agent started")
	time.Sleep(1 * time.Second)
	f, err := os.Create("server_mem.prof")
	if err != nil {
		log.Fatal().Err(err).Msg("could not create memory profile")
	}
	runtime.GC()
	if err = pprof.WriteHeapProfile(f); err != nil {
		log.Fatal().Err(err).Msg("could not write memory profile")
	}
	err = f.Close()
	if err != nil {
		log.Fatal().Err(err).Msg("could not close memory profile")
	}

	go agent.Start()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit

	log.Info().Msg("Shutdown Agent ...")

	// run func with timeout
	ctx, cancel = context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	agent.Stop(ctx)
}

func waitForConnection(ctx context.Context, conn *grpc.ClientConn) {
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()
	for {
		select {
		case <-ctx.Done():
			log.Fatal().Msg("failed to connect to server")
		default:
			state := conn.GetState()
			if state == connectivity.Ready {
				return
			}
			log.Info().Msgf("connection state: %s with address %s", state.String(), conn.Target())
			time.Sleep(1 * time.Second)
		}
	}
}

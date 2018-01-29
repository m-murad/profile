package main

import (
	"context"
	"flag"
	"fmt"
	"net"

	"github.com/caarlos0/env"
	"github.com/golang/glog"
	geopb "github.com/pondohva/profile/proto/geoip"
	pb "github.com/pondohva/profile/proto/profile"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
)

func getGeoData(ctx context.Context) (*geopb.GetGeoDataByIPResponse, error) {
	conn, err := grpc.Dial(Cfg.GeoIP, grpc.WithInsecure())
	md, _ := metadata.FromIncomingContext(ctx)
	ip := md["x-forwarded-for"][0]

	if err != nil {
		return nil, fmt.Errorf("can't connect to geodata %v", err)
	}
	defer conn.Close()

	c := geopb.NewGeoipClient(conn)

	r, err := c.GetGeoDataByIP(ctx, &geopb.GetGeoDataByIPRequest{IP: ip})

	if err != nil {
		return nil, fmt.Errorf("can't return geodata %v", err)
	}
	return r, nil
}

func getUserName(ctx context.Context) string {
	return "anonymous"
}

type server struct{}

type config struct {
	GeoIP   string `env:"GIZMO_GEOIP_URL"`
	Address string `env:"GIZMO_PROFILE_ADDRESS" envDefault:"0.0.0.0"`
	Port    int    `env:"GIZMO_PROFILE_PORT" envDefault:"52082"`
}

// Cfg contains env variables
var Cfg config

func (s *server) GetProfile(ctx context.Context, in *pb.GetProfileByIDRequest) (*pb.GetProfileByIDResponse, error) {

	geodata, err := getGeoData(ctx)

	if err != nil {
		return nil, err
	}

	username := getUserName(ctx)
	return &pb.GetProfileByIDResponse{Country: geodata.CountryName, UserName: username, ID: "1"}, nil
}

func main() {
	flag.CommandLine.Parse([]string{})
	defer glog.Flush()

	err := env.Parse(&Cfg)

	lis, err := net.Listen("tcp", fmt.Sprintf("%v:%v", Cfg.Address, Cfg.Port))

	if err != nil {
		glog.Fatalf("can't bind to %v:%v", Cfg.Address, Cfg.Port)
	}

	s := grpc.NewServer()

	pb.RegisterProfileServer(s, &server{})
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		glog.Fatalf("failed to serve: %v", err)
	}
}

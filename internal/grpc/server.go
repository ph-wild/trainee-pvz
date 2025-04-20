package pvz_proto

import (
	"context"
	"net"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"

	"trainee-pvz/internal/repository"
)

type PVZGRPCServer struct {
	UnimplementedPVZServiceServer
	repo repository.PVZRepository
}

func NewPVZGRPCServer(repo repository.PVZRepository) PVZServiceServer {
	return &PVZGRPCServer{repo: repo}
}

func (s *PVZGRPCServer) GetPVZList(ctx context.Context, req *GetPVZListRequest) (*GetPVZListResponse, error) {
	data, err := s.repo.ListWithReceptionsAndProducts(ctx, nil, nil, 1, 100)
	if err != nil {
		return nil, err
	}

	var resp GetPVZListResponse
	for _, item := range data {
		resp.Pvzs = append(resp.Pvzs, &PVZ{
			Id:               item.PVZ.ID,
			City:             item.PVZ.City,
			RegistrationDate: timestamppb.New(item.PVZ.RegistrationDate),
		})
	}

	return &resp, nil
}

func StartGRPCServer(repo *repository.PVZRepository, port string) error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return errors.Wrap(err, "can't listen port")
	}

	s := grpc.NewServer()
	RegisterPVZServiceServer(s, NewPVZGRPCServer(*repo))
	reflection.Register(s)

	return s.Serve(lis)
}

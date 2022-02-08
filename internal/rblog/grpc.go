package rblog

import (
	"google.golang.org/grpc"
	"log"
	"net"
)

type grpcServe struct {
	address string
	*grpc.Server
}

func (s grpcServe) Run() {
	listen, err := net.Listen("tcp", s.address)

	if err != nil {
		log.Fatalln(err)
	}

	go func() {
		if err := s.Serve(listen); err != nil {
			log.Fatalln(err)
		}

		s.GracefulStop()
	}()

	log.Printf("grpc serve is running port: %s", s.address)
}

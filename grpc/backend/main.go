package main

import (
	"flag"
	"fmt"
	"github.com/Sirupsen/logrus"
	pb "github.com/arjanvaneersel/grpc/api"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"net"
	"io/ioutil"
	"github.com/arjanvaneersel/flite"
)

func main() {
	port := flag.Int("p", 8080, "port to listen to")
	flag.Parse()

	logrus.Infof("Listening to port %d", *port)

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		logrus.Fatalf("Could not listen to port %d: %v", *port, err)
	}

	s := grpc.NewServer()
	pb.RegisterTextToSpeechServer(s, server{})
	if err := s.Serve(l); err != nil {
		logrus.Fatalf("Could not serve: %v", err)
	}
}

type server struct{}

func (server) Say(ctx context.Context, text *pb.Text) (*pb.Speech, error) {
	f, err := ioutil.TempFile("", "")
	if err != nil {
		return nil, fmt.Errorf("Error creating temp file: %v", err)
	}
	if err := f.Close(); err != nil {
		return nil, fmt.Errorf("Error closing temp file %s: %v", f.Name(), err)
	}

	if err := flite.TextToSpeech( f.Name(), text.Text); err != nil {
		return nil, fmt.Errorf("Flite failed: %v", err)
	}

	data, err := ioutil.ReadFile(f.Name())
	if err != nil {
		return nil, fmt.Errorf("Error reading temp file: %v", err)
	}
	return &pb.Speech{Audio: data}, nil
}

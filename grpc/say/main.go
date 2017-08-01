package main

import (
	"flag"
	"log"
	pb "github.com/arjanvaneersel/grpc/api"
	"google.golang.org/grpc"
	"golang.org/x/net/context"
	"io/ioutil"
	"os"
	"fmt"
)

func main() {
	backend := flag.String("b", "localhost:8080", "address of backend")
	output := flag.String("o", "output.wav", "output file")
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Printf("usage:\n\t%s \"text to speak\"", os.Args[0])
		os.Exit(1)
	}

	conn, err := grpc.Dial(*backend, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Couldn't connect to %s: %v", *backend, err)
	}
	defer conn.Close()

	client := pb.NewTextToSpeechClient(conn)
	txt := &pb.Text{Text: flag.Arg(0)}
	res, err := client.Say(context.Background(), txt)
	if err != nil {
		log.Fatalf("Couldn't say %s: %v", txt, err)
	}
	if err := ioutil.WriteFile(*output, res.Audio, 0666); err != nil {
		log.Fatalf("Couldn't write file %s: %v", *output, err)
	}
}

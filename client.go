package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"time"

	// The Protobuf generated file
	creator "github.com/appscode/haproxy-grpc-sample/codenamecreator"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

func main() {
	address := flag.String("address", "", "server address")
	crt := flag.String("crt", "", "path to cert file")
	flag.Parse()

	log.Println(*address, *crt)

	option := grpc.WithInsecure()
	if len(*crt) > 0 {
		creds, err := credentials.NewClientTLSFromFile(*crt, "")
		if err != nil {
			log.Fatalf("Failed to load TLS certificate")
		}
		option = grpc.WithTransportCredentials(creds)
	}

	conn, err := grpc.Dial(*address, option)
	if err != nil {
		log.Fatalf("Did not connect, %v", err)
	}
	defer conn.Close()

	client := creator.NewCodenameCreatorClient(conn)
	ctx := context.Background()

	// Optional: add some metadata
	ctx = metadata.AppendToOutgoingContext(ctx, "mysecretpassphrase", "abc123")

	getCodenamesStreamingExample(ctx, client)
	// getSingleCodenameAndExitExample(ctx, client, "Science")
}

func getSingleCodenameAndExitExample(ctx context.Context, client creator.CodenameCreatorClient, category string) {
	result, err := client.GetCodename(ctx, &creator.NameRequest{Category: category})
	if err != nil {
		log.Fatalf("Could not get result, %v", err)
	}

	log.Printf("Codename result: %s", result)
}

func getCodenamesStreamingExample(ctx context.Context, client creator.CodenameCreatorClient) {
	fmt.Println("Generating codenames...")
	stream, err := client.KeepGettingCodenames(ctx)

	if err != nil {
		log.Fatalf("Could not get stream, %v", err)
	}

	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				break
			}

			if err != nil {
				log.Fatalf("%v", err)
			}

			log.Printf("Received: %s\n", in.Name)
		}
	}()

	category := "Science"
	for {
		if category == "Science" {
			category = "Animals"
		} else {
			category = "Science"
		}

		err := stream.Send(&creator.NameRequest{Category: category})
		if err != nil {
			log.Fatalf("%v", err)
		}
		time.Sleep(10 * time.Second)
	}
}

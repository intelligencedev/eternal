package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"strings"

	pb "mlxeternal/pkg/eproto"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewTextGeneratorClient(conn)

	stream, err := c.GenerateTextStream(context.Background(), &pb.TextRequest{Prompt: "write a hello world in golang."})
	if err != nil {
		log.Fatalf("could not initiate stream: %v", err)
	}

	var responseBuilder strings.Builder
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			// Flush the remaining content in buffer if any
			if responseBuilder.Len() > 0 {
				printResponse(&responseBuilder)
			}
			break
		}
		if err != nil {
			log.Fatalf("error while receiving: %v", err)
		}

		// Append the received text to the buffer
		responseBuilder.WriteString(resp.Response)

		printResponse(&responseBuilder)
	}
}

// printResponse processes the buffer, prints complete messages, and retains any incomplete message
func printResponse(responseBuilder *strings.Builder) {
	content := responseBuilder.String()
	messages := strings.Split(content, "\n")
	for i, msg := range messages {
		if i == len(messages)-1 && !strings.HasSuffix(content, "\n") {
			// If the last part does not end with a newline, it's an incomplete message
			responseBuilder.Reset()
			responseBuilder.WriteString(msg)
		} else if len(msg) > 0 {
			// Print the complete message
			fmt.Printf("%s\n", msg)
		}
	}
	if strings.HasSuffix(content, "\n") {
		// If the last part ends with a newline, clear the buffer
		responseBuilder.Reset()
	}
}

package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"grpchello/proto"
	"log"
	"os"
)

var (
	tls                = flag.Bool("tls", true, "Connection uses TLS if true, else plain TCP")
	caFile             = flag.String("ca_file", "", "The file containing the CA root cert file")
	serverAddr         = flag.String("addr", "grpctest2.bitautotech.com:443", "The server address in the format of host:port")
	serverHostOverride = flag.String("server_host_override", "grpctest2.bitautotech.com", "The server name used to verify the hostname returned by the TLS handshake")
)

func main() {
	//var serviceHost = "grpctest2.bitautotech.com:9081"
	flag.Parse()
	var opts []grpc.DialOption
	if *tls {
		if *caFile == "" {
			*caFile = Path("x509/bitautotech.pem")
		}
		creds, err := credentials.NewClientTLSFromFile(*caFile, *serverHostOverride)
		if err != nil {
			log.Fatalf("Failed to create TLS credentials %v", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()

	client := proto.NewEchoClient(conn)
	rsp, err := client.Ping(context.TODO(), &proto.Empty{})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(rsp)
	rsp1, err := client.Reverse(context.TODO(), &proto.Content{
		Text: "ssshhh",
	})
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(rsp1)

	fmt.Println("按回车键退出程序...")
	in := bufio.NewReader(os.Stdin)
	_, _, _ = in.ReadLine()
}

package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/opentracing/opentracing-go"

	"strconv"

	zipkin "github.com/openzipkin/zipkin-go-opentracing"

	"grpclb/services/server1/pb"

	grpclb "grpclb/etcdv3"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var (
	serv = flag.String("service", "hello_service", "service name")
	reg  = flag.String("reg", "http://127.0.0.1:2379", "register etcd address")
)

const (
	// Our service name.
	serviceName = "zipkinclient1"

	// Host + port of our service.
	hostPort = "127.0.0.1:50000"

	// Endpoint to send Zipkin spans to.
	zipkinHTTPEndpoint = "http://localhost:9411/api/v1/spans"

	// Debug mode.
	debug = false

	// Base endpoint of our SVC1 service.
	svc1Endpoint = "http://localhost:61001"

	// same span can be set to true for RPC style spans (Zipkin V1) vs Node style (OpenTracing)
	sameSpan = true
)

func main() {
	flag.Parse()

	// Create our HTTP collector.
	collector, err := zipkin.NewHTTPCollector(zipkinHTTPEndpoint)
	if err != nil {
		fmt.Printf("unable to create Zipkin HTTP collector: %+v", err)
		os.Exit(-1)
	}

	// Create our recorder.
	recorder := zipkin.NewRecorder(collector, debug, hostPort, serviceName)

	// Create our tracer.
	tracer, err := zipkin.NewTracer(
		recorder,
	)
	if err != nil {
		fmt.Printf("unable to create Zipkin tracer: %+v", err)
		os.Exit(-1)
	}

	// Explicitely set our tracer to be the default tracer.
	opentracing.InitGlobalTracer(tracer)

	r := grpclb.NewResolver(*serv)
	b := grpc.RoundRobin(r)

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	// Create Root Span for duration of the interaction with svc1

	conn, err := grpc.DialContext(ctx, *reg, grpc.WithInsecure(), grpc.WithBalancer(b))
	if err != nil {
		panic(err)
	}
	client := pb.NewGreeterClient(tracer, conn)

	ticker := time.NewTicker(time.Second)
	for t := range ticker.C {

		resp, err := client.SayHello(ctx, &pb.HelloRequest{Name: "world " + strconv.Itoa(t.Second())})
		if err == nil {
			fmt.Printf("%v: Reply is %s\n", t, resp.Message)
		}
	}
}

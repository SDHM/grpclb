package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	grpclb "grpclb/etcdv3"
	"grpclb/tracing"

	"grpclb/services/server2/pb"

	zipkin "github.com/openzipkin/zipkin-go-opentracing"
)

var (
	serv = flag.String("servicea", "serverA", "servicea name")

	port = flag.Int("port", 50003, "listening port")
	reg  = flag.String("reg", "http://127.0.0.1:2379", "register etcd address")
)

type tracingType struct{}

func main() {
	flag.Parse()

	zipkinHTTPEndpoint := "http://localhost:9411/api/v1/spans"
	hostname, _ := os.Hostname()
	fmt.Println("hostname:", hostname)
	InitTracer(zipkinHTTPEndpoint, hostname, "servera")

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", *port))
	if err != nil {
		panic(err)
	}

	err = grpclb.Register(*serv, "127.0.0.1", *port, *reg, time.Second*5, 15)
	if err != nil {
		panic(err)
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)
	go func() {
		s := <-ch
		log.Printf("receive signal '%v'", s)
		grpclb.UnRegister()
		os.Exit(1)
	}()

	log.Printf("starting hello service at %d", *port)

	var interceptor grpc.UnaryServerInterceptor
	interceptor = func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		fmt.Println("\n\nserver method called begin:", info.FullMethod)

		funcTrac := tracing.FromGRPCRequest(opentracing.GlobalTracer(), "service "+info.FullMethod)
		// Retrieve gRPC metadata.
		md, ok := metadata.FromContext(ctx)
		if !ok {
			md = metadata.MD{}
		}

		ctx = funcTrac(ctx, &md)

		span := opentracing.SpanFromContext(ctx)
		ext.Component.Set(span, "servera")
		ext.PeerService.Set(span, "servera")
		// fmt.Println("span:", span)
		defer span.Finish()

		// 这里进行traceing
		if err != nil {
			return
		}
		defer fmt.Printf("server method called end:%s\n\n\n", info.FullMethod)
		// 继续处理请求
		resp, err = handler(ctx, req)

		return
	}

	s := grpc.NewServer(grpc.UnaryInterceptor(interceptor))
	// s := grpc.NewServer()
	pb.RegisterServerAServer(s, &serverA{})
	s.Serve(lis)
}

// server is used to implement helloworld.GreeterServer.
type serverA struct{}

// SayHello implements helloworld.GreeterServer
func (s *serverA) ServerAFunc(ctx context.Context, in *pb.ServierARequest) (*pb.ServierAReply, error) {
	fmt.Printf("%v: Receive is %s\n", time.Now(), in.Name)
	time.Sleep(time.Millisecond * 10)
	return &pb.ServierAReply{Message: "Hello " + in.Name}, nil
}

func InitTracer(zipkinURL string, hostPort string, serviceName string) {
	collector, err := zipkin.NewHTTPCollector(zipkinURL)
	if err != nil {
		log.Fatalf("unable to create Zipkin HTTP collector: %v", err)
		return
	}
	fmt.Println("serviceName:", serviceName)
	tracer, err := zipkin.NewTracer(
		zipkin.NewRecorder(collector, false, "hostPort", serviceName),
	)
	if err != nil {
		log.Fatalf("unable to create Zipkin tracer: %v", err)
		return
	}
	opentracing.InitGlobalTracer(tracer)
	return
}

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

	"grpclb/services/server1/pb"
	s2pb "grpclb/services/server2/pb"
	s3pb "grpclb/services/server3/pb"

	zipkin "github.com/openzipkin/zipkin-go-opentracing"
)

var (
	serv  = flag.String("service", "hello_service", "service name")
	serva = flag.String("servicea", "serverA", "servicea name")
	servb = flag.String("serviceb", "serverB", "serviceb name")
	port  = flag.Int("port", 50001, "listening port")
	reg   = flag.String("reg", "http://127.0.0.1:2379", "register etcd address")

	client2 s2pb.ServerAClient
	client3 s3pb.ServerBClient
)

type tracingType struct{}

func main() {
	flag.Parse()

	zipkinHTTPEndpoint := "http://localhost:9411/api/v1/spans"
	hostname, _ := os.Hostname()
	InitTracer(zipkinHTTPEndpoint, hostname, "server")

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
		ext.Component.Set(span, "server")
		ext.SpanKind.Set(span, "resource")
		ext.PeerService.Set(span, "server")
		ext.PeerHostname.Set(span, "localhost")
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

	conn3, err := grpc.Dial("localhost:50005", grpc.WithInsecure(), grpc.WithBlock(), grpc.WithBackoffConfig(grpc.BackoffConfig{
		MaxDelay: 10 * time.Second,
	}), grpc.WithTimeout(time.Minute))

	client3 = s3pb.NewServerBClient(opentracing.GlobalTracer(), conn3)

	conn2, err := grpc.Dial("localhost:50003", grpc.WithInsecure(), grpc.WithBlock(), grpc.WithBackoffConfig(grpc.BackoffConfig{
		MaxDelay: 10 * time.Second,
	}), grpc.WithTimeout(time.Minute))

	client2 = s2pb.NewServerAClient(opentracing.GlobalTracer(), conn2)

	s := grpc.NewServer(grpc.UnaryInterceptor(interceptor))
	// s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	s.Serve(lis)
}

// server is used to implement helloworld.GreeterServer.
type server struct{}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	fmt.Printf("%v: Receive is %s\n", time.Now(), in.Name)
	time.Sleep(time.Millisecond * 10)

	s.CallServerA(opentracing.GlobalTracer(), ctx)

	s.CallServerB(opentracing.GlobalTracer(), ctx)
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}

func InitTracer(zipkinURL string, hostPort string, serviceName string) {
	collector, err := zipkin.NewHTTPCollector(zipkinURL)
	if err != nil {
		log.Fatalf("unable to create Zipkin HTTP collector: %v", err)
		return
	}
	fmt.Println("serviceName:", serviceName)
	tracer, err := zipkin.NewTracer(
		zipkin.NewRecorder(collector, false, hostPort, serviceName),
	)
	if err != nil {
		log.Fatalf("unable to create Zipkin tracer: %v", err)
		return
	}
	opentracing.InitGlobalTracer(tracer)
	return
}

func (s *server) CallServerA(tracer opentracing.Tracer, ctx context.Context) {

	fmt.Println("callservera begin")
	resp, err := client2.ServerAFunc(ctx, &s2pb.ServierARequest{Name: "world "})
	fmt.Println("callservera end")
	if err == nil {
		fmt.Printf(" Reply is %s\n", resp.Message)
	}
}

func (s *server) CallServerB(tracer opentracing.Tracer, ctx context.Context) {

	fmt.Println("callserverb begin")
	resp, err := client3.ServerBFunc(ctx, &s3pb.ServierBRequest{Name: "world "})
	fmt.Println("callserverb end")
	if err == nil {
		fmt.Printf(" Reply is %s\n", resp.Message)
	}
}

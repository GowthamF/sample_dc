package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	pb "dc.com/m/v2/protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func buildConfig(kubeconfig string) (*rest.Config, error) {
	if kubeconfig != "" {
		cfg, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, err
		}
		return cfg, nil
	}

	cfg, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

type server struct {
	pb.UnimplementedGreeterServer
}

func (s *server) SayHelloAgain(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	var message string = "Hello" + in.GetName()
	return &pb.HelloReply{Message: &message}, nil
}

var (
	port     = flag.Int("port", 50051, "The server port")
	addr     = flag.String("addr", "localhost:50051", "the address to connect to")
	name     = flag.String("name", "world", "Name to greet")
	isServer = flag.Bool("server", false, "True or False")
)

func getServerMessage() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r1, err1 := c.SayHelloAgain(ctx, &pb.HelloRequest{Name: name})

	if err1 != nil {
		log.Fatalf("could not greet: %v", err1)
	}
	log.Printf("Greeting: %s", r1.GetMessage())
}

func createServer() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func main() {
	flag.Parse()
	// if !*isServer {
	createServer()
	// }
	// if *isServer {
	getServerMessage()
	// }

	// klog.InitFlags(nil)
	// lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	// var kubeconfig string
	// var leaseLockName string
	// var leaseLockNamespace string
	// var id string

	// flag.StringVar(&kubeconfig, "kubeconfig", "", "absolute path to the kubeconfig file")
	// flag.StringVar(&id, "id", uuid.New().String(), "the holder identity name")
	// flag.StringVar(&leaseLockName, "lease-lock-name", "", "the lease lock resource name")
	// flag.StringVar(&leaseLockNamespace, "lease-lock-namespace", "", "the lease lock resource namespace")
	// flag.Parse()

	// if leaseLockName == "" {
	// 	klog.Fatal("unable to get lease lock resource name (missing lease-lock-name flag).")
	// }
	// if leaseLockNamespace == "" {
	// 	klog.Fatal("unable to get lease lock resource namespace (missing lease-lock-namespace flag).")
	// }

	// // leader election uses the Kubernetes API by writing to a
	// // lock object, which can be a LeaseLock object (preferred),
	// // a ConfigMap, or an Endpoints (deprecated) object.
	// // Conflicting writes are detected and each client handles those actions
	// // independently.

	// s := grpc.NewServer()

	// pb.RegisterGreeterServer(s, &server{})

	// if err := s.Serve(lis); err != nil {
	// 	log.Fatalf("failed to serve: %v", err)
	// }

	// config, err := buildConfig(kubeconfig)
	// if err != nil {
	// 	klog.Fatal(err)
	// }
	// client := clientset.NewForConfigOrDie(config)

	// run := func(ctx context.Context) {
	// 	// complete your controller loop here
	// 	klog.Info("Controller loop...")

	// 	select {}
	// }

	// // use a Go context so we can tell the leaderelection code when we
	// // want to step down
	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()

	// // listen for interrupts or the Linux SIGTERM signal and cancel
	// // our context, which the leader election code will observe and
	// // step down
	// ch := make(chan os.Signal, 1)
	// signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	// go func() {
	// 	<-ch
	// 	klog.Info("Received termination, signaling shutdown")
	// 	cancel()
	// }()

	// // we use the Lease lock type since edits to Leases are less common
	// // and fewer objects in the cluster watch "all Leases".
	// lock := &resourcelock.LeaseLock{
	// 	LeaseMeta: metav1.ObjectMeta{
	// 		Name:      leaseLockName,
	// 		Namespace: leaseLockNamespace,
	// 	},
	// 	Client: client.CoordinationV1(),
	// 	LockConfig: resourcelock.ResourceLockConfig{
	// 		Identity: id,
	// 	},
	// }

	// // start the leader election code loop
	// leaderelection.RunOrDie(ctx, leaderelection.LeaderElectionConfig{
	// 	Lock: lock,
	// 	// IMPORTANT: you MUST ensure that any code you have that
	// 	// is protected by the lease must terminate **before**
	// 	// you call cancel. Otherwise, you could have a background
	// 	// loop still running and another process could
	// 	// get elected before your background loop finished, violating
	// 	// the stated goal of the lease.
	// 	ReleaseOnCancel: true,
	// 	LeaseDuration:   60 * time.Second,
	// 	RenewDeadline:   15 * time.Second,
	// 	RetryPeriod:     5 * time.Second,
	// 	Callbacks: leaderelection.LeaderCallbacks{
	// 		OnStartedLeading: func(ctx context.Context) {
	// 			// we're notified when we start - this is where you would
	// 			// usually put your code
	// 			run(ctx)
	// 		},
	// 		OnStoppedLeading: func() {
	// 			// we can do cleanup here
	// 			klog.Infof("leader lost: %s", id)
	// 			os.Exit(0)
	// 		},
	// 		OnNewLeader: func(identity string) {
	// 			// we're notified when new leader elected
	// 			if identity == id {
	// 				// I just got the lock
	// 				return
	// 			}
	// 			klog.Infof("new leader elected: %s", identity)
	// 		},
	// 	},
	// })
}

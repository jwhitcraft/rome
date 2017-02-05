// Copyright © 2017 Jon Whitcraft
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"log"
	"net"

	"fmt"
	"os"
	"os/signal"
	"syscall"

	pb "github.com/jwhitcraft/rome/cesar"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50051"
)

var (
	attributes *pb.BuildAttrResponse
)

type server struct{}

func (s *server) BuildFile(ctx context.Context, in *pb.FileRequest) (*pb.FileResponse, error) {
	return &pb.FileResponse{File: "/tmp/" + in.Path, Success: true}, nil
}

func (s *server) CreateSymLink(ctx context.Context, in *pb.CreateSymLinkRequest) (*pb.FileResponse, error) {
	return &pb.FileResponse{File: "/tmp/" + in.OrginalFile, Success: true}, nil
}

func (s *server) SetBuildAttributes(ctx context.Context, in *pb.SetBuildAttrRequest) (*pb.BuildAttrResponse, error) {
	attributes = &pb.BuildAttrResponse{Version: in.Version, Clean: in.Clean, Flavor: in.Flavor, Folder: in.Folder}

	return attributes, nil
}

func (s *server) GetBuildAttributes(ctx context.Context, in *pb.GetBuildAttrRequest) (*pb.BuildAttrResponse, error) {
	if attributes == nil {
		attributes = &pb.BuildAttrResponse{}
	}

	return attributes, nil
}

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run Rome as a service on a remote machine",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		errc := make(chan error)

		// Interrupt handler
		go func() {
			c := make(chan os.Signal, 1)
			signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
			errc <- fmt.Errorf("%s", <-c)
		}()

		// gRPC transport
		go func() {
			lis, err := net.Listen("tcp", port)
			if err != nil {
				log.Fatalf("failed to listen: %v", err)
			}
			s := grpc.NewServer()
			pb.RegisterCesarServer(s, &server{})
			// Register reflection service on gRPC server.
			reflection.Register(s)
			errc <- s.Serve(lis)
		}()

		<-errc
	},
}

func init() {
	RootCmd.AddCommand(serverCmd)
}

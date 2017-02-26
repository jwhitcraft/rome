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
	"fmt"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/go-kit/kit/log"
	pb "github.com/jwhitcraft/rome/aqueduct"
	"github.com/jwhitcraft/rome/build"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":47600"
)

var (
	attributes  *pb.BuildAttrResponse
	buildRoot   string
	buildFolder string
	logger      log.Logger
	logfile     string
)

type server struct{}

func (s *server) BuildFile(ctx context.Context, in *pb.FileRequest) (*pb.FileResponse, error) {
	target := filepath.Join(buildFolder, in.Target)
	logger.Log("msg", "Building File"+target)
	file := build.CreateRemoteFile(target, in.Contents)
	err := file.Process(attributes.Flavor, attributes.Version)
	if err != nil {
		return nil, err
	}

	return &pb.FileResponse{File: target}, nil
}

func (s *server) CreateSymLink(ctx context.Context, in *pb.CreateSymLinkRequest) (*pb.FileResponse, error) {
	target := filepath.Join(buildRoot, attributes.Folder, in.Target)
	logger.Log("msg", fmt.Sprintf("Symlinking %s to %s", in.OriginFile, target))
	file := build.CreateSymLink(target, in.OriginFile)
	err := file.Process(attributes.Flavor, attributes.Version)
	if err != nil {
		return nil, err
	}
	return &pb.FileResponse{File: target}, nil
}

func (s *server) SetBuildAttributes(ctx context.Context, in *pb.SetBuildAttrRequest) (*pb.BuildAttrResponse, error) {
	attributes = &pb.BuildAttrResponse{Version: in.Version, Clean: in.Clean, Flavor: in.Flavor, Folder: in.Folder}

	buildFolder = filepath.Join(buildRoot, attributes.Folder)
	if attributes.Clean {
		logger.Log("msg", fmt.Sprintf("Cleaning %s", buildFolder))
		build.CleanBuild(build.TargetDirectory{Path: buildFolder})
	}

	return attributes, nil
}

func (s *server) CleanCache(ctx context.Context, in *pb.CleanCacheRequest) (*pb.CleanCacheResponse, error) {

	err := build.CleanCache(buildFolder, cleanCacheItems)
	if err != nil {
		return nil, err
	}

	return &pb.CleanCacheResponse{}, nil
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
	Long: `Running Rome as a service, allows building of files on a remote host, when the files are on a different
machine.   This allows code to live locally, but be built in a VM or a Container where this service is running.`,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		logOutput := os.Stdout
		if logfile != "" {
			logOutput, err = os.Create(logfile)
			if err == nil {
				defer logOutput.Close()
			}
		}
		logger = log.NewLogfmtLogger(logOutput)
		logger = log.NewContext(logger).With("ts", log.DefaultTimestampUTC)
		logger.Log("msg", fmt.Sprintf("Start %s", cmd.Short))
		defer logger.Log("msg", "Shutting down server")
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
				logger.Log("failed to listen: %v", err)
			}
			s := grpc.NewServer(
				grpc.MaxMsgSize(1024 * 1024 * 50),
			)
			pb.RegisterAqueductServer(s, &server{})
			// Register reflection service on gRPC server.
			reflection.Register(s)
			errc <- s.Serve(lis)
		}()

		<-errc
	},
}

func init() {
	RootCmd.AddCommand(serverCmd)

	serverCmd.Flags().StringVarP(
		&buildRoot,
		"build-root",
		"r",
		"/var/www/html",
		"What is the default root, for to build the files at on the remote server",
	)

	serverCmd.Flags().StringVar(
		&logfile,
		"logfile",
		"",
		"Send the output to a log file, if nothing is passed, it will default to standard out",
	)
}

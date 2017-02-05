package build

import (
	"os"
	"path"

	pb "github.com/jwhitcraft/rome/cesar"
	"golang.org/x/net/context"
)

type SymLink struct {
	Target     string
	OriginFile string
}

func CreateSymLink(target, origin string) *SymLink {
	return &SymLink{Target: target, OriginFile: origin}
}

func (s *SymLink) Process(version, flavor string) error {
	// make sure the folder exists first
	os.MkdirAll(path.Dir(s.Target), 0775)
	// link it up!
	return os.Symlink(s.OriginFile, s.Target)
}

func (s *SymLink) GetTarget() string {
	return s.Target
}

func (s *SymLink) SendToCesar(cesar pb.CesarClient) (*pb.FileResponse, error) {
	return cesar.CreateSymLink(context.Background(), &pb.CreateSymLinkRequest{
		Name:        s.Target,
		OrginalFile: s.OriginFile,
	})
}
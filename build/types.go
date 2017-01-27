package build

type Directory interface {
	Clean(path string) error
	Dir() string
}

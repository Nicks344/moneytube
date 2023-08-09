package serverutils

import (
	"net/http"
	"strings"
)

type fileSystem struct {
	fs http.FileSystem
}

// Open opens file
func (fs fileSystem) Open(path string) (http.File, error) {
	f, err := fs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if s.IsDir() {
		index := strings.TrimSuffix(path, "/") + "/index.html"
		if _, err := fs.fs.Open(index); err != nil {
			return nil, err
		}
	}

	return f, nil
}

func FileServerWrapperHandler(path string, filePath string) http.Handler {
	return http.StripPrefix(path, http.FileServer(http.Dir(filePath)))
}

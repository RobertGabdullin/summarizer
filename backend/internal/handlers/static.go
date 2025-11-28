package handlers

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const distDir = "../frontend/dist"

type Static struct {
	fileHandler http.Handler
}

func NewStatic() *Static {
	return &Static{
		fileHandler: http.FileServer(http.Dir(distDir)),
	}
}

func (s *Static) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	upath := r.URL.Path

	fpath := filepath.Join(distDir, filepath.FromSlash(upath))
	info, err := os.Stat(fpath)
	if err == nil && !info.IsDir() {
		if strings.HasSuffix(fpath, "index.html") {
			w.Header().Set("Cache-Control", "no-cache")
		} else {
			w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		}
		s.fileHandler.ServeHTTP(w, r)
		return
	}

	index := filepath.Join(distDir, "index.html")
	http.ServeFile(w, r, index)
}

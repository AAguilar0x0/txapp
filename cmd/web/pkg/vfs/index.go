package vfs

import (
	"crypto/sha256"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"path/filepath"
	"strconv"
)

func calculateChecksum(fsys fs.FS, path string) (string, error) {
	f, err := fsys.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

type VersionedFileServer struct {
	fs          fs.FS
	urlPath     string
	checksums   map[string]string
	fileHandler http.Handler
	maxAge      uint
}

func New(fileSystem fs.FS, urlPath string, maxAgeInSeconds uint) (*VersionedFileServer, error) {
	vfs := &VersionedFileServer{
		fs:          fileSystem,
		urlPath:     urlPath,
		checksums:   make(map[string]string),
		fileHandler: http.FileServer(http.FS(fileSystem)),
		maxAge:      maxAgeInSeconds,
	}

	err := fs.WalkDir(fileSystem, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}

		baseFile := filepath.Base(path)

		if checksum, err := calculateChecksum(fileSystem, path); err == nil {
			vfs.checksums[baseFile] = checksum
		} else {
			vfs.checksums[baseFile] = baseFile
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return vfs, nil
}

func (d *VersionedFileServer) filesVersion(file string) string {
	baseFile := filepath.Base(file)
	if checksum, exists := d.checksums[baseFile]; exists {
		return checksum
	}
	return ""
}

func (d *VersionedFileServer) URLWithVersion(file string) string {
	fileURL := d.urlPath + "/" + file
	if checksum, exists := d.checksums[file]; exists {
		return fileURL + "?v=" + checksum
	}
	return fileURL
}

func (d *VersionedFileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	etag := d.filesVersion(r.URL.Path)
	if etag == "" {
		http.NotFound(w, r)
		return
	}

	if match := r.Header.Get("If-None-Match"); match != "" && match == etag {
		w.WriteHeader(http.StatusNotModified)
		return
	}

	w.Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%s", strconv.Itoa(int(d.maxAge))))
	w.Header().Set("ETag", etag)

	d.fileHandler.ServeHTTP(w, r)
}

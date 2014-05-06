package pack

import (
	"archive/zip"
	"bytes"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func sanitizedName(filename string) string {
	if len(filename) > 1 && filename[1] == ':' &&
		runtime.GOOS == "windows" {
		filename = filename[2:]
	}
	filename = filepath.ToSlash(filename)
	filename = filepath.Clean(filename)
	return strings.Replace(filename, "../", "", -1)
}

type Archiever interface {
	Add(filename string) error
	Close() error
}

type Zip struct {
	*zip.Writer
}

func CreateZip(filename string) (*Zip, error) {
	fd, err := os.Create(filename)
	if err != nil {
		return nil, err
	}
	zipper := zip.NewWriter(fd)
	return &Zip{Writer: zipper}, nil
}

func (z *Zip) Add(filename string) error {
	info, err := os.Lstat(filename)
	if err != nil {
		return err
	}
	// content
	var reader io.Reader
	if info.Mode()&os.ModeSymlink != 0 {
		target, err := os.Readlink(filename)
		if err != nil {
			return err
		}
		reader = bytes.NewBuffer([]byte(target))
	} else {
		file, err := os.Open(filename)
		if err != nil {
			return err
		}
		defer file.Close()
		reader = file
	}
	// header
	hdr, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}
	hdr.Name = sanitizedName(filename)
	if info.IsDir() {
		hdr.Name += "/"
	}
	hdr.Method = zip.Deflate // compress method
	writer, err := z.CreateHeader(hdr)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		_, err = io.Copy(writer, reader)
		return err
	}
	return nil
}

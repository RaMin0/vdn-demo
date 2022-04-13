package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
)

type localLister struct{ path string }

func (l *localLister) List() (es []string) {
	var (
		idRegexp = regexp.MustCompile("(S\\d{2}E\\d{2,3})")
	)

	fs, err := ioutil.ReadDir(l.path)
	if err != nil {
		return nil
	}
	for _, f := range fs {
		es = append(es, idRegexp.FindString(f.Name()))
	}
	return es
}

type localUploader struct{ path string }

func (u *localUploader) Upload(id string) {
	src, err := os.Open(filepath.Join(".tmp", id))
	if err != nil {
		return
	}
	defer src.Close()
	dst, err := os.Create(filepath.Join(u.path, fmt.Sprintf("%s.mp4", id)))
	if err != nil {
		return
	}
	defer dst.Close()

	io.Copy(dst, src)
}

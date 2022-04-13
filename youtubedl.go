package main

import (
	"os/exec"
	"path/filepath"
)

type youtubedlDownloader struct{ extraArgs []string }

func (d *youtubedlDownloader) Download(id, url string) {
	exec.Command("youtube-dl", append(append(append([]string{},
		"--ignore-errors",
		"--merge-output-format", "mp4",
		"--output", filepath.Join(".tmp", id),
	), d.extraArgs...), url)...).Run()
}

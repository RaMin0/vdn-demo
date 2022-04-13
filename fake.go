package main

import (
	"fmt"
	"time"
)

type fakeLister struct{ n int }

func (l *fakeLister) List() []string {
	a := make([]string, l.n)
	for i := range a {
		a[i] = fmt.Sprintf("E%02d", i+1)
	}
	return a
}

type fakeFetcher struct{ n int }

func (f *fakeFetcher) Fetch() <-chan FetcherEpisode {
	ch := make(chan FetcherEpisode)
	go func() {
		defer close(ch)
		for i := 0; i < f.n; i++ {
			time.Sleep(100)
			ch <- FetcherEpisode{
				ID:  fmt.Sprintf("E%02d", i+1),
				URL: fmt.Sprintf("https://example.com/E%02d.mp4", i+1),
			}
		}
	}()
	return ch
}

type fakeDownloader struct{}

func (*fakeDownloader) Download(string, string) {
	time.Sleep(1 * time.Second)
}

type fakeUploader struct{}

func (*fakeUploader) Upload(string) {
	time.Sleep(2 * time.Second)
}

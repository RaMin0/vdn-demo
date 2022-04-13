package main

import (
	"sync"

	"github.com/fatih/color"
)

type Series struct {
	Name string
	Lister
	Fetcher
	Downloader
	Uploader
}

type Lister interface {
	List() (localEpisodes []string)
}

type Fetcher interface {
	Fetch() (remoteEpisodesCh <-chan FetcherEpisode)
}
type FetcherEpisode struct {
	ID  string
	URL string
}

type Downloader interface {
	Download(id, remoteURL string)
}

type Uploader interface {
	Upload(destinationPath string)
}

func main() {
	s := Series{
		Name:       "fake",
		Lister:     &fakeLister{2},
		Fetcher:    &fakeFetcher{9},
		Downloader: &fakeDownloader{},
		Uploader:   &fakeUploader{},

		// Name:       "local+okru+youtubedl",
		// Lister:     &localLister{"dl"},
		// Fetcher:    &okruFetcher{"TODO: search query"},
		// Downloader: &youtubedlDownloader{},
		// Uploader:   &localUploader{"dl"},
	}

	// List
	log.Printf("%v %q", color.WhiteString("Listing"), s.Name)
	localEpisodes := list(s)

	// Fetch
	log.Printf("%v %q", color.BlueString("Fetching"), s.Name)
	remoteEpisodes := fetch(s)

	// Download
	downloadedEpisodes := download(s, localEpisodes, remoteEpisodes, 3)

	// Upload
	upload(s, downloadedEpisodes)
}

func list(s Lister) map[string]bool {
	es := s.List()
	esMap := map[string]bool{}
	for _, id := range es {
		esMap[id] = true
	}
	return esMap
}
func fetch(s Fetcher) <-chan FetcherEpisode {
	ch := make(chan FetcherEpisode)
	go func() {
		defer close(ch)
		for e := range s.Fetch() {
			log.Printf("%v %s", color.BlueString("Fetched"), e.ID)
			ch <- e
		}
	}()
	return ch
}
func download(s Downloader, localEpisodesMap map[string]bool, remoteEpisodes <-chan FetcherEpisode, limit int) <-chan string {
	ch := make(chan string)
	go func() {
		defer close(ch)

		var wg sync.WaitGroup
		defer wg.Wait()

		guard := make(chan bool, limit)
		for e := range remoteEpisodes {
			if _, ok := localEpisodesMap[e.ID]; ok {
				continue
			}

			wg.Add(1)

			go func(e FetcherEpisode) {
				defer wg.Done()

				log.Printf("%v %s", color.MagentaString("Queueing"), e.ID)

				guard <- true
				log.Printf("%v %s", color.YellowString("Downloading"), e.ID)
				s.Download(e.ID, e.URL)
				log.Printf("%v %s", color.YellowString("Downloaded"), e.ID)
				<-guard

				ch <- e.ID
			}(e)
		}
	}()
	return ch
}
func upload(s Uploader, downloadedEpisodes <-chan string) {
	for id := range downloadedEpisodes {
		log.Printf("%v %s", color.GreenString("Uploading"), id)
		s.Upload(id)
		log.Printf("%v %s", color.GreenString("Uploaded"), id)
	}
}

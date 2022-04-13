package main

import (
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type okruFetcher struct{ query string }

func (f *okruFetcher) Fetch() <-chan FetcherEpisode {
	ch := make(chan FetcherEpisode)
	go func() {
		defer close(ch)

		var (
			idRegexp = regexp.MustCompile("(S\\d{2}E\\d{2,3})")
		)

		for page := 1; ; page++ {
			q := url.Values{}
			q.Add("cmd", "VideoSearchResultMoviesBlock")
			q.Add("st.cmd", "video")
			q.Add("st.m", "SEARCH")
			q.Add("st.ft", "search")
			q.Add("st.v.sq", f.query)
			q.Add("st.page", strconv.Itoa(page))
			q.Add("gwt.requested", "")
			res, err := http.Get("https://ok.ru/dk?" + q.Encode())
			if err != nil {
				break
			}
			defer res.Body.Close()

			doc, err := goquery.NewDocumentFromReader(res.Body)
			if err != nil {
				break
			}
			sel := doc.Find(".video-card_n")
			if sel.Length() == 0 {
				break
			}
			sel.Each(func(_ int, s *goquery.Selection) {
				var e FetcherEpisode
				for _, a := range s.Nodes[0].Attr {
					switch a.Key {
					case "onclick":
						id := strings.SplitN(strings.SplitN(a.Val, "/video/", 2)[1], "'", 2)[0]
						e.URL = "https://ok.ru/video/" + id
					case "title":
						e.ID = idRegexp.FindString(a.Val)
					}
				}
				ch <- e
			})
		}
	}()
	return ch
}

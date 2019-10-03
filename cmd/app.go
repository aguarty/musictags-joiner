package main

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/patrickmn/go-cache"
)

func (a *application) loadCache(genre string) (*mycache, error) {
	a.storage.Unlock()
	defer a.storage.Lock()

	var genreFile string
	if strings.Index(genre, " ") > 0 {
		genreFile = r.Replace(genre)
	} else {
		genreFile = genre
	}
	cch := &mycache{
		data:    cache.New(25*time.Minute, 100*time.Minute),
		expTime: time.Now(),
	}

	if _, err := os.Stat(StortagsPath + genreFile); os.IsNotExist(err) {
		if err := cch.tagGetTopArtistsGetData(genre, a.cfg); err != nil {
			return cch, err
		}
		cch.data.SaveFile(StortagsPath + genreFile)
	} else {
		if err := cch.data.LoadFile(StortagsPath + genreFile); err != nil {
			return cch, err
		}
	}
	return cch, nil
}

func (a *application) joining(req []string, res *[]resGetTopArtists) error {
	var (
		err      error
		mincash  string
		mincount int = 99000
	)
	a.storage.Lock()
	defer a.storage.Unlock()

	for _, v := range req {
		if _, ok := a.storage.cache[v]; !ok {
			a.storage.cache[v], err = a.loadCache(v)
			if err != nil {
				a.logger.Error(err.Error())
			}
		} else {
			a.storage.cache[v].expTime = time.Now()
		}
		if a.storage.cache[v].data.ItemCount() < mincount {
			mincount = a.storage.cache[v].data.ItemCount()
			mincash = v
		}
	}

	for artist, url := range a.storage.cache[mincash].data.Items() {
		isExist := 0
		for key, _ := range a.storage.cache {
			if key != mincash && containsString(req, key) {
				if _, ok := a.storage.cache[key].data.Get(artist); ok {
					isExist++
				}
			}
		}
		if isExist == len(req)-1 {
			*res = append(*res, resGetTopArtists{
				Name: artist,
				Url:  url.Object.(string),
			})
		}
	}

	return err

}

func (c *mycache) tagGetTopArtistsGetData(genre string, cfg config) error {

	reqHeaders := make(map[string]string)
	reqHeaders["Content-Type"] = "application/json"

	urlParams := url.Values{}
	urlParams.Add("method", MethodTagGetTopArtists)
	urlParams.Add("api_key", cfg.Apikey)
	urlParams.Add("format", "json")
	urlParams.Add("tag", genre)
	urlParams.Add("limit", "500")
	urlParams.Add("page", "0")
	count := 0
	for i := 1; ; i++ {
		m := tagGetTopArtistsData{}
		urlParams.Set("page", strconv.Itoa(i))
		uri := constructUrl(UrlApiBase, urlParams)
		err, _, _ := doRequest("GET", uri, &m, reqHeaders, nil)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		if len(m.Data.Artist) > 0 {
			for _, v := range m.Data.Artist {
				c.data.Add(v.Name, v.Url, cache.NoExpiration)
				count++
			}
		} else {
			break
		}
		if strconv.Itoa(i) == m.Data.Attr.TotalPages {
			break
		}

	}
	return nil
}

func (a *application) storageBouncer(wg *sync.WaitGroup) {

	var interval time.Duration
	interval = 1 * time.Minute

	for {
		select {
		case <-a.ctx.Done():
			a.logger.Info("storage bouncer stoping...")
			wg.Done()
			return
		case <-time.After(interval):
			a.storage.Lock()
			for k, v := range a.storage.cache {
				if time.Since(v.expTime) >= cachettl {
					a.logger.Debug("bounce cache - ", k)
					delete(a.storage.cache, k)
				}
			}
			a.storage.Unlock()
		}
	}
}

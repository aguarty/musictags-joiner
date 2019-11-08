package storage

import (
	"context"
	"fmt"
	"musictags-joiner/internal/utils"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/patrickmn/go-cache"
)

type mycache struct {
	data    *cache.Cache
	expTime time.Time
}

type Storage struct {
	apikey       string
	stortagsPath string
	sync.RWMutex
	cache map[string]*mycache
}

const (
	UrlApiBase string        = "http://ws.audioscrobbler.com/2.0/"
	cachettl   time.Duration = 1 * time.Hour

	MethodTagGetTopArtists string = "tag.gettopartists"
	MethodTagGetTopTags    string = "tag.getTopTags"
)

var r = strings.NewReplacer(" ", "_")

// NewStorage -
func NewStorage(apikey, stortagsPath string) *Storage {
	cache := make(map[string]*mycache)
	return &Storage{apikey: apikey, stortagsPath: stortagsPath, cache: cache}
}

func (s *Storage) loadCache(genre string) (*mycache, error) {
	s.RUnlock()
	defer s.RLock()

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

	if _, err := os.Stat(s.stortagsPath + genreFile); os.IsNotExist(err) {
		if err := cch.tagGetTopArtistsGetData(genre, s.apikey); err != nil {
			return cch, err
		}
		cch.data.SaveFile(s.stortagsPath + genreFile)
	} else {
		if err := cch.data.LoadFile(s.stortagsPath + genreFile); err != nil {
			return cch, err
		}
	}
	return cch, nil
}

func (c *mycache) tagGetTopArtistsGetData(genre string, apikey string) error {

	reqHeaders := make(map[string]string)
	reqHeaders["Content-Type"] = "application/json"

	urlParams := url.Values{}
	urlParams.Add("method", MethodTagGetTopArtists)
	urlParams.Add("api_key", apikey)
	urlParams.Add("format", "json")
	urlParams.Add("tag", genre)
	urlParams.Add("limit", "500")
	urlParams.Add("page", "0")

	count := 0
	for i := 1; ; i++ {
		m := tagGetTopArtistsData{}
		urlParams.Set("page", strconv.Itoa(i))
		uri := utils.ConstructUrl(UrlApiBase, urlParams)
		err, _, _ := utils.DoRequest("GET", uri, &m, reqHeaders, nil)
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

func (s *Storage) Joining(req []string) (res []ResGetTopArtists, err error) {
	var (
		mincash  string
		mincount int = 99000
	)
	s.RLock()
	defer s.RUnlock()

	for _, v := range req {
		if _, ok := s.cache[v]; !ok {
			s.cache[v], err = s.loadCache(v)
			if err != nil {
				return
			}
		} else {
			s.cache[v].expTime = time.Now()
		}
		if s.cache[v].data.ItemCount() < mincount {
			mincount = s.cache[v].data.ItemCount()
			mincash = v
		}
	}

	for artist, url := range s.cache[mincash].data.Items() {
		isExist := 0
		for key, _ := range s.cache {
			if key != mincash && utils.ContainsString(req, key) {
				if _, ok := s.cache[key].data.Get(artist); ok {
					isExist++
				}
			}
		}
		if isExist == len(req)-1 {
			res = append(res, ResGetTopArtists{
				Name: artist,
				Url:  url.Object.(string),
			})
		}
	}

	return
}

func (s *Storage) StorageBouncer(wg *sync.WaitGroup, ctx context.Context) {

	var interval time.Duration
	interval = 1 * time.Minute

	for {
		select {
		case <-ctx.Done():
			wg.Done()
			return
		case <-time.After(interval):
			s.Lock()
			for k, v := range s.cache {
				if time.Since(v.expTime) >= cachettl {
					delete(s.cache, k)
				}
			}
			s.Unlock()
		}
	}
}

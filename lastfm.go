package main

import "github.com/shkh/lastfm-go/lastfm"

type basicTrack struct {
	Title  string
	Artist string
}

type anonFunc func() interface{}

func cacheWith(key string, provider anonFunc) interface{} {
	if CacheIsExist(key) && false {
		return CacheGet(key)
	} else {
		result := provider()
		CachePut(key, result, 3600)
		return result
	}
}

func artistTopTracks(name string) []basicTrack {
	result := cacheWith("artistTopTracks:"+name, func() interface{} {
		r, err := api.Artist.GetTopTracks(lastfm.P{
			"artist": name,
		})
		if err != nil {
			return nil
		}

		results := make([]basicTrack, 0)

		for _, m := range r.Tracks {

			results = append(results, basicTrack{
				Title:  m.Name,
				Artist: name,
			})
		}

		return results
	})
	return result.([]basicTrack)
}

func artistGetSimilar(name string) []string {
	result := cacheWith("artistGetSimilar:"+name, func() interface{} {
		r, err := api.Artist.GetSimilar(lastfm.P{
			"artist": name,
		})
		if err != nil {
			return nil
		}

		results := make([]string, 0)

		for _, m := range r.Similars {
			results = append(results, m.Name)
		}

		return results
	})
	return result.([]string)
}

func tagTopTracks(name string) []basicTrack {
	result := cacheWith("tagTopTracks:"+name, func() interface{} {
		r, err := api.Tag.GetTopTracks(lastfm.P{
			"tag": name,
		})
		if err != nil {
			return nil
		}

		results := make([]basicTrack, 0)

		for _, m := range r.Tracks {
			results = append(results, basicTrack{
				Title:  m.Name,
				Artist: m.Artist.Name,
			})
		}

		return results
	})
	return result.([]basicTrack)
}

func trackGetSimilar(title, artist string) []basicTrack {
	result := cacheWith("trackGetSimilar:"+title+":"+artist, func() interface{} {
		r, err := api.Track.GetSimilar(lastfm.P{
			"track":  title,
			"artist": artist,
		})
		if err != nil {
			return nil
		}

		results := make([]basicTrack, 0)

		for _, m := range r.Tracks {
			results = append(results, basicTrack{
				Title:  m.Name,
				Artist: m.Artist.Name,
			})
		}

		return results
	})
	return result.([]basicTrack)
}

package storage

type tagGetTopArtistsData struct {
	Data struct {
		Artist []artist `json:"artist"`
		Attr   struct {
			Tag        string `json:"tag"`
			Page       string `json:"page"`
			PerPage    string `json:"perPages"`
			TotalPages string `json:"totalPages"`
			Total      string `json:"total"`
		} `json:"@attr"`
	} `json:"topartists"`
}

type artist struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type ResGetTopArtists struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

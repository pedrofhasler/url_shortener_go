package url

type mock_database struct {
	urls   map[string]*Url
	clicks map[string]int
}

func NewMockDatabase() *mock_database {
	return &mock_database{
		make(map[string]*Url),
		make(map[string]int),
	}
}

func (r *mock_database) LookForUrl(url string) *Url {
	for _, u := range r.urls {
		if u.Destination == url {
			return u
		}
	}
	return nil
}

func (r *mock_database) LookForId(id string) *Url {
	return r.urls[id]
}

func (r *mock_database) IdMatch(id string) bool {
	_, match := r.urls[id]
	return match
}

func (r *mock_database) Save(url Url) error {
	r.urls[url.Id] = &url
	return nil
}

func (r *mock_database) LookForClicks(id string) int {
	return r.clicks[id]
}

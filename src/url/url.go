package url

import (
	"math/rand"
	"net/url"
	"time"
)

const (
	newUrlSize = 5
	symbols    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890_-+"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Url struct {
	Id          string    `json:"id"`
	CreatedAt   time.Time `json:"createdAt"`
	Destination string    `json:"destination"`
}

type Stats struct {
	Url    *Url `json:"url"`
	Clicks int  `json:"clicks"`
}

type IUrl interface {
	IdMatch(id string) bool
	LookForId(id string) *Url
	LookForUrl(url string) *Url
	Save(url Url) error
	RegisterClick(id string)
	LookForClicks(id string) int
}

var repo IUrl

func RepositoryConfiguration(r IUrl) {
	repo = r
}

func RegisterClick(id string) {
	repo.RegisterClick(id)
}

func LookOrCreateNewUrl(destination string) (
	u *Url,
	new bool,
	err error) {

	if u = repo.LookForUrl(destination); u != nil {
		return u, false, nil
	}

	if _, err = url.ParseRequestURI(destination); err != nil {
		return nil, false, err
	}

	url := Url{makeId(), time.Now(), destination}
	repo.Save(url)
	return &url, true, nil
}

func makeId() string {
	newId := func() string {
		id := make([]byte, newUrlSize, newUrlSize)
		for i := range id {
			id[i] = symbols[rand.Intn(len(symbols))]
		}
		return string(id)
	}

	for {
		if id := newId(); !repo.IdMatch(id) {
			return id
		}
	}
}

func LookUp(id string) *Url {
	return repo.LookForId(id)
}

func (u *Url) Stats() *Stats {
	clicks := repo.LookForClicks(u.Id)
	return &Stats{u, clicks}
}

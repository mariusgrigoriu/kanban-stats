package trello

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
)

type List struct {
	ID    string
	Name  string
	Cards []Card
}

type Card struct {
	Labels []Label
}

type Label struct {
	ID      string
	IDBoard string
	Name    string
	Color   string
	Uses    int
}

type Client interface {
	GetLists(boardID string) (lists []List)
}

type NetworkClient struct {
	Key   string
	Token string
}

func (c NetworkClient) GetLists(boardID string) (lists []List) {
	query := url.Values{
		"key":         {c.Key},
		"token":       {c.Token},
		"cards":       {"open"},
		"card_fields": {"labels"},
	}
	boardUrl := url.URL{
		Scheme:   "https",
		Host:     "api.trello.com",
		Path:     "1/board/" + boardID + "/lists",
		RawQuery: query.Encode(),
	}

	response, err := http.Get(boardUrl.String())
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		log.Fatal("Trello response: " + response.Status)
	}

	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&lists)
	if err != nil {
		log.Fatal(err)
	}
	return
}

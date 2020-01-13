package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"sync"
)

const apiToken = "xoxb-51553666834-887269487187-8gqXxYVqA0hAW67Zscmbc6d7"

type Command int

const (
	play command = iota
	flip command = iota
	bid  command = iota
	liar command = iota
)

//map[response_url:[https://hooks.slack.com/commands/T1HG9KLQJ/886014563746/HarTowFqJ7KSgU9602DtxwZf] team_id:[T1HG9KLQJ] team_domain:[aviatodev] channel_id:[C1HEJF0DR] text:[] token:[uywRTb17a4I313NX3DLGx9R3] user_id:[U78H0HHPG] command:[/hello] user_name:[lee.painton] channel_name:[random] trigger_id:[900657739207.51553666834.5eff0db620d4fda7f5752e0352182959]]
type message struct {
	command     command
	channelID   channel
	responseURL *url.URL
	text        string
	userID      player
	userName    string
}

type task struct {
	message
}

type Face uint8

const (
	One   Face = iota
	Two   Face = iota
	Three Face = iota
	Four  Face = iota
	Five  Face = iota
	Six   Face = iota
)

type server struct {
	*table
	token string
	sync.Mutex
}

func (s *server) apiHandler(w http.ResponseWriter, r *http.Request) {
	s.Lock()

	/*m, err := parseMsg(r)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), 400)
		return
	}*/

	switch r.URL.Path {
	case "/play":
		//m.command = play
		fmt.Fprintln(w, "Joined game")
	case "/flip":
		//m.command = flip
		fmt.Fprintln(w, "Flipped game")
	case "/bid":
		//m.command = bid
		fmt.Fprintln(w, "Bidding")
	case "/liar":
		//m.command = liar
		fmt.Fprintln(w, "Liar!")
	default:
		fmt.Fprintln(w, "Error!")
	}

	s.Unlock()
}

func parseMsg(r *http.Request) (*message, error) {
	r.ParseForm()
	for _, f := range []string{
		"channel_id",
		"text",
		"user_id",
		"user_name",
	} {
		if v := r.Form[f]; len(v) == 0 {
			err := fmt.Errorf("missing field %q in api call", f)
			return nil, err
		}
	}

	return &message{
		channelID: channel(r.Form["channel_id"][0]),
		text:      r.Form["text"][0],
		userID:    player(r.Form["user_id"][0]),
		userName:  r.Form["user_name"][0],
	}, nil
}

func main() {
	s := &server{
		token: apiToken,
	}

	http.HandleFunc("/", s.apiHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

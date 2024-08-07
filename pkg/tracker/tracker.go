package tracker

import (
	"encoding/json"
	"net/http"
)

const TrackerUrl = "https://us.forums.blizzard.com/en/wow/groups/blizzard-tracker/posts.json"

type Post struct {
	Id        int
	Excerpt   string
	Truncated bool
	TopicId   int
}

func GetPosts() ([]Post, error) {
	resp, err := http.Get(TrackerUrl)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	posts := make([]Post, 0)
	jDecoder := json.NewDecoder(resp.Body)
	decoderErr := jDecoder.Decode(&posts)

	if decoderErr != nil {
		return nil, decoderErr
	}

	return posts, nil

}

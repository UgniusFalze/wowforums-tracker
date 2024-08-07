package tracker

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sync"
)

const TrackerUrl = "https://us.forums.blizzard.com/en/wow/groups/blizzard-tracker/posts.json"

const TopicUrl = "https://us.forums.blizzard.com/en/wow/t"

type Post struct {
	Id        int
	Excerpt   string
	Truncated bool
	Topic_id  int
}

type topic struct {
	Post_stream post_stream
}

type post_stream struct {
	Posts []TopicPost
}

type TopicPost struct {
	Cooked           string
	Display_username string
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

func GetTopicContent(topicId int, postId int) (*topic, error) {
	url, err := url.Parse(fmt.Sprintf("%s/%d/posts.json?post_ids[]=%d", TopicUrl, topicId, postId))
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(url.String())

	if err != nil {
		return nil, err
	}

	if resp.StatusCode > 299 {
		return nil, fmt.Errorf("error: %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	posts := &topic{}
	jDecoder := json.NewDecoder(resp.Body)
	decoderErr := jDecoder.Decode(&posts)

	if decoderErr != nil {
		return nil, decoderErr
	}

	return posts, nil
}

type TopicResult struct {
	Topic TopicPost
	Error error
}

func GetPostsTopics(posts []Post) []TopicResult {
	var wg sync.WaitGroup
	topics := make([]TopicResult, len(posts))
	for i, param := range posts {
		if param.Truncated {
			wg.Add(1)
			go func(i int, param Post) {
				defer wg.Done()
				topic, err := GetTopicContent(param.Topic_id, param.Id)
				if err != nil {
					topics[i] = TopicResult{Error: err}
				} else {
					topics[i] = TopicResult{Topic: topic.Post_stream.Posts[0]}
				}
			}(i, param)
		} else {
			topics[i] = TopicResult{Topic: TopicPost{Cooked: param.Excerpt}}
		}
	}

	wg.Wait()

	return topics
}

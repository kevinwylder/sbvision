package sources_test

import (
	"testing"

	"github.com/kevinwylder/sbvision"

	"github.com/kevinwylder/sbvision/media/sources"
)

func TestSources(t *testing.T) {
	fails := []string{
		"https://www.reddit.com/r/funny/comments/ased23/my_taste_in_music",
		"https://www.reddit.com/r/skateboarding/ased23/my_taste_in_music",
		"https://www.youtube.com/watch",
		"https://www.google.com/search?q=recursion",
	}

	for _, url := range fails {
		if _, err := sources.FindVideoSource(url); err == nil {
			t.Fail()
		}
	}

	reddit := []string{
		"https://www.reddit.com/r/skateboarding/comments/g0fo0y/5050_manual_shuv_out_from_a_little_while_ago/",
	}

	for _, url := range reddit {
		source, err := sources.FindVideoSource(url)
		if err != nil {
			t.Fatal(err)
		}
		if source.Type() != sbvision.RedditVideo {
			t.Fail()
		}
	}

	youtube := []string{
		"https://youtu.be/6s4Bx7mzNkM",
		"https://www.youtube.com/watch?v=yNsS9uqsR3k",
	}

	for _, url := range youtube {
		source, err := sources.FindVideoSource(url)
		if err != nil {
			t.Fail()
		}
		if source.Type() != sbvision.YoutubeVideo {
			t.Fail()
		}
	}
}

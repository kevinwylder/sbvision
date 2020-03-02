package reddit

type listing struct {
	Data struct {
		Children []struct {
			Data entry `json:"data"`
		} `json:"children"`
	} `json:"data"`
}

type entry struct {
	IsVideo    bool   `json:"is_video"`
	Domain     string `json:"domain"`
	YoutubeURL string `json:"url"`
	RedditLink struct {
		RedditVideo struct {
			URL string `json:"fallback_url"`
		} `json:"reddit_video"`
	} `json:"secure_media"`
}

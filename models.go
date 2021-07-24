package main

type searchListResponse struct {
	Kind string `json:"kind"`
	NextPageToken string `json:"nextPageToken"`
	Items []searchItem
}

type searchItem struct {
	Kind string `json:"kind"`
	Snippet snippet
	ID id
}

type snippet struct {
	PublishedAt string `json:"publishedAt"`
	ChannelID string `json:"channelId"`
	Title string `json:"title"`
	ChannelTitle string `json:"channelTitle"`
	CategoryID string `json:"categoryId"`
}

type id struct {
	VideoID string `json:"videoId"`
}

type videoListResponse struct {
	Kind string `json:"kind"`
	Items []videoItem
}

type videoItem struct {
	Kind	string `json:"kind"`
	ID		string `json:"id"`
	Snippet	snippet
	Statistics statistics
	Status	status
	ContentDetails contentDetails
}

type statistics struct {
	ViewCount string `json:"viewCount"`
	LikeCount string `json:"likeCount"`
	DislikeCount string `json:"dislikeCount"`
	FavoriteCount string `json:"favoriteCount"`
	CommentCount string `json:"commentCount"`
}

type status struct {
	PrivacyStatus string `json:"privacyStatus"`
}

type contentDetails struct {
	Duration string `json:"duration"`
}

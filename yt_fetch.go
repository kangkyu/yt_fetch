package main

import (
	"html/template"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

func main() {
	port := "8080"
	if value, ok := os.LookupEnv("PORT"); ok {
		port = value
	}
	http.HandleFunc("/", pageHandler)
	http.HandleFunc("/fetches", fetchHandler)
	log.Println("Listen on localhost:"+port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func pageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	t, err := template.ParseFiles("page.html")
	if err != nil {
		http.Error(w, "file not found", 404)
		return
	}
	t.Execute(w, "")
}

func fetchHandler(rw http.ResponseWriter, r *http.Request) {
	channelID := r.FormValue("uuid")

	// fetch

	rw.Header().Set("Content-Type", "text/csv")
	rw.Header().Set("Content-Disposition", "attachment;filename=result.csv")
	w := csv.NewWriter(rw)

	generateCSV(w, channelID)

	w.Flush()
	if err := w.Error(); err != nil {
		log.Fatal(err)
	}
}

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

func generateCSV(w *csv.Writer, uuid string) {

	header := []string{
		"kind",
		"publishedAt",
		"channelId",
		"channelTitle",
		"id",
		"title",
		"categoryId",
		"viewCount",
		"likeCount",
		"dislikeCount",
		"favoriteCount",
		"commentCount",
		"privacyStatus",
		"duration",
	}
	if err := w.Write(header); err != nil {
		log.Fatalln("error writing record to csv:", err)
	}

	var u *url.URL
	var sl searchListResponse
	var nextPageToken string

	for {
		u = searchURL(nextPageToken, uuid)
		sl = printVideos(u, w)
		nextPageToken = sl.NextPageToken

		if len(nextPageToken) == 0 || len(sl.Items) == 0 {
			break
		}
	}
}

func printVideos(su *url.URL, w *csv.Writer) searchListResponse {

	sl := searchListFromSearchURL(su)
	if len(sl.Items) == 0 {
		return sl
	}

	vu := videosURL(sl)
	vl := videoListFromVideosURL(vu)

	vl.print(w) // print

	return sl
}

func videoListFromVideosURL(u *url.URL) videoListResponse {

	fmt.Printf("%s\n", u.String())

	resp, err := http.Get(u.String()) //
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	s := string(body)
	bs := []byte(s)

	var videoList = videoListResponse{}
	err = json.Unmarshal(bs, &videoList)
	if err != nil {
		fmt.Println(err)
	}

	return videoList
}

func searchListFromSearchURL(uuu *url.URL) searchListResponse {

	fmt.Printf("%s\n", uuu.String())

	resp, err := http.Get(uuu.String()) //
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	s := string(body)
	bs := []byte(s)

	var searchList = searchListResponse{}
	err = json.Unmarshal(bs, &searchList)
	if err != nil {
		fmt.Println(err)
	}

	return searchList
}

func searchURL(nextPageToken string, uuid string) *url.URL {

	u, err := url.Parse("https://www.googleapis.com/youtube/v3/search")
	if err != nil {
		log.Fatal(err)
	}

	v := url.Values{}
	v.Set("key", os.Getenv("YT_API_KEY"))
	v.Add("part", "snippet")
	v.Add("type", "video")
	v.Add("maxResults", "50")
	v.Add("order", "date")
	v.Add("channelId", uuid)

	if len(nextPageToken) != 0 {
		v.Set("pageToken", nextPageToken)
	}

	u.RawQuery = v.Encode()
	return u
}


func videosURL(searchList searchListResponse) *url.URL {

	videoIDString := strings.Join(*searchList.videoIDs(), ",")

	u, err := url.Parse("https://www.googleapis.com/youtube/v3/videos")
	if err != nil {
		log.Fatal(err)
	}

	v := url.Values{}
	v.Set("key", os.Getenv("YT_API_KEY"))
	v.Add("part", "snippet,statistics,status,contentDetails")
	v.Add("id", videoIDString)

	u.RawQuery = v.Encode()
	return u
}

func (videoList videoListResponse) print(w *csv.Writer) {

	for _, item := range videoList.Items {
		parsedTime, err := time.Parse("2006-01-02T15:04:05Z07:00", item.Snippet.PublishedAt)
		if err != nil {
			fmt.Println(err)
		}
		record := []string{
			item.Kind,
			parsedTime.Local().Format("2006-01-02"),
			item.Snippet.ChannelID,
			item.Snippet.ChannelTitle,
			item.ID,
			item.Snippet.Title,
			item.Snippet.CategoryID,
			item.Statistics.ViewCount,
			item.Statistics.LikeCount,
			item.Statistics.DislikeCount,
			item.Statistics.FavoriteCount,
			item.Statistics.CommentCount,
			item.Status.PrivacyStatus,
			item.ContentDetails.Duration,
		}
		if err := w.Write(record); err != nil {
			log.Fatalln("error writing record to csv:", err)
		}
	}
}

func (searchList searchListResponse) videoIDs() *[]string {
	// items.collect{|item| item.video_id}.compact.uniq.join(",")
	keys := make(map[string]bool)
	ids := []string{}
	for _, item := range searchList.Items {
		id := item.ID.VideoID
		if _, value := keys[id]; !value {
			keys[id] = true
			ids = append(ids, id)
		}
	}
	return &ids
}

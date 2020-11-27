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

func fetchHandler(w http.ResponseWriter, r *http.Request) {
	channelID := r.FormValue("uuid")
	// TODO: need validation channelID presence

	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment;filename=videosof"+channelID+".csv")
	cw := csv.NewWriter(w)

	err := generateCSV(cw, channelID)
	if err != nil {
		log.Println(err)
		fmt.Fprint(w, "could not generate CSV:\n")
		fmt.Fprint(w, err.Error())
		return
	}

	cw.Flush()
	if err := cw.Error(); err != nil {
		http.Error(w, "internal error", 500)
		return
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

func generateCSV(cw *csv.Writer, uuid string) error {

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
	if err := cw.Write(header); err != nil {
		return fmt.Errorf("error writing record to csv: %v", err)
	}

	var nextPageToken string

	for {
		su := searchURL(nextPageToken, uuid)
		sl, err := searchListFromSearchURL(su)
		if err != nil {
			return err
		}

		vu := videosURL(sl.videoIDs())
		vl, err := videoListFromVideosURL(vu)
		if err != nil {
			return err
		}

		vl.print(cw)

		nextPageToken = sl.NextPageToken

		if len(nextPageToken) == 0 || len(sl.Items) == 0 {
			break
		}
	}
	return nil
}

func videoListFromVideosURL(vu *url.URL) (videoListResponse, error) {

	var videoList = videoListResponse{}

	resp, err := http.Get(vu.String())
	if err != nil {
		return videoList, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return videoList, err
	}

	s := string(body)
	bs := []byte(s)

	err = json.Unmarshal(bs, &videoList)
	if err != nil {
		return videoList, err
	}

	return videoList, nil
}

func searchListFromSearchURL(su *url.URL) (searchListResponse, error) {

	var searchList = searchListResponse{}
	// fmt.Printf("%s\n", su.String())

	resp, err := http.Get(su.String())
	if err != nil {
		return searchList, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return searchList, err
	}

	s := string(body)
	bs := []byte(s)

	if resp.StatusCode != 200 {
		return searchList, fmt.Errorf(s)
	}

	err = json.Unmarshal(bs, &searchList)
	if err != nil {
		return searchList, err
	}

	return searchList, nil
}

func searchURL(nextPageToken string, uuid string) *url.URL {

	u, _ := url.Parse("https://www.googleapis.com/youtube/v3/search")

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


func videosURL(videoIDs []string) *url.URL {

	u, _ := url.Parse("https://www.googleapis.com/youtube/v3/videos")

	v := url.Values{}
	v.Set("key", os.Getenv("YT_API_KEY"))
	v.Add("part", "snippet,statistics,status,contentDetails")
	v.Add("id", strings.Join(videoIDs, ","))

	u.RawQuery = v.Encode()
	return u
}

func (videoList videoListResponse) print(cw *csv.Writer) error {

	for _, item := range videoList.Items {
		parsedTime, err := time.Parse("2006-01-02T15:04:05Z07:00", item.Snippet.PublishedAt)
		if err != nil {
			return err
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
		if err := cw.Write(record); err != nil {
			return fmt.Errorf("error writing record to csv: %v", err)
		}
	}
	return nil
}

func (searchList searchListResponse) videoIDs() []string {
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
	return ids
}

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

	if err := generateCSV(w, channelID); err != nil {
		fmt.Fprint(w, "could not generate CSV:\n")
		fmt.Fprint(w, err.Error())
		// TODO: error response should be 400 or 500 by cases, how do you tell?
		// http.Error(w, "internal error", 500)
	}
}

func generateCSV(w http.ResponseWriter, uuid string) error {
	cw := csv.NewWriter(w)
	defer cw.Flush()

	if err := printHeaderRow(cw); err != nil {
		return err
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

		err = printVideos(cw, vl)
		if err != nil {
			return err
		}

		nextPageToken = sl.NextPageToken

		if len(nextPageToken) == 0 || len(sl.Items) == 0 {
			break
		}
	}
	return cw.Error()
}

func videoListFromVideosURL(vu *url.URL) (videoListResponse, error) {

	var videoList = videoListResponse{}

	resp, err := http.Get(vu.String())
	if err != nil {
		return videoList, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&videoList)
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

	if resp.StatusCode != 200 {
		return searchList, fmt.Errorf("youtube error response: %v", string(body))
	}

	err = json.Unmarshal(body, &searchList)
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

func printHeaderRow(cw *csv.Writer) error {
	header := make([]string, 0, len(fields))

	for i := 0; i < len(fields); i++ {
		for field, value := range fields {
			if value == i {
				header = append(header, field)
			}
		}
	}
	if err := cw.Write(header); err != nil {
		return fmt.Errorf("error writing header to csv: %v", err)
	}
	return nil
}

var fields = map[string]int{
	"kind": 0,
	"publishedAt": 1,
	"channelId": 2,
	"channelTitle": 3,
	"id": 4,
	"title": 5,
	"categoryId": 6,
	"viewCount": 7,
	"likeCount": 8,
	"dislikeCount": 9,
	"favoriteCount": 10,
	"commentCount": 11,
	"privacyStatus": 12,
	"duration": 13,
}

func printVideos(cw *csv.Writer, videoList videoListResponse) error {
	row := make([]string, 14)

	for _, item := range videoList.Items {

		row[fields["kind"]] = item.Kind
		publishedOn, err := parseDate(item.Snippet.PublishedAt)
		if err != nil {
			return fmt.Errorf("error parsing publishedAt (can it be?)", err)
		}
		row[fields["publishedAt"]] = publishedOn
		row[fields["channelId"]] = item.Snippet.ChannelID
		row[fields["channelTitle"]] = item.Snippet.ChannelTitle
		row[fields["id"]] = item.ID
		row[fields["title"]] = item.Snippet.Title
		row[fields["categoryId"]] = item.Snippet.CategoryID
		row[fields["viewCount"]] = item.Statistics.ViewCount
		row[fields["likeCount"]] = item.Statistics.LikeCount
		row[fields["dislikeCount"]] = item.Statistics.DislikeCount
		row[fields["favoriteCount"]] = item.Statistics.FavoriteCount
		row[fields["commentCount"]] = item.Statistics.CommentCount
		row[fields["privacyStatus"]] = item.Status.PrivacyStatus
		row[fields["duration"]] = item.ContentDetails.Duration
		if err := cw.Write(row); err != nil {
			return fmt.Errorf("error writing a row to csv: %v", err)
		}
	}
	return nil
}

func parseDate(ts string) (string, error) {
	t, err := time.Parse("2006-01-02T15:04:05Z07:00", ts)
	if err != nil {
		return "", err
	}
	return t.Local().Format("2006-01-02"), nil
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

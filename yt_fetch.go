package main

import (
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

type searchListResponse struct {
    Kind          string `json:"kind"`
    NextPageToken string `json:"nextPageToken"`
    Items         []searchItem
}

type searchItem struct {
    Kind    string `json:"kind"`
    Snippet snippet
    ID      id
}

type snippet struct {
    PublishedAt string `json:"publishedAt"`
    ChannelID   string `json:"channelId"`
}

type id struct {
    VideoID string `json:"videoId"`
}

type videoListResponse struct {
    Kind  string `json:"kind"`
    Items []videoItem
}

type videoItem struct {
    Kind       string `json:"kind"`
    ID         string `json:"id"`
    Snippet    snippet
    Statistics statistics
}

type statistics struct {
    ViewCount string `json:"viewCount"`
}

func main() {

    file, err := os.Create("result.csv")
    if err != nil {
        log.Fatal("Cannot create file", err)
    }
    defer file.Close()

    w := csv.NewWriter(file)

    generateCSV(w)

    w.Flush()
    if err := w.Error(); err != nil {
        log.Fatal(err)
    }
}

func generateCSV(w *csv.Writer) {

    header := []string{"Kind", "PublishedAt", "ChannelID", "ID", "ViewCount"}
    if err := w.Write(header); err != nil {
        log.Fatalln("error writing record to csv:", err)
    }

    var u *url.URL
    var sl searchListResponse
    var nextPageToken string

    for {
        u = searchURL(nextPageToken)
        sl = printVideos(u, w)
        nextPageToken = sl.NextPageToken

        if len(nextPageToken) == 0 || len(sl.Items) == 0 {
            break
        }
    }
}

func printVideos(uuu *url.URL, w *csv.Writer) searchListResponse {

    sl := searchListFromSearchURL(uuu)
    if len(sl.Items) == 0 {
        return sl
    }

    u := videosURL(sl)
    vl := videoListFromVideosURL(u)

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

func searchURL(nextPageToken string) *url.URL {

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
    v.Add("channelId", os.Args[1])

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
    v.Add("part", "snippet,statistics")
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
            item.ID, item.Statistics.ViewCount}
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

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
    u.RawQuery = v.Encode()

    nextPageToken, itemCount := videosFromURL(u.String(), w)
    for {
        q := u.Query()
        q.Set("pageToken", nextPageToken)
        u.RawQuery = q.Encode()

        nextPageToken, itemCount = videosFromURL(u.String(), w)
        if len(nextPageToken) == 0 || itemCount == 0 {
            break
        }
    }

    w.Flush()

    if err := w.Error(); err != nil {
        log.Fatal(err)
    }
}

func videosFromURL(uuu string, w *csv.Writer) (string, int) {
    fmt.Printf("%s\n", uuu)

    resp, err := http.Get(uuu)
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

    if len(searchList.Items) == 0 {
        return searchList.NextPageToken, 0
    }

    videoIDString := videoIDString(searchList.Items)

    u3, err := url.Parse("https://www.googleapis.com/youtube/v3/videos")
    if err != nil {
        log.Fatal(err)
    }
    v3 := url.Values{}
    v3.Set("key", os.Getenv("YT_API_KEY"))
    v3.Add("part", "snippet,statistics")
    v3.Add("id", videoIDString)
    u3.RawQuery = v3.Encode()

    fmt.Printf("%s\n", u3.String())

    resp3, err := http.Get(u3.String())
    if err != nil {
        log.Fatal(err)
    }
    defer resp3.Body.Close()
    body3, err := ioutil.ReadAll(resp3.Body)

    s3 := string(body3)
    bs3 := []byte(s3)

    var videoList = videoListResponse{}
    err = json.Unmarshal(bs3, &videoList)
    if err != nil {
        fmt.Println(err)
    }


    header := []string{"Kind", "PublishedAt", "ChannelID", "ID", "ViewCount"}
    if err := w.Write(header); err != nil {
        log.Fatalln("error writing record to csv:", err)
    }
    for _, item := range videoList.Items {
        parsedTime, err := time.Parse("2006-01-02T15:04:05Z07:00", item.Snippet.PublishedAt)
        if err != nil {
            fmt.Println(err)
        }
        record := []string{item.Kind, parsedTime.Local().Format("2006-01-02"), item.Snippet.ChannelID, item.ID, item.Statistics.ViewCount}
        if err := w.Write(record); err != nil {
            log.Fatalln("error writing record to csv:", err)
        }
    }

    return searchList.NextPageToken, len(searchList.Items)
}

func videoIDString(items []searchItem) string {
    ids := make([]string, 0, len(items))
    for _, item := range items {
        ids = append(ids, item.ID.VideoID)
    }
    return strings.Join(ids, ",")
}

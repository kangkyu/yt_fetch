package main

import (
    "fmt"
    "net/url"
    "net/http"
    "io/ioutil"
    "os"
    "log"
    "encoding/json"
    "encoding/csv"
    "strings"
)

type searchListResponse struct {
    Kind string `json:"kind"`
    NextPageToken  string `json:"nextPageToken"`
    Items []searchItem
}

type searchItem struct {
    Kind string `json:"kind"`
    Snippet snippet
    Id id
}

type snippet struct {
    PublishedAt string `json:"publishedAt"`
    ChannelId string `json:"channelId"`
}

type id struct {
    VideoId string `json:"videoId"`
}

type videoListResponse struct {
    Kind string `json:"kind"`
    Items []videoItem
}

type videoItem struct {
    Kind string `json:"kind"`
    Id string `json:"id"`
    Snippet snippet
    Statistics statistics
}

type statistics struct {
    ViewCount string `json:"viewCount"`
}

func main() {
    w := csv.NewWriter(os.Stdout)

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

    nextPageToken, itemCount := VideosFromUrl(u.String(), w)
    for {
        q := u.Query()
        q.Set("pageToken", nextPageToken)
        u.RawQuery = q.Encode()

        nextPageToken, itemCount = VideosFromUrl(u.String(), w)
        if len(nextPageToken) == 0 || itemCount == 0 {
            break
        }
    }

    w.Flush()

    if err := w.Error(); err != nil {
        log.Fatal(err)
    }
}

func VideosFromUrl(uuu string, w *csv.Writer) (string, int) {
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

    if len(searchList.Items) == 0 { return searchList.NextPageToken, 0 }

    videoIdString := VideoIdString(searchList.Items)

    u3, err := url.Parse("https://www.googleapis.com/youtube/v3/videos")
    if err != nil {
        log.Fatal(err)
    }
    v3 := url.Values{}
    v3.Set("key", os.Getenv("YT_API_KEY"))
    v3.Add("part", "snippet,statistics")
    v3.Add("id", videoIdString)
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

    for _, item := range videoList.Items {
        record := []string{ item.Kind, item.Snippet.PublishedAt, item.Snippet.ChannelId, item.Id, item.Statistics.ViewCount }
        if err := w.Write(record); err != nil {
            log.Fatalln("error writing record to csv:", err)
        }
    }

    return searchList.NextPageToken, len(searchList.Items)
}

func VideoIdString(items []searchItem) string {
    ids := make([]string, 0, len(items))
    for _, item := range items {
        ids = append(ids, item.Id.VideoId)
    }
    return strings.Join(ids, ",")
}

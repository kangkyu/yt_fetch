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

    fmt.Printf("%s\n", u.String())


    resp, err := http.Get(u.String())
    if err != nil {
        log.Fatal(err)
    }
    defer resp.Body.Close()

    // fmt.Println("Response status:", resp.Status)

    body, err := ioutil.ReadAll(resp.Body)
    // fmt.Println(string(body))

    s := string(body)
    bs := []byte(s)

    var searchList = searchListResponse{}
    err = json.Unmarshal(bs, &searchList)
    if err != nil {
        fmt.Println(err)
    }
    // fmt.Println(searchList)


    videoIdString := VideoIdString(searchList.Items)
    // fmt.Println(videoIdString)

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
    // fmt.Println(string(body3))

    s3 := string(body3)
    bs3 := []byte(s3)

    var videoList = videoListResponse{}
    err = json.Unmarshal(bs3, &videoList)
    if err != nil {
        fmt.Println(err)
    }
    // fmt.Println(videoList)


    for _, item := range videoList.Items {
        record := []string{ item.Kind, item.Snippet.PublishedAt, item.Snippet.ChannelId, item.Id, item.Statistics.ViewCount }
        if err := w.Write(record); err != nil {
            log.Fatalln("error writing record to csv:", err)
        }
    }

    // loop! while there are results
    var nextPageToken string = searchList.NextPageToken
    for {
        q := u.Query()
        q.Set("pageToken", nextPageToken)
        u.RawQuery = q.Encode()

        fmt.Printf("%s\n", u.String())


        resp2, err := http.Get(u.String())
        if err != nil {
            log.Fatal(err)
        }
        defer resp2.Body.Close()

        // fmt.Println("Response status:", resp2.Status)

        body2, err := ioutil.ReadAll(resp2.Body)
        // fmt.Println(string(body2))

        s2 := string(body2)
        bs2 := []byte(s2)

        var searchList2 = searchListResponse{}
        err = json.Unmarshal(bs2, &searchList2)
        if err != nil {
           fmt.Println(err)
        }
        // fmt.Println(searchList)
        // searchList.Items = append(searchList.Items, searchList2.Items...)

        videoIdString2 := VideoIdString(searchList2.Items)
        // fmt.Println(videoIdString2)

        u4, err := url.Parse("https://www.googleapis.com/youtube/v3/videos")
        if err != nil {
            log.Fatal(err)
        }
        v4 := url.Values{}
        v4.Set("key", os.Getenv("YT_API_KEY"))
        v4.Add("part", "snippet,statistics")
        v4.Add("id", videoIdString2)
        u4.RawQuery = v4.Encode()

        fmt.Printf("%s\n", u4.String())

        resp4, err := http.Get(u4.String())
        if err != nil {
            log.Fatal(err)
        }
        defer resp4.Body.Close()
        body4, err := ioutil.ReadAll(resp4.Body)
        // fmt.Println(string(body4))

        s4 := string(body4)
        bs4 := []byte(s4)

        var videoList2 = videoListResponse{}
        err = json.Unmarshal(bs4, &videoList2)
        if err != nil {
            fmt.Println(err)
        }
        // fmt.Println(videoList2)


        for _, item := range videoList2.Items {
            record := []string{ item.Kind, item.Snippet.PublishedAt, item.Snippet.ChannelId, item.Id, item.Statistics.ViewCount }
            if err := w.Write(record); err != nil {
                log.Fatalln("error writing record to csv:", err)
            }
        }

        nextPageToken = searchList2.NextPageToken
        if len(searchList2.NextPageToken) == 0 || len(searchList2.Items) == 0 {
            break
        }
    }


    w.Flush()

    if err := w.Error(); err != nil {
        log.Fatal(err)
    }
}

func VideoIdString(items []searchItem) string {
    ids := make([]string, 0, len(items))
    for _, item := range items {
        ids = append(ids, item.Id.VideoId)
    }
    return strings.Join(ids, ",")
}

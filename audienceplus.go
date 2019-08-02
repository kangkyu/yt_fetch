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
)

type searchListResponse struct {
    Kind string `json:"kind"`
    NextPageToken  string `json:"nextPageToken"`
    Items []item
}

type item struct {
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

func main() {
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

    fmt.Println("Response status:", resp.Status)

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

        fmt.Println("Response status:", resp2.Status)

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
        searchList.Items = append(searchList.Items, searchList2.Items...)
        nextPageToken = searchList2.NextPageToken

        if len(searchList2.NextPageToken) == 0 || len(searchList2.Items) == 0 {
            break
        }
    }

    w := csv.NewWriter(os.Stdout)

    for _, item := range searchList.Items {
        record := []string{ item.Kind, item.Snippet.PublishedAt, item.Snippet.ChannelId, item.Id.VideoId }
        if err := w.Write(record); err != nil {
            log.Fatalln("error writing record to csv:", err)
        }
    }

    w.Flush()

    if err := w.Error(); err != nil {
        log.Fatal(err)
    }
}

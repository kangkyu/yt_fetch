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
}

type snippet struct {
    PublishedAt string `json:"publishedAt"`
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

    w := csv.NewWriter(os.Stdout)

    for _, item := range searchList.Items {
        record := []string{ item.Kind, item.Snippet.PublishedAt }
        if err := w.Write(record); err != nil {
            log.Fatalln("error writing record to csv:", err)
        }
    }

    w.Flush()

    if err := w.Error(); err != nil {
        log.Fatal(err)
    }
}

package main

import (
	"fmt"
	"net/url"
	"net/http"
	"io/ioutil"
	"os"
	"log"
)

func main() {
	u, err := url.Parse("https://www.googleapis.com/youtube/v3/search")
	if err != nil {
		log.Fatal(err)
	}
	v := url.Values{}
	v.Set("key", os.Getenv("YT_API_KEY"))
	v.Add("q", "penguin")
	v.Add("part", "snippet")
	v.Add("type", "video")
	v.Add("maxResults", "2")
	u.RawQuery = v.Encode()
	fmt.Printf("%s\n", u.String())
	resp, err := http.Get(u.String())
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}

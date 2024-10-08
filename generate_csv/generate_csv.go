package generate_csv

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"sync"
	"time"

	yt "github.com/kangkyu/youtube_api"
)

func GenerateCSV(w http.ResponseWriter, uuid string) error {
	cw := csv.NewWriter(w)
	defer cw.Flush()

	if err := printHeaderRow(cw); err != nil {
		return err
	}

	client := yt.NewClient(uuid)
	fetcher := client.NewPaginatedFetcher()

	// Begin with a channel to receive items
	itemChan := make(chan yt.VideoItem)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			response, err := fetcher.FetchNextPage()
			if err != nil {
				fmt.Println("Error fetching data:", err)
				close(itemChan)
				return
			}

			for _, item := range response.Items {
				itemChan <- item
			}

			if !fetcher.HasNextPage() {
				close(itemChan)
				return
			}
		}
	}()

	for item := range itemChan {
		printVideo(cw, item)
	}

	wg.Wait()

	return cw.Error()
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
	"kind":          0,
	"publishedAt":   1,
	"channelId":     2,
	"channelTitle":  3,
	"id":            4,
	"title":         5,
	"categoryId":    6,
	"viewCount":     7,
	"likeCount":     8,
	"dislikeCount":  9,
	"favoriteCount": 10,
	"commentCount":  11,
	"privacyStatus": 12,
	"duration":      13,
}

func printVideo(cw *csv.Writer, item yt.VideoItem) error {
	row := make([]string, len(fields))
	row[fields["kind"]] = item.Kind
	publishedOn, err := parseDate(item.Snippet.PublishedAt)
	if err != nil {
		return fmt.Errorf("error parsing publishedAt: %v", err)
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
	return nil
}

func parseDate(ts string) (string, error) {
	t, err := time.Parse("2006-01-02T15:04:05Z07:00", ts)
	if err != nil {
		return "", err
	}
	return t.Local().Format("2006-01-02"), nil
}

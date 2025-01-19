# yt_fetch

Test application for [youtube_api](https://github.com/kangkyu/youtube_api) module.

```sh
$ go build .
$ ./yt_fetch
```
And then, open `localhost:8080` page.

## How to use

+ Fill the input field a channel ID (something like 'UCmKH1rkv9OJ1FjM64qbeLVw')
+ Click the "Generate CSV" button to download

## Environment Variables

+ `YT_API_KEY`: It's there at [Google developer console](https://console.developers.google.com). See APIs > Credentials

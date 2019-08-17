# Audience Plus Go

+ get "/search" response
+ parse JSON
+ generate CSV string
+ get a channel ID as ARGV (for example, "./yt_fetch UCHPA-41shIcxUHAbhdQUKSA")
+ get all videos with nextPageToken ("pageToken" option)
+ get "/videos" response, for metrics
+ save the CSV as a file
+ major refactor for readability and maintainability
+ add all the columns
+ send email with the CSV file attached
+ build a web application
  - server
  - form submit page
  - download button for CSV file
  - use goroutine for multiple fetches by multiple users
  - user login

Sample response of video search:

```json
{
 "kind": "youtube#searchListResponse",
 "etag": "\"Bdx4f4ps3xCOOo1WZ91nTLkRZ_c/AHXPcsSnHb7Jef_R4a143DkxHDo\"",
 "nextPageToken": "CAIQAA",
 "regionCode": "US",
 "pageInfo": {
  "totalResults": 1000000,
  "resultsPerPage": 2
 },
 "items": [
  {
   "kind": "youtube#searchResult",
   "etag": "\"Bdx4f4ps3xCOOo1WZ91nTLkRZ_c/tfxEnPjX9PHHbe_UOyZWL5rF-1E\"",
   "id": {
    "kind": "youtube#video",
    "videoId": "IvkfpgjBt5k"
   },
   "snippet": {
    "publishedAt": "2018-12-28T17:00:02.000Z",
    "channelId": "UCwmZiChSryoWQCZMIQezgTg",
    "title": "Penguin Chicks&#39; Stand Off Against Predator | Spy In The Snow | BBC Earth",
    "description": "When a petrel attackes them, emperor penguing chicks stand together against it. Exclusive preview from #SpyInTheSnow, out in the UK December 30th.",
    "thumbnails": {
     "default": {
      "url": "https://i.ytimg.com/vi/IvkfpgjBt5k/default.jpg",
      "width": 120,
      "height": 90
     },
     "medium": {
      "url": "https://i.ytimg.com/vi/IvkfpgjBt5k/mqdefault.jpg",
      "width": 320,
      "height": 180
     },
     "high": {
      "url": "https://i.ytimg.com/vi/IvkfpgjBt5k/hqdefault.jpg",
      "width": 480,
      "height": 360
     }
    },
    "channelTitle": "BBC Earth",
    "liveBroadcastContent": "none"
   }
  },
  {
   "kind": "youtube#searchResult",
   "etag": "\"Bdx4f4ps3xCOOo1WZ91nTLkRZ_c/f_4Zgx-1w9T_tgjofY2K-IzjtXo\"",
   "id": {
    "kind": "youtube#video",
    "videoId": "0syrKzXfKuI"
   },
   "snippet": {
    "publishedAt": "2013-07-05T17:05:40.000Z",
    "channelId": "UCZftwGKErylOICBSDLZ8lBw",
    "title": "PINGU FULL ESPISODES 1 3   YouTube",
    "description": "Please Subscribe & Like Facebook page & Also see the pingu series there ........ https://www.facebook.com/pages/Pingu/554908081221022 Regard pingu ...",
    "thumbnails": {
     "default": {
      "url": "https://i.ytimg.com/vi/0syrKzXfKuI/default.jpg",
      "width": 120,
      "height": 90
     },
     "medium": {
      "url": "https://i.ytimg.com/vi/0syrKzXfKuI/mqdefault.jpg",
      "width": 320,
      "height": 180
     },
     "high": {
      "url": "https://i.ytimg.com/vi/0syrKzXfKuI/hqdefault.jpg",
      "width": 480,
      "height": 360
     }
    },
    "channelTitle": "pingu netwok",
    "liveBroadcastContent": "none"
   }
  }
 ]
}
```

Sample response of video videos list
```json
{
 "kind": "youtube#videoListResponse",
 "etag": "\"Bdx4f4ps3xCOOo1WZ91nTLkRZ_c/ZwuIE09sDe2ZeuNsc3Rld0lEvWg\"",
 "pageInfo": {
  "totalResults": 2,
  "resultsPerPage": 2
 },
 "items": [
  {
   "kind": "youtube#video",
   "etag": "\"Bdx4f4ps3xCOOo1WZ91nTLkRZ_c/_oh7H4lj5Bqespe-aNJUMushCKQ\"",
   "id": "zO_B_SLhAJk",
   "snippet": {
    "publishedAt": "2019-08-01T13:30:02.000Z",
    "channelId": "UCqnbDFdCpuN8CMEg0VuEBqA",
    "title": "The Second 2019 Democratic Debate: Key Moments, Day 2 | NYT News",
    "description": "Former Vice President Joseph R. Biden Jr. and Senator Kamala Harris sparred while fending off attacks from fellow candidates on health care and criminal justice reform.\n\nRead more: https://nyti.ms/2YdYXn0\n\nSubscribe: http://bit.ly/U8Ys7n\nMore from The New York Times Video:  http://nytimes.com/video\n----------\nWhether it's reporting on conflicts abroad and political divisions at home, or covering the latest style trends and scientific developments, New York Times video journalists provide a revealing and unforgettable view of the world. It's all the news that's fit to watch.",
    "thumbnails": {
     "default": {
      "url": "https://i.ytimg.com/vi/zO_B_SLhAJk/default.jpg",
      "width": 120,
      "height": 90
     },
     "medium": {
      "url": "https://i.ytimg.com/vi/zO_B_SLhAJk/mqdefault.jpg",
      "width": 320,
      "height": 180
     },
     "high": {
      "url": "https://i.ytimg.com/vi/zO_B_SLhAJk/hqdefault.jpg",
      "width": 480,
      "height": 360
     },
     "standard": {
      "url": "https://i.ytimg.com/vi/zO_B_SLhAJk/sddefault.jpg",
      "width": 640,
      "height": 480
     },
     "maxres": {
      "url": "https://i.ytimg.com/vi/zO_B_SLhAJk/maxresdefault.jpg",
      "width": 1280,
      "height": 720
     }
    },
    "channelTitle": "The New York Times",
    "tags": [
     "debate day 2",
     "biden debate",
     "debate highlights",
     "kamala harris and joe biden",
     "bided kid",
     "biden called kamala a kid",
     "u.s. news",
     "debate news",
     "New York times",
     "nytimes highlights"
    ],
    "categoryId": "25",
    "liveBroadcastContent": "none",
    "localized": {
     "title": "The Second 2019 Democratic Debate: Key Moments, Day 2 | NYT News",
     "description": "Former Vice President Joseph R. Biden Jr. and Senator Kamala Harris sparred while fending off attacks from fellow candidates on health care and criminal justice reform.\n\nRead more: https://nyti.ms/2YdYXn0\n\nSubscribe: http://bit.ly/U8Ys7n\nMore from The New York Times Video:  http://nytimes.com/video\n----------\nWhether it's reporting on conflicts abroad and political divisions at home, or covering the latest style trends and scientific developments, New York Times video journalists provide a revealing and unforgettable view of the world. It's all the news that's fit to watch."
    },
    "defaultAudioLanguage": "en"
   },
   "contentDetails": {
    "duration": "PT2M46S",
    "dimension": "2d",
    "definition": "hd",
    "caption": "true",
    "licensedContent": true,
    "projection": "rectangular"
   },
   "status": {
    "uploadStatus": "processed",
    "privacyStatus": "public",
    "license": "youtube",
    "embeddable": true,
    "publicStatsViewable": false
   },
   "statistics": {
    "viewCount": "26290",
    "likeCount": "259",
    "dislikeCount": "172",
    "favoriteCount": "0",
    "commentCount": "208"
   }
  },
  {
   "kind": "youtube#video",
   "etag": "\"Bdx4f4ps3xCOOo1WZ91nTLkRZ_c/MAwSMC5npLUDmWxPQyhRc84MT2M\"",
   "id": "pcge_OHwfSc",
   "snippet": {
    "publishedAt": "2019-07-31T13:30:02.000Z",
    "channelId": "UCqnbDFdCpuN8CMEg0VuEBqA",
    "title": "The Second  2019 Democratic Debate: Key Moments, Day 1 | NYT News",
    "description": "The leading progressives, Senators Bernie Sanders and Elizabeth Warren, fended off attacks from underdog moderate challengers.\n\nRead more: https://nyti.ms/2GDUUpv\nSubscribe: http://bit.ly/U8Ys7n\nMore from The New York Times Video:  http://nytimes.com/video\n----------\nWhether it's reporting on conflicts abroad and political divisions at home, or covering the latest style trends and scientific developments, New York Times video journalists provide a revealing and unforgettable view of the world. It's all the news that's fit to watch.",
    "thumbnails": {
     "default": {
      "url": "https://i.ytimg.com/vi/pcge_OHwfSc/default.jpg",
      "width": 120,
      "height": 90
     },
     "medium": {
      "url": "https://i.ytimg.com/vi/pcge_OHwfSc/mqdefault.jpg",
      "width": 320,
      "height": 180
     },
     "high": {
      "url": "https://i.ytimg.com/vi/pcge_OHwfSc/hqdefault.jpg",
      "width": 480,
      "height": 360
     },
     "standard": {
      "url": "https://i.ytimg.com/vi/pcge_OHwfSc/sddefault.jpg",
      "width": 640,
      "height": 480
     },
     "maxres": {
      "url": "https://i.ytimg.com/vi/pcge_OHwfSc/maxresdefault.jpg",
      "width": 1280,
      "height": 720
     }
    },
    "channelTitle": "The New York Times",
    "tags": [
     "democrat debates",
     "debate highlights",
     "dem debates highlights",
     "warren",
     "sanders",
     "key moments democrats",
     "news",
     "new york times",
     "nytimes video",
     "video from the nytimes",
     "cnn democrat debate"
    ],
    "categoryId": "25",
    "liveBroadcastContent": "none",
    "localized": {
     "title": "The Second  2019 Democratic Debate: Key Moments, Day 1 | NYT News",
     "description": "The leading progressives, Senators Bernie Sanders and Elizabeth Warren, fended off attacks from underdog moderate challengers.\n\nRead more: https://nyti.ms/2GDUUpv\nSubscribe: http://bit.ly/U8Ys7n\nMore from The New York Times Video:  http://nytimes.com/video\n----------\nWhether it's reporting on conflicts abroad and political divisions at home, or covering the latest style trends and scientific developments, New York Times video journalists provide a revealing and unforgettable view of the world. It's all the news that's fit to watch."
    },
    "defaultAudioLanguage": "en"
   },
   "contentDetails": {
    "duration": "PT3M7S",
    "dimension": "2d",
    "definition": "hd",
    "caption": "true",
    "licensedContent": true,
    "projection": "rectangular"
   },
   "status": {
    "uploadStatus": "processed",
    "privacyStatus": "public",
    "license": "youtube",
    "embeddable": true,
    "publicStatsViewable": false
   },
   "statistics": {
    "viewCount": "242576",
    "likeCount": "1181",
    "dislikeCount": "1397",
    "favoriteCount": "0",
    "commentCount": "1425"
   }
  }
 ]
}
```

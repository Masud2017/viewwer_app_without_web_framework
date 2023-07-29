//package main
//
//import (
//	"encoding/json"
//	"errors"
//	"fmt"
//	"github.com/anaskhan96/soup"
//	"log"
//	"math/rand"
//	"regexp"
//	"strconv"
//	"strings"
//	"sync"
//	"time"
//)
//
//var (
//	requestCounter int64     = 0 // counter for request in TikTok
//	globalTime     time.Time     // it will be used to synchronise concurrent go routines
//)
//
//type VideoInfo struct {
//	ViewCount int
//	Platform  string
//	Username  string
//	Duration  string
//}
//
//func ProcessUrl(url string) VideoInfo {
//	var videoInfo VideoInfo
//
//	// portion of the code that checks whether the url from YouTube.
//	platform_youtube_pattern := ".*\\.youtube\\..*"
//	platform_youtube_pattern_compiled, _ := regexp.Compile(platform_youtube_pattern)
//
//	is_youtube := platform_youtube_pattern_compiled.MatchString(url)
//
//	if is_youtube {
//		videoInfo.Platform = "Youtube"
//	}
//
//	// portion of the code that cehcks whether the url from instagram
//	platform_intagram_pattern := ".*\\.instagram\\..*"
//	platform_instagram_pattern_compiled, _ := regexp.Compile(platform_intagram_pattern)
//
//	is_insta := platform_instagram_pattern_compiled.MatchString(url)
//
//	if is_insta {
//		videoInfo.Platform = "Instagram"
//	}
//
//	// portion of the code that cehcks whether the url from TikTok
//	platform_tiktok_pattern := ".*\\.tiktok\\..*"
//	platform_tiktok_pattern_compiled, _ := regexp.Compile(platform_tiktok_pattern)
//
//	is_tiktok := platform_tiktok_pattern_compiled.MatchString(url)
//
//	if is_tiktok {
//		videoInfo.Platform = "Tiktok"
//	}
//
//	return videoInfo
//}
//
//func getYoutubeVideoIdFromUrl(url string) string {
//	videoIdPattern, err := regexp.Compile("(v=.*)")
//	if err != nil {
//		log.Fatalln(err)
//	}
//
//	videoId := videoIdPattern.FindString(url)
//	videoId = videoId[2:]
//	if strings.Contains(videoId, "&") {
//		videoId = strings.Split(videoId, "&")[0]
//	}
//	fmt.Println(videoId)
//	return videoId
//}
//
//func getYoutubeVideoDuration(url string) string {
//	queryUrl := fmt.Sprintf("https://www.youtube.com/results?search_query=%s", getYoutubeVideoIdFromUrl(url))
//	soupObj, _ := soup.Get(queryUrl)
//
//	htmlContent := soup.HTMLParse(soupObj)
//	//fmt.Println(htmlContent.Find("title"))
//	htmlTextContent := htmlContent.FullText()
//	jsonFirstPattern, _ := regexp.Compile("{\"responseContext\"")
//	firstIdx := jsonFirstPattern.FindStringIndex(htmlTextContent)[0]
//
//	jsonSecondPattern, _ := regexp.Compile("\"targetId\":\"search-page\"};if")
//	matchArr := jsonSecondPattern.FindStringIndex(htmlTextContent)
//	secondIdx := matchArr[len(matchArr)-1]
//	data := htmlTextContent[firstIdx : secondIdx-3]
//
//	var jsonData map[string]interface{}
//
//	error := json.Unmarshal([]byte(data), &jsonData)
//	if error != nil {
//		fmt.Println(error)
//	}
//
//	youtubeVideoDuration := jsonData["contents"].(map[string]interface{})["twoColumnSearchResultsRenderer"].(map[string]interface{})["primaryContents"].(map[string]interface{})["sectionListRenderer"].(map[string]interface{})["contents"].([]interface{})[0].(map[string]interface{})["itemSectionRenderer"].(map[string]interface{})["contents"].([]interface{})[0].(map[string]interface{})["videoRenderer"].(map[string]interface{})["lengthText"].(map[string]interface{})["accessibility"].(map[string]interface{})["accessibilityData"].(map[string]interface{})["label"]
//	fmt.Println(youtubeVideoDuration)
//
//	return youtubeVideoDuration.(string)
//}
//
//func ScrapeYoutubeData(videoInfo *VideoInfo, url string) error {
//
//	soupObj, err := soup.Get(url)
//
//	if err != nil {
//		fmt.Println("An error happnd while trying get the url")
//		return errors.New("Error happening while trying to call \"soup.Get(url)\" ")
//	}
//
//	htmlContent := soup.HTMLParse(soupObj)
//
//	// video view
//	link := htmlContent.Find("meta", "itemprop", "interactionCount")
//	videoView := link.Attrs()["content"]
//
//	videoInfo.ViewCount, _ = strconv.Atoi(videoView)
//
//	// channel name
//	channelNameLink := htmlContent.Find("span", "itemprop", "author").Find("link", "itemprop", "name")
//	channelName := channelNameLink.Attrs()["content"]
//	videoInfo.Username = channelName
//
//	videoInfo.Duration = getYoutubeVideoDuration(url)
//
//	return nil
//}
//
//func GetTiktokVideoId(url string) string {
//	pattern := "\\/video\\/(\\w+)"
//	pattern_compiled, _ := regexp.Compile(pattern)
//	res := pattern_compiled.FindString(url)
//	videoId := strings.Split(res, "/")[2]
//
//	return videoId
//}
//
//func ScrapeTiktokData(videoInfo *VideoInfo, url string) error {
//	rand.Seed(time.Now().UnixNano())
//	soupObj, err := soup.Get(url)
//
//	if err != nil {
//		log.Fatalf("%s", err)
//	}
//
//	htmlContent := soup.HTMLParse(soupObj)
//	if htmlContent.Find("title").FullText() == "Access Denied" {
//		//waitTime := rand.Intn(20-5+1) + 1
//		log.Println("Access denied, So waiting some time and sending the request again 🙂 Waiting for ", 10, " second")
//
//		//time.Sleep(time.Second * time.Duration(waitTime))
//		time.Sleep(time.Second * 10)
//
//		ScrapeTiktokData(videoInfo, url)
//
//	}
//
//	content := htmlContent.Find("script", "id", "SIGI_STATE").FullText()
//
//	if len(content) > 0 {
//		var jsonData map[string]interface{}
//
//		error := json.Unmarshal([]byte(content), &jsonData)
//		if error != nil {
//			fmt.Println(error)
//		}
//
//		// fetching the view count
//		tiktokVideoId := GetTiktokVideoId(url)
//
//		rawStatData := jsonData["ItemModule"].(map[string]interface{})[tiktokVideoId].(map[string]interface{})["stats"]
//		rawViewCount := rawStatData.(map[string]interface{})["playCount"]
//		viewCount, _ := rawViewCount.(float64)
//
//		videoInfo.ViewCount = int(viewCount)
//
//		// fetching the channel name (in this case username)
//		rawAuthorData := jsonData["ItemModule"].(map[string]interface{})[tiktokVideoId].(map[string]interface{})["author"]
//		authorDataString, _ := rawAuthorData.(string)
//
//		videoInfo.Username = authorDataString
//
//		videoDuration := jsonData["ItemModule"].(map[string]interface{})[tiktokVideoId].(map[string]interface{})["video"].(map[string]interface{})["duration"].(float64)
//		videoDurationConverted := time.Duration(videoDuration * 1e9)
//
//		durationString := videoDurationConverted.String()
//
//		videoInfo.Duration = durationString
//	}
//
//	return nil
//}
//
//func GetData(url string) (VideoInfo, error) {
//	videoInfo := ProcessUrl(url) // this will only populate the platform field
//
//	switch videoInfo.Platform {
//	case "Youtube":
//		err := ScrapeYoutubeData(&videoInfo, url) // it will populate the videInfo with video data not going to return anything
//		return videoInfo, err
//
//	case "Instagram":
//		// For instagram thing
//
//	case "Tiktok":
//		err := ScrapeTiktokData(&videoInfo, url)
//		return videoInfo, err
//	}
//
//	return videoInfo, nil
//}
//func GetViewData(url string, wg *sync.WaitGroup, mutex *sync.Mutex) <-chan VideoInfoWithErr {
//	defer wg.Done()
//	videoInfoWithErrChannel := make(chan VideoInfoWithErr)
//
//	if requestCounter == 15 {
//		mutex.Lock()
//		// we need to make the wait random
//		waitTime := rand.Intn(20-5+1) + 1
//		fmt.Printf("The request was sent 15 times so waiting for %d second before starting again \n", waitTime)
//		globalTime = time.Now().Add(time.Second * time.Duration(waitTime))
//		mutex.Unlock()
//	}
//
//	if globalTime != (time.Time{}) {
//		currentTime := time.Now()
//
//		for {
//			if currentTime.After(globalTime) {
//				mutex.Lock()
//				globalTime = time.Time{} // cleaning the global time
//				requestCounter = 0
//				mutex.Unlock()
//				break
//			} else {
//				currentTime = time.Now()
//			}
//		}
//	}
//
//	go func() {
//		defer close(videoInfoWithErrChannel)
//		videoInfo, err := GetData(url)
//		videoInfoWithErrChannel <- VideoInfoWithErr{videoInfo, err}
//	}()
//
//	mutex.Lock()
//	//fmt.Println("Value of Counter is ; ", requestCounter)
//	requestCounter++
//	mutex.Unlock()
//	return videoInfoWithErrChannel
//}
//
//type VideoInfoWithErr struct {
//	VideoInfo VideoInfo
//	Err       error
//}
//
//// test code
//func main() {
//	//url := "https://www.tiktok.com/@islamic_way_for_muslim/video/7246968701021408514"
//	url := "https://www.youtube.com/watch?v=153ZT6tgZiw&t=1s&ab_channel=Zahed%27sTake"
//
//	var channelArr []<-chan VideoInfoWithErr
//
//	wg := sync.WaitGroup{}
//	mutex := sync.Mutex{}
//	wg.Add(100)
//
//	for x := 0; x < 100; x++ {
//		tempChannel := GetViewData(url, &wg, &mutex)
//
//		channelArr = append(channelArr, tempChannel)
//	}
//
//	for _, channelItem := range channelArr {
//		videoItem := <-channelItem
//		fmt.Println(videoItem.VideoInfo)
//	}
//	wg.Wait()
//
//}

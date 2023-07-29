package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/anaskhan96/soup"
	"log"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	requestCounter int64     = 0 // counter for request in TikTok
	globalTime     time.Time     // it will be used to synchronise concurrent go routines
)

//type VideoViewTrackingInfo struct {
//	Platform  string
//	Username  string
//	Url       string
//	Id        int
//	SourceId  int
//	AudioId   int
//	VoiceId   int
//	WriterId  int
//	SponsorId int
//}

type VideoViewTrackingInfo struct {
	Platform string
	Username string
	Url      string
	Id       int
}

type VideoMetrics struct {
	ViewCount    int
	Platform     string
	Username     string
	Duration     string
	LikeCount    float64
	CommentCount float64
	Title        string
}

func ProcessUrl(url string) VideoMetrics {
	var videoInfo VideoMetrics

	// portion of the code that checks whether the url from YouTube.
	platform_youtube_pattern := ".*\\.youtube\\..*"
	platform_youtube_pattern_compiled, _ := regexp.Compile(platform_youtube_pattern)

	is_youtube := platform_youtube_pattern_compiled.MatchString(url)

	if is_youtube {
		videoInfo.Platform = "Youtube"

		//fmt.Println(url)
		//channelNamePattern, _ := regexp.Compile(`(ab_channel=)([^&]+)`)
		//channelName := channelNamePattern.FindString(url)
		//// fmt.Println(tern)
		//
		//videoInfo.Username = channelName[11:]

	}

	// portion of the code that cehcks whether the url from instagram
	platform_intagram_pattern := ".*\\.instagram\\..*"
	platform_instagram_pattern_compiled, _ := regexp.Compile(platform_intagram_pattern)

	is_insta := platform_instagram_pattern_compiled.MatchString(url)

	if is_insta {
		videoInfo.Platform = "Instagram"
	}

	// portion of the code that cehcks whether the url from TikTok
	platform_tiktok_pattern := ".*\\.tiktok\\..*"
	platform_tiktok_pattern_compiled, _ := regexp.Compile(platform_tiktok_pattern)

	is_tiktok := platform_tiktok_pattern_compiled.MatchString(url)

	if is_tiktok {
		videoInfo.Platform = "Tiktok"
		//tiktokUserNamePattern, _ := regexp.Compile(`@([^\/]+)`)
		//regexResult := tiktokUserNamePattern.FindString(url)
		//videoInfo.Username = regexResult[1:]
	}

	return videoInfo
}

func getYoutubeVideoIdFromUrl(url string) string {
	videoIdPattern, err := regexp.Compile("(v=.*)")
	if err != nil {
		log.Fatalln("something went wrong while compiling regular expression for videoIdPatter", err)
	}

	videoId := videoIdPattern.FindString(url)
	videoId = videoId[2:]
	if strings.Contains(videoId, "&") {
		videoId = strings.Split(videoId, "&")[0]
	}
	//fmt.Println(videoId)
	return videoId
}

func getYoutubeVideoDuration(url string) string {
	queryUrl := fmt.Sprintf("https://www.youtube.com/results?search_query=%s", getYoutubeVideoIdFromUrl(url))
	soupObj, err := soup.Get(queryUrl)
	if err != nil {
		log.Fatalln("Some thing went wrong while requesting for the url that you have given. PLease check your internet connection or the url ", err)
	}

	htmlContent := soup.HTMLParse(soupObj)
	//fmt.Println(htmlContent.Find("title"))
	htmlTextContent := htmlContent.FullText()
	jsonFirstPattern, err := regexp.Compile("{\"responseContext\"")
	if err != nil {
		log.Fatalln("Something wrong while compiling the regular expression for responseContext json keyword ", err)
	}
	firstIdx := jsonFirstPattern.FindStringIndex(htmlTextContent)[0]

	jsonSecondPattern, err := regexp.Compile("\"targetId\":\"search-page\"};if")

	if err != nil {
		log.Fatalln("Something went wrong while compiling the regular expression for \"targetId\": \"search-page\"} ", err)

	}

	matchArr := jsonSecondPattern.FindStringIndex(htmlTextContent)
	secondIdx := matchArr[len(matchArr)-1]
	data := htmlTextContent[firstIdx : secondIdx-3]

	var jsonData map[string]interface{}

	error := json.Unmarshal([]byte(data), &jsonData)
	if error != nil {
		log.Fatalln("Something went wrong while unmarshalling the json data ", err)
	}

	if jsonData["contents"].(map[string]interface{})["twoColumnSearchResultsRenderer"].(map[string]interface{})["primaryContents"].(map[string]interface{})["sectionListRenderer"].(map[string]interface{})["contents"].([]interface{})[0].(map[string]interface{})["itemSectionRenderer"].(map[string]interface{})["contents"].([]interface{})[0].(map[string]interface{})["videoRenderer"] == nil {
		fmt.Println("😅 Didn't get the youtube video duration so retrying .. 😅")
		return getYoutubeVideoDuration(url)
	}

	youtubeVideoDuration := jsonData["contents"].(map[string]interface{})["twoColumnSearchResultsRenderer"].(map[string]interface{})["primaryContents"].(map[string]interface{})["sectionListRenderer"].(map[string]interface{})["contents"].([]interface{})[0].(map[string]interface{})["itemSectionRenderer"].(map[string]interface{})["contents"].([]interface{})[0].(map[string]interface{})["videoRenderer"].(map[string]interface{})["lengthText"].(map[string]interface{})["accessibility"].(map[string]interface{})["accessibilityData"].(map[string]interface{})["label"]
	fmt.Println(youtubeVideoDuration)

	return youtubeVideoDuration.(string)
}

func getYoutubeLikeCount(url string) float64 {
	soupObj, err := soup.Get(url)
	if err != nil {
		fmt.Println("Facing issue while requesting for the url that you have provided. Please check your internet connection or check the url ", err)
	}

	htmlContent := soup.HTMLParse(soupObj)
	htmlTextContent := htmlContent.FullText()
	likeCountPattern, _ := regexp.Compile("likeCount\":\"\\d*\"")
	likeCountString := likeCountPattern.FindString(htmlTextContent)
	likeCountString = strings.Split(likeCountString, ":")[1]

	likeCountStr := likeCountString[1 : len(likeCountString)-1]
	likeCount, _ := strconv.ParseFloat(likeCountStr, 64)

	return likeCount
}

func getApproximateValueFromYoutubeCommentCount(commentCountStr string) float64 {

	initialVal, _ := strconv.ParseFloat(commentCountStr[:len(commentCountStr)-1], 64)
	var approxVal float64

	switch commentCountStr[len(commentCountStr)-1:] {
	case "K":
		approxVal = initialVal * 1000
		break
	case "M":
		approxVal = initialVal * 1000000
		break
	case "B":
		approxVal = initialVal * 1000000000
		break
	case "T":
		approxVal = initialVal * 1000000000000
		break
	default:
		val, _ := strconv.ParseFloat(commentCountStr, 64)
		return val

	}

	return approxVal
}

func getYoutubeCommentCount(url string) float64 {
	soup.Header("accept-language", "en-US,en;q=0.9,en-GB;q=0.8")
	soupObj, err := soup.Get(url)
	if err != nil {
		fmt.Println(err)
	}

	htmlContent := soup.HTMLParse(soupObj)
	htmlTextContent := htmlContent.FullText()

	commentCountPattern, _ := regexp.Compile(`contextualInfo":{"runs":\[{"text":"(\d*)?(\w*)?"}]}`)
	commentCountJson := commentCountPattern.FindString(htmlTextContent)
	commentCountJson = commentCountJson[16:]

	var jsonData map[string]interface{}

	error := json.Unmarshal([]byte(commentCountJson), &jsonData)
	if error != nil {
		fmt.Println(error)
	}

	commentCountStr := jsonData["runs"].([]interface{})[0].(map[string]interface{})["text"]

	fmt.Println(commentCountStr)
	return getApproximateValueFromYoutubeCommentCount(commentCountStr.(string))
}

func ScrapeYoutubeData(videoInfo *VideoMetrics, url string) error {

	soupObj, err := soup.Get(url)

	if err != nil {
		fmt.Println("An error happend while trying get the url")
		return errors.New("Error happening while trying to call \"soup.Get(url)\". Please check your internet connection ! ")
	}

	htmlContent := soup.HTMLParse(soupObj)

	// video view
	link := htmlContent.Find("meta", "itemprop", "interactionCount")
	videoView := link.Attrs()["content"]

	videoInfo.ViewCount, _ = strconv.Atoi(videoView)

	// channel name
	// channelNameLink := htmlContent.Find("span", "itemprop", "author").Find("link", "itemprop", "name")
	// channelName := channelNameLink.Attrs()["content"]
	// videoInfo.Username = channelName

	videoInfo.Duration = getYoutubeVideoDuration(url)
	videoInfo.LikeCount = getYoutubeLikeCount(url)
	videoInfo.CommentCount = getYoutubeCommentCount(url)

	// getting youtube video title
	//titleLink := htmlContent.Find("title")
	//title := titleLink.Text()
	//videoInfo.Title = title

	//channelName, title, err := getTitleAndUsername()

	return nil
}

func GetTiktokVideoId(url string) string {
	pattern := "\\/video\\/(\\w+)"
	pattern_compiled, _ := regexp.Compile(pattern)
	res := pattern_compiled.FindString(url)
	videoId := strings.Split(res, "/")[2]

	return videoId
}

func getCaptionOfTiktok(url string) (string, error) {
	soupObj, err := soup.Get(url)

	if err != nil {
		//fmt.Println(""err)
		return "", errors.New("please check your internet connection or the url that you are provided")
	}

	htmlContent := soup.HTMLParse(soupObj)

	htmlTextContent := htmlContent.Find("meta", "property", "og:description")
	// IMPORTANT: check for nil ptr
	if htmlTextContent.Pointer == nil {
		return "", errors.New("cannot find the tiktok caption for url: " + url)
	}
	caption := htmlTextContent.Attrs()["content"]

	fmt.Println("Testing the value of tiktok caption ", caption)
	return caption, err
}

func ScrapeTiktokData(videoInfo *VideoMetrics, url string) error {
	rand.Seed(time.Now().UnixNano())
	soupObj, err := soup.Get(url)

	if err != nil {
		log.Fatalf("%s", err)
		return errors.New("please check your internet connection or the url that you are provided")
	}

	htmlContent := soup.HTMLParse(soupObj)
	if htmlContent.Find("title").FullText() == "Access Denied" {
		//waitTime := rand.Intn(20-5+1) + 1
		log.Println("Access denied, So waiting some time and sending the request again 🙂 Waiting for ", 10, " second")

		//time.Sleep(time.Second * time.Duration(waitTime))
		time.Sleep(time.Second * 10)

		ScrapeTiktokData(videoInfo, url)

	}

	content := htmlContent.Find("script", "id", "SIGI_STATE")
	if content.Pointer == nil {
		return errors.New("nil pointer in search for SIGI_STATE")
	}
	if len(content.FullText()) > 0 {
		var jsonData map[string]interface{}

		error := json.Unmarshal([]byte(content.FullText()), &jsonData)
		if error != nil {
			return errors.New("got error while unmarshalling the json data ")
		}

		// fetching the view count
		tiktokVideoId := GetTiktokVideoId(url)

		lvl1, ok := jsonData["ItemModule"].(map[string]interface{})
		if !ok {
			return errors.New("cannot find the ItemModule")
		}

		lvl2, ok := lvl1[tiktokVideoId].(map[string]interface{})
		if !ok {
			return errors.New("cannot find the tiktok video id")
		}

		rawStatData := lvl2["stats"].(map[string]interface{})
		rawViewCount := rawStatData["playCount"]
		viewCount, _ := rawViewCount.(float64)

		videoInfo.ViewCount = int(viewCount)

		// fetching the channel name (in this case username)
		// rawAuthorData := jsonData["ItemModule"].(map[string]interface{})[tiktokVideoId].(map[string]interface{})["author"]
		// authorDataString, _ := rawAuthorData.(string)

		// videoInfo.Username = authorDataString

		videoDuration := jsonData["ItemModule"].(map[string]interface{})[tiktokVideoId].(map[string]interface{})["video"].(map[string]interface{})["duration"].(float64)
		videoDurationConverted := time.Duration(videoDuration * 1e9)

		durationString := videoDurationConverted.String()

		videoInfo.Duration = durationString

		videoInfo.LikeCount = rawStatData["diggCount"].(float64)
		videoInfo.CommentCount = rawStatData["commentCount"].(float64)
		caption, err := getCaptionOfTiktok(url)
		if err != nil {
			fmt.Println(err)
			errorMessage := fmt.Sprintf("% %", "got error while trying to get captoin of tiktok ", err)
			return errors.New(errorMessage)
		}
		videoInfo.Title = caption
	}

	return nil
}

/*
*
@param url - string - url string of the video
@param platform - string - Platform name of the video : eg. Youtube,Tiktok,Instagram
@returns title, username, and error
*/
func getTitleAndUsername(url string, platform string) (string, string, error) {
	soup.Header("accept-language", "en-US,en;q=0.9,en-GB;q=0.8")
	var userName string
	var title string
	soupObj, err := soup.Get(url)
	if err != nil {
		log.Println(err)

	}
	htmlContent := soup.HTMLParse(soupObj)
	switch platform {
	case "Youtube":
		channelNameLink := htmlContent.Find("span", "itemprop", "author").Find("link", "itemprop", "name")
		channelName := channelNameLink.Attrs()["content"]
		userName = channelName

		titleLink := htmlContent.Find("title")
		videoTitle := titleLink.Text()
		title = videoTitle

		return title, userName, nil
		break
	case "Tiktok":
		// fetching the channel name (in this case username)
		tiktokUserNamePattern, _ := regexp.Compile(`@([^\/]+)`)
		regexResult := tiktokUserNamePattern.FindString(url)
		userName = regexResult[1:]

		caption, err := getCaptionOfTiktok(url)
		if err != nil {
			log.Println(err)
		}

		title = caption
		return title, userName, nil
		break
	}

	return "", "", nil
}

func GetData(url string) (VideoMetrics, error) {
	videoInfo := ProcessUrl(url) // this will only populate the platform field

	switch videoInfo.Platform {
	case "Youtube":
		err := ScrapeYoutubeData(&videoInfo, url) // it will populate the videInfo with video data not going to return anything
		return videoInfo, err

	case "Instagram":
		// For instagram thing

	case "Tiktok":
		err := ScrapeTiktokData(&videoInfo, url)
		return videoInfo, err
	}

	return videoInfo, nil
}
func GetViewData(q VideoViewTrackingInfo, wg *sync.WaitGroup, mutex *sync.Mutex) <-chan VideoInfoWithErr {
	defer wg.Done()
	videoInfoWithErrChannel := make(chan VideoInfoWithErr)

	if requestCounter == 15 {
		mutex.Lock()
		// we need to make the wait random
		waitTime := rand.Intn(20-5+1) + 1
		fmt.Printf("The request was sent 15 times so waiting for %d second before starting again \n", waitTime)
		globalTime = time.Now().Add(time.Second * time.Duration(waitTime))
		mutex.Unlock()
	}

	if globalTime != (time.Time{}) {
		currentTime := time.Now()

		for {
			if currentTime.After(globalTime) {
				mutex.Lock()
				globalTime = time.Time{} // cleaning the global time
				requestCounter = 0
				mutex.Unlock()
				break
			} else {
				currentTime = time.Now()
			}
		}
	}

	go func() {
		defer close(videoInfoWithErrChannel)
		videoMetrics, err := GetData(q.Url)
		videoInfoWithErrChannel <- VideoInfoWithErr{videoMetrics, q, err}
	}()

	mutex.Lock()
	//fmt.Println("Value of Counter is ; ", requestCounter)
	requestCounter++
	mutex.Unlock()
	return videoInfoWithErrChannel
}

type VideoInfoWithErr struct {
	VideoMetrics VideoMetrics
	VideoInfo    VideoViewTrackingInfo
	Err          error
}

// test code
func main() {
	//url := "https://www.tiktok.com/@prince_abdullah_1_2_3/video/7215628737218497794"
	url := "https://www.youtube.com/watch?v=2RERHdL0ZWY&ab_channel=ChamokHasan"
	//	url := "https://www.tiktok.com/@reddit.guy/video/7245740863991663914"

	var channelArr []<-chan VideoInfoWithErr

	wg := sync.WaitGroup{}
	mutex := sync.Mutex{}
	wg.Add(1)

	for x := 0; x < 1; x++ {
		videoViewTrackingInfo := VideoViewTrackingInfo{Url: url}

		tempChannel := GetViewData(videoViewTrackingInfo, &wg, &mutex)

		channelArr = append(channelArr, tempChannel)
	}

	for _, channelItem := range channelArr {
		videoItem := <-channelItem
		fmt.Println(videoItem.VideoMetrics)
	}
	wg.Wait()
	//
	//title, userName, err := getTitleAndUsername(url, "Tiktok")
	//if err != nil {
	//	log.Println(err)
	//}
	//
	//fmt.Printf("Username : %s ; title : %s\n", userName, title)

}

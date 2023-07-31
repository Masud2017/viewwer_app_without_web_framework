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

const (
	pageNotFound                        = 404
	nullPointer                         = 700 // custom error code will start from 700
	accessDenied                        = 701
	invalidUriOrInternetConnectionIssue = 702
	allOk                               = 0
)

var errorCodeWithMessageMap = map[int]string{
	404: "Page not found please check your url and try again",
	700: "Found pointer while accessing the node",
	701: "Access defined please check your url or wait for some time before retry",
	702: "The url you have provide might be invalid or you have network connectivity issue.",
	0:   "Everything is ok",
}

// VideoInfo are information about the video that do not change
type VideoInfo struct {
	Url      string
	Platform string
	Username string
	Caption  string
	// Duration in seconds
	Duration int
}

// VideoMetrics are information about the video that change over time
type VideoMetrics struct {
	ViewCount    int
	LikeCount    float64
	CommentCount float64
}

type VideoAnalyticsError struct {
	Code    int
	Message string
}

// VideoInfoResponse is a combination of VideoInfo and VideoMetrics, as well as any error that was returned while fetching the data
type VideoInfoResponse struct {
	VideoInfo
	VideoMetrics
	Err VideoAnalyticsError
}

type VideoAnalyticsProvider interface {
	GetVideoInfo(url string) (VideoInfo, error)

	// GetFullVideoMetrics returns a channel of VideoInfoResponse, which contains the populated VideoInfo and VideoMetrics with information from the respective service (tiktok, instagram, etc)
	GetFullVideoMetrics(info VideoInfo, wg *sync.WaitGroup, mutex *sync.Mutex) <-chan VideoInfoResponse
}

func ProcessUrl(url string) VideoInfo {
	var videoInfo VideoInfo

	// portion of the code that checks whether the url from YouTube.
	platform_youtube_pattern := ".*\\.youtube\\..*"
	platform_youtube_pattern_compiled, _ := regexp.Compile(platform_youtube_pattern)

	is_youtube := platform_youtube_pattern_compiled.MatchString(url)

	if is_youtube {
		videoInfo.Platform = "Youtube"
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

	// checking for nil pointer issue for each item in the json
	if jsonData == nil {
		log.Fatalln("Json Data is nil so retrying again.")
		return getYoutubeVideoDuration(url)
	}
	if jsonData["contents"] == nil {
		log.Fatal("Content keyword is not available in the jsonData so retrying again.")
		return getYoutubeVideoDuration(url)
	}
	if jsonData["contents"].(map[string]interface{})["twoColumnSearchResultsRenderer"] == nil {
		log.Fatal("twoColumnSearchResultsRenderer keyword is not available so trying again.")
		return getYoutubeVideoDuration(url)
	}

	if jsonData["contents"].(map[string]interface{})["twoColumnSearchResultsRenderer"].(map[string]interface{})["primaryContents"] == nil {
		log.Fatal("primaryContents keyword is not found so retrying again.")
		return getYoutubeVideoDuration(url)
	}

	if jsonData["contents"].(map[string]interface{})["twoColumnSearchResultsRenderer"].(map[string]interface{})["primaryContents"].(map[string]interface{})["sectionListRenderer"] == nil {
		log.Fatalln("SectionListRenderer keyword is not found so retrying again.")
		return getYoutubeVideoDuration(url)
	}

	if jsonData["contents"].(map[string]interface{})["twoColumnSearchResultsRenderer"].(map[string]interface{})["primaryContents"].(map[string]interface{})["sectionListRenderer"].(map[string]interface{})["contents"].([]interface{})[0].(map[string]interface{})["itemSectionRenderer"].(map[string]interface{})["contents"] == nil {
		log.Fatal("contents keyword is not found so retrying again.")
		return getYoutubeVideoDuration(url)
	}

	if jsonData["contents"].(map[string]interface{})["twoColumnSearchResultsRenderer"].(map[string]interface{})["primaryContents"].(map[string]interface{})["sectionListRenderer"].(map[string]interface{})["contents"].([]interface{})[0].(map[string]interface{})["itemSectionRenderer"].(map[string]interface{})["contents"].([]interface{})[0].(map[string]interface{})["videoRenderer"] == nil {
		fmt.Println("😅 Didn't get the youtube video duration so retrying .. 😅")
		return getYoutubeVideoDuration(url)
	}

	youtubeVideoDuration := jsonData["contents"].(map[string]interface{})["twoColumnSearchResultsRenderer"].(map[string]interface{})["primaryContents"].(map[string]interface{})["sectionListRenderer"].(map[string]interface{})["contents"].([]interface{})[0].(map[string]interface{})["itemSectionRenderer"].(map[string]interface{})["contents"].([]interface{})[0].(map[string]interface{})["videoRenderer"].(map[string]interface{})["lengthText"].(map[string]interface{})["accessibility"].(map[string]interface{})["accessibilityData"].(map[string]interface{})["label"]
	fmt.Println(youtubeVideoDuration)

	return youtubeVideoDuration.(string)
}
func getIntFromYoutubeDuration(duration string) int {

	minute := strings.Split(duration, ",")[0]
	second := strings.Split(duration, ",")[1]
	minute_digit := strings.Split(minute, " ")[0]
	second_digit := strings.Split(second, " ")[1]

	minute_in_second, _ := strconv.Atoi(minute_digit)
	minute_in_second = minute_in_second * 60
	second_digit_num, _ := strconv.Atoi(second_digit)
	//fmt.Println(second_digit)

	return minute_in_second + second_digit_num
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

	commentCountPattern, err := regexp.Compile(`contextualInfo":{"runs":\[{"text":"(\d*)?(\w*)?"}]}`)

	if err != nil {
		log.Fatalln("facing issue while compiling the regular expression for contextualInfo ", err)
	}

	commentCountJson := commentCountPattern.FindString(htmlTextContent)
	commentCountJson = commentCountJson[16:]

	var jsonData map[string]interface{}

	error := json.Unmarshal([]byte(commentCountJson), &jsonData)
	if error != nil {
		fmt.Println(error)
	}

	if jsonData["runs"] == nil {
		log.Fatalln("jsonData[\"runs\"] is null so retrying again .")
		return getYoutubeCommentCount(url)
	}

	if jsonData["runs"].([]interface{})[0].(map[string]interface{})["text"] != nil {
		log.Fatalln("jsonData[\"runs\"].([]interface{})[0].(map[string]interface{})[\"text\"] is null so retrying again.")

		return getYoutubeCommentCount(url)
	}

	commentCountStr := jsonData["runs"].([]interface{})[0].(map[string]interface{})["text"]

	//fmt.Println(commentCountStr)
	return getApproximateValueFromYoutubeCommentCount(commentCountStr.(string))
}

func ScrapeYoutubeData(url string) (VideoMetrics, VideoAnalyticsError) {
	var videoMetrics VideoMetrics

	soupObj, err := soup.Get(url)

	if err != nil {
		log.Fatalln("An error happened while trying get the url")
		//return VideoMetrics{}, errors.New("Error happening while trying to call \"soup.Get(url)\". Please check your internet connection ! ")
		return VideoMetrics{}, VideoAnalyticsError{invalidUriOrInternetConnectionIssue, errorCodeWithMessageMap[invalidUriOrInternetConnectionIssue]}
	}

	htmlContent := soup.HTMLParse(soupObj)

	// video view
	link := htmlContent.Find("meta", "itemprop", "interactionCount")
	videoView := link.Attrs()["content"]

	videoMetrics.ViewCount, _ = strconv.Atoi(videoView)

	videoMetrics.LikeCount = getYoutubeLikeCount(url)
	videoMetrics.CommentCount = getYoutubeCommentCount(url)

	return videoMetrics, VideoAnalyticsError{allOk, errorCodeWithMessageMap[allOk]}
}

func GetTiktokVideoId(url string) string {
	pattern := "\\/video\\/(\\w+)"
	pattern_compiled, err := regexp.Compile(pattern)

	if err != nil {
		log.Fatalln("Error found while compiling the regular expression for tiktok video id ", err)
	}

	res := pattern_compiled.FindString(url)
	videoId := strings.Split(res, "/")[2]

	return videoId
}

func getCaptionOfTiktok(url string) (string, error) {
	soupObj, err := soup.Get(url)

	if err != nil {
		//fmt.Println(""err)
		return "", errors.New("please check your internet connection or the url that you are provided ")
	}

	htmlContent := soup.HTMLParse(soupObj)

	htmlTextContent := htmlContent.Find("meta", "property", "og:description")
	// IMPORTANT: check for nil ptr
	if htmlTextContent.Pointer == nil {
		return "", errors.New("cannot find the tiktok caption for url: " + url)
	}
	caption := htmlTextContent.Attrs()["content"]

	//fmt.Println("Testing the value of tiktok caption ", caption)
	return caption, err
}

func GetTiktokVideoDuration(url string) (int, error) {
	soupObj, _ := soup.Get(url)

	htmlContent := soup.HTMLParse(soupObj)

	if htmlContent.Find("title").FullText() == "Access Denied" {
		//waitTime := rand.Intn(20-5+1) + 1
		log.Println("Access denied, So waiting some time and sending the request again 🙂 Waiting for ", 10, " second")

		//time.Sleep(time.Second * time.Duration(waitTime))
		time.Sleep(time.Second * 10)

		return GetTiktokVideoDuration(url)

	}

	content := htmlContent.Find("script", "id", "SIGI_STATE")
	if content.Pointer == nil {
		return 0, errors.New("nil pointer in search for SIGI_STATE")
	}
	if len(content.FullText()) > 0 {
		var jsonData map[string]interface{}

		error := json.Unmarshal([]byte(content.FullText()), &jsonData)
		if error != nil {
			return 0, errors.New("got error while unmarshalling the json data ")
		}

		// fetching the view count
		tiktokVideoId := GetTiktokVideoId(url)

		if jsonData == nil {
			log.Fatalln("jsonData is nil so retrying again.")
			return GetTiktokVideoDuration(url)
		}

		if jsonData["ItemModule"] == nil {
			log.Fatalln("jsonData[\"ItemModule\"] is null so retrying again")
			return GetTiktokVideoDuration(url)
		}

		if jsonData["ItemModule"].(map[string]interface{})[tiktokVideoId] == nil {
			log.Fatalln("jsonData[\"ItemModule\"].(map[string]interface{})[tiktokVideoId] is null so retrying agian ")
			return GetTiktokVideoDuration(url)
		}

		if jsonData["ItemModule"].(map[string]interface{})[tiktokVideoId].(map[string]interface{})["video"] == nil {
			log.Fatalln("jsonData[\"ItemModule\"].(map[string]interface{})[tiktokVideoId].(map[string]interface{})[\"video\"] is null so retrying again")
			return GetTiktokVideoDuration(url)
		}

		if jsonData["ItemModule"].(map[string]interface{})[tiktokVideoId].(map[string]interface{})["video"].(map[string]interface{})["duration"] == nil {
			log.Fatalln("jsonData[\"ItemModule\"].(map[string]interface{})[tiktokVideoId].(map[string]interface{})[\"video\"].(map[string]interface{})[\"duration\"] is null so retrying again")
			return GetTiktokVideoDuration(url)
		}

		videoDuration := jsonData["ItemModule"].(map[string]interface{})[tiktokVideoId].(map[string]interface{})["video"].(map[string]interface{})["duration"].(float64)
		videoDurationConverted := time.Duration(videoDuration * 1e9)

		durationString := videoDurationConverted.String()

		durationFiltered := durationString[:len(durationString)-1]
		durationInt, _ := strconv.Atoi(durationFiltered)
		return durationInt, nil
	}
	return 0, nil
}

func ScrapeTiktokData(url string) (VideoMetrics, VideoAnalyticsError) {
	var videoMetrics VideoMetrics
	rand.Seed(time.Now().UnixNano())
	soupObj, err := soup.Get(url)

	if err != nil {
		log.Fatalf("%s", err)
		//return VideoMetrics{}, errors.New("please check your internet connection or the url that you are provided")

		return VideoMetrics{}, VideoAnalyticsError{invalidUriOrInternetConnectionIssue, errorCodeWithMessageMap[invalidUriOrInternetConnectionIssue]}
	}

	htmlContent := soup.HTMLParse(soupObj)
	if htmlContent.Find("title").FullText() == "Access Denied" {
		//waitTime := rand.Intn(20-5+1) + 1
		log.Println("Access denied, So waiting some time and sending the request again 🙂 Waiting for ", 10, " second")

		//time.Sleep(time.Second * time.Duration(waitTime))
		time.Sleep(time.Second * 10)

		ScrapeTiktokData(url)

	}

	content := htmlContent.Find("script", "id", "SIGI_STATE")
	if content.Pointer == nil {
		//return VideoMetrics{}, errors.New("nil pointer in search for SIGI_STATE")
		return VideoMetrics{}, VideoAnalyticsError{nullPointer, errorCodeWithMessageMap[nullPointer]}
	}
	if len(content.FullText()) > 0 {
		var jsonData map[string]interface{}

		error := json.Unmarshal([]byte(content.FullText()), &jsonData)
		if error != nil {
			//return VideoMetrics{}, errors.New("got error while unmarshalling the json data ")
			return VideoMetrics{}, VideoAnalyticsError{nullPointer, errorCodeWithMessageMap[nullPointer]}
		}

		// fetching the view count
		tiktokVideoId := GetTiktokVideoId(url)

		lvl1, ok := jsonData["ItemModule"].(map[string]interface{})
		if !ok {
			//return VideoMetrics{}, errors.New("cannot find the ItemModule")
			return VideoMetrics{}, VideoAnalyticsError{nullPointer, errorCodeWithMessageMap[nullPointer]}
		}

		lvl2, ok := lvl1[tiktokVideoId].(map[string]interface{})
		if !ok {
			//return VideoMetrics{}, errors.New("cannot find the tiktok video id")
			return VideoMetrics{}, VideoAnalyticsError{nullPointer, errorCodeWithMessageMap[nullPointer]}
		}

		rawStatData := lvl2["stats"].(map[string]interface{})
		rawViewCount := rawStatData["playCount"]
		viewCount, _ := rawViewCount.(float64)

		videoMetrics.ViewCount = int(viewCount)

		videoMetrics.LikeCount = rawStatData["diggCount"].(float64)
		videoMetrics.CommentCount = rawStatData["commentCount"].(float64)
	}

	return VideoMetrics{}, VideoAnalyticsError{}
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

func GetData(url string) VideoInfoResponse {
	videoInfo := ProcessUrl(url) // this will only populate the platform field

	switch videoInfo.Platform {
	case "Youtube":
		videoMetrics, err := ScrapeYoutubeData(url) // it will populate the videInfo with video data not going to return anything
		return VideoInfoResponse{VideoInfo{}, videoMetrics, err}

	case "Instagram":
		// For instagram thing

	case "Tiktok":
		videoMetrics, err := ScrapeTiktokData(url)
		return VideoInfoResponse{VideoInfo{}, videoMetrics, err}
	}

	return VideoInfoResponse{VideoInfo{}, VideoMetrics{}, VideoAnalyticsError{}}
}

type VideoAnalyticsProviderImpl struct {
	ErrorCode    int
	ErrorMessage string
}

func (videoAnalyticsProviderImpl *VideoAnalyticsProviderImpl) GetVideoInfo(url string) (VideoInfo, error) {
	var videoInfo VideoInfo = ProcessUrl(url)
	videoInfo.Url = url

	title, username, _ := getTitleAndUsername(url, videoInfo.Platform)

	videoInfo.Caption = title
	videoInfo.Username = username

	switch videoInfo.Platform {
	case "Youtube":
		videoInfo.Duration = getIntFromYoutubeDuration(getYoutubeVideoDuration(url))
		break
	case "Tiktok":
		duration, _ := GetTiktokVideoDuration(url)
		videoInfo.Duration = duration
		break
	}

	return videoInfo, nil
}

func (videoAnalyticsProviderImpl *VideoAnalyticsProviderImpl) GetFullVideoMetrics(info VideoInfo, wg *sync.WaitGroup, mutex *sync.Mutex) <-chan VideoInfoResponse {
	defer wg.Done()
	videoInfoWithErrChannel := make(chan VideoInfoResponse)

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
		videoInfoResponse := GetData(info.Url)
		videoInfoResponse.VideoInfo = info
		videoInfoWithErrChannel <- videoInfoResponse
	}()

	mutex.Lock()
	//fmt.Println("Value of Counter is ; ", requestCounter)
	requestCounter++
	mutex.Unlock()
	return videoInfoWithErrChannel
}

// test code
func main() {
	//url := "https://www.tiktok.com/@prince_abdullah_1_2_3/video/7215628737218497794"
	url := "https://www.youtube.com/watch?v=2RERHdL0ZWY&ab_channel=ChamokHasan"
	//	url := "https://www.tiktok.com/@reddit.guy/video/7245740863991663914"

	var channelArr []<-chan VideoInfoResponse
	var videoAnalyticsProvider VideoAnalyticsProviderImpl

	wg := sync.WaitGroup{}
	mutex := sync.Mutex{}
	wg.Add(1)

	for x := 0; x < 1; x++ {
		videoInfo, _ := videoAnalyticsProvider.GetVideoInfo(url)

		tempChannel := videoAnalyticsProvider.GetFullVideoMetrics(videoInfo, &wg, &mutex)

		channelArr = append(channelArr, tempChannel)
	}

	for _, channelItem := range channelArr {
		videoItem := <-channelItem
		fmt.Println(videoItem)
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

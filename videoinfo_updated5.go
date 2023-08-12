package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/anaskhan96/soup"
	"log"
	"math/rand"
	"os"
	"os/exec"
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
	allOk = 0

	pageNotFound                        = 404
	nullPointer                         = 700 // custom error code will start from 700
	accessDenied                        = 701
	invalidUriOrInternetConnectionIssue = 702
	jsonUnmarshallIssue                 = 703
)

var errorCodeWithMessageMap = map[int]string{
	404: "Page not found please check your url and try again",
	700: "Found null pointer while accessing the node",
	701: "Access defined please check your url or wait for some time before retry",
	702: "The url you have provide might be invalid or you have network connectivity issue.",
	0:   "Everything is ok",
	703: "Facing issue while unmarshalling json data.",
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
	if htmlContent.Pointer == nil {
		log.Fatalln("Youtube comment count data not found so retrying again")
		return getYoutubeCommentCount(url)
	}
	htmlTextContent := htmlContent.FullText()

	commentCountPattern, err := regexp.Compile(`contextualInfo":{"runs":\[{"text":"(\d*)?(\w*)?"}]}`)

	if err != nil {
		log.Fatalln("facing issue while compiling the regular expression for contextualInfo ", err)
	}

	commentCountJson := commentCountPattern.FindString(htmlTextContent)
	if len(commentCountJson) == 0 {
		//log.Fatalln("comment count json data not found so retrying ")
		return getApproximateValueFromYoutubeCommentCount("0")
	}
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

	if jsonData["runs"].([]interface{})[0].(map[string]interface{})["text"] == nil {
		log.Fatalln("jsonData[\"runs\"].([]interface{})[0].(map[string]interface{})[\"text\"] is null so retrying again.")

		return getYoutubeCommentCount(url)
	}

	commentCountStr := jsonData["runs"].([]interface{})[0].(map[string]interface{})["text"]

	//fmt.Println(commentCountStr)
	return getApproximateValueFromYoutubeCommentCount(commentCountStr.(string))
}

func ScrapeYoutubeData(url string) (VideoMetrics, VideoAnalyticsError) {
	// setting up the proxy
	os.Setenv("HTTP_PROXY", "http://someip:someport")
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
			return VideoMetrics{}, VideoAnalyticsError{jsonUnmarshallIssue, errorCodeWithMessageMap[jsonUnmarshallIssue]}
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

func extractMediaID(url string) (string, error) {
	re := regexp.MustCompile(`/([a-zA-Z0-9_-]+)\/?$`)
	match := re.FindStringSubmatch(url)
	if len(match) < 2 {
		return "", fmt.Errorf("could not extract media ID")
	}
	return match[1], nil
}

func shortcodeToMediaID(shortcode string) int {
	alphabet := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_"
	mediaID := 0

	for _, letter := range shortcode {
		mediaID = (mediaID * 64) + strings.IndexRune(alphabet, letter)
	}

	return mediaID
}

func ScrapeInstagram(url string, cookie string) (VideoMetrics, VideoAnalyticsError) {
	//shortCode, _ := extractMediaID(url)
	//mediaid := shortcodeToMediaID(shortCode)
	//
	//// setting up the headers to the soup
	//soup.Header("authority", "www.instagram.com")
	//soup.Header("method", "GET")
	//path := fmt.Sprintf("/api/v1/media/%d/info/", mediaid)
	//fmt.Println(path)
	//soup.Header("path", "/api/v1/media/3033137866545317138/info/")
	//soup.Header("scheme", "https")
	//soup.Header("accept", "*/*")
	//soup.Header("accept-encoding", "gzip, deflate, br")
	//soup.Header("accept-language", "en-US,en;q=0.9")
	//soup.Header("cookie", "mid=ZIlxKgALAAFScymuGZ63XopJbgmQ; ig_did=3D7F7CB4-D713-4AC0-B289-1CB257D7146C; ig_nrcb=1; datr=KHGJZGAHYNaqVNT8a_ULEH_J; fbm_124024574287414=base_domain=.instagram.com; ds_user_id=6428581218; csrftoken=gIR0LRMzGwtZBswOMsjUFlvTtIYvHD8b; shbid=\"3375\\0546428581218\\0541722876021:01f73ddcd6511a5f2ae24cb30cd970559f9ee04ddd2c67e1d86cba68c3d1af338a9b9200\"; shbts=\"1691340021\\0546428581218\\0541722876021:01f73e78321dfa01034452fefe26ceed79a8d4c6fa1cac2520f0ce956f70dcf86433a578\"; sessionid=6428581218%3AA5GFgB30EGGQFW%3A22%3AAYfA3idt6kMo1YUWRLruKZsif8s1QKN0Qk54Ib3quTE; fbsr_124024574287414=lU0NJQRAp26tlgSlWEmuwDoqtz93XJWm7LscH6tRaB8.eyJ1c2VyX2lkIjoiMTAwMDA1ODY3NjEwNzE1IiwiY29kZSI6IkFRRGR4SXpZLV9rbjFzdnV1emUzb0VwS2RZOVFBSWRTaTVaMGc3QTlzM0hIakZYQ3lxbFFwZmFQcWdZM0hWRDdVaEREYl9lR1JuMlFXVmhvZFJMZG9JSnloWHpqekhCd3BfQW9oUzV3X3d2bV9BVnVleW5rd0s3cUZzeWhFc3pHWk1FSk9CSUJvMDc3VkN3SHdIMDYzRmV1ZWkzOUZwYVJGaE9FMVlJelVpY1NQcFpBXzF6MS01YXg0OXlUUUN2UTIwWVRJUEp6eFROTGpIeXo4RDVHeWhZR3BiMzVscURZc2M1dXh0MTFFX1VRYUo2RnRqV0hRVVZkVWoxcmRWak1GeEJLdnNydVYtT1Rkb1ZrbnlCODZvU1V5NVV4LWRaLV9hcm0ySlctd3NFOVJQLWlkbWVDS1lMY2JsTE10dmowUGtjcWx1T2dxTGNnWUpmMUVKUjM3NTMtIiwib2F1dGhfdG9rZW4iOiJFQUFCd3pMaXhuallCTzF2enNqWGJzRzZ5TDVEbXVMTkNWVHFzM0ViblpDQkJCOXdublZCT1QzWkFFVndhRzBaQUdhanVVMVBMWkM3N2NxTGFCY25kSWZtTXA2T29ITFJQek1aQ1lTZlcydExNNDdaQUlLaW50T3JqZVFrTmdrcWhGcjV5dDRYbFIzdTEzRlVaQmZUaFZ6ZVZLUU1EQnVLSXhUMHNtZ1JzekhlVVd2NkMxd09ReVdNc2l1MHdmUm5aQTBNZ3F1Z1pEIiwiYWxnb3JpdGhtIjoiSE1BQy1TSEEyNTYiLCJpc3N1ZWRfYXQiOjE2OTEzNDAzNjd9; rur=\"PRN\\0546428581218\\0541722876370:01f703c1d3cc873a19efcd7747e238391f399dcb59297f69f30accb544200f02748f852f\"")
	//soup.Header("referer", "https://www.instagram.com/reel/CoX3tHCOEUS/")
	//soup.Header("sec-ch-prefers-color-scheme", "light")
	//soup.Header("sec-ch-ua", "\"Not?A_Brand\";v=\"8\", \"Chromium\";v=\"108\"")
	//soup.Header("sec-ch-ua-mobile", "?0")
	//soup.Header("sec-ch-ua-platform", "\"Windows\"")
	//soup.Header("sec-fetch-dest", "empty")
	//soup.Header("sec-fetch-mode", "cors")
	//soup.Header("sec-fetch-site", "same-origin")
	//soup.Header("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36")
	//soup.Header("viewport-width", "1280")
	//soup.Header("x-asbd-id", "198387")
	//soup.Header("x-csrftoken", "sNEl62eIgXH7oKUtgNU0WbBGfGsPY01v")
	//soup.Header("x-ig-app-id", "936619743392459")
	//soup.Header("x-ig-www-claim", "0")
	//soup.Header("x-requested-with", "XMLHttpRequest")
	//
	////targetUrl := fmt.Sprintf("https://www.instagram.com/api/v1/media/%d/info/", mediaid)
	//soupObj, err := soup.Get("https://www.instagram.com/api/v1/media/3033137866545317138/info/")
	//if err != nil {
	//	fmt.Println("Something went wrong while trying to get the url")
	//}
	//
	////htmlContent := soup.HTMLParse(soupObj)
	////
	//fmt.Println(soupObj)

	cmd := exec.Command("python", "python_solution.py")
	out, err := cmd.Output()
	fmt.Println(err)
	jsonOutput := string(out)

	var jsonData map[string]interface{}
	error := json.Unmarshal([]byte(jsonOutput), &jsonData)

	if error != nil {
		fmt.Println(error)
	}

	playCount := jsonData["items"].([]interface{})[0].(map[string]interface{})["play_count"].(float64)
	likeCount := jsonData["items"].([]interface{})[0].(map[string]interface{})["like_count"].(float64)
	commentCount := jsonData["items"].([]interface{})[0].(map[string]interface{})["comment_count"].(float64)
	videoMetrics := VideoMetrics{ViewCount: int(playCount), LikeCount: likeCount, CommentCount: commentCount}
	fmt.Println(playCount)

	//return VideoMetrics{}, VideoAnalyticsError{}
	return videoMetrics, VideoAnalyticsError{}
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
		if channelNameLink.Pointer == nil {
			return "", "", errors.New("found null pointer while searching for span itemprop=author and link itemprop = name")
		}
		channelName := channelNameLink.Attrs()["content"]
		userName = channelName

		titleLink := htmlContent.Find("title")
		if titleLink.Pointer == nil {
			return "", "", errors.New("found null pointer while searching for youtube video title")
		}
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

	case "Instagram":
		var userName string
		meta := htmlContent.Find("meta", "name", "description")

		if meta.Pointer == nil {
			fmt.Println("Null pointer error")
		}

		metaDescription := meta.Attrs()["content"]
		regularExpression, err := regexp.Compile(`comments\s?-\s?\w*`)
		if err != nil {
			fmt.Println(err)
		}
		extractedMetaDescription := regularExpression.FindString(metaDescription)

		if len(extractedMetaDescription) > 0 {
			splitecExtractedMetaDescription := strings.Split(extractedMetaDescription, "-")
			userName = strings.Trim(splitecExtractedMetaDescription[1], " ")
		}

		return "", userName, nil

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
		videoMetrics, err := ScrapeInstagram(url, "'")
		return VideoInfoResponse{VideoInfo{}, videoMetrics, err}

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
	case "Instagram":
		cmd := exec.Command("python", "python_solution.py")
		out, _ := cmd.Output()
		jsonOutput := string(out)

		var jsonData map[string]interface{}
		error := json.Unmarshal([]byte(jsonOutput), &jsonData)

		if error != nil {
			fmt.Println(error)
		}

		playCount := jsonData["items"].([]interface{})[0].(map[string]interface{})["username"].(string)

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
	////url := "https://www.youtube.com/watch?v=2RERHdL0ZWY&ab_channel=ChamokHasan"
	////	url := "https://www.tiktok.com/@reddit.guy/video/7245740863991663914"
	//
	//urlArray := []string{"https://www.youtube.com/watch?v=JTjmkai8iMI&ab_channel=Askthereddit",
	//	"https://www.youtube.com/watch?v=4aV2OThWXXI&ab_channel=Askthereddit",
	//	"https://www.youtube.com/watch?v=XqCAsTNX7Rw&ab_channel=Askthereddit"}
	//
	//var channelArr []<-chan VideoInfoResponse
	//var videoAnalyticsProvider VideoAnalyticsProviderImpl
	//
	//wg := sync.WaitGroup{}
	//mutex := sync.Mutex{}
	//wg.Add(3)
	//
	//for x := 0; x < 3; x++ {
	//	videoInfo, _ := videoAnalyticsProvider.GetVideoInfo(urlArray[x])
	//
	//	tempChannel := videoAnalyticsProvider.GetFullVideoMetrics(videoInfo, &wg, &mutex)
	//
	//	channelArr = append(channelArr, tempChannel)
	//}
	//
	//for _, channelItem := range channelArr {
	//	videoItem := <-channelItem
	//	fmt.Println(videoItem)
	//}
	//wg.Wait()
	//
	//title, userName, err := getTitleAndUsername(url, "Tiktok")
	//if err != nil {
	//	log.Println(err)
	//}
	//
	//fmt.Printf("Username : %s ; title : %s\n", userName, title)

	ScrapeInstagram("https://www.instagram.com/reel/CoX3tHCOEUS/", `mid=ZIlxKgALAAFScymuGZ63XopJbgmQ; ig_did=3D7F7CB4-D713-4AC0-B289-1CB257D7146C; ig_nrcb=1; datr=KHGJZGAHYNaqVNT8a_ULEH_J; fbm_124024574287414=base_domain=.instagram.com; ds_user_id=6428581218; csrftoken=gIR0LRMzGwtZBswOMsjUFlvTtIYvHD8b; shbid="3375\0546428581218\0541722876021:01f73ddcd6511a5f2ae24cb30cd970559f9ee04ddd2c67e1d86cba68c3d1af338a9b9200"; shbts="1691340021\0546428581218\0541722876021:01f73e78321dfa01034452fefe26ceed79a8d4c6fa1cac2520f0ce956f70dcf86433a578"; sessionid=6428581218%3AA5GFgB30EGGQFW%3A22%3AAYfA3idt6kMo1YUWRLruKZsif8s1QKN0Qk54Ib3quTE; fbsr_124024574287414=lU0NJQRAp26tlgSlWEmuwDoqtz93XJWm7LscH6tRaB8.eyJ1c2VyX2lkIjoiMTAwMDA1ODY3NjEwNzE1IiwiY29kZSI6IkFRRGR4SXpZLV9rbjFzdnV1emUzb0VwS2RZOVFBSWRTaTVaMGc3QTlzM0hIakZYQ3lxbFFwZmFQcWdZM0hWRDdVaEREYl9lR1JuMlFXVmhvZFJMZG9JSnloWHpqekhCd3BfQW9oUzV3X3d2bV9BVnVleW5rd0s3cUZzeWhFc3pHWk1FSk9CSUJvMDc3VkN3SHdIMDYzRmV1ZWkzOUZwYVJGaE9FMVlJelVpY1NQcFpBXzF6MS01YXg0OXlUUUN2UTIwWVRJUEp6eFROTGpIeXo4RDVHeWhZR3BiMzVscURZc2M1dXh0MTFFX1VRYUo2RnRqV0hRVVZkVWoxcmRWak1GeEJLdnNydVYtT1Rkb1ZrbnlCODZvU1V5NVV4LWRaLV9hcm0ySlctd3NFOVJQLWlkbWVDS1lMY2JsTE10dmowUGtjcWx1T2dxTGNnWUpmMUVKUjM3NTMtIiwib2F1dGhfdG9rZW4iOiJFQUFCd3pMaXhuallCTzF2enNqWGJzRzZ5TDVEbXVMTkNWVHFzM0ViblpDQkJCOXdublZCT1QzWkFFVndhRzBaQUdhanVVMVBMWkM3N2NxTGFCY25kSWZtTXA2T29ITFJQek1aQ1lTZlcydExNNDdaQUlLaW50T3JqZVFrTmdrcWhGcjV5dDRYbFIzdTEzRlVaQmZUaFZ6ZVZLUU1EQnVLSXhUMHNtZ1JzekhlVVd2NkMxd09ReVdNc2l1MHdmUm5aQTBNZ3F1Z1pEIiwiYWxnb3JpdGhtIjoiSE1BQy1TSEEyNTYiLCJpc3N1ZWRfYXQiOjE2OTEzNDAzNjd9; rur="PRN\0546428581218\0541722876370:01f703c1d3cc873a19efcd7747e238391f399dcb59297f69f30accb544200f02748f852f"`)

}

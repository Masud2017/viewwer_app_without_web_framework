//package main
//
//import (
//	"encoding/json"
//	"errors"
//	"fmt"
//	"github.com/anaskhan96/soup"
//	"io"
//	"log"
//	"math/rand"
//	"net/http"
//	"regexp"
//	"strconv"
//	"strings"
//	"sync"
//	"time"
//)
//
//var tiktokRequestCounter int = 0 // counter for request in tiktok
//
//type VideoInfo struct {
//	ViewCount int
//	Platform  string
//	Username  string
//}
//
//func ProcessUrl(url string) VideoInfo {
//	var videoInfo VideoInfo
//
//	// portion of the code that checks whether the url from youtube.
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
//	// portion of the code that cehcks whether the url from tiktok
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
//	return nil
//}
//
//// func ScrapeInstagramData(videoInfo *VideoInfo,urll string) error {
//// 	shortCode := strings.Split(urll,"/")[4]
//// 	url := fmt.Sprintf("https://instagram-scraper-2022.p.rapidapi.com/ig/post_info/?shortcode=%s",shortCode)
//
//// 	req, _ := http.NewRequest("GET", url, nil)
//
//// 	req.Header.Add("X-RapidAPI-Key", "d34345206emshadd9b00e3b03f6fp1f97a4jsn83cf7dddaef2")
//// 	req.Header.Add("X-RapidAPI-Host", "instagram-scraper-2022.p.rapidapi.com")
//
//// 	res, _ := http.DefaultClient.Do(req)
//
//// 	if (res.Status != "200 OK") {
//// 		return errors.New("Fetching info with the instagram-scrapper-2022 from rapid api is failed")
//// 	}
//
//// 	defer res.Body.Close()
//// 	body, _ := io.ReadAll(res.Body)
//
//// 	// fmt.Println(res)
//// 	// fmt.Println(string(body))
//// 	responseData := string(body)
//// 	var jsonData map[string] interface{}
//
//// 	err := json.Unmarshal([]byte(responseData), &jsonData)
//// 	if err != nil {
//// 		fmt.Println(err)
//// 	}
//
//// 	rawPlayCount,ok := jsonData["video_play_count"]
//
//// 	if !ok {
//// 		fmt.Println("Something went wrong while trying to get the video play count")
//// 	}
//
//// 	playCount,ok :=  rawPlayCount.(float64)
//// 	if (!ok) {
//// 		fmt.Println("Getting error while trying to get float64 value from raw palycount")
//// 	}
//// 	fmt.Println(playCount)
//
//// 	videoInfo.ViewCount = int(playCount)
//
//// 	rawFullName, ok := jsonData["owner"].(map[string]interface{})["full_name"] // jsonData.(map[string]interface{}) is called type assertion
//
//// 	if (!ok) {
//// 		fmt.Println("Something went wrong while trying to fetch the full name data from the unmarshalled json data")
//// 	}
//
//// 	fullName,ok := rawFullName.(string)
//
//// 	if (!ok) {
//// 		fmt.Println("Something went wrong while trying to get string value from raw full name")
//// 	}
//
//// 	videoInfo.Username = fullName
//
//// 	return nil
//// }
//
//func ScrapeInstagramDataAlternative(videoInfo *VideoInfo, urll string) error {
//	url := fmt.Sprintf("https://instagram110.p.rapidapi.com/v2/instagram/post/info?query=%s&related_posts=false", urll)
//
//	req, _ := http.NewRequest("GET", url, nil)
//
//	req.Header.Add("X-RapidAPI-Key", "d34345206emshadd9b00e3b03f6fp1f97a4jsn83cf7dddaef2")
//	req.Header.Add("X-RapidAPI-Host", "instagram110.p.rapidapi.com")
//
//	res, _ := http.DefaultClient.Do(req)
//
//	if res.Status != "200 OK" {
//		return errors.New("Fetching info with the instagram110.p from rapid api is failed")
//	}
//
//	defer res.Body.Close()
//	body, _ := io.ReadAll(res.Body)
//
//	// fmt.Println(res)
//	// fmt.Println(string(body))
//	responseData := string(body)
//	var jsonData map[string]interface{}
//
//	err := json.Unmarshal([]byte(responseData), &jsonData)
//	if err != nil {
//		fmt.Println(err)
//	}
//
//	rawPlayCount, ok := jsonData["video_plays_count"]
//
//	if !ok {
//		fmt.Println("Something went wrong while trying to get the video play count")
//	}
//
//	playCount, ok := rawPlayCount.(float64)
//
//	if !ok {
//		fmt.Println("Getting error while trying to get float64 value from raw palycount")
//	}
//
//	videoInfo.ViewCount = int(playCount)
//
//	videoInfo.Username = "Not found"
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
//		log.Println("Access denied, So waiting some time and sending the request again 🙂 Wating for ", 10, " second")
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
//		// fetching the channel name (in this case user name)
//		rawAuthorData := jsonData["ItemModule"].(map[string]interface{})[tiktokVideoId].(map[string]interface{})["author"]
//		authorDataString, _ := rawAuthorData.(string)
//
//		videoInfo.Username = authorDataString
//
//	}
//
//	return nil
//}
//
//func ScrapeInstagramData(videoInfo *VideoInfo, url string) error {
//
//	req, err := http.NewRequest("GET", url, nil)
//
//	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/77.0.3865.120 Safari/537.36")
//	req.Header.Set("X-Requested-With", "XMLHttpRequest")
//	req.Header.Set("cookie", "mid=ZIlxKgALAAFScymuGZ63XopJbgmQ; ig_did=3D7F7CB4-D713-4AC0-B289-1CB257D7146C; ig_nrcb=1; datr=KHGJZGAHYNaqVNT8a_ULEH_J; fbm_124024574287414=base_domain=.instagram.com; ds_user_id=6428581218; shbid=\"3375\\0546428581218\\0541718418125:01f7fc8df5a78e07f4b0fbabdb2d972264b57e55bd727c853611a7918397ccd88dbcca93\"; shbts=\"1686882125\\0546428581218\\0541718418125:01f772b273faa88c03bd166c9bcfa6cd2bf29a9e18dbc4197bb11c50447e39a27a78e32c\"; fbsr_124024574287414=qLO34O776IjuoscRUlmawVuP4dmuwlMmFoyZC0zgTTk.eyJ1c2VyX2lkIjoiMTAwMDA1ODY3NjEwNzE1IiwiY29kZSI6IkFRQmJmOU5lSVJkaXR5eDhkMkhfTGJzZVVpa044TjlvRXpWeHN3cS1sbFYzM0ppZjJ4dG1YeFp0VUd2NTNrR1J1OFJFT0xFMG8zUnFROEVYMVZxeUFhQ25Cc00yZzRnTEE5eTJjU3RlT0ZaWmhRTllDRnBlRDZQRWU1cEtoazEyMXp5OF9CVkVWYU9oNnZOUWVZTmMyQVJQYUtQZDFwR1JDZHY2UWZOT3hMVi0wNXpLX1NSeXNVbmFmMVlaRTlZSmdWSjF0bmdFR3FmM1B0a3V5b09uemM2Yi1rU2NSVVJGZlBRcDBPWE5xNTVoYVRvQU11SzhIenoySVVhYTZYNlUtSzllblNiVlpzLWNkSWRPZzdqVHBBLW9WdmlTT3JTYXZIWVBiS3JrcF9ENFJKUHRuYzA2eTFkbVEtOFFEZzIxa2pEWGRycGkxWFdtRURNN0JxTDBBQ2FXIiwib2F1dGhfdG9rZW4iOiJFQUFCd3pMaXhuallCQUtZYjc3VDlaQ1pBcnJHbk9XbTdvYzdXcW9aQXM5NnptY3lNV3l3WVFsdjlWZHY1dWRTWkM4WGZ4VGNRWDFWWFl5WkJ6UDFaQmpnVG9VQ1dmSTdpcXRucURLV2V5eGZuN3FxemxaQVFzWkFaQTFZaGNxWG5GZlEwWkM5RGtmUFEybXFxZnN4R0FyNUIzVnRSYTF3ankxWkNDTmt2bWNoQlc3UHJhU1pBV1JpYWhKMGhiOE5Tcm9JUEFYMFpEIiwiYWxnb3JpdGhtIjoiSE1BQy1TSEEyNTYiLCJpc3N1ZWRfYXQiOjE2ODcwMDI1MTN9; csrftoken=gIR0LRMzGwtZBswOMsjUFlvTtIYvHD8b; sessionid=6428581218%3AA5GFgB30EGGQFW%3A22%3AAYf8lrg0fdE7nhLVyEDtzZ-QBYEyrmtGTnVTpviGSg; fbsr_124024574287414=qLO34O776IjuoscRUlmawVuP4dmuwlMmFoyZC0zgTTk.eyJ1c2VyX2lkIjoiMTAwMDA1ODY3NjEwNzE1IiwiY29kZSI6IkFRQmJmOU5lSVJkaXR5eDhkMkhfTGJzZVVpa044TjlvRXpWeHN3cS1sbFYzM0ppZjJ4dG1YeFp0VUd2NTNrR1J1OFJFT0xFMG8zUnFROEVYMVZxeUFhQ25Cc00yZzRnTEE5eTJjU3RlT0ZaWmhRTllDRnBlRDZQRWU1cEtoazEyMXp5OF9CVkVWYU9oNnZOUWVZTmMyQVJQYUtQZDFwR1JDZHY2UWZOT3hMVi0wNXpLX1NSeXNVbmFmMVlaRTlZSmdWSjF0bmdFR3FmM1B0a3V5b09uemM2Yi1rU2NSVVJGZlBRcDBPWE5xNTVoYVRvQU11SzhIenoySVVhYTZYNlUtSzllblNiVlpzLWNkSWRPZzdqVHBBLW9WdmlTT3JTYXZIWVBiS3JrcF9ENFJKUHRuYzA2eTFkbVEtOFFEZzIxa2pEWGRycGkxWFdtRURNN0JxTDBBQ2FXIiwib2F1dGhfdG9rZW4iOiJFQUFCd3pMaXhuallCQUtZYjc3VDlaQ1pBcnJHbk9XbTdvYzdXcW9aQXM5NnptY3lNV3l3WVFsdjlWZHY1dWRTWkM4WGZ4VGNRWDFWWFl5WkJ6UDFaQmpnVG9VQ1dmSTdpcXRucURLV2V5eGZuN3FxemxaQVFzWkFaQTFZaGNxWG5GZlEwWkM5RGtmUFEybXFxZnN4R0FyNUIzVnRSYTF3ankxWkNDTmt2bWNoQlc3UHJhU1pBV1JpYWhKMGhiOE5Tcm9JUEFYMFpEIiwiYWxnb3JpdGhtIjoiSE1BQy1TSEEyNTYiLCJpc3N1ZWRfYXQiOjE2ODcwMDI1MTN9; rur=\"PRN\\0546428581218\\0541718538529:01f7f5f54aa7c7af32fb0a15cf4646560d5d2a7c0d931d3a848cfa49a297332d14fd7122\"")
//
//	res, err := http.DefaultClient.Do(req)
//
//	if err != nil {
//		fmt.Println(err)
//	}
//
//	body, _ := io.ReadAll(res.Body)
//
//	fmt.Println(string(body))
//
//	return nil
//}
//
//func GetViewData(url string) (VideoInfo, error) {
//	videoInfo := ProcessUrl(url) // this will only populate the platform field
//
//	switch videoInfo.Platform {
//	case "Youtube":
//		err := ScrapeYoutubeData(&videoInfo, url) // it will populate the videInfo with video data not going to return anything
//		return videoInfo, err
//
//	case "Instagram":
//		err := ScrapeInstagramData(&videoInfo, url)
//		var error2 error
//		if err != nil {
//			error2 = ScrapeInstagramDataAlternative(&videoInfo, url)
//		}
//		if error2 != nil {
//			return videoInfo, errors.New("Both the scrapper failed to fetch the data")
//		}
//
//	case "Tiktok":
//		err := ScrapeTiktokData(&videoInfo, url)
//		return videoInfo, err
//	}
//
//	return videoInfo, nil
//}
//
//type VideoInfoWithErr struct {
//	VideoInfo VideoInfo
//	Err       error
//}
//
//func GetViewDataWithChan(url string, videoInfoWithErrChannel chan VideoInfoWithErr, wg *sync.WaitGroup) {
//	defer wg.Done()
//	videoInfo, err := GetViewData(url)
//	videoInfoWithErr := VideoInfoWithErr{videoInfo, err}
//
//	videoInfoWithErrChannel <- videoInfoWithErr // sending the data to the channel
//}
//
//func GetViewDataBatch(urlArr []string) []VideoInfoWithErr {
//	var channelList []chan VideoInfoWithErr
//	var videoInfoWithErrList []VideoInfoWithErr
//	wg := sync.WaitGroup{}
//
//	wg.Add(len(urlArr))
//	for x := 0; x < len(urlArr); x++ {
//		if tiktokRequestCounter == 15 {
//			for _, chanelItem := range channelList {
//				tempVideoInfoWithErr := <-chanelItem
//				videoInfoWithErrList = append(videoInfoWithErrList, tempVideoInfoWithErr)
//
//				log.Println(tempVideoInfoWithErr.VideoInfo)
//			}
//			channelList = nil // cleaning the slice
//
//			waitTime := rand.Intn(20-5+1) + 1
//			fmt.Printf("Total 15 request sended at once so waiting %d second before sending next 15\n", waitTime)
//			time.Sleep(time.Second * time.Duration(waitTime))
//			tiktokRequestCounter = 0
//		}
//		tempChannel := make(chan VideoInfoWithErr)
//		go GetViewDataWithChan(urlArr[x], tempChannel, &wg)
//		channelList = append(channelList, tempChannel)
//		tiktokRequestCounter++
//	}
//
//	wg.Wait()
//	return videoInfoWithErrList
//}
//
//// test code
//func main() {
//	//url := "https://www.tiktok.com/@asktheredditor/video/7244312458221931819"
//	url := "https://www.youtube.com/watch?v=153ZT6tgZiw&t=1s&ab_channel=Zahed%27sTake"
//	var urlArr []string
//
//	for x := 0; x < 100; x++ {
//		urlArr = append(urlArr, url)
//	}
//
//	videoInfoWithErr := GetViewDataBatch(urlArr)
//
//	for _, x := range videoInfoWithErr {
//		fmt.Println(x)
//	}
//
//}

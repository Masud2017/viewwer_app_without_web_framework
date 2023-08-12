package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

func main() {
	//cookies := "mid=ZIlxKgALAAFScymuGZ63XopJbgmQ; ig_did=3D7F7CB4-D713-4AC0-B289-1CB257D7146C; ig_nrcb=1; datr=KHGJZGAHYNaqVNT8a_ULEH_J; fbm_124024574287414=base_domain=.instagram.com; ds_user_id=6428581218; csrftoken=gIR0LRMzGwtZBswOMsjUFlvTtIYvHD8b; shbid=\"3375\\0546428581218\\0541722876021:01f73ddcd6511a5f2ae24cb30cd970559f9ee04ddd2c67e1d86cba68c3d1af338a9b9200\"; shbts=\"1691340021\\0546428581218\\0541722876021:01f73e78321dfa01034452fefe26ceed79a8d4c6fa1cac2520f0ce956f70dcf86433a578\"; sessionid=6428581218%3AA5GFgB30EGGQFW%3A22%3AAYfA3idt6kMo1YUWRLruKZsif8s1QKN0Qk54Ib3quTE; fbsr_124024574287414=qsMra1d-oRJizTTslVZnYphz7DhuD4lEBJcruxwRQS0.eyJ1c2VyX2lkIjoiMTAwMDA1ODY3NjEwNzE1IiwiY29kZSI6IkFRQXdqMk1zQ09ySkk2VkxXT3FKcmdiTDFsT3V0M2cwRlhINDRjTVZxa2RwTWZTMW5zTFFBeEt3RmI0ZUlCSkhjVnJCTmkzMGlLYW5oR1V2VEp2OE1DRy1XQmFyTi1PMWU4RGhKVTQ4c1JFZXFzRGRsZlJsVWh2ZnRPWFRYZVRmX1hRLS1UV0NHVWtfQWJDRldCWnFUeDdldXlFR1N1aENGeWZSRGVfUTJwaWdGQ1FsWEtUM1dwRmFYeVdzT1F6MXpuNWlpenlyQ2dHRWZ3LVVRWE5LTzdsNVpQemNiSzVfWlZISXg5MUdHRXgwSzFUY3ZtRU1Cc3V2Z1hLUTEybG1vbE02YS05QUpPZWt4emJfeWF2T2cwV1pvRGk3bTVqdHBwRjRkUGdYVE9QVERibVR5ci1JTXFzQ0dieHR6aUlXTklHSXdGQU9kUTNqc3JuUTVLVElaNGt0Iiwib2F1dGhfdG9rZW4iOiJFQUFCd3pMaXhuallCTzJGVld5YWdHaTVWMmJrc1dneTRUN2xSRFlWc3l0STlXeGd2bWkyYkY1TDVMNDJ5RUsyMk5ib1pCdmZXaWd5eVhsd1FxTWFyMmNMcVVpWkNQbnQ4WkI0bEFXSGZ0bWIwZG9ObXg1a01XNkhUdGdGSGh6Zk0xOTRwQUZldVA1dnpaQm5GQ0JZWEJHTjh5RFA0M2IwZ2NXV0swMVBIWDhNMFBJdnpWaEVSZklBR0xyVjhGQTdySWl3WkQiLCJhbGdvcml0aG0iOiJITUFDLVNIQTI1NiIsImlzc3VlZF9hdCI6MTY5MTM0MTY3MX0; rur=\"PRN\\0546428581218\\0541722877674:01f77821349cc860fa7a0f34e0b12149124cfe2e1282027855b2b797e50a530e43bb55eb\""
	//
	//headers := map[string]string{
	//	"authority":                   "www.instagram.com",
	//	"method":                      "GET",
	//	"path":                        "/api/v1/media/3035187576999415510/info/",
	//	"scheme":                      "https",
	//	"accept":                      "*/*",
	//	"accept-encoding":             "gzip, deflate, br",
	//	"accept-language":             "en-US,en;q=0.9",
	//	"cookie":                      "mid=ZIlxKgALAAFScymuGZ63XopJbgmQ; ig_did=3D7F7CB4-D713-4AC0-B289-1CB257D7146C; ig_nrcb=1; datr=KHGJZGAHYNaqVNT8a_ULEH_J; fbm_124024574287414=base_domain=.instagram.com; ds_user_id=6428581218; csrftoken=gIR0LRMzGwtZBswOMsjUFlvTtIYvHD8b; shbid=\"3375\\0546428581218\\0541722876021:01f73ddcd6511a5f2ae24cb30cd970559f9ee04ddd2c67e1d86cba68c3d1af338a9b9200\"; shbts=\"1691340021\\0546428581218\\0541722876021:01f73e78321dfa01034452fefe26ceed79a8d4c6fa1cac2520f0ce956f70dcf86433a578\"; fbsr_124024574287414=HQ-ZqQUwohukR0xcRkn9uaaSH4dxoHq6R8kM09F-zSE.eyJ1c2VyX2lkIjoiMTAwMDA1ODY3NjEwNzE1IiwiY29kZSI6IkFRQ3ZELUlIeWtUZ0k2MUtYc2RxcWotejdmS05wRGY4TElvbVlVWVBvdnh6a3h4UGRwOFVlYWdsZ1FCVXZUVFFabnhYUmoyZzZaMWlfSHhBVUNkRkE0UEVzY1o3NE5Oa2xzNVY0bFdtUHI5X2FZMlkyUDk2YTU5S3NZYlViSnhwMS1BbUJSUzBESVVELS1vZ3hRM3UxclI4d2xfeWdISHRqOTR2d2dnT1h0UktHV0h0QWR1dDJEaUZ5SDhvNlc4WTN4WWZjMXM5N0RDcXZOMFN0QmRzN011eWhLOEluVm9vUFFtRkJ3YkdvaTNuYVRqVl9NeDBzMmtSeGZSaGdMOEsyS1VXYlRTSjNTT3BiWVJZeE9QNGx6U2xIQlZOOGRBcDJMblRJTHdqcVJsNWItVEREeDlUejM1SFc0bkJlREVPUVQ4RnFRNER4Q1Bhb0FWYUp2RC1xS3VoIiwib2F1dGhfdG9rZW4iOiJFQUFCd3pMaXhuallCT3hjSjl6ak9LY2lBWkJPbXRRNzhXbWxXVDA0YzFiYXZLdFgyeE9QQVpDeDF1V2FlQ3VIVDg1bVR5ZTBiUWlZWkFKbVdQN1pCN01QZUF4T2dGUGtNa2k0NHJ6dFpBSGFMdzBHYVFOZE9qQVdYMWxMaUhZZzdua0F4SFF4aVlxQ2JwN2I3YnFnaHN2WkJER0d1bjh6TEVCaDlOSGJVbVd6MGdNMFpDR2ptcmFiRE1zUDZVWkNka2JvMFpCUklaRCIsImFsZ29yaXRobSI6IkhNQUMtU0hBMjU2IiwiaXNzdWVkX2F0IjoxNjkxNTk4MTE1fQ; sessionid=6428581218%3AA5GFgB30EGGQFW%3A22%3AAYfazLyUI_7iwb8A-IOttymVT7ifzhvqtSPnALDUFG4; rur=\"PRN\\0546428581218\\0541723134127:01f793dab9580299f8ccdfd2d6994025ee4104598615cb8a6b59344c11d7da32b97d149d\"",
	//	"referer":                     "https://www.instagram.com/reel/CoX3tHCOEUS/",
	//	"sec-ch-prefers-color-scheme": "light",
	//	"sec-ch-ua":                   "\"Not?A_Brand\";v=\"8\", \"Chromium\";v=\"108\"",
	//	"sec-ch-ua-mobile":            "?0",
	//	"sec-ch-ua-platform":          "\"Windows\"",
	//	"sec-fetch-dest":              "empty",
	//	"sec-fetch-mode":              "cors",
	//	"sec-fetch-site":              "same-origin",
	//	"user-agent":                  "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36",
	//	"viewport-width":              "1280",
	//	"x-asbd-id":                   "198387",
	//	"x-ig-app-id":                 "936619743392459",
	//	"x-ig-www-claim":              "0",
	//	"x-requested-with":            "XMLHttpRequest",
	//}
	//
	//url := "https://www.instagram.com/api/v1/media/3035187576999415510/info/"
	//
	//req, err := http.NewRequest("GET", url, nil)
	//if err != nil {
	//	fmt.Println("Error creating request:", err)
	//	return
	//}
	//
	//for key, value := range headers {
	//	fmt.Println(key, " ", value)
	//	req.Header.Add(key, value)
	//}
	//
	//client := &http.Client{}
	//resp, err := client.Do(req)
	//if err != nil {
	//	fmt.Println("Error sending request:", err)
	//	return
	//}
	//defer resp.Body.Close()
	//
	//fmt.Println("Status Code:", resp.StatusCode)
	//
	//body, err := ioutil.ReadAll(resp.Body)
	////
	////fmt.Println(string(body))
	////fmt.Println(soup.HTMLParse(string(body)))
	//
	////fmt.Println(req.Header)
	//
	//var jsonData map[string]interface{}
	//error := json.Unmarshal([]byte(body), &jsonData)
	//fmt.Println("error : ", error)
	//fmt.Println(jsonData)
	//fmt.Println(string(body))

	cmd := exec.Command("python", "python_solution.py")
	out, _ := cmd.Output()
	jsonOutput := string(out)

	var jsonData map[string]interface{}
	error := json.Unmarshal([]byte(jsonOutput), &jsonData)

	if error != nil {
		fmt.Println(error)
	}

	playCount := jsonData["items"].([]interface{})[0].(map[string]interface{})["video_duration"]
	fmt.Println(playCount)
}

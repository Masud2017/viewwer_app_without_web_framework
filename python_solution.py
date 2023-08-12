import requests
import json

headers = {
    "authority": "www.instagram.com",
    "method": "GET",
    "path": "/api/v1/media/3033137866545317138/info/",
    "scheme": "https",
    "accept": "*/*",
    "accept-encoding": "gzip, deflate, br",
    "accept-language": "en-US,en;q=0.9",
    "cookie": "mid=ZIlxKgALAAFScymuGZ63XopJbgmQ; ig_did=3D7F7CB4-D713-4AC0-B289-1CB257D7146C; ig_nrcb=1; datr=KHGJZGAHYNaqVNT8a_ULEH_J; fbm_124024574287414=base_domain=.instagram.com; ds_user_id=6428581218; csrftoken=gIR0LRMzGwtZBswOMsjUFlvTtIYvHD8b; sessionid=6428581218%3AA5GFgB30EGGQFW%3A22%3AAYfazLyUI_7iwb8A-IOttymVT7ifzhvqtSPnALDUFG4; shbid=\"3375\0546428581218\0541723136658:01f7d502488b704637c9b56ae8344379dbc500d5615381d25241dff12f4c85d87537228a\"; shbts=\"1691600658\0546428581218\0541723136658:01f7a286932db917f076d30fffdd6a4bda3052c402c1e37ca1d02995cea87ff9fe8f165a\"; fbsr_124024574287414=M1i4OW6E1KsfunLaSdCxZx8XALsrVbnB6YmxwwDX0uY.eyJ1c2VyX2lkIjoiMTAwMDA1ODY3NjEwNzE1IiwiY29kZSI6IkFRQ0h5SXhkQ3RXQURGU2psZzNMWTdfbXdaVXFrVXAza0NobmxlOUpUM3JKSG96SkFxMDJFYkxINnhhNlRyYjc3S2k0SHVTMUdJSDd3YVdCd2NMTVJsQ3BCU3pBOG5CcnBwbTdDMGNTbFhQRkVJeU5hQ0FPTE5BUVlnaWJLMVJQcWJQUll1MlRrWmpPOFhHaU5oV3F6S3JaeF8xTEQteU1jQmttMGtFa056QVJmemtJTlRHcXBSMWJNOXRTOFd5Q0lKT2JjNWdvWnFCbm5DeFRRcmx6d2JjTC1Dbng3bXJmQUhyQUF0bm1sRmtDTVc1c004MVlEU25ja2dPMmtKRzRKME1HbVhSNkhuUHdLZzRYM3UyUFhHN1NNTU5oRUpKcWVzbFNlNElXdkk1OGZ2eGJ0LUFxdHVWcWQ3M3hZVHY2bk0zdFZ3UEMwRG1QUGFkUHMxbmgtakRBIiwib2F1dGhfdG9rZW4iOiJFQUFCd3pMaXhuallCT1pCd3pxcHhGc0RzMFpDMEttdlNVNkFqa2NaQUhlV2ZDbGpWOGVqWUNHRGx6Tkx4VjQ5UVVJNVhDWkF0VjZ6em4zc1A0cUxFaHVHZk1nNEpONWdoTmI1aGc4OFB2TW5aQWRHaTRQbnFFZlpBWW11cEh5cm13YUg4dUtCODd5YVR5ekFUc3NzckJuTGppa1pCb042WVpCYmgyT1huenI0Rk1jNnFraXhRQ2RaQjBpczBpaGtOck5xbTVaQlpBY1pEIiwiYWxnb3JpdGhtIjoiSE1BQy1TSEEyNTYiLCJpc3N1ZWRfYXQiOjE2OTE2MDA2NTl9; rur=\"PRN\0546428581218\0541723136660:01f71e13eaa9cb48703a684876543c42d2ebe255257f3d63688387a1f98601635793c9b2\"",
    "referer": "https://www.instagram.com/reel/CoX3tHCOEUS/",
    "sec-ch-prefers-color-scheme": "light",
    "sec-ch-ua": "\"Not?A_Brand\";v=\"8\", \"Chromium\";v=\"108\"",
    "sec-ch-ua-mobile": "?0",
    "sec-ch-ua-platform": "\"Windows\"",
    "sec-fetch-dest": "empty",
    "sec-fetch-mode": "cors",
    "sec-fetch-site": "same-origin",
    "user-agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36",
    "viewport-width": "1280",
    "x-asbd-id": "198387",

    "x-ig-app-id": "936619743392459",
    "x-ig-www-claim": "0",
    "x-requested-with": "XMLHttpRequest",
}


req = requests.get("https://www.instagram.com/api/v1/media/3033137866545317138/info/",headers = headers)

print(req.text)

import requests
import sys
import re

cookies = sys.argv[1]
url = sys.argv[2]

def shortcode_to_mediaid(shortcode):
    alphabet = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_'
    mediaid = 0

    for letter in shortcode:
        mediaid = (mediaid * 64) + alphabet.index(letter)

    return mediaid

def extract_media_id_from_url (url):
    res = re.search(r'([a-zA-Z0-9_-]+)\/?$',url,re.IGNORECASE).group()
    if res[len(res)-1] == "/":
        res = res[:len(res)-1]
        pass
    return res

media_url = "https://www.instagram.com/api/v1/media/"+str(shortcode_to_mediaid(extract_media_id_from_url(url)))+"/info/"
path = "/api/v1/media/"+str(shortcode_to_mediaid(extract_media_id_from_url(url)))+"/info/"

headers = {
    "authority": "www.instagram.com",
    "method": "GET",
    "path": path,
    "scheme": "https",
    "accept": "*/*",
    "accept-encoding": "gzip, deflate, br",
    "accept-language": "en-US,en;q=0.9",
    "cookie": cookies,
    "referer": url,
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


req = requests.get(media_url,headers = headers)

print(req.text)

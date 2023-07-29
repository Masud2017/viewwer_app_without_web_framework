
function process_url(url) {
    if (url.includes("tiktok")) {
        return {"platform":"Tiktok"};
    } else if(url.includes("youtube")) {
        return {"platform":"Youtube"};
    } else if (url.includes("instagram")) {
        return {"platform":"Instagram"};
    }
    return "Not detected";
}

function getEmbeddedHTMLForYoutube(url) {
    let videoIdRegex = /(v=.*)/g;

    let occuranceIdx = url.search(videoIdRegex);

    let extracted_string = url.slice(occuranceIdx,url.length - 1);
    let videoId = extracted_string.split("&")[0].split("=")[1];

    let embedded_html = String.raw`<iframe title="YouTube video player" class="youtube-player" type="text/html" width="640"
                height="390" src="http://www.youtube.com/embed/${videoId}" frameborder="0" allowFullScreen></iframe>`;

    return embedded_html;
}

function getEmbeddedHTMLForInstagram(url) {
    let embedded_html = String.raw`<iframe src="${url}embed/" width="612" height="710" frameborder="0" scrolling="no" allowtransparency="true"></iframe>`;
    return embedded_html
}

function getTiktokVideoId (url) {
    url = url + "?"
    const pattern = /\/video\/(\d+)\?/;
    let res = url.search(pattern)
    return url.slice(res,url.length - 1).split("/")[2].split("?")[0]
}
function getEmbeddedHTMLForTiktok(url) {
    let videoId = getTiktokVideoId(url)
    let embedded_text_template = String.raw`<blockquote className="tiktok-embed" cite="${url}" data-video-id="${videoId}" style="max-width: 605px; min-width: 325px;"> <iframe name="__tt_embed__v6828268207359413509" sandbox="allow-popups allow-popups-to-escape-sandbox allow-scripts allow-top-navigation allow-same-origin"src="https://www.tiktok.com/embed/v2/${videoId}?lang=en-US&amp;"style="width: 100%; height: 707px; display: block; visibility: unset; max-height: 707px;"></iframe></blockquote>`
        // <script async="" src="https://www.tiktok.com/embed.js"></script>

    return embedded_text_template.replace("\n","").replace(" ","");
}

function getEmbeddedHTML(url) {
    let processed_data = process_url(url);

    switch (processed_data.platform) {
        case "Youtube":
            return getEmbeddedHTMLForYoutube(url);
            break;
        case "Tiktok":
            let embeddedHTMLForTiktok = getEmbeddedHTMLForTiktok(url);
           return embeddedHTMLForTiktok

            break;
        case "Instagram":
            return getEmbeddedHTMLForInstagram(url);
            break;
    }
}
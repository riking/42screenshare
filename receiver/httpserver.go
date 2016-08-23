package receiver

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

func SetupHTTP() {
	http.HandleFunc("/image.png", httpShowImage)
	http.HandleFunc("/", httpShowHTML)
}

var imageBytes []byte
var imageLastUpdated time.Time
var imageMutex sync.Mutex

func httpShowImage(w http.ResponseWriter, r *http.Request) {
	imageMutex.Lock()
	defer imageMutex.Unlock()

	// Can't use L-M as it has a resolution of 1 second
	etag := fmt.Sprintf("t-%d-%d", imageLastUpdated.Unix(), imageLastUpdated.UnixNano())

	w.Header().Set("ETag", etag)
	if r.Header.Get("If-None-Match") == etag {
		w.WriteHeader(http.StatusNotModified)
		return
	}

	w.Write(imageBytes)
}

var html = []byte(`<!DOCTYPE html>
<html>
<head>
<script src="https://ajax.googleapis.com/ajax/libs/jquery/3.1.0/jquery.min.js"></script>
</head>
<body>
<img id=sshot>
<script>
var shotImg = document.getElementById('sshot');
var imgTemplate = "/image.png";
var REFRESH_DELAY = 100;
function reloadImage() {
	shotImg.src = imgTemplate + "#" + (new Date().getTime());
}
shotImg.addEventListener("load", reloadImage);
setTimeout(reloadImage, 0);
</script>
</body>
</html>`)

func httpShowHTML(w http.ResponseWriter, r *http.Request) {
	if r.URL.EscapedPath() != "/" {
		http.NotFound(w, r)
		return
	}
	w.Write(html)
}

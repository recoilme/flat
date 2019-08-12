package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"time"

	"golang.org/x/net/html/charset"
)

func main() {
	s, err := GetUtf8("https://news.ycombinator.com/", 0, "")
	//log.Println(err, s)
	_, _ = s, err
	initBinaries()
}
func initBinaries() {
	log.Println(os.Getenv("PATH"))
	var err error
	goBin, err := exec.LookPath("go")
	if err != nil {
		fmt.Println("hover: Failed to lookup `go` executable. Please install Go.\nhttps://golang.org/doc/install")
		os.Exit(1)
	}
	flutterBin, err := exec.LookPath("flutter")
	if err != nil {
		fmt.Println(err, "hover: Failed to lookup `flutter` executable. Please install flutter.\nhttps://flutter.dev/docs/get-started/install")
		os.Exit(1)
	}
	_, _ = goBin, flutterBin

}

// GetUtf8 return utf8 string from url
func GetUtf8(geturl string, t time.Duration, ua string) (s string, err error) {
	if t == 0 {
		t = time.Second * 10
	}
	ctx, cncl := context.WithTimeout(context.Background(), t)
	defer cncl()
	q := geturl
	if !strings.HasPrefix(geturl, "http") {
		q = "http://" + geturl
	}
	req, err := http.NewRequest(http.MethodGet, q, nil)
	if err != nil {
		return s, err
	}
	//Host
	u, err := url.Parse(q)
	if err == nil && len(u.Host) > 2 {
		req.Header.Set("Host", u.Host)
	}
	var defHeaders = make(map[string]string)
	defHeaders["User-Agent"] = "Mozilla/5.0 (compatible; YandexBot/3.0; +http://yandex.com/bots)"
	defHeaders["Accept"] = "text/html,application/xhtml+xml,application/xml,application/rss+xml;q=0.9,image/webp,*/*;q=0.8"
	defHeaders["Accept-Language"] = "ru-RU,ru;q=0.8,en-US;q=0.5,en;q=0.3"
	for k, v := range defHeaders {
		req.Header.Set(k, v)
	}
	if ua != "" {
		req.Header.Set("User-Agent", ua)
	}
	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		return s, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		contentType := resp.Header.Get("Content-Type")
		utf8, err := charset.NewReader(resp.Body, contentType)
		if err != nil {
			return s, err
		}
		body, err := ioutil.ReadAll(utf8)
		if err != nil {
			return s, err
		}
		return string(body), err
	}
	return s, fmt.Errorf("Error, statusCode:%d", resp.StatusCode)
}

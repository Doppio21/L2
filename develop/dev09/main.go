package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/net/html"
)

func Save(dir string, link string, r io.Reader) error {
	u, err := url.Parse(link)
	if err != nil {
		return err
	}

	var (
		filePath  string
		queryPath = u.Path
	)

	switch queryPath {
	case "", "/":
		filePath = "/index.html"
	default:
		splitted := strings.Split(queryPath, "/")
		filePath = splitted[len(splitted)-1]
	}

	if err := os.MkdirAll(filepath.Dir(dir), 0755); err != nil {
		return err
	}

	dir = strings.TrimSuffix(dir, filePath)
	f, err := os.OpenFile(filepath.Join(dir, filePath),
		os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0755)
	if err != nil {
		return err
	}

	tr := io.TeeReader(r, f)
	_, err = io.ReadAll(tr)
	return err
}

func Download(curDir, curLink string, r bool, links map[string]int) error {
	resp, err := http.Get(curLink)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if !r {
		return Save(curDir, curLink, resp.Body)
	}

	curDir = filepath.Join(curDir, resp.Request.Host)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	buff := bytes.NewReader(body)
	getLinks(buff, links)
	if _, err := buff.Seek(0, io.SeekStart); err != nil {
		return err
	}

	for link, cnt := range links {
		if cnt > 0 {
			continue
		}

		links[link]++
		if strings.HasSuffix(curLink, link) {
			continue
		}

		nextDir := curDir
		u, err := url.Parse(link)
		if err != nil {
			return err
		}

		if strings.Contains(link, "http") &&
			!strings.Contains(link, u.Host) {
			continue
		}

		startIdx := 0
		splitted := strings.Split(u.Path, "/")
		for idx, dir := range splitted {
			if strings.HasSuffix(nextDir, dir) {
				startIdx = idx
				break
			}
		}

		for i := startIdx + 1; i < len(splitted); i++ {
			nextDir = filepath.Join(nextDir, splitted[i])
		}

		nextLink := link
		if !strings.HasPrefix(link, "http") {
			temp := strings.TrimRight(curLink, "/")
			nextLink = temp + "/" + strings.TrimLeft(link, "/")
		}

		if err := Download(nextDir, nextLink, false, links); err != nil {
			return err
		}
	}

	return Save(curDir, curLink, buff)
}

func getLinks(body io.Reader, links map[string]int) {
	z := html.NewTokenizer(body)

	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			return
		case html.StartTagToken, html.EndTagToken:
			token := z.Token()
			for _, attr := range token.Attr {
				if attr.Key == "href" || attr.Key == "src" {
					links[attr.Val] = 0
				}
			}
		}
	}
}

func main() {
	var (
		path string
		r    bool
	)

	link := os.Args[len(os.Args)-1]
	pathf, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		return
	}

	flag.StringVar(&path, "P", pathf, "выбрать каталог, в который будут загружаться файлы")
	flag.BoolVar(&r, "r", false, "рекурсивная работа утилиты")
	flag.Parse()

	if err := Download(path, link, r, map[string]int{path: 1}); err != nil {
		fmt.Println(err)
	}
}

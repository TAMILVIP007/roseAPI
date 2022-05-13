package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var (
	c = &http.Client{}
)

func Tpb(w http.ResponseWriter, r *http.Request) {
	if !blockWrongMethod(w, r, "GET") {
		return
	}
	r.Header.Set("X-Start-Time", fmt.Sprint(time.Now().UnixNano()))
	query := r.URL.Query()
	if query.Get("help") != "" {
		w.Write([]byte(strings.ReplaceAll(_help_["tpb"], "{}", r.URL.Hostname())))
		return
	}
	i := query.Get("i")
	q := query.Get("q")
	if q == "" {
		http.Error(w, "missing query", http.StatusBadRequest)
		return
	}
	url := "https://tpb23.ukpass.co/apibay/q.php" + "?q=" + url.QueryEscape(q)
	resp, err := c.Get(url)
	if !ERR(err, w) {
		return
	}
	WriteJson(w, r, resp.Body, i)
}

func Google(w http.ResponseWriter, r *http.Request) {
	if !blockWrongMethod(w, r, "GET") {
		return
	}
	r.Header.Set("X-Start-Time", fmt.Sprint(time.Now().UnixNano()))
	query := r.URL.Query()
	if query.Get("help") != "" {
		w.Write([]byte(strings.ReplaceAll(_help_["google"], "{}", r.URL.Hostname())))
		return
	}
	q := query.Get("q")
	i := query.Get("i")
	if query.Get("q") == "" {
		http.Error(w, "missing query", http.StatusBadRequest)
		return
	}
	URL := "https://www.google.com/search?q=" + url.QueryEscape(q)
	resp, err := c.Get(URL)
	if !ERR(err, w) {
		return
	}
	doc, _ := goquery.NewDocumentFromReader(resp.Body)
	var results []GoogleResult
	doc.Find("div").Each(func(i int, s *goquery.Selection) {
		if s.HasClass("ZINbbc luh4tb xpd O9g5cc uUPGi") {
			var result GoogleResult
			s.Find("div").Each(func(i int, s *goquery.Selection) {
				if s.HasClass("egMi0 kCrYT") {
					result.Title = s.Text()
					url := strings.Split(strings.Replace(s.Find("a").AttrOr("href", ""), "/url?q=", "", 1), "&")[0]
					result.Url = url
				} else if s.HasClass("BNeawe s3v9rd AP7Wnd") {
					result.Description = s.Text()
				}
			})
			results = append(results, result)
		}
	})
	if !ERR(err, w) {
		return
	}
	var data string
	if len(results) != 0 {
		data = EncodeJson(results)
	} else {
		data = `[]`
	}
	WriteJson(w, r, data, i)
}

func Youtube(w http.ResponseWriter, r *http.Request) {
	if !blockWrongMethod(w, r, "GET") {
		return
	}
	r.Header.Set("X-Start-Time", fmt.Sprint(time.Now().UnixNano()))
	query := r.URL.Query()
	if query.Get("help") != "" {
		w.Write([]byte(strings.ReplaceAll(_help_["youtube"], "{}", r.URL.Hostname())))
		return
	}
	q := query.Get("q")
	i := query.Get("i")
	if q == "" {
		http.Error(w, "missing query", http.StatusBadRequest)
		return
	}
	URL := "https://www.youtube.com/results?search_query=" + url.QueryEscape(q)
	resp, err := c.Get(URL)
	if !ERR(err, w) {
		return
	}
	var exp, _ = regexp.Compile(`ytInitialData = [\s\S]*]`)
	b, _ := ioutil.ReadAll(resp.Body)
	match := exp.FindStringSubmatch(string(b))
	var d string
	if len(match) != 0 {
		d = match[0]
		d = strings.Replace(d, "ytInitialData = ", "", 1)
		d = strings.Split(d, ";</script>")[0]
	}
	pData := ParseYoutubeRAW(d)
	WriteJson(w, r, pData, i)
}

func ImDB(w http.ResponseWriter, r *http.Request) {
	if !blockWrongMethod(w, r, "GET") {
		return
	}
	r.Header.Set("X-Start-Time", fmt.Sprint(time.Now().UnixNano()))
	query := r.URL.Query()
	if query.Get("help") != "" {
		w.Write([]byte(strings.ReplaceAll(_help_["imdb"], "{}", r.URL.Hostname())))
		return
	}
	q := query.Get("q")
	i := query.Get("i")
	id := query.Get("id")
	if q == "" && id == "" {
		http.Error(w, "missing query", http.StatusBadRequest)
		return
	}
	if q != "" {
		firstLetter := strings.ToLower(string(q[0]))
		URL := "https://v2.sg.media-imdb.com/suggestion/titles/" + firstLetter + "/" + url.QueryEscape(q) + ".json"
		resp, err := c.Get(URL)
		if !ERR(err, w) {
			return
		}
		var data ImDBData
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}
		var results []ImDBResult
		for _, r := range data.D {
			results = append(results, ImDBResult{Title: r.L, Year: fmt.Sprint(r.Y), ID: r.ID, Actors: r.S, Rank: fmt.Sprint(r.Rank), Link: "https://www.imdb.com/title/" + r.ID, Poster: r.I.ImageURL})
		}
		if !ERR(err, w) {
			return
		}
		var result string
		if len(results) != 0 {
			result = EncodeJson(results)
		} else {
			result = `[]`
		}
		WriteJson(w, r, result, i)
	} else if id != "" {
		w.Write([]byte(`{"title":"` + id + `"}`))
	}
}

func ChatBot(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("X-Start-Time", fmt.Sprint(time.Now().UnixNano()))
	var API = "https://icap.iconiq.ai/talk?&botkey=icH-VVd4uNBhjUid30-xM9QhnvAaVS3wVKA3L8w2mmspQ-hoUB3ZK153sEG3MX-Z8bKchASVLAo~&channel=7&sessionid=482070240&client_name=uuiprod-un18e6d73c-user-19422&id=true"
	query := r.URL.Query()
	q := query.Get("message")
	if query.Get("help") != "" {
		w.Write([]byte(strings.ReplaceAll(_help_["chatbot"], "{}", r.URL.Hostname())))
		return
	}
	if q == "" {
		http.Error(w, "missing 'message'", http.StatusBadRequest)
		return
	}
	req, err := http.PostForm(API, url.Values{"input": {q}})
	if !ERR(err, w) {
		return
	}
	defer req.Body.Close()
	var resp map[string]interface{}
	json.NewDecoder(req.Body).Decode(&resp)
	msg := resp["responses"].([]interface{})[0].(string)
	d := `{"message": "` + msg + `"}`
	WriteJson(w, r, d, "")

}

func Lyrics(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("X-Start-Time", fmt.Sprint(time.Now().UnixNano()))
	query := r.URL.Query()
	if query.Get("help") != "" {
		w.Write([]byte(strings.ReplaceAll(_help_["lyrics"], "{}", r.URL.Hostname())))
		return
	}
	q := query.Get("q")
	if q == "" {
		http.Error(w, "missing query", http.StatusBadRequest)
		return
	}
	seURL := "https://cse.google.com/cse/element/v1?rsz=filtered_cse&num=10&hl=en&source=gcsc&gss=.com&cselibv=3e1664f444e6eb06&cx=15ba6306c8bf0c5d0&q=" + q + "&safe=off&cse_tok=AJvRUv3bw29E-03lEFZhaQV4UDN7:1652443252075&exp=csqr,cc&callback=google.search.cse.api10882"
	resp, err := c.Get(seURL)
	if !ERR(err, w) {
		return
	}
	var b []byte
	if b, err = ioutil.ReadAll(resp.Body); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	bstring := strings.Replace(string(b), `google.search.cse.api10882(`, "", 1)
	bstring = strings.Replace(bstring, ");", "", 1)
	bstring = strings.Replace(bstring, "/*O_o*/", "", -1)
	rd := strings.NewReader(bstring)
	var d map[string]interface{}
	if err := json.NewDecoder(rd).Decode(&d); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	var lyricURL string
	for _, c := range d["results"].([]interface{}) {
		if c.(map[string]interface{})["url"].(string) != "" {
			if strings.Contains(c.(map[string]interface{})["url"].(string), "lyrics.com") {
				lyricURL = c.(map[string]interface{})["url"].(string)
				break
			}
		}
	}
	if lyricURL == "" {
		http.Error(w, "lyrics not found", http.StatusNotFound)
		return
	}
	resp, err = c.Get(lyricURL)
	if !ERR(err, w) {
		return
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	var t string
	doc.Find(".lyric-body-text").Each(func(i int, s *goquery.Selection) {
		t = s.Text()
		fmt.Println(t)
	})
	w.Write([]byte(t))

}

func init() {
	http.HandleFunc("/tpb", Tpb)
	http.HandleFunc("/google", Google)
	http.HandleFunc("/youtube", Youtube)
	http.HandleFunc("/imdb", ImDB)
	http.HandleFunc("/chatbot", ChatBot)
	http.HandleFunc("/lyrics", Lyrics)
}

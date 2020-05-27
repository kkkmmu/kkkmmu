package main

import (
	"bytes"
	"crypto/tls"
	"io/ioutil"
	"log"
	"net/http"
	//"net/http/httputil"
	"os"
	"strings"
	"time"
	//"errors"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"regexp"
	"sync"
)

var Downloaded map[string]string = make(map[string]string, 10000)
var Filtered map[string]string = make(map[string]string, 10000)

type Movie struct {
	Name    string            `json:Name`
	Actress string            `json:Actress`
	Date    string            `json:Date`
	Tags    string            `json:Tags`
	Magnets map[string]string `json:Magnets`
	Code    string            `json:Code`
	Length  string            `json:Length`
	Link    string            `json:Link`
	Cover   string            `json:Cover`
	Samples map[string]string `json:Sampes`
}

func (m Movie) String() string {
	res := fmt.Sprintf("======================================================================================\n")
	res += fmt.Sprintf("%20s \n\t%5s \n\t%10s \n\t%5s \n\t%5s: \n", m.Name, m.Actress, m.Date, m.Code, m.Length)
	res += fmt.Sprintf("\t%20s \n\t%20s\n", m.Link, m.Cover)
	res += fmt.Sprintf("\t")
	res += fmt.Sprintf("%10s ", m.Tags)
	res += fmt.Sprintf("\n")
	for _, mag := range m.Magnets {
		res += fmt.Sprintf("\t%s\n", mag)
	}
	for _, sam := range m.Samples {
		res += fmt.Sprintf("\t%s\n", sam)
	}

	return res
}

func (m Movie) IsFiltered() bool {

	if strings.Contains(m.Tags, "粪") ||
		strings.Contains(m.Tags, "糞") ||
		strings.Contains(m.Tags, "大便") ||
		strings.Contains(m.Tags, "肛") ||
		strings.Contains(m.Tags, "中年") ||
		strings.Contains(m.Tags, "胖") {
		Filtered[m.Link] = m.Link
		return true
	}

	if strings.Contains(m.Name, "粪") ||
		strings.Contains(m.Name, "糞") ||
		strings.Contains(m.Name, "大便") ||
		strings.Contains(m.Name, "肛") ||
		strings.Contains(m.Name, "中年") ||
		strings.Contains(m.Name, "胖") {
		Filtered[m.Link] = m.Link
		return true
	}

	if strings.Contains(m.Date, "2018") ||
		strings.Contains(m.Date, "2017") ||
		strings.Contains(m.Date, "2016") ||
		strings.Contains(m.Date, "2015") ||
		strings.Contains(m.Date, "2014") ||
		strings.Contains(m.Date, "2013") {
		return false
	}

	Filtered[m.Link] = m.Link
	return true
}

var MagR = regexp.MustCompile(`href="(?P<mag>[[:alnum:]\:\?=_\-&\.]+)"`)

var ROOT = "art"

func main() {
	data, err := ioutil.ReadFile("list2.txt")
	if err != nil {
		log.Println("Error happened when read file: ", err.Error())
		return
	}
	d := NewDownloader(ROOT)
	sr := bytes.NewBufferString(string(data))
	d.Start()
	go func() {
		for {
			line, err := sr.ReadString('\n')
			if err != nil {
				log.Println("Error when readline: ", err.Error())
				return
			}
			//Read string will return the "\n", we must strip this when do http reqeuest"
			d.Download(line[:len(line)-1])
		}
	}()

	for {
		select {
		case <-d.done:
			log.Println("Download Finished")
		case link := <-d.success:
			log.Println("Download success: ", link)
		}
	}
}

type Downloader struct {
	root      string
	ratelimit <-chan time.Time
	client    *http.Client
	queue     chan string
	done      chan bool
	failed    chan string
	success   chan string
	userAgent string
	doc       *goquery.Document
}

func NewDownloader(root string) *Downloader {
	return &Downloader{
		root:      root,
		ratelimit: time.Tick(time.Millisecond * 1000),
		client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		},
		queue:   make(chan string, 500),
		done:    make(chan bool),
		failed:  make(chan string),
		success: make(chan string),
		//	userAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/48.0.2564.109 Safari/537.36",
		userAgent: "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Safari/537.36",
	}
}

func (d *Downloader) Download(link string) {
	d.queue <- d.verifyLink(link)
}

func (d *Downloader) verifyLink(link string) string {
	/*
		if !strings.HasPrefix(link, "http") {
			return "http://" + link
		}
	*/
	return link
}

func (d *Downloader) Start() {
	if _, err := os.Stat(d.root); err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(d.root, 0666)
		}
	}
	go d.download()
}

func (d *Downloader) download() {
	os.Remove("javlist.txt")
	os.Remove("javinfo.txt")
	var (
		wg sync.WaitGroup
	)
	for link := range d.queue {
		if _, ok := Filtered[link]; ok {
			continue
		}
		if _, ok := Downloaded[link]; ok {
			continue
		}
		<-d.ratelimit
		go func(link string) {
			wg.Add(1)
			var movie Movie
			//fields := strings.Split(link, "/")
			//name := fields[len(fields)-1]

			//fmt.Printf("Downloading: %s -----> %s\n", link, name)
			req, err := http.NewRequest("GET", link, nil)
			if err != nil {
				log.Println("Create new reqeust error ", err.Error())
				return
			}
			req.Header.Add("User-Agent", d.userAgent)
			resp, err := d.client.Do(req)
			if err != nil {
				log.Println("Download ", link, " error: ", err.Error())
				return
			}

			movie.Link = link

			doc, err := goquery.NewDocumentFromResponse(resp)
			if err != nil {
				log.Println(err.Error())
				return
			}

			d.doc = doc

			titles, err := d.GetText("div.container>h3")
			if err != nil {
				log.Println(err)
				return
			}

			//log.Println(titles[0])

			if len(titles) > 0 {
				movie.Name = titles[0]
			}

			covers, err := d.GetAttr("div.container>div.movie>div.screencap>a.bigImage>img", "src")
			if err != nil {
				log.Println(err)
			}

			//log.Println(covers[0])
			movie.Cover = covers[0]

			codes, err := d.GetText(`span[style="color:#CC0000;"]`)
			if err != nil {
				log.Println(err)
			}
			//log.Println(codes[0])

			movie.Code = codes[0]

			dates, err := d.GetText(`div.container>div.movie>div.info>p:nth-child(2)`)
			if err != nil {
				log.Println(err)
			}
			//	log.Println(dates[0])

			movie.Date = dates[0]

			lengths, err := d.GetText(`div.container>div.movie>div.info>p:nth-child(3)`)
			if err != nil {
				log.Println(err)
			}
			//log.Println(lengths[0])

			movie.Length = lengths[0]

			tags, err := d.GetText(`div.container>div.movie>div.info>p>span.genre>a`)
			if err != nil {
				log.Println(err)
			}

			for _, t := range tags {
				//log.Println(t)
				if strings.Contains(t, "六合彩") ||
					strings.Contains(t, "下单") ||
					strings.Contains(t, "北京") ||
					strings.Contains(t, "网投") ||
					strings.Contains(t, "下app") ||
					strings.Contains(t, "送") ||
					strings.Contains(t, "存") ||
					strings.Contains(t, "反水") {
					continue
				}
				movie.Tags += t + " "
			}

			//		actresses, err := d.GetAttr(`div.container>div.movie>div.col-md-3>ul>div>li>a>img`, "title")
			actresses, err := d.GetAttr(`div.container>div.movie>div.info>ul>div>li>div>a`, "title")
			//actresses, err := d.GetHtml(`div.container>div.movie>div.col-md-3>p:nth-child(10)>span>a`)
			if err != nil {
				log.Println(err)
			}
			if len(actresses) > 0 {
				//log.Println(actresses[0])
				movie.Actress = actresses[0]
			} else {
				//log.Println("Unknown actor")
				movie.Actress = "Unknown actor"
			}

			scripts, err := d.GetText(`body>script:nth-child(9)`)
			if err != nil {
				log.Println(err)
			}

			samples, err := d.GetAttr("a.sample-box", "href")
			if err != nil {
				log.Println(err)
			}
			movie.Samples = make(map[string]string, len(samples))
			for _, sample := range samples {
				movie.Samples[sample] = sample
			}

			//log.Println(scripts[0])
			params := scripts[0]
			params = strings.Replace(params, ";", "&", -1)
			params = strings.Replace(params, "var", "", -1)
			params = strings.Replace(params, " ", "", -1)
			params = strings.Replace(params, "\n", "", -1)
			params = strings.Replace(params, "\t", "", -1)
			params = strings.Replace(params, "\r", "", -1)
			params = strings.Replace(params, "'", "", -1)
			params += "lang=zh"
			//log.Println(params)

			magneturl := "https://www.javbus2.pw/ajax/uncledatoolsbyajax.php?" + params
			magneturl = fmt.Sprintf("%s&floor=%d", magneturl, time.Now().Unix()%100)
			//log.Println(magneturl)

			req2, err := http.NewRequest("GET", magneturl, nil)
			if err != nil {
				log.Println(err)
				return
			}

			cookies := resp.Cookies()
			for _, cie := range cookies {
				//log.Printf("%+v\n", cie)
				req2.AddCookie(cie)
			}

			//existm := &http.Cookie{Name:"existmag", Value:"all"}
			//req2.AddCookie(existm)

			req2.Header.Set("referer", link)
			req2.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.132 Safari/537.36")
			req2.Header.Set("x-requested-with", "XMLHttpRequest")
			//req2.Header.Set("accept", "*/*")
			//req2.Header.Set("accept-encoding","gzip, deflate, br")
			//req2.Header.Set("accept-language","zh-CN,zh;q=0.9")

			resp2, err := d.client.Do(req2)
			if err != nil {
				log.Println("Failed to download: ", link)
				return
			}

			magnet, err := ioutil.ReadAll(resp2.Body)
			if err != nil {
				log.Println("Cannot get the response body")
				return
			}

			defer resp.Body.Close()

			mags := MagR.FindAllStringSubmatch(string(magnet), -1)

			movie.Magnets = make(map[string]string, len(mags))
			for _, mag := range mags {
				movie.Magnets[mag[1]] = mag[1]
			}

			if movie.IsFiltered() {
				return
			}

			go d.DownloadMovie(&movie)

			d.success <- link
		}(link)
	}

	wg.Wait()
}

func (d *Downloader) GetHtml(rule string) ([]string, error) {
	var (
		res = make([]string, 0) //for leaf
		wg  sync.WaitGroup
		mu  sync.Mutex
	)

	d.doc.Find(rule).Each(func(ix int, sl *goquery.Selection) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			content, _ := sl.Html()
			mu.Lock()
			res = append(res, content)
			mu.Unlock()

		}()
	})

	wg.Wait()
	return res, nil
}

func (d *Downloader) GetText(rule string) ([]string, error) {
	var (
		res = make([]string, 0) //for leaf
		wg  sync.WaitGroup
		mu  sync.Mutex
	)

	d.doc.Find(rule).Each(func(ix int, sl *goquery.Selection) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			res = append(res, sl.Text())
			mu.Unlock()
		}()
	})
	wg.Wait()
	return res, nil
}

func (d *Downloader) GetAttr(rule, attr string) ([]string, error) {
	var (
		res = make([]string, 0) //for leaf
		wg  sync.WaitGroup
		mu  sync.Mutex
	)

	d.doc.Find(rule).Each(func(ix int, sl *goquery.Selection) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			attr, ok := sl.Attr(attr)
			if ok {
				mu.Lock()
				res = append(res, attr)
				mu.Unlock()
			}
		}()
	})
	wg.Wait()
	return res, nil
}

func (d *Downloader) DownloadMovie(movie *Movie) error {
	res, err := json.Marshal(movie)
	if err != nil {
		return err
	}

	AppendToFile("javlist.txt", res)
	AppendToFile("javinfo.txt", []byte(fmt.Sprintf("%s", movie)))

	fields := strings.Split(movie.Link, "/")
	name := fields[len(fields)-1]
	d.DownloadImage(movie.Cover, movie.Actress+"_"+name+".jpg")
	var i = 0
	for _, s := range movie.Samples {

		d.DownloadImage(s, fmt.Sprintf("%s_%s_%d.jpg", movie.Actress, name, i))
		i++
	}

	Downloaded[movie.Link] = movie.Code

	return nil

}

func (d *Downloader) DownloadImage(link, name string) error {
	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		log.Println("Create new reqeust error ", err.Error())
		return err
	}
	req.Header.Add("User-Agent", d.userAgent)
	resp, err := d.client.Do(req)
	if err != nil {
		log.Println("Download ", link, " error: ", err.Error())
		return err
	}

	//dump, _ := httputil.DumpRequestOut(req, true)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Cannot get the response body")
		return err
	}
	defer resp.Body.Close()

	file, err := os.Create("img/" + name)
	//file, err := os.Create(d.root+name)
	if err != nil {
		log.Println("Cannot create new file for link: ", link, " error: ", err.Error())
		return err
	}
	defer file.Close()
	file.Write(body)

	return nil

}

func AppendToFile(name string, data []byte) {
	file, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	file.Write(data)
}

func SaveToFile(name string, data []byte) {
	file, err := os.Create(name)
	if err != nil {
		log.Println("Cannot create file: ", name, " ", err.Error())
		return
	}

	file.Write(data)
	file.Sync()
	defer file.Close()
}

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	content, err := ioutil.ReadFile("downloaded.txt")
	if err != nil {
		fmt.Println(err)
	}

	err = json.Unmarshal(content, &Downloaded)
	if err != nil {
		fmt.Println(err)
	}

	content, err = ioutil.ReadFile("filtered.txt")
	if err != nil {
		fmt.Println(err)
	}

	err = json.Unmarshal(content, &Filtered)
	if err != nil {
		fmt.Println(err)
	}

	go func() {
		tick := time.Tick(time.Duration(time.Second * 30))
		for range tick {
			ms, err := json.Marshal(Downloaded)
			if err != nil {
				panic(err)
			}

			SaveToFile("downloaded.txt", ms)

			ms, err = json.Marshal(Filtered)
			if err != nil {
				panic(err)
			}

			SaveToFile("filtered.txt", ms)
		}
	}()
}

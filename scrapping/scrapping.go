package scrapping

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type extractedBooth struct {
	title string
	link  string
}

func Scrapping(term string) string {
	var result string
	var booths []extractedBooth
	c := make(chan []extractedBooth)

	var baseURL string = "https://www.albamon.com/search/Recruit?Keyword=" + term + "&SiCode=I000&GuCode=I220&PageSize=20"

	totalPages := getPages(baseURL)

	for i := 1; i <= totalPages; i++ {
		go getPage(i, baseURL, c)
	}

	for i := 1; i <= totalPages; i++ {
		extractedBooths := <-c
		booths = append(booths, extractedBooths...)
	}

	for i := 0; i < len(booths); i++ {
		// fmt.Println(booths[i].title)
		slice := strings.Split(booths[i].title, " ")
		for j := 0; j < len(slice)-2; j++ {
			result += slice[j] + " "
		}
		result += "\n"
	}

	return strings.TrimSpace(result)
}

func getPage(page int, url string, mainC chan<- []extractedBooth) {
	var booths []extractedBooth
	c := make(chan extractedBooth)

	pageURL := url + "&Page=" + strconv.Itoa(page)
	fmt.Println("Requesting:", pageURL)

	res, err := http.Get(pageURL)
	CheckErr(err)
	checkStatus(res)

	doc, err := goquery.NewDocumentFromReader(res.Body)
	defer res.Body.Close()
	CheckErr(err)

	searchBooths := doc.Find(".smResult>.booth")
	searchBooths.Each(func(i int, booth *goquery.Selection) {
		go extractBooth(booth, c)
	})

	for i := 0; i < searchBooths.Length(); i++ {
		booth := <-c
		booths = append(booths, booth)
	}

	mainC <- booths
}

func extractBooth(booth *goquery.Selection, c chan<- extractedBooth) {
	title, _ := booth.Find(".list>dt>a").Attr("title")
	link, _ := booth.Find(".list>dt>a").Attr("href")

	c <- extractedBooth{
		title: title,
		link:  link,
	}
}

func getPages(url string) int {
	pages := 0
	res, err := http.Get(url)
	CheckErr(err)
	checkStatus(res)

	doc, err := goquery.NewDocumentFromReader(res.Body)
	defer res.Body.Close()
	CheckErr(err)

	doc.Find(".listPaging").Each(func(i int, page *goquery.Selection) {
		pages = page.Find("a").Length()
	})

	return pages + 1
}

// CheckErr returns err and kill program
func CheckErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func checkStatus(res *http.Response) {
	if res.StatusCode != 200 {
		log.Fatalln("Request failed with Status:", res.StatusCode)
	}
}

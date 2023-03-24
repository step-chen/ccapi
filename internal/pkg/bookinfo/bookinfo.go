package bookinfo

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	gojsonq "github.com/thedevsaddam/gojsonq/v2"

	cclog "ccapi/internal/pkg/log"
)

type Book struct {
	Status      int      `json:"status"`
	ISBN        string   `json:"isbn"`
	Title       string   `json:"title"`
	SubTitle    string   `json:"subtitle"`
	OriTitle    string   `json:"orititle"`
	Author      []string `json:"author"`
	Publisher   string   `json:"publisher"`
	Published   string   `json:"published"`
	PageCount   int      `json:"pageCount"`
	Rating      float32  `json:"rating"`
	Designed    string   `json:"designed"`
	Price       string   `json:"price"`
	URL         string   `json:"url"`
	Translator  []string `json:"translator"`
	CoverURL    string   `json:"coverUrl"`
	Description string   `json:"description"`
	AuthorIntro string   `json:"authorIntro"`
}

func GetByIsbnFromDouban(c *gin.Context) {
	var bookInfo Book
	pstrBody, strIsbn, surl := fetchBodyString(c)
	if pstrBody == nil || *pstrBody == "" {
		c.IndentedJSON(http.StatusNotFound, gin.H{"Loction": surl})
		return
	}

	bookInfo.ISBN = strIsbn
	bookInfo.Status = http.StatusOK

	fetchBasicInfo(&bookInfo, pstrBody)
	fetchSubTitle(&bookInfo, pstrBody)
	fetchOriginalTitle(&bookInfo, pstrBody)
	fetchRating(&bookInfo, pstrBody)
	fetchTranslator(&bookInfo, pstrBody)
	fetchPublisher(&bookInfo, pstrBody)
	fetchCoverURL(&bookInfo, pstrBody)
	fetchPublished(&bookInfo, pstrBody)
	fetchPageCount(&bookInfo, pstrBody)
	fetchPrice(&bookInfo, pstrBody)
	fetchDesigned(&bookInfo, pstrBody)
	fetchDescription(&bookInfo, pstrBody)
	fetchAuthorIntro(&bookInfo, pstrBody)

	c.IndentedJSON(http.StatusOK, bookInfo)

	/* bytRespJson, err := json.Marshal(bookInfo)
	if err != nil {
		cclog.LogErr("%s | URL: %s", err.Error(), surl)
		return
	}

	fmt.Println(string(bytRespJson[:])) */

	// c.IndentedJSON(http.StatusOK, gin.H{"Loction": strLoc})*/
}

func fetchBodyString(c *gin.Context) (pstrBody *string, strIsbn string, surl string) {
	strIsbn = c.Param("isbn")
	if strIsbn == "" {
		return nil, "", ""
	}

	var strBuilder strings.Builder
	strBuilder.WriteString("https://book.douban.com/isbn/")
	strBuilder.WriteString(strIsbn)
	strBuilder.WriteString("/")
	surl = strBuilder.String()

	req, _ := http.NewRequest("GET", surl, nil)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.130 Safari/537.36 OPR/66.0.3515.115")
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		cclog.LogErr("%s | URL: %s", err.Error(), surl)
		return nil, strIsbn, surl
	}

	if resp != nil {
		defer resp.Body.Close()
	}
	cclog.Log("%s: %d | URL: %s", "DouBan status", resp.StatusCode, surl)

	if resp.StatusCode != 200 {
		return nil, "", surl
	}

	bytBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		cclog.LogErr("%s", err.Error())
		return nil, strIsbn, surl
	}

	strBody := string(bytBody[:])
	strBody = strings.Replace(strBody, "\n", "", -1)
	strBody = strings.Replace(strBody, "\r", "", -1)

	return &strBody, strIsbn, surl
}

func fetchBasicInfo(bookInfo *Book, pstrBody *string) {
	patternString := `application/ld\+json">(.*?)</script>`
	reg := regexp.MustCompile(patternString)
	pattern := reg.FindStringSubmatch(*pstrBody)

	if pattern != nil && len(pattern) > 1 {
		title := gojsonq.New().FromString(pattern[1]).Find("name")
		if title != nil && title != "" {
			bookInfo.Title = strings.TrimSpace(fmt.Sprint(title))
		}

		arrAuthor := gojsonq.New().FromString(pattern[1]).Find("author").([]interface{})
		for _, val := range arrAuthor {
			t := val.(map[string]interface{})["name"]
			if t != nil && t != "" {
				bookInfo.Author = append(bookInfo.Author, strings.TrimSpace(fmt.Sprint(t)))
			}
		}

		doubanUrl := gojsonq.New().FromString(pattern[1]).Find("url")
		if doubanUrl != nil && doubanUrl != "" {
			bookInfo.URL = strings.TrimSpace(fmt.Sprint(doubanUrl))
		}
	}
}

func fetchSubTitle(bookInfo *Book, pstrBody *string) {
	patternString := `副标题:</span>(.*?)<br/>`
	reg := regexp.MustCompile(patternString)
	pattern := reg.FindStringSubmatch(*pstrBody)

	if pattern != nil && len(pattern) > 1 {
		bookInfo.SubTitle = strings.TrimSpace(pattern[1])
	}
}

func fetchOriginalTitle(bookInfo *Book, pstrBody *string) {
	patternString := `原作名:</span>(.*?)<br/>`
	reg := regexp.MustCompile(patternString)
	pattern := reg.FindStringSubmatch(*pstrBody)

	if pattern != nil && len(pattern) > 1 {
		bookInfo.OriTitle = strings.TrimSpace(pattern[1])
	}
}

func fetchRating(bookInfo *Book, pstrBody *string) {
	patternString := `v:average">(.*?)</`
	reg := regexp.MustCompile(patternString)
	pattern := reg.FindStringSubmatch(*pstrBody)

	if pattern != nil && len(pattern) > 1 {
		rating, err := strconv.ParseFloat(strings.TrimSpace(pattern[1]), 32)
		if err != nil {
			cclog.LogErr("%s", err.Error())
		} else {
			bookInfo.Rating = float32(rating)
		}
	}
}

func fetchTranslator(bookInfo *Book, pstrBody *string) {
	patternString := `译者</span>:(.*?)</span>`
	reg := regexp.MustCompile(patternString)
	pattern := reg.FindStringSubmatch(*pstrBody)

	if pattern != nil && len(pattern) > 1 {
		patternString = `">(.*?)</a>`
		reg = regexp.MustCompile(patternString)
		arrPattern := reg.FindAllStringSubmatch(pattern[1], -1)
		if arrPattern != nil {
			for _, val := range arrPattern {
				if len(val) > 1 {
					bookInfo.Translator = append(bookInfo.Translator, strings.TrimSpace(val[1]))
				}
			}
		}
	}
}

func fetchPublisher(bookInfo *Book, pstrBody *string) {
	patternString := `出版社:</span>(.*?)">(.*?)<`
	reg := regexp.MustCompile(patternString)
	pattern := reg.FindStringSubmatch(*pstrBody)

	if pattern != nil && len(pattern) > 2 {
		bookInfo.Publisher = strings.TrimSpace(pattern[2])
	}
}

func fetchCoverURL(bookInfo *Book, pstrBody *string) {
	patternString := `data-pic="(.*?)"`
	reg := regexp.MustCompile(patternString)
	pattern := reg.FindStringSubmatch(*pstrBody)

	if pattern != nil && len(pattern) > 1 {
		bookInfo.CoverURL = strings.TrimSpace(pattern[1])
	}
}

func fetchPublished(bookInfo *Book, pstrBody *string) {
	patternString := `出版年:</span>(.*?)<br/>`
	reg := regexp.MustCompile(patternString)
	pattern := reg.FindStringSubmatch(*pstrBody)

	if pattern != nil && len(pattern) > 1 {
		bookInfo.Published = strings.TrimSpace(pattern[1])
	}
}

func fetchPageCount(bookInfo *Book, pstrBody *string) {
	patternString := `页数:</span>(.*?)<br/>`
	reg := regexp.MustCompile(patternString)
	pattern := reg.FindStringSubmatch(*pstrBody)

	if pattern != nil && len(pattern) > 1 {
		count, err := strconv.Atoi(strings.TrimSpace(pattern[1]))
		if err != nil {
			cclog.LogErr("%s", err.Error())
		} else {
			bookInfo.PageCount = count
		}
	}
}

func fetchPrice(bookInfo *Book, pstrBody *string) {
	patternString := `定价:</span>(.*?)<br/>`
	reg := regexp.MustCompile(patternString)
	pattern := reg.FindStringSubmatch(*pstrBody)

	if pattern != nil && len(pattern) > 1 {
		bookInfo.Price = strings.TrimSpace(pattern[1])
	}
}

func fetchDesigned(bookInfo *Book, pstrBody *string) {
	patternString := `装帧:</span>(.*?)<br/>`
	reg := regexp.MustCompile(patternString)
	pattern := reg.FindStringSubmatch(*pstrBody)

	if pattern != nil && len(pattern) > 1 {
		bookInfo.Designed = strings.TrimSpace(pattern[1])
	}
}

func fetchDescription(bookInfo *Book, pstrBody *string) {
	patternString := `内容简介(.*?)<h2>`
	reg := regexp.MustCompile(patternString)
	arrPattern := reg.FindStringSubmatch(*pstrBody)

	if arrPattern != nil && len(arrPattern) > 1 {
		patternString = `class="intro">(.*?)</div>`
		reg = regexp.MustCompile(patternString)
		arrContent := reg.FindAllStringSubmatch(arrPattern[1], -1)
		for _, val := range arrContent {
			if len(val) > 1 {
				if strings.Contains(val[1], "展开全部") {
					continue
				} else {
					val[1] = strings.Replace(val[1], "\u003c/p\u003e", "\n", -1)
					val[1] = strings.Replace(val[1], "\u003cp\u003e", "", -1)
					bookInfo.Description = strings.TrimSpace(val[1])
				}
			}
		}
	}
}

func fetchAuthorIntro(bookInfo *Book, pstrBody *string) {
	patternString := `作者简介(.*?)<h2>`
	reg := regexp.MustCompile(patternString)
	arrPattern := reg.FindStringSubmatch(*pstrBody)

	if arrPattern != nil && len(arrPattern) > 1 {
		patternString = `class="intro">(.*?)</div>`
		reg = regexp.MustCompile(patternString)
		arrContent := reg.FindAllStringSubmatch(arrPattern[1], -1)
		for _, val := range arrContent {
			if len(val) > 1 {
				if strings.Contains(val[1], "展开全部") {
					continue
				} else {
					val[1] = strings.Replace(val[1], "\u003c/p\u003e", "\n", -1)
					val[1] = strings.Replace(val[1], "\u003cp\u003e", "", -1)
					bookInfo.AuthorIntro = strings.TrimSpace(val[1])
				}
			}
		}
	}
}

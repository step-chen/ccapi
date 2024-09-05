package bookinfo

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	gojsonq "github.com/thedevsaddam/gojsonq/v2"

	cclog "github.com/step-chen/ccapi/internal/pkg/log"
)

func fetchBodyByNameFromDouban(strName string, nCount int) (strBodies []string, strSearchUrl string) {
	strSearchUrl = buildUrl("https://search.douban.com/book/subject_search?", strName)

	strBody := fetchUrlBodyContent(strSearchUrl)
	if strBody == "" {
		return nil, strSearchUrl
	}

	patternString := `"url": "https://book.douban.com/subject/(.*?)/"`
	reg := regexp.MustCompile(patternString)
	pattern := reg.FindAllStringSubmatch(strBody, nCount)

	for _, v := range pattern {
		strUrl := buildUrl("https://book.douban.com/subject/", v[1])
		strBody := fetchUrlBodyContent(strUrl)

		strBody = strings.Replace(strBody, "\n", "", -1)
		strBody = strings.Replace(strBody, "\r", "", -1)

		strBodies = append(strBodies, strBody)
	}

	return strBodies, strSearchUrl
}

func fetchBodyByIsbnFromDouban(strIsbn string) (strBody string, strUrl string) {
	var strBuilder strings.Builder
	strBuilder.WriteString("https://book.douban.com/isbn/")
	strBuilder.WriteString(strIsbn)
	strBuilder.WriteString("/")
	strUrl = strBuilder.String()

	req, _ := http.NewRequest("GET", strUrl, nil)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.130 Safari/537.36 OPR/66.0.3515.115")
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		cclog.LogErr("%s | URL: %s", err.Error(), strUrl)
		return "", strUrl
	}

	if resp != nil {
		defer resp.Body.Close()
	}
	cclog.Log("%s: %d | URL: %s", "DouBan status", resp.StatusCode, strUrl)

	if resp.StatusCode != 200 {
		return "", strUrl
	}

	bytBody, err := io.ReadAll(resp.Body)
	if err != nil {
		cclog.LogErr("%s", err.Error())
		return "", strUrl
	}

	strBody = string(bytBody[:])
	strBody = strings.Replace(strBody, "\n", "", -1)
	strBody = strings.Replace(strBody, "\r", "", -1)

	return strBody, strUrl
}

func fetchBasicInfoFromDouban(bookInfo *Book, strBody string) {
	patternString := `application/ld\+json">(.*?)</script>`
	reg := regexp.MustCompile(patternString)
	pattern := reg.FindStringSubmatch(strBody)

	if len(pattern) > 1 {
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

		doubanISBN := gojsonq.New().FromString(pattern[1]).Find("isbn")
		if doubanISBN != nil && doubanISBN != "" {
			bookInfo.ISBN = strings.TrimSpace(fmt.Sprint(doubanISBN))
		}
	}
}

func fetchSubTitleFromDouban(bookInfo *Book, strBody string) {
	patternString := `副标题:</span>(.*?)<br/>`
	reg := regexp.MustCompile(patternString)
	pattern := reg.FindStringSubmatch(strBody)

	if len(pattern) > 1 {
		bookInfo.SubTitle = strings.TrimSpace(pattern[1])
	}
}

func fetchOriginalTitleFromDouban(bookInfo *Book, strBody string) {
	patternString := `原作名:</span>(.*?)<br/>`
	reg := regexp.MustCompile(patternString)
	pattern := reg.FindStringSubmatch(strBody)

	if len(pattern) > 1 {
		bookInfo.OriTitle = strings.TrimSpace(pattern[1])
	}
}

func fetchRatingFromDouban(bookInfo *Book, strBody string) {
	patternString := `v:average">(.*?)</`
	reg := regexp.MustCompile(patternString)
	pattern := reg.FindStringSubmatch(strBody)

	if len(pattern) > 1 {
		rating, err := strconv.ParseFloat(strings.TrimSpace(pattern[1]), 32)
		if err != nil {
			cclog.LogErr("%s", err.Error())
		} else {
			bookInfo.Rating = float32(rating)
		}
	}
}

func fetchTranslatorFromDouban(bookInfo *Book, strBody string) {
	patternString := `译者</span>:(.*?)</span>`
	reg := regexp.MustCompile(patternString)
	pattern := reg.FindStringSubmatch(strBody)

	if len(pattern) > 1 {
		patternString = `">(.*?)</a>`
		reg = regexp.MustCompile(patternString)
		arrPattern := reg.FindAllStringSubmatch(pattern[1], -1)

		for _, val := range arrPattern {
			if len(val) > 1 {
				bookInfo.Translator = append(bookInfo.Translator, strings.TrimSpace(val[1]))
			}
		}
	}
}

func fetchPublisherFromDouban(bookInfo *Book, strBody string) {
	patternString := `出版社:</span>(.*?)">(.*?)<`
	reg := regexp.MustCompile(patternString)
	pattern := reg.FindStringSubmatch(strBody)

	if len(pattern) > 2 {
		bookInfo.Publisher = strings.TrimSpace(pattern[2])
	}
}

func fetchCoverURLFromDouban(bookInfo *Book, strBody string) {
	patternString := `data-pic="(.*?)"`
	reg := regexp.MustCompile(patternString)
	pattern := reg.FindStringSubmatch(strBody)

	if len(pattern) > 1 {
		bookInfo.CoverURL = strings.TrimSpace(pattern[1])
	}
}

func fetchPublishedFromDouban(bookInfo *Book, strBody string) {
	patternString := `出版年:</span>(.*?)<br/>`
	reg := regexp.MustCompile(patternString)
	pattern := reg.FindStringSubmatch(strBody)

	if len(pattern) > 1 {
		bookInfo.Published = strings.TrimSpace(pattern[1])
	}
}

func fetchPageCountFromDouban(bookInfo *Book, strBody string) {
	patternString := `页数:</span>(.*?)<br/>`
	reg := regexp.MustCompile(patternString)
	pattern := reg.FindStringSubmatch(strBody)

	if len(pattern) > 1 {
		count, err := strconv.Atoi(strings.TrimSpace(pattern[1]))
		if err != nil {
			cclog.LogErr("%s", err.Error())
		} else {
			bookInfo.PageCount = count
		}
	}
}

func fetchPriceFromDouban(bookInfo *Book, strBody string) {
	patternString := `定价:</span>(.*?)<br/>`
	reg := regexp.MustCompile(patternString)
	pattern := reg.FindStringSubmatch(strBody)

	if len(pattern) > 1 {
		bookInfo.Price = strings.TrimSpace(pattern[1])
	}
}

func fetchDesignedFromDouban(bookInfo *Book, strBody string) {
	patternString := `装帧:</span>(.*?)<br/>`
	reg := regexp.MustCompile(patternString)
	pattern := reg.FindStringSubmatch(strBody)

	if len(pattern) > 1 {
		bookInfo.Designed = strings.TrimSpace(pattern[1])
	}
}

func fetchDescriptionFromDouban(bookInfo *Book, strBody string) {
	patternString := `内容简介(.*?)<h2>`
	reg := regexp.MustCompile(patternString)
	arrPattern := reg.FindStringSubmatch(strBody)

	if len(arrPattern) > 1 {
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

func fetchAuthorIntroFromDouban(bookInfo *Book, strBody string) {
	patternString := `作者简介(.*?)<h2>`
	reg := regexp.MustCompile(patternString)
	arrPattern := reg.FindStringSubmatch(strBody)

	if len(arrPattern) > 1 {
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

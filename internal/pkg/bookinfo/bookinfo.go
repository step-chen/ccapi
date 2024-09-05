package bookinfo

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// https://search.douban.com/book/subject_search?search_text=%E8%A6%81%E5%91%BD%E8%BF%98%E6%98%AF%E8%A6%81%E7%81%B5%E9%AD%82
// https://book.douban.com/subject/36164018/

func GetByNameFromDouban(c *gin.Context) {
	nCount := 1
	strName := buildParams("search_text", c.Param("name"))
	strCount := c.Param("count")
	if strName == "" {
		return
	}
	if strCount != "" {
		count, err := strconv.Atoi(strCount)
		if err == nil && count > 0 {
			nCount = count
		}
	}

	fmt.Printf("GetByNameFromDouban: %s, %d\n", strName, nCount)
	var bookInfos []Book
	strBodies, strSearchUrl := fetchBodyByNameFromDouban(strName, nCount)
	if len(strBodies) == 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"Loction": strSearchUrl})
		return
	}
	for _, strBody := range strBodies {
		bookInfo := getInfoByBodyContent(strBody)
		bookInfos = append(bookInfos, bookInfo)
	}

	c.Header("Access-Control-Allow-Origin", "*")
	c.JSON(http.StatusOK, bookInfos)

	// fmt.Println(strBodies, strUrls, strSearchUrl)
}

func GetByIsbnFromDouban(c *gin.Context) {
	strIsbn := c.Param("isbn")
	if strIsbn == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"Loction": "No ISBN provided."})
		return
	}

	strBody, strUrl := fetchBodyByIsbnFromDouban(strIsbn)
	if strBody == "" {
		c.IndentedJSON(http.StatusNotFound, gin.H{"Loction": strUrl})
		return
	}

	bookInfo := getInfoByBodyContent(strBody)

	c.Header("Access-Control-Allow-Origin", "*")
	c.JSON(http.StatusOK, bookInfo)
}

func getInfoByBodyContent(strBody string) (bookInfo Book) {
	bookInfo.Status = http.StatusOK

	fetchBasicInfoFromDouban(&bookInfo, strBody)
	fetchSubTitleFromDouban(&bookInfo, strBody)
	fetchOriginalTitleFromDouban(&bookInfo, strBody)
	fetchRatingFromDouban(&bookInfo, strBody)
	fetchTranslatorFromDouban(&bookInfo, strBody)
	fetchPublisherFromDouban(&bookInfo, strBody)
	fetchCoverURLFromDouban(&bookInfo, strBody)
	fetchPublishedFromDouban(&bookInfo, strBody)
	fetchPageCountFromDouban(&bookInfo, strBody)
	fetchPriceFromDouban(&bookInfo, strBody)
	fetchDesignedFromDouban(&bookInfo, strBody)
	fetchDescriptionFromDouban(&bookInfo, strBody)
	fetchAuthorIntroFromDouban(&bookInfo, strBody)

	return bookInfo
}

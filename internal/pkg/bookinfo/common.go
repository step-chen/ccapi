package bookinfo

import (
	"io"
	"net/http"
	"net/url"
	"strings"

	cclog "github.com/step-chen/ccapi/internal/pkg/log"
)

func buildUrl(baseUrl string, str ...string) string {
	var strBuilder strings.Builder
	strBuilder.WriteString(baseUrl)
	for _, s := range str {
		strBuilder.WriteString(s)
	}
	return strBuilder.String()
}

func buildParams(k string, v string) string {
	if v == "" {
		return ""
	}

	params := url.Values{}
	params.Add(k, v)

	return params.Encode()
}

func fetchUrlBodyContent(strUrl string) (strBody string) {
	req, _ := http.NewRequest("GET", strUrl, nil)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.130 Safari/537.36 OPR/66.0.3515.115")
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		cclog.LogErr("%s | URL: %s", err.Error(), strUrl)
		return ""
	}

	if resp != nil {
		defer resp.Body.Close()
	}
	cclog.Log("%s: %d | URL: %s", "DouBan status", resp.StatusCode, strUrl)

	if resp.StatusCode != 200 {
		return ""
	}

	bytBody, err := io.ReadAll(resp.Body)
	if err != nil {
		cclog.LogErr("%s", err.Error())
		return ""
	}

	return string(bytBody[:])
}

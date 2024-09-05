# ccapi

## Description
This is a GO implementation of a REST API that retrieves book information using ISBN or book title. Currently, it only supports fetching from Douban.
这是一个 GO 实现的通过ISBN或书名获取书籍信息的 REST API，当前仅支持从豆瓣获取。

## Usage
```bash
CGO_ENABLED=0 GOOS=linux go build -o ./apps/ccapi -ldflags '-s -w --extldflags "-static -fpic"' main.go
./apps/ccapi
```
## API
```bash
curl -X GET "http://localhost:5003/book/douban/isbn/9787542679307"
curl -X GET "http://localhost:5003/book/douban/name/要命还是要灵魂/2"
```
package bookinfo

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

package web

type Page struct {
	Title       string
	Description string
}

func NewPage(title string) Page {
	return Page{
		Title: title,
	}
}

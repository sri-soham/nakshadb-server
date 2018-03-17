package helper

import (
	"fmt"
	"html/template"
	"strings"
)

type Paginator struct {
	page        int
	per_page    int
	count       int
	url         string
	first       int
	last        int
	total_pages int
}

func (p *Paginator) Process() {
	p.total_pages = p.count / p.per_page
	if p.count%p.per_page > 0 {
		p.total_pages++
	}
	p.first = p.page - 2
	if p.first < 1 {
		p.first = 1
	}
	p.last = p.first + 4
	if p.last > p.total_pages {
		p.last = p.total_pages
		if p.last-4 > 0 {
			p.first = p.last - 4
		}
	}
}

func (p *Paginator) Links() template.HTML {
	links := make([]string, 0)
	if p.total_pages > 3 {
		if p.page > 2 {
			links = append(links, p.getLink(1, "&ltrif;&ltrif;"))
		}
		if p.page > 1 {
			links = append(links, p.getLink(p.page-1, "&ltrif;"))
		}
		for i := p.first; i <= p.last; i++ {
			links = append(links, p.getLink(i, fmt.Sprintf("%v", i)))
		}
		if p.page < p.total_pages {
			links = append(links, p.getLink(p.page+1, "&rtrif;"))
		}
		if p.page < (p.total_pages - 1) {
			links = append(links, p.getLink(p.total_pages, "&rtrif;&rtrif;"))
		}
	} else {
		for i := 1; i <= p.total_pages; i++ {
			links = append(links, p.getLink(i, fmt.Sprintf("%v", i)))
		}
	}

	return template.HTML(strings.Join(links, ""))
}

func (p *Paginator) Text() string {
	min := (p.page - 1) * p.per_page
	max := min + p.per_page
	min += 1
	if max > p.count {
		max = p.count
	}

	return fmt.Sprintf("%v - %v of %v", min, max, p.count)
}

func (p *Paginator) getLink(page int, text string) string {
	if page == p.page {
		return fmt.Sprintf("<span class=\"current-page\">%v</span>", fmt.Sprintf("%v", page))
	} else {
		href := strings.Replace(p.url, "{page}", fmt.Sprintf("%v", page), 1)
		return fmt.Sprintf("<a href=\"%v\">%s</a>", href, text)
	}
}

func MakePaginator(page int, per_page int, count int, url string) Paginator {
	paginator := Paginator{
		page,
		per_page,
		count,
		url,
		0,
		0,
		0,
	}
	paginator.Process()

	return paginator
}

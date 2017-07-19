package main

type Pagestore interface {
	NewPage(url string, finished bool) (*PageData, error)
	GetPage(id string) (*PageData, error)
	SavePage(p *PageData) error
	FinishPage(p *PageData) error
}

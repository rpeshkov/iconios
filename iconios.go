package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"

	"strconv"

	"github.com/mediocregopher/radix.v2/redis"
	uuid "github.com/satori/go.uuid"
)

type RedisPagestore struct {
	c *redis.Client
}

func (r *RedisPagestore) FinishPage(p *PageData) error {
	key := "page:" + p.ID.String()
	exists, err := r.c.Cmd("EXISTS", key).Int()

	if err != nil {
		return err
	}

	if exists == 0 {
		return fmt.Errorf("Page '%s' not exist", key)
	}

	return r.c.Cmd("HSET", key, "finished", "1").Err
}

func (r *RedisPagestore) NewPage(url string, finished bool) (*PageData, error) {
	p := &PageData{uuid.NewV4(), url, finished}
	if err := r.SavePage(p); err != nil {
		return nil, err
	}

	return p, nil
}

func (r *RedisPagestore) SavePage(p *PageData) error {
	key := "page:" + p.ID.String()
	return r.c.Cmd("HMSET", key, "url", p.URL, "finished", false).Err
}

func (r *RedisPagestore) GetPage(id string) (*PageData, error) {
	key := "page:" + id
	exists, err := r.c.Cmd("EXISTS", key).Int()

	if err != nil {
		return nil, err
	}

	if exists == 0 {
		return nil, fmt.Errorf("Page '%s' does not exist", key)
	}

	m, err := r.c.Cmd("HGETALL", key).Map()
	if err != nil {
		return nil, fmt.Errorf("get page data failed: %v", err)
	}

	f, err := strconv.ParseBool(m["finished"])
	if err != nil {
		f = false
	}

	return &PageData{ID: uuid.FromStringOrNil(id), URL: m["url"], Finished: f}, nil
}

func NewRedisPagestore(network, addr string) (*RedisPagestore, error) {
	c, err := redis.Dial(network, addr)

	if err != nil {
		return nil, fmt.Errorf("unable to create RedisPagestore: %v", err)
	}

	return &RedisPagestore{c}, nil
}

type Pagestore interface {
	NewPage(url string, finished bool) (*PageData, error)
	GetPage(id string) (*PageData, error)
	SavePage(p *PageData) error
	FinishPage(p *PageData) error
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	t, err := createTemplates()

	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load templates: %v", err)
		os.Exit(2)
	}

	ps, err := NewRedisPagestore("tcp", "localhost:6379")

	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize pagestore: %v", err)
		os.Exit(2)
	}

	staticHandler := http.FileServer(http.Dir("static"))
	storageHandler := http.FileServer(http.Dir("storage"))

	http.Handle("/static/", http.StripPrefix("/static/", staticHandler))
	http.Handle("/storage/", http.StripPrefix("/storage/", storageHandler))
	http.Handle("/finish/", finishHandler(t, ps))
	http.Handle("/opn/", openHandler(t, ps))
	http.Handle("/", indexHandler(t, ps))

	http.ListenAndServe(":"+port, nil)
}

func createTemplates() (*template.Template, error) {
	t, err := template.ParseGlob(filepath.Join("tmpl", "*.gohtml"))

	if err != nil {
		return nil, fmt.Errorf("unable to load index template: %v", err)
	}

	return t, nil
}

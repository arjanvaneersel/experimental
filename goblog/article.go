package goblog

import (
	"time"
	"github.com/Machiel/slugify"
	"log"
	"errors"
	"google.golang.org/appengine/datastore"
)

type ArticleTranslation struct {
	Title string `json:"id"`
	Slug string `json:"slug"`
	Body string `json:"body"`
}

type ArticleTranslations map[string]*ArticleTranslation

type Article struct{
	ID *datastore.Key `json:"id" datastore:"-"`
	Title string `json:"title"`
	Slug string `json:"slug"`
	Body string `json:"body"`
	Created time.Time `json:"created"`
	Modified time.Time `json:"modified"`
	Published time.Time `json:"published"`
	Translations ArticleTranslations `json:"translations"`
}


func (a *Article) SetSlug() {
	if a.Slug == "" {
		a.Slug = slugify.Slugify(a.Title)
	} else {
		a.Slug = slugify.Slugify(a.Slug)
	}

	for lang, translation := range a.Translations {
		if translation.Slug == "" {
			log.Println(lang)
			translation.Slug = slugify.Slugify(translation.Title)
		} else {
			log.Println(lang)
			translation.Slug = slugify.Slugify(translation.Slug)
		}
	}
}

func (a *Article) Validates() error {
	return errors.New("Not implemented")
}

type ArticleManager interface {
	Create(article *Article) error
	GetByID(id interface{}) (*Article, error)
	GetBySlug(slug string) (*Article, error)
	GetAll(count int, order string) (*[]Article, error)
	Update(article *Article) error
	Delete(id interface{}) error
}
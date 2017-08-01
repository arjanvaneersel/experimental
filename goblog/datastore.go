package goblog

import (
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/datastore"
	"time"
	"net/http"
	"golang.org/x/net/context"
)

const ReturnAll int = -1

type DatastoreArticleManager struct {
	ctx context.Context
}

func (ds *DatastoreArticleManager) Create(article *Article) error {
	var err error
	log.Debugf(ds.ctx, "Saving article")
	if article.ID == nil {
		article.ID = datastore.NewIncompleteKey(ds.ctx, "Article", nil)
	}
	//ToDo: Add Author

	article.SetSlug()
	now := time.Now()
	article.Created = now
	article.Modified = now
	if article.Published.IsZero() {
		article.Published = now
	}

	//ToDo: Validation
	article.ID, err = datastore.Put(ds.ctx, article.ID, nil)
	if err != nil {
		return err
	}
	return nil
}

func (ds *DatastoreArticleManager) Update(article *Article) error {
	var err error

	log.Debugf(ds.ctx, "Updating article")
	if article.ID == nil {
		article.ID = datastore.NewIncompleteKey(ds.ctx, "Article", nil)
	}

	//ToDo: Validation
	article.ID, err = datastore.Put(ds.ctx, article.ID, article)
	if err != nil {
		return err
	}

	return nil
}

func (ds *DatastoreArticleManager) GetByID(id interface{}) (*Article, error) {
	key := id.(*datastore.Key)
	a := &Article{}

	err := datastore.Get(ds.ctx, key, &a)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (ds *DatastoreArticleManager) GetBySlug(slug string) (*Article, error) {
	q := datastore.NewQuery("Article").Filter("slug", slug)

	var results []Article
	keys, err := q.GetAll(ds.ctx, &results)
	if err != nil {
		return nil, err
	}

	a := &results[0]
	a.ID = keys[0]

	return a, nil
}

func (ds *DatastoreArticleManager) GetAll(limit int, order string) (*[]Article, error) {
	q := datastore.NewQuery("Article").Order(order).Limit(limit)

	var results []Article
	keys, err := q.GetAll(ds.ctx, &results)
	if err != nil {
		return nil, err
	}


	for i := 0; i < len(results); i++ {
		results[i].ID = keys[i]
	}

	return &results, nil
}

func (ds *DatastoreArticleManager) Delete(id interface{}) error {
	err := datastore.Delete(ds.ctx, id.(*datastore.Key))
	return err
}

func SetContext(ds *DatastoreArticleManager, ctx context.Context) {
	ds.ctx = ctx
}

func NewContext(ds *DatastoreArticleManager, r *http.Request) {
	ds.ctx = appengine.NewContext(r)
}

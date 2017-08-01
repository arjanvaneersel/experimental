package goblog

import (
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/user"
	"google.golang.org/appengine/log"
	"golang.org/x/net/context"
	"errors"
	"crypto/md5"
	"io"
	"strings"
	"fmt"
)

type User struct {
	ID *datastore.Key `json:"id" datastore:"-"`
	UserID string `json:"-"`
	DisplayNaem string `json:"display_name"`
	AvatarURL string `json:"avatar_url"`
	Score int `json:"score"`
}

func GetAppEngineUser(ctx context.Context) (*User, error) {
	u := user.Current(ctx)
	if u == nil {
		return nil, errors.New("Not logged in")
	}

	var user User
	user.ID = datastore.NewKey(ctx, "User", u.ID, 0, nil)
	err := datastore.Get(ctx, user.ID, &user)
	if err != nil && err != datastore.ErrNoSuchEntity {
		return nil, err
	}
	if err == datastore.ErrNoSuchEntity {
		user.UserID = u.ID
		user.DisplayNaem = u.String()
		user.AvatarURL = gravatarURL(u.Email)
		log.Infof(ctx, "Registering new user")
		user.ID, err = datastore.Put(ctx, user.ID, &user)
		if err != nil {
			return nil, err
		}
	}
	return &user, nil
}

func gravatarURL(email string) string {
	m := md5.New()
	io.WriteString(m, strings.ToLower(email))
	return fmt.Sprintf("//www.gravatar.com/avatar/%x", m.Sum(nil))
}
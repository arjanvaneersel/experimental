package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path"
)

var ErrNoAvatarURL = errors.New("Chat: Unable to get an avatar URL")

type Avatar interface {
	GetAvatarURL(ChatUser) (string, error)
}

type TryAvatars []Avatar
func (a TryAvatars) GetAvatarURL(c ChatUser) (string, error) {
	for _, avatar := range a {
		if url, err := avatar.GetAvatarURL(c); err == nil {
			return url, nil
		}
	}
	return "", ErrNoAvatarURL
}

type AuthAvatar struct{}
func (AuthAvatar) GetAvatarURL(c ChatUser) (string, error) {
	url := c.AvatarURL()
	if len(url) == 0 {
		return "", ErrNoAvatarURL
	}

	return url, nil
}

var UseAuthAvatar AuthAvatar

type GravatarAvatar struct{}
func (GravatarAvatar) GetAvatarURL(c ChatUser) (string, error) {
	return fmt.Sprintf("//www.gravatar.com/avatar/%s", c.UniqueID()), nil
}

var UseGravatar GravatarAvatar

type FileSystemAvatar struct{}
func (FileSystemAvatar) GetAvatarURL(c ChatUser) (string, error) {
	files, err := ioutil.ReadDir("avatars")
	if err != nil {
		return "", ErrNoAvatarURL
	}

	for _, file := range files {
		if file.IsDir() { continue }
		if match, _ := path.Match(c.UniqueID() + "*", file.Name()); match {
			return "/avatars/" + file.Name(), nil
		}

	}

	return "", ErrNoAvatarURL
}

var UseFileSystemAvatar FileSystemAvatar
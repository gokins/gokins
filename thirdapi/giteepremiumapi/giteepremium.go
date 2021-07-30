package giteepremiumapi

import (
	"github.com/gokins/gokins/thirdapi"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func New(uri string) (*thirdapi.Client, error) {
	base, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}
	if !strings.HasSuffix(base.Path, "/") {
		base.Path = base.Path + "/"
	}
	client := &wrapper{new(thirdapi.Client)}
	client.BaseURL = base
	c := &http.Client{
		Timeout: time.Second * 8,
	}
	client.HttpClient = c
	client.Repositories = &RepositoryService{client}
	return client.Client, nil
}

func NewDefault() *thirdapi.Client {
	client, _ := New(BaseApiGiteePremium)
	return client
}

type wrapper struct {
	*thirdapi.Client
}

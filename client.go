package shotgun_api

import "net/http"

type ShotgunClient interface {
	Do(req *http.Request) (*http.Response, error)
}

var (
	Client ShotgunClient
)

func init() {
	Client = &http.Client{}
}

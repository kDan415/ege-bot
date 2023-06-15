package ege

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

const (
	apiUrl = "https://checkege.rustest.ru/api/exam"
	domain = "checkege.rustest.ru"
	agent  = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36"
)

type Grabber struct {
	url        *url.URL
	httpClient *http.Client
	lastResult *Credentials
}

func New(config Config) (*Grabber, error) {
	var err error
	g := &Grabber{}
	if g.url, err = url.Parse(apiUrl); err != nil {
		return nil, NewError(ConstructorError, 0, err)
	}
	g.httpClient = &http.Client{
		Timeout: time.Second * 10,
	}

	g.UpdateCookie(config.Participant)

	return g, err
}

func (g *Grabber) UpdateCookie(value string) {
	log.Print("Обновление куков")
	var cookies []*http.Cookie
	jar, _ := cookiejar.New(nil)
	cookies = append(cookies, &http.Cookie{
		Name:   "Participant",
		Value:  value,
		Path:   "/",
		Domain: domain,
	})
	jar.SetCookies(g.url, cookies)

	g.httpClient.Jar = jar
}

func (g *Grabber) GetResult() (*Credentials, error) {
	log.Print("запрос к серверу")
	credentials := &Credentials{}

	request, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		return nil, NewError(CreateRequestError, 0, err)
	}
	request.Header.Set("User-Agent", agent)
	request.Close = true

	resp, err := g.httpClient.Do(request)
	if err != nil {
		log.Print("Http client error: ", err)
		return nil, NewError(HttpError, 0, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, NewError(HttpError, resp.StatusCode, nil)
	}

	if err = json.NewDecoder(resp.Body).Decode(credentials); err != nil {
		return nil, NewError(JsonDecodeError, 0, err)
	}
	credentials.Time = time.Now()

	return credentials, nil
}

// GetIfDifferent returns a value only if it differs from the previous request to the EGE server, otherwise nil.
func (g *Grabber) GetIfDifferent() (*Credentials, error) {
	newResult, err := g.GetResult()
	if err != nil {
		return nil, err
	}

	if !compare(g.lastResult, newResult) {
		g.lastResult = newResult
		return newResult, nil
	}
	g.lastResult.Time = time.Now()

	return nil, nil
}

func (g *Grabber) GetLastResult() *Credentials {
	return g.lastResult
}

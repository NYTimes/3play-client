package threeplay

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const ThreePlayHost = "api.3playmedia.com"

type File struct {
	ID                   uint   `json:"id"`
	ProjectID            uint   `json:"project_id"`
	BatchID              uint   `json:"batch_id"`
	Duration             uint   `json:"duration"`
	Attribute1           string `json:"attribute1"`
	Attribute2           string `json:"attribute2"`
	Attribute3           string `json:"attribute3"`
	VideoID              string `json:"video_id"`
	Name                 string `json:"name"`
	CallbackURL          string `json:"callback_url"`
	Description          string `json:"description"`
	CreatedAt            string `json:"created_at"`
	UpdatedAt            string `json:"updated_at"`
	WordCount            uint   `json:"word_count"`
	ThumbnailURL         string `json:"thumbnail_url"`
	LanguageID           int    `json:"language_id"`
	DefaultServiceTypeID int    `json:"default_service_type_id"`
	Downloaded           bool   `json:"downloaded"`
	State                string `json:"state"`
	TurnaroundLevelID    int    `json:"turnaround_level_id"`
	Deadline             string `json:"deadline"`
	BatchName            string `json:"batch_name"`
	ErrorDescription     string `json:"error_description"`
}

type FilesPage struct {
	Files   []File `json:"files"`
	Summary `json:"summary"`
}

type Summary struct {
	CurrentPage  json.Number `json:"current_page"`
	PerPage      json.Number `json:"per_page"`
	TotalEntries json.Number `json:"total_entries"`
	TotalPages   json.Number `json:"total_pages"`
}

type HTTPClient interface {
	Get(string) (*http.Response, error)
	PostForm(string, url.Values) (*http.Response, error)
}

type Client struct {
	apiKey    string
	apiSecret string
	client    HTTPClient
}

type Error struct {
	IsError bool              `json:"iserror"`
	Errors  map[string]string `json:"errors"`
}

func NewClient(apiKey, apiSecret string) *Client {
	return &Client{
		apiKey:    apiKey,
		apiSecret: apiSecret,
		client: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}

func NewClientWithHTTPClient(apiKey, apiSecret string, client HTTPClient) *Client {
	return &Client{
		apiKey:    apiKey,
		apiSecret: apiSecret,
		client:    client,
	}
}

func (c Client) buildUrl(endpoint string, querystring url.Values) string {
	querystring.Add("apikey", c.apiKey)

	url := url.URL{
		Scheme:   "https",
		Host:     ThreePlayHost,
		Path:     endpoint,
		RawQuery: querystring.Encode(),
	}

	return url.String()
}

func (c Client) fetchAndParse(endpoint string, ref interface{}) error {
	apiError := &Error{}
	response, err := c.client.Get(endpoint)

	if err != nil {
		return err
	}

	responseData, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return err
	}

	err = json.Unmarshal(responseData, apiError)

	if err != nil {
		return err
	}

	if apiError.IsError {
		return errors.New("API Error")
	}

	err = json.Unmarshal(responseData, ref)

	if err != nil {
		return err
	}

	return nil
}

func (c Client) GetFiles(params ...url.Values) (*FilesPage, error) {
	querystring := url.Values{}
	if len(params) > 0 {
		querystring = params[0]
	}

	filesPage := &FilesPage{}
	endpoint := c.buildUrl("/files", querystring)
	if err := c.fetchAndParse(endpoint, filesPage); err != nil {
		return nil, err
	} else {
		return filesPage, nil
	}
}

func (c Client) GetFile(id uint) (*File, error) {
	file := &File{}
	endpoint := c.buildUrl(fmt.Sprintf("/files/%d", id), url.Values{})
	if err := c.fetchAndParse(endpoint, file); err != nil {
		return nil, err
	} else {
		return file, nil
	}

}

//UploadFile uploads a file to threeplay using the file's URL.
func (c *Client) UploadFileFromURL(fileURL string, options url.Values) (string, error) {
	endpoint := fmt.Sprintf("https://%s/files", ThreePlayHost)

	data := url.Values{}
	data.Set("apikey", c.apiKey)
	data.Set("api_secret_key", c.apiSecret)
	data.Set("link", fileURL)

	for key, val := range options {
		data[key] = val
	}

	response, err := c.client.PostForm(endpoint, data)
	if err != nil {
		return "", err
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	apiError := &Error{}
	json.Unmarshal(responseData, apiError)
	if apiError.IsError {
		return "", errors.New("API Error")
	}

	return string(responseData), nil
}

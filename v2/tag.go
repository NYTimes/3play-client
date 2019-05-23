package v2

import (
	"errors"
	"fmt"
	"github.com/NYTimes/threeplay/common"
	"net/url"
)

// GetTags gets the list of tags of a file
func (c *ClientV2) GetTags(fileID uint) ([]string, error) {

	endpoint := fmt.Sprintf("https://%s/files/%d/tags?apikey=%s", common.ThreePlayHost, fileID, c.apiKey)
	response, err := c.httpClient.Get(endpoint)
	if err != nil {
		return nil, err
	}

	var tags []string
	if err := parseResponse(response, &tags); err != nil {
		return nil, err
	}

	return tags, nil
}

type addTagResult struct {
	Result bool     `json:"result"`
	Tags   []string `json:"media_file_tags"`
}

// AddTag adds a tag to a file
func (c *ClientV2) AddTag(fileID uint, tag string) ([]string, error) {
	endpoint := fmt.Sprintf("https://%s/files/%d/tags", common.ThreePlayHost, fileID)

	data := url.Values{}
	data.Set("apikey", c.apiKey)
	data.Set("api_secret_key", c.apiSecret)
	data.Set("name", tag)

	response, err := c.httpClient.PostForm(endpoint, data)
	if err != nil {
		return nil, err
	}

	result := &addTagResult{}
	if err := parseResponse(response, result); err != nil {
		return nil, err
	}

	if !result.Result {
		return nil, errors.New("adding Tag Failed")
	}

	return result.Tags, nil
}

// RemoveTag removes a tag of a file
func (c *ClientV2) RemoveTag(fileID uint, tag string) ([]string, error) {
	endpoint := fmt.Sprintf("https://%s/files/%d/tags/%s", common.ThreePlayHost, fileID, tag)

	data := url.Values{}
	data.Set("apikey", c.apiKey)
	data.Set("api_secret_key", c.apiSecret)
	data.Set("_method", "delete")

	response, err := c.httpClient.PostForm(endpoint, data)
	if err != nil {
		return nil, err
	}

	var tags []string
	if err := parseResponse(response, &tags); err != nil {
		return nil, err
	}

	return tags, nil
}
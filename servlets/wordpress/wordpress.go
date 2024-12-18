// wordpress.go
package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	pdk "github.com/extism/go-pdk"
)

// WordPressClient handles WordPress API interactions
type WordPressClient struct {
	baseURL    string
	authHeader string
	isWPCom    bool
	httpClient *http.Client
}

// Enhanced post data structure
type PostData struct {
	Title      string `json:"title"`
	Content    string `json:"content"`
	Status     string `json:"status,omitempty"`
	Categories []int  `json:"categories,omitempty"`
	Tags       []int  `json:"tags,omitempty"`
	Excerpt    string `json:"excerpt,omitempty"`
	Featured   bool   `json:"sticky,omitempty"`
	Date       string `json:"date,omitempty"`
	Format     string `json:"format,omitempty"`
}

// Category/Tag data structures
type TermData struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Parent      int    `json:"parent,omitempty"`
}

type Term struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Helper methods for category/tag resolution
func (c *WordPressClient) resolveTerms(items []interface{}, endpoint string) ([]int, error) {
	var resolvedIDs []int
	var namesToCreate []string

	// First pass: collect IDs and names
	for _, item := range items {
		switch v := item.(type) {
		case float64:
			resolvedIDs = append(resolvedIDs, int(v))
		case string:
			namesToCreate = append(namesToCreate, v)
		}
	}

	if len(namesToCreate) == 0 {
		return resolvedIDs, nil
	}

	// Get existing terms
	resp, err := c.makeRequest("GET", endpoint+"?per_page=100", nil)
	if err != nil {
		return nil, err
	}

	var existingTerms []Term
	if err := json.Unmarshal(resp, &existingTerms); err != nil {
		return nil, fmt.Errorf("failed to parse terms: %v", err)
	}

	// Create lookup map
	termMap := make(map[string]int)
	for _, term := range existingTerms {
		termMap[strings.ToLower(term.Name)] = term.ID
	}

	// Resolve or create terms
	for _, name := range namesToCreate {
		if id, exists := termMap[strings.ToLower(name)]; exists {
			resolvedIDs = append(resolvedIDs, id)
			continue
		}

		// Create new term
		termData := TermData{Name: name}
		body, err := json.Marshal(termData)
		if err != nil {
			return nil, err
		}

		resp, err := c.makeRequest("POST", endpoint, body)
		if err != nil {
			return nil, err
		}

		var newTerm Term
		if err := json.Unmarshal(resp, &newTerm); err != nil {
			return nil, fmt.Errorf("failed to parse new term: %v", err)
		}

		resolvedIDs = append(resolvedIDs, newTerm.ID)
	}

	return resolvedIDs, nil
}

// OAuthResponse represents the OAuth2 token response
type OAuthResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

// NewWordPressClient creates a new WordPress API client
func NewWordPressClient() (*WordPressClient, error) {
	// Get website URL
	url, ok := pdk.GetConfig("website_url")
	if !ok {
		return nil, errors.New("website_url config required")
	}

	// Clean the URL by removing protocol and trailing slashes
	siteDomain := url
	siteDomain = strings.TrimPrefix(siteDomain, "https://")
	siteDomain = strings.TrimPrefix(siteDomain, "http://")
	siteDomain = strings.TrimSuffix(siteDomain, "/")

	// Get on_prem status
	onPrem, ok := pdk.GetConfig("on_prem")
	if !ok {
		return nil, errors.New("on_prem config required")
	}

	isWPCom := onPrem == "no" || onPrem == "false" || onPrem == "0"

	// Initialize client
	client := &WordPressClient{
		httpClient: &http.Client{},
		isWPCom:    isWPCom,
	}

	// Set base URL based on whether it's WordPress.com or self-hosted
	if isWPCom {
		client.baseURL = fmt.Sprintf("https://public-api.wordpress.com/wp/v2/sites/%s", siteDomain)

		// Get OAuth credentials
		clientID, ok := pdk.GetConfig("client_id")
		if !ok {
			return nil, errors.New("client_id config required for WordPress.com sites")
		}
		clientSecret, ok := pdk.GetConfig("client_secret")
		if !ok {
			return nil, errors.New("client_secret config required for WordPress.com sites")
		}
		username, ok := pdk.GetConfig("username")
		if !ok {
			return nil, errors.New("username config required for WordPress.com sites")
		}
		password, ok := pdk.GetConfig("app_password")
		if !ok {
			return nil, errors.New("app_password config required for WordPress.com sites")
		}

		// Get OAuth token
		token, err := client.getOAuthToken(clientID, clientSecret, username, password)
		if err != nil {
			return nil, fmt.Errorf("failed to get OAuth token: %v", err)
		}

		client.authHeader = fmt.Sprintf("Bearer %s", token)
	} else {
		client.baseURL = fmt.Sprintf("%s/wp-json/wp/v2", url)

		// Get basic auth credentials
		username, ok := pdk.GetConfig("username")
		if !ok {
			return nil, errors.New("username config required")
		}
		appPassword, ok := pdk.GetConfig("app_password")
		if !ok {
			return nil, errors.New("app_password config required")
		}

		// Create basic auth header
		auth := fmt.Sprintf("%s:%s", username, appPassword)
		client.authHeader = fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(auth)))
	}

	return client, nil
}

func (c *WordPressClient) makeRequest(method, endpoint string, body []byte) ([]byte, error) {
	url := c.baseURL + endpoint
	pdk.Log(pdk.LogInfo, fmt.Sprintf("Making %s request to %s", method, url))

	// Convert string method to pdk.HTTPMethod
	var pdkMethod pdk.HTTPMethod
	switch method {
	case "GET":
		pdkMethod = pdk.MethodGet
	case "POST":
		pdkMethod = pdk.MethodPost
	case "PUT":
		pdkMethod = pdk.MethodPut
	case "DELETE":
		pdkMethod = pdk.MethodDelete
	default:
		return nil, fmt.Errorf("unsupported HTTP method: %s", method)
	}

	// Create request using PDK
	req := pdk.NewHTTPRequest(pdkMethod, url)
	req.SetHeader("Authorization", c.authHeader)
	req.SetHeader("Content-Type", "application/json")

	if len(body) > 0 {
		req.SetBody(body)
	}

	// Send request
	res := req.Send()

	// Check status code
	if res.Status() < 200 || res.Status() >= 300 {
		return nil, fmt.Errorf("WordPress API error %d: %s", res.Status(), string(res.Body()))
	}

	return res.Body(), nil
}

func (c *WordPressClient) getOAuthToken(clientID, clientSecret, username, password string) (string, error) {
	req := pdk.NewHTTPRequest(pdk.MethodPost, "https://public-api.wordpress.com/oauth2/token")
	req.SetHeader("Content-Type", "application/x-www-form-urlencoded")

	body := fmt.Sprintf(
		"client_id=%s&client_secret=%s&grant_type=password&username=%s&password=%s",
		clientID, clientSecret, username, password,
	)
	req.SetBody([]byte(body))

	res := req.Send()

	if res.Status() != 200 {
		return "", fmt.Errorf("OAuth token request failed with status %d: %s", res.Status(), string(res.Body()))
	}

	var oauthResp OAuthResponse
	if err := json.Unmarshal(res.Body(), &oauthResp); err != nil {
		return "", fmt.Errorf("failed to parse OAuth response: %v", err)
	}

	return oauthResp.AccessToken, nil
}

func (c *WordPressClient) CreatePost(args map[string]interface{}) (CallToolResult, error) {
	postData := PostData{
		Title:   args["title"].(string),
		Content: args["content"].(string),
	}

	if status, ok := args["status"].(string); ok {
		postData.Status = status
	}
	if categories, ok := args["categories"].([]interface{}); ok {
		resolvedCats, err := c.resolveTerms(categories, "/categories")
		if err != nil {
			return CallToolResult{}, fmt.Errorf("failed to resolve categories: %v", err)
		}
		postData.Categories = resolvedCats
	}
	if tags, ok := args["tags"].([]interface{}); ok {
		resolvedTags, err := c.resolveTerms(tags, "/tags")
		if err != nil {
			return CallToolResult{}, fmt.Errorf("failed to resolve tags: %v", err)
		}
		postData.Tags = resolvedTags
	}
	if featured, ok := args["featured"].(bool); ok {
		postData.Featured = featured
	}
	if excerpt, ok := args["excerpt"].(string); ok {
		postData.Excerpt = excerpt
	}

	body, err := json.Marshal(postData)
	if err != nil {
		return CallToolResult{}, err
	}

	resp, err := c.makeRequest("POST", "/posts", body)
	if err != nil {
		return CallToolResult{}, err
	}

	return CallToolResult{
		Content: []Content{{Type: ContentTypeText, Text: ptr(string(resp))}},
	}, nil
}

func (c *WordPressClient) EditPost(args map[string]interface{}) (CallToolResult, error) {
	postID, ok := args["post_id"].(float64)
	if !ok {
		return CallToolResult{}, errors.New("post_id parameter required")
	}

	var postData PostData
	if title, ok := args["title"].(string); ok {
		postData.Title = title
	}
	if content, ok := args["content"].(string); ok {
		postData.Content = content
	}
	if status, ok := args["status"].(string); ok {
		postData.Status = status
	}
	if categories, ok := args["categories"].([]interface{}); ok {
		resolvedCats, err := c.resolveTerms(categories, "/categories")
		if err != nil {
			return CallToolResult{}, fmt.Errorf("failed to resolve categories: %v", err)
		}
		postData.Categories = resolvedCats
	}
	if tags, ok := args["tags"].([]interface{}); ok {
		resolvedTags, err := c.resolveTerms(tags, "/tags")
		if err != nil {
			return CallToolResult{}, fmt.Errorf("failed to resolve tags: %v", err)
		}
		postData.Tags = resolvedTags
	}
	if featured, ok := args["featured"].(bool); ok {
		postData.Featured = featured
	}
	if excerpt, ok := args["excerpt"].(string); ok {
		postData.Excerpt = excerpt
	}

	body, err := json.Marshal(postData)
	if err != nil {
		return CallToolResult{}, err
	}

	resp, err := c.makeRequest("PUT", fmt.Sprintf("/posts/%d", int(postID)), body)
	if err != nil {
		return CallToolResult{}, err
	}

	return CallToolResult{
		Content: []Content{{Type: ContentTypeText, Text: ptr(string(resp))}},
	}, nil
}

func (c *WordPressClient) GetPost(args map[string]interface{}) (CallToolResult, error) {
	postID, ok := args["post_id"].(float64)
	if !ok {
		return CallToolResult{}, errors.New("post_id parameter required")
	}

	resp, err := c.makeRequest("GET", fmt.Sprintf("/posts/%d", int(postID)), nil)
	if err != nil {
		return CallToolResult{}, err
	}

	return CallToolResult{
		Content: []Content{{Type: ContentTypeText, Text: ptr(string(resp))}},
	}, nil
}

func (c *WordPressClient) DeletePost(args map[string]interface{}) (CallToolResult, error) {
	postID, ok := args["post_id"].(float64)
	if !ok {
		return CallToolResult{}, errors.New("post_id parameter required")
	}

	force := false
	if forceVal, ok := args["force"].(bool); ok {
		force = forceVal
	}

	endpoint := fmt.Sprintf("/posts/%d", int(postID))
	if force {
		endpoint += "?force=true"
	}

	resp, err := c.makeRequest("DELETE", endpoint, nil)
	if err != nil {
		return CallToolResult{}, err
	}

	return CallToolResult{
		Content: []Content{{Type: ContentTypeText, Text: ptr(string(resp))}},
	}, nil
}

func (c *WordPressClient) ListPosts(args map[string]interface{}) (CallToolResult, error) {
	queryParams := []string{}

	if page, ok := args["page"].(float64); ok {
		queryParams = append(queryParams, fmt.Sprintf("page=%d", int(page)))
	}
	if perPage, ok := args["per_page"].(float64); ok {
		queryParams = append(queryParams, fmt.Sprintf("per_page=%d", int(perPage)))
	}
	if status, ok := args["status"].(string); ok {
		queryParams = append(queryParams, fmt.Sprintf("status=%s", status))
	}
	if category, ok := args["category"].(float64); ok {
		queryParams = append(queryParams, fmt.Sprintf("categories=%d", int(category)))
	}
	if tag, ok := args["tag"].(float64); ok {
		queryParams = append(queryParams, fmt.Sprintf("tags=%d", int(tag)))
	}

	endpoint := "/posts"
	if len(queryParams) > 0 {
		endpoint += "?" + strings.Join(queryParams, "&")
	}

	resp, err := c.makeRequest("GET", endpoint, nil)
	if err != nil {
		return CallToolResult{}, err
	}

	return CallToolResult{
		Content: []Content{{Type: ContentTypeText, Text: ptr(string(resp))}},
	}, nil
}

func (c *WordPressClient) CreateCategory(args map[string]interface{}) (CallToolResult, error) {
	termData := TermData{
		Name: args["name"].(string),
	}

	if desc, ok := args["description"].(string); ok {
		termData.Description = desc
	}
	if parent, ok := args["parent"].(float64); ok {
		termData.Parent = int(parent)
	}

	body, err := json.Marshal(termData)
	if err != nil {
		return CallToolResult{}, err
	}

	resp, err := c.makeRequest("POST", "/categories", body)
	if err != nil {
		return CallToolResult{}, err
	}

	return CallToolResult{
		Content: []Content{{Type: ContentTypeText, Text: ptr(string(resp))}},
	}, nil
}

func (c *WordPressClient) CreateTag(args map[string]interface{}) (CallToolResult, error) {
	termData := TermData{
		Name: args["name"].(string),
	}

	if desc, ok := args["description"].(string); ok {
		termData.Description = desc
	}

	body, err := json.Marshal(termData)
	if err != nil {
		return CallToolResult{}, err
	}

	resp, err := c.makeRequest("POST", "/tags", body)
	if err != nil {
		return CallToolResult{}, err
	}

	return CallToolResult{
		Content: []Content{{Type: ContentTypeText, Text: ptr(string(resp))}},
	}, nil
}

func (c *WordPressClient) SchedulePost(args map[string]interface{}) (CallToolResult, error) {
	postID, ok := args["post_id"].(float64)
	if !ok {
		return CallToolResult{}, errors.New("post_id parameter required")
	}

	date, ok := args["date"].(string)
	if !ok {
		return CallToolResult{}, errors.New("date parameter required")
	}

	postData := PostData{
		Status: "future",
		Date:   date,
	}

	body, err := json.Marshal(postData)
	if err != nil {
		return CallToolResult{}, err
	}

	resp, err := c.makeRequest("PUT", fmt.Sprintf("/posts/%d", int(postID)), body)
	if err != nil {
		return CallToolResult{}, err
	}

	return CallToolResult{
		Content: []Content{{Type: ContentTypeText, Text: ptr(string(resp))}},
	}, nil
}

func (c *WordPressClient) ListCategories() (CallToolResult, error) {
	resp, err := c.makeRequest("GET", "/categories", nil)
	if err != nil {
		return CallToolResult{}, err
	}

	return CallToolResult{
		Content: []Content{{Type: ContentTypeText, Text: ptr(string(resp))}},
	}, nil
}

func (c *WordPressClient) ListTags() (CallToolResult, error) {
	resp, err := c.makeRequest("GET", "/tags", nil)
	if err != nil {
		return CallToolResult{}, err
	}

	return CallToolResult{
		Content: []Content{{Type: ContentTypeText, Text: ptr(string(resp))}},
	}, nil
}

func (c *WordPressClient) ListComments(args map[string]interface{}) (CallToolResult, error) {
	endpoint := "/comments"
	if postID, ok := args["post_id"].(float64); ok {
		endpoint = fmt.Sprintf("/comments?post=%d", int(postID))
	}

	resp, err := c.makeRequest("GET", endpoint, nil)
	if err != nil {
		return CallToolResult{}, err
	}

	return CallToolResult{
		Content: []Content{{Type: ContentTypeText, Text: ptr(string(resp))}},
	}, nil
}

func (c *WordPressClient) ApproveComment(args map[string]interface{}) (CallToolResult, error) {
	commentID, ok := args["comment_id"].(float64)
	if !ok {
		return CallToolResult{}, errors.New("comment_id parameter required")
	}

	body, err := json.Marshal(map[string]string{"status": "approved"})
	if err != nil {
		return CallToolResult{}, err
	}

	resp, err := c.makeRequest("POST", fmt.Sprintf("/comments/%d", int(commentID)), body)
	if err != nil {
		return CallToolResult{}, err
	}

	return CallToolResult{
		Content: []Content{{Type: ContentTypeText, Text: ptr(string(resp))}},
	}, nil
}

func (c *WordPressClient) RemoveComment(args map[string]interface{}) (CallToolResult, error) {
	commentID, ok := args["comment_id"].(float64)
	if !ok {
		return CallToolResult{}, errors.New("comment_id parameter required")
	}

	resp, err := c.makeRequest("DELETE", fmt.Sprintf("/comments/%d", int(commentID)), nil)
	if err != nil {
		return CallToolResult{}, err
	}

	return CallToolResult{
		Content: []Content{{Type: ContentTypeText, Text: ptr(string(resp))}},
	}, nil
}

// Helper function for creating string pointers
func ptr(s string) *string {
	return &s
}

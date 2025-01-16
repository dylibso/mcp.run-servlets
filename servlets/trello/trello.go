package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	pdk "github.com/extism/go-pdk"
)

// TrelloClient handles Trello API interactions
type TrelloClient struct {
	apiKey string
	token  string
}

// NewTrelloClient creates a new Trello API client
func NewTrelloClient(apiKey string, token string) *TrelloClient {
	return &TrelloClient{
		apiKey: apiKey,
		token:  token,
	}
}

// GetAuthURL returns a URL for authorizing with Trello
func GetAuthURL(apiKey string) string {
	return fmt.Sprintf("https://trello.com/1/authorize?key=%s&response_type=token&expiration=30days&scope=read,write", apiKey)
}

// makeRequest makes an authenticated request to Trello API
func (c *TrelloClient) makeRequest(method pdk.HTTPMethod, path string, queryParams url.Values, body []byte) pdk.HTTPResponse {
	if queryParams == nil {
		queryParams = url.Values{}
	}
	queryParams.Set("key", c.apiKey)
	queryParams.Set("token", c.token)

	url := fmt.Sprintf("https://api.trello.com/1%s?%s", path, queryParams.Encode())

	req := pdk.NewHTTPRequest(method, url)
	if body != nil {
		req.SetHeader("Content-Type", "application/json")
		req.SetBody(body)
	}

	return req.Send()
}

// GetBoardMembers returns all members of a board
func (c *TrelloClient) GetBoardMembers(boardID string) ([]byte, error) {
	resp := c.makeRequest(pdk.MethodGet, fmt.Sprintf("/boards/%s/members", boardID), nil, nil)
	if resp.Status() != 200 {
		return nil, fmt.Errorf("failed to get board members: %s", resp.Body())
	}
	return resp.Body(), nil
}

// GetCardMembers returns all members assigned to a card
func (c *TrelloClient) GetCardMembers(cardID string) ([]byte, error) {
	resp := c.makeRequest(pdk.MethodGet, fmt.Sprintf("/cards/%s/members", cardID), nil, nil)
	if resp.Status() != 200 {
		return nil, fmt.Errorf("failed to get card members: %s", resp.Body())
	}
	return resp.Body(), nil
}

// AddCardMember assigns a member to a card
func (c *TrelloClient) AddCardMember(cardID string, memberID string) ([]byte, error) {
	params := map[string]interface{}{
		"value": memberID,
	}
	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	resp := c.makeRequest(pdk.MethodPost, fmt.Sprintf("/cards/%s/idMembers", cardID), nil, body)
	if resp.Status() != 200 {
		return nil, fmt.Errorf("failed to add card member: %s", resp.Body())
	}
	return resp.Body(), nil
}

// RemoveCardMember removes a member from a card
func (c *TrelloClient) RemoveCardMember(cardID string, memberID string) ([]byte, error) {
	resp := c.makeRequest(pdk.MethodDelete, fmt.Sprintf("/cards/%s/idMembers/%s", cardID, memberID), nil, nil)
	if resp.Status() != 200 {
		return nil, fmt.Errorf("failed to remove card member: %s", resp.Body())
	}
	return resp.Body(), nil
}

// ListBoards retrieves all boards for the authenticated user
func (c *TrelloClient) ListBoards(filter string) ([]byte, error) {
	params := url.Values{}
	if filter != "" {
		params.Set("filter", filter)
	}

	resp := c.makeRequest(pdk.MethodGet, "/members/me/boards", params, nil)
	if resp.Status() != 200 {
		return nil, fmt.Errorf("failed to list boards: %s", resp.Body())
	}

	return resp.Body(), nil
}

// GetBoard retrieves a board's details
func (c *TrelloClient) GetBoard(id string, fields []string) ([]byte, error) {
	params := url.Values{}
	if len(fields) > 0 {
		params.Set("fields", strings.Join(fields, ","))
	}

	resp := c.makeRequest(pdk.MethodGet, fmt.Sprintf("/boards/%s", id), params, nil)
	if resp.Status() != 200 {
		return nil, fmt.Errorf("failed to get board: %s", resp.Body())
	}

	return resp.Body(), nil
}

// CreateBoard creates a new board
func (c *TrelloClient) CreateBoard(name string, desc string) ([]byte, error) {
	params := map[string]interface{}{
		"name": name,
		"desc": desc,
	}
	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	resp := c.makeRequest(pdk.MethodPost, "/boards", nil, body)
	if resp.Status() != 200 {
		return nil, fmt.Errorf("failed to create board: %s", resp.Body())
	}

	return resp.Body(), nil
}

// CreateList creates a new list on a board
func (c *TrelloClient) CreateList(boardID string, name string, pos string) ([]byte, error) {
	params := map[string]interface{}{
		"name":    name,
		"idBoard": boardID,
		"pos":     pos,
	}
	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	resp := c.makeRequest(pdk.MethodPost, "/lists", nil, body)
	if resp.Status() != 200 {
		return nil, fmt.Errorf("failed to create list: %s", resp.Body())
	}

	return resp.Body(), nil
}

// CreateCard creates a new card in a list
func (c *TrelloClient) CreateCard(listID string, name string, desc string) ([]byte, error) {
	params := map[string]interface{}{
		"name":   name,
		"desc":   desc,
		"idList": listID,
	}
	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	resp := c.makeRequest(pdk.MethodPost, "/cards", nil, body)
	if resp.Status() != 200 {
		return nil, fmt.Errorf("failed to create card: %s", resp.Body())
	}

	return resp.Body(), nil
}

// MoveCard moves a card to a different list and/or position
func (c *TrelloClient) MoveCard(cardID string, listID string, position string) ([]byte, error) {
	params := map[string]interface{}{
		"idList": listID,
	}
	if position != "" {
		params["pos"] = position
	}

	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	resp := c.makeRequest(pdk.MethodPut, fmt.Sprintf("/cards/%s", cardID), nil, body)
	if resp.Status() != 200 {
		return nil, fmt.Errorf("failed to move card: %s", resp.Body())
	}

	return resp.Body(), nil
}

// GetListCards returns cards in a list with pagination
func (c *TrelloClient) GetListCards(listID string, limit int, page int) ([]byte, error) {
	params := url.Values{}

	// Add pagination parameters
	if limit > 0 {
		params.Set("limit", fmt.Sprintf("%d", limit))
	}
	if page > 0 {
		params.Set("page", fmt.Sprintf("%d", page))
	}

	resp := c.makeRequest(pdk.MethodGet, fmt.Sprintf("/lists/%s/cards", listID), params, nil)
	if resp.Status() != 200 {
		return nil, fmt.Errorf("failed to get list cards: %s", resp.Body())
	}
	return resp.Body(), nil
}

// UpdateCard updates a card's properties
func (c *TrelloClient) UpdateCard(cardID string, args map[string]interface{}) ([]byte, error) {
	body, err := json.Marshal(args)
	if err != nil {
		return nil, err
	}

	resp := c.makeRequest(pdk.MethodPut, fmt.Sprintf("/cards/%s", cardID), nil, body)
	if resp.Status() != 200 {
		return nil, fmt.Errorf("failed to update card: %s", resp.Body())
	}
	return resp.Body(), nil
}

// AddComment adds a comment to a card
func (c *TrelloClient) AddComment(cardID string, text string) ([]byte, error) {
	params := url.Values{}
	params.Set("text", text)

	resp := c.makeRequest(pdk.MethodPost, fmt.Sprintf("/cards/%s/actions/comments", cardID), params, nil)
	if resp.Status() != 200 {
		return nil, fmt.Errorf("failed to add comment: %s", resp.Body())
	}
	return resp.Body(), nil
}

// DeleteCard deletes a card
func (c *TrelloClient) DeleteCard(cardID string) error {
	resp := c.makeRequest(pdk.MethodDelete, fmt.Sprintf("/cards/%s", cardID), nil, nil)
	if resp.Status() != 200 {
		return fmt.Errorf("failed to delete card: %s", resp.Body())
	}
	return nil
}

// ArchiveAllCards archives all cards in a list
func (c *TrelloClient) ArchiveAllCards(listID string) ([]byte, error) {
	resp := c.makeRequest(pdk.MethodPost, fmt.Sprintf("/lists/%s/archiveAllCards", listID), nil, nil)
	if resp.Status() != 200 {
		return nil, fmt.Errorf("failed to archive cards: %s", resp.Body())
	}
	return resp.Body(), nil
}

// MoveList moves a list to a different board
func (c *TrelloClient) MoveList(listID string, targetBoardID string) ([]byte, error) {
	params := url.Values{}
	params.Set("value", targetBoardID)

	resp := c.makeRequest(pdk.MethodPut, fmt.Sprintf("/lists/%s/idBoard", listID), params, nil)
	if resp.Status() != 200 {
		return nil, fmt.Errorf("failed to move list: %s", resp.Body())
	}
	return resp.Body(), nil
}

// GetBoardLists returns all lists on a board with optional filtering
func (c *TrelloClient) GetBoardLists(boardID string, filter string) ([]byte, error) {
	params := url.Values{}
	if filter != "" {
		params.Set("filter", filter) // open, closed, all
	}

	resp := c.makeRequest(pdk.MethodGet, fmt.Sprintf("/boards/%s/lists", boardID), params, nil)
	if resp.Status() != 200 {
		return nil, fmt.Errorf("failed to get board lists: %s", resp.Body())
	}
	return resp.Body(), nil
}

// GetBoardCards returns all cards on a board
func (c *TrelloClient) GetBoardCards(boardID string) ([]byte, error) {
	resp := c.makeRequest(pdk.MethodGet, fmt.Sprintf("/boards/%s/cards", boardID), nil, nil)
	if resp.Status() != 200 {
		return nil, fmt.Errorf("failed to get board cards: %s", resp.Body())
	}
	return resp.Body(), nil
}

// Labels
func (c *TrelloClient) GetBoardLabels(boardID string) ([]byte, error) {
	resp := c.makeRequest(pdk.MethodGet, fmt.Sprintf("/boards/%s/labels", boardID), nil, nil)
	if resp.Status() != 200 {
		return nil, fmt.Errorf("failed to get board labels: %s", resp.Body())
	}
	return resp.Body(), nil
}

func (c *TrelloClient) CreateLabel(boardID string, name string, color string) ([]byte, error) {
	params := map[string]interface{}{
		"name":    name,
		"color":   color,
		"idBoard": boardID,
	}
	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	resp := c.makeRequest(pdk.MethodPost, "/labels", nil, body)
	if resp.Status() != 200 {
		return nil, fmt.Errorf("failed to create label: %s", resp.Body())
	}
	return resp.Body(), nil
}

func (c *TrelloClient) DeleteLabel(labelID string) error {
	resp := c.makeRequest(pdk.MethodDelete, fmt.Sprintf("/labels/%s", labelID), nil, nil)
	if resp.Status() != 200 {
		return fmt.Errorf("failed to delete label: %s", resp.Body())
	}
	return nil
}

// Checklists
func (c *TrelloClient) CreateChecklist(cardID string, name string) ([]byte, error) {
	params := map[string]interface{}{
		"name":   name,
		"idCard": cardID,
	}
	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	resp := c.makeRequest(pdk.MethodPost, "/checklists", nil, body)
	if resp.Status() != 200 {
		return nil, fmt.Errorf("failed to create checklist: %s", resp.Body())
	}
	return resp.Body(), nil
}

func (c *TrelloClient) CreateChecklistItem(checklistID string, name string) ([]byte, error) {
	params := url.Values{}
	params.Set("name", name)

	resp := c.makeRequest(pdk.MethodPost, fmt.Sprintf("/checklists/%s/checkItems", checklistID), params, nil)
	if resp.Status() != 200 {
		return nil, fmt.Errorf("failed to create checklist item: %s", resp.Body())
	}
	return resp.Body(), nil
}

// Comments
func (c *TrelloClient) GetCardComments(cardID string) ([]byte, error) {
	params := url.Values{}
	params.Set("filter", "commentCard")

	resp := c.makeRequest(pdk.MethodGet, fmt.Sprintf("/cards/%s/actions", cardID), params, nil)
	if resp.Status() != 200 {
		return nil, fmt.Errorf("failed to get card comments: %s", resp.Body())
	}
	return resp.Body(), nil
}

func (c *TrelloClient) DeleteComment(cardID string, commentID string) error {
	resp := c.makeRequest(pdk.MethodDelete, fmt.Sprintf("/cards/%s/actions/%s/comments", cardID, commentID), nil, nil)
	if resp.Status() != 200 {
		return fmt.Errorf("failed to delete comment: %s", resp.Body())
	}
	return nil
}

func (c *TrelloClient) AddBoardMember(boardID string, email string, fullName string) ([]byte, error) {
	params := map[string]interface{}{
		"email":    email,
		"fullName": fullName,
	}
	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	resp := c.makeRequest(pdk.MethodPut, fmt.Sprintf("/boards/%s/members", boardID), nil, body)
	if resp.Status() != 200 {
		return nil, fmt.Errorf("failed to add board member: %s", resp.Body())
	}
	return resp.Body(), nil
}

func (c *TrelloClient) RemoveBoardMember(boardID string, memberID string) error {
	resp := c.makeRequest(pdk.MethodDelete, fmt.Sprintf("/boards/%s/members/%s", boardID, memberID), nil, nil)
	if resp.Status() != 200 {
		return fmt.Errorf("failed to remove board member: %s", resp.Body())
	}
	return nil
}

func (c *TrelloClient) AddCardComment(cardID string, text string) ([]byte, error) {
	params := url.Values{}
	params.Set("text", text)

	resp := c.makeRequest(pdk.MethodPost, fmt.Sprintf("/cards/%s/actions/comments", cardID), params, nil)
	if resp.Status() != 200 {
		return nil, fmt.Errorf("failed to add comment: %s", resp.Body())
	}
	return resp.Body(), nil
}

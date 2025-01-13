package main

import (
	"fmt"

	pdk "github.com/extism/go-pdk"
)

var apiKey string

// loadCredentials loads the API key from config
func loadCredentials() error {
	if apiKey != "" {
		return nil
	}

	key, ok := pdk.GetConfig("api_key")
	if !ok {
		return fmt.Errorf("api_key config required")
	}
	apiKey = key
	return nil
}

// Call implements the servlet entrypoint
func Call(input CallToolRequest) (CallToolResult, error) {
	if err := loadCredentials(); err != nil {
		return CallToolResult{}, err
	}

	args := input.Params.Arguments.(map[string]interface{})

	// Get token from request config if provided
	token, hasToken := args["token"].(string)
	client := NewTrelloClient(apiKey, token)

	if !hasToken {
		input.Params.Name = "get_auth_url"
	}

	switch input.Params.Name {
	case "get_auth_url":
		authURL := GetAuthURL(apiKey)
		instructionsMsg := fmt.Sprintf(
			"1. Visit this URL in your browser:\n%s\n\n"+
				"2. Click 'Allow' to grant access\n"+
				"3. Copy the token shown on the next page\n"+
				"4. Pass in the token as the 'token' parameter to other Trello tools",
			authURL,
		)

		contents := []Content{{
			Type: ContentTypeText,
			Text: some(instructionsMsg),
		}}

		if !hasToken {
			contents = append(contents, Content{
				Type: ContentTypeText,
				Text: some("Note: You must provide a 'token' parameter to use Trello tools"),
			})
		}

		return CallToolResult{
			Content: contents,
		}, nil

	case "list_boards":
		filter, _ := args["filter"].(string)

		resp, err := client.ListBoards(filter)
		if err != nil {
			return CallToolResult{}, err
		}

		return CallToolResult{
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(string(resp)),
			}},
		}, nil

	case "get_board":
		boardID, _ := args["board_id"].(string)
		fields := []string{}
		if f, ok := args["fields"].([]interface{}); ok {
			for _, field := range f {
				fields = append(fields, field.(string))
			}
		}

		resp, err := client.GetBoard(boardID, fields)
		if err != nil {
			return CallToolResult{}, err
		}

		return CallToolResult{
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(string(resp)),
			}},
		}, nil

	case "create_board":
		name, _ := args["name"].(string)
		desc, _ := args["description"].(string)

		resp, err := client.CreateBoard(name, desc)
		if err != nil {
			return CallToolResult{}, err
		}

		return CallToolResult{
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(string(resp)),
			}},
		}, nil

	case "create_list":
		boardID, _ := args["board_id"].(string)
		name, _ := args["name"].(string)
		pos, _ := args["position"].(string)

		resp, err := client.CreateList(boardID, name, pos)
		if err != nil {
			return CallToolResult{}, err
		}

		return CallToolResult{
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(string(resp)),
			}},
		}, nil

	case "create_card":
		listID, _ := args["list_id"].(string)
		name, _ := args["name"].(string)
		desc, _ := args["description"].(string)

		resp, err := client.CreateCard(listID, name, desc)
		if err != nil {
			return CallToolResult{}, err
		}

		return CallToolResult{
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(string(resp)),
			}},
		}, nil

	case "move_card":
		cardID, _ := args["card_id"].(string)
		listID, _ := args["list_id"].(string)
		position, _ := args["position"].(string)

		resp, err := client.MoveCard(cardID, listID, position)
		if err != nil {
			return CallToolResult{}, err
		}

		return CallToolResult{
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(string(resp)),
			}},
		}, nil

	case "get_board_members":
		boardID, _ := args["board_id"].(string)
		resp, err := client.GetBoardMembers(boardID)
		if err != nil {
			return CallToolResult{}, err
		}
		return CallToolResult{
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(string(resp)),
			}},
		}, nil

	case "get_board_lists":
		boardID, _ := args["board_id"].(string)
		filter, _ := args["filter"].(string)
		resp, err := client.GetBoardLists(boardID, filter)
		if err != nil {
			return CallToolResult{}, err
		}
		return CallToolResult{
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(string(resp)),
			}},
		}, nil

	case "get_card_members":
		cardID, _ := args["card_id"].(string)
		resp, err := client.GetCardMembers(cardID)
		if err != nil {
			return CallToolResult{}, err
		}
		return CallToolResult{
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(string(resp)),
			}},
		}, nil

	case "add_card_member":
		cardID, _ := args["card_id"].(string)
		memberID, _ := args["member_id"].(string)
		resp, err := client.AddCardMember(cardID, memberID)
		if err != nil {
			return CallToolResult{}, err
		}
		return CallToolResult{
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(string(resp)),
			}},
		}, nil

	case "remove_card_member":
		cardID, _ := args["card_id"].(string)
		memberID, _ := args["member_id"].(string)
		resp, err := client.RemoveCardMember(cardID, memberID)
		if err != nil {
			return CallToolResult{}, err
		}
		return CallToolResult{
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(string(resp)),
			}},
		}, nil

	default:
		return CallToolResult{}, fmt.Errorf("unknown tool: %s", input.Params.Name)
	}
}

// Describe implements the servlet description
func Describe() (ListToolsResult, error) {
	return ListToolsResult{
		Tools: []ToolDescription{
			{
				Name:        "get_auth_url",
				Description: "Get a Trello authorization URL that will display an API token. This token is required for other Trello tools.",
				InputSchema: map[string]interface{}{
					"type":       "object",
					"properties": map[string]interface{}{},
				},
			},
			{
				Name:        "get_board",
				Description: "Get details about a Trello board",
				InputSchema: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"token": map[string]interface{}{
							"type":        "string",
							"description": "Trello API token",
						},
						"board_id": map[string]interface{}{
							"type":        "string",
							"description": "The ID of the board to retrieve",
						},
						"fields": map[string]interface{}{
							"type":        "array",
							"description": "Board fields to include in response",
							"items": map[string]interface{}{
								"type": "string",
							},
						},
					},
					"required": []string{"token", "board_id"},
				},
			},
			{
				Name:        "list_boards",
				Description: "List all boards for the authenticated user",
				InputSchema: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"token": map[string]interface{}{
							"type":        "string",
							"description": "Trello API token",
						},
						"filter": map[string]interface{}{
							"type":        "string",
							"description": "Filter boards. Valid values: all, open, closed, members, organization, public, starred (default: all)",
							"enum":        []string{"all", "open", "closed", "members", "organization", "public", "starred"},
						},
					},
					"required": []string{"token"},
				},
			},
			{
				Name:        "create_board",
				Description: "Create a new Trello board",
				InputSchema: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"token": map[string]interface{}{
							"type":        "string",
							"description": "Trello API token",
						},
						"name": map[string]interface{}{
							"type":        "string",
							"description": "Name of the board",
						},
						"description": map[string]interface{}{
							"type":        "string",
							"description": "Description of the board",
						},
					},
					"required": []string{"token", "name"},
				},
			},
			{
				Name:        "create_list",
				Description: "Create a new list on a board",
				InputSchema: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"token": map[string]interface{}{
							"type":        "string",
							"description": "Trello API token",
						},
						"board_id": map[string]interface{}{
							"type":        "string",
							"description": "ID of the board to add list to",
						},
						"name": map[string]interface{}{
							"type":        "string",
							"description": "Name of the list",
						},
						"position": map[string]interface{}{
							"type":        "string",
							"description": "Position of list (top, bottom, or a positive number)",
						},
					},
					"required": []string{"token", "board_id", "name"},
				},
			},
			{
				Name:        "create_card",
				Description: "Create a new card in a list",
				InputSchema: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"token": map[string]interface{}{
							"type":        "string",
							"description": "Trello API token",
						},
						"list_id": map[string]interface{}{
							"type":        "string",
							"description": "ID of the list to add card to",
						},
						"name": map[string]interface{}{
							"type":        "string",
							"description": "Name of the card",
						},
						"description": map[string]interface{}{
							"type":        "string",
							"description": "Description of the card",
						},
					},
					"required": []string{"token", "list_id", "name"},
				},
			},
			{
				Name:        "move_card",
				Description: "Move a card to a different list and/or position",
				InputSchema: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"token": map[string]interface{}{
							"type":        "string",
							"description": "Trello API token",
						},
						"card_id": map[string]interface{}{
							"type":        "string",
							"description": "ID of the card to move",
						},
						"list_id": map[string]interface{}{
							"type":        "string",
							"description": "ID of the destination list",
						},
						"position": map[string]interface{}{
							"type":        "string",
							"description": "Optional - Position in the list (top, bottom, or a positive number)",
						},
					},
					"required": []string{"token", "card_id", "list_id"},
				},
			},
			{
				Name:        "get_board_members",
				Description: "Get all members of a board",
				InputSchema: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"token": map[string]interface{}{
							"type":        "string",
							"description": "Trello API token",
						},
						"board_id": map[string]interface{}{
							"type":        "string",
							"description": "ID of the board",
						},
					},
					"required": []string{"token", "board_id"},
				},
			},
			{
				Name:        "get_board_lists",
				Description: "Get all lists on a board",
				InputSchema: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"token": map[string]interface{}{
							"type":        "string",
							"description": "Trello API token",
						},
						"board_id": map[string]interface{}{
							"type":        "string",
							"description": "ID of the board",
						},
						"filter": map[string]interface{}{
							"type":        "string",
							"description": "Filter lists: open, closed, or all",
							"enum":        []string{"open", "closed", "all"},
						},
					},
					"required": []string{"token", "board_id"},
				},
			},
			{
				Name:        "get_card_members",
				Description: "Get all members assigned to a card",
				InputSchema: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"token": map[string]interface{}{
							"type":        "string",
							"description": "Trello API token",
						},
						"card_id": map[string]interface{}{
							"type":        "string",
							"description": "ID of the card",
						},
					},
					"required": []string{"token", "card_id"},
				},
			},
			{
				Name:        "add_card_member",
				Description: "Assign a member to a card",
				InputSchema: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"token": map[string]interface{}{
							"type":        "string",
							"description": "Trello API token",
						},
						"card_id": map[string]interface{}{
							"type":        "string",
							"description": "ID of the card",
						},
						"member_id": map[string]interface{}{
							"type":        "string",
							"description": "ID of the member to assign",
						},
					},
					"required": []string{"token", "card_id", "member_id"},
				},
			},
			{
				Name:        "remove_card_member",
				Description: "Remove a member from a card",
				InputSchema: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"token": map[string]interface{}{
							"type":        "string",
							"description": "Trello API token",
						},
						"card_id": map[string]interface{}{
							"type":        "string",
							"description": "ID of the card",
						},
						"member_id": map[string]interface{}{
							"type":        "string",
							"description": "ID of the member to remove",
						},
					},
					"required": []string{"token", "card_id", "member_id"},
				},
			},
		},
	}, nil
}

func some[T any](t T) *T {
	return &t
}

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

	if !hasToken && input.Params.Name != "auth_get_url" {
		input.Params.Name = "auth_get_url"
	}

	switch input.Params.Name {
	// Auth operations
	case "auth_get_url":
		authURL := GetAuthURL(apiKey)
		instructionsMsg :=
			"1. Visit this URL in your browser:\n" +
				"2. Click 'Allow' to grant access\n" +
				"3. Copy the token shown on the next page\n" +
				"4. Pass in the token as the 'token' parameter to other Trello tools"

		contents := []Content{
			{
				Type: ContentTypeText,
				Text: &authURL,
			},
			{
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

	// Board operations
	case "board_list":
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

	case "board_get":
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

	case "board_create":
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

	case "board_get_members":
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

	case "board_add_member":
		boardID, _ := args["board_id"].(string)
		email, _ := args["email"].(string)
		fullName, _ := args["full_name"].(string)
		resp, err := client.AddBoardMember(boardID, email, fullName)
		if err != nil {
			return CallToolResult{}, err
		}
		return CallToolResult{
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(string(resp)),
			}},
		}, nil

	case "board_remove_member":
		boardID, _ := args["board_id"].(string)
		memberID, _ := args["member_id"].(string)
		err := client.RemoveBoardMember(boardID, memberID)
		if err != nil {
			return CallToolResult{}, err
		}
		return CallToolResult{
			Content: []Content{{
				Type: ContentTypeText,
				Text: some("Member removed successfully"),
			}},
		}, nil

	case "board_get_labels":
		boardID, _ := args["board_id"].(string)
		resp, err := client.GetBoardLabels(boardID)
		if err != nil {
			return CallToolResult{}, err
		}
		return CallToolResult{
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(string(resp)),
			}},
		}, nil

	case "board_get_lists":
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

	case "board_get_cards":
		boardID, _ := args["board_id"].(string)
		resp, err := client.GetBoardCards(boardID)
		if err != nil {
			return CallToolResult{}, err
		}
		return CallToolResult{
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(string(resp)),
			}},
		}, nil

	// List operations
	case "list_create":
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

	case "list_move":
		listID, _ := args["list_id"].(string)
		boardID, _ := args["board_id"].(string)
		resp, err := client.MoveList(listID, boardID)
		if err != nil {
			return CallToolResult{}, err
		}
		return CallToolResult{
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(string(resp)),
			}},
		}, nil

	case "list_get_cards":
		listID, _ := args["list_id"].(string)
		limit := 50 // default value
		if l, ok := args["limit"].(float64); ok {
			limit = int(l)
		}
		page := 0 // default value
		if p, ok := args["page"].(float64); ok {
			page = int(p)
		}
		resp, err := client.GetListCards(listID, limit, page)
		if err != nil {
			return CallToolResult{}, err
		}
		return CallToolResult{
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(string(resp)),
			}},
		}, nil

	case "list_archive_cards":
		listID, _ := args["list_id"].(string)
		resp, err := client.ArchiveAllCards(listID)
		if err != nil {
			return CallToolResult{}, err
		}
		return CallToolResult{
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(string(resp)),
			}},
		}, nil

	// Card operations
	case "card_create":
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

	case "card_move":
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

	case "card_get_members":
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

	case "card_add_member":
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

	case "card_remove_member":
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

	case "card_get_comments":
		cardID, _ := args["card_id"].(string)
		resp, err := client.GetCardComments(cardID)
		if err != nil {
			return CallToolResult{}, err
		}
		return CallToolResult{
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(string(resp)),
			}},
		}, nil

	case "card_add_comment":
		cardID, _ := args["card_id"].(string)
		text, _ := args["text"].(string)
		resp, err := client.AddCardComment(cardID, text)
		if err != nil {
			return CallToolResult{}, err
		}
		return CallToolResult{
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(string(resp)),
			}},
		}, nil

	// Label operations
	case "label_create":
		boardID, _ := args["board_id"].(string)
		name, _ := args["name"].(string)
		color, _ := args["color"].(string)
		resp, err := client.CreateLabel(boardID, name, color)
		if err != nil {
			return CallToolResult{}, err
		}
		return CallToolResult{
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(string(resp)),
			}},
		}, nil

	case "label_delete":
		labelID, _ := args["label_id"].(string)
		err := client.DeleteLabel(labelID)
		if err != nil {
			return CallToolResult{}, err
		}
		return CallToolResult{
			Content: []Content{{
				Type: ContentTypeText,
				Text: some("Label deleted successfully"),
			}},
		}, nil

	// Checklist operations
	case "checklist_create":
		cardID, _ := args["card_id"].(string)
		name, _ := args["name"].(string)
		resp, err := client.CreateChecklist(cardID, name)
		if err != nil {
			return CallToolResult{}, err
		}
		return CallToolResult{
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(string(resp)),
			}},
		}, nil

	case "checklist_add_item":
		checklistID, _ := args["checklist_id"].(string)
		name, _ := args["name"].(string)
		resp, err := client.CreateChecklistItem(checklistID, name)
		if err != nil {
			return CallToolResult{}, err
		}
		return CallToolResult{
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(string(resp)),
			}},
		}, nil

	// Comment operations
	case "comment_delete":
		cardID, _ := args["card_id"].(string)
		commentID, _ := args["comment_id"].(string)
		err := client.DeleteComment(cardID, commentID)
		if err != nil {
			return CallToolResult{}, err
		}
		return CallToolResult{
			Content: []Content{{
				Type: ContentTypeText,
				Text: some("Comment deleted successfully"),
			}},
		}, nil

	default:
		return CallToolResult{}, fmt.Errorf("unknown tool: %s", input.Params.Name)
	}
}

func Describe() (ListToolsResult, error) {
	return ListToolsResult{
		Tools: []ToolDescription{
			// Auth tools
			{
				Name:        "auth_get_url",
				Description: "Get a Trello authorization URL that will display an API token. This token is required for other Trello tools.",
				InputSchema: map[string]interface{}{
					"type":       "object",
					"properties": map[string]interface{}{},
				},
			},
			// Board tools
			{
				Name:        "board_list",
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
				Name:        "board_get",
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
				Name:        "board_create",
				Description: "Create a new Trello board. Note: when creating a board, some lists are automatically created. Check the existing lists before creating new ones.",
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
				Name:        "board_get_members",
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
				Name:        "board_add_member",
				Description: "Add a member to a board by email",
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
						"email": map[string]interface{}{
							"type":        "string",
							"description": "Email of the user to add",
						},
						"full_name": map[string]interface{}{
							"type":        "string",
							"description": "Full name of the user",
						},
					},
					"required": []string{"token", "board_id", "email", "full_name"},
				},
			},
			{
				Name:        "board_remove_member",
				Description: "Remove a member from a board",
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
						"member_id": map[string]interface{}{
							"type":        "string",
							"description": "ID of the member to remove",
						},
					},
					"required": []string{"token", "board_id", "member_id"},
				},
			},
			{
				Name:        "board_get_labels",
				Description: "Get all labels on a board",
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
				Name:        "board_get_lists",
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
				Name:        "board_get_cards",
				Description: "Get all cards on a board",
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

			// List tools
			{
				Name:        "list_create",
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
				Name:        "list_move",
				Description: "Move a list to a different board",
				InputSchema: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"token": map[string]interface{}{
							"type":        "string",
							"description": "Trello API token",
						},
						"list_id": map[string]interface{}{
							"type":        "string",
							"description": "ID of the list to move",
						},
						"board_id": map[string]interface{}{
							"type":        "string",
							"description": "ID of the target board",
						},
					},
					"required": []string{"token", "list_id", "board_id"},
				},
			},
			{
				Name:        "list_get_cards",
				Description: "Get all cards in a list with pagination",
				InputSchema: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"token": map[string]interface{}{
							"type":        "string",
							"description": "Trello API token",
						},
						"list_id": map[string]interface{}{
							"type":        "string",
							"description": "ID of the list",
						},
						"limit": map[string]interface{}{
							"type":        "integer",
							"description": "Maximum number of cards to return per page (default: 50, max: 1000)",
							"minimum":     1,
							"maximum":     1000,
							"default":     50,
						},
						"page": map[string]interface{}{
							"type":        "integer",
							"description": "Page number to return (0-based)",
							"minimum":     0,
							"default":     0,
						},
					},
					"required": []string{"token", "list_id"},
				},
			},
			{
				Name:        "list_archive_cards",
				Description: "Archive all cards in a list",
				InputSchema: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"token": map[string]interface{}{
							"type":        "string",
							"description": "Trello API token",
						},
						"list_id": map[string]interface{}{
							"type":        "string",
							"description": "ID of the list",
						},
					},
					"required": []string{"token", "list_id"},
				},
			},

			// Card tools
			{
				Name:        "card_create",
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
				Name:        "card_move",
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
				Name:        "card_get_members",
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
				Name:        "card_add_member",
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
				Name:        "card_remove_member",
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
			{
				Name:        "card_get_comments",
				Description: "Get all comments on a card",
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
				Name:        "card_add_comment",
				Description: "Add a comment to a card",
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
						"text": map[string]interface{}{
							"type":        "string",
							"description": "Text of the comment",
						},
					},
					"required": []string{"token", "card_id", "text"},
				},
			},

			// Label tools
			{
				Name:        "label_create",
				Description: "Create a new label on a board",
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
						"name": map[string]interface{}{
							"type":        "string",
							"description": "Name of the label",
						},
						"color": map[string]interface{}{
							"type":        "string",
							"description": "Color of the label (red, yellow, green, blue, purple, orange, black)",
							"enum":        []string{"red", "yellow", "green", "blue", "purple", "orange", "black"},
						},
					},
					"required": []string{"token", "board_id", "name", "color"},
				},
			},
			{
				Name:        "label_delete",
				Description: "Delete a label",
				InputSchema: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"token": map[string]interface{}{
							"type":        "string",
							"description": "Trello API token",
						},
						"label_id": map[string]interface{}{
							"type":        "string",
							"description": "ID of the label to delete",
						},
					},
					"required": []string{"token", "label_id"},
				},
			},

			// Checklist tools
			{
				Name:        "checklist_create",
				Description: "Create a new checklist on a card",
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
						"name": map[string]interface{}{
							"type":        "string",
							"description": "Name of the checklist",
						},
					},
					"required": []string{"token", "card_id", "name"},
				},
			},
			{
				Name:        "checklist_add_item",
				Description: "Add an item to a checklist",
				InputSchema: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"token": map[string]interface{}{
							"type":        "string",
							"description": "Trello API token",
						},
						"checklist_id": map[string]interface{}{
							"type":        "string",
							"description": "ID of the checklist",
						},
						"name": map[string]interface{}{
							"type":        "string",
							"description": "Name of the checklist item",
						},
					},
					"required": []string{"token", "checklist_id", "name"},
				},
			},

			// Comment tools

			{
				Name:        "comment_delete",
				Description: "Delete a comment from a card",
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
						"comment_id": map[string]interface{}{
							"type":        "string",
							"description": "ID of the comment to delete",
						},
					},
					"required": []string{"token", "card_id", "comment_id"},
				},
			},
		},
	}, nil
}

func some[T any](t T) *T {
	return &t
}

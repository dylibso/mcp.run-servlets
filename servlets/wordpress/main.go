// main.go
package main

import (
	"errors"
	"fmt"
)

// Call handles all tool requests
func Call(input CallToolRequest) (CallToolResult, error) {
	if input.Params.Arguments == nil {
		return CallToolResult{}, errors.New("arguments must be provided")
	}

	args, ok := input.Params.Arguments.(map[string]interface{})
	if !ok {
		return CallToolResult{}, errors.New("invalid arguments format")
	}

	client, err := NewWordPressClient()
	if err != nil {
		return CallToolResult{}, fmt.Errorf("failed to create WordPress client: %v", err)
	}

	switch input.Params.Name {
	case "wp_post_create":
		return client.CreatePost(args)
	case "wp_post_edit":
		return client.EditPost(args)
	case "wp_post_list":
		return client.ListPosts(args)
	case "wp_post_get":
		return client.GetPost(args)
	case "wp_post_delete":
		return client.DeletePost(args)
	case "wp_post_schedule":
		return client.SchedulePost(args)
	case "wp_category_list":
		return client.ListCategories()
	case "wp_category_create":
		return client.CreateCategory(args)
	case "wp_tag_list":
		return client.ListTags()
	case "wp_tag_create":
		return client.CreateTag(args)
	case "wp_comment_list":
		return client.ListComments(args)
	case "wp_comment_approve":
		return client.ApproveComment(args)
	case "wp_comment_delete":
		return client.RemoveComment(args)
	default:
		return CallToolResult{}, fmt.Errorf("unknown tool: %s", input.Params.Name)
	}
}

// Describe implements the tool description
func Describe() (ListToolsResult, error) {
	return ListToolsResult{
		Tools: []ToolDescription{
			{
				Name:        "wp_post_create",
				Description: "Create a new WordPress post",
				InputSchema: map[string]interface{}{
					"type":     "object",
					"required": []string{"title", "content"},
					"properties": map[string]interface{}{
						"title": map[string]interface{}{
							"type":        "string",
							"description": "Post title",
						},
						"content": map[string]interface{}{
							"type":        "string",
							"description": "Post content",
						},
						"status": map[string]interface{}{
							"type":        "string",
							"description": "Post status (draft, publish, private, etc.)",
							"default":     "draft",
						},
						"categories": map[string]interface{}{
							"type":        "array",
							"description": "Array of category IDs",
							"items": map[string]interface{}{
								"type": "number",
							},
						},
						"tags": map[string]interface{}{
							"type":        "array",
							"description": "Array of tag IDs",
							"items": map[string]interface{}{
								"type": "number",
							},
						},
						"featured": map[string]interface{}{
							"type":        "boolean",
							"description": "Whether to mark the post as sticky",
							"default":     false,
						},
						"excerpt": map[string]interface{}{
							"type":        "string",
							"description": "Post excerpt/summary",
						},
					},
				},
			},
			{
				Name:        "wp_post_edit",
				Description: "Edit an existing WordPress post",
				InputSchema: map[string]interface{}{
					"type":     "object",
					"required": []string{"post_id"},
					"properties": map[string]interface{}{
						"post_id": map[string]interface{}{
							"type":        "number",
							"description": "ID of the post to edit",
						},
						"title": map[string]interface{}{
							"type":        "string",
							"description": "New post title",
						},
						"content": map[string]interface{}{
							"type":        "string",
							"description": "New post content",
						},
						"status": map[string]interface{}{
							"type":        "string",
							"description": "New post status",
						},
						"categories": map[string]interface{}{
							"type":        "array",
							"description": "Array of category IDs to set",
							"items": map[string]interface{}{
								"type": "number",
							},
						},
						"tags": map[string]interface{}{
							"type":        "array",
							"description": "Array of tag IDs to set",
							"items": map[string]interface{}{
								"type": "number",
							},
						},
						"featured": map[string]interface{}{
							"type":        "boolean",
							"description": "Whether to mark the post as sticky",
						},
						"excerpt": map[string]interface{}{
							"type":        "string",
							"description": "New post excerpt/summary",
						},
					},
				},
			},
			{
				Name:        "wp_post_get",
				Description: "Get a single WordPress post by ID",
				InputSchema: map[string]interface{}{
					"type":     "object",
					"required": []string{"post_id"},
					"properties": map[string]interface{}{
						"post_id": map[string]interface{}{
							"type":        "number",
							"description": "ID of the post to retrieve",
						},
					},
				},
			},
			{
				Name:        "wp_post_delete",
				Description: "Delete a WordPress post",
				InputSchema: map[string]interface{}{
					"type":     "object",
					"required": []string{"post_id"},
					"properties": map[string]interface{}{
						"post_id": map[string]interface{}{
							"type":        "number",
							"description": "ID of the post to delete",
						},
						"force": map[string]interface{}{
							"type":        "boolean",
							"description": "Whether to bypass trash and force deletion",
							"default":     false,
						},
					},
				},
			},
			{
				Name:        "wp_post_list",
				Description: "List WordPress posts with pagination and filtering",
				InputSchema: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"page": map[string]interface{}{
							"type":        "number",
							"description": "Page number",
							"default":     1,
						},
						"per_page": map[string]interface{}{
							"type":        "number",
							"description": "Posts per page",
							"default":     10,
						},
						"status": map[string]interface{}{
							"type":        "string",
							"description": "Filter by post status",
							"enum":        []string{"publish", "draft", "private", "pending", "future"},
						},
						"category": map[string]interface{}{
							"type":        "number",
							"description": "Filter by category ID",
						},
						"tag": map[string]interface{}{
							"type":        "number",
							"description": "Filter by tag ID",
						},
					},
				},
			},
			{
				Name:        "wp_post_schedule",
				Description: "Schedule a post for future publication",
				InputSchema: map[string]interface{}{
					"type":     "object",
					"required": []string{"post_id", "date"},
					"properties": map[string]interface{}{
						"post_id": map[string]interface{}{
							"type":        "number",
							"description": "ID of the post to schedule",
						},
						"date": map[string]interface{}{
							"type":        "string",
							"description": "ISO 8601 formatted date (e.g. 2024-12-31T15:30:00)",
						},
					},
				},
			},
			{
				Name:        "wp_category_create",
				Description: "Create a new WordPress category",
				InputSchema: map[string]interface{}{
					"type":     "object",
					"required": []string{"name"},
					"properties": map[string]interface{}{
						"name": map[string]interface{}{
							"type":        "string",
							"description": "Category name",
						},
						"description": map[string]interface{}{
							"type":        "string",
							"description": "Category description",
						},
						"parent": map[string]interface{}{
							"type":        "number",
							"description": "ID of parent category",
						},
					},
				},
			},
			{
				Name:        "wp_category_list",
				Description: "List all WordPress categories",
				InputSchema: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"dummy": map[string]interface{}{
							"type":        "string",
							"description": "Dummy parameter to make the schema non-empty. Ignored.",
						},
					},
				},
			},
			{
				Name:        "wp_tag_create",
				Description: "Create a new WordPress tag",
				InputSchema: map[string]interface{}{
					"type":     "object",
					"required": []string{"name"},
					"properties": map[string]interface{}{
						"name": map[string]interface{}{
							"type":        "string",
							"description": "Tag name",
						},
						"description": map[string]interface{}{
							"type":        "string",
							"description": "Tag description",
						},
					},
				},
			},
			{
				Name:        "wp_tag_list",
				Description: "List all WordPress tags",
				InputSchema: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"dummy": map[string]interface{}{
							"type":        "string",
							"description": "Dummy parameter to make the schema non-empty. Ignored.",
						},
					},
				},
			},
			{
				Name:        "wp_comment_list",
				Description: "List WordPress comments, optionally filtered by post",
				InputSchema: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"post_id": map[string]interface{}{
							"type":        "number",
							"description": "Optional post ID to filter comments",
						},
					},
				},
			},
			{
				Name:        "wp_comment_approve",
				Description: "Approve a WordPress comment",
				InputSchema: map[string]interface{}{
					"type":     "object",
					"required": []string{"comment_id"},
					"properties": map[string]interface{}{
						"comment_id": map[string]interface{}{
							"type":        "number",
							"description": "ID of the comment to approve",
						},
					},
				},
			},
			{
				Name:        "wp_comment_delete",
				Description: "Remove a WordPress comment",
				InputSchema: map[string]interface{}{
					"type":     "object",
					"required": []string{"comment_id"},
					"properties": map[string]interface{}{
						"comment_id": map[string]interface{}{
							"type":        "number",
							"description": "ID of the comment to remove",
						},
					},
				},
			},
		},
	}, nil
}

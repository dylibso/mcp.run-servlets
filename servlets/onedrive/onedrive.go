package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/extism/go-pdk"
)

// Helper function to get the base URL for drive operations
func getDriveBaseURL(driveID string, path string) string {
	if driveID == "" {
		// Use default drive
		return fmt.Sprintf("https://graph.microsoft.com/v1.0/me/drive/%s", path)
	} else {
		// Use specified drive
		return fmt.Sprintf("https://graph.microsoft.com/v1.0/drives/%s/%s", driveID, path)
	}
}

// Helper function to add optional query parameters
func addPaginationParams(params *[]string, args map[string]interface{}) {
	// Handle top parameter (with a maximum of 999)
	if top, ok := args["top"].(float64); ok && top > 0 {
		if top > 999 {
			top = 999 // Maximum value
		}
		*params = append(*params, fmt.Sprintf("$top=%d", int(top)))
	}

	// Handle orderby parameter
	if orderby, ok := args["orderby"].(string); ok && orderby != "" {
		*params = append(*params, fmt.Sprintf("$orderby=%s", orderby))
	}
}

// List items in the root of the drive
func listDriveItems(token string, args map[string]interface{}) (CallToolResult, error) {
	driveID, _ := args["drive_id"].(string)
	baseURL := getDriveBaseURL(driveID, "root/children")
	params := []string{}

	// Add pagination and sorting parameters
	addPaginationParams(&params, args)

	// Build final URL
	requestURL := baseURL
	if len(params) > 0 {
		requestURL = fmt.Sprintf("%s?%s", baseURL, strings.Join(params, "&"))
	}

	req := pdk.NewHTTPRequest(pdk.MethodGet, requestURL)
	req.SetHeader("Authorization", fmt.Sprintf("Bearer %s", token))
	req.SetHeader("Accept", "application/json")

	resp := req.Send()
	if resp.Status() != 200 {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(fmt.Sprintf("Failed to list drive items: %d %s", resp.Status(), string(resp.Body()))),
			}},
		}, nil
	}

	return CallToolResult{
		Content: []Content{{
			Type: ContentTypeText,
			Text: some(string(resp.Body())),
		}},
	}, nil
}

// Get recently accessed files
func getRecentFiles(token string, args map[string]interface{}) (CallToolResult, error) {
	driveID, _ := args["drive_id"].(string)
	baseURL := getDriveBaseURL(driveID, "recent")
	params := []string{}

	// Add pagination parameter
	if top, ok := args["top"].(float64); ok && top > 0 {
		if top > 999 {
			top = 999 // Maximum value
		}
		params = append(params, fmt.Sprintf("$top=%d", int(top)))
	}

	// Build final URL
	requestURL := baseURL
	if len(params) > 0 {
		requestURL = fmt.Sprintf("%s?%s", baseURL, strings.Join(params, "&"))
	}

	req := pdk.NewHTTPRequest(pdk.MethodGet, requestURL)
	req.SetHeader("Authorization", fmt.Sprintf("Bearer %s", token))
	req.SetHeader("Accept", "application/json")

	resp := req.Send()
	if resp.Status() != 200 {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(fmt.Sprintf("Failed to get recent files: %d %s", resp.Status(), string(resp.Body()))),
			}},
		}, nil
	}

	return CallToolResult{
		Content: []Content{{
			Type: ContentTypeText,
			Text: some(string(resp.Body())),
		}},
	}, nil
}

// List files shared with the user
func getSharedWithMe(token string, args map[string]interface{}) (CallToolResult, error) {
	baseURL := "https://graph.microsoft.com/v1.0/me/drive/sharedWithMe"
	params := []string{}

	// Add pagination parameter
	if top, ok := args["top"].(float64); ok && top > 0 {
		if top > 999 {
			top = 999 // Maximum value
		}
		params = append(params, fmt.Sprintf("$top=%d", int(top)))
	}

	// Build final URL
	requestURL := baseURL
	if len(params) > 0 {
		requestURL = fmt.Sprintf("%s?%s", baseURL, strings.Join(params, "&"))
	}

	req := pdk.NewHTTPRequest(pdk.MethodGet, requestURL)
	req.SetHeader("Authorization", fmt.Sprintf("Bearer %s", token))
	req.SetHeader("Accept", "application/json")

	resp := req.Send()
	if resp.Status() != 200 {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(fmt.Sprintf("Failed to get shared files: %d %s", resp.Status(), string(resp.Body()))),
			}},
		}, nil
	}

	return CallToolResult{
		Content: []Content{{
			Type: ContentTypeText,
			Text: some(string(resp.Body())),
		}},
	}, nil
}

// Search OneDrive for files and folders
func searchDrive(token string, args map[string]interface{}) (CallToolResult, error) {
	// Get search query
	query, ok := args["query"].(string)
	if !ok || query == "" {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some("Search query is required"),
			}},
		}, nil
	}

	driveID, _ := args["drive_id"].(string)
	var baseURL string

	if driveID == "" {
		// Use proper search function format for default drive
		baseURL = fmt.Sprintf("https://graph.microsoft.com/v1.0/me/drive/search(q='%s')", url.QueryEscape(query))
	} else {
		// Use proper search function format for specified drive
		baseURL = fmt.Sprintf("https://graph.microsoft.com/v1.0/drives/%s/search(q='%s')", driveID, url.QueryEscape(query))
	}

	params := []string{}

	// Add pagination parameter
	if top, ok := args["top"].(float64); ok && top > 0 {
		if top > 999 {
			top = 999 // Maximum value
		}
		params = append(params, fmt.Sprintf("$top=%d", int(top)))
	}

	// Build final URL
	requestURL := baseURL
	if len(params) > 0 {
		requestURL = fmt.Sprintf("%s?%s", baseURL, strings.Join(params, "&"))
	}

	req := pdk.NewHTTPRequest(pdk.MethodGet, requestURL)
	req.SetHeader("Authorization", fmt.Sprintf("Bearer %s", token))
	req.SetHeader("Accept", "application/json")

	resp := req.Send()
	if resp.Status() != 200 {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(fmt.Sprintf("Failed to search drive: %d %s", resp.Status(), string(resp.Body()))),
			}},
		}, nil
	}

	return CallToolResult{
		Content: []Content{{
			Type: ContentTypeText,
			Text: some(string(resp.Body())),
		}},
	}, nil
}

// Create a new folder in OneDrive
func createFolder(token string, args map[string]interface{}) (CallToolResult, error) {
	// Get folder name
	name, ok := args["name"].(string)
	if !ok || name == "" {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some("Folder name is required"),
			}},
		}, nil
	}

	driveID, _ := args["drive_id"].(string)

	// Determine parent folder
	var baseURL string
	parent, hasParent := args["parent"].(string)

	if hasParent && parent != "" {
		baseURL = getDriveBaseURL(driveID, fmt.Sprintf("items/%s/children", parent))
	} else {
		baseURL = getDriveBaseURL(driveID, "root/children")
	}

	// Create folder payload
	folderPayload := map[string]interface{}{
		"name":                              name,
		"folder":                            map[string]interface{}{},
		"@microsoft.graph.conflictBehavior": "rename",
	}

	jsonData, err := json.Marshal(folderPayload)
	if err != nil {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(fmt.Sprintf("Failed to marshal folder payload: %s", err)),
			}},
		}, nil
	}

	req := pdk.NewHTTPRequest(pdk.MethodPost, baseURL)
	req.SetHeader("Authorization", fmt.Sprintf("Bearer %s", token))
	req.SetHeader("Content-Type", "application/json")
	req.SetHeader("Accept", "application/json")
	req.SetBody(jsonData)

	resp := req.Send()
	if resp.Status() != 201 {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(fmt.Sprintf("Failed to create folder: %d %s", resp.Status(), string(resp.Body()))),
			}},
		}, nil
	}

	return CallToolResult{
		Content: []Content{{
			Type: ContentTypeText,
			Text: some(string(resp.Body())),
		}},
	}, nil
}

// List all drives
func listDrives(token string) (CallToolResult, error) {
	requestURL := "https://graph.microsoft.com/v1.0/me/drives"

	req := pdk.NewHTTPRequest(pdk.MethodGet, requestURL)
	req.SetHeader("Authorization", fmt.Sprintf("Bearer %s", token))
	req.SetHeader("Accept", "application/json")

	resp := req.Send()
	if resp.Status() != 200 {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(fmt.Sprintf("Failed to list drives: %d %s", resp.Status(), string(resp.Body()))),
			}},
		}, nil
	}

	return CallToolResult{
		Content: []Content{{
			Type: ContentTypeText,
			Text: some(string(resp.Body())),
		}},
	}, nil
}

// List children of a specific folder
func listFolderChildren(token string, args map[string]interface{}) (CallToolResult, error) {
	// Get folder ID
	folderId, ok := args["folder_id"].(string)
	if !ok || folderId == "" {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some("Folder ID is required"),
			}},
		}, nil
	}

	driveID, _ := args["drive_id"].(string)
	baseURL := getDriveBaseURL(driveID, fmt.Sprintf("items/%s/children", folderId))
	params := []string{}

	// Add pagination and sorting parameters
	addPaginationParams(&params, args)

	// Build final URL
	requestURL := baseURL
	if len(params) > 0 {
		requestURL = fmt.Sprintf("%s?%s", baseURL, strings.Join(params, "&"))
	}

	req := pdk.NewHTTPRequest(pdk.MethodGet, requestURL)
	req.SetHeader("Authorization", fmt.Sprintf("Bearer %s", token))
	req.SetHeader("Accept", "application/json")

	resp := req.Send()
	if resp.Status() != 200 {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(fmt.Sprintf("Failed to list folder children: %d %s", resp.Status(), string(resp.Body()))),
			}},
		}, nil
	}

	return CallToolResult{
		Content: []Content{{
			Type: ContentTypeText,
			Text: some(string(resp.Body())),
		}},
	}, nil
}

// Upload file (for files smaller than 4MB)
func uploadFile(token string, args map[string]interface{}) (CallToolResult, error) {
	// Get file details
	name, ok := args["name"].(string)
	if !ok || name == "" {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some("File name is required"),
			}},
		}, nil
	}

	content, ok := args["content"].(string)
	if !ok {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some("File content is required"),
			}},
		}, nil
	}

	contentType := "application/octet-stream" // Default content type
	if ct, ok := args["content_type"].(string); ok && ct != "" {
		contentType = ct
	}

	driveID, _ := args["drive_id"].(string)

	// Determine parent folder
	var baseURL string
	parentId, hasParent := args["parent_id"].(string)

	if hasParent && parentId != "" {
		if driveID == "" {
			baseURL = fmt.Sprintf("https://graph.microsoft.com/v1.0/me/drive/items/%s:/%s:/content", parentId, url.PathEscape(name))
		} else {
			baseURL = fmt.Sprintf("https://graph.microsoft.com/v1.0/drives/%s/items/%s:/%s:/content", driveID, parentId, url.PathEscape(name))
		}
	} else {
		if driveID == "" {
			baseURL = fmt.Sprintf("https://graph.microsoft.com/v1.0/me/drive/root:/%s:/content", url.PathEscape(name))
		} else {
			baseURL = fmt.Sprintf("https://graph.microsoft.com/v1.0/drives/%s/root:/%s:/content", driveID, url.PathEscape(name))
		}
	}

	req := pdk.NewHTTPRequest(pdk.MethodPut, baseURL)
	req.SetHeader("Authorization", fmt.Sprintf("Bearer %s", token))
	req.SetHeader("Content-Type", contentType)
	req.SetBody([]byte(content))

	resp := req.Send()
	if resp.Status() != 201 && resp.Status() != 200 {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(fmt.Sprintf("Failed to upload file: %d %s", resp.Status(), string(resp.Body()))),
			}},
		}, nil
	}

	return CallToolResult{
		Content: []Content{{
			Type: ContentTypeText,
			Text: some(string(resp.Body())),
		}},
	}, nil
}

// Get information about a specific item (file or folder)
func getItem(token string, args map[string]interface{}) (CallToolResult, error) {
	// Get item ID
	itemID, ok := args["item_id"].(string)
	if !ok || itemID == "" {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some("Item ID is required"),
			}},
		}, nil
	}

	driveID, _ := args["drive_id"].(string)
	baseURL := getDriveBaseURL(driveID, fmt.Sprintf("items/%s", itemID))

	req := pdk.NewHTTPRequest(pdk.MethodGet, baseURL)
	req.SetHeader("Authorization", fmt.Sprintf("Bearer %s", token))
	req.SetHeader("Accept", "application/json")

	resp := req.Send()
	if resp.Status() != 200 {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(fmt.Sprintf("Failed to get item: %d %s", resp.Status(), string(resp.Body()))),
			}},
		}, nil
	}

	return CallToolResult{
		Content: []Content{{
			Type: ContentTypeText,
			Text: some(string(resp.Body())),
		}},
	}, nil
}

// Get information about a specific drive or the default drive
func getDriveInfo(token string, args map[string]interface{}) (CallToolResult, error) {
	var baseURL string
	driveID, hasDriveID := args["drive_id"].(string)

	if hasDriveID && driveID != "" {
		baseURL = fmt.Sprintf("https://graph.microsoft.com/v1.0/drives/%s", driveID)
	} else {
		baseURL = "https://graph.microsoft.com/v1.0/me/drive"
	}

	req := pdk.NewHTTPRequest(pdk.MethodGet, baseURL)
	req.SetHeader("Authorization", fmt.Sprintf("Bearer %s", token))
	req.SetHeader("Accept", "application/json")

	resp := req.Send()
	if resp.Status() != 200 {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(fmt.Sprintf("Failed to get drive info: %d %s", resp.Status(), string(resp.Body()))),
			}},
		}, nil
	}

	return CallToolResult{
		Content: []Content{{
			Type: ContentTypeText,
			Text: some(string(resp.Body())),
		}},
	}, nil
}

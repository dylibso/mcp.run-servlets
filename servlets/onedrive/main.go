package main

import (
	"fmt"

	"github.com/extism/go-pdk"
)

// Tool definitions
var (
	ListDriveItemsTool = ToolDescription{
		Name:        "list-drive-items",
		Description: "List items in the root of a drive",
		InputSchema: schema{
			"type": "object",
			"properties": props{
				"drive_id": prop("string", "The ID of the drive (omit to use default drive)"),
				"top":      prop("integer", "Number of items to return (max 999)"),
				"orderby":  prop("string", "Property to sort by (e.g., name, lastModifiedDateTime)"),
			},
		},
	}

	ListRecentFilesTool = ToolDescription{
		Name:        "recent-files",
		Description: "Get a list of recently accessed files",
		InputSchema: schema{
			"type": "object",
			"properties": props{
				"drive_id": prop("string", "The ID of the drive (omit to use default drive)"),
				"top":      prop("integer", "Number of items to return (max 999)"),
			},
		},
	}

	ListSharedWithMeTool = ToolDescription{
		Name:        "shared-with-me",
		Description: "List files and folders that have been shared with you",
		InputSchema: schema{
			"type": "object",
			"properties": props{
				"top": prop("integer", "Number of items to return (max 999)"),
			},
		},
	}

	SearchDriveTool = ToolDescription{
		Name:        "search",
		Description: "Search for files and folders in a drive",
		InputSchema: schema{
			"type": "object",
			"properties": props{
				"query":    prop("string", "Search query string"),
				"drive_id": prop("string", "The ID of the drive (omit to use default drive)"),
				"top":      prop("integer", "Number of items to return (max 999)"),
			},
			"required": []string{"query"},
		},
	}

	CreateFolderTool = ToolDescription{
		Name:        "create-folder",
		Description: "Create a new folder in a drive",
		InputSchema: schema{
			"type": "object",
			"properties": props{
				"name":     prop("string", "The name of the folder to create"),
				"drive_id": prop("string", "The ID of the drive (omit to use default drive)"),
				"parent":   prop("string", "The parent folder ID (optional, defaults to root)"),
			},
			"required": []string{"name"},
		},
	}

	ListChildrenOfDriveTool = ToolDescription{
		Name:        "list-drive-children",
		Description: "List all drives available to the user",
		InputSchema: schema{
			"type":       "object",
			"properties": props{},
		},
	}

	ListFolderChildrenTool = ToolDescription{
		Name:        "list-folder-children",
		Description: "List children items of a specific folder",
		InputSchema: schema{
			"type": "object",
			"properties": props{
				"folder_id": prop("string", "The ID of the folder to list contents"),
				"drive_id":  prop("string", "The ID of the drive (omit to use default drive)"),
				"top":       prop("integer", "Number of items to return (max 999)"),
				"orderby":   prop("string", "Property to sort by (e.g., name, lastModifiedDateTime)"),
			},
			"required": []string{"folder_id"},
		},
	}

	UploadFileTool = ToolDescription{
		Name:        "upload-file",
		Description: "Upload a small file to OneDrive (less than 4MB)",
		InputSchema: schema{
			"type": "object",
			"properties": props{
				"name":         prop("string", "The name of the file to upload"),
				"content":      prop("string", "The content of the file as a string"),
				"drive_id":     prop("string", "The ID of the drive (omit to use default drive)"),
				"parent_id":    prop("string", "The parent folder ID (optional, defaults to root)"),
				"content_type": prop("string", "The content type of the file (e.g., text/plain)"),
			},
			"required": []string{"name", "content"},
		},
	}

	GetItemTool = ToolDescription{
		Name:        "get-item",
		Description: "Get information about a specific item (file or folder) by ID",
		InputSchema: schema{
			"type": "object",
			"properties": props{
				"item_id":  prop("string", "The ID of the item to get information about"),
				"drive_id": prop("string", "The ID of the drive (omit to use default drive)"),
			},
			"required": []string{"item_id"},
		},
	}

	GetDriveInfoTool = ToolDescription{
		Name:        "get-drive-info",
		Description: "Get information about a specific drive by ID or default drive",
		InputSchema: schema{
			"type": "object",
			"properties": props{
				"drive_id": prop("string", "The ID of the drive (omit to get info about default drive)"),
			},
		},
	}

	OneDriveTools = []ToolDescription{
		ListDriveItemsTool,
		ListRecentFilesTool,
		ListSharedWithMeTool,
		SearchDriveTool,
		CreateFolderTool,
		ListChildrenOfDriveTool,
		ListFolderChildrenTool,
		UploadFileTool,
		GetItemTool,
		GetDriveInfoTool,
	}
)

// Called when the tool is invoked
func Call(input CallToolRequest) (CallToolResult, error) {
	token, ok := pdk.GetConfig("OAUTH_TOKEN")
	if !ok {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some("No OAUTH_TOKEN configured"),
			}},
		}, nil
	}

	args := input.Params.Arguments.(map[string]interface{})
	pdk.Log(pdk.LogDebug, fmt.Sprint("Args: ", args))

	switch input.Params.Name {
	case "list-drive-items":
		return listDriveItems(token, args)
	case "recent-files":
		return getRecentFiles(token, args)
	case "shared-with-me":
		return getSharedWithMe(token, args)
	case "search":
		return searchDrive(token, args)
	case "create-folder":
		return createFolder(token, args)
	case "list-drive-children":
		return listDrives(token)
	case "list-folder-children":
		return listFolderChildren(token, args)
	case "upload-file":
		return uploadFile(token, args)
	case "get-item":
		return getItem(token, args)
	case "get-drive-info":
		return getDriveInfo(token, args)
	default:
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some("Unknown tool " + input.Params.Name),
			}},
		}, nil
	}
}

// Describe the tools provided by this servlet
func Describe() (ListToolsResult, error) {
	return ListToolsResult{
		Tools: OneDriveTools,
	}, nil
}

// Helper function to create a pointer
func some[T any](t T) *T {
	return &t
}

// Schema related types for tool description
type SchemaProperty struct {
	Type        string  `json:"type"`
	Description string  `json:"description,omitempty"`
	Items       *schema `json:"items,omitempty"`
}

func prop(tpe, description string) SchemaProperty {
	return SchemaProperty{Type: tpe, Description: description}
}

type schema = map[string]interface{}
type props = map[string]SchemaProperty

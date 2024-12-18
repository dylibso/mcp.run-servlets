mod pdk;

use extism_pdk::*;
use pdk::{types::ToolDescription, *};
use serde_json::Value;
use std::collections::BTreeMap;

const NOTION_API_BASE: &str = "https://api.notion.com/v1";
const COMMON_ID_DESCRIPTION: &str = "It should be a 32-character string (excluding hyphens) formatted as 8-4-4-4-12 with hyphens (-).";

fn make_notion_request(
    endpoint: &str,
    method: &str,
    body: Option<Value>,
    token: &str,
) -> Result<Value, Error> {
    let url = format!("{}{}", NOTION_API_BASE, endpoint);

    let mut request = HttpRequest::new(&url);
    request.method = Some(method.to_string());

    let mut headers = BTreeMap::new();
    headers.insert("Authorization".to_string(), format!("Bearer {}", token));
    headers.insert("Notion-Version".to_string(), "2022-06-28".to_string());
    headers.insert("Content-Type".to_string(), "application/json".to_string());
    request.headers = headers;

    let body_bytes = body.map(|b| serde_json::to_vec(&b)).transpose()?;

    let response = http::request(&request, body_bytes.as_deref())?;

    if response.status_code() < 200 || response.status_code() >= 300 {
        return Err(Error::msg(format!(
            "Notion API error: {} - {}",
            response.status_code(),
            String::from_utf8_lossy(&response.body())
        )));
    }

    serde_json::from_slice(&response.body()).map_err(Error::from)
}

pub(crate) fn call(input: types::CallToolRequest) -> Result<types::CallToolResult, Error> {
    let args = input
        .params
        .arguments
        .as_ref()
        .ok_or_else(|| Error::msg("No arguments provided"))?;

    let token = config::get("NOTION_TOKEN")?
        .ok_or_else(|| Error::msg("NOTION_TOKEN not found in config"))?;

    // Log the request (following TypeScript pattern)
    eprintln!("Executing Notion tool: {}", input.params.name);

    let result = match input.params.name.as_str() {
        // Block operations
        "notion_append_block_children" => {
            let block_id = args
                .get("block_id")
                .and_then(|v| v.as_str())
                .ok_or_else(|| Error::msg("block_id is required and must be a string"))?;
            let children = args
                .get("children")
                .ok_or_else(|| Error::msg("children array is required"))?;

            // Validate that children is an array
            if !children.is_array() {
                return Err(Error::msg("children must be an array of block objects"));
            }

            make_notion_request(
                &format!("/blocks/{}/children", block_id),
                "PATCH",
                Some(serde_json::json!({ "children": children })),
                &token,
            )?
        }

        "notion_retrieve_block" => {
            let block_id = args
                .get("block_id")
                .and_then(|v| v.as_str())
                .ok_or_else(|| Error::msg("block_id is required and must be a string"))?;

            make_notion_request(&format!("/blocks/{}", block_id), "GET", None, &token)?
        }

        "notion_retrieve_block_children" => {
            let block_id = args
                .get("block_id")
                .and_then(|v| v.as_str())
                .ok_or_else(|| Error::msg("block_id is required and must be a string"))?;

            let mut params = Vec::new();
            if let Some(start_cursor) = args.get("start_cursor").and_then(|v| v.as_str()) {
                params.push(("start_cursor", start_cursor.to_string()));
            }
            if let Some(page_size) = args.get("page_size").and_then(|v| v.as_u64()) {
                params.push(("page_size", page_size.to_string()));
            }

            let query_string = if !params.is_empty() {
                format!(
                    "?{}",
                    params
                        .iter()
                        .map(|(k, v)| format!("{}={}", k, v))
                        .collect::<Vec<_>>()
                        .join("&")
                )
            } else {
                String::new()
            };

            make_notion_request(
                &format!("/blocks/{}/children{}", block_id, query_string),
                "GET",
                None,
                &token,
            )?
        }

        "notion_delete_block" => {
            let block_id = args
                .get("block_id")
                .and_then(|v| v.as_str())
                .ok_or_else(|| Error::msg("block_id is required and must be a string"))?;

            make_notion_request(&format!("/blocks/{}", block_id), "DELETE", None, &token)?
        }

        // Page operations
        "notion_retrieve_page" => {
            let page_id = args
                .get("page_id")
                .and_then(|v| v.as_str())
                .ok_or_else(|| Error::msg("page_id is required and must be a string"))?;

            make_notion_request(&format!("/pages/{}", page_id), "GET", None, &token)?
        }

        "notion_update_page_properties" => {
            let page_id = args
                .get("page_id")
                .and_then(|v| v.as_str())
                .ok_or_else(|| Error::msg("page_id is required and must be a string"))?;
            let properties = args
                .get("properties")
                .ok_or_else(|| Error::msg("properties object is required"))?;

            make_notion_request(
                &format!("/pages/{}", page_id),
                "PATCH",
                Some(serde_json::json!({ "properties": properties })),
                &token,
            )?
        }

        // Database operations
        "notion_create_database" => {
            let parent = args
                .get("parent")
                .ok_or_else(|| Error::msg("parent object is required"))?;
            let title = args
                .get("title")
                .ok_or_else(|| Error::msg("title array is required"))?;
            let properties = args
                .get("properties")
                .ok_or_else(|| Error::msg("properties object is required"))?;

            make_notion_request(
                "/databases",
                "POST",
                Some(serde_json::json!({
                    "parent": parent,
                    "title": title,
                    "properties": properties
                })),
                &token,
            )?
        }

        "notion_query_database" => {
            let database_id = args
                .get("database_id")
                .and_then(|v| v.as_str())
                .ok_or_else(|| Error::msg("database_id is required and must be a string"))?;

            let mut query_body = serde_json::Map::new();
            if let Some(filter) = args.get("filter") {
                query_body.insert("filter".to_string(), filter.clone());
            }
            if let Some(sorts) = args.get("sorts") {
                query_body.insert("sorts".to_string(), sorts.clone());
            }
            if let Some(start_cursor) = args.get("start_cursor") {
                query_body.insert("start_cursor".to_string(), start_cursor.clone());
            }
            if let Some(page_size) = args.get("page_size") {
                query_body.insert("page_size".to_string(), page_size.clone());
            }

            make_notion_request(
                &format!("/databases/{}/query", database_id),
                "POST",
                Some(Value::Object(query_body)),
                &token,
            )?
        }

        // Comments operations
        "notion_create_comment" => {
            let mut body = serde_json::Map::new();

            if let Some(parent) = args.get("parent") {
                body.insert("parent".to_string(), parent.clone());
            }
            if let Some(discussion_id) = args.get("discussion_id") {
                body.insert("discussion_id".to_string(), discussion_id.clone());
            }
            if let Some(rich_text) = args.get("rich_text") {
                body.insert("rich_text".to_string(), rich_text.clone());
            } else {
                return Err(Error::msg("rich_text array is required"));
            }

            if !body.contains_key("parent") && !body.contains_key("discussion_id") {
                return Err(Error::msg(
                    "Either parent.page_id or discussion_id must be provided",
                ));
            }

            make_notion_request("/comments", "POST", Some(Value::Object(body)), &token)?
        }

        "notion_retrieve_comments" => {
            let block_id = args
                .get("block_id")
                .and_then(|v| v.as_str())
                .ok_or_else(|| Error::msg("block_id is required and must be a string"))?;

            let mut params = vec![("block_id", block_id.to_string())];
            if let Some(start_cursor) = args.get("start_cursor").and_then(|v| v.as_str()) {
                params.push(("start_cursor", start_cursor.to_string()));
            }
            if let Some(page_size) = args.get("page_size").and_then(|v| v.as_u64()) {
                params.push(("page_size", page_size.to_string()));
            }

            let query_string = params
                .iter()
                .map(|(k, v)| format!("{}={}", k, v))
                .collect::<Vec<_>>()
                .join("&");

            make_notion_request(&format!("/comments?{}", query_string), "GET", None, &token)?
        }

        // Search operations
        "notion_search" => {
            let mut search_body = serde_json::Map::new();

            if let Some(query) = args.get("query") {
                search_body.insert("query".to_string(), query.clone());
            }
            if let Some(filter) = args.get("filter") {
                search_body.insert("filter".to_string(), filter.clone());
            }
            if let Some(sort) = args.get("sort") {
                search_body.insert("sort".to_string(), sort.clone());
            }
            if let Some(start_cursor) = args.get("start_cursor") {
                search_body.insert("start_cursor".to_string(), start_cursor.clone());
            }
            if let Some(page_size) = args.get("page_size") {
                search_body.insert("page_size".to_string(), page_size.clone());
            }

            make_notion_request("/search", "POST", Some(Value::Object(search_body)), &token)?
        }

        // Error for unknown tool
        _ => return Err(Error::msg(format!("Unknown tool: {}", input.params.name))),
    };

    Ok(types::CallToolResult {
        content: vec![types::Content {
            r#type: types::ContentType::Text,
            text: Some(serde_json::to_string_pretty(&result)?),
            annotations: None,
            data: None,
            mime_type: Some("application/json".to_string()),
        }],
        is_error: None,
    })
}

pub(crate) fn describe() -> Result<types::ListToolsResult, Error> {
    // Define rich text object schema
    let rich_text_schema = serde_json::json!({
        "type": "object",
        "description": "A rich text object.",
        "properties": {
            "type": {
                "type": "string",
                "description": "The type of this rich text object.",
                "enum": ["text", "mention", "equation"]
            },
            "text": {
                "type": "object",
                "description": "Object containing text content and optional link info. Required if type is 'text'.",
                "properties": {
                    "content": {
                        "type": "string",
                        "description": "The actual text content."
                    },
                    "link": {
                        "type": "object",
                        "description": "Optional link object with a 'url' field.",
                        "properties": {
                            "url": {
                                "type": "string",
                                "description": "The URL the text links to."
                            }
                        }
                    }
                }
            },
            "annotations": {
                "type": "object",
                "description": "Styling information for the text.",
                "properties": {
                    "bold": { "type": "boolean" },
                    "italic": { "type": "boolean" },
                    "strikethrough": { "type": "boolean" },
                    "underline": { "type": "boolean" },
                    "code": { "type": "boolean" },
                    "color": {
                        "type": "string",
                        "description": "Color for the text.",
                        "enum": [
                            "default", "blue", "blue_background", "brown", "brown_background",
                            "gray", "gray_background", "green", "green_background",
                            "orange", "orange_background", "pink", "pink_background",
                            "purple", "purple_background", "red", "red_background",
                            "yellow", "yellow_background"
                        ]
                    }
                }
            }
        },
        "required": ["type"]
    });

    // Define block object schema
    let block_schema = serde_json::json!({
        "type": "object",
        "description": "A Notion block object.",
        "properties": {
            "object": {
                "type": "string",
                "description": "Should be 'block'.",
                "enum": ["block"]
            },
            "type": {
                "type": "string",
                "description": "Type of the block.",
                "enum": [
                    "paragraph", "heading_1", "heading_2", "heading_3",
                    "bulleted_list_item", "numbered_list_item", "to_do",
                    "toggle", "child_page", "child_database", "embed",
                    "callout", "quote", "equation", "divider",
                    "table_of_contents", "column", "column_list",
                    "link_preview", "synced_block", "template",
                    "link_to_page", "audio", "bookmark", "breadcrumb",
                    "code", "file", "image", "pdf", "video"
                ]
            }
        },
        "required": ["object", "type"]
    });

    Ok(types::ListToolsResult {
        tools: vec![
            // Block operations
            ToolDescription {
                name: "notion_append_block_children".to_string(),
                description: "Append new children blocks to a specified parent block in Notion. Requires insert content capabilities.".to_string(),
                input_schema: serde_json::json!({
                    "type": "object",
                    "properties": {
                        "block_id": {
                            "type": "string",
                            "description": "The ID of the parent block. ".to_string() + COMMON_ID_DESCRIPTION
                        },
                        "children": {
                            "type": "array",
                            "description": "Array of block objects to append. Each block must follow the Notion block schema.",
                            "items": block_schema
                        }
                    },
                    "required": ["block_id", "children"],
                    "additionalProperties": false
                }).as_object().unwrap().clone(),
            },
            ToolDescription {
                name: "notion_retrieve_block".to_string(),
                description: "Retrieve a block from Notion".to_string(),
                input_schema: serde_json::json!({
                    "type": "object",
                    "properties": {
                        "block_id": {
                            "type": "string",
                            "description": "The ID of the block to retrieve. ".to_string() + COMMON_ID_DESCRIPTION
                        }
                    },
                    "required": ["block_id"],
                    "additionalProperties": false
                }).as_object().unwrap().clone(),
            },
            // Pages
            ToolDescription {
                name: "notion_retrieve_page".to_string(),
                description: "Retrieve a page from Notion".to_string(),
                input_schema: serde_json::json!({
                    "type": "object",
                    "properties": {
                        "page_id": {
                            "type": "string",
                            "description": "The ID of the page to retrieve. ".to_string() + COMMON_ID_DESCRIPTION
                        }
                    },
                    "required": ["page_id"],
                    "additionalProperties": false
                }).as_object().unwrap().clone(),
            },
            ToolDescription {
                name: "notion_update_page_properties".to_string(),
                description: "Update properties of a page or an item in a Notion database".to_string(),
                input_schema: serde_json::json!({
                    "type": "object",
                    "properties": {
                        "page_id": {
                            "type": "string",
                            "description": "The ID of the page or database item to update. ".to_string() + COMMON_ID_DESCRIPTION
                        },
                        "properties": {
                            "type": "object",
                            "description": "Properties to update. These correspond to the columns or fields in the database."
                        }
                    },
                    "required": ["page_id", "properties"],
                    "additionalProperties": false
                }).as_object().unwrap().clone(),
            },
            // Databases
            ToolDescription {
                name: "notion_create_database".to_string(),
                description: "Create a database in Notion".to_string(),
                input_schema: serde_json::json!({
                    "type": "object",
                    "properties": {
                        "parent": {
                            "type": "object",
                            "description": "Parent object of the database"
                        },
                        "title": {
                            "type": "array",
                            "description": "Title of database as it appears in Notion. An array of rich text objects.",
                            "items": rich_text_schema
                        },
                        "properties": {
                            "type": "object",
                            "description": "Property schema of database. The keys are the names of properties as they appear in Notion and the values are property schema objects."
                        }
                    },
                    "required": ["parent", "title", "properties"],
                    "additionalProperties": false
                }).as_object().unwrap().clone(),
            },
            ToolDescription {
                name: "notion_query_database".to_string(),
                description: "Query a database in Notion".to_string(),
                input_schema: serde_json::json!({
                    "type": "object",
                    "properties": {
                        "database_id": {
                            "type": "string",
                            "description": "The ID of the database to query. ".to_string() + COMMON_ID_DESCRIPTION
                        },
                        "filter": {
                            "type": "object",
                            "description": "Filter conditions"
                        },
                        "sorts": {
                            "type": "array",
                            "description": "Sort conditions"
                        },
                        "start_cursor": {
                            "type": "string",
                            "description": "Pagination cursor for next page of results"
                        },
                        "page_size": {
                            "type": "integer",
                            "description": "Number of results per page (max 100)"
                        }
                    },
                    "required": ["database_id"],
                    "additionalProperties": false
                }).as_object().unwrap().clone(),
            },
            // Users
            ToolDescription {
                name: "notion_list_users".to_string(),
                description: "List all users in the Notion workspace. **Note:** This function requires upgrading to the Notion Enterprise plan and using an Organization API key to avoid permission errors.".to_string(),
                input_schema: serde_json::json!({
                    "type": "object",
                    "properties": {
                        "start_cursor": {
                            "type": "string",
                            "description": "Pagination start cursor for listing users"
                        },
                        "page_size": {
                            "type": "integer",
                            "description": "Number of users to retrieve (max 100)"
                        }
                    },
                    "additionalProperties": false
                }).as_object().unwrap().clone(),
            },
            // Comments
            ToolDescription {
                name: "notion_create_comment".to_string(),
                description: "Create a comment in Notion. This requires the integration to have 'insert comment' capabilities. You can either specify a page parent or a discussion_id, but not both.".to_string(),
                input_schema: serde_json::json!({
                    "type": "object",
                    "properties": {
                        "parent": {
                            "type": "object",
                            "description": "Parent object that specifies the page to comment on",
                            "properties": {
                                "page_id": {
                                    "type": "string",
                                    "description": "The ID of the page to comment on. ".to_string() + COMMON_ID_DESCRIPTION
                                }
                            }
                        },
                        "discussion_id": {
                            "type": "string",
                            "description": "The ID of an existing discussion thread to add a comment to. ".to_string() + COMMON_ID_DESCRIPTION
                        },
                        "rich_text": {
                            "type": "array",
                            "description": "Array of rich text objects representing the comment content.",
                            "items": rich_text_schema
                        }
                    },
                    "required": ["rich_text"],
                    "additionalProperties": false
                }).as_object().unwrap().clone(),
            },
            // Search
            ToolDescription {
                name: "notion_search".to_string(),
                description: "Search pages or databases by title in Notion".to_string(),
                input_schema: serde_json::json!({
                    "type": "object",
                    "properties": {
                        "query": {
                            "type": "string",
                            "description": "Text to search for in page or database titles"
                        },
                        "filter": {
                            "type": "object",
                            "description": "Filter results by object type (page or database)",
                            "properties": {
                                "property": {
                                    "type": "string",
                                    "description": "Must be 'object'"
                                },
                                "value": {
                                    "type": "string",
                                    "description": "Either 'page' or 'database'"
                                }
                            }
                        },
                        "sort": {
                            "type": "object",
                            "description": "Sort configuration for search results",
                            "properties": {
                                "direction": {
                                    "type": "string",
                                    "enum": ["ascending", "descending"]
                                },
                                "timestamp": {
                                    "type": "string",
                                    "enum": ["last_edited_time"]
                                }
                            }
                        },
                        "start_cursor": {
                            "type": "string",
                            "description": "Pagination cursor"
                        },
                        "page_size": {
                            "type": "integer",
                            "description": "Number of results to return (max 100)"
                        }
                    },
                    "additionalProperties": false
                }).as_object().unwrap().clone(),
            },
        ],
    })
}

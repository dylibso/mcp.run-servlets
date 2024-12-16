mod pdk;

use extism_pdk::*;
use pdk::{types::ToolDescription, *};
use serde_json::Value;
use std::collections::BTreeMap;

const NOTION_API_BASE: &str = "https://api.notion.com/v1";

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

    let operation = args
        .get("operation")
        .and_then(|op| op.as_str())
        .ok_or_else(|| Error::msg("operation is required"))?;

    let token = config::get("NOTION_TOKEN")?
        .ok_or_else(|| Error::msg("NOTION_TOKEN not found in config"))?;

    let result = match operation {
        // Block operations
        "append_block_children" => {
            let block_id = args
                .get("block_id")
                .and_then(|v| v.as_str())
                .ok_or_else(|| Error::msg("block_id is required"))?;
            let children = args
                .get("children")
                .ok_or_else(|| Error::msg("children array is required"))?;

            make_notion_request(
                &format!("/blocks/{}/children", block_id),
                "PATCH",
                Some(serde_json::json!({ "children": children })),
                &token,
            )?
        }
        "retrieve_block" => {
            let block_id = args
                .get("block_id")
                .and_then(|v| v.as_str())
                .ok_or_else(|| Error::msg("block_id is required"))?;

            make_notion_request(&format!("/blocks/{}", block_id), "GET", None, &token)?
        }
        "retrieve_block_children" => {
            let block_id = args
                .get("block_id")
                .and_then(|v| v.as_str())
                .ok_or_else(|| Error::msg("block_id is required"))?;

            let mut params = Vec::new();
            if let Some(start_cursor) = args.get("start_cursor").and_then(|v| v.as_str()) {
                params.push(("start_cursor", start_cursor.to_string()));
            }
            if let Some(page_size) = args.get("page_size").and_then(|v| v.as_u64()) {
                params.push(("page_size", page_size.to_string()));
            }

            let query_string = if !params.is_empty() {
                let params: Vec<String> =
                    params.iter().map(|(k, v)| format!("{}={}", k, v)).collect();
                format!("?{}", params.join("&"))
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
        "delete_block" => {
            let block_id = args
                .get("block_id")
                .and_then(|v| v.as_str())
                .ok_or_else(|| Error::msg("block_id is required"))?;

            make_notion_request(&format!("/blocks/{}", block_id), "DELETE", None, &token)?
        }
        // Page operations
        "retrieve_page" => {
            let page_id = args
                .get("page_id")
                .and_then(|v| v.as_str())
                .ok_or_else(|| Error::msg("page_id is required"))?;

            make_notion_request(&format!("/pages/{}", page_id), "GET", None, &token)?
        }
        "update_page_properties" => {
            let page_id = args
                .get("page_id")
                .and_then(|v| v.as_str())
                .ok_or_else(|| Error::msg("page_id is required"))?;
            let properties = args
                .get("properties")
                .ok_or_else(|| Error::msg("properties are required"))?;

            make_notion_request(
                &format!("/pages/{}", page_id),
                "PATCH",
                Some(serde_json::json!({ "properties": properties })),
                &token,
            )?
        }
        // Database operations
        "create_database" => {
            let parent = args
                .get("parent")
                .ok_or_else(|| Error::msg("parent is required"))?;
            let title = args
                .get("title")
                .ok_or_else(|| Error::msg("title is required"))?;
            let properties = args
                .get("properties")
                .ok_or_else(|| Error::msg("properties are required"))?;

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
        "query_database" => {
            let database_id = args
                .get("database_id")
                .and_then(|v| v.as_str())
                .ok_or_else(|| Error::msg("database_id is required"))?;

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
        "retrieve_database" => {
            let database_id = args
                .get("database_id")
                .and_then(|v| v.as_str())
                .ok_or_else(|| Error::msg("database_id is required"))?;

            make_notion_request(&format!("/databases/{}", database_id), "GET", None, &token)?
        }
        "update_database" => {
            let database_id = args
                .get("database_id")
                .and_then(|v| v.as_str())
                .ok_or_else(|| Error::msg("database_id is required"))?;

            let mut update_body = serde_json::Map::new();
            if let Some(title) = args.get("title") {
                update_body.insert("title".to_string(), title.clone());
            }
            if let Some(description) = args.get("description") {
                update_body.insert("description".to_string(), description.clone());
            }
            if let Some(properties) = args.get("properties") {
                update_body.insert("properties".to_string(), properties.clone());
            }

            make_notion_request(
                &format!("/databases/{}", database_id),
                "PATCH",
                Some(Value::Object(update_body)),
                &token,
            )?
        }
        "create_database_item" => {
            let database_id = args
                .get("database_id")
                .and_then(|v| v.as_str())
                .ok_or_else(|| Error::msg("database_id is required"))?;
            let properties = args
                .get("properties")
                .ok_or_else(|| Error::msg("properties are required"))?;

            make_notion_request(
                "/pages",
                "POST",
                Some(serde_json::json!({
                    "parent": { "database_id": database_id },
                    "properties": properties
                })),
                &token,
            )?
        }
        // Comment operations
        "create_comment" => {
            let mut body = serde_json::Map::new();

            if let Some(parent) = args.get("parent") {
                body.insert("parent".to_string(), parent.clone());
            }
            if let Some(discussion_id) = args.get("discussion_id") {
                body.insert("discussion_id".to_string(), discussion_id.clone());
            }
            if let Some(rich_text) = args.get("rich_text") {
                body.insert("rich_text".to_string(), rich_text.clone());
            }

            if !body.contains_key("parent") && !body.contains_key("discussion_id") {
                return Err(Error::msg(
                    "Either parent.page_id or discussion_id must be provided",
                ));
            }

            make_notion_request("/comments", "POST", Some(Value::Object(body)), &token)?
        }
        "retrieve_comments" => {
            let block_id = args
                .get("block_id")
                .and_then(|v| Some(v.to_string()))
                .ok_or_else(|| Error::msg("block_id is required"))?;

            let mut params = vec![("block_id", block_id)];
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
        // User operations
        "list_users" => {
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

            make_notion_request(&format!("/users{}", query_string), "GET", None, &token)?
        }
        "retrieve_user" => {
            let user_id = args
                .get("user_id")
                .and_then(|v| v.as_str())
                .ok_or_else(|| Error::msg("user_id is required"))?;

            make_notion_request(&format!("/users/{}", user_id), "GET", None, &token)?
        }
        "retrieve_bot_user" => make_notion_request("/users/me", "GET", None, &token)?,
        // Search operation
        "search" => {
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
        _ => return Err(Error::msg(format!("Unknown operation: {}", operation))),
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
    Ok(types::ListToolsResult {
        tools: vec![
            ToolDescription {
                name: "notion".to_string(),
                description: "Interact with the Notion API to manage blocks, pages, databases, users, comments, and search functionality".to_string(),
                input_schema: serde_json::json!({
                    "type": "object",
                    "properties": {
                        "operation": {
                            "type": "string",
                            "description": "The operation to perform",
                            "enum": [
                                // Block operations
                                "append_block_children",
                                "retrieve_block",
                                "retrieve_block_children",
                                "delete_block",
                                // Page operations
                                "retrieve_page",
                                "update_page_properties",
                                // Database operations
                                "create_database",
                                "query_database",
                                "retrieve_database",
                                "update_database",
                                "create_database_item",
                                // Comment operations
                                "create_comment",
                                "retrieve_comments",
                                // User operations
                                "list_users",
                                "retrieve_user",
                                "retrieve_bot_user",
                                // Search operation
                                "search"
                            ]
                        },
                        // Block operation parameters
                        "block_id": {
                            "type": "string",
                            "description": "The ID of the block. Required for block operations."
                        },
                        "children": {
                            "type": "array",
                            "description": "Array of block objects to append. Required for append_block_children.",
                            "items": {
                                "type": "object",
                                "description": "A block object following Notion's block schema"
                            }
                        },
                        // Page operation parameters
                        "page_id": {
                            "type": "string",
                            "description": "The ID of the page. Required for page operations."
                        },
                        "properties": {
                            "type": "object",
                            "description": "Properties to update. Required for update_page_properties and create_database_item."
                        },
                        // Database operation parameters
                        "database_id": {
                            "type": "string",
                            "description": "The ID of the database. Required for database operations."
                        },
                        "parent": {
                            "type": "object",
                            "description": "Parent object specifying where to create the database. Required for create_database."
                        },
                        "title": {
                            "type": "array",
                            "description": "Title for database creation or update.",
                            "items": {
                                "type": "object",
                                "description": "Rich text object for title"
                            }
                        },
                        "filter": {
                            "type": "object",
                            "description": "Filter conditions for database queries or search"
                        },
                        "sorts": {
                            "type": "array",
                            "description": "Sort conditions for database queries"
                        },
                        // Comment operation parameters
                        "rich_text": {
                            "type": "array",
                            "description": "Rich text content for comments",
                            "items": {
                                "type": "object",
                                "description": "Rich text object"
                            }
                        },
                        "discussion_id": {
                            "type": "string",
                            "description": "ID of discussion thread for comments"
                        },
                        // Pagination parameters
                        "start_cursor": {
                            "type": "string",
                            "description": "Pagination cursor for subsequent requests"
                        },
                        "page_size": {
                            "type": "integer",
                            "description": "Number of results to return per request (max 100)"
                        },
                        // Search parameters
                        "query": {
                            "type": "string",
                            "description": "Search query string"
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
                        }
                    },
                    "required": ["operation"],
                    "additionalProperties": false
                })
                .as_object()
                .unwrap()
                .clone(),
            }
        ],
    })
}

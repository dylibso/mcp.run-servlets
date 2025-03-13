mod pdk;

use std::collections::BTreeMap;

use extism_pdk::*;
use htmd::HtmlToMarkdown;
use json::Value;
use pdk::types::{
    CallToolRequest, CallToolResult, Content, ContentType, ListToolsResult, ToolDescription,
};
use serde_json::json;

pub(crate) fn call(input: CallToolRequest) -> Result<CallToolResult, Error> {
    match input.params.name.as_str() {
        "fetch" => fetch(input),
        _ => Ok(CallToolResult {
            is_error: Some(true),
            content: vec![Content {
                annotations: None,
                text: Some(format!("Unknown tool: {}", input.params.name)),
                mime_type: None,
                r#type: ContentType::Text,
                data: None,
            }],
        }),
    }
}

fn fetch(input: CallToolRequest) -> Result<CallToolResult, Error> {
    let args = input.params.arguments.unwrap_or_default();
    if let Some(Value::String(url)) = args.get("url") {
        // Create HTTP request
        let mut req = HttpRequest {
            url: url.clone(),
            headers: BTreeMap::new(),
            method: Some("GET".to_string()),
        };

        // Add a user agent header to be polite
        req.headers
            .insert("User-Agent".to_string(), "fetch-tool/1.0".to_string());

        // Perform the request
        let res = http::request::<()>(&req, None)?;

        // Convert response body to string
        let body = res.body();
        let html = String::from_utf8_lossy(body.as_slice());

        let converter = HtmlToMarkdown::builder()
            .skip_tags(vec!["script", "style"])
            .build();

        // Convert HTML to markdown
        match converter.convert(&html) {
            Ok(markdown) => Ok(CallToolResult {
                is_error: None,
                content: vec![Content {
                    annotations: None,
                    text: Some(markdown),
                    mime_type: Some("text/markdown".to_string()),
                    r#type: ContentType::Text,
                    data: None,
                }],
            }),
            Err(e) => Ok(CallToolResult {
                is_error: Some(true),
                content: vec![Content {
                    annotations: None,
                    text: Some(format!("Failed to convert HTML to markdown: {}", e)),
                    mime_type: None,
                    r#type: ContentType::Text,
                    data: None,
                }],
            }),
        }
    } else {
        Ok(CallToolResult {
            is_error: Some(true),
            content: vec![Content {
                annotations: None,
                text: Some("Please provide a url".into()),
                mime_type: None,
                r#type: ContentType::Text,
                data: None,
            }],
        })
    }
}

// Called by mcpx to understand how and why to use this tool
pub(crate) fn describe() -> Result<ListToolsResult, Error> {
    Ok(ListToolsResult{
        tools: vec![
            ToolDescription {
                name: "fetch".into(),
                description:  "Enables to open and access arbitrary text URLs. Fetches the contents of a URL and returns its contents converted to markdown".into(),
                input_schema: json!({
                    "type": "object",
                    "properties": {
                        "url": {
                            "type": "string",
                            "description": "The URL to fetch",
                        },
                    },
                    "required": ["url"],
                })
                .as_object()
                .unwrap()
                .clone(),
            },
        ],
    })
}

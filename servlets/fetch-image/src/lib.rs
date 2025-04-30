mod pdk;

use std::collections::BTreeMap;

use base64::Engine;
use extism_pdk::*;
use json::Value;
use pdk::types::{
    CallToolRequest, CallToolResult, Content, ContentType, ListToolsResult, ToolDescription,
};
use serde_json::json;

pub(crate) fn call(input: CallToolRequest) -> Result<CallToolResult, Error> {
    match input.params.name.as_str() {
        "fetch-image" => fetch_image(input),
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

fn fetch_image(input: CallToolRequest) -> Result<CallToolResult, Error> {
    let args = input.params.arguments.unwrap_or_default();
    let opt_mime_type = args.get("mime-type");
    let mut mime_type = "image/png".to_string();
    if opt_mime_type.is_some() {
        if let Some(Value::String(mime)) = opt_mime_type {
            match mime.as_str() {
                "image/png" | "image/jpeg" | "image/gif" | "image/webp" | "image/svg+xml" => {
                    mime_type = mime.clone();
                }
                _ => {
                    return Ok(CallToolResult {
                        is_error: Some(true),
                        content: vec![Content {
                            annotations: None,
                            text: Some("Invalid mime type".into()),
                            mime_type: None,
                            r#type: ContentType::Text,
                            data: None,
                        }],
                    });
                }
            }
        }
    }

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

        // Let's filter by content type, we only want images
        req.headers.insert("Accept".to_string(), mime_type.clone());

        // Perform the request
        let res = http::request::<()>(&req, None)?;

        // Convert response body to string
        let body = res.body();

        let encoded_image = base64::engine::general_purpose::STANDARD.encode(body.as_slice());

        Ok(CallToolResult {
            is_error: None,
            content: vec![Content {
                annotations: None,
                data: Some(encoded_image),
                mime_type: Some(mime_type),
                r#type: ContentType::Image,
                text: None,
            }],
        })
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
                name: "fetch-image".into(),
                description:  "Enables to read images URLs. Fetches the contents of a URL pointing to an image and returns its contents converted to base64".into(),
                input_schema: json!({
                    "type": "object",
                    "properties": {
                        "url": {
                            "type": "string",
                            "description": "The URL of the image to fetch",
                        },
                        "mime-type": {
                            "type": "string",
                            "description": "The mime type to filter by, it must be of the form image/png, image/jpeg, etc. \
                            If the URL ends with an image extension, it should match with the mime type \
                            (e.g. if the URL ends with .png, the mime type should be image/png, if the URL ends with .jpg,\
                            the mime type should be image/jpeg, if the URL ends with .gif, the mime type should be image/gif, etc.).
                            If the mime-type is not provided it will default to image/png",
                        },
                    },
                })
                .as_object()
                .unwrap()
                .clone(),
            },

        ],
    })
}

mod pdk;

use std::collections::BTreeMap;

use extism_pdk::*;
use htmd::HtmlToMarkdown;
use json::{Map, Value};
use pdk::types::{CallToolRequest, CallToolResult, Content, ContentType, ToolDescription};

pub(crate) fn call(input: CallToolRequest) -> Result<CallToolResult, Error> {
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
pub(crate) fn describe() -> Result<ToolDescription, Error> {
    /*
    {
        name: "fetch",
        description: "Fetches a url and returns it's contents converted to markdown",
        inputSchema: {
          type: "object",
          properties: {
            url: {
              type: "string",
              description: "The URL to fetch",
            },
          },
          required: ["url"],
        },
    */
    let mut url_prop: Map<String, Value> = Map::new();
    url_prop.insert("type".into(), "string".into());
    url_prop.insert("description".into(), "The URL to fetch".into());

    let mut props: Map<String, Value> = Map::new();
    props.insert("url".into(), url_prop.into());

    let mut schema: Map<String, Value> = Map::new();
    schema.insert("type".into(), "object".into());
    schema.insert("properties".into(), Value::Object(props));
    schema.insert("required".into(), Value::Array(vec!["url".into()]));

    Ok(ToolDescription {
        name: "fetch".into(),
        description:
            "Fetches the contents of a URL and returns it's contents converted to markdown".into(),
        input_schema: schema,
    })
}


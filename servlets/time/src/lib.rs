mod pdk;

use extism_pdk::*;
use pdk::types::{CallToolResult, Content, ContentType, ToolDescription};
use pdk::*;
use serde_json::json;
use std::error::Error as StdError;

use chrono::Utc;

#[derive(Debug)]
struct CustomError(String);

impl std::fmt::Display for CustomError {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(f, "{}", self.0)
    }
}

impl StdError for CustomError {}

// Called when the tool is invoked.
pub(crate) fn call(input: types::CallToolRequest) -> Result<types::CallToolResult, Error> {
    let args = input.params.arguments.unwrap_or_default();
    let name = args.get("name").unwrap().as_str().unwrap();
    match name {
        "get_time_utc" => {
            let now = Utc::now();
            let timestamp = now.timestamp().to_string();
            let rfc2822 = now.to_rfc2822().to_string();
            Ok(CallToolResult {
                content: vec![Content {
                    text: Some(json!({
                        "utc_time": timestamp,
                        "utc_time_rfc2822": rfc2822,
                    }).to_string()),
                    r#type: ContentType::Text,
                    ..Default::default()
                }],
                is_error: Some(false),
            })
        }
        "parse_time" => {
            let time = args.get("time_rfc2822").unwrap().as_str().unwrap();
            let t = chrono::DateTime::parse_from_rfc2822(time).unwrap();
            let timestamp = t.timestamp().to_string();
            let rfc2822 = t.to_rfc2822().to_string();
            Ok(CallToolResult {
                content: vec![Content {
                    text: Some(json!({
                        "utc_time": timestamp,
                        "utc_time_rfc2822": rfc2822,
                    }).to_string()),
                    r#type: ContentType::Text,
                    ..Default::default()
                }],
                is_error: Some(false),
            })
        }
        "time_offset" => {
            let t1 = args.get("timestamp").unwrap().as_i64().unwrap();
            let offset = args.get("offset").unwrap().as_i64().unwrap();
            let t1 = chrono::DateTime::from_timestamp(t1, 0).unwrap();
            let t2 = t1 + chrono::Duration::seconds(offset);
            let timestamp = t2.timestamp().to_string();
            let rfc2822 = t2.to_rfc2822().to_string();
            Ok(CallToolResult {
                content: vec![Content {
                    text: Some(json!({
                        "utc_time": timestamp,
                        "utc_time_rfc2822": rfc2822,
                    }).to_string()),
                    r#type: ContentType::Text,
                    ..Default::default()
                }],
                is_error: Some(false),
            })
        }
        _ => Err(Error::new(CustomError("unknown command".to_string()))),
    }
}

// Called by mcpx to understand how and why to use this tool.
// Note: these imports are NOT available in this context: config_get
pub(crate) fn describe() -> Result<types::ListToolsResult, Error> {
    Ok(types::ListToolsResult { tools: vec![ToolDescription {
        name: "time".into(),
        description: "Time operations plugin. It provides the following operations:
        
- `get_time_utc`: Returns the current time in the UTC timezone. Takes no parameters.
- `parse_time`: Takes a `time_rfc2822` string in RFC2822 format and returns the timestamp in UTC timezone.
- `time_offset`: Takes integer `timestamp` and `offset` parameters. Adds a time offset to a given timestamp and returns the new timestamp in UTC timezone.
                
Always use this tool to compute time operations, especially when it is necessary
to compute time differences or offsets.".into(),
        input_schema: json!({
            "type": "object",
            "required": ["name"],
            "properties": {
                "name": {
                    "type": "string",
                    "description": "The name of the operation to perform. ",
                    "enum": ["get_time_utc", "time_offset",  "parse_time"],
                },
                "timestamp": {
                    "type": "integer",
                    "description": "The timestamp used for `time_offset`.",
                },
                "offset" : {
                    "type": "integer",
                    "description": "The offset to add to the time in seconds. ",
                },
                "time_rfc2822": {
                    "type": "string",
                    "description": "The time in RFC2822 format used in `parse_time`",
                },
            },
        })
        .as_object()
        .unwrap()
        .clone(),
    }]})
}

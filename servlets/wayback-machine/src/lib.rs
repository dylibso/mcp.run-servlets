mod pdk;

use extism_pdk::*;
use pdk::types::{CallToolResult, Content, ContentType, ToolDescription};
use pdk::*;
use query_string_builder::QueryString;
use serde_json::json;

// Called when the tool is invoked.
pub(crate) fn call(request: types::CallToolRequest) -> Result<types::CallToolResult, Error> {
    let args = request.params.arguments.unwrap_or_default();
    let qs = QueryString::simple()
        .with_value("url", args.get("url").unwrap())
        .with_opt_value("timestamp", args.get("timestamp"));

    let wmurl = format!("http://archive.org/wayback/available{qs}");
    let req = HttpRequest::new(wmurl);
    let res = http::request::<()>(&req, None)?;

    // http://archive.org/wayback/available?url=
    Ok(CallToolResult {
        content: vec![Content {
            text: Some(std::str::from_utf8(&res.body()).unwrap().to_string()),
            r#type: ContentType::Text,
            ..Default::default()
        }],
        is_error: Some(false),
    })
}

// Called by mcpx to understand how and why to use this tool.
// Note: these imports are NOT available in this context: config_get
pub(crate) fn describe() -> Result<types::ToolDescription, Error> {
    Ok(ToolDescription {
        name: "wayback-machine".into(),
        description: "Looks up a URL on the Wayback Machine \
                     and returns the snapshots that are available".into(),
        input_schema: json!({
            "type": "object",
            "required": ["name"],
            "properties": {
                "url": {
                    "type": "string",
                    "description": "The URL to lookup on the Wayback Machine",
                },
                "timestamp": {
                    "type": "string",
                    "description": "Optional timestamp. The format of the timestamp is 1-14 digits (YYYYMMDDhhmmss) ex:
                    timestamp=20060101 would return the closest snapshot to January 1, 2006. You can omit any parts of the timestamp,
                    ex: timestamp=2006 would return the closest snapshot to January 1, 2006."
                },
                "required" : ["url"]
            }
        })
            .as_object()
            .unwrap()
            .clone(),
    })
}

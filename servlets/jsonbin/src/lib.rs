mod pdk;

use extism_pdk::*;
use pdk::*;


pub(crate) fn call(input: types::CallToolRequest) -> Result<types::CallToolResult, Error> {
        return Err(Error::msg("Unknown tool name"));
    }

    let api_key = config::get("API_KEY")?        
        .expect("Missing `API_KEY` must be provided");

    let base_url = config::get("BASE_URL")?        
        .expect("Missing `BASE_URL` must be provided");


    let action = input
        .params
        .arguments
        .as_ref()
        .and_then(|args| args.get("action"))
        .and_then(|action| action.as_str())
        .ok_or_else(|| Error::msg("Argument `action` must be provided"))?;

    let path = input
        .params
        .arguments
        .as_ref()
        .and_then(|args| args.get("path"))
        .and_then(|action| action.as_str())
        .ok_or_else(|| Error::msg("Argument `path` must be provided"))?;

    let body = input
        .params
        .arguments
        .as_ref()
        .and_then(|args| args.get("body"))
        .and_then(|query| query.as_str());

    let result: String;

    match action {
        "post" => {
            let req = HttpRequest::new(format!("{}/{}", base_url, path))
                .with_method("POST")
                .with_header("Authorization", format!("token {}", api_key));
            let resp = http::request(&req, body)?;
            result = resp.json()?
        }
        "get" => {
            let req = HttpRequest::new(format!("{}/{}", base_url, path))
                .with_method("GET")
                .with_header("Authorization", format!("token {}", api_key));
            let resp = http::request(&req, body)?;
            result = resp.json()?
        }

        _ => return Err(Error::msg("Invalid action. Use 'post' or 'get'")),
    }

    Ok(types::CallToolResult {
        content: vec![types::Content {
            r#type: types::ContentType::Text,
            text: Some(result),
            annotations: None,
            data: None,
            mime_type: None,
        }],
        is_error: None,
    })
}

pub(crate) fn describe() -> Result<types::ListToolsResult, Error> {
    Ok(types::ListToolsResult {
        tools: vec![types::ToolDescription {
            name: "jsonbin".to_string(),
            description: "A tool to store and retrieve JSON payloads".to_string(),
            input_schema: serde_json::json!({
                "type": "object",
                "properties": {
                    "action": {
                        "type": "string",
                        "description": "The action to perform (post/get)",
                        "enum": ["post", "get"]
                    },
                    "path": {
                        "type": "string", 
                        "description": "the identifier of the jsonbin"
                    },
                    "body": {
                        "type": "string",
                        "description": "A JSON-formatted string with the body to send (required for post)"
                    },
                },
                "required": ["action", "path"]
            })
            .as_object()
            .unwrap()
            .clone(),
        }],
    })
}


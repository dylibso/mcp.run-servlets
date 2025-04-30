mod pdk;

use extism_pdk::*;
use pdk::*;


pub(crate) fn call(input: types::CallToolRequest) -> Result<types::CallToolResult, Error> {
    if input.params.name != "pantry" {
        return Err(Error::msg("Unknown tool name"));
    }

    let pantry_id = config::get("PANTRY_ID")?
        .expect("Missing `PANTRY_ID` must be provided");

    let base_url = format!("https://getpantry.cloud/apiv1/pantry/{}/basket", pantry_id);

    let action = input
        .params
        .arguments
        .as_ref()
        .and_then(|args| args.get("action"))
        .and_then(|action| action.as_str())
        .ok_or_else(|| Error::msg("Argument `action` must be provided"))?;

    let basket = input
        .params
        .arguments
        .as_ref()
        .and_then(|args| args.get("basket"))
        .and_then(|action| action.as_str())
        .ok_or_else(|| Error::msg("Argument `basket` must be provided"))?;

    let body = input
        .params
        .arguments
        .as_ref()
        .and_then(|args| args.get("body"))
        .and_then(|query| query.as_str());

    let result: String;

    match action {
        "post" => {
            let req = HttpRequest::new(format!("{}/{}", base_url, basket))
                .with_method("POST")
                .with_header("Content-Type", "application/json");
            let resp = http::request(&req, body)?;
            result = String::from_utf8(resp.body())?
        }
        "get" => {
            let req = HttpRequest::new(format!("{}/{}", base_url, basket))
                .with_method("GET");
            let resp = http::request(&req, body)?;
            result = String::from_utf8(resp.body())?
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
            name: "pantry".to_string(),
            description: "A tool to store and retrieve JSON payloads".to_string(),
            input_schema: serde_json::json!({
                "type": "object",
                "properties": {
                    "action": {
                        "type": "string",
                        "description": "The action to perform (post/get)",
                        "enum": ["post", "get"]
                    },
                    "basket": {
                        "type": "string", 
                        "description": "the identifier of the basket"
                    },
                    "body": {
                        "type": "string",
                        "description": "A JSON-formatted string with the body to send (required for post)"
                    },
                },
                "required": ["action", "basket"]
            })
            .as_object()
            .unwrap()
            .clone(),
        }],
    })
}


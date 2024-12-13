mod pdk;
use extism_pdk::*;
use pdk::*;
use url::form_urlencoded;

pub(crate) fn call(input: types::CallToolRequest) -> Result<types::CallToolResult, Error> {
    let query = input
        .params
        .arguments
        .as_ref()
        .and_then(|args| args.get("query"))
        .and_then(|query| query.as_str())
        .ok_or_else(|| Error::msg("Argument `query` must be provided"))?;

    let app_id = config::get("WOLFRAM_APP_ID")?
        .ok_or_else(|| Error::msg("config WOLFRAM_APP_ID must be set"))?;

    let url = format!(
        "https://www.wolframalpha.com/api/v1/llm-api?appid={}&input={}",
        form_urlencoded::byte_serialize(app_id.as_bytes()).collect::<String>(),
        form_urlencoded::byte_serialize(query.as_bytes()).collect::<String>(),
    );

    let req = HttpRequest::new(&url).with_method("GET");
    let res = http::request::<()>(&req, None)?;

    if res.status_code() != 200 {
        return Err(Error::msg(format!(
            "Wolfram Alpha API returned status code: {}",
            res.status()
        )));
    }

    // Get response body as string
    let response_text = String::from_utf8(res.body().to_vec())
        .map_err(|e| Error::msg(format!("Failed to parse response: {}", e)))?;

    Ok(types::CallToolResult {
        content: vec![types::Content {
            r#type: types::ContentType::Text,
            text: Some(response_text),
            annotations: None,
            data: None,
            mime_type: None,
        }],
        is_error: None,
    })
}

pub(crate) fn describe() -> Result<types::ToolDescription, Error> {
    Ok(types::ToolDescription {
        name: "wolfram_llm".to_string(),
        description: "Query the Wolfram Alpha LLM API".to_string(),
        input_schema: serde_json::json!({
            "type": "object",
            "properties": {
                "query": {
                    "type": "string",
                    "description": "the question or prompt to send to Wolfram Alpha's LLM"
                }
            },
            "required": ["query"]
        })
        .as_object()
        .unwrap()
        .clone(),
    })
}

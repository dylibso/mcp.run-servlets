mod pdk;

use extism_pdk::{http, HttpRequest, Error, config};
use uuid::Uuid;
use std::collections::BTreeMap;
use pdk::*;

// Called when the tool is invoked.
// If you support multiple tools, you must switch on the input.params.name to detect which tool is being called.
pub(crate) fn call(_input: types::CallToolRequest) -> Result<types::CallToolResult, Error> {
    let tool_name = _input.params.name.as_str();

    match tool_name {
        "create_minions" => {
            let minion_endpoint = config::get("MINION_ENDPOINT")
            .expect("'MINION_ENDPOINT' key set in config").unwrap();

            let prompts = _input
                .params
                .arguments
                .as_ref()
                .and_then(|args| args.get("prompts"))
                .and_then(|p| p.as_array())
                .ok_or_else(|| Error::msg("Argument `prompts` must be provided as an array"))?;

            // Generate mob_id and minion_ids
            let mob_id = Uuid::new_v4().to_string();
            let mut minion_ids = Vec::new();

            for prompt in prompts {
                let minion_id = Uuid::new_v4().to_string();
                minion_ids.push(minion_id.clone());

                let payload = serde_json::json!({
                    "prompt": prompt,
                    "minion_id": minion_id,
                    "mob_id": mob_id
                });

                // Build HTTP request for this minion
                let mut req = HttpRequest {
                    url: minion_endpoint.clone(),
                    headers: BTreeMap::new(),
                    method: Some("POST".to_string()),
                };
                req.headers.insert("User-Agent".to_string(), "minions-tool/1.0".to_string());
                req.headers.insert("Content-Type".to_string(), "application/json".to_string());

                // Send HTTP request for this minion
                let body = payload.to_string();
                let res = http::request::<Vec<u8>>(&req, Some(body.as_bytes().to_vec()))
                    .map_err(|e| Error::msg(format!("HTTP error: {:?}", e)))?;

                if res.status_code() < 200 || res.status_code() >= 300 {
                    return Err(Error::msg(format!("HTTP POST failed: status {}", res.status_code())));
                }
            }

            let result = serde_json::json!({
                "minion_ids": minion_ids,
                "mob_id": mob_id
            });
            Ok(types::CallToolResult {
                content: vec![types::Content {
                    r#type: types::ContentType::Text,
                    text: Some(result.to_string()),
                    annotations: None,
                    data: None,
                    mime_type: Some("application/json".to_string()),
                }],
                is_error: None,
            })
        }
        "check_minion_state" => {
            // Parse minion_id
            let minion_id = _input
                .params
                .arguments
                .as_ref()
                .and_then(|args| args.get("minion_id"))
                .and_then(|id| id.as_str())
                .ok_or_else(|| Error::msg("Argument `minion_id` must be provided"))?;

            // Get Pantry configuration from environment/config
            let pantry_id = config::get("PANTRY_ID")
                .expect("'PANTRY_ID' key set in config").unwrap();
            // The basket name is exactly the minion_id
            let basket_name = minion_id;

            // Build the GET URL
            let url = format!(
                "https://getpantry.cloud/apiv1/pantry/{}/basket/{}",
                pantry_id, basket_name
            );

            // Build HTTP GET request
            let mut req = HttpRequest {
                url,
                headers: BTreeMap::new(),
                method: Some("GET".to_string()),
            };
            req.headers.insert("User-Agent".to_string(), "minions-tool/1.0".to_string());

            // Send HTTP GET request
            let res = http::request::<Vec<u8>>(&req, None)
                .map_err(|e| Error::msg(format!("HTTP error: {:?}", e)))?;

            if res.status_code() < 200 || res.status_code() >= 300 {
                return Err(Error::msg(format!("HTTP GET failed: status {}", res.status_code())));
            }
            // Convert response body to string
            let body = String::from_utf8(res.body())
                .map_err(|e| Error::msg(format!("Invalid UTF-8 in response: {:?}", e)))?;

            Ok(types::CallToolResult {
                content: vec![types::Content {
                    r#type: types::ContentType::Text,
                    text: Some(body),
                    annotations: None,
                    data: None,
                    mime_type: Some("application/json".to_string()),
                }],
                is_error: None,
            })
        }
        "check_mob_state" => {
            // Placeholder: parse mob_id and return dummy state
            let mob_id = _input
                .params
                .arguments
                .as_ref()
                .and_then(|args| args.get("mob_id"))
                .and_then(|id| id.as_str())
                .ok_or_else(|| Error::msg("Argument `mob_id` must be provided"))?;
            let result = serde_json::json!({
                "mob_id": mob_id,
                "state": "pending"
            });
            Ok(types::CallToolResult {
                content: vec![types::Content {
                    r#type: types::ContentType::Text,
                    text: Some(result.to_string()),
                    annotations: None,
                    data: None,
                    mime_type: Some("application/json".to_string()),
                }],
                is_error: None,
            })
        }
        _ => {
            // Fallback: existing greet logic
            let name = _input
                .params
                .arguments
                .as_ref()
                .and_then(|args| args.get("name"))
                .and_then(|name| name.as_str())
                .ok_or_else(|| Error::msg("Argument `name` must be provided"))?;
            Ok(types::CallToolResult {
                content: vec![types::Content {
                    r#type: types::ContentType::Text,
                    text: Some(format!("Hello {}!!!", name)),
                    annotations: None,
                    data: None,
                    mime_type: None,
                }],
                is_error: None,
            })
        }
    }
}

// Called by mcpx to understand how and why to use this tool.
// Note: Your servlet configs will not be set when this function is called,
// so do not rely on config in this function
pub(crate) fn describe() -> Result<types::ListToolsResult, Error> {
    Ok(types::ListToolsResult {
        tools: vec![
            types::ToolDescription {
                name: "create_minions".to_string(),
                description: "Spawn a set of minions that perform actions in background. Returns minion IDs and a global mob ID.".to_string(),
                input_schema: serde_json::json!({
                    "type": "object",
                    "properties": {
                        "prompts": {
                            "type": "array",
                            "items": { "type": "string" },
                            "description": "List of prompts, one per minion"
                        }
                    },
                    "required": ["prompts"]
                })
                .as_object()
                .unwrap()
                .clone(),
            },
            types::ToolDescription {
                name: "check_minion_state".to_string(),
                description: "Check the completion or status of a single minion by its ID.".to_string(),
                input_schema: serde_json::json!({
                    "type": "object",
                    "properties": {
                        "minion_id": {
                            "type": "string",
                            "description": "The identifier of the minion to check"
                        }
                    },
                    "required": ["minion_id"]
                })
                .as_object()
                .unwrap()
                .clone(),
            },
            // types::ToolDescription {
            //     name: "check_mob_state".to_string(),
            //     description: "Check the aggregate completion or status of the mob (all minions) by the global mob ID.".to_string(),
            //     input_schema: serde_json::json!({
            //         "type": "object",
            //         "properties": {
            //             "mob_id": {
            //                 "type": "string",
            //                 "description": "The identifier of the mob to check"
            //             }
            //         },
            //         "required": ["mob_id"]
            //     })
            //     .as_object()
            //     .unwrap()
            //     .clone(),
            // },
        ],
    })
}

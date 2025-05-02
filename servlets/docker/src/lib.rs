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
        "list_containers" => {
            let docker_endpoint = 
            config::get("DOCKER_API_ENDPOINT")
                .expect("'DOCKER_API_ENDPOINT' key set in config").unwrap();
            let url = format!("{}/containers/json", docker_endpoint);
            let mut req = HttpRequest {
                url,
                headers: BTreeMap::new(),
                method: Some("GET".to_string()),
            };
            req.headers.insert("User-Agent".to_string(), "docker-servlet/1.0".to_string());
            let res = http::request::<Vec<u8>>(&req, None)
                .map_err(|e| Error::msg(format!("HTTP error: {:?}", e)))?;
            if res.status_code() < 200 || res.status_code() >= 300 {
                return Err(Error::msg(format!("HTTP GET failed: status {}", res.status_code())));
            }
            let body = String::from_utf8_lossy(&res.body()).to_string();
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
        "exec_in_container" => {
            let docker_endpoint = config::get("DOCKER_API_ENDPOINT")
                .expect("'DOCKER_API_ENDPOINT' key set in config").unwrap();
            let args = _input.params.arguments.as_ref().ok_or_else(|| Error::msg("Arguments must be provided"))?;
            let container_id = args.get("container_id").and_then(|v| v.as_str()).ok_or_else(|| Error::msg("Argument 'container_id' must be provided"))?;
            let cmd = args.get("cmd").and_then(|v| v.as_array()).ok_or_else(|| Error::msg("Argument 'cmd' must be provided as array"))?;
            let cmd_vec: Vec<String> = cmd.iter().filter_map(|v| v.as_str().map(|s| s.to_string())).collect();
            // Step 1: Create exec instance
            let url = format!("{}/containers/{}/exec", docker_endpoint, container_id);
            let mut req = HttpRequest {
                url,
                headers: BTreeMap::new(),
                method: Some("POST".to_string()),
            };
            req.headers.insert("User-Agent".to_string(), "docker-servlet/1.0".to_string());
            req.headers.insert("Content-Type".to_string(), "application/json".to_string());
            let payload = serde_json::json!({
                "AttachStdout": true,
                "AttachStderr": true,
                "Tty": false,
                "Cmd": cmd_vec
            });
            let res = http::request::<Vec<u8>>(&req, Some(payload.to_string().as_bytes().to_vec()))
                .map_err(|e| Error::msg(format!("HTTP error: {:?}", e)))?;
            if res.status_code() < 200 || res.status_code() >= 300 {
                return Err(Error::msg(format!("HTTP POST exec failed: status {}", res.status_code())));
            }
            let exec_resp: serde_json::Value = serde_json::from_slice(&res.body())
                .map_err(|e| Error::msg(format!("Failed to parse exec response: {:?}", e)))?;
            let exec_id = exec_resp.get("Id").and_then(|v| v.as_str()).ok_or_else(|| Error::msg("No exec ID returned"))?;
            // Step 2: Start exec instance
            let url = format!("{}/exec/{}/start", docker_endpoint, exec_id);
            let mut req = HttpRequest {
                url,
                headers: BTreeMap::new(),
                method: Some("POST".to_string()),
            };
            req.headers.insert("User-Agent".to_string(), "docker-servlet/1.0".to_string());
            req.headers.insert("Content-Type".to_string(), "application/json".to_string());
            let payload = serde_json::json!({
                "Detach": false,
                "Tty": false
            });
            let res = http::request::<Vec<u8>>(&req, Some(payload.to_string().as_bytes().to_vec()))
                .map_err(|e| Error::msg(format!("HTTP error: {:?}", e)))?;
            if res.status_code() < 200 || res.status_code() >= 300 {
                return Err(Error::msg(format!("HTTP POST exec/start failed: status {}", res.status_code())));
            }
            let body = String::from_utf8_lossy(&res.body()).trim().to_string();
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
        _ => {
            Err(Error::msg("Unknown tool name"))
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
                name: "list_containers".to_string(),
                description: "List all Docker containers (running by default).".to_string(),
                input_schema: serde_json::json!({
                    "type": "object",
                    "properties": {},
                    "required": []
                })
                .as_object()
                .unwrap()
                .clone(),
            },
            types::ToolDescription {
                name: "exec_in_container".to_string(),
                description: "Execute a command inside a running Docker container.".to_string(),
                input_schema: serde_json::json!({
                    "type": "object",
                    "properties": {
                        "container_id": {
                            "type": "string",
                            "description": "The ID or name of the container"
                        },
                        "cmd": {
                            "type": "array",
                            "items": { "type": "string" },
                            "description": "Command to execute (as array of arguments)"
                        }
                    },
                    "required": ["container_id", "cmd"]
                })
                .as_object()
                .unwrap()
                .clone(),
            },
        ],
    })
}

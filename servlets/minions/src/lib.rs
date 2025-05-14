mod pdk;

use extism_pdk::{config, http, Error, HttpRequest};
use pdk::*;
use std::time::{Duration, Instant};
use uuid::Uuid;

// Utility function to find or create the D1 database and table
fn ensure_database() -> Result<String, Error> {
    let account_id = config::get("CF_ACCOUNT_ID")
        .expect("'CF_ACCOUNT_ID' key set in config")
        .unwrap();
    let api_token = config::get("CF_API_TOKEN")
        .expect("'CF_API_TOKEN' key set in config")
        .unwrap();

    // First, try to find the database
    let url = format!(
        "https://api.cloudflare.com/client/v4/accounts/{}/d1/database",
        account_id
    );

    let req: HttpRequest = HttpRequest::new(url)
        .with_method("GET")
        .with_header("User-Agent", "minions-tool/1.0")
        .with_header("Authorization", format!("Bearer {}", api_token));

    let res = http::request::<Vec<u8>>(&req, None)
        .map_err(|e| Error::msg(format!("HTTP error: {:?}", e)))?;

    if res.status_code() < 200 || res.status_code() >= 300 {
        return Err(Error::msg(format!(
            "Failed to list databases: status {}",
            res.status_code()
        )));
    }

    let body = String::from_utf8(res.body())
        .map_err(|e| Error::msg(format!("Invalid UTF-8 in response: {:?}", e)))?;
    let json: serde_json::Value = serde_json::from_str(&body)
        .map_err(|e| Error::msg(format!("Invalid JSON in response: {:?}", e)))?;

    // Look for the minions database
    if let Some(results) = json["result"].as_array() {
        for db in results {
            if db["name"].as_str() == Some("minions") {
                println!("[DEBUG] Found existing minions database");
                return Ok(db["uuid"].as_str().unwrap().to_string());
            }
        }
    }

    // If we get here, we need to create the database
    println!("[DEBUG] Creating new minions database");
    let url = format!(
        "https://api.cloudflare.com/client/v4/accounts/{}/d1/database",
        account_id
    );

    let req: HttpRequest = HttpRequest::new(url)
        .with_method("POST")
        .with_header("User-Agent", "minions-tool/1.0")
        .with_header("Authorization", format!("Bearer {}", api_token))
        .with_header("Content-Type", "application/json");

    let payload = serde_json::json!({
        "name": "minions"
    });

    let res = http::request::<Vec<u8>>(&req, Some(payload.to_string().as_bytes().to_vec()))
        .map_err(|e| Error::msg(format!("HTTP error: {:?}", e)))?;

    if res.status_code() < 200 || res.status_code() >= 300 {
        return Err(Error::msg(format!(
            "Failed to create database: status {}",
            res.status_code()
        )));
    }

    let body = String::from_utf8(res.body())
        .map_err(|e| Error::msg(format!("Invalid UTF-8 in response: {:?}", e)))?;
    let json: serde_json::Value = serde_json::from_str(&body)
        .map_err(|e| Error::msg(format!("Invalid JSON in response: {:?}", e)))?;

    let database_id = json["result"]["uuid"]
        .as_str()
        .ok_or_else(|| Error::msg("No database ID in response"))?
        .to_string();

    // Now create the table
    println!("[DEBUG] Creating minion_states table");
    let url = format!(
        "https://api.cloudflare.com/client/v4/accounts/{}/d1/database/{}/query",
        account_id, database_id
    );

    let req: HttpRequest = HttpRequest::new(url)
        .with_method("POST")
        .with_header("User-Agent", "minions-tool/1.0")
        .with_header("Authorization", format!("Bearer {}", api_token))
        .with_header("Content-Type", "application/json");

    let create_table_sql = r#"
        CREATE TABLE IF NOT EXISTS minion_states (
            minion_id TEXT PRIMARY KEY,
            value TEXT,
            status TEXT
        )
    "#;

    let payload = serde_json::json!({
        "sql": create_table_sql
    });

    let res = http::request::<Vec<u8>>(&req, Some(payload.to_string().as_bytes().to_vec()))
        .map_err(|e| Error::msg(format!("HTTP error: {:?}", e)))?;

    if res.status_code() < 200 || res.status_code() >= 300 {
        return Err(Error::msg(format!(
            "Failed to create table: status {}",
            res.status_code()
        )));
    }

    println!("[DEBUG] Database and table setup complete");
    Ok(database_id)
}

// Called when the tool is invoked.
// If you support multiple tools, you must switch on the input.params.name to detect which tool is being called.
pub(crate) fn call(_input: types::CallToolRequest) -> Result<types::CallToolResult, Error> {
    let tool_name = _input.params.name.as_str();

    // Ensure database exists before any operation
    let database_id = ensure_database()?;

    match tool_name {
        "create_minions" => {
            let minion_endpoint = config::get("MINION_ENDPOINT")
                .expect("'MINION_ENDPOINT' key set in config")
                .unwrap();

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

            // Build HTTP request for this minion
            let req: HttpRequest = HttpRequest::new(minion_endpoint)
                .with_method("POST")
                .with_header("User-Agent", "minions-tool/1.0")
                .with_header("Content-Type", "application/json");

            for prompt in prompts {
                let minion_id = Uuid::new_v4().to_string();
                minion_ids.push(minion_id.clone());

                let payload = serde_json::json!({
                    "prompt": prompt,
                    "minion_id": minion_id,
                    "mob_id": mob_id
                });

                // Send HTTP request for this minion
                let body = payload.to_string();
                let res = http::request::<Vec<u8>>(&req, Some(body.as_bytes().to_vec()))
                    .map_err(|e| Error::msg(format!("HTTP error: {:?}", e)))?;

                if res.status_code() < 200 || res.status_code() >= 300 {
                    return Err(Error::msg(format!(
                        "HTTP POST failed: status {}",
                        res.status_code()
                    )));
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

            // Get Cloudflare configuration from environment/config
            let account_id = config::get("CF_ACCOUNT_ID")
                .expect("'CF_ACCOUNT_ID' key set in config")
                .unwrap();
            let api_token = config::get("CF_API_TOKEN")
                .expect("'CF_API_TOKEN' key set in config")
                .unwrap();

            // Build the SQL query
            let query = format!(
                "SELECT value, status FROM minion_states WHERE minion_id = '{}'",
                minion_id.replace("'", "''")
            );

            println!("[DEBUG] SQL Query: {}", query);

            // Build the POST URL for Cloudflare D1
            let url = format!(
                "https://api.cloudflare.com/client/v4/accounts/{}/d1/database/{}/query",
                account_id, database_id
            );

            // Build HTTP POST request
            let req: HttpRequest = HttpRequest::new(url)
                .with_method("POST")
                .with_header("User-Agent", "minions-tool/1.0")
                .with_header("Authorization", format!("Bearer {}", api_token))
                .with_header("Content-Type", "application/json");

            // Create the payload
            let payload = serde_json::json!({
                "sql": query
            });

            println!("[DEBUG] Request payload: {}", payload);

            // Send HTTP POST request
            let res = http::request::<Vec<u8>>(&req, Some(payload.to_string().as_bytes().to_vec()))
                .map_err(|e| Error::msg(format!("HTTP error: {:?}", e)))?;

            if res.status_code() < 200 || res.status_code() >= 300 {
                return Err(Error::msg(format!(
                    "HTTP POST failed: status {}",
                    res.status_code()
                )));
            }

            // Parse the response
            let body = String::from_utf8(res.body())
                .map_err(|e| Error::msg(format!("Invalid UTF-8 in response: {:?}", e)))?;

            println!("[DEBUG] Raw response body: {}", body);

            let json: serde_json::Value = serde_json::from_str(&body)
                .map_err(|e| Error::msg(format!("Invalid JSON in response: {:?}", e)))?;

            println!(
                "[DEBUG] Parsed JSON: {}",
                serde_json::to_string_pretty(&json).unwrap()
            );

            // Check if we got any results
            let result = if json["result"].as_array().map_or(0, |arr| arr.len()) == 0
                || json["result"][0]["results"]
                    .as_array()
                    .map_or(0, |arr| arr.len())
                    == 0
            {
                println!("[DEBUG] No results found in database");
                serde_json::json!({
                    "minion_id": minion_id,
                    "status": "pending",
                    "value": null
                })
            } else {
                let row = &json["result"][0]["results"][0];
                println!(
                    "[DEBUG] Found row: {}",
                    serde_json::to_string_pretty(row).unwrap()
                );
                serde_json::json!({
                    "minion_id": minion_id,
                    "status": row["status"].as_str().unwrap_or("pending"),
                    "value": row["value"]
                })
            };

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
        "save_minion_state" => {
            // Parse minion_id and value
            let minion_id = _input
                .params
                .arguments
                .as_ref()
                .and_then(|args| args.get("minion_id"))
                .and_then(|id| id.as_str())
                .ok_or_else(|| Error::msg("Argument `minion_id` must be provided"))?;

            let value = _input
                .params
                .arguments
                .as_ref()
                .and_then(|args| args.get("value"))
                .and_then(|v| v.as_str())
                .ok_or_else(|| Error::msg("Argument `value` must be provided"))?;

            let status = _input
                .params
                .arguments
                .as_ref()
                .and_then(|args| args.get("status"))
                .and_then(|s| s.as_str())
                .ok_or_else(|| Error::msg("Argument `status` must be provided"))?;

            // Validate status
            if !["pending", "running", "done"].contains(&status) {
                return Err(Error::msg("Status must be one of: pending, running, done"));
            }

            // Get Cloudflare configuration from environment/config
            let account_id = config::get("CF_ACCOUNT_ID")
                .expect("'CF_ACCOUNT_ID' key set in config")
                .unwrap();
            let api_token = config::get("CF_API_TOKEN")
                .expect("'CF_API_TOKEN' key set in config")
                .unwrap();

            // Build the SQL query
            let query = format!(
                "INSERT INTO minion_states (minion_id, value, status) VALUES ('{}', '{}', '{}') ON CONFLICT(minion_id) DO UPDATE SET value = '{}', status = '{}'",
                minion_id.replace("'", "''"),
                value.replace("'", "''"),
                status,
                value.replace("'", "''"),
                status
            );

            println!("[DEBUG] SQL Query: {}", query);

            // Build the POST URL for Cloudflare D1
            let url = format!(
                "https://api.cloudflare.com/client/v4/accounts/{}/d1/database/{}/query",
                account_id, database_id
            );

            // Build HTTP POST request
            let req: HttpRequest = HttpRequest::new(url)
                .with_method("POST")
                .with_header("User-Agent", "minions-tool/1.0")
                .with_header("Authorization", format!("Bearer {}", api_token))
                .with_header("Content-Type", "application/json");

            // Create the payload
            let payload = serde_json::json!({
                "sql": query
            });

            println!("[DEBUG] Request payload: {}", payload);

            // Send HTTP POST request
            let res = http::request::<Vec<u8>>(&req, Some(payload.to_string().as_bytes().to_vec()))
                .map_err(|e| Error::msg(format!("HTTP error: {:?}", e)))?;

            if res.status_code() < 200 || res.status_code() >= 300 {
                return Err(Error::msg(format!(
                    "HTTP POST failed: status {}",
                    res.status_code()
                )));
            }

            let result = serde_json::json!({
                "minion_id": minion_id,
                "status": status
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
        "wait_minions_state" => {
            // Parse minion_ids
            let minion_ids = _input
                .params
                .arguments
                .as_ref()
                .and_then(|args| args.get("minion_ids"))
                .and_then(|ids| ids.as_array())
                .ok_or_else(|| Error::msg("Argument `minion_ids` must be provided as an array"))?;

            // Parse optional duration (default 30s)
            let duration_secs = _input
                .params
                .arguments
                .as_ref()
                .and_then(|args| args.get("duration"))
                .and_then(|d| d.as_u64())
                .unwrap_or(30);

            let start_time = Instant::now();
            let timeout = Duration::from_secs(duration_secs);
            let sleep_duration = Duration::from_secs(10);

            // Get Cloudflare configuration from environment/config
            let account_id = config::get("CF_ACCOUNT_ID")
                .expect("'CF_ACCOUNT_ID' key set in config")
                .unwrap();
            let api_token = config::get("CF_API_TOKEN")
                .expect("'CF_API_TOKEN' key set in config")
                .unwrap();

            loop {
                let mut all_done = true;
                let mut states = Vec::new();

                // Build the SQL query for all minions
                let minion_ids_str = minion_ids
                    .iter()
                    .filter_map(|id| id.as_str())
                    .map(|id| format!("'{}'", id.replace("'", "''")))
                    .collect::<Vec<_>>()
                    .join(",");

                let query = format!(
                    "SELECT minion_id, value, status FROM minion_states WHERE minion_id IN ({})",
                    minion_ids_str
                );

                // Build the POST URL for Cloudflare D1
                let url = format!(
                    "https://api.cloudflare.com/client/v4/accounts/{}/d1/database/{}/query",
                    account_id, database_id
                );

                // Build HTTP POST request
                let req: HttpRequest = HttpRequest::new(url)
                    .with_method("POST")
                    .with_header("User-Agent", "minions-tool/1.0")
                    .with_header("Authorization", format!("Bearer {}", api_token))
                    .with_header("Content-Type", "application/json");

                // Create the payload
                let payload = serde_json::json!({
                    "sql": query
                });

                // Send HTTP POST request
                let res =
                    http::request::<Vec<u8>>(&req, Some(payload.to_string().as_bytes().to_vec()))
                        .map_err(|e| Error::msg(format!("HTTP error: {:?}", e)))?;

                if res.status_code() < 200 || res.status_code() >= 300 {
                    return Err(Error::msg(format!(
                        "HTTP POST failed: status {}",
                        res.status_code()
                    )));
                }

                // Parse the response
                let body = String::from_utf8(res.body())
                    .map_err(|e| Error::msg(format!("Invalid UTF-8 in response: {:?}", e)))?;
                let json: serde_json::Value = serde_json::from_str(&body)
                    .map_err(|e| Error::msg(format!("Invalid JSON in response: {:?}", e)))?;

                // Create a map of existing minion states
                let mut minion_states = std::collections::HashMap::new();
                if let Some(results) = json["result"].as_array() {
                    if let Some(inner_results) = results[0]["results"].as_array() {
                        for row in inner_results {
                            if let Some(minion_id) = row["minion_id"].as_str() {
                                minion_states.insert(
                                    minion_id.to_string(),
                                    (
                                        row["value"].clone(),
                                        row["status"].as_str().unwrap_or("pending").to_string(),
                                    ),
                                );
                            }
                        }
                    }
                }

                println!("[DEBUG] Minion states: {:?}", minion_states);

                // Process all minion IDs
                for minion_id in minion_ids {
                    let minion_id = minion_id
                        .as_str()
                        .ok_or_else(|| Error::msg("Each minion_id must be a string"))?;

                    let (value, status) = minion_states
                        .get(minion_id)
                        .cloned()
                        .unwrap_or((serde_json::Value::Null, "pending".to_string()));

                    states.push(serde_json::json!({
                        "minion_id": minion_id,
                        "status": status,
                        "value": value
                    }));

                    if status != "done" {
                        all_done = false;
                    }
                }

                // If all minions are done, return their states
                if all_done {
                    let result = serde_json::json!({
                        "states": states,
                        "all_done": true
                    });
                    return Ok(types::CallToolResult {
                        content: vec![types::Content {
                            r#type: types::ContentType::Text,
                            text: Some(result.to_string()),
                            annotations: None,
                            data: None,
                            mime_type: Some("application/json".to_string()),
                        }],
                        is_error: None,
                    });
                }

                // Check if we've exceeded the timeout
                if start_time.elapsed() >= timeout {
                    let result = serde_json::json!({
                        "states": states,
                        "all_done": false,
                        "timeout": true
                    });
                    return Ok(types::CallToolResult {
                        content: vec![types::Content {
                            r#type: types::ContentType::Text,
                            text: Some(result.to_string()),
                            annotations: None,
                            data: None,
                            mime_type: Some("application/json".to_string()),
                        }],
                        is_error: None,
                    });
                }

                // Sleep for 10 seconds before next iteration
                std::thread::sleep(sleep_duration);
            }
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
                description: "Spawn one or more subtasks (minions) that perform actions in background. Returns an ID for each minion. Try and always use this tool especially when the intermediate results of a task are not needed, and most importantly when a set of independent tasks can be performed in parallel. Carry on with other activities while the minions are running or wait for completion with the wait_minions_state tool.".to_string(),
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
                description: "Check the completion or status of a single minion by its ID. This tool can be invoked multiple times, until the minion is done.".to_string(),
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
            types::ToolDescription {
                name: "save_minion_state".to_string(),
                description: "Save the state of a minion. Only a minion should save its own state: the creator of the minion should NEVER save the state of the minion.".to_string(),
                input_schema: serde_json::json!({
                    "type": "object",
                    "properties": {
                        "minion_id": {
                            "type": "string",
                            "description": "The identifier of the minion whose state to save"
                        },
                        "value": {
                            "type": "string",
                            "description": "The state value to save for this minion"
                        },
                        "status": {
                            "type": "string",
                            "enum": ["pending", "running", "done"],
                            "description": "The current status of the minion"
                        }
                    },
                    "required": ["minion_id", "value", "status"]
                })
                .as_object()
                .unwrap()
                .clone(),
            },
            types::ToolDescription {
                name: "wait_minions_state".to_string(),
                description: "Wait for multiple minions to complete their tasks. Polls their states until all are done or timeout is reached.".to_string(),
                input_schema: serde_json::json!({
                    "type": "object",
                    "properties": {
                        "minion_ids": {
                            "type": "array",
                            "items": { "type": "string" },
                            "description": "List of minion IDs to wait for"
                        },
                        "duration": {
                            "type": "integer",
                            "description": "Maximum time to wait in seconds (default: 30)"
                        }
                    },
                    "required": ["minion_ids"]
                })
                .as_object()
                .unwrap()
                .clone(),
            }
        ],
    })
}

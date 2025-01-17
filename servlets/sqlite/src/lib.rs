mod pdk;

use extism_pdk::*;
use pdk::*;
use sqlite::{Connection, State};

pub(crate) fn call(input: types::CallToolRequest) -> Result<types::CallToolResult, Error> {
    if input.params.name != "sqlite" {
        return Err(Error::msg("Unknown tool name"));
    }

    let action = input
        .params
        .arguments
        .as_ref()
        .and_then(|args| args.get("action"))
        .and_then(|action| action.as_str())
        .ok_or_else(|| Error::msg("Argument `action` must be provided"))?;

    let query = input
        .params
        .arguments
        .as_ref()
        .and_then(|args| args.get("query"))
        .and_then(|query| query.as_str())
        .ok_or_else(|| Error::msg("Argument `query` must be provided"))?;

    let db_path = input
        .params
        .arguments
        .as_ref()
        .and_then(|args| args.get("db_path"))
        .and_then(|db_path| db_path.as_str())
        .ok_or_else(|| Error::msg("Argument `db_path` must be provided"))?;

    let conn = Connection::open(db_path)?;

    let mut result = String::new();

    match action {
        "query" => {
            let mut statement = conn.prepare(query)?;
            while let State::Row = statement.next()? {
                let mut row = Vec::new();
                for i in 0..statement.column_count() {
                    // Get the value based on column type
                    let value = match statement.column_type(i)? {
                        sqlite::Type::Null => "NULL".to_string(),
                        sqlite::Type::Integer => statement.read::<i64, _>(i)?.to_string(),
                        sqlite::Type::Float => statement.read::<f64, _>(i)?.to_string(),
                        sqlite::Type::String => statement.read::<String, _>(i)?,
                        sqlite::Type::Binary => "<BLOB>".to_string(), // or handle binary data as needed
                    };
                    row.push(value);
                }
                result.push_str(&row.join(","));
                result.push('\n');
            }
        }
        "execute" => {
            conn.execute(query)?;
            result = "Operation executed successfully".to_string();
        }
        _ => return Err(Error::msg("Invalid action. Use 'query' or 'execute'")),
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
            name: "sqlite".to_string(),
            description: "A tool to interact with SQLite databases".to_string(),
            input_schema: serde_json::json!({
                "type": "object",
                "properties": {
                    "action": {
                        "type": "string",
                        "description": "The action to perform (query/execute)",
                        "enum": ["query", "execute"]
                    },
                    "query": {
                        "type": "string",
                        "description": "The SQL query to execute"
                    },
                    "db_path": {
                        "type": "string", 
                        "description": "Path to the SQLite database file"
                    }
                },
                "required": ["action", "query", "db_path"]
            })
            .as_object()
            .unwrap()
            .clone(),
        }],
    })
}


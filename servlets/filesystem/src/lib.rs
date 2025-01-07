use extism_pdk::*;
use pdk::types::{CallToolRequest, CallToolResult, Content, ContentType, ListToolsResult, ToolDescription};
use serde_json::json;
use std::{error::Error as StdError, fs, path::PathBuf};
use walkdir::WalkDir;
mod pdk;

#[derive(Debug)]
struct CustomError(String);

impl std::fmt::Display for CustomError {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(f, "{}", self.0)
    }
}

impl StdError for CustomError {}

pub(crate) fn call(request: CallToolRequest) -> Result<CallToolResult, Error> {
    match request.params.name.as_str() {
        "read_file" => read_file(request),
        "write_file" => write_file(request),
        "list_directory" => list_directory(request),
        "search_files" => search_files(request),
        _ => Err(Error::new(CustomError("unknown tool".to_string()))),
    }
}

fn read_file(request: CallToolRequest) -> Result<CallToolResult, Error> {
    let args = request.params.arguments.unwrap_or_default();
    let path = PathBuf::from(
        args["path"]
            .as_str()
            .ok_or(Error::new(CustomError("path required".to_string())))?,
    );
    let content = fs::read_to_string(&path)?;
    Ok(CallToolResult {
        content: vec![Content {
            text: Some(content),
            r#type: ContentType::Text,
            ..Default::default()
        }],
        is_error: Some(false),
    })
}

fn write_file(request: CallToolRequest) -> Result<CallToolResult, Error> {
    let args = request.params.arguments.unwrap_or_default();
    let path = PathBuf::from(
        args["path"]
            .as_str()
            .ok_or(Error::new(CustomError("path required".to_string())))?,
    );
    let content = args["content"]
        .as_str()
        .ok_or(Error::new(CustomError("content required".to_string())))?;

    fs::write(&path, content)?;

    Ok(CallToolResult {
        content: vec![Content {
            text: Some(format!("Successfully wrote to {}", path.display())),
            r#type: ContentType::Text,
            ..Default::default()
        }],
        is_error: Some(false),
    })
}

fn list_directory(request: CallToolRequest) -> Result<CallToolResult, Error> {
    let args = request.params.arguments.unwrap_or_default();
    let path = PathBuf::from(
        args["path"]
            .as_str()
            .ok_or(Error::new(CustomError("path required".to_string())))?,
    );
    let entries = fs::read_dir(&path)?;
    let mut result = String::new();

    for entry in entries.flatten() {
        let prefix = if entry.file_type()?.is_dir() {
            "[DIR]"
        } else {
            "[FILE]"
        };
        result.push_str(&format!(
            "{} {}\n",
            prefix,
            entry.file_name().to_string_lossy()
        ));
    }

    Ok(CallToolResult {
        content: vec![Content {
            text: Some(result),
            r#type: ContentType::Text,
            ..Default::default()
        }],
        is_error: Some(false),
    })
}

fn search_files(request: CallToolRequest) -> Result<CallToolResult, Error> {
    let args = request.params.arguments.unwrap_or_default();
    let path = PathBuf::from(
        args["path"]
            .as_str()
            .ok_or(Error::new(CustomError("path required".to_string())))?,
    );
    let pattern = args["pattern"]
        .as_str()
        .ok_or(Error::new(CustomError("pattern required".to_string())))?
        .to_lowercase();

    let results: Vec<String> = WalkDir::new(&path)
        .into_iter()
        .filter_map(|e| e.ok())
        .filter(|e| {
            e.file_name()
                .to_string_lossy()
                .to_lowercase()
                .contains(&pattern)
        })
        .map(|e| e.path().display().to_string())
        .collect();

    Ok(CallToolResult {
        content: vec![Content {
            text: Some(results.join("\n")),
            r#type: ContentType::Text,
            ..Default::default()
        }],
        is_error: Some(false),
    })
}

pub(crate) fn describe() -> Result<ListToolsResult, Error> {
    Ok(ListToolsResult {
        tools: vec![
            ToolDescription {
                name: "read_file".into(),
                description: "Read contents of a file".into(),
                input_schema: json!({
                    "type": "object",
                    "required": ["path"],
                    "properties": {
                        "path": {
                            "type": "string",
                            "description": "File system path"
                        }
                    }
                })
                .as_object()
                .unwrap()
                .clone(),
            },
            ToolDescription {
                name: "write_file".into(),
                description: "Write content to a file".into(),
                input_schema: json!({
                    "type": "object",
                    "required": ["path", "content"],
                    "properties": {
                        "path": {
                            "type": "string",
                            "description": "File system path"
                        },
                        "content": {
                            "type": "string",
                            "description": "Content to write to file"
                        }
                    }
                })
                .as_object()
                .unwrap()
                .clone(),
            },
            ToolDescription {
                name: "list_directory".into(),
                description: "List contents of a directory".into(),
                input_schema: json!({
                    "type": "object",
                    "required": ["path"],
                    "properties": {
                        "path": {
                            "type": "string",
                            "description": "File system path"
                        }
                    }
                })
                .as_object()
                .unwrap()
                .clone(),
            },
            ToolDescription {
                name: "search_files".into(),
                description: "Search for files matching a pattern".into(),
                input_schema: json!({
                    "type": "object",
                    "required": ["path", "pattern"],
                    "properties": {
                        "path": {
                            "type": "string",
                            "description": "File system path"
                        },
                        "pattern": {
                            "type": "string",
                            "description": "Search pattern"
                        }
                    }
                })
                .as_object()
                .unwrap()
                .clone(),
            },
        ],
    })
}

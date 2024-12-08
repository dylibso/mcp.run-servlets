use extism_pdk::*;
use pdk::types::{CallToolRequest, CallToolResult, Content, ContentType, ToolDescription};
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
    let args = request.params.arguments.unwrap_or_default();
    let name = args.get("name").unwrap().as_str().unwrap();
    match name {
        "read_file" => {
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
        "write_file" => {
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
        "list_directory" => {
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
        "search_files" => {
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
        _ => Err(Error::new(CustomError("unknown command".to_string()))),
    }
}

pub(crate) fn describe() -> Result<ToolDescription, Error> {
    Ok(ToolDescription {
        name: "filesystem".into(),
        description: "File system operations plugin".into(),
        input_schema: json!({
            "type": "object",
            "required": ["name"],
            "properties": {
                "name": {
                    "type": "string",
                    "description": "The name of the operation to perform",
                    "enum": ["read_file", "write_file", "list_directory", "search_files"]
                },
                "path": {
                    "type": "string",
                    "description": "File system path"
                },
                "pattern": {
                    "type": "string",
                    "description": "Search pattern"
                },
                "content": {
                    "type": "string",
                    "description": "Content to write to file"
                },
            }
        })
        .as_object()
        .unwrap()
        .clone(),
    })
}



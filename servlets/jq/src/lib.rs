mod pdk;

use anyhow::anyhow;
use extism_pdk::*;
use jaq_core::{load, Ctx, RcIter};
use jaq_json::Val;
use pdk::types::{CallToolResult, Content, ContentType, ListToolsResult, ToolDescription};
use pdk::*;
use serde_json::{json, Value};

// Called when the tool is invoked.
// If you support multiple tools, you must switch on the input.params.name to detect which tool is being called.
pub(crate) fn call(input: types::CallToolRequest) -> Result<CallToolResult, Error> {
    match input.params.name.as_str() {
        "apply" => apply(input),
        _ => Ok(CallToolResult {
            is_error: Some(true),
            content: vec![Content {
                annotations: None,
                text: Some(format!("Unknown tool: {}", input.params.name)),
                mime_type: None,
                r#type: ContentType::Text,
                data: None,
            }],
        }),
    }
}

fn apply(input: types::CallToolRequest) -> Result<CallToolResult, Error> {
    let args = input.params.arguments.unwrap_or_default();
    if let (Some(Value::String(expr)), Some(Value::String(payload))) =
        (args.get("expression"), args.get("payload"))
    {
        use load::{Arena, File, Loader};
        let input: Value = serde_json::from_str(payload.as_str())?;
        let program = File {
            code: expr.as_str(),
            path: (),
        };

        let loader = Loader::new(jaq_std::defs().chain(jaq_json::defs()));
        let arena = Arena::default();

        // parse the filter
        let modules = loader.load(&arena, program).unwrap();

        // compile the filter
        let filter = jaq_core::Compiler::default()
            .with_funs(jaq_std::funs().chain(jaq_json::funs()))
            .compile(modules)
            .unwrap();

        let inputs = RcIter::new(core::iter::empty());

        // iterator over the output values

        let mut is_err = false;
        let mut results: Vec<String> = vec![];

        for output in filter.run((Ctx::new([], &inputs), Val::from(input))) {
            is_err = is_err || output.is_err();
            let value = output.map_err(|e| anyhow!("error: {e}"))?;
            results.push(value.to_string());
        }

        return Ok(CallToolResult {
            is_error: Some(is_err),
            content: vec![Content {
                annotations: None,
                text: Some(serde_json::to_string(&results)?),
                mime_type: None,
                r#type: ContentType::Text,
                data: None,
            }],
        });
    }

    Ok(CallToolResult {
        is_error: Some(true),
        content: vec![Content {
            annotations: None,
            text: Some("Please provide a url".into()),
            mime_type: None,
            r#type: ContentType::Text,
            data: None,
        }],
    })
}

// Called by mcpx to understand how and why to use this tool.
// Note: Your servlet configs will not be set when this function is called,
// so do not rely on config in this function
pub(crate) fn describe() -> Result<types::ListToolsResult, Error> {
    Ok(ListToolsResult {
        tools: vec![ToolDescription {
            name: "apply".into(),
            description: "Apply a jq filter to a JSON payload".into(),
            input_schema: json!({
                "type": "object",
                "properties": {
                    "expression": {
                        "type": "string",
                        "description": "The jq filter expression to apply"
                    },
                    "payload": {
                        "type": "string",
                        "description": "The JSON payload to apply the filter to"
                    }
                },
                "required": ["expression", "payload"]
            })
            .as_object()
            .unwrap()
            .clone(),
        }],
    })
}

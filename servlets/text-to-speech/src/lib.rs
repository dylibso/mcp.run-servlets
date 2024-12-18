mod pdk;

use std::hash::{DefaultHasher, Hash, Hasher};

use extism_pdk::*;
use pdk::*;
use types::{CallToolResult, Content, ContentType, ListToolsResult, ToolDescription};

// Called when the tool is invoked.
// If you support multiple tools, you must switch on the input.params.name to detect which tool is being called.
// The name will match one of the tool names returned from "describe".
pub(crate) fn call(input: types::CallToolRequest) -> Result<types::CallToolResult, Error> {
    let args = input.params.arguments.unwrap_or_default();
    let api_key: String = config::get("ELEVEN_LABS_API_KEY")
        .unwrap()
        .expect("ELEVEN_LABS_API_KEY should be set");
    let text = args["text"].as_str().unwrap();
    let path = args
        .get("path")
        .map(|x| x.as_str().unwrap())
        .unwrap_or("/tmp");

    info!("Converting text to speech: {text}");
    let voice = args
        .get("voice")
        .map(|x| x.as_str().unwrap())
        .unwrap_or_else(|| "flq6f7yk4E4fJM5XTYuZ");
    info!("With voice: {voice}");
    let body = serde_json::json!({
        "text": text,
    });

    info!("Making HTTP request");
    let res = http::request(
        &HttpRequest::new(format!(
            "https://api.elevenlabs.io/v1/text-to-speech/{voice}"
        ))
        .with_method("POST")
        .with_header("xi-api-key", api_key)
        .with_header("Content-Type", "application/json"),
        Some(body),
    )?;
    info!("Got HTTP response");

    let mut out = CallToolResult::default();
    let audio = res.body();
    if res.status_code() != 200 {
        out.is_error = Some(true);
        let mut c = Content::default();
        c.text = Some(
            serde_json::json!({
                "status_code": res.status_code(),
                "message": String::from_utf8(audio)?
            })
            .to_string(),
        );
        c.r#type = ContentType::Text;
        out.content.push(c);
        return Ok(out);
    }

    let mut h = DefaultHasher::new();
    text.hash(&mut h);

    let mut output = std::path::PathBuf::new();
    output.push(path);
    output.push(format!("text-to-speech.{}.mp3", h.finish()));
    info!("Saving audio to file: {}", output.display());
    std::fs::write(&output, audio)?;
    info!("Saved audio file");
    let mut c = Content::default();
    c.text = Some(
        serde_json::json!({
            "outputPath": output.to_str().unwrap()
        })
        .to_string(),
    );
    c.mime_type = Some("application/json".to_string());
    c.r#type = ContentType::Text;
    out.content = vec![c];
    Ok(out)
}

// Called by mcpx to understand how and why to use this tool.
// Note: Your servlet configs will not be set when this function is called,
// so do not rely on config in this function
pub(crate) fn describe() -> Result<types::ListToolsResult, Error> {
    Ok(ListToolsResult {
        tools: vec![ToolDescription {
            description: "Text-to-speech, generate audio of text using the Eleven Labs API. The LLM will output the path on disk where the audio file is saved."
                .to_string(),
            name: "text-to-speech".to_string(),
            input_schema: serde_json::json!({
                "type": "object",
                "required": ["text"],
                "properties": {
                    "text": {
                        "type": "string",
                        "description": "The text to be converted to audio",
                    },
                    "voice": {
                        "type": "string",
                        "description": "Eleven Labs voice ID"
                    },
                }
            })
            .as_object()
            .cloned()
            .unwrap(),
        }],
    })
}

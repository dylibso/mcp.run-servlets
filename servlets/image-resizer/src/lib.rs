mod pdk;

use crate::types::ListToolsResult;
use crate::types::ToolDescription;
use base64::{
    engine::general_purpose::STANDARD, engine::general_purpose::URL_SAFE_NO_PAD, Engine as _,
};
use extism_pdk::*;
use image::GenericImageView;
use image::ImageReader;
use pdk::*;
use serde_json::{Map, Value};
use std::io::Cursor;

// Called when the tool is invoked.
// If you support multiple tools, you must switch on the input.params.name to detect which tool is being called.
// The name will match one of the tool names returned from "describe".
pub(crate) fn call(_input: types::CallToolRequest) -> Result<types::CallToolResult, Error> {
    // load the params
    let b64_image = _input
        .params
        .arguments
        .as_ref()
        .and_then(|args| args.get("data"))
        .and_then(|data| data.as_str())
        .ok_or_else(|| Error::msg("Argument `data` must be provided"))?;
    let image_data = match URL_SAFE_NO_PAD.decode(b64_image) {
        Ok(data) => data,
        _ => STANDARD.decode(b64_image)?
    };
    let image = ImageReader::new(Cursor::new(image_data))
        .with_guessed_format()?
        .decode()?;
    let scale = _input
        .params
        .arguments
        .as_ref()
        .and_then(|args| args.get("scale"))
        .and_then(|scale| scale.as_number())
        .ok_or_else(|| Error::msg("Argument `scale` must be provided"))?;
    let scale = scale.as_f64().unwrap();

    // scale
    let (oldw, oldh) = image.dimensions();
    let neww = ((oldw as f64) * scale) as u32;
    let newh = ((oldh as f64) * scale) as u32;
    let image = image.resize(neww, newh, image::imageops::FilterType::Nearest);

    // return the result
    let mut result_bytes: Vec<u8> = Vec::new();
    image.write_to(&mut Cursor::new(&mut result_bytes), image::ImageFormat::Png)?;
    let result_text = URL_SAFE_NO_PAD.encode(result_bytes.clone());
    let result_image = STANDARD.encode(result_bytes);
    Ok(types::CallToolResult {
        content: vec![
            types::Content {
                r#type: types::ContentType::Image,
                text: None,
                annotations: None,
                data: Some(result_image),
                mime_type: Some("image/png".into()),
            },
            types::Content {
                r#type: types::ContentType::Text,
                text: Some(result_text),
                annotations: None,
                data: None,
                mime_type: None,
            },
        ],
        is_error: None,
    })
}

// Called by mcpx to understand how and why to use this tool.
// Note: Your servlet configs will not be set when this function is called,
// so do not rely on config in this function
pub(crate) fn describe() -> Result<types::ListToolsResult, Error> {
    let mut data_prop: Map<String, Value> = Map::new();
    data_prop.insert("type".into(), "string".into());
    data_prop.insert(
        "description".into(),
        "base64url data of image file to resize".into(),
    );

    let mut scale_prop: Map<String, Value> = Map::new();
    scale_prop.insert("type".into(), "number".into());
    scale_prop.insert(
        "description".into(),
        "Amount to scale image, for example, 1 to keep the same, 2 to double the size, etc".into(),
    );

    let mut props: Map<String, Value> = Map::new();
    props.insert("data".into(), data_prop.into());
    props.insert("scale".into(), scale_prop.into());

    let mut schema: Map<String, Value> = Map::new();
    schema.insert("type".into(), "object".into());
    schema.insert("properties".into(), Value::Object(props));
    schema.insert(
        "required".into(),
        Value::Array(vec!["data".into(), "scale".into()]),
    );

    Ok(ListToolsResult {
        tools: vec![ToolDescription {
            name: "image-resizer".into(),
            description: "Resize an image file".into(),
            input_schema: schema,
        }],
    })
}

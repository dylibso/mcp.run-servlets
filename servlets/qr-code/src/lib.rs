mod pdk;

use base64::Engine;
use extism_pdk::*;
use pdk::types::*;
use qrcode_png::{Color, QrCode, QrCodeEcc};
use serde_json::{json, Map, Value};

// Called when the tool is invoked.
pub(crate) fn call(input: CallToolRequest) -> Result<CallToolResult, Error> {
    extism_pdk::log!(
        LogLevel::Info,
        "called with args: {:?}",
        input.params.arguments
    );
    let args = input.params.arguments.unwrap_or_default();
    let ecc = to_ecc(
        args.get("ecc")
            .cloned()
            .unwrap_or_else(|| json!(4))
            .as_number()
            .unwrap()
            .is_u64() as u8,
    );

    let data = match args.get("data") {
        Some(v) => v.as_str().unwrap(),
        None => return Err(Error::msg("`data` must be available")),
    };

    let mut code = QrCode::new(data, ecc)?;
    code.margin(10);
    code.zoom(10);

    let b = code.generate(Color::Grayscale(0, 255))?;
    let data = base64::engine::general_purpose::STANDARD.encode(b);

    Ok(CallToolResult {
        is_error: None,
        content: vec![Content {
            annotations: None,
            text: None,
            mime_type: Some("image/png".into()),
            r#type: ContentType::Image,
            data: Some(data),
        }],
    })
}

fn to_ecc(num: u8) -> QrCodeEcc {
    if num < 4 {
        return unsafe { std::mem::transmute::<u8, QrCodeEcc>(num) };
    }

    QrCodeEcc::High
}

// Called by mcpx to understand how and why to use this tool
pub(crate) fn describe() -> Result<ListToolsResult, Error> {
    /*
    { tools: [{
        name: "qr_as_png",
        description: "Convert a URL to a QR code PNG",
        inputSchema: {
          type: "object",
          properties: {
            data: {
              type: "string",
              description: "data to convert to a QR code PNG",
            },
            ecc: {
              type: "number",
              description: "Error correction level (1-4, default to 4 unless user specifies)",
            },
          },
          required: ["data"],
        }]
    }
    */
    let mut data_prop: Map<String, Value> = Map::new();
    data_prop.insert("type".into(), "string".into());
    data_prop.insert(
        "description".into(),
        "data to convert to a QR code PNG".into(),
    );

    let mut ecc_prop: Map<String, Value> = Map::new();
    ecc_prop.insert("type".into(), "number".into());
    ecc_prop.insert(
        "description".into(),
        "Error correction level (range from 1 [low] to 4 [high], default to 4 unless user specifies)".into(),
    );

    let mut props: Map<String, Value> = Map::new();
    props.insert("data".into(), data_prop.into());
    props.insert("ecc".into(), ecc_prop.into());

    let mut schema: Map<String, Value> = Map::new();
    schema.insert("type".into(), "object".into());
    schema.insert("properties".into(), Value::Object(props));
    schema.insert("required".into(), Value::Array(vec!["data".into()]));

    Ok(ListToolsResult {
        tools: vec![ToolDescription {
            name: "qr-code".into(),
            description:
                "Convert data like a message or URL to a QR code (resulting in a PNG file)".into(),
            input_schema: schema,
        }],
    })
}

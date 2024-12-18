mod pdk;

use extism_pdk::*;
use pdk::*;
use types::{CallToolResult, Content, ContentType};

// Called when the tool is invoked.
// If you support multiple tools, you must switch on the input.params.name to detect which tool is being called.
// The name will match one of the tool names returned from "describe".
pub(crate) fn call(input: types::CallToolRequest) -> Result<types::CallToolResult, Error> {
    let args = input.params.arguments.unwrap_or_default();
    match input.params.name.as_str() {
        "ipfs-add" => {
            let endpoint = args
                .get("endpoint")
                .and_then(|x| x.as_str())
                .unwrap_or_else(|| "http://localhost:5001");
            let mut data = args["data"]
                .as_str()
                .map(|x| x.to_string())
                .unwrap_or_else(|| args["data"].to_string());
            let pin = args
                .get("pin")
                .map(|x| x.as_bool().unwrap_or_default())
                .unwrap_or_default();
            let req = HttpRequest::new(format!("{endpoint}/api/v0/add?pin={pin}"))
                .with_method("POST")
                .with_header("Abspath", "file.txt")
                .with_header(
                    "Content-Disposition",
                    r#"form-data; name="file"; filename="file.txt""#,
                )
                .with_header("Content-Type", "multipart/form-data; boundary=$$aaaaaaaaaa");
            data.insert_str(
                0,
                "\n\n--$$aaaaaaaaaa\nContent-Disposition: form-data; name=\"text\"\n\n",
            );
            data.push_str("\n--$$aaaaaaaaaa--\n");
            let res = http::request(&req, Some(data))?;
            if res.status_code() != 200 {
                let s = String::from_utf8(res.body()).unwrap();
                error!("Request error: {s}");
                return Err(Error::msg(s));
            }
            let res: serde_json::Map<String, serde_json::Value> = res.json()?;
            let hash = res["Hash"].as_str().unwrap().to_string();

            let mut res = CallToolResult::default();
            let mut c = Content::default();
            c.text = Some(hash);
            c.r#type = ContentType::Text;
            res.content = vec![c];
            Ok(res)
        }
        "ipfs-get" => {
            let endpoint = args
                .get("endpoint")
                .and_then(|x| x.as_str())
                .unwrap_or_else(|| "https://ipfs.io");
            let cid = args["cid"]
                .as_str()
                .map(|x| x.to_string())
                .unwrap_or_else(|| args["data"].to_string());
            let req = HttpRequest::new(format!("{endpoint}/ipfs/{cid}")).with_method("GET");
            let res = http::request::<()>(&req, None)?;

            let mut c = Content::default();
            match String::from_utf8(res.body()) {
                Ok(t) => {
                    c.text = Some(t);
                    c.r#type = ContentType::Text;
                }
                Err(_) => return Err(Error::msg("Got non-text data")),
            }
            let mut res = CallToolResult::default();
            res.content = vec![c];
            Ok(res)
        }
        name => Err(Error::msg(format!("Invalid tool name {name}"))),
    }
}

// Called by mcpx to understand how and why to use this tool.
// Note: Your servlet configs will not be set when this function is called,
// so do not rely on config in this function
pub(crate) fn describe() -> Result<types::ListToolsResult, Error> {
    todo!("Implement describe")
}

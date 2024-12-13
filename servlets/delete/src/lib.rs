mod pdk;

use extism_pdk::*;
use pdk::*;

// Called when the tool is invoked.
// If you support multiple tools, you must switch on the input.params.name to detect which tool is being called.
// The name will match one of the tool names returned from "describe".
pub(crate) fn call(_input: types::CallToolRequest) -> Result<types::CallToolResult, Error> {
    todo!("Implement call")
}

// Called by mcpx to understand how and why to use this tool.
// Note: Your servlet configs will not be set when this function is called,
// so do not rely on config in this function
pub(crate) fn describe() -> Result<types::ListToolsResult, Error> {
    todo!("Implement describe")
}

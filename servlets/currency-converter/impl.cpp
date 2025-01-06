#include "extism-pdk.hpp"
#include "pdk.gen.hpp"

static inline pdk::CallToolResult errorResult(std::string message) {
  return pdk::CallToolResult {
    .content = {
      {.text = std::move(message), .type = pdk::ContentType::text},
    },
    .isError = true
  };
}

/**
 * Called when the tool is invoked.
 * If you support multiple tools, you must switch on the input.params.name to
 * detect which tool is being called. The name will match one of the tool names
 * returned from "describe".
 *
 * @param input The incoming tool request from the LLM
 * @return The servlet's response to the given tool call
 */
std::expected<pdk::CallToolResult, pdk::Error>
impl::call(pdk::CallToolRequest &&input) {
  if (!input.params.arguments) {
    return errorResult("No arguments");
  }
  const auto& args = *input.params.arguments;
  const auto amount = args["amount"].as<double>();
  const auto from = args["from"].as_string_view();
  const auto to = args["to"].as_string_view();
  auto response = extism::http_request<char>(R"(
  {
    "method": "GET",
    "url":"https://api.fxratesapi.com/latest"
  }
  )");
  if (!response) {
    return errorResult("Failed to fetch latest rates");
  }
  const auto body = jsoncons::json::parse(response->body().string());
  if (!body["success"].as<bool>()) {
    return errorResult("Failed to fetch latest rates (2)");
  }
  const auto from_rate = body["rates"][from].as<double>();
  const auto to_rate = body["rates"][to].as<double>();
  const auto value = (amount / from_rate) * to_rate;
  const auto output_string = std::to_string(value);
  return pdk::CallToolResult {
    .content = {
      {.text = std::move(output_string), .type = pdk::ContentType::text},
    }
  };
}

/**
 * Called by mcpx to understand how and why to use this tool.
 * Note: Your servlet configs will not be set when this function is called,
 * so do not rely on config in this function
 *
 * @return The tools' descriptions, supporting multiple tools from a single
 * servlet.
 */
std::expected<pdk::ListToolsResult, pdk::Error> impl::describe() {
  return pdk::ListToolsResult{
    .tools = {
      pdk::ToolDescription{
        .inputSchema = jsoncons::json::parse(R"(
    {
      "type": "object",
      "properties": {
        "amount": {
          "type": "number",
          "description": "The amount of currency to convert."
        },
        "from": {
          "type": "string",
          "description": "The input type of currency to convert, the three letter ISO 4217 code, for example: USD or CAD or EUR."
        },
        "to": {
          "type": "string",
          "description": "The output type of currency to convert to, the three letter ISO 4217 code, for example: USD or CAD or EUR."
        }
      },
      "required": ["amount", "to", "from"]
    }
    )"),
        .description = "Currency converter",
        .name = "currency-converter",
      }
    }
  };
}

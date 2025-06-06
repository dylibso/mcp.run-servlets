// THIS FILE WAS GENERATED BY `xtp-cpp-bindgen`. DO NOT EDIT.
#include <chrono>
#include <cstddef>
#include <expected>
#include <jsoncons/json.hpp>
#include <memory>
#include <span>
#include <stdint.h>
#include <string>
#include <string_view>
#include <unordered_map>
#include <vector>

namespace pdk {

//
enum class ContentType { text, image, resource };

// The sender or recipient of messages and data in a conversation.
enum class Role { assistant, user };

// Describes the capabilities and expected paramters of the tool function
struct ToolDescription {
  // The JSON schema describing the argument input
  jsoncons::json inputSchema;
  // A description of the tool
  std::string description;
  // The name of the tool. It should match the plugin / binding name.
  std::string name;
};

// Provides one or more descriptions of the tools available in this servlet.
struct ListToolsResult {
  // The list of ToolDescription objects provided by this servlet.
  std::vector<ToolDescription> tools;
};

//
struct TextResourceContents {
  // The MIME type of this resource, if known.
  std::optional<std::string> mimeType;
  // The text of the item. This must only be set if the item can actually be
  // represented as text (not binary data).
  std::string text;
  // The URI of this resource.
  std::string uri;
};

// A text annotation
struct TextAnnotation {
  // Describes who the intended customer of this object or data is.  It can
  // include multiple entries to indicate content useful for multiple audiences
  // (e.g., `["user", "assistant"]`).
  std::vector<Role> audience;
  // Describes how important this data is for operating the server.  A value of
  // 1 means "most important," and indicates that the data is effectively
  // required, while 0 means "least important," and indicates that the data is
  // entirely optional.
  float priority;
};

// A content response. For text content set type to ContentType.Text and set the
// `text` property For image content set type to ContentType.Image and set the
// `data` and `mimeType` properties
struct Content {
  //
  std::optional<TextAnnotation> annotations;
  // The base64-encoded image data.
  std::optional<std::string> data;
  // The MIME type of the image. Different providers may support different image
  // types.
  std::optional<std::string> mimeType;
  // The text content of the message.
  std::optional<std::string> text;
  //
  ContentType type;
};

// The server's response to a tool call.  Any errors that originate from the
// tool SHOULD be reported inside the result object, with `isError` set to true,
// _not_ as an MCP protocol-level error response. Otherwise, the LLM would not
// be able to see that an error occurred and self-correct.  However, any errors
// in _finding_ the tool, an error indicating that the server does not support
// tool calls, or any other exceptional conditions, should be reported as an MCP
// error response.
struct CallToolResult {
  //
  std::vector<Content> content;
  // Whether the tool call ended in an error.  If not set, this is assumed to be
  // false (the call was successful).
  std::optional<bool> isError;
};

//
struct Params {
  //
  std::optional<jsoncons::json> arguments;
  //
  std::string name;
};

// Used by the client to invoke a tool provided by the server.
struct CallToolRequest {
  //
  Params params;
  //
  std::optional<std::string> method;
};

//
struct BlobResourceContents {
  // A base64-encoded string representing the binary data of the item.
  std::string blob;
  // The MIME type of this resource, if known.
  std::optional<std::string> mimeType;
  // The URI of this resource.
  std::string uri;
};

// host function errors
enum class Error { extism, host_null, not_json, json_null, not_implemented };

} // namespace pdk

namespace impl {

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
call(pdk::CallToolRequest &&input);

/**
 * Called by mcpx to understand how and why to use this tool.
 * Note: Your servlet configs will not be set when this function is called,
 * so do not rely on config in this function
 *
 * @return The tools' descriptions, supporting multiple tools from a single
 * servlet.
 */
std::expected<pdk::ListToolsResult, pdk::Error> describe();

} // namespace impl

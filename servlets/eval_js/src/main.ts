import { CallToolRequest, CallToolResult, ContentType, ToolDescription } from "./pdk";

/**
 * Called when the tool is invoked.
 *
 * @param {CallToolRequest} input - The incoming tool request from the LLM
 * @returns {CallToolResult} The servlet's response to the given tool call
 */
export function callImpl(input: CallToolRequest): CallToolResult {
  const code = input.params.arguments?.code
  if (!code) {
    throw new Error("Argument `code` must be provided")
  }
  return {
    content: [
      {
        type: ContentType.Text,
        text: eval(code).toString()
      }
    ]
  }

}

/**
 * Called by mcpx to understand how and why to use this tool
 *
 * @returns {ToolDescription} The tool's description
 */
export function describeImpl(): ToolDescription {
  return {
    name: "eval_js",
    description: "Evaluate some javascript using `eval()` in a sandbox.",
    inputSchema: {
      type: "object",
      properties: {
        code: {
          type: "string",
          description: "The javascript code to eval. This code gets passed into `eval()` and the result is stringified. Do not use console.log to emit the result.",
        },
      },
      required: ["code"],
    },
  }
}

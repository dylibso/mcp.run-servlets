import {
  CallToolRequest,
  CallToolResult,
  ContentType,
  ListToolsResult,
} from "./pdk";

import Ajv from "ajv";

/**
 * Called when the tool is invoked.
 * If you support multiple tools, you must switch on the input.params.name to detect which tool is being called.
 *
 * @param {CallToolRequest} input - The incoming tool request from the LLM
 * @returns {CallToolResult} The servlet's response to the given tool call
 */
export function callImpl(input: CallToolRequest): CallToolResult {
  const { name } = input.params;

  switch (name) {
    case "validate":
      return validate(input);
    default:
      throw new Error(`Unknown tool: ${name}`);
  }
}

function validate(input: CallToolRequest): CallToolResult {
  const { schema, document } = input.params.arguments;
  if (!schema) {
    throw new Error("Missing required argument: schema");
  }

  if (!document) {
    throw new Error("Missing required argument: document");
  }

  const ajv = new Ajv();
  const validate = ajv.compile(schema);
  const valid = validate(document);

  if (!valid) {
    return {
      content: [
        {
          type: ContentType.Text,
          mimeType: "application/json",
          text: JSON.stringify({ valid: false, errors: validate.errors}, null, 2)
        }
      ]
    }
  } else {
    return {
      content: [
        {
          type: ContentType.Text,
          mimeType: "application/json",
          text: JSON.stringify({ valid: true }, null, 2)
        }
      ]
    }
  }
}

/**
 * Called by mcpx to understand how and why to use this tool.
 * Note: Your servlet configs will not be set when this function is called,
 * so do not rely on config in this function
 *
 * @returns {ListToolsResult} The tools' descriptions, supporting multiple tools from a single servlet.
 */
export function describeImpl(): ListToolsResult {
  return {
    tools: [
      {
        name: "validate",
        description: "Validates a json document against a json schema. Returns a json object with a 'valid' property indicating if the document is valid or not. If the document is not valid, the 'errors' property will contain an array of validation errors.",
        inputSchema: {
          type: "object",
          properties: {
            schema: {
              type: "object",
              description: "The json schema to validate against.",
            },
            document: {
              type: "object",
              description: "The document to validate.",
            }
          },
          required: ["schema", "document"],
        },
      },
    ],
  };
}


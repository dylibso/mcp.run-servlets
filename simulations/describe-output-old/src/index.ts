import { Test } from "@dylibso/xtp-test";

export function test(): I32 {
  const output = Test.call("describe", undefined);
  let outcome = true;
  let reason = "parsed properly with expected values";

  try {
    // First get the raw JSON data
    const rawData = JSON.parse(output.text());

    // Then manually construct our result using the class constructors
    const tools = Array.isArray(rawData.tools)
      ? rawData.tools.map((toolData: ToolDescription) =>
        new ToolDescription(
          toolData.name || "",
          toolData.description || "",
          toolData.inputSchema || { type: "object", properties: {} },
        )
      )
      : [];

    const result = new ListToolsResult(tools);

    // Rest of validation code...
    if (result.tools.length === 0) {
      throw new Error("describe must provide at least one tool");
    }

    // Check tool descriptions and names
    const invalidTools = result.tools.filter(
      (tool) => !tool.description || !tool.name,
    );
    if (invalidTools.length > 0) {
      throw new Error(
        `Found tools with missing name or description: ${
          invalidTools.map((t) => t.name || "unnamed").join(", ")
        }`,
      );
    }

    // Validate each tool's schema
    result.tools.forEach((tool) => {
      const schema = tool.inputSchema;

      if (!schema.type) {
        throw new Error(`Tool ${tool.name}: missing inputSchema.type`);
      }

      if (!schema.properties) {
        throw new Error(`Tool ${tool.name}: missing inputSchema.properties`);
      }

      // Check required fields against properties
      if (schema.required) {
        const missingFields = schema.required.filter(
          (field) =>
            !Object.prototype.hasOwnProperty.call(schema.properties, field),
        );

        if (missingFields.length > 0) {
          throw new Error(
            `Tool ${tool.name}: required fields missing from properties: ${
              missingFields.join(", ")
            }`,
          );
        }
      }
    });
  } catch (e: any) {
    outcome = false;
    reason = String(e);
  }

  Test.assert("describe returns ListToolsResult JSON", outcome, reason);
  return 0;
}

export class ToolDescription {
  description: string;
  inputSchema: {
    type: string;
    properties: Record<string, unknown>;
    required?: string[];
  };
  name: string;

  constructor(name: string, description: string, inputSchema: {
    type: string;
    properties: Record<string, unknown>;
    required?: string[];
  }) {
    this.name = name;
    this.description = description;
    this.inputSchema = inputSchema;
  }

  static fromJson(obj: any): ToolDescription {
    if (!obj) {
      throw new Error("Cannot create ToolDescription from null or undefined");
    }
    return new ToolDescription(
      obj.name || "",
      obj.description || "",
      obj.inputSchema || { type: "object", properties: {} },
    );
  }

  static toJson(obj: ToolDescription): any {
    return {
      name: obj.name,
      description: obj.description,
      inputSchema: obj.inputSchema,
    };
  }
}

export class ListToolsResult {
  tools: Array<ToolDescription>;

  constructor(tools: Array<ToolDescription>) {
    this.tools = tools;
  }

  static fromJson(obj: any): ListToolsResult {
    if (!obj) {
      throw new Error("Cannot create ListToolsResult from null or undefined");
    }
    const tools = Array.isArray(obj.tools)
      ? obj.tools.map((tool: any) => ToolDescription.fromJson(tool))
      : [];
    return new ListToolsResult(tools);
  }

  static toJson(obj: ListToolsResult): any {
    return {
      tools: obj.tools.map((tool) => ToolDescription.toJson(tool)),
    };
  }
}

import { Test } from "@dylibso/xtp-test";
import Ajv from "ajv";
import draft7 from 'ajv/lib/refs/json-schema-draft-07.json';

export function test() {
  const result: ListToolsResult = Test.call("describe", "").json();

  for (const tool of result.tools) {
    if (!tool?.name) {
      Test.assert(`Tool name is provided`, false, tool.name);
      continue;
    }

    Test.group(tool.name, () => {
      const nameIsNotEmpty = !!tool.name && typeof tool.name === "string" && tool.name.length > 0;
      Test.assert(`Tool name is not empty: ${tool.name}`, nameIsNotEmpty, tool.name);

      const nameHasNoWhitespace = tool.name.trim() === tool.name;
      Test.assert(`Tool name has no leading/trailing whitespace`, nameHasNoWhitespace, tool.name);

      const descriptionIsNotEmpty = !!tool.description && typeof tool.description === "string" && tool.description.length > 0;
      Test.assert(`Tool description is not empty: ${tool.description}`, descriptionIsNotEmpty, tool.description);

      const descHasNoWhitespace = tool.description.trim() === tool.description;
      Test.assert(`Tool description has no leading/trailing whitespace`, descHasNoWhitespace, tool.description);

      const inputSchemaIsAnObject = !!tool.inputSchema && typeof tool.inputSchema === "object";
      Test.assert(`Tool input schema is an object`, inputSchemaIsAnObject, '');

      const inputSchemaHasProperties = !!tool.inputSchema && !!(tool.inputSchema as any).properties && typeof (tool.inputSchema as any).properties === "object";
      Test.assert(`Tool inputSchema.properties must be an object`, inputSchemaHasProperties, '');

      const inputSchemaHasTypeObject = !!tool.inputSchema && !!(tool.inputSchema as any).type && (tool.inputSchema as any).type === "object";
      Test.assert(`Tool inputSchema.type must be 'object'`, inputSchemaHasTypeObject, '');

      validateInputSchema(tool);
    });
  }
}

interface ListToolsResult {
  /** The list of ToolDescription objects provided by this servlet. */
  tools: ToolDescription[];
}

interface ToolDescription {
  /** A description of the tool */
  description: string;

  /** The JSON schema describing the argument input */
  inputSchema: unknown;

  /** The name of the tool. It should match the plugin / binding name. */
  name: string;
}

function validateInputSchema(tool: ToolDescription) {
  const ajv = new Ajv({
    strict: false,
    strictTypes: true,
    meta: true,
    allowUnionTypes: true,
    allErrors: true,
  });

  // the meta schema is draft-07 with some modifications
  const metaschema = {
    ...draft7,
    $id: 'tool-input',

    // input schema can't have additional properties
    additionalProperties: false,

    // properties must have type and description
    properties: {
      ...draft7.properties,
      properties: {
        type: "object",
        additionalProperties: {
          allOf: [
            { "$ref": "#" },
            { type: "object", required: ["type", "description"] }
          ]
        }
      }
    }
  };

  const validate = ajv.compile(metaschema);
  const valid = validate(tool.inputSchema);
  if (!valid) {
    Test.assert(`Tool input schema is valid`, false, JSON.stringify(validate.errors, null, 2));
    return
  }

  try {
    ajv.compile(tool.inputSchema as any);
  } catch (e: any) {
    Test.assert(`Tool input schema compiles`, false, e.message);
  }

  Test.assert(`Tool input schema is valid`, true, '');
}

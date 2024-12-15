import * as main from "./main";

import { CallToolRequest, CallToolResult, ToolDescription } from "./pdk";

export function call(): number {
  const jsonInput = Host.inputString();
  
  let input;
  if (jsonInput) {
    const untypedInput = JSON.parse(jsonInput);
    input = CallToolRequest.fromJson(untypedInput);
  } else {
    input = new CallToolRequest();
  }

  const output = main.callImpl(input);

  const untypedOutput = CallToolResult.toJson(output);
  Host.outputString(JSON.stringify(untypedOutput));

  return 0;
}

export function describe(): number {
  const output = main.describeImpl();

  const untypedOutput = ToolDescription.toJson(output);
  Host.outputString(JSON.stringify(untypedOutput));

  return 0;
}

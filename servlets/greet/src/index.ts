import * as main from "./main";

import { CallToolRequest, CallToolResult, ToolDescription } from "./pdk";

export function call(): number {
  const untypedInput = JSON.parse(Host.inputString());
  const input = CallToolRequest.fromJson(untypedInput);

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

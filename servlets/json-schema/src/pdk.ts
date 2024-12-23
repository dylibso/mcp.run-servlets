function isNull(v: any): boolean {
  return v === undefined || v === null;
}

function cast(caster: (v: any) => any, v: any): any {
  if (isNull(v)) return v;
  return caster(v);
}

function castArray(caster: (v: any) => any) {
  return (v?: Array<any>) => {
    if (isNull(v)) return v;
    caster = cast.bind(null, caster); // bind to null-preserving logic in `cast`
    return v!.map(caster);
  };
}

function castMap(caster: (v: any) => any) {
  return (v?: any) => {
    if (isNull(v)) return v;

    caster = cast.bind(null, caster); // bind to null-preserving logic in `cast`
    const newMap: any = {};
    for (const k in v) {
      newMap[k] = caster(v![k]);
    }
    return newMap;
  };
}

function dateToJson(v?: Date): string | undefined | null {
  if (v === undefined || v === null) return v;
  return v.toISOString();
}
function dateFromJson(v?: string): Date | undefined | null {
  if (v === undefined || v === null) return v;
  return new Date(v);
}

function bufferToJson(v?: ArrayBuffer): string | undefined | null {
  if (v === undefined || v === null) return v;
  return Host.arrayBufferToBase64(v);
}
function bufferFromJson(v?: string): ArrayBuffer | undefined | null {
  if (v === undefined || v === null) return v;
  return Host.base64ToArrayBuffer(v);
}

/**
 *
 */
export class BlobResourceContents {
  /**
   * A base64-encoded string representing the binary data of the item.
   */
  // @ts-expect-error TS2564
  blob: string;

  /**
   * The MIME type of this resource, if known.
   */
  mimeType?: string;

  /**
   * The URI of this resource.
   */
  // @ts-expect-error TS2564
  uri: string;

  static fromJson(obj: any): BlobResourceContents {
    return {
      ...obj,
    };
  }

  static toJson(obj: BlobResourceContents): any {
    return {
      ...obj,
    };
  }
}

/**
 * Used by the client to invoke a tool provided by the server.
 */
export class CallToolRequest {
  method?: string;

  // @ts-expect-error TS2564
  params: Params;

  static fromJson(obj: any): CallToolRequest {
    return {
      ...obj,
      params: cast(Params.fromJson, obj.params),
    };
  }

  static toJson(obj: CallToolRequest): any {
    return {
      ...obj,
      params: cast(Params.toJson, obj.params),
    };
  }
}

/**
 * The server's response to a tool call.
 *
 * Any errors that originate from the tool SHOULD be reported inside the result
 * object, with `isError` set to true, _not_ as an MCP protocol-level error
 * response. Otherwise, the LLM would not be able to see that an error occurred
 * and self-correct.
 *
 * However, any errors in _finding_ the tool, an error indicating that the
 * server does not support tool calls, or any other exceptional conditions,
 * should be reported as an MCP error response.
 */
export class CallToolResult {
  // @ts-expect-error TS2564
  content: Array<Content>;

  /**
   * Whether the tool call ended in an error.
   *
   * If not set, this is assumed to be false (the call was successful).
   */
  isError?: boolean;

  static fromJson(obj: any): CallToolResult {
    return {
      ...obj,
      content: cast(castArray(Content.fromJson), obj.content),
    };
  }

  static toJson(obj: CallToolResult): any {
    return {
      ...obj,
      content: cast(castArray(Content.toJson), obj.content),
    };
  }
}

/**
 * A content response.
 * For text content set type to ContentType.Text and set the `text` property
 * For image content set type to ContentType.Image and set the `data` and `mimeType` properties
 */
export class Content {
  annotations?: TextAnnotation;

  /**
   * The base64-encoded image data.
   */
  data?: string;

  /**
   * The MIME type of the image. Different providers may support different image types.
   */
  mimeType?: string;

  /**
   * The text content of the message.
   */
  text?: string;

  // @ts-expect-error TS2564
  type: ContentType;

  static fromJson(obj: any): Content {
    return {
      ...obj,
      annotations: cast(TextAnnotation.fromJson, obj.annotations),
    };
  }

  static toJson(obj: Content): any {
    return {
      ...obj,
      annotations: cast(TextAnnotation.toJson, obj.annotations),
    };
  }
}

/**
 *
 */
export enum ContentType {
  Text = "text",
  Image = "image",
  Resource = "resource",
}

/**
 * Provides one or more descriptions of the tools available in this servlet.
 */
export class ListToolsResult {
  /**
   * The list of ToolDescription objects provided by this servlet.
   */
  // @ts-expect-error TS2564
  tools: Array<ToolDescription>;

  static fromJson(obj: any): ListToolsResult {
    return {
      ...obj,
      tools: cast(castArray(ToolDescription.fromJson), obj.tools),
    };
  }

  static toJson(obj: ListToolsResult): any {
    return {
      ...obj,
      tools: cast(castArray(ToolDescription.toJson), obj.tools),
    };
  }
}

/**
 *
 */
export class Params {
  arguments?: any;

  // @ts-expect-error TS2564
  name: string;

  static fromJson(obj: any): Params {
    return {
      ...obj,
    };
  }

  static toJson(obj: Params): any {
    return {
      ...obj,
    };
  }
}

/**
 * The sender or recipient of messages and data in a conversation.
 */
export enum Role {
  Assistant = "assistant",
  User = "user",
}

/**
 * A text annotation
 */
export class TextAnnotation {
  /**
   * Describes who the intended customer of this object or data is.
   *
   * It can include multiple entries to indicate content useful for multiple audiences (e.g., `["user", "assistant"]`).
   */
  audience?: Array<Role>;

  /**
   * Describes how important this data is for operating the server.
   *
   * A value of 1 means "most important," and indicates that the data is
   * effectively required, while 0 means "least important," and indicates that
   * the data is entirely optional.
   */
  priority?: number;

  static fromJson(obj: any): TextAnnotation {
    return {
      ...obj,
    };
  }

  static toJson(obj: TextAnnotation): any {
    return {
      ...obj,
    };
  }
}

/**
 *
 */
export class TextResourceContents {
  /**
   * The MIME type of this resource, if known.
   */
  mimeType?: string;

  /**
   * The text of the item. This must only be set if the item can actually be represented as text (not binary data).
   */
  // @ts-expect-error TS2564
  text: string;

  /**
   * The URI of this resource.
   */
  // @ts-expect-error TS2564
  uri: string;

  static fromJson(obj: any): TextResourceContents {
    return {
      ...obj,
    };
  }

  static toJson(obj: TextResourceContents): any {
    return {
      ...obj,
    };
  }
}

/**
 * Describes the capabilities and expected paramters of the tool function
 */
export class ToolDescription {
  /**
   * A description of the tool
   */
  // @ts-expect-error TS2564
  description: string;

  /**
   * The JSON schema describing the argument input
   */
  inputSchema: any;

  /**
   * The name of the tool. It should match the plugin / binding name.
   */
  // @ts-expect-error TS2564
  name: string;

  static fromJson(obj: any): ToolDescription {
    return {
      ...obj,
    };
  }

  static toJson(obj: ToolDescription): any {
    return {
      ...obj,
    };
  }
}

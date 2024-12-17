import { CallToolRequest, CallToolResult, ListToolsResult, ContentType, ToolDescription } from "./pdk";

class SlackClient {
  private botHeaders: { Authorization: string; "Content-Type": string };
  private teamId: string;

  constructor(botToken: string, teamId: string) {
    this.botHeaders = {
      Authorization: `Bearer ${botToken}`,
      "Content-Type": "application/json",
    };
    this.teamId = teamId;
  }

  getChannels(limit: number = 100, cursor?: string): any {
    const params = new URLSearchParams({
      types: "public_channel",
      exclude_archived: "true",
      limit: Math.min(limit, 200).toString(),
      team_id: this.teamId,
    });

    if (cursor) {
      params.append("cursor", cursor);
    }

    return Http.request({
      url: `https://slack.com/api/conversations.list?${params}`,
      headers: this.botHeaders
    }).body
  }

  getChannelHistory(channel_id: string, limit: number = 10): any {
    const params = new URLSearchParams({
      channel: channel_id,
      limit: limit.toString(),
    });
  
    return Http.request({
      url: `https://slack.com/api/conversations.history?${params}`,
      headers: this.botHeaders
    }).body
  }

  getUsers(limit: number = 100, cursor?: string): any {
    const params = new URLSearchParams({
      limit: Math.min(limit, 200).toString(),
      team_id: this.teamId,
    });

    if (cursor) {
      params.append("cursor", cursor);
    }

    return Http.request({
      url: `https://slack.com/api/users.list?${params}`,
      headers: this.botHeaders
    }).body
  }
}

function ErrorContent(s: string) {
  return {
    content: [
      {
        text: s,
        type: ContentType.Text
      },
    ],
    isError: true
  }
}

/**
 * Called when the tool is invoked.
 * If you support multiple tools, you must switch on the input.params.name to detect which tool is being called.
 * The name will match one of the tool names returned from "describe".
 *
 * @param {CallToolRequest} input - The incoming tool request from the LLM
 * @returns {CallToolResult} The servlet's response to the given tool call
 */
export function callImpl(input: CallToolRequest): CallToolResult {
  const botToken = Config.get('SLACK_BOT_TOKEN');
  const teamId = Config.get('SLACK_TEAM_ID');
  if (!botToken) {
    return ErrorContent('Config SLACK_BOT_TOKEN not provided');
  }
  if (!teamId) {
    return ErrorContent('Config SLACK_TEAM_ID not provided')
  }
  const slackClient = new SlackClient(botToken as string, teamId as string);

  let text = 'unset response'
  if (input.params.name === 'slack_list_channels') {
    const {cursor, limit} = input.params.arguments || {}
    text = slackClient.getChannels(limit, cursor)
  } else if (input.params.name === 'slack_get_channel_history') {
    const { channel_id, limit } = input.params.arguments || {}
    if (!channel_id) {
      return ErrorContent('channel_id not provided')
    }
    text = slackClient.getChannelHistory(channel_id, limit);
  } else if (input.params.name === 'slack_get_users') {
    const {cursor, limit} = input.params.arguments || {}
    text = slackClient.getUsers(limit, cursor)
  } else {
    return ErrorContent(`Unknown command ${input.params.name}`);
  }

  return {
    content: [
      {
        text: text,
        type: ContentType.Text
      },
    ]
  };
}

/**
 * Called by mcpx to understand how and why to use this tool.
 * Note: Your servlet configs will not be set when this function is called,
 * so do not rely on config in this function
 *
 * @returns {ListToolsResult} The tools' descriptions, supporting multiple tools from a single servlet.
 */
export function describeImpl(): ListToolsResult {
  const listChannelsTool: ToolDescription = {
    name: "slack_list_channels",
    description: "List public channels in the workspace with pagination",
    inputSchema: {
      type: "object",
      properties: {
        limit: {
          type: "number",
          description:
            "Maximum number of channels to return (default 100, max 200)",
          default: 100,
        },
        cursor: {
          type: "string",
          description: "Pagination cursor for next page of results",
        },
      },
    },
  };

  const getChannelHistoryTool: ToolDescription = {
    name: "slack_get_channel_history",
    description: "Get recent messages from a channel. The names of the 'U' ids are available with slack_get_users",
    inputSchema: {
      type: "object",
      properties: {
        channel_id: {
          type: "string",
          description: "The ID of the channel",
        },
        limit: {
          type: "number",
          description: "Number of messages to retrieve (default 10)",
          default: 10,
        },
      },
      required: ["channel_id"],
    },
  };

  const getUsersTool: ToolDescription = {
    name: "slack_get_users",
    description:
      "Get a list of all users in the workspace with their basic profile information including their names",
    inputSchema: {
      type: "object",
      properties: {
        cursor: {
          type: "string",
          description: "Pagination cursor for next page of results",
        },
        limit: {
          type: "number",
          description: "Maximum number of users to return (default 100, max 200)",
          default: 100,
        },
      },
    },
  };

  return {
    tools: [listChannelsTool, getChannelHistoryTool, getUsersTool]
  };
}

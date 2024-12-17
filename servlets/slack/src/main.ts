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

  getChannels(limit: number = 100, cursor?: string): string {
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

  postMessage(channel_id: string, text: string): string {
    return Http.request({
      url: "https://slack.com/api/chat.postMessage",
      method: "POST",
      headers: this.botHeaders,
    }, JSON.stringify({
      channel: channel_id,
      text: text,
    })).body
  }

  postReply(
    channel_id: string,
    thread_ts: string,
    text: string,
  ): string {
    return Http.request({
      url:"https://slack.com/api/chat.postMessage",
      method: "POST",
      headers: this.botHeaders,
    }, JSON.stringify({
        channel: channel_id,
        thread_ts: thread_ts,
        text: text,
      })
    ).body
  }

  addReaction(
    channel_id: string,
    timestamp: string,
    reaction: string,
  ): string {
    return Http.request({
      url: "https://slack.com/api/reactions.add",
      method: "POST",
      headers: this.botHeaders,
    }, JSON.stringify({
        channel: channel_id,
        timestamp: timestamp,
        name: reaction,
      })
    ).body
  }

  getChannelHistory(channel_id: string, limit: number = 10): string {
    const params = new URLSearchParams({
      channel: channel_id,
      limit: limit.toString(),
    });
  
    return Http.request({
      url: `https://slack.com/api/conversations.history?${params}`,
      headers: this.botHeaders
    }).body
  }

  getThreadReplies(channel_id: string, thread_ts: string): string {
    const params = new URLSearchParams({
      channel: channel_id,
      ts: thread_ts,
    });

    return Http.request({
      url: `https://slack.com/api/conversations.replies?${params}`,
      headers: this.botHeaders
    },
    ).body;
  }

  getUsers(limit: number = 100, cursor?: string): string {
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

  getUserProfile(user_id: string): string {
    const params = new URLSearchParams({
      user: user_id,
      include_labels: "true",
    });

    return Http.request({
      url: `https://slack.com/api/users.profile.get?${params}`,
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

  let contentText = 'unset response'
  if (input.params.name === 'slack_list_channels') {
    const {cursor, limit} = input.params.arguments || {}
    contentText = slackClient.getChannels(limit, cursor)
  } else if (input.params.name === 'slack_post_message') {
    const {channel_id, text } = input.params.arguments || {}
    if (!channel_id || !text) {
      return ErrorContent('channel_id or text is missing')
    }
    contentText = slackClient.postMessage(channel_id, text)
  } else if (input.params.name === 'slack_reply_to_thread') {
    const {channel_id, thread_ts, text } = input.params.arguments || {}
    if (!channel_id || !thread_ts || !text) {
      return ErrorContent('channel_id or thread_ts or text is missing')
    }
    contentText = slackClient.postReply(channel_id, thread_ts, text)
  } else if (input.params.name === 'slack_add_reaction') {
    const {channel_id, timestamp, reaction} = input.params.arguments || {}
    if (!channel_id || !timestamp || !reaction) {
      return ErrorContent('channel_id or timestamp or reaction is missing')
    }
    contentText = slackClient.addReaction(channel_id, timestamp, reaction)
  } else if (input.params.name === 'slack_get_channel_history') {
    const { channel_id, limit } = input.params.arguments || {}
    if (!channel_id) {
      return ErrorContent('channel_id not provided')
    }
    contentText = slackClient.getChannelHistory(channel_id, limit);
  } else if (input.params.name === 'slack_get_thread_replies') {
    const {channel_id, thread_ts} = input.params.arguments || {}
    if (!channel_id || !thread_ts) {
      return ErrorContent('channel_id or thread_ts is missing')
    }
    contentText = slackClient.getThreadReplies(channel_id, thread_ts)
  } else if (input.params.name === 'slack_get_users') {
    const {cursor, limit} = input.params.arguments || {}
    contentText = slackClient.getUsers(limit, cursor)
  } else if (input.params.name === 'slack_get_user_profile') {
    const user_id = input.params.arguments?.user_id
    if (!user_id) {
      return ErrorContent('user_id is missing')
    }
    contentText = slackClient.getUserProfile(user_id)
  } else {
    return ErrorContent(`Unknown command ${input.params.name}`);
  }

  return {
    content: [
      {
        text: contentText,
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

  const postMessageTool: ToolDescription = {
    name: "slack_post_message",
    description: "Post a new message to a Slack channel",
    inputSchema: {
      type: "object",
      properties: {
        channel_id: {
          type: "string",
          description: "The ID of the channel to post to",
        },
        text: {
          type: "string",
          description: "The message text to post",
        },
      },
      required: ["channel_id", "text"],
    },
  };

  const replyToThreadTool: ToolDescription = {
    name: "slack_reply_to_thread",
    description: "Reply to a specific message thread in Slack",
    inputSchema: {
      type: "object",
      properties: {
        channel_id: {
          type: "string",
          description: "The ID of the channel containing the thread",
        },
        thread_ts: {
          type: "string",
          description: "The timestamp of the parent message",
        },
        text: {
          type: "string",
          description: "The reply text",
        },
      },
      required: ["channel_id", "thread_ts", "text"],
    },
  };

  const addReactionTool: ToolDescription = {
    name: "slack_add_reaction",
    description: "Add a reaction emoji to a message",
    inputSchema: {
      type: "object",
      properties: {
        channel_id: {
          type: "string",
          description: "The ID of the channel containing the message",
        },
        timestamp: {
          type: "string",
          description: "The timestamp of the message to react to",
        },
        reaction: {
          type: "string",
          description: "The name of the emoji reaction (without ::)",
        },
      },
      required: ["channel_id", "timestamp", "reaction"],
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

  const getThreadRepliesTool: ToolDescription = {
    name: "slack_get_thread_replies",
    description: "Get all replies in a message thread",
    inputSchema: {
      type: "object",
      properties: {
        channel_id: {
          type: "string",
          description: "The ID of the channel containing the thread",
        },
        thread_ts: {
          type: "string",
          description: "The timestamp of the parent message",
        },
      },
      required: ["channel_id", "thread_ts"],
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

  const getUserProfileTool: ToolDescription = {
    name: "slack_get_user_profile",
    description: "Get detailed profile information for a specific user",
    inputSchema: {
      type: "object",
      properties: {
        user_id: {
          type: "string",
          description: "The ID of the user",
        },
      },
      required: ["user_id"],
    },
  };

  return {
    tools: [listChannelsTool, postMessageTool, replyToThreadTool, addReactionTool, getChannelHistoryTool, getThreadRepliesTool, getUsersTool, getUserProfileTool]
  };
}

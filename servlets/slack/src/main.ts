import {
  CallToolRequest,
  CallToolResult,
  ContentType,
  ListToolsResult,
  ToolDescription,
} from "./pdk";

interface MessageOptions {
  text?: string;
  blocks?: any[];
  attachments?: any[];
  thread_ts?: string;
  markdown_text?: string;
  mrkdwn?: boolean;
  parse?: string;
  link_names?: boolean;
  unfurl_links?: boolean;
  unfurl_media?: boolean;
  metadata?: any;
  icon_emoji?: string;
  icon_url?: string;
  username?: string;
  reply_broadcast?: boolean;
}

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
      headers: this.botHeaders,
    }).body;
  }

  postMessage(channel_id: string, options: MessageOptions): string {
    // Ensure at least one of text, blocks, or attachments is provided
    if (
      !options.text && !options.blocks && !options.attachments &&
      !options.markdown_text
    ) {
      throw new Error(
        "At least one of text, blocks, attachments, or markdown_text must be provided",
      );
    }

    // Construct the message payload
    const payload: any = {
      channel: channel_id,
      ...options,
    };

    // If markdown_text is provided, use it as text and enable mrkdwn
    if (options.markdown_text) {
      payload.text = options.markdown_text;
      payload.mrkdwn = true;
      delete payload.markdown_text; // Remove the custom field before sending
    }

    // Ensure blocks and attachments are properly stringified
    if (payload.blocks && typeof payload.blocks === "string") {
      try {
        payload.blocks = JSON.parse(payload.blocks);
      } catch (e) {
        throw new Error("Invalid blocks JSON string");
      }
    }

    if (payload.attachments && typeof payload.attachments === "string") {
      try {
        payload.attachments = JSON.parse(payload.attachments);
      } catch (e) {
        throw new Error("Invalid attachments JSON string");
      }
    }

    if (payload.metadata && typeof payload.metadata === "string") {
      try {
        payload.metadata = JSON.parse(payload.metadata);
      } catch (e) {
        throw new Error("Invalid metadata JSON string");
      }
    }

    return Http.request({
      url: "https://slack.com/api/chat.postMessage",
      method: "POST",
      headers: this.botHeaders,
    }, JSON.stringify(payload)).body;
  }

  postReply(
    channel_id: string,
    thread_ts: string,
    text: string,
  ): string {
    return Http.request(
      {
        url: "https://slack.com/api/chat.postMessage",
        method: "POST",
        headers: this.botHeaders,
      },
      JSON.stringify({
        channel: channel_id,
        thread_ts: thread_ts,
        text: text,
      }),
    ).body;
  }

  addReaction(
    channel_id: string,
    timestamp: string,
    reaction: string,
  ): string {
    return Http.request(
      {
        url: "https://slack.com/api/reactions.add",
        method: "POST",
        headers: this.botHeaders,
      },
      JSON.stringify({
        channel: channel_id,
        timestamp: timestamp,
        name: reaction,
      }),
    ).body;
  }

  getChannelHistory(channel_id: string, limit: number = 10): string {
    const params = new URLSearchParams({
      channel: channel_id,
      limit: limit.toString(),
    });

    return Http.request({
      url: `https://slack.com/api/conversations.history?${params}`,
      headers: this.botHeaders,
    }).body;
  }

  getThreadReplies(channel_id: string, thread_ts: string): string {
    const params = new URLSearchParams({
      channel: channel_id,
      ts: thread_ts,
    });

    return Http.request({
      url: `https://slack.com/api/conversations.replies?${params}`,
      headers: this.botHeaders,
    }).body;
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
      headers: this.botHeaders,
    }).body;
  }

  getUserProfile(user_id: string): string {
    const params = new URLSearchParams({
      user: user_id,
      include_labels: "true",
    });

    return Http.request({
      url: `https://slack.com/api/users.profile.get?${params}`,
      headers: this.botHeaders,
    }).body;
  }
}

function ErrorContent(s: string) {
  return {
    content: [
      {
        text: s,
        type: ContentType.Text,
      },
    ],
    isError: true,
  };
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
  const botToken = Config.get("SLACK_BOT_TOKEN");
  const teamId = Config.get("SLACK_TEAM_ID");
  if (!botToken) {
    return ErrorContent("Config SLACK_BOT_TOKEN not provided");
  }
  if (!teamId) {
    return ErrorContent("Config SLACK_TEAM_ID not provided");
  }

  const slackClient = new SlackClient(botToken as string, teamId as string);
  let contentText = "unset response";
  switch (input.params.name) {
    case "slack_list_channels": {
      const { cursor, limit } = input.params.arguments || {};
      contentText = slackClient.getChannels(limit, cursor);
      break;
    }
    case "slack_post_message": {
      const { channel_id, ...messageOptions } = input.params.arguments || {};
      if (!channel_id) {
        return ErrorContent("channel_id is missing");
      }
      contentText = slackClient.postMessage(channel_id, messageOptions);
      break;
    }
    case "slack_reply_to_thread": {
      const { channel_id, thread_ts, text } = input.params.arguments || {};
      if (!channel_id || !thread_ts || !text) {
        return ErrorContent("channel_id or thread_ts or text is missing");
      }
      contentText = slackClient.postReply(channel_id, thread_ts, text);
      break;
    }
    case "slack_add_reaction": {
      const { channel_id, timestamp, reaction } = input.params.arguments || {};
      if (!channel_id || !timestamp || !reaction) {
        return ErrorContent("channel_id or timestamp or reaction is missing");
      }
      contentText = slackClient.addReaction(channel_id, timestamp, reaction);
      break;
    }
    case "slack_get_channel_history": {
      const { channel_id, limit } = input.params.arguments || {};
      if (!channel_id) {
        return ErrorContent("channel_id not provided");
      }
      contentText = slackClient.getChannelHistory(channel_id, limit);
      break;
    }
    case "slack_get_thread_replies": {
      const { channel_id, thread_ts } = input.params.arguments || {};
      if (!channel_id || !thread_ts) {
        return ErrorContent("channel_id or thread_ts is missing");
      }
      contentText = slackClient.getThreadReplies(channel_id, thread_ts);
      break;
    }
    case "slack_get_users": {
      const { cursor, limit } = input.params.arguments || {};
      contentText = slackClient.getUsers(limit, cursor);
      break;
    }
    case "slack_get_user_profile": {
      const user_id = input.params.arguments?.user_id;
      if (!user_id) {
        return ErrorContent("user_id is missing");
      }
      contentText = slackClient.getUserProfile(user_id);
      break;
    }
    default:
      return ErrorContent(`Unknown command ${input.params.name}`);
  }

  return {
    content: [
      {
        text: contentText,
        type: ContentType.Text,
      },
    ],
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
    description:
      "Post a new message to a Slack channel with support for rich formatting, blocks, and attachments. Either text, blocks, attachments, or markdown_text is required along with the channel_id.",
    inputSchema: {
      type: "object",
      properties: {
        channel_id: {
          type: "string",
          description: "The ID of the channel to post to",
        },
        text: {
          type: "string",
          description:
            "The message text to post. Required unless blocks or attachments are provided. When blocks/attachments are used, this becomes fallback text for notifications",
        },
        blocks: {
          type: "string",
          description:
            "A JSON string containing an array of Slack blocks for rich message formatting",
        },
        attachments: {
          type: "string",
          description:
            "A JSON string containing an array of Slack message attachments",
        },
        markdown_text: {
          type: "string",
          description:
            "Text formatted with markdown. Cannot be used with blocks. Limited to 12,000 characters",
        },
        thread_ts: {
          type: "string",
          description:
            "Timestamp of another message to reply to, making this message a thread reply",
        },
        reply_broadcast: {
          type: "boolean",
          description:
            "When replying to a thread, whether to also send the message to the channel",
        },
        parse: {
          type: "string",
          description:
            "Change how message text is treated. Can be 'none' or 'full'",
        },
        mrkdwn: {
          type: "boolean",
          description:
            "Whether to parse markdown-like syntax in the message. Defaults to true",
        },
        link_names: {
          type: "boolean",
          description: "Find and link user groups",
        },
        unfurl_links: {
          type: "boolean",
          description:
            "Whether to enable unfurling of primarily text-based content",
        },
        unfurl_media: {
          type: "boolean",
          description: "Whether to enable unfurling of media content",
        },
        metadata: {
          type: "string",
          description:
            "JSON string with event_type and event_payload fields for message metadata",
        },
        icon_emoji: {
          type: "string",
          description:
            "Emoji to use as the icon for this message. Overrides icon_url",
        },
        icon_url: {
          type: "string",
          description: "URL to an image to use as the icon for this message",
        },
        username: {
          type: "string",
          description: "Set your bot's user name for this message",
        },
      },
      required: ["channel_id"],
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
    description:
      "Get recent messages from a channel. The names of the 'U' ids are available with slack_get_users",
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
          description:
            "Maximum number of users to return (default 100, max 200)",
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
    tools: [
      listChannelsTool,
      postMessageTool,
      replyToThreadTool,
      addReactionTool,
      getChannelHistoryTool,
      getThreadRepliesTool,
      getUsersTool,
      getUserProfileTool,
    ],
  };
}

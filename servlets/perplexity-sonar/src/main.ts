// src/main.ts
import {
  CallToolRequest,
  CallToolResult,
  ContentType,
  ListToolsResult,
} from "./pdk";

interface ChatCompletionRequest {
  model: string;
  messages: Array<{
    role: "system" | "user" | "assistant";
    content: string;
  }>;
  max_tokens?: number;
  temperature?: number;
  top_p?: number;
  search_domain_filter?: string[];
  return_images?: boolean;
  return_related_questions?: boolean;
  search_recency_filter?: string;
  top_k?: number;
  stream?: boolean;
  presence_penalty?: number;
  frequency_penalty?: number;
}

interface ChatCompletionResponse {
  id: string;
  model: string;
  object: string;
  created: number;
  citations: string[];
  choices: Array<{
    index: number;
    finish_reason: string;
    message: {
      role: string;
      content: string;
    };
  }>;
  usage: {
    prompt_tokens: number;
    completion_tokens: number;
    total_tokens: number;
  };
}

export function callImpl(input: CallToolRequest): CallToolResult {
  const apiKey = Config.get("PERPLEXITY_API_KEY");
  if (!apiKey) {
    return {
      content: [{
        type: ContentType.Text,
        text: "Missing API key configuration",
      }],
      isError: true,
    };
  }

  const args = input.params.arguments;

  // Build the request body
  const requestBody: ChatCompletionRequest = {
    model: args.model || "sonar",
    messages: args.messages,
    temperature: args.temperature,
    max_tokens: args.max_tokens,
    top_p: args.top_p,
    search_domain_filter: args.search_domain_filter,
    return_images: args.return_images,
    return_related_questions: args.return_related_questions,
    search_recency_filter: args.search_recency_filter,
    top_k: args.top_k,
    stream: args.stream,
    presence_penalty: args.presence_penalty,
    frequency_penalty: args.frequency_penalty,
  };

  // Make the API request
  const response = Http.request({
    method: "POST",
    url: "https://api.perplexity.ai/chat/completions",
    headers: {
      "Authorization": `Bearer ${apiKey}`,
      "Content-Type": "application/json",
    },
  }, JSON.stringify(requestBody));

  // Parse and return the response
  if (response.status !== 200) {
    return {
      isError: true,
      content: [
        {
          type: ContentType.Text,
          text: response.body,
        },
      ],
    };
  }

  return {
    content: [{
      type: ContentType.Text,
      text: response.body,
    }],
  };
}

export function describeImpl(): ListToolsResult {
  return {
    tools: [{
      name: "perplexity-chat",
      description: "Make a chat completion request to the Perplexity API",
      inputSchema: {
        type: "object",
        properties: {
          model: {
            type: "string",
            description: "The name of the model to use (e.g. 'sonar')",
          },
          messages: {
            type: "array",
            items: {
              type: "object",
              properties: {
                role: {
                  type: "string",
                  enum: ["system", "user", "assistant"],
                },
                content: {
                  type: "string",
                },
              },
              required: ["role", "content"],
            },
            description: "The messages to send to the model",
          },
          temperature: {
            type: "number",
            description: "Controls randomness (0-2)",
            minimum: 0,
            maximum: 2,
          },
          max_tokens: {
            type: "integer",
            description: "Maximum number of tokens to generate",
          },
          top_p: {
            type: "number",
            description: "Controls diversity via nucleus sampling (0-1)",
            minimum: 0,
            maximum: 1,
          },
          search_domain_filter: {
            type: "array",
            items: {
              type: "string",
            },
            description: "Limit citations to specific domains",
          },
          return_images: {
            type: "boolean",
            description: "Whether to return images in the response",
          },
          return_related_questions: {
            type: "boolean",
            description: "Whether to return related questions",
          },
          search_recency_filter: {
            type: "string",
            enum: ["month", "week", "day", "hour"],
            description: "Filter search results by recency",
          },
        },
        required: ["messages"],
      },
    }],
  };
}

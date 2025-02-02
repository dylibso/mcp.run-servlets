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

  args.model ||= "sonar";

  // Build the request body
  const requestBody: ChatCompletionRequest = {
    model: args.model,
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
      text: formatAs(
        args.output_format || "markdown",
        response.body,
      ),
    }],
  };
}

interface PerplexityChoice {
  index: number;
  finish_reason: string;
  message: {
    role: string;
    content: string;
  };
  delta?: {
    role: string;
    content: string;
  };
}

interface PerplexityResponse {
  id: string;
  model: string;
  created: number;
  usage: {
    prompt_tokens: number;
    completion_tokens: number;
    total_tokens: number;
  };
  citations: string[];
  object: string;
  choices: PerplexityChoice[];
}

function formatPerplexityResponse(response: PerplexityResponse): string {
  // Extract all content from choices
  const contents = response.choices
    .map((choice) => choice.message.content)
    .filter(Boolean)
    .join("\n\n");

  // Format citations as markdown links
  const citationLinks = response.citations
    .map((url, index) => `[${index + 1}] ${url}`)
    .join("\n");

  // Combine content and citations with a separator
  const formattedOutput = `
${contents}

---
References:
${citationLinks}
`.trim();

  return formattedOutput;
}

function formatAs(format: string, data: string) {
  if (format === "markdown") {
    try {
      const response = JSON.parse(data) as PerplexityResponse;
      return formatPerplexityResponse(response);
    } catch (err) {
      return "There was an error converting the response to markdown, please try again but set the output_format to json.";
    }
  }

  return data;
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
            description:
              "The name of the model to use (e.g. 'sonar', 'sonar-pro', 'sonar-reasoning')",
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
          output_format: {
            type: "string",
            enum: ["json", "markdown"],
            description:
              "Format the output from the API call in JSON or Markdown - if Markdown, just report it back directly as-is.",
          },
        },
        required: ["messages"],
      },
    }],
  };
}

# Perplexity Sonar API

A servlet for accessing Perplexity's AI Sonar API through mcp.run. This servlet
enables you to have conversations with Perplexity's AI models, with support for
various parameters like temperature control, token limits, and search filters.

## Configuration

This servlet requires the following configuration:

- `PERPLEXITY_API_KEY`: Your Perplexity API key. You can obtain one from
  [Perplexity's website](https://www.perplexity.ai/).

## Usage

The servlet exposes a single tool called `perplexity-chat` with the following
capabilities:

### Basic Chat Request

Here's a simple example of using the chat completion endpoint:

```javascript
{
  "params": {
    "name": "perplexity-chat",
    "arguments": {
      "messages": [
        {
          "role": "user",
          "content": "What is the capital of France?"
        }
      ]
    }
  }
}
```

### Advanced Configuration

You can customize the request with various parameters:

```javascript
{
  "params": {
    "name": "perplexity-chat",
    "arguments": {
      "model": "sonar",  // Default model
      "messages": [
        {
          "role": "system",
          "content": "Be precise and concise."
        },
        {
          "role": "user",
          "content": "How many stars are there in our galaxy?"
        }
      ],
      "temperature": 0.2,
      "max_tokens": 100,
      "top_p": 0.9,
      "search_domain_filter": ["space.com", "nasa.gov"],
      "return_images": false,
      "return_related_questions": true,
      "search_recency_filter": "month"
    }
  }
}
```

### Parameters

| Parameter                  | Type    | Description                                                       |
| -------------------------- | ------- | ----------------------------------------------------------------- |
| `model`                    | string  | The name of the model to use (e.g., 'sonar')                      |
| `messages`                 | array   | Array of message objects with `role` and `content`                |
| `temperature`              | number  | Controls randomness (0-2)                                         |
| `max_tokens`               | integer | Maximum number of tokens to generate                              |
| `top_p`                    | number  | Controls diversity via nucleus sampling (0-1)                     |
| `search_domain_filter`     | array   | Limit citations to specific domains                               |
| `return_images`            | boolean | Whether to return images in the response                          |
| `return_related_questions` | boolean | Whether to return related questions                               |
| `search_recency_filter`    | string  | Filter search results by recency ("month", "week", "day", "hour") |

### Message Roles

The `messages` array supports three roles:

- `system`: Sets the behavior of the AI
- `user`: Your messages to the AI
- `assistant`: Previous AI responses (for context)

### Response Format

The API returns a response in this format:

```javascript
{
  "id": "3c90c3cc-0d44-4b50-8888-8dd25736052a",
  "model": "sonar",
  "object": "chat.completion",
  "created": 1724369245,
  "citations": [
    "https://example.com/source1",
    "https://example.com/source2"
  ],
  "choices": [
    {
      "index": 0,
      "finish_reason": "stop",
      "message": {
        "role": "assistant",
        "content": "The AI's response text here"
      }
    }
  ],
  "usage": {
    "prompt_tokens": 14,
    "completion_tokens": 70,
    "total_tokens": 84
  }
}
```

## Error Handling

The servlet will return an error response if:

- The API key is missing or invalid
- The request format is incorrect
- The API returns an error
- Network issues occur

Error responses include an `isError: true` flag and an error message in the
content.

import { CallToolRequest, CallToolResult, ContentType, ListToolsResult } from "./pdk";

export function callImpl(input: CallToolRequest): CallToolResult {
  const { to, from, body } = input.params.arguments || {};
  
  if (!to || !from || !body) {
    throw new Error("Arguments 'to', 'from', and 'body' must be provided");
  }

  const accountSid = Config.get('TWILIO_ACCOUNT_SID')
  const authToken = Config.get('TWILIO_AUTH_TOKEN')

  if (!accountSid || !authToken) {
    throw new Error("Twilio credentials not configured");
  }

  const auth = Host.arrayBufferToBase64((new TextEncoder()).encode(`${accountSid}:${authToken}`).buffer);
  const url = `https://api.twilio.com/2010-04-01/Accounts/${accountSid}/Messages.json`;

  const response = Http.request({
    url,
    method: 'POST',
    headers: {
      'Authorization': `Basic ${auth}`,
      'Content-Type': 'application/x-www-form-urlencoded'
    },
  },
    new URLSearchParams({
      To: to,
      From: from,
      Body: body
    }).toString()
  )

  return {
    content: [{
      type: ContentType.Text,
      text: response.body
    }]
  };
}

export function describeImpl(): ListToolsResult {
  return {
    tools: [{
      name: "send_sms",
      description: "Send an SMS message via Twilio",
      inputSchema: {
        type: "object",
        properties: {
          to: {
            type: "string",
            description: "The destination phone number"
          },
          from: {
            type: "string", 
            description: "The sender phone number (must be a Twilio number)"
          },
          body: {
            type: "string",
            description: "The message content"
          }
        },
        required: ["to", "from", "body"]
      }
    }]
  };
}


import { CallToolRequest, CallToolResult, ContentType, ListToolsResult, ToolDescription } from "./pdk";
import { decode } from 'base64-arraybuffer';
// Types
type TranscriptStatus = 'queued' | 'processing' | 'completed' | 'error';

interface Transcript {
  id?: string;
  status: TranscriptStatus;
  text: string;
  audio_url: string;
  error?: string;
}

interface UploadResponse {
  upload_url: string;
}

// AssemblyAI Client
class AssemblyAIClient {
  private apiKey: string;
  private baseURL: string;

  constructor(apiKey: string) {
    this.apiKey = apiKey;
    this.baseURL = 'https://api.assemblyai.com/v2';
  }

  upload(data: ArrayBuffer): string {
    console.log('Uploading audio data, length: ' + data.byteLength);
    const response = Http.request({
      method: 'POST',
      url: `${this.baseURL}/upload`,
      headers: {
        'Authorization': this.apiKey,
        'Content-Type': 'application/octet-stream'
      }
    }, data);

    if (response.status !== 200) {
      throw new Error(`Upload failed with status ${response.status}`);
    }

    const result = JSON.parse(response.body) as UploadResponse;
    return result.upload_url;
  }

  submitTranscript(audioUrl: string): Transcript {
    const response = Http.request({
      method: 'POST',
      url: `${this.baseURL}/transcript`,
      headers: {
        'Authorization': this.apiKey,
        'Content-Type': 'application/json'
      }},
      JSON.stringify({
        audio_url: audioUrl
      }));

    if (response.status !== 200) {
      throw new Error(`Transcript submission failed with status ${response.status}`);
    }

    return JSON.parse(response.body) as Transcript;
  }

  getTranscript(transcriptId: string): Transcript {
    const response = Http.request({
      method: 'GET',
      url: `${this.baseURL}/transcript/${transcriptId}`,
      headers: {
        'Authorization': this.apiKey
      }
    });

    if (response.status !== 200) {
      throw new Error(`Get transcript failed with status ${response.status}`);
    }

    return JSON.parse(response.body) as Transcript;
  }
}

function transcribeAudio(apiKey: string, audioData: string): Transcript {
  const client = new AssemblyAIClient(apiKey);

  const audioBuffer = decode(audioData);
  console.log('Audio decoded successfully with length: ' + audioBuffer.byteLength);

  // Upload the audio file
  const uploadUrl = client.upload(audioBuffer);
  console.log('Audio uploaded successfully');

  // Submit for transcription
  let transcript = client.submitTranscript(uploadUrl);
  console.log('Transcription submitted');

  // Poll until completion
  while (transcript.status !== 'completed' && transcript.status !== 'error') {
    console.log(`Checking status: ${transcript.status}. Text: ${transcript.text}`);

    if (!transcript.id) {
      throw new Error('Transcript ID not found');
    }

    transcript = client.getTranscript(transcript.id);

    if (transcript.status === 'error') {
      throw new Error(`Transcription failed: ${transcript.error}`);
    }
  }

  console.log('Transcription completed');
  return transcript;
}

// Main servlet functions
export function callImpl(input: CallToolRequest): CallToolResult {
  if (!input.params.arguments) {
    throw new Error('Arguments must be provided');
  }

  const args = input.params.arguments as { audio?: string };

  if (!args.audio) {
    throw new Error('Audio parameter must be provided as base64 string');
  }

  const apiKey = Config.get('ASSEMBLYAI_API_KEY');
  if (!apiKey) {
    throw new Error('ASSEMBLYAI_API_KEY config must be set!');
  }

  const transcript = transcribeAudio(apiKey, args.audio);
  console.log(`Transcript: ${transcript.text}`);

  return {
    content: [{
      type: ContentType.Text,
      text: transcript.text
    }]
  };
}

export function describeImpl(): ListToolsResult {
  return {
    tools: [
      {
        name: 'transcribe',
        description: 'Transcribes an audio file using AssemblyAI. Supports mp3, wav, and other common audio formats.',
        inputSchema: {
          type: 'object',
          required: ['audio'],
          properties: {
            audio: {
              type: 'string',
              description: 'Base64 encoded audio file content'
            }
          }
        }
      }
    ],
  }
}
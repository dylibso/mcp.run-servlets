from typing import Optional, List  # noqa: F401
from datetime import datetime  # noqa: F401
import extism  # noqa: F401 # pyright: ignore
import traceback
from io import StringIO, BytesIO
import sys
from base64 import b64encode

from pdk_types import (
    BlobResourceContents,
    CallToolRequest,
    CallToolResult,
    Content,
    ContentType,
    ListToolsResult,
    Params,
    Role,
    TextAnnotation,
    TextResourceContents,
    ToolDescription,
)  # noqa: F401


from typing import List, Optional  # noqa: F401

class StdoutCapture:
    def __init__(self):
        self.str = StringIO()
        self.buffer = BytesIO()

    def write(self, data):
        # Handle both str and bytes input
        if isinstance(data, str):
            self.str.write(data)
        else:
            self.buffer.write(data)

    def flush(self):
        pass

    def getvalue(self):
        return [self.str.getvalue(), self.buffer.getvalue()]

    def close(self):
        self.str.close()
        self.buffer.close()

def guess_mime_type(binary_data):
    signatures = {
        b'\xFF\xD8\xFF': 'image/jpeg',                    # JPEG
        b'\x89PNG\r\n\x1A\n': 'image/png',               # PNG
        b'GIF87a': 'image/gif',                          # GIF87a
        b'GIF89a': 'image/gif',                          # GIF89a
        b'BM': 'image/bmp',                              # BMP
        b'\x00\x00\x01\x00': 'image/x-icon',            # ICO
        b'RIFF': 'image/webp',                           # WebP
    }

    for signature, mime_type in signatures.items():
        if binary_data.startswith(signature):
            return mime_type

    raise ValueError("Unknown image format: could not detect mime type from file signature")


# Called when the tool is invoked.
# If you support multiple tools, you must switch on the input.params.name to detect which tool is being called.
def call(input: CallToolRequest) -> CallToolResult:
    try:
        code = input.params.arguments['code']
        output = StdoutCapture()
        old_stdout = sys.stdout
        sys.stdout = output
        exec(code)
        sys.stdout = old_stdout
        result = output.getvalue()
        output.close()
        isError = False
    except Exception as e:
        result = "\n".join([
            "Traceback (most recent call last):",
            traceback.format_tb(e.__traceback__)[0].rstrip(),
            f"{type(e).__name__}: {e}"
        ])
        return CallToolResult(
            content=[
                Content(
                    text=result,
                    type=ContentType.Text,
                )
            ],
            isError=True,
        )
    content = []
    if len(result[0]) or len(result[1]) == 0:
        content = [
            Content(
                type=ContentType.Text,
                text=result[0],
            )
        ]
    if len(result[1]):
        try:
            mimeType = guess_mime_type(result[1])
            content.append(
                Content(
                    type=ContentType.Image,
                    mimeType=mimeType,
                    data=b64encode(result[1]).decode('ascii')
                )
            )
        except:
            content.append(
                Content(
                    type=ContentType.Text,
                    text=f'Unknown image format: {result[1]}',
                )
            )
            isError=True
    return CallToolResult(
        content=content,
        isError=isError
    )

# Called by mcpx to understand how and why to use this tool.
# Note: Your servlet configs will not be set when this function is called,
# so do not rely on config in this function
def describe() -> ListToolsResult:
    return ListToolsResult(
        [
            ToolDescription(
                name="eval-py",
                description="Evaluate some Python using `exec()` in a sandbox.",
                inputSchema={
                    "type": "object",
                    "properties": {
                        "code": {
                            "type": "string",
                            "description": "The Python code to eval. This code gets passed into `exec()` and the stdout output is returned. Unless it's an image file, only text can be outputted, not binary.",
                        },
                    },
                    "required": ["code"],
                },
            )
        ]
    )

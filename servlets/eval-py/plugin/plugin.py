from typing import Optional, List  # noqa: F401
from datetime import datetime  # noqa: F401
import extism  # noqa: F401 # pyright: ignore
import traceback
from io import StringIO
import sys

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


# Called when the tool is invoked.
# If you support multiple tools, you must switch on the input.params.name to detect which tool is being called.
def call(input: CallToolRequest) -> CallToolResult:
    try:
        code = input.params.arguments['code']
        output = StringIO()
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
        isError = True
    return CallToolResult(
        content=[
            Content(
                text=result,
                mimeType="text/plain",
                type=ContentType.Text,
                annotations=None,
                data=None,
            )
        ],
        isError=isError,
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
                            "description": "The Python code to eval. This code gets passed into `exec()` and the stdout output is returned.",
                        },
                    },
                    "required": ["code"],
                },
            )
        ]
    )

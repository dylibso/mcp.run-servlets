from typing import Optional, List  # noqa: F401
from datetime import datetime  # noqa: F401
import extism  # noqa: F401 # pyright: ignore
from urllib.parse import urlencode

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

class Obsidian:
    def __init__(self):
        self.host = extism.Config.get_str('OBSIDIAN_API_URL')
        api_key = extism.Config.get_str('OBSIDIAN_API_KEY')
        self.headers = {"Authorization": f"Bearer {api_key}"}

    def get(self, path):
        extism.log(extism.LogLevel.Info, f"get {self.host}{path}")
        res = extism.Http.request(f"{self.host}{path}", 'GET', None, self.headers)
        extism.log(extism.LogLevel.Info, f"-> {res.status_code}")
        resres = res.data_str()
        extism.log(extism.LogLevel.Debug, f"-> {resres}")
        return resres
    
    def post(self, path, body = None, extraheaders = {}):
        headers = self.headers | extraheaders
        extism.log(extism.LogLevel.Info, f"post {self.host}{path}")
        res = extism.Http.request(f"{self.host}{path}", 'POST', body, headers)
        extism.log(extism.LogLevel.Info, f"-> {res.status_code}")
        resres = res.data_str()
        extism.log(extism.LogLevel.Debug, f"-> {resres}")
        return resres

    def list_files_in_vault(self):
        return self.get('/vault/')

    def list_files_in_dir(self, path):
        return self.get(f"/vault/{path}/")

    def get_file_contents(self, path):
        return self.get(f"/vault/{path}")
    
    def search(self, query: str, context_length: int = 100):
        params = {
            'query': query,
            'contextLength': context_length
        }
        query_string = urlencode(params)
        return self.post(f"/search/simple/?{query_string}")
    
    def append_content(self, path, content):
        return self.post(f"/vault/{path}", content, {'Content-Type': 'text/markdown'})
    
    def patch_content(filepath, operation, target_type, target, content):
        headers = {
            'Content-Type': 'text/markdown',
            'Operation': operation,
            'Target-Type': target_type,
            'Target': urllib.parse.quote(target),
        }
        return self.post(f"/vault/{path}", content, headers)
    
    def complex_search(self, query):
        headers = {
            'Content-Type': 'application/vnd.olrapi.jsonlogic+json',
        }
        query_string = urlencode(params)
        return self.post(f"/search/", json.dumps(query), headers)

def errorReturn(message):
    return CallToolResult(
                content=[
                    Content(
                        text=message,
                        mimeType="text/plain",
                        type=ContentType.Text,
                        annotations=None,
                        data=None,
                    )
                ],
                isError=True,
            )


# Called when the tool is invoked.
# If you support multiple tools, you must switch on the input.params.name to detect which tool is being called.
def call(input) -> CallToolResult:
    try:
        fname = input['params']['name']
    except:
        raise Exception("params name must be provided")
    obsidian = Obsidian()
    match fname:
        case "list_files_in_vault":
            contentText = obsidian.list_files_in_vault()
        case "list_files_in_dir":
            try:
                dirpath = input['params']['arguments']['dirpath']
            except:
                return errorReturn("Argument dirpath not provided")
            contentText = obsidian.list_files_in_dir(dirpath)
        case "get_file_contents":
            try:
                filepath = input['params']['arguments']['filepath']
            except:
                return errorReturn("Argument filepath not provided")
            contentText = obsidian.get_file_contents(filepath)
        case "simple_search":
            try:
                query = input['params']['arguments']['query']
            except:
                return errorReturn("Argument query not provided")
            context_length = input['params']['arguments'].get('context_length')
            contentText = obsidian.search(query, context_length)
        case "append_content":
            try:
                filepath = input['params']['arguments']['filepath']
                content = input['params']['arguments']['content']
            except:
                return errorReturn("Argument filepath or content not provided")
            contentText = obsidian.append_content(filepath, content)
        case "patch_content":
            try:
                filepath = input['params']['arguments']['filepath']
                operation = input['params']['arguments']['operation']
                target_type = input['params']['arguments']['target_type']
                target = input['params']['arguments']['target']
                content = input['params']['arguments']['content']
            except:
                return errorReturn("Arguments missing")
            contentText = obsidian.patch_content(filepath, operation, target_type, target, content)
        case "complex_search":
            try:
                query = input['params']['arguments']['query']
            except:
                return errorReturn("Argument query not provided")
            contentText = obsidian.complex_search(query)
        case _:
            return errorReturn(f"Unknown tool {fname}")
    return CallToolResult(
        content=[
            Content(
                text=contentText,
                mimeType="text/plain",
                type=ContentType.Text,
                annotations=None,
                data=None,
            )
        ],
        isError=False,
    )


# Called by mcpx to understand how and why to use this tool.
# Note: Your servlet configs will not be set when this function is called,
# so do not rely on config in this function
def describe() -> ListToolsResult:
    return ListToolsResult(
        [
            ToolDescription(
                name="list_files_in_vault",
                description="Lists all files and directories in the root directory of your Obsidian vault.",
                inputSchema={
                    "type": "object",
                },
            ),
            ToolDescription(
                name="list_files_in_dir",
                description="Lists all files and directories that exist in a specific Obsidian directory.",
                inputSchema={
                    "type": "object",
                    "properties": {
                        "dirpath": {
                            "type": "string",
                            "description": "Path to list files from (relative to your vault root). Note that empty directories will not be returned."
                        },
                    },
                    "required": ["dirpath"]
                },
            ),
            ToolDescription(
                name="get_file_contents",
                description="Return the content of a single file in your vault.",
                inputSchema={
                    "type": "object",
                    "properties": {
                        "filepath": {
                            "type": "string",
                            "description": "Path to the relevant file (relative to your vault root).",
                            "format": "path"
                        },
                    },
                    "required": ["filepath"]
                },
            ),
            ToolDescription(
                name="simple_search",
                description="""Simple search for documents matching a specified text query across all files in the vault. 
                Use this tool when you want to do a simple text search""",
                inputSchema={
                    "type": "object",
                    "properties": {
                        "query": {
                            "type": "string",
                            "description": "Text to a simple search for in the vault."
                        },
                        "context_length": {
                            "type": "integer",
                            "description": "How much context to return around the matching string (default: 100)",
                            "default": 100
                        }
                    },
                    "required": ["query"]
                },
            ),
            ToolDescription(
                name="append_content",
                description="Append content to a new or existing file in the vault.",
                inputSchema={
                    "type": "object",
                    "properties": {
                        "filepath": {
                            "type": "string",
                            "description": "Path to the file (relative to vault root)",
                            "format": "path"
                        },
                        "content": {
                            "type": "string",
                            "description": "Content to append to the file"
                        }
                    },
                    "required": ["filepath", "content"]
                },
            ),
            ToolDescription(
                name="patch_content",
                description="Insert content into an existing note relative to a heading, block reference, or frontmatter field.",
                inputSchema={
                    "type": "object",
                    "properties": {
                        "filepath": {
                            "type": "string",
                            "description": "Path to the file (relative to vault root)",
                            "format": "path"
                        },
                        "operation": {
                            "type": "string",
                            "description": "Operation to perform (append, prepend, or replace)",
                            "enum": ["append", "prepend", "replace"]
                        },
                        "target_type": {
                            "type": "string",
                            "description": "Type of target to patch",
                            "enum": ["heading", "block", "frontmatter"]
                        },
                        "target": {
                            "type": "string", 
                            "description": "Target identifier (heading path, block reference, or frontmatter field)"
                        },
                        "content": {
                            "type": "string",
                            "description": "Content to insert"
                        }
                    },
                    "required": ["filepath", "operation", "target_type", "target", "content"]
                },
            ),
            ToolDescription(
                name="complex_search",
                description="""Complex search for documents using a JsonLogic query. 
                Supports standard JsonLogic operators plus 'glob' and 'regexp' for pattern matching. Results must be non-falsy.

                Use this tool when you want to do a complex search, e.g. for all documents with certain tags etc.""",
                inputSchema={
                    "type": "object",
                    "properties": {
                        "query": {
                            "type": "object",
                            "description": "JsonLogic query object. Example: {\"glob\": [\"*.md\", {\"var\": \"path\"}]} matches all markdown files"
                        }
                    },
                    "required": ["query"]
                },
            ),
        ]
    )

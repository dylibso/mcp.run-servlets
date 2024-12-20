# THIS FILE WAS GENERATED BY `xtp-python-bindgen`. DO NOT EDIT.

from typing import Optional, List  # noqa: F401
from datetime import datetime  # noqa: F401
import extism  # pyright: ignore
import plugin
import json

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


# Imports

# Exports
# The implementations for these functions is in `plugin.py`


# Called when the tool is invoked.
# If you support multiple tools, you must switch on the input.params.name to detect which tool is being called.
@extism.plugin_fn
def call():
    data = json.loads(extism.input_str())
    res = plugin.call(data)
    extism.output(res)


# Called by mcpx to understand how and why to use this tool.
# Note: Your servlet configs will not be set when this function is called,
# so do not rely on config in this function
@extism.plugin_fn
def describe():
    res = plugin.describe()
    extism.output(res)
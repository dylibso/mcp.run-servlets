# THIS FILE WAS GENERATED BY `xtp-python-bindgen`. DO NOT EDIT.

from __future__ import annotations
from enum import Enum  # noqa: F401
from typing import Optional, List  # noqa: F401
from datetime import datetime  # noqa: F401
from dataclasses import dataclass  # noqa: F401
from dataclass_wizard import JSONWizard, skip_if_field, IS  # noqa: F401
from dataclass_wizard.type_def import JSONObject
from base64 import b64encode, b64decode


@dataclass
class BlobResourceContents(JSONWizard):
    # A base64-encoded string representing the binary data of the item.
    blob: str

    # The URI of this resource.
    uri: str

    # The MIME type of this resource, if known.
    mimeType: Optional[str] = skip_if_field(IS(None), default=None)

    @classmethod
    def _pre_from_dict(cls, o: JSONObject) -> JSONObject:
        return o

    def _pre_dict(self):
        return


@dataclass
class CallToolRequest(JSONWizard):
    params: Params

    method: Optional[str] = skip_if_field(IS(None), default=None)

    @classmethod
    def _pre_from_dict(cls, o: JSONObject) -> JSONObject:
        return o

    def _pre_dict(self):
        return


@dataclass
class CallToolResult(JSONWizard):
    content: List[Content]

    # Whether the tool call ended in an error.
    #
    # If not set, this is assumed to be false (the call was successful).
    isError: Optional[bool] = skip_if_field(IS(None), default=None)

    @classmethod
    def _pre_from_dict(cls, o: JSONObject) -> JSONObject:
        return o

    def _pre_dict(self):
        return


@dataclass
class Content(JSONWizard):
    type: ContentType

    annotations: Optional[TextAnnotation] = skip_if_field(IS(None), default=None)

    # The base64-encoded image data.
    data: Optional[str] = skip_if_field(IS(None), default=None)

    # The MIME type of the image. Different providers may support different image types.
    mimeType: Optional[str] = skip_if_field(IS(None), default=None)

    # The text content of the message.
    text: Optional[str] = skip_if_field(IS(None), default=None)

    @classmethod
    def _pre_from_dict(cls, o: JSONObject) -> JSONObject:
        return o

    def _pre_dict(self):
        return


class ContentType(Enum):
    Text = "text"
    Image = "image"
    Resource = "resource"


@dataclass
class ListToolsResult(JSONWizard):
    # The list of ToolDescription objects provided by this servlet.
    tools: List[ToolDescription]

    @classmethod
    def _pre_from_dict(cls, o: JSONObject) -> JSONObject:
        return o

    def _pre_dict(self):
        return


@dataclass
class Params(JSONWizard):
    name: str

    arguments: Optional[dict] = skip_if_field(IS(None), default=None)

    @classmethod
    def _pre_from_dict(cls, o: JSONObject) -> JSONObject:
        return o

    def _pre_dict(self):
        return


class Role(Enum):
    Assistant = "assistant"
    User = "user"


@dataclass
class TextAnnotation(JSONWizard):
    # Describes who the intended customer of this object or data is.
    #
    # It can include multiple entries to indicate content useful for multiple audiences (e.g., `["user", "assistant"]`).
    audience: Optional[List[Role]] = skip_if_field(IS(None), default=None)

    # Describes how important this data is for operating the server.
    #
    # A value of 1 means "most important," and indicates that the data is
    # effectively required, while 0 means "least important," and indicates that
    # the data is entirely optional.
    priority: Optional[float] = skip_if_field(IS(None), default=None)

    @classmethod
    def _pre_from_dict(cls, o: JSONObject) -> JSONObject:
        return o

    def _pre_dict(self):
        return


@dataclass
class TextResourceContents(JSONWizard):
    # The text of the item. This must only be set if the item can actually be represented as text (not binary data).
    text: str

    # The URI of this resource.
    uri: str

    # The MIME type of this resource, if known.
    mimeType: Optional[str] = skip_if_field(IS(None), default=None)

    @classmethod
    def _pre_from_dict(cls, o: JSONObject) -> JSONObject:
        return o

    def _pre_dict(self):
        return


@dataclass
class ToolDescription(JSONWizard):
    # A description of the tool
    description: str

    # The JSON schema describing the argument input
    inputSchema: dict

    # The name of the tool. It should match the plugin / binding name.
    name: str

    @classmethod
    def _pre_from_dict(cls, o: JSONObject) -> JSONObject:
        return o

    def _pre_dict(self):
        return

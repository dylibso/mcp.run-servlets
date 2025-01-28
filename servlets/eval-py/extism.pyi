from typing import Any, Callable, Dict, List, Optional, Type, TypeVar, Union, overload
from enum import Enum

class LogLevel(Enum):
    """Log levels for Extism plugin logging."""
    Trace = 0
    Debug = 1
    Info = 2
    Warn = 3
    Error = 4

class MemoryHandle:
    """Represents a handle to memory in the Extism runtime."""
    offset: int
    length: int
    
    def __init__(self, offset: int, length: int) -> None: ...

class Memory:
    """Memory management utilities."""
    @staticmethod
    def find(offset: int) -> Optional[MemoryHandle]: ...
    
    @staticmethod
    def bytes(handle: MemoryHandle) -> bytes: ...
    
    @staticmethod
    def string(handle: MemoryHandle) -> str: ...
    
    @staticmethod
    def free(handle: MemoryHandle) -> None: ...
    
    @staticmethod
    def alloc(data: bytes) -> MemoryHandle: ...

class HttpRequest:
    """Represents an HTTP request."""
    url: str
    method: Optional[str]
    headers: Optional[Dict[str, str]]
    
    def __init__(
        self,
        url: str,
        method: Optional[str] = None,
        headers: Optional[Dict[str, str]] = None
    ) -> None: ...

class HttpResponse:
    """Represents an HTTP response."""
    _inner: Any  # Internal representation
    
    def status_code(self) -> int: ...
    def data_bytes(self) -> bytes: ...
    def data_str(self) -> str: ...
    def data_json(self) -> Any: ...
    def headers(self) -> Dict[str, str]: ...

class Http:
    """HTTP utility functions."""
    @staticmethod
    def request(
        url: str,
        meth: str = "GET",
        body: Optional[Union[bytes, str]] = None,
        headers: Optional[Dict[str, str]] = None
    ) -> HttpResponse: ...

class Var:
    """Variable storage utilities."""
    @staticmethod
    def get_bytes(key: str) -> Optional[bytes]: ...
    
    @staticmethod
    def get_str(key: str) -> Optional[str]: ...
    
    @staticmethod
    def get_json(key: str) -> Any: ...
    
    @staticmethod
    def set(key: str, value: Union[bytes, str]) -> None: ...

class Config:
    """Configuration utilities."""
    @staticmethod
    def get_str(key: str) -> Optional[str]: ...
    
    @staticmethod
    def get_json(key: str) -> Any: ...

def plugin_fn(func: Callable) -> Callable:
    """Decorator for functions that will be called by Extism."""
    ...

def shared_fn(func: Callable) -> Callable:
    """Decorator for exports that won't be called directly by Extism."""
    ...

def import_fn(module: str, name: str) -> Callable:
    """Decorator for importing functions from the host."""
    ...

def input_str() -> str:
    """Get input as string."""
    ...

def input_bytes() -> bytes:
    """Get input as bytes."""
    ...

def output_str(result: str) -> None:
    """Set string output."""
    ...

def output_bytes(result: bytes) -> None:
    """Set bytes output."""
    ...

def input_json(t: Optional[Type] = None) -> Any:
    """Get input as JSON."""
    ...

def output_json(x: Any) -> None:
    """Set JSON output."""
    ...

T = TypeVar('T')
def input(t: Type[T] = None) -> Optional[T]:
    """Get typed input."""
    ...

def output(x: Any = None) -> None:
    """Set typed output."""
    ...

def log(level: LogLevel, msg: Union[str, bytes]) -> None:
    """Log a message at the specified level."""
    ...

def set_error(msg: str) -> None:
    """Set an error message."""
    ...
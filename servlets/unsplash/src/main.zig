const schema = @import("schema.zig");
const extism = @import("extism-pdk");
const http = extism.http;
const plugin = @import("pdk.zig")._plugin;
const std = @import("std");
const json = std.json;
const eql = std.mem.eql;
const allocPrint = std.fmt.allocPrint;
const allocator = std.heap.wasm_allocator;

/// Called when the tool is invoked.
/// If you support multiple tools, you must switch on the input.params.name to detect which tool is being called.
/// It takes CallToolRequest as input (The incoming tool request from the LLM)
/// And returns CallToolResult (The servlet's response to the given tool call)
pub fn call(input: schema.CallToolRequest) !schema.CallToolResult {
    const name = input.params.name;
    if (eql(u8, name, "photos")) {
        return getPhotos(input.params.arguments);
    } else if (eql(u8, name, "photos_id")) {
        return getPhotosId(input.params.arguments);
    } else if (eql(u8, name, "search_photos")) {
        return searchPhotos(input.params.arguments);
    }
    return error.PluginFunctionNotImplemented;
}

fn getPhotos(args: ?json.ArrayHashMap(std.json.Value)) !schema.CallToolResult {
    const apiKey = try plugin.getConfig("API_KEY") orelse return error.MissingConfig;
    const page = args.?.map.get("page") orelse json.Value{ .integer = 1 };
    const per_page = args.?.map.get("per_page") orelse json.Value{ .integer = 10 };
    const url = try allocPrint(
        allocator,
        "https://api.unsplash.com/photos?page={d}&per_page={d}",
        .{ page.integer, per_page.integer },
    );
    var req = http.HttpRequest.init("GET", url);
    defer req.deinit(allocator);
    try req.setHeader(
        allocator,
        "Authorization",
        try allocPrint(allocator, "Client-ID {s}", .{apiKey}),
    );
    const resp = try plugin.request(req, null);
    const body = try resp.body(allocator);
    if (resp.status != 200) {
        return callToolError(try allocPrint(
            allocator,
            "Error {d}: {s}",
            .{ resp.status, body },
        ));
    }
    return schema.CallToolResult{
        .content = try allocator.dupe(schema.Content, &.{.{
            .type = schema.ContentType.text,
            .text = body,
        }}),
    };
}

fn getPhotosId(arguments: ?json.ArrayHashMap(std.json.Value)) !schema.CallToolResult {
    const apiKey = try plugin.getConfig("API_KEY") orelse return error.MissingConfig;
    const args = arguments orelse return callToolError("missing arguments");
    const id = args.map.get("id") orelse return error.MissingArgument;
    const url = try allocPrint(
        allocator,
        "https://api.unsplash.com/photos/{s}",
        .{id.string},
    );
    var req = http.HttpRequest.init("GET", url);
    defer req.deinit(allocator);
    try req.setHeader(
        allocator,
        "Authorization",
        try allocPrint(allocator, "Client-ID {s}", .{apiKey}),
    );
    const resp = try plugin.request(req, null);
    const body = try resp.body(allocator);
    if (resp.status != 200) {
        return callToolError(try allocPrint(
            allocator,
            "Error {d}: {s}",
            .{ resp.status, body },
        ));
    }
    return schema.CallToolResult{
        .content = try allocator.dupe(schema.Content, &.{.{
            .type = schema.ContentType.text,
            .text = body,
        }}),
    };
}

// query Search terms.
// page	Page number to retrieve. (Optional; default: 1)
// per_page	Number of items per page. (Optional; default: 10)
// order_by	How to sort the photos. (Optional; default: relevant). Valid values are latest and relevant.
// collections	Collection ID(‘s) to narrow search. Optional. If multiple, comma-separated.
// content_filter	Limit results by content safety. (Optional; default: low). Valid values are low and high.
// color	Filter results by color. Optional. Valid values are: black_and_white, black, white, yellow, orange, red, purple, magenta, green, teal, and blue.
// orientation	Filter by photo orientation. Optional. (Valid values: landscape, portrait, squarish)

fn searchPhotos(arguments: ?json.ArrayHashMap(std.json.Value)) !schema.CallToolResult {
    const apiKey = try plugin.getConfig("API_KEY") orelse return error.MissingConfig;
    const args = arguments orelse return callToolError("missing arguments");
    const query = args.map.get("query") orelse return callToolError("missing query");
    const page = args.map.get("page") orelse json.Value{ .integer = 1 };
    const per_page = args.map.get("per_page") orelse json.Value{ .integer = 10 };
    const order_by = args.map.get("order_by") orelse json.Value{ .string = "relevant" };
    const content_filter = args.map.get("content_filter") orelse json.Value{ .string = "low" };
    var color: []u8 = "";
    if (args.map.get("color") != null) {
        const c = args.map.get("color").?.string;
        color = try allocPrint(allocator, "&color={s}", .{c});
    }
    var orientation: []u8 = "";
    if (args.map.get("orientation") != null) {
        const c = args.map.get("orientation").?.string;
        orientation = try allocPrint(allocator, "&orientation={s}", .{c});
    }
    const url = try allocPrint(
        allocator,
        "https://api.unsplash.com/search/photos?page={d}&per_page={d}&order_by={s}&content_filter={s}{s}{s}&query={s}",
        .{ page.integer, per_page.integer, order_by.string, content_filter.string, color, orientation, query.string },
    );
    var req = http.HttpRequest.init("GET", url);
    defer req.deinit(allocator);
    try req.setHeader(
        allocator,
        "Authorization",
        try allocPrint(allocator, "Client-ID {s}", .{apiKey}),
    );
    const resp = try plugin.request(req, null);
    const body = try resp.body(allocator);
    if (resp.status != 200) {
        return callToolError(try allocPrint(
            allocator,
            "Error {d}: {s}",
            .{ resp.status, body },
        ));
    }
    return schema.CallToolResult{
        .content = try allocator.dupe(schema.Content, &.{.{
            .type = schema.ContentType.text,
            .text = body,
        }}),
    };
}

fn callToolError(msg: []const u8) !schema.CallToolResult {
    return schema.CallToolResult{
        .isError = true,
        .content = try allocator.dupe(schema.Content, &.{.{
            .type = schema.ContentType.text,
            .text = msg,
        }}),
    };
}

pub fn describe() !schema.ListToolsResult {
    const tools = @embedFile("tools.json");
    const r = try json.parseFromSlice(schema.ListToolsResult, allocator, tools, .{});
    return r.value;
}

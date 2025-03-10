const std = @import("std");

// if you're not using WASI, change this to .freestanding
const BUILD_TARGET_OS = .wasi;

pub fn build(b: *std.Build) !void {
    const optimize = b.standardOptimizeOption(.{});
    const target = b.standardTargetOptions(.{
        .default_target = .{ .abi = .musl, .os_tag = BUILD_TARGET_OS, .cpu_arch = .wasm32 },
    });
    const pdk_module = b.dependency("extism-pdk", .{ .target = target, .optimize = optimize }).module("extism-pdk");

    var plugin = b.addExecutable(.{
        .name = "plugin",
        .root_source_file = b.path("src/pdk.zig"),
        .target = target,
        .optimize = optimize,
    });
    if (BUILD_TARGET_OS == .wasi) {
        plugin.wasi_exec_model = .reactor;
    }
    plugin.rdynamic = true;
    plugin.entry = .disabled; // or add an empty `pub fn main() void {}` to your code
    plugin.root_module.addImport("extism-pdk", pdk_module);
    b.installArtifact(plugin);
}

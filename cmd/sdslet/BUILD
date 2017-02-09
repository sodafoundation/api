package(default_visibility = ["//visibility:public"])

licenses(["notice"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_binary",
    "go_library",
)

go_binary(
    name = "sdslet",
    library = ":go_default_library",
    tags = ["automanaged"],
)

go_library(
    name = "go_default_library",
    srcs = ["sdslet.go"],
    tags = ["automanaged"],
    deps = [
        "//pkg/grpc/server:go_default_library",
		"//vendor:openstack/golang-client/volume",
		"//vendor:openstack/golang-client/share",
    ],
)

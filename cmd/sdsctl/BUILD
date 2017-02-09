package(default_visibility = ["//visibility:public"])

licenses(["notice"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_binary",
    "go_library",
)

go_binary(
    name = "sdsctl",
    library = ":go_default_library",
    tags = ["automanaged"],
)

go_library(
    name = "go_default_library",
    srcs = ["sdsctl.go"],
    tags = ["automanaged"],
    deps = [
        "//pkg/pkg/api:go_default_library",
		"//pkg/pkg/api/volumes:go_default_library",
		"//pkg/pkg/api/shares:go_default_library",
		"//github.com/spf13/cobra",
    ],
)

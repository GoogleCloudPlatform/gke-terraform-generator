# Copyright 2018 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     https://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

#
# TODO do a toolchain for this.  At this point this is not working.
#

#http_archive(
#    name = "terraform",
#    urls = ["https://releases.hashicorp.com/terraform/0.11.14/terraform_0.11.14_darwin_amd64.zip"],
#    sha256 = "829bdba148afbd61eab4aafbc6087838f0333d8876624fe2ebc023920cfc2ad5",
#    build_file_content = """
#filegroup(
#  name = "binary",
#  srcs = ["terraform"],
#  visibility = ["//visibility:public"],
#)
#        """
#)

http_archive(
    name = "io_bazel_rules_go",
    urls = ["https://github.com/bazelbuild/rules_go/archive/bf51fd0698feb9b495bc7ae17f77782138746455.tar.gz"],
    sha256 = "6eb05605adb5687bd45d960a153b3ab4052c239f1fd357173dac58241a5df17b",
    strip_prefix = "rules_go-bf51fd0698feb9b495bc7ae17f77782138746455",
)

http_archive(
    name = "bazel_gazelle",
    urls = ["https://github.com/bazelbuild/bazel-gazelle/releases/download/0.17.0/bazel-gazelle-0.17.0.tar.gz"],
    sha256 = "3c681998538231a2d24d0c07ed5a7658cb72bfb5fd4bf9911157c0e9ac6a2687",
)

load("@io_bazel_rules_go//go:deps.bzl", "go_rules_dependencies", "go_register_toolchains")

go_rules_dependencies()

go_register_toolchains(nogo = "@//pkg/analyzer:my_nogo")

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies", "go_repository")

gazelle_dependencies()

go_repository(
    name = "com_github_armon_consul_api",
    commit = "eb2c6b5be1b6",
    importpath = "github.com/armon/consul-api",
)

go_repository(
    name = "com_github_burntsushi_toml",
    importpath = "github.com/BurntSushi/toml",
    tag = "v0.3.1",
)

go_repository(
    name = "com_github_coreos_etcd",
    importpath = "github.com/coreos/etcd",
    tag = "v3.3.10",
)

go_repository(
    name = "com_github_coreos_go_etcd",
    importpath = "github.com/coreos/go-etcd",
    tag = "v2.0.0",
)

go_repository(
    name = "com_github_coreos_go_semver",
    importpath = "github.com/coreos/go-semver",
    tag = "v0.2.0",
)

go_repository(
    name = "com_github_cpuguy83_go_md2man",
    importpath = "github.com/cpuguy83/go-md2man",
    tag = "v1.0.10",
)

go_repository(
    name = "com_github_creasty_defaults",
    importpath = "github.com/creasty/defaults",
    tag = "v1.3.0",
)

go_repository(
    name = "com_github_davecgh_go_spew",
    importpath = "github.com/davecgh/go-spew",
    tag = "v1.1.1",
)

go_repository(
    name = "com_github_fsnotify_fsnotify",
    importpath = "github.com/fsnotify/fsnotify",
    tag = "v1.4.7",
)

go_repository(
    name = "com_github_go_playground_locales",
    importpath = "github.com/go-playground/locales",
    tag = "v0.12.1",
)

go_repository(
    name = "com_github_go_playground_universal_translator",
    importpath = "github.com/go-playground/universal-translator",
    tag = "v0.16.0",
)

go_repository(
    name = "com_github_hashicorp_hcl",
    importpath = "github.com/hashicorp/hcl",
    tag = "v1.0.0",
)

go_repository(
    name = "com_github_imdario_mergo",
    importpath = "github.com/imdario/mergo",
    tag = "v0.3.7",
)

go_repository(
    name = "com_github_inconshreveable_mousetrap",
    importpath = "github.com/inconshreveable/mousetrap",
    tag = "v1.0.0",
)

go_repository(
    name = "com_github_leodido_go_urn",
    importpath = "github.com/leodido/go-urn",
    tag = "v1.1.0",
)

go_repository(
    name = "com_github_magiconair_properties",
    importpath = "github.com/magiconair/properties",
    tag = "v1.8.0",
)

go_repository(
    name = "com_github_mitchellh_go_homedir",
    importpath = "github.com/mitchellh/go-homedir",
    tag = "v1.1.0",
)

go_repository(
    name = "com_github_mitchellh_mapstructure",
    importpath = "github.com/mitchellh/mapstructure",
    tag = "v1.1.2",
)

go_repository(
    name = "com_github_pelletier_go_toml",
    importpath = "github.com/pelletier/go-toml",
    tag = "v1.2.0",
)

go_repository(
    name = "com_github_pmezard_go_difflib",
    importpath = "github.com/pmezard/go-difflib",
    tag = "v1.0.0",
)

go_repository(
    name = "com_github_russross_blackfriday",
    importpath = "github.com/russross/blackfriday",
    tag = "v1.5.2",
)

go_repository(
    name = "com_github_spf13_afero",
    importpath = "github.com/spf13/afero",
    tag = "v1.1.2",
)

go_repository(
    name = "com_github_spf13_cast",
    importpath = "github.com/spf13/cast",
    tag = "v1.3.0",
)

go_repository(
    name = "com_github_spf13_cobra",
    importpath = "github.com/spf13/cobra",
    tag = "v0.0.4",
)

go_repository(
    name = "com_github_spf13_jwalterweatherman",
    importpath = "github.com/spf13/jwalterweatherman",
    tag = "v1.0.0",
)

go_repository(
    name = "com_github_spf13_pflag",
    importpath = "github.com/spf13/pflag",
    tag = "v1.0.3",
)

go_repository(
    name = "com_github_spf13_viper",
    importpath = "github.com/spf13/viper",
    tag = "v1.3.2",
)

go_repository(
    name = "com_github_stretchr_testify",
    importpath = "github.com/stretchr/testify",
    tag = "v1.2.2",
)

go_repository(
    name = "com_github_ugorji_go_codec",
    commit = "d75b2dcb6bc8",
    importpath = "github.com/ugorji/go/codec",
)

go_repository(
    name = "com_github_xordataexchange_crypt",
    commit = "b2862e3d0a77",
    importpath = "github.com/xordataexchange/crypt",
)

go_repository(
    name = "in_gopkg_check_v1",
    commit = "20d25e280405",
    importpath = "gopkg.in/check.v1",
)

go_repository(
    name = "in_gopkg_go_playground_assert_v1",
    importpath = "gopkg.in/go-playground/assert.v1",
    tag = "v1.2.1",
)

go_repository(
    name = "in_gopkg_go_playground_validator_v9",
    importpath = "gopkg.in/go-playground/validator.v9",
    tag = "v9.29.0",
)

go_repository(
    name = "in_gopkg_yaml_v2",
    importpath = "gopkg.in/yaml.v2",
    tag = "v2.2.2",
)

go_repository(
    name = "io_k8s_klog",
    importpath = "k8s.io/klog",
    tag = "v0.3.3",
)

go_repository(
    name = "org_golang_x_crypto",
    commit = "c2843e01d9a2",
    importpath = "golang.org/x/crypto",
)

go_repository(
    name = "org_golang_x_sys",
    commit = "d0b11bdaac8a",
    importpath = "golang.org/x/sys",
)

go_repository(
    name = "org_golang_x_text",
    importpath = "golang.org/x/text",
    tag = "v0.3.0",
)

go_repository(
    name = "org_golang_x_net",
    commit = "3b0461eec859",
    importpath = "golang.org/x/net",
)

go_repository(
    name = "org_golang_x_sync",
    commit = "112230192c58",
    importpath = "golang.org/x/sync",
)

go_repository(
    name = "org_golang_x_tools",
    commit = "c8855242db9c",
    importpath = "golang.org/x/tools",
)

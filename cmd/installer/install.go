//go:build mage

package main

const (
	// https://golang.org/dl/
	goVersion = "1.17.6"

	// https://github.com/protocolbuffers/protobuf/releases
	protocVersion = "3.19.3"

	// https://repo1.maven.org/maven2/io/grpc/protoc-gen-grpc-java/
	protocGenGRPCJavaVersion = "1.43.2"

	// https://github.com/nvm-sh/nvm/releases
	nvmVersion = "0.39.1"

	// https://www.azul.com/downloads/zulu-community/?version=java-11-lts&os=ubuntu&architecture=x86-64-bit&package=jdk
	zuluVersion    = "11.52.13"
	zuluJDKVersion = "11.0.13"

	// https://maven.apache.org/download.cgi
	mavenVersion = "3.8.4"

	// https://github.com/restic/restic/releases
	resticVersion = "0.12.1"

	// https://github.com/jeffkaufman/icdiff/tags
	icdiffVersion = "2.0.4"
)

var (
	aptPackageNames = []string{
		"apt-utils", "locales", "wget", "less", "vim", "nano", "zip", "unzip", "xz-utils", "htop", "gcc", "make",
		"telnet", "netcat", "socat", "docker.io", "libfontconfig", "postgresql-client", "iputils-ping", "libxml2-utils",
		"curl", "git", "ca-certificates",
	}

	goToolModules = []string{
		"github.com/uudashr/gopkgs/v2/cmd/gopkgs@latest",
		"github.com/ramya-rao-a/go-outline@latest",
		"github.com/cweill/gotests/gotests@latest",
		"github.com/fatih/gomodifytags@latest",
		"github.com/josharian/impl@latest",
		"github.com/haya14busa/goplay/cmd/goplay@latest",
		"github.com/go-delve/delve/cmd/dlv@latest",
		"github.com/golangci/golangci-lint/cmd/golangci-lint@latest",
		"golang.org/x/tools/gopls@latest",
		"github.com/dvyukov/go-fuzz/go-fuzz@latest",
		"github.com/dvyukov/go-fuzz/go-fuzz-build@latest",
		"golang.org/x/perf/cmd/benchstat@latest",
		"github.com/orijtech/structslop/cmd/structslop@latest",
		"github.com/blizzy78/textsimilarity/cmd/textsimilarity@latest",
	}

	nodeLTSNames = []string{"gallium", "dubnium"}

	npmPackageNames = []string{
		"serve",
		"gatsby",
	}

	volumeFolders = []string{
		".vscode-server/extensions",
		"workspaces",
		".m2",
		"sophora-repo",
		".ssh",
		"restic-repos",
	}

	protocGoModuleURLs = []string{
		"google.golang.org/protobuf/cmd/protoc-gen-go",
		"google.golang.org/grpc/cmd/protoc-gen-go-grpc",
	}
)

const tz = "Europe/Berlin"

const (
	goURL = "https://golang.org/dl/go" + goVersion + ".linux-amd64.tar.gz"

	protocURL            = "https://github.com/protocolbuffers/protobuf/releases/download/v" + protocVersion + "/protoc-" + protocVersion + "-linux-x86_64.zip"
	protocGenGRPCJavaURL = "https://repo1.maven.org/maven2/io/grpc/protoc-gen-grpc-java/" + protocGenGRPCJavaVersion + "/protoc-gen-grpc-java-" + protocGenGRPCJavaVersion + "-linux-x86_64.exe"

	nvmInstallURL = "https://raw.githubusercontent.com/nvm-sh/nvm/v" + nvmVersion + "/install.sh"
	npmInstallURL = "https://www.npmjs.com/install.sh"

	zuluJDKURL = "https://cdn.azul.com/zulu/bin/zulu" + zuluVersion + "-ca-jdk" + zuluJDKVersion + "-linux_x64.tar.gz"

	mavenURL = "https://dlcdn.apache.org/maven/maven-3/" + mavenVersion + "/binaries/apache-maven-" + mavenVersion + "-bin.tar.gz"

	resticURL = "https://github.com/restic/restic/releases/download/v" + resticVersion + "/restic_" + resticVersion + "_linux_amd64.bz2"

	icdiffURL = "https://github.com/jeffkaufman/icdiff/archive/refs/tags/release-" + icdiffVersion + ".tar.gz"
)

//go:build mage

package main

const (
	// https://github.com/protocolbuffers/protobuf/releases
	protocVersion = "3.19.4"

	// https://repo1.maven.org/maven2/io/grpc/protoc-gen-grpc-java/
	protocGenGRPCJavaVersion = "1.45.1"

	// https://github.com/nvm-sh/nvm/releases
	nvmVersion = "0.39.1"

	// https://www.azul.com/downloads/?version=java-11-lts&os=linux&architecture=x86-64-bit&package=jdk
	zuluVersion    = "11.54.25"
	zuluJDKVersion = "11.0.14.1"

	// https://maven.apache.org/download.cgi
	mavenVersion = "3.8.5"

	// https://github.com/restic/restic/releases
	resticVersion = "0.13.0"

	// https://github.com/jeffkaufman/icdiff/tags
	icdiffVersion = "2.0.4"
)

var (
	// https://golang.org/dl/
	goVersions = []string{"1.18", "1.16.15"} // first is default

	pacmanPackageNames = []string{
		"which", "wget", "vim", "nano", "zip", "unzip", "htop", "gcc", "make", "gnu-netcat", "socat", "docker", "fontconfig",
		"postgresql", "git", "graphviz", "inetutils", "openssh", "man-db", "man-pages", "diffutils", "bash-completion",
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
		"golang.org/x/perf/cmd/benchstat@latest",
		"github.com/orijtech/structslop/cmd/structslop@latest",
		"github.com/blizzy78/textsimilarity/cmd/textsimilarity@latest",
	}

	nodeLTSNames = []string{"gallium", "dubnium"} // first is default

	npmPackageNames = []string{"serve", "gatsby"}

	volumeFolders = []string{
		".vscode-server/extensions",
		"workspaces",
		".bash_history_dir",
		".bashrc_dir",
		".gitconfig_dir",
		".m2",
		"sophora-repo",
		".ssh",
		"restic-repos",
		".containerrunner",
		".cache/go-build/fuzz",
	}

	protocGoModuleURLs = []string{
		"google.golang.org/protobuf/cmd/protoc-gen-go@latest",
		"google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest",
	}

	gitCompletionURL = "https://raw.githubusercontent.com/git/git/master/contrib/completion/git-completion.bash"

	gitPromptURL    = "https://raw.githubusercontent.com/git/git/master/contrib/completion/git-prompt.sh"
	gitPromptBashRC = `export PS1="\[\033[30;44m\] \w \[\033[00m\]\[\033[30;42m\]\$(__git_ps1 \" %s \")\[\033[00m\] "` + "\n"
)

const tz = "Europe/Berlin"

const (
	protocURL            = "https://github.com/protocolbuffers/protobuf/releases/download/v" + protocVersion + "/protoc-" + protocVersion + "-linux-x86_64.zip"
	protocGenGRPCJavaURL = "https://repo1.maven.org/maven2/io/grpc/protoc-gen-grpc-java/" + protocGenGRPCJavaVersion + "/protoc-gen-grpc-java-" + protocGenGRPCJavaVersion + "-linux-x86_64.exe"

	nvmInstallURL = "https://raw.githubusercontent.com/nvm-sh/nvm/v" + nvmVersion + "/install.sh"
	npmInstallURL = "https://www.npmjs.com/install.sh"

	zuluJDKURL = "https://cdn.azul.com/zulu/bin/zulu" + zuluVersion + "-ca-jdk" + zuluJDKVersion + "-linux_x64.tar.gz"

	mavenURL = "https://dlcdn.apache.org/maven/maven-3/" + mavenVersion + "/binaries/apache-maven-" + mavenVersion + "-bin.tar.gz"

	resticURL = "https://github.com/restic/restic/releases/download/v" + resticVersion + "/restic_" + resticVersion + "_linux_amd64.bz2"

	icdiffURL = "https://github.com/jeffkaufman/icdiff/archive/refs/tags/release-" + icdiffVersion + ".tar.gz"
)

var goURL = "https://golang.org/dl/go" + goVersions[0] + ".linux-amd64.tar.gz"

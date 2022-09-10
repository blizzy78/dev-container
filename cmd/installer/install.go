//go:build mage

package main

const (
	// https://github.com/protocolbuffers/protobuf/releases
	protocVersion = "21.5"

	// https://repo1.maven.org/maven2/io/grpc/protoc-gen-grpc-java/
	protocGenGRPCJavaVersion = "1.49.0"

	// https://www.azul.com/downloads/?version=java-11-lts&os=linux&architecture=x86-64-bit&package=jdk
	zuluVersion    = "11.58.23"
	zuluJDKVersion = "11.0.16.1"

	// https://maven.apache.org/download.cgi
	mavenVersion = "3.8.6"
)

var (
	// https://golang.org/dl/
	goVersions = []string{"1.19.1", "1.16.15"} // first is default

	pacmanPackageNames = []string{
		"which", "wget", "vim", "nano", "zip", "unzip", "htop", "gcc", "make", "gnu-netcat", "socat", "docker", "fontconfig",
		"postgresql", "git", "graphviz", "inetutils", "openssh", "man-db", "man-pages", "diffutils", "bash-completion", "fakeroot",
		"restic", "dnsutils",

		// dependencies for Chromium in react-snap
		"libxcomposite", "libxcursor", "libxdamage", "libxi", "libxtst", "libxss", "libxrandr", "alsa-lib", "atk", "at-spi2-atk", "gtk3", "nss",
	}

	yayPackageNames = []string{
		"icdiff", "nvm-git", "mage-bin",
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
		"github.com/blizzy78/textsimilarity/cmd/textsimilarity@latest",
	}

	nodeLTSNames = []string{"gallium", "dubnium"} // first is default

	npmPackageNames = []string{"serve"}

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

	npmInstallURL = "https://www.npmjs.com/install.sh"

	zuluJDKURL = "https://cdn.azul.com/zulu/bin/zulu" + zuluVersion + "-ca-jdk" + zuluJDKVersion + "-linux_x64.tar.gz"

	mavenURL = "https://dlcdn.apache.org/maven/maven-3/" + mavenVersion + "/binaries/apache-maven-" + mavenVersion + "-bin.tar.gz"

	yayURL = "https://aur.archlinux.org/yay-bin.git"
)

var goURL = "https://golang.org/dl/go" + goVersions[0] + ".linux-amd64.tar.gz"

//go:build mage

package main

const (
	// https://repo1.maven.org/maven2/io/grpc/protoc-gen-grpc-java/
	protocGenGRPCJavaVersion = "1.68.1"

	// https://maven.apache.org/download.cgi
	mavenVersion = "3.9.9"
)

type zuluVersion struct {
	jdkMajorVersion string
	jdkVersion      string
	version         string
	tag             string
}

var (
	// first is default
	zuluVersions = []zuluVersion{
		// https://www.azul.com/downloads/?version=java-21-lts&os=linux&architecture=x86-64-bit&package=jdk-crac#zulu
		{
			jdkMajorVersion: "21",
			jdkVersion:      "21.0.5",
			version:         "21.38.21",
			tag:             "ca-crac",
		},

		// https://www.azul.com/downloads/?version=java-23&os=linux&architecture=x86-64-bit&package=jdk-crac#zulu
		{
			jdkMajorVersion: "23",
			jdkVersion:      "23.0.1",
			version:         "23.30.13",
			tag:             "ca-crac",
		},

		// https://www.azul.com/downloads/?version=java-17-lts&os=linux&architecture=x86-64-bit&package=jdk-crac#zulu
		{
			jdkMajorVersion: "17",
			jdkVersion:      "17.0.13",
			version:         "17.54.21",
			tag:             "ca-crac",
		},

		// https://www.azul.com/downloads/?version=java-11-lts&os=linux&architecture=x86-64-bit&package=jdk#zulu
		{
			jdkMajorVersion: "11",
			jdkVersion:      "11.0.25",
			version:         "11.76.21",
			tag:             "ca",
		},
	}

	// https://golang.org/dl/
	goVersions = []string{"1.23.3"} // first is default

	pacmanPackageNames = []string{
		"which", "wget", "vim", "nano", "zip", "unzip", "htop", "gcc", "make", "gnu-netcat", "socat", "docker", "docker-buildx", "fontconfig",
		"postgresql", "git", "graphviz", "inetutils", "openssh", "man-db", "man-pages", "diffutils", "fakeroot", "restic", "dnsutils", "ack",
		"imagemagick", "zsh", "patch", "protobuf", "podman", "kubectl", "helm", "helmfile", "k9s", "hyperfine", "jq", "base-devel", "zoxide",
		"fzf", "act", "gnupg", "pwgen", "python-pipx", "git-delta", "lazygit", "cmctl",

		// dependencies for Playwright
		"nss", "nspr", "atk", "at-spi2-atk", "libdrm", "libxkbcommon", "at-spi2-core", "libxcomposite", "libxdamage", "libxfixes", "libxrandr",
		"mesa", "alsa-lib", "libxcursor", "gtk3",
	}

	yayPackageNames = []string{
		"icdiff", "nvm-git", "google-cloud-cli", "google-cloud-cli-gke-gcloud-auth-plugin",
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
		"github.com/blizzy78/xmlquery@latest",
		"github.com/errata-ai/vale/v2/cmd/vale@latest",
	}

	// https://github.com/nodejs/release#release-schedule
	// first is default
	nodeVersionNames = []string{"lts/jod", "21"}

	npmVersion      = "10.8.3"
	npmPackageNames = []string{"serve", "pnpm", "yarn"}

	volumeFolders = []string{
		".vscode-server/extensions",
		"workspaces",
		".zsh_history_dir",
		".zshrc_dir",
		".gitconfig_dir",
		".m2",
		"sophora-repo",
		".ssh",
		"restic-repos",
		".containerrunner",
		".cache/go-build/fuzz",
		".kube",
		".config",
		".zoxide",
		".supermaven",
	}

	protocGoModuleURLs = []string{
		"google.golang.org/protobuf/cmd/protoc-gen-go@latest",
		"google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest",
	}

	gitCompletionBashURL = "https://raw.githubusercontent.com/git/git/master/contrib/completion/git-completion.bash"
	gitCompletionZSHURL  = "https://raw.githubusercontent.com/git/git/master/contrib/completion/git-completion.zsh"

	gitPromptURL   = "https://raw.githubusercontent.com/git/git/master/contrib/completion/git-prompt.sh"
	gitPromptZSHRC = `setopt PROMPT_SUBST` + "\n" +
		`PS1=$'%{\e[30;44m%} %d %{\e[00m\e[30;42m%}\$(__git_ps1 \" %s \")%{\e[00m%}%(?..%{\e[30;41m%} %? %{\e[00m%}) '` + "\n"
)

const tz = "Europe/Berlin"

const (
	protocGenGRPCJavaURL = "https://repo1.maven.org/maven2/io/grpc/protoc-gen-grpc-java/" + protocGenGRPCJavaVersion + "/protoc-gen-grpc-java-" + protocGenGRPCJavaVersion + "-linux-x86_64.exe"

	mavenURL = "https://dlcdn.apache.org/maven/maven-3/" + mavenVersion + "/binaries/apache-maven-" + mavenVersion + "-bin.tar.gz"

	yayURL = "https://aur.archlinux.org/yay-bin.git"
)

var goURL = "https://golang.org/dl/go" + goVersions[0] + ".linux-amd64.tar.gz"

func zuluJDKURL(version zuluVersion) string {
	return "https://cdn.azul.com/zulu/bin/zulu" + version.version + "-" + version.tag + "-jdk" + version.jdkVersion + "-linux_x64.tar.gz"
}

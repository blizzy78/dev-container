//go:build mage

package main

const (
	// https://repo1.maven.org/maven2/io/grpc/protoc-gen-grpc-java/
	protocGenGRPCJavaVersion = "1.59.0"

	// https://maven.apache.org/download.cgi
	mavenVersion = "3.9.5"
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
		// https://www.azul.com/downloads/?version=java-17-lts&architecture=x86-64-bit&package=jdk-crac#zulu
		{
			jdkMajorVersion: "17",
			jdkVersion:      "17.0.8.1",
			version:         "17.44.55",
			tag:             "ca-crac",
		},

		// https://www.azul.com/downloads/?version=java-11-lts&os=linux&architecture=x86-64-bit&package=jdk#zulu
		{
			jdkMajorVersion: "11",
			jdkVersion:      "11.0.21",
			version:         "11.68.17",
			tag:             "ca",
		},
	}

	// https://golang.org/dl/
	goVersions = []string{"1.21.3"} // first is default

	pacmanPackageNames = []string{
		"which", "wget", "vim", "nano", "zip", "unzip", "htop", "gcc", "make", "gnu-netcat", "socat", "docker", "docker-buildx", "fontconfig",
		"postgresql", "git", "graphviz", "inetutils", "openssh", "man-db", "man-pages", "diffutils", "bash-completion", "fakeroot",
		"restic", "dnsutils", "ack", "imagemagick", "zsh", "patch", "protobuf", "podman", "kubectl", "helm", "helmfile", "k9s",
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
	nodeVersionNames = []string{"21", "20"}

	npmPackageNames = []string{"serve", "pnpm", "yarn"}

	volumeFolders = []string{
		".vscode-server/extensions",
		".vscode-server-insiders/extensions",
		"workspaces",
		".bash_history_dir",
		".bashrc_dir",
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
	}

	protocGoModuleURLs = []string{
		"google.golang.org/protobuf/cmd/protoc-gen-go@latest",
		"google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest",
	}

	gitCompletionBashURL = "https://raw.githubusercontent.com/git/git/master/contrib/completion/git-completion.bash"
	gitCompletionZSHURL  = "https://raw.githubusercontent.com/git/git/master/contrib/completion/git-completion.zsh"

	gitPromptURL    = "https://raw.githubusercontent.com/git/git/master/contrib/completion/git-prompt.sh"
	gitPromptBashRC = `export PS1="\[\033[30;44m\] \w \[\033[00m\]\[\033[30;42m\]\$(__git_ps1 \" %s \")\[\033[00m\] "` + "\n"
	gitPromptZSHRC  = `setopt PROMPT_SUBST` + "\n" +
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

//go:build mage

package main

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

var (
	sudoPacmanInstall = sh.RunCmd("sudo", "pacman", "-S", "--noconfirm", "--needed")
	yayInstall        = sh.RunCmd("yay", "-S", "--noconfirm", "--needed")
	ln                = sh.RunCmd("ln", "-s")
	sudoLn            = sh.RunCmd("sudo", "ln", "-s")
	g0                = sh.RunCmd("/go/bin/go")
)

var goMu sync.Mutex

var Default = Install

func Install(ctx context.Context) {
	mg.CtxDeps(ctx,
		timezone,
		caCertificates,
		pacmanPackages,
	)

	// need to do this separately because of different working directory
	mg.CtxDeps(ctx, yay)

	// these all depend on yay
	mg.CtxDeps(ctx,
		bashrc,
		volumes,
		dockerGroup,
		yayPackages,
		installGo,
		installGoSecondary,
		goModules,
		protocGenGRPCJava,
		protocGoModules,
		nvm,
		nodeJS,
		npmPackages,
		jdk,
		maven,
		gitCompletion,
		gitPrompt,
		zoxide,
	)

	mg.CtxDeps(ctx, manPages)
}

func timezone() error {
	if err := sudoLn("/usr/share/zoneinfo/"+tz, "/etc/localtime"); err != nil {
		return fmt.Errorf("sudo ln /etc/localtime: %w", err)
	}

	return nil
}

func caCertificates(ctx context.Context) error {
	mg.CtxDeps(ctx, timezone)

	if err := sh.Run("sudo", "update-ca-trust"); err != nil {
		return fmt.Errorf("sudo update-ca-trust: %w", err)
	}

	return nil
}

func pacmanPackages(ctx context.Context) error {
	mg.CtxDeps(ctx, timezone, caCertificates)

	if err := sudoPacmanInstall(pacmanPackageNames...); err != nil {
		return fmt.Errorf("pacman install packages: %w", err)
	}

	return nil
}

func yayPackages(ctx context.Context) error {
	mg.CtxDeps(ctx, timezone, caCertificates, yay)

	if err := yayInstall(yayPackageNames...); err != nil {
		return fmt.Errorf("yay install packages: %w", err)
	}

	return nil
}

func installGo(ctx context.Context) error {
	mg.CtxDeps(ctx, timezone, caCertificates, bashrc)

	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("Getwd: %w", err)
	}

	if err := downloadAndUnTarGZIPTo(ctx, goURL, wd); err != nil {
		return fmt.Errorf("download and extract Go: %w", err)
	}

	if err := os.Rename("go", "go-"+goVersions[0]); err != nil {
		return fmt.Errorf("rename Go folder: %w", err)
	}

	if err := ln("go-"+goVersions[0], "go"); err != nil {
		return fmt.Errorf("ln Go folder: %w", err)
	}

	if err := sudoLn(wd+"/go", "/go"); err != nil {
		return fmt.Errorf("sudo ln Go folder: %w", err)
	}

	if err := sudoLn("/go/bin/go", "/usr/bin/go"); err != nil {
		return fmt.Errorf("sudo ln /go/bin/go: %w", err)
	}

	if err := appendText(".bashrc", "export PATH=\"$PATH:/go/bin\"\n"); err != nil {
		return fmt.Errorf("add Go PATH to .bashrc: %w", err)
	}

	if err := appendText(".zshrc", "export PATH=\"$PATH:/go/bin\"\n"); err != nil {
		return fmt.Errorf("add Go PATH to .zshrc: %w", err)
	}

	return nil
}

func installGoSecondary(ctx context.Context) error {
	mg.CtxDeps(ctx, timezone, caCertificates, installGo)

	for _, ver := range goVersions[1:] {
		if err := g0("install", "golang.org/dl/go"+ver+"@latest"); err != nil {
			return fmt.Errorf("install Go dl v%s: %w", ver, err)
		}

		if err := sh.Run("/go/bin/go"+ver, "download"); err != nil {
			return fmt.Errorf("download Go v%s: %w", ver, err)
		}
	}

	return nil
}

func goModules(ctx context.Context) error {
	mg.CtxDeps(ctx, timezone, caCertificates, installGo, installGoSecondary)

	for _, mod := range goToolModules {
		err := func() error {
			goMu.Lock()
			defer goMu.Unlock()

			return g0("install", mod)
		}()

		if err != nil {
			return fmt.Errorf("install Go module %s: %w", mod, err)
		}
	}

	if err := ln("dlv", "/home/vscode/go/bin/dlv-dap"); err != nil {
		return fmt.Errorf("ln dlv-dap: %w", err)
	}

	return nil
}

func protocGenGRPCJava(ctx context.Context) error {
	mg.CtxDeps(ctx, timezone, caCertificates)

	name := "protoc-gen-grpc-java-" + protocGenGRPCJavaVersion
	if err := downloadAs(ctx, protocGenGRPCJavaURL, name); err != nil {
		return fmt.Errorf("download protoc-gen-grpc-java: %w", err)
	}

	if err := ln(name, "protoc-gen-grpc-java"); err != nil {
		return fmt.Errorf("ln protoc-gen-grpc-java: %w", err)
	}

	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("Getwd: %w", err)
	}

	if err := sudoLn(wd+"/protoc-gen-grpc-java", "/usr/bin/protoc-gen-grpc-java"); err != nil {
		return fmt.Errorf("sudo ln /usr/bin/protoc-gen-grpc-java: %w", err)
	}

	return nil
}

func protocGoModules(ctx context.Context) error {
	mg.CtxDeps(ctx, timezone, caCertificates, installGo)

	for _, mod := range protocGoModuleURLs {
		err := func() error {
			goMu.Lock()
			defer goMu.Unlock()

			return g0("install", mod)
		}()

		if err != nil {
			return fmt.Errorf("install protoc Go module %s: %w", mod, err)
		}
	}

	return nil
}

func nvm(ctx context.Context) error {
	mg.CtxDeps(ctx, timezone, caCertificates, bashrc, pacmanPackages, yayPackages)

	if err := appendText(".bashrc", ". /usr/share/nvm/init-nvm.sh\nnvm use --silent default\n"); err != nil {
		return fmt.Errorf("add init-nvm.sh to .bashrc: %w", err)
	}

	if err := appendText(".zshrc", ". /usr/share/nvm/init-nvm.sh\nnvm use --silent default\n"); err != nil {
		return fmt.Errorf("add init-nvm.sh to .zshrc: %w", err)
	}

	return nil
}

func nodeJS(ctx context.Context) error {
	mg.CtxDeps(ctx, timezone, caCertificates, pacmanPackages, nvm)

	s := ". /usr/share/nvm/init-nvm.sh\n"

	for _, v := range nodeVersionNames {
		if strings.HasPrefix(v, "lts/") {
			v = strings.TrimPrefix(v, "lts/")
			s += "nvm install --lts=" + v + "\n"
			continue
		}

		s += "nvm install " + v + "\n"
	}

	s += "nvm alias default " + nodeVersionNames[0]

	if err := bashStdin(strings.NewReader(s), "-e"); err != nil {
		return fmt.Errorf("run node install script: %w", err)
	}

	return nil
}

func npmPackages(ctx context.Context) error {
	mg.CtxDeps(ctx, timezone, caCertificates, nodeJS)

	s := ". /usr/share/nvm/init-nvm.sh\n" +
		"nvm use --silent default\n" +
		"npm install -g --no-audit --no-fund npm\n" +
		"npm install -g --no-audit --no-fund " + strings.Join(npmPackageNames, " ")

	if err := bashStdin(strings.NewReader(s), "-e"); err != nil {
		return fmt.Errorf("npm install packages: %w", err)
	}

	return nil
}

func jdk(ctx context.Context) error {
	mg.CtxDeps(ctx, timezone, caCertificates, bashrc)

	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("Getwd: %w", err)
	}

	for _, version := range zuluVersions {
		if err := downloadAndUnTarGZIPTo(ctx, zuluJDKURL(version), wd); err != nil {
			return fmt.Errorf("download and extract JDK: %w", err)
		}

		name := "zulu" + version.version + "-" + version.tag + "-jdk" + version.jdkVersion + "-linux_x64"
		if err := ln(name, "jdk-"+version.jdkMajorVersion); err != nil {
			return fmt.Errorf("ln JDK "+version.jdkMajorVersion+" folder: %w", err)
		}
	}

	defaultVersion := zuluVersions[0]
	if err := ln("jdk-"+defaultVersion.jdkMajorVersion, "jdk"); err != nil {
		return fmt.Errorf("ln JDK folder: %w", err)
	}

	path := wd + "/jdk/bin"
	dir := os.DirFS(path)
	files, err := fs.Glob(dir, "*")
	if err != nil {
		return fmt.Errorf("find JDK binaries: %w", err)
	}

	for _, f := range files {
		if err = sudoLn(path+"/"+f, "/usr/bin/"+f); err != nil {
			return fmt.Errorf("sudo ln /usr/bin/%s: %w", f, err)
		}
	}

	if err := appendText(".bashrc", "export JAVA_HOME=\""+wd+"/jdk\"\n"); err != nil {
		return fmt.Errorf("add JAVA_HOME to .bashrc: %w", err)
	}

	if err := appendText(".zshrc", "export JAVA_HOME=\""+wd+"/jdk\"\n"); err != nil {
		return fmt.Errorf("add JAVA_HOME to .zshrc: %w", err)
	}

	return nil
}

func maven(ctx context.Context) error {
	mg.CtxDeps(ctx, timezone, caCertificates)

	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("Getwd: %w", err)
	}

	if err := downloadAndUnTarGZIPTo(ctx, mavenURL, wd); err != nil {
		return fmt.Errorf("download and extract Maven: %w", err)
	}

	name := "apache-maven-" + mavenVersion
	if err := ln(name, "maven"); err != nil {
		return fmt.Errorf("ln Maven folder: %w", err)
	}

	if err := sudoLn(wd+"/maven/bin/mvn", "/usr/bin/mvn"); err != nil {
		return fmt.Errorf("sudo ln /usr/bin/maven: %w", err)
	}

	return nil
}

func volumes(ctx context.Context) error {
	mg.CtxDeps(ctx, timezone)

	for _, f := range volumeFolders {
		if err := mkdir(f); err != nil {
			return fmt.Errorf("create volume folder %s: %w", f, err)
		}
	}

	for _, f := range []string{".bash_history", ".zsh_history", ".gitconfig"} {
		if err := ln(f+"_dir/"+f, f); err != nil {
			return fmt.Errorf("ln volume folder %s_dir: %w", f, err)
		}
	}

	return nil
}

func dockerGroup(ctx context.Context) error {
	mg.CtxDeps(ctx, timezone, pacmanPackages)

	if err := sh.Run("sudo", "usermod", "-G", "docker", "vscode"); err != nil {
		return fmt.Errorf("sudo usermod: %w", err)
	}

	return nil
}

func gitCompletion(ctx context.Context) error {
	mg.CtxDeps(ctx, timezone, caCertificates, bashrc)

	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("Getwd: %w", err)
	}

	if err = downloadAs(ctx, gitCompletionBashURL, wd+"/.git-completion.sh"); err != nil {
		return fmt.Errorf("download git-completion.sh: %w", err)
	}

	if err = os.Mkdir(wd+"/.zsh", os.ModePerm); err != nil {
		return fmt.Errorf("mkdir .zsh: %w", err)
	}

	if err = downloadAs(ctx, gitCompletionZSHURL, wd+"/.zsh/_git"); err != nil {
		return fmt.Errorf("download git-completion.zsh: %w", err)
	}

	if err := appendText(".bashrc", ". ~/.git-completion.sh\n"); err != nil {
		return fmt.Errorf("add git-completion.sh to .bashrc: %w", err)
	}

	if err := appendText(".zshrc", "fpath=(~/.zsh $fpath)\nzstyle ':completion:*:*:git:*' script ~/.git-completion.sh\n"); err != nil {
		return fmt.Errorf("add git-completion.zsh to .zshrc: %w", err)
	}

	return nil
}

func gitPrompt(ctx context.Context) error {
	mg.CtxDeps(ctx, timezone, caCertificates, bashrc)

	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("Getwd: %w", err)
	}

	if err = downloadAs(ctx, gitPromptURL, wd+"/.git-prompt.sh"); err != nil {
		return fmt.Errorf("download git-prompt: %w", err)
	}

	if err := appendText(".bashrc", ". ~/.git-prompt.sh\n"+gitPromptBashRC); err != nil {
		return fmt.Errorf("add git-prompt to .bashrc: %w", err)
	}

	if err := appendText(".zshrc", ". ~/.git-prompt.sh\n"+gitPromptZSHRC); err != nil {
		return fmt.Errorf("add git-prompt to .zshrc: %w", err)
	}

	return nil
}

func bashrc(ctx context.Context) error {
	mg.CtxDeps(ctx, timezone)

	if err := appendText(".bashrc", "[[ -f ~/.bashrc_dir/.bashrc ]] && . ~/.bashrc_dir/.bashrc\n"); err != nil {
		return fmt.Errorf("add ~/.bashrc_dir/.bashrc to .bashrc: %w", err)
	}

	if err := appendText(".zshrc", "[[ -f ~/.zshrc_dir/.zshrc ]] && . ~/.zshrc_dir/.zshrc\n"); err != nil {
		return fmt.Errorf("add ~/.zshrc_dir/.zshrc to .zshrc: %w", err)
	}

	return nil
}

func manPages(ctx context.Context) error {
	mg.CtxDeps(ctx, timezone, pacmanPackages)

	if err := sh.Run("sudo", "mandb", "-c"); err != nil {
		return fmt.Errorf("sudo mandb: %w", err)
	}

	return nil
}

func yay(ctx context.Context) error {
	mg.CtxDeps(ctx, timezone, caCertificates, pacmanPackages)

	return doInTempDir(func() error {
		if err := sh.Run("git", "clone", yayURL); err != nil {
			return fmt.Errorf("git clone: %w", err)
		}

		wd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("get wd: %w", err)
		}

		repoName := path.Base(strings.TrimSuffix(yayURL, ".git"))

		return doInDir(wd+"/"+repoName, func() error {
			if err := sh.Run("makepkg", "-srci", "--noconfirm"); err != nil {
				return fmt.Errorf("makepkg: %w", err)
			}

			return nil
		})
	})
}

func zoxide(ctx context.Context) error {
	mg.CtxDeps(ctx, bashrc, pacmanPackages)

	if err := appendText(".bashrc", "export _ZO_DATA_DIR=/home/vscode/.zoxide\neval \"$(zoxide init --cmd cd bash)\"\n"); err != nil {
		return fmt.Errorf("add zoxide to .bashrc: %w", err)
	}

	if err := appendText(".zshrc", "export _ZO_DATA_DIR=/home/vscode/.zoxide\neval \"$(zoxide init --cmd cd zsh)\"\n"); err != nil {
		return fmt.Errorf("add zoxide to .zshrc: %w", err)
	}

	return nil
}

func doInTempDir(f func() error) error {
	dir, err := os.MkdirTemp("", "")
	if err != nil {
		return fmt.Errorf("make temp dir: %w", err)
	}
	defer os.RemoveAll(dir)

	return doInDir(dir, f)
}

func doInDir(dir string, f func() error) error {
	oldWd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get wd: %w", err)
	}

	if err := os.Chdir(dir); err != nil {
		return fmt.Errorf("change dir: %w", err)
	}
	defer os.Chdir(oldWd)

	return f()
}

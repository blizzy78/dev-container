//go:build mage

package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

var (
	aptInstall = sh.RunCmd("sudo", "apt", "install", "-y")
	ln         = sh.RunCmd("ln", "-s")
	sudoLn     = sh.RunCmd("sudo", "ln", "-s")
	g0         = sh.RunCmd("/go/bin/go")
	npmInstall = sh.RunCmd("npm", "install")
)

var (
	systemMu sync.Mutex
	goMu     sync.Mutex
	npmMu    sync.Mutex
)

var Default = Install

func Install(ctx context.Context) {
	mg.CtxDeps(ctx, timezone)

	mg.CtxDeps(ctx, aptPackages, caCertificates)

	mg.CtxDeps(ctx,
		installGo,
		goTools,
		mage,
		protoc,
		protocGenGRPCJava,
		protocGoModules,
		npm,
		npmPackages,
		jdk,
		maven,
		volumes,
		locales,
		gatsby,
		restic,
	)
}

func aptPackages() error {
	systemMu.Lock()
	defer systemMu.Unlock()

	if err := sh.Run("sudo", "apt", "update"); err != nil {
		return fmt.Errorf("sudo apt update: %w", err)
	}

	if err := aptInstall(aptPackageNames...); err != nil {
		return fmt.Errorf("apt install packages: %w", err)
	}

	return nil
}

func caCertificates(ctx context.Context) error {
	mg.CtxDeps(ctx, aptPackages, installGo)

	if err := sh.Run("sudo", "update-ca-certificates"); err != nil {
		return fmt.Errorf("sudo update-ca-certificates: %w", err)
	}

	return nil
}

func installGo(ctx context.Context) error {
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("Getwd: %w", err)
	}

	if err := downloadAndUnTarGZIPTo(ctx, goURL, wd); err != nil {
		return fmt.Errorf("download and extract Go: %w", err)
	}

	if err := os.Rename("go", "go-"+goVersion); err != nil {
		return fmt.Errorf("rename Go folder: %w", err)
	}

	if err := ln("go-"+goVersion, "go"); err != nil {
		return fmt.Errorf("ln Go folder: %w", err)
	}

	if err := sudoLn(wd+"/go", "/go"); err != nil {
		return fmt.Errorf("sudo ln Go folder: %w", err)
	}

	if err := sudoLn("/go/bin/go", "/usr/bin/go"); err != nil {
		return fmt.Errorf("sudo ln /go/bin/go: %w", err)
	}

	file, err := os.OpenFile(".bashrc", os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		return fmt.Errorf("open .bashrc: %w", err)
	}
	defer file.Close()

	if _, err = file.WriteString("export PATH=\"$PATH:/go/bin\"\n"); err != nil {
		return fmt.Errorf("write Go PATH to .bashrc: %w", err)
	}

	return nil
}

func goTools(ctx context.Context) error {
	mg.CtxDeps(ctx, aptPackages, installGo)

	for _, u := range goToolURLs {
		err := func() error {
			goMu.Lock()
			defer goMu.Unlock()

			return g0("install", u)
		}()

		if err != nil {
			return fmt.Errorf("install Go tools: %w", err)
		}
	}

	if err := ln("dlv", "/home/vscode/go/bin/dlv-dap"); err != nil {
		return fmt.Errorf("ln dlv-dap: %w", err)
	}

	return nil
}

func mage(ctx context.Context) error {
	mg.CtxDeps(ctx, aptPackages, installGo)

	if err := sh.Run("git", "clone", "https://github.com/magefile/mage"); err != nil {
		return fmt.Errorf("git clone mage: %w", err)
	}
	defer sh.Rm("mage")

	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("Getwd: %w", err)
	}

	goMu.Lock()
	defer goMu.Unlock()

	c := exec.Command("go", "run", "bootstrap.go")
	c.Dir = wd + "/mage"
	o, err := c.CombinedOutput()
	if mg.Verbose() {
		_, _ = os.Stdout.Write(o)
	}

	if err != nil {
		return fmt.Errorf("go run bootstrap.go: %w", err)
	}

	return nil
}

func protoc(ctx context.Context) error {
	dir := "protoc-" + protocVersion
	if err := mkdir(dir); err != nil {
		return fmt.Errorf("create protoc folder: %w", err)
	}

	if err := downloadAndUnZipTo(ctx, protocURL, dir); err != nil {
		return fmt.Errorf("download and extract protoc: %w", err)
	}

	if err := ln(dir, "protoc"); err != nil {
		return fmt.Errorf("ln protoc folder: %w", err)
	}

	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("Getwd: %w", err)
	}

	if err := sudoLn(wd+"/protoc/bin/protoc", "/usr/bin/protoc"); err != nil {
		return fmt.Errorf("sudo ln /usr/bin/protoc: %w", err)
	}

	return nil
}

func protocGenGRPCJava(ctx context.Context) error {
	mg.CtxDeps(ctx, protoc)

	name := "protoc-gen-grpc-java-" + protocGenGRPCJavaVersion
	if err := downloadAs(ctx, protocGenGRPCJavaURL, "protoc/bin/"+name); err != nil {
		return fmt.Errorf("download protoc-gen-grpc-java: %w", err)
	}

	if err := ln(name, "protoc/bin/protoc-gen-grpc-java"); err != nil {
		return fmt.Errorf("ln protoc-gen-grpc-java: %w", err)
	}

	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("Getwd: %w", err)
	}

	if err := sudoLn(wd+"/protoc/bin/protoc-gen-grpc-java", "/usr/bin/protoc-gen-grpc-java"); err != nil {
		return fmt.Errorf("sudo ln /usr/bin/protoc-gen-grpc-java: %w", err)
	}

	return nil
}

func protocGoModules(ctx context.Context) error {
	mg.CtxDeps(ctx, installGo)

	goMu.Lock()
	defer goMu.Unlock()

	if err := g0(append([]string{"get"}, protocGoModuleURLs...)...); err != nil {
		return fmt.Errorf("go get protoc modules: %w", err)
	}

	return nil
}

func nvm(ctx context.Context) error {
	b, err := download(ctx, nvmInstallURL)
	if err != nil {
		return fmt.Errorf("download nvm install script: %w", err)
	}

	if err := bashStdin(bytes.NewReader(b)); err != nil {
		return fmt.Errorf("run nvm install script: %w", err)
	}

	return nil
}

func nodeJS(ctx context.Context) error {
	mg.CtxDeps(ctx, nvm)

	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("Getwd: %w", err)
	}

	s := "export NVM_DIR=\"" + wd + "/.nvm\"\n" +
		". ${NVM_DIR}/nvm.sh\n"

	for _, v := range nodeLTSNames {
		s += "nvm install --lts=" + v + "\n"
	}

	s += "nvm alias default lts/" + nodeLTSNames[0] + "\n" +
		"nvm use default\n" +
		"sudo ln -s $(which node) /usr/bin/node"

	if err := bashStdin(strings.NewReader(s), "-e"); err != nil {
		return fmt.Errorf("run node install script: %w", err)
	}

	return nil
}

func npm(ctx context.Context) error {
	if err := downloadAs(ctx, npmInstallURL, "install-npm.sh"); err != nil {
		return fmt.Errorf("download npm install script: %w", err)
	}
	defer sh.Rm("install-npm.sh")

	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("Getwd: %w", err)
	}

	mg.CtxDeps(ctx, nodeJS)

	s := "export NVM_DIR=\"" + wd + "/.nvm\"\n" +
		". ${NVM_DIR}/nvm.sh\n" +
		"bash install-npm.sh\n" +
		"sudo ln -s $(which npm) /usr/bin/npm\n" +
		"sudo ln -s $(which npx) /usr/bin/npx"

	if err := bashStdin(strings.NewReader(s), "-e"); err != nil {
		return fmt.Errorf("run npm install script: %w", err)
	}

	return nil
}

func npmPackages(ctx context.Context) error {
	mg.CtxDeps(ctx, npm)

	npmMu.Lock()
	defer npmMu.Unlock()

	if err := npmInstall(append([]string{"-g"}, npmPackageNames...)...); err != nil {
		return fmt.Errorf("npm install packages: %w", err)
	}

	return nil
}

func jdk(ctx context.Context) error {
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("Getwd: %w", err)
	}

	name := "zulu" + zuluVersion + "-ca-jdk" + zuluJDKVersion + "-linux_x64"
	if err := downloadAndUnTarGZIPTo(ctx, zuluJDKURL, wd); err != nil {
		return fmt.Errorf("download and extract JDK: %w", err)
	}

	if err := ln(name, "jdk"); err != nil {
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

	return nil
}

func maven(ctx context.Context) error {
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

func volumes() error {
	for _, f := range volumeFolders {
		if err := mkdir(f); err != nil {
			return fmt.Errorf("create volume folder %s: %w", f, err)
		}
	}

	for _, f := range []string{".bash_history", ".gitconfig"} {
		if err := mkdir(f + "_dir"); err != nil {
			return fmt.Errorf("create volume folder %s_dir: %w", f, err)
		}
		if err := ln(f+"_dir/"+f, f); err != nil {
			return fmt.Errorf("ln volume folder %s_dir: %w", f, err)
		}
	}

	return nil
}

// https://stackoverflow.com/a/20693661
func timezone() error {
	systemMu.Lock()
	defer systemMu.Unlock()

	err := func() error {
		f, err := os.OpenFile("preseed.txt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
		if err != nil {
			return fmt.Errorf("open preseed.txt: %w", err)
		}
		defer f.Close()

		parts := strings.SplitN(tz, "/", 2)
		s := "tzdata tzdata/Areas select " + parts[0] + "\n" +
			"tzdata tzdata/Zones/" + parts[0] + " select " + parts[1] + "\n"
		if _, err = io.WriteString(f, s); err != nil {
			return fmt.Errorf("write preseed.txt: %w", err)
		}

		return nil
	}()

	if err != nil {
		return err
	}

	defer sh.Rm("preseed.txt")

	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("Getwd: %w", err)
	}

	if err := sh.Run("sudo", "debconf-set-selections", wd+"/preseed.txt"); err != nil {
		return fmt.Errorf("sudo debconf-set-selections: %w", err)
	}

	if err := sh.Run("sudo", "env", "DEBIAN_FRONTEND=noninteractive", "DEBCONF_NONINTERACTIVE_SEEN=true", "apt", "install", "-y", "tzdata"); err != nil {
		return fmt.Errorf("sudo apt install tzdata: %w", err)
	}

	return nil
}

func locales(ctx context.Context) error {
	mg.CtxDeps(ctx, aptPackages)

	systemMu.Lock()
	defer systemMu.Unlock()

	for _, l := range []string{"en_US", "en_US.UTF-8"} {
		if err := sh.Run("sudo", "locale-gen", l); err != nil {
			return fmt.Errorf("sudo locale-gen %s: %w", l, err)
		}
	}

	return nil
}

func gatsby(ctx context.Context) error {
	mg.CtxDeps(ctx, npmPackages)

	if err := sh.Run("npx", "gatsby", "telemetry", "--disable"); err != nil {
		return fmt.Errorf("disable Gatsby telemetry: %w", err)
	}

	return nil
}

func restic(ctx context.Context) error {
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("Getwd: %w", err)
	}

	if err := downloadAndUnBZip2To(ctx, resticURL, wd+"/restic_"+resticVersion+"_linux_amd64"); err != nil {
		return fmt.Errorf("download and extract restic: %w", err)
	}

	if err := os.Chmod(wd+"/restic_"+resticVersion+"_linux_amd64", 0755); err != nil {
		return fmt.Errorf("chmod restic: %w", err)
	}

	if err := ln(wd+"/restic_"+resticVersion+"_linux_amd64", wd+"/restic"); err != nil {
		return fmt.Errorf("ln restic: %w", err)
	}

	if err := sudoLn(wd+"/restic", "/usr/bin/restic"); err != nil {
		return fmt.Errorf("sudo ln /usr/bin/restic: %w", err)
	}

	return nil
}

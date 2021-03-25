//+build mage

package main

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

const (
	protocVersion            = "3.15.6"
	protocGenGRPCJavaVersion = "1.36.0"

	nvmVersion = "0.37.2"

	zuluVersion    = "11.45.27"
	zuluJDKVersion = "11.0.10"

	mavenVersion = "3.6.3"

	tz = "Europe/Berlin"
)

var (
	toolNames = []string{
		"apt-utils", "wget", "less", "vim", "nano", "zip", "unzip", "xz-utils", "htop", "gcc", "make",
	}

	goToolURLs = []string{
		"github.com/uudashr/gopkgs/v2/cmd/gopkgs",
		"github.com/ramya-rao-a/go-outline",
		"github.com/cweill/gotests/gotests",
		"github.com/fatih/gomodifytags",
		"github.com/josharian/impl",
		"github.com/haya14busa/goplay/cmd/goplay",
		"github.com/go-delve/delve/cmd/dlv",
		"github.com/golangci/golangci-lint/cmd/golangci-lint",
		"golang.org/x/tools/gopls",
	}

	postCSSPackages = []string{
		"postcss@latest", "postcss-cli@latest",
	}

	volumeFolders = []string{
		".vscode-server/extensions",
		"workspaces",
		".m2",
		"sophora-repo",
	}
)

const (
	protocURL            = "https://github.com/protocolbuffers/protobuf/releases/download/v" + protocVersion + "/protoc-" + protocVersion + "-linux-x86_64.zip"
	protocGenGRPCJavaURL = "https://repo1.maven.org/maven2/io/grpc/protoc-gen-grpc-java/" + protocGenGRPCJavaVersion + "/protoc-gen-grpc-java-" + protocGenGRPCJavaVersion + "-linux-x86_64.exe"

	nvmInstallURL = "https://raw.githubusercontent.com/nvm-sh/nvm/v" + nvmVersion + "/install.sh"
	npmInstallURL = "https://www.npmjs.com/install.sh"

	zuluJDKURL = "https://cdn.azul.com/zulu/bin/zulu" + zuluVersion + "-ca-jdk" + zuluJDKVersion + "-linux_x64.tar.gz"

	mavenURL = "https://mirror.netcologne.de/apache.org/maven/maven-3/" + mavenVersion + "/binaries/apache-maven-" + mavenVersion + "-bin.tar.gz"
)

var (
	aptInstall = sh.RunCmd("sudo", "apt", "install", "-y")
	ln         = sh.RunCmd("ln", "-s")
	sudoLn     = sh.RunCmd("sudo", "ln", "-s")
	g0         = sh.RunCmd("/go/bin/go")
	npmInstall = sh.RunCmd("npm", "install")
)

var Default = Install

func Install(ctx context.Context) {
	mg.SerialCtxDeps(ctx,
		tools,
		goTools,
		mage,
		protoc,
		protocGenGRPCJava,
		npm,
		postCSS,
		jdk,
		maven,
		volumes,
		timezone,
	)
}

func tools() error {
	return aptInstall(toolNames...)
}

func goTools() error {
	return g0(append([]string{"get"}, goToolURLs...)...)
}

func mage() error {
	if err := sh.Run("git", "clone", "https://github.com/magefile/mage"); err != nil {
		return err
	}
	defer sh.Rm("mage")

	if err := os.Chdir("mage"); err != nil {
		return err
	}
	defer os.Chdir("..")

	return sh.Run("go", "run", "bootstrap.go")
}

func protoc(ctx context.Context) error {
	dir := "protoc-" + protocVersion
	if err := mkdir(dir); err != nil {
		return err
	}

	if err := downloadAndUnZipTo(ctx, protocURL, dir); err != nil {
		return err
	}

	if err := ln(dir, "protoc"); err != nil {
		return err
	}

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	return sudoLn(wd+"/protoc/bin/protoc", "/usr/bin/protoc")
}

func protocGenGRPCJava(ctx context.Context) error {
	mg.SerialCtxDeps(ctx, protoc)

	name := "protoc-gen-grpc-java-" + protocGenGRPCJavaVersion
	if err := downloadAs(ctx, protocGenGRPCJavaURL, "protoc/bin/"+name); err != nil {
		return err
	}

	if err := ln(name, "protoc/bin/protoc-gen-grpc-java"); err != nil {
		return err
	}

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	return sudoLn(wd+"/protoc/bin/protoc-gen-grpc-java", "/usr/bin/protoc-gen-grpc-java")
}

func nvm(ctx context.Context) error {
	if err := downloadAs(ctx, nvmInstallURL, "install-nvm.sh"); err != nil {
		return err
	}
	defer sh.Rm("install-nvm.sh")

	return sh.Run("bash", "install-nvm.sh")
}

func nodeJS(ctx context.Context) error {
	mg.SerialCtxDeps(ctx, nvm)

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	s := "export NVM_DIR=\"" + wd + "/.nvm\"\n" +
		". ${NVM_DIR}/nvm.sh\n" +
		"nvm install node\n" +
		"sudo ln -s $(which node) /usr/bin/node"
	c := exec.Cmd{
		Path: "/usr/bin/bash",
		Args: []string{
			"/usr/bin/bash",
			"-e",
		},
		Stdin: strings.NewReader(s),
	}
	o, err := c.CombinedOutput()
	_, _ = os.Stdout.Write(o)
	return err
}

func npm(ctx context.Context) error {
	mg.SerialCtxDeps(ctx, nodeJS)

	if err := downloadAs(ctx, npmInstallURL, "install-npm.sh"); err != nil {
		return err
	}
	defer sh.Rm("install-npm.sh")

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	s := "export NVM_DIR=\"" + wd + "/.nvm\"\n" +
		". ${NVM_DIR}/nvm.sh\n" +
		"bash install-npm.sh\n" +
		"sudo ln -s $(which npm) /usr/bin/npm\n" +
		"sudo ln -s $(which npx) /usr/bin/npx"
	c := exec.Cmd{
		Path: "/usr/bin/bash",
		Args: []string{
			"/usr/bin/bash",
			"-e",
		},
		Stdin: strings.NewReader(s),
	}
	o, err := c.CombinedOutput()
	_, _ = os.Stdout.Write(o)
	return err
}

func postCSS(ctx context.Context) error {
	mg.SerialCtxDeps(ctx, npm)
	return npmInstall(append([]string{"-g"}, postCSSPackages...)...)
}

func jdk(ctx context.Context) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	name := "zulu" + zuluVersion + "-ca-jdk" + zuluJDKVersion + "-linux_x64"
	if err := downloadAndUnTarGZIPTo(ctx, zuluJDKURL, wd); err != nil {
		return err
	}

	if err := ln(name, "jdk"); err != nil {
		return err
	}

	path := wd + "/jdk/bin"
	dir := os.DirFS(path)
	files, err := fs.Glob(dir, "*")
	if err != nil {
		return err
	}

	for _, f := range files {
		if err = sudoLn(path+"/"+f, "/usr/bin/"+f); err != nil {
			return err
		}
	}

	return nil
}

func maven(ctx context.Context) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	name := "apache-maven-" + mavenVersion
	if err := downloadAndUnTarGZIPTo(ctx, mavenURL, wd); err != nil {
		return err
	}

	if err := ln(name, "maven"); err != nil {
		return err
	}

	return sudoLn(wd+"/maven/bin/mvn", "/usr/bin/mvn")
}

func volumes() error {
	for _, f := range volumeFolders {
		if err := mkdir(f); err != nil {
			return err
		}
	}

	for _, f := range []string{".bash_history", ".gitconfig"} {
		if err := mkdir(f + "_dir"); err != nil {
			return err
		}
		if err := ln(f+"_dir/"+f, f); err != nil {
			return err
		}
	}

	return nil
}

// https://stackoverflow.com/a/20693661
func timezone() error {
	err := func() error {
		f, err := os.OpenFile("preseed.txt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
		if err != nil {
			return err
		}
		defer f.Close()

		parts := strings.SplitN(tz, "/", 2)
		s := "tzdata tzdata/Areas select " + parts[0] + "\n" +
			"tzdata tzdata/Zones/" + parts[0] + " select " + parts[1] + "\n"
		if _, err = io.WriteString(f, s); err != nil {
			return err
		}

		return nil
	}()

	if err != nil {
		return err
	}

	defer sh.Rm("preseed.txt")

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	if err := sh.Run("sudo", "debconf-set-selections", wd+"/preseed.txt"); err != nil {
		return err
	}

	return sh.Run("sudo", "env", "DEBIAN_FRONTEND=noninteractive", "DEBCONF_NONINTERACTIVE_SEEN=true", "apt", "install", "-y", "tzdata")
}

func mkdir(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}

func downloadAndUnZipTo(ctx context.Context, url string, dest string) error {
	b, err := download(ctx, url)
	if err != nil {
		return err
	}
	return unZipTo(b, dest)
}

func downloadAndUnTarGZIPTo(ctx context.Context, url string, dest string) error {
	b, err := download(ctx, url)
	if err != nil {
		return err
	}
	return unTarGZIPTo(b, dest)
}

func unZipTo(b []byte, dest string) error {
	r, err := zip.NewReader(bytes.NewReader(b), int64(len(b)))
	if err != nil {
		return err
	}

	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name)

		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("%s: illegal file path", fpath)
		}

		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(fpath, os.ModePerm); err != nil {
				return err
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		createFile := func() error {
			file, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer file.Close()

			fr, err := f.Open()
			if err != nil {
				return err
			}
			defer fr.Close()

			_, err = io.Copy(file, fr)
			if err != nil {
				return err
			}

			return nil
		}

		if err := createFile(); err != nil {
			return err
		}
	}

	return nil
}

func unTarGZIPTo(b []byte, dest string) error {
	r, err := gzip.NewReader(bytes.NewReader(b))
	if err != nil {
		return err
	}
	defer r.Close()

	tr := tar.NewReader(r)

	for {
		h, err := tr.Next()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return err
		}

		if h.Typeflag != tar.TypeReg && h.Typeflag == tar.TypeDir {
			continue
		}

		fpath := filepath.Join(dest, h.Name)

		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("%s: illegal file path", fpath)
		}

		if h.Typeflag == tar.TypeDir {
			if err := os.MkdirAll(fpath, os.ModePerm); err != nil {
				return err
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		createFile := func() error {
			file, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, h.FileInfo().Mode())
			if err != nil {
				return err
			}
			defer file.Close()

			_, err = io.Copy(file, tr)
			if err != nil {
				return err
			}

			return nil
		}

		if err := createFile(); err != nil {
			return err
		}
	}

	return nil
}

func downloadAs(ctx context.Context, url string, path string) error {
	b, err := download(ctx, url)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(b)
	if err != nil {
		return err
	}

	return nil
}

func download(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return io.ReadAll(res.Body)
}

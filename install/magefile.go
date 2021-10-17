//go:build mage

package main

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

const tz = "Europe/Berlin"

var (
	aptPackageNames = []string{
		"apt-utils", "locales", "wget", "less", "vim", "nano", "zip", "unzip", "xz-utils", "htop", "gcc", "make",
		"telnet", "netcat", "socat", "docker.io", "libfontconfig", "postgresql-client", "iputils-ping", "libxml2-utils",
		"curl", "git", "ca-certificates",
	}

	goToolURLs = []string{
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
	}

	npmPackageNames = []string{
		"postcss@latest",
		"postcss-cli@latest",
		"serve",
		"gatsby",
	}

	volumeFolders = []string{
		".vscode-server/extensions",
		"workspaces",
		".m2",
		"sophora-repo",
		".ssh",
	}

	protocGoModuleURLs = []string{
		"google.golang.org/protobuf/cmd/protoc-gen-go",
		"google.golang.org/grpc/cmd/protoc-gen-go-grpc",
	}
)

const (
	goURL = "https://golang.org/dl/go" + goVersion + ".linux-amd64.tar.gz"

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
	)
}

func aptPackages() error {
	systemMu.Lock()
	defer systemMu.Unlock()

	if err := sh.Run("sudo", "apt", "update"); err != nil {
		return err
	}

	return aptInstall(aptPackageNames...)
}

func caCertificates(ctx context.Context) error {
	mg.CtxDeps(ctx, aptPackages, installGo)
	return sh.Run("sudo", "update-ca-certificates")
}

func installGo(ctx context.Context) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	if err := downloadAndUnTarGZIPTo(ctx, goURL, wd); err != nil {
		return err
	}

	if err := os.Rename("go", "go-"+goVersion); err != nil {
		return err
	}

	if err := ln("go-"+goVersion, "go"); err != nil {
		return err
	}

	if err := sudoLn(wd+"/go", "/go"); err != nil {
		return err
	}

	if err := sudoLn("/go/bin/go", "/usr/bin/go"); err != nil {
		return err
	}

	file, err := os.OpenFile(".bashrc", os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err = file.WriteString("export PATH=\"$PATH:/go/bin\"\n"); err != nil {
		return err
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
			return err
		}
	}

	return ln("dlv", "/home/vscode/go/bin/dlv-dap")
}

func mage(ctx context.Context) error {
	mg.CtxDeps(ctx, aptPackages, installGo)

	if err := sh.Run("git", "clone", "https://github.com/magefile/mage"); err != nil {
		return err
	}
	defer sh.Rm("mage")

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	goMu.Lock()
	defer goMu.Unlock()

	c := exec.Command("go", "run", "bootstrap.go")
	c.Dir = wd + "/mage"
	o, err := c.CombinedOutput()
	if mg.Verbose() {
		_, _ = os.Stdout.Write(o)
	}
	return err
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
	mg.CtxDeps(ctx, protoc)

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

func protocGoModules(ctx context.Context) error {
	mg.CtxDeps(ctx, installGo)

	goMu.Lock()
	defer goMu.Unlock()
	return g0(append([]string{"get"}, protocGoModuleURLs...)...)
}

func nvm(ctx context.Context) error {
	b, err := download(ctx, nvmInstallURL)
	if err != nil {
		return err
	}
	return bashStdin(bytes.NewReader(b))
}

func nodeJS(ctx context.Context) error {
	mg.CtxDeps(ctx, nvm)

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	s := "export NVM_DIR=\"" + wd + "/.nvm\"\n" +
		". ${NVM_DIR}/nvm.sh\n" +
		"nvm install node\n" +
		"sudo ln -s $(which node) /usr/bin/node"
	return bashStdin(strings.NewReader(s), "-e")
}

func npm(ctx context.Context) error {
	if err := downloadAs(ctx, npmInstallURL, "install-npm.sh"); err != nil {
		return err
	}
	defer sh.Rm("install-npm.sh")

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	mg.CtxDeps(ctx, nodeJS)

	s := "export NVM_DIR=\"" + wd + "/.nvm\"\n" +
		". ${NVM_DIR}/nvm.sh\n" +
		"bash install-npm.sh\n" +
		"sudo ln -s $(which npm) /usr/bin/npm\n" +
		"sudo ln -s $(which npx) /usr/bin/npx"
	return bashStdin(strings.NewReader(s), "-e")
}

func npmPackages(ctx context.Context) error {
	mg.CtxDeps(ctx, npm)
	npmMu.Lock()
	defer npmMu.Unlock()
	return npmInstall(append([]string{"-g"}, npmPackageNames...)...)
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

	if err := downloadAndUnTarGZIPTo(ctx, mavenURL, wd); err != nil {
		return err
	}

	name := "apache-maven-" + mavenVersion
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
	systemMu.Lock()
	defer systemMu.Unlock()

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

func locales(ctx context.Context) error {
	mg.CtxDeps(ctx, aptPackages)

	systemMu.Lock()
	defer systemMu.Unlock()

	if err := sh.Run("sudo", "locale-gen", "en_US"); err != nil {
		return err
	}

	return sh.Run("sudo", "locale-gen", "en_US.UTF-8")
}

func gatsby(ctx context.Context) error {
	mg.CtxDeps(ctx, npmPackages)
	return sh.Run("npx", "gatsby", "telemetry", "--disable")
}

func bashStdin(r io.Reader, args ...string) error {
	c := exec.Command("bash", args...)
	c.Stdin = r
	o, err := c.CombinedOutput()
	if mg.Verbose() {
		_, _ = os.Stdout.Write(o)
	}
	return err
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
			if err := mkdir(fpath); err != nil {
				return err
			}
			continue
		}

		if err := mkdir(filepath.Dir(fpath)); err != nil {
			return err
		}

		createFile := func() error {
			fr, err := f.Open()
			if err != nil {
				return err
			}
			defer fr.Close()

			return copyToFile(fr, fpath, f.Mode())
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
			if err := mkdir(fpath); err != nil {
				return err
			}
			continue
		}

		if err := mkdir(filepath.Dir(fpath)); err != nil {
			return err
		}

		if err := copyToFile(tr, fpath, h.FileInfo().Mode()); err != nil {
			return err
		}
	}

	return nil
}

func copyToFile(r io.Reader, path string, perm fs.FileMode) error {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, r)
	return err
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
	c := http.Client{
		Transport: &http.Transport{
			ResponseHeaderTimeout: 10 * time.Second,

			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Minute)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return io.ReadAll(res.Body)
}

//go:build mage

package main

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/bzip2"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func downloadAndUnZipTo(ctx context.Context, url string, dest string) error {
	b, err := download(ctx, url)
	if err != nil {
		return fmt.Errorf("download: %w", err)
	}

	if err := unZipTo(b, dest); err != nil {
		return fmt.Errorf("unzip: %w", err)
	}

	return nil
}

func downloadAndUnTarGZIPTo(ctx context.Context, url string, dest string) error {
	b, err := download(ctx, url)
	if err != nil {
		return fmt.Errorf("download: %w", err)
	}

	if err := unTarGZIPTo(b, dest); err != nil {
		return fmt.Errorf("untargz: %w", err)
	}

	return nil
}

func downloadAndUnBZip2To(ctx context.Context, url string, dest string) error {
	b, err := download(ctx, url)
	if err != nil {
		return fmt.Errorf("download: %w", err)
	}

	if err := unBZip2To(b, dest); err != nil {
		return fmt.Errorf("unbzip2: %w", err)
	}

	return nil
}

func unZipTo(b []byte, dest string) error {
	r, err := zip.NewReader(bytes.NewReader(b), int64(len(b)))
	if err != nil {
		return fmt.Errorf("open zip: %w", err)
	}

	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name)

		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("unzip: %s: illegal file path", fpath)
		}

		if f.FileInfo().IsDir() {
			if err := mkdir(fpath); err != nil {
				return fmt.Errorf("%s: create folder: %w", fpath, err)
			}
			continue
		}

		if err := mkdir(filepath.Dir(fpath)); err != nil {
			return fmt.Errorf("%s: create folder: %w", filepath.Dir(fpath), err)
		}

		createFile := func() error {
			fr, err := f.Open()
			if err != nil {
				return fmt.Errorf("unzip: %s: open: %w", f.Name, err)
			}
			defer fr.Close()

			if err = copyToFile(fr, fpath, f.Mode()); err != nil {
				return fmt.Errorf("unzip: %s: copy: %w", f.Name, err)
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
		return fmt.Errorf("open gzip: %w", err)
	}
	defer r.Close()

	tr := tar.NewReader(r)

	for {
		h, err := tr.Next()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return fmt.Errorf("untar: %w", err)
		}

		if h.Typeflag != tar.TypeReg && h.Typeflag == tar.TypeDir {
			continue
		}

		fpath := filepath.Join(dest, h.Name)

		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("untar: %s: illegal file path", fpath)
		}

		if h.Typeflag == tar.TypeDir {
			if err := mkdir(fpath); err != nil {
				return fmt.Errorf("%s: create folder: %w", fpath, err)
			}
			continue
		}

		if err := mkdir(filepath.Dir(fpath)); err != nil {
			return fmt.Errorf("%s: create folder: %w", filepath.Dir(fpath), err)
		}

		if err := copyToFile(tr, fpath, h.FileInfo().Mode()); err != nil {
			return fmt.Errorf("untar: %s: copy: %w", h.Name, err)
		}
	}

	return nil
}

func unBZip2To(b []byte, dest string) error {
	r := bzip2.NewReader(bytes.NewReader(b))

	if err := copyToFile(r, dest, fs.ModePerm); err != nil {
		return fmt.Errorf("copy: %w", err)
	}

	return nil
}

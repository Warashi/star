package main

import (
	"archive/tar"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/mattn/go-isatty"
	"github.com/pkg/errors"
)

func create(w io.WriteCloser, src string) error {
	defer w.Close()

	// ensure the src actually exists before trying to tar it
	if _, err := os.Stat(src); err != nil {
		return errors.Errorf("Unable to tar file: %v", err)
	}

	tw := tar.NewWriter(w)
	defer tw.Close()

	// walk path
	return filepath.Walk(src, func(file string, fi os.FileInfo, err error) error {

		// return on any error
		if err != nil {
			return err
		}

		// create a new dir/file header
		header, err := tar.FileInfoHeader(fi, fi.Name())
		if err != nil {
			return err
		}

		// update the name to correctly reflect the desired destination when untaring
		header.Name = file

		// write the header
		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		// return on non-regular files
		if !fi.Mode().IsRegular() {
			return nil
		}

		// open files for taring
		f, err := os.Open(file)
		defer f.Close()
		if err != nil {
			return err
		}

		// copy file data into tar writer
		if _, err := io.Copy(tw, f); err != nil {
			return err
		}

		return nil
	})

}

func extract(r io.ReadCloser, dst string) error {
	defer r.Close()

	tr := tar.NewReader(r)

	for {
		header, err := tr.Next()

		switch {
		// if no more files are found return
		case err == io.EOF:
			return nil

		// return any other error
		case err != nil:
			return err

		// if the header is nil, just skip it (not sure how this happens)
		case header == nil:
			continue
		}

		// the target location where the dir/file should be created
		target := filepath.Join(dst, header.Name)

		// check the file type
		switch header.Typeflag {

		// if it's a dir and it doesn't exist create it
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					return err
				}
			}

		// if it's a file create it
		case tar.TypeReg:
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			defer f.Close()

			// copy over contents
			if _, err := io.Copy(f, tr); err != nil {
				return err
			}
		}
	}
}

func main() {
	var c, x, v bool

	flag.BoolVar(&c, "c", false, "create archive mode")
	flag.BoolVar(&x, "x", false, "extract archive mode")
	flag.BoolVar(&v, "v", false, "print version")
	flag.Parse()

	if v {
		fmt.Println("Simple TAR archiver")
		fmt.Println("version: " + version)
		return
	}

	if (c && x) || (!c && !x) {
		log.Fatal("one and only one of create or extract archive mode")
	}

	switch {
	case c:
		if flag.NArg() != 1 {
			log.Fatal("when create tar, specify just one directory as argument")
		}
		stdoutIsatty := isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd())
		if stdoutIsatty {
			log.Fatal("when create tar, stdout should not be a tty")
		}
		if err := create(os.Stdout, flag.Arg(0)); err != nil {
			log.Fatalf("error occured when creating tar: %v", err)
		}

	case x:
		var dst string
		stdinIsatty := isatty.IsTerminal(os.Stdin.Fd()) || isatty.IsCygwinTerminal(os.Stdin.Fd())
		if stdinIsatty {
			log.Fatal("when extract tar, tar data should be inputted from stdin")
		}
		if flag.NArg() == 0 {
			dst = "."
		} else if flag.NArg() == 1 {
			dst = flag.Arg(0)
		} else {
			log.Fatal("too many arguments")
		}

		if err := extract(os.Stdin, dst); err != nil {
			log.Fatalf("error occured when extracting tar: %v", err)
		}
	}
}

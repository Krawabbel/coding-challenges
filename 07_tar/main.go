package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
)

func mustSucceed(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {

	archive := flag.String("file", "", "use archive file or device ARCHIVE")
	doExtract := flag.Bool("extract", false, "extract files from an archive")
	doCreate := flag.Bool("create", false, "create a new archive")

	flag.Parse()

	if *doExtract && *doCreate {
		panic("you may not specify more than one 'extract' or 'create' option")
	}

	if *doExtract {
		mustSucceed(extract(*archive))
	}

	if *doCreate {
		paths := flag.Args()
		mustSucceed(create(*archive, paths))
	}
}

func create(archivepath string, filepaths []string) error {

	tar := make([]byte, 0)

	for _, path := range filepaths {

		raw, err := pack(path)
		if err != nil {
			return err
		}

		tar = append(tar, raw...)
	}

	tar = append(tar, make([]byte, 1024)...)

	return os.WriteFile(archivepath, tar, os.FileMode(0644))
}

func extract(archivename string) error {

	tar, err := os.ReadFile(archivename)
	if err != nil {
		return err
	}

	files, err := unpack(tar)
	if err != nil {
		return err
	}

	for _, f := range files {
		name := f.info.filename()
		perm := fs.FileMode(f.info.fileMode)
		data := f.data

		fmt.Println(name, perm)

		err := os.WriteFile(name, data, perm)
		if err != nil {
			return err
		}

	}

	return nil
}

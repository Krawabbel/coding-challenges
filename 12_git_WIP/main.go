package main

import (
	"fmt"
	"os"
)

const (
	DESCRIPTION = "Unnamed repository; edit this file 'description' to name the repository."

	EXCLUDE = `# git ls-files --others --exclude-from=.git/info/exclude
# Lines that start with '#' are comments.
# For a project mostly in C, the following would be a good set of
# exclude patterns (uncomment them if you want to use them):
# *.[oa]
# *~`

	CONFIG = `[core]
	repositoryformatversion = 0
	filemode = true
	bare = false
	logallrefupdates = true`

	PERMMASK = 0755
)

func cmdInit(args []string) error {
	fmt.Println("args", args)

	if err := os.MkdirAll(".git/objects/info", PERMMASK); err != nil {
		return err
	}

	if err := os.MkdirAll(".git/objects/pack", PERMMASK); err != nil {
		return err
	}

	if err := os.MkdirAll(".git/refs/heads", PERMMASK); err != nil {
		return err
	}

	if err := os.MkdirAll(".git/refs/tags", PERMMASK); err != nil {
		return err
	}

	if err := os.MkdirAll(".git/hooks", PERMMASK); err != nil {
		return err
	}

	if err := os.MkdirAll(".git/info", PERMMASK); err != nil {
		return err
	}

	if err := os.WriteFile(".git/info/exclude", []byte(EXCLUDE), PERMMASK); err != nil {
		return err
	}

	if err := os.WriteFile(".git/HEAD", []byte("ref:refs/heads/main"), PERMMASK); err != nil {
		return err
	}

	if err := os.WriteFile(".git/description", []byte(DESCRIPTION), PERMMASK); err != nil {
		return err
	}

	if err := os.WriteFile(".git/config", []byte(CONFIG), PERMMASK); err != nil {
		return err
	}

	return nil
}

func run(args []string) error {
	cmd := args[0]
	switch cmd {
	case "init":
		return cmdInit(args[1:])
	}

	return fmt.Errorf("unexpected command '%s'", cmd)
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		panic(err)
	}
}

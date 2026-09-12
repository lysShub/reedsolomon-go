package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

// emit appends a step-summary section to $GITHUB_STEP_SUMMARY when set
// (GitHub Actions), otherwise to stdout (local runs).
func emit(kind, name, body string) {
	var w io.Writer = os.Stdout
	if p := os.Getenv("GITHUB_STEP_SUMMARY"); p != "" {
		f, err := os.OpenFile(p, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0o644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "open summary: %v\n", err)
		} else {
			defer f.Close()
			w = f
		}
	}
	if name != "" {
		fmt.Fprintf(w, "## %s: %s\n\n", kind, name)
	} else {
		fmt.Fprintf(w, "## %s\n\n", kind)
	}
	fmt.Fprintf(w, "```\n%s```\n", body)
}

// run runs a command in dir, streaming output like a direct invocation.
func run(dir, name string, args ...string) error {
	fmt.Fprintf(os.Stderr, "+ (cd %s) %s %v\n", dir, name, args)
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// runToFile runs a command in dir, streaming output like a direct invocation
// while also writing stdout to out (relative to the current working directory).
func runToFile(dir, out, name string, args ...string) error {
	fmt.Fprintf(os.Stderr, "+ (cd %s) %s %v > %s\n", dir, name, args, out)
	f, err := os.Create(out)
	if err != nil {
		return err
	}
	defer f.Close()

	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdout = io.MultiWriter(f, os.Stdout)
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

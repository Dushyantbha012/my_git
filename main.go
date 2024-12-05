package main

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"os"
	"path"
	"strings"
)

// Usage: your_program.sh <command> <arg1> <arg2> ...
func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: mygit <command> [<args>...]\n")
		os.Exit(1)
	}

	switch command := os.Args[1]; command {
	case "usage:":
		fmt.Fprintf(os.Stderr, "will implement it\n")
	case "init":
		//Uncomment this block to pass the first stage!

		for _, dir := range []string{".git", ".git/objects", ".git/refs"} {
			if err := os.MkdirAll(dir, 0755); err != nil {
				fmt.Fprintf(os.Stderr, "Error creating directory: %s\n", err)
			}
		}

		headFileContents := []byte("ref: refs/heads/main\n")
		if err := os.WriteFile(".git/HEAD", headFileContents, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing file: %s\n", err)
		}

		fmt.Println("Initialized git directory")
	case "cat-file":
		if len(os.Args) < 4 {
			fmt.Fprintf(os.Stderr, "usage: mygit <command> [<args>...]\n")
			os.Exit(1)
		}
		if os.Args[2] != "-p" {
			fmt.Fprintf(os.Stderr, "error, second params needs to be '-p'\n")
			os.Exit(1)
		}
		dirname := os.Args[3][:2]
		filename := os.Args[3][2:]
		objectPath := path.Join(".git", "objects", dirname, filename)
		fileBytes, err := os.ReadFile(objectPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "couldn't read the file")
			os.Exit(1)
		}
		bytesReader := bytes.NewReader(fileBytes)
		r, err := zlib.NewReader(bytesReader)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error creating zlib new reader")
			os.Exit(1)
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "couldn't read the file")
			os.Exit(1)
		}
		defer r.Close()
		var decompressedData bytes.Buffer
		_, err = decompressedData.ReadFrom(r)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error decompressing data")
			os.Exit(1)
		}
		resultString := decompressedData.String()
		splits := strings.SplitN(resultString, "\x00", 2)
		_, content := splits[0], splits[1]
		fmt.Print(content)
	default:
		fmt.Fprintf(os.Stderr, "Unknown command %s\n", command)
		os.Exit(1)
	}
}

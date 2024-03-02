package main

import (
	"fmt"
	"os"

	multiline_text_replace "github.com/evgeniums/multiline_text_replace/pkg"
	"github.com/jessevdk/go-flags"
)

func main() {

	opts := multiline_text_replace.Options{}

	_, err := flags.Parse(&opts)
	if err != nil {
		os.Exit(1)
	}

	err = multiline_text_replace.ReplaceText(opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
	}

	fmt.Println("Success")
}

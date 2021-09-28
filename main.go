package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"pwdgen/pgen"
	"runtime"
	"strings"

	"github.com/atotto/clipboard"
)

var (
	length       = flag.Int("length", 16, "length of password specified ")
	numDigits    = flag.Int("digits", 4, "number of digits in the password ")
	numSymbols   = flag.Int("symbols", 6, "number of symbols in the password")
	noUppercase  = flag.Bool("no-upper", true, "excludes uppercase letters ")
	allowRepeat  = flag.Bool("allow-repeat", true, "allow repeated characters")
	loops        = flag.Int("copies", 1, "copies of password")
	noClipboard  = flag.Bool("no-clipboard", false, "do not copy to clipboard")
	printVersion = flag.Bool("version", false, "print version and exit ")
	version      = "0.0.1"
	buildTime    = "2021-09-28"
)

const banner = `                                    
______        ______   ____ _____ _   _ 
|  _ \ \      / |  _ \ / ___| ____| \ | |
| |_) \ \ /\ / /| | | | |  _|  _| |  \| |
|  __/ \ V  V / | |_| | |_| | |___| |\  |
|_|     \_/\_/  |____/ \____|_____|_| \_|

Simple password generator, version %s, build %s
`

func main() {
	var passwords []string
	flag.Usage = func() {
		_, _ = fmt.Fprintf(os.Stderr, (fmt.Sprintf(banner, version, buildTime)))
		flag.PrintDefaults()
	}
	flag.Parse()
	// print version and exit
	if *printVersion {
		_, _ = fmt.Fprintf(os.Stderr, (fmt.Sprintf(banner, version, buildTime)))
		os.Exit(0)
	}

	// generate passwords by loop times
	for i := 0; i < *loops; i++ {
		var result, err = pgen.Generate(*length, *numDigits, *numSymbols, *noUppercase, *allowRepeat)
		if err != nil {
			log.Fatal(err)
		}
		passwords = append(passwords, result)
	}
	// @see https://stackoverflow.com/questions/28799110/how-to-join-a-slice-into-a-string
	var result = strings.Join(passwords, "\n")
	// copy password to clipboard
	isInLinuxTerminal := runtime.GOOS == "linux" && os.Getenv("DISPLAY") == ""
	if !isInLinuxTerminal && !*noClipboard {
		if err := clipboard.WriteAll(result); err != nil {
			log.Fatal(err)
		}
	}
	// print generated password
	fmt.Println(result)
}

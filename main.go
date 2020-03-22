package main

import (
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

// available avatar servers:
// https://seccdn.libravatar.org/avatar
// https://secure.gravatar.com/avatar

var appVersion = "dev"
var avatarServer = "https://secure.gravatar.com/avatar"

func main() {
	var showVersion bool
	var email, page string
	flag.BoolVar(&showVersion, "v", false, "Show app version info")
	flag.Parse()

	if showVersion {
		printVersion(os.Stdout)
		return
	}

	paramNum := flag.NArg()
	if paramNum < 1 {
		fmt.Fprintf(os.Stderr, "email param required\n")
		os.Exit(-1)
	} else if paramNum == 2 {
		email = flag.Arg(0)
		page = flag.Arg(1)
	} else if paramNum == 1 {
		email = flag.Arg(0)
	}

	// page param currently not used now, below code keep go compiler not complain
	if false {
		fmt.Printf(page)
	}

	envServer := os.Getenv("CGIT_AVATAR_SERVER")
	envServer = strings.TrimSpace(envServer)
	if envServer != "" &&
		len(envServer) > 4 &&
		(envServer[0:4] == "http" || envServer[0:2] == "//") {
		avatarServer = strings.TrimRight(envServer, "/")
	}

	if input, err := ioutil.ReadAll(os.Stdin); err != nil {
		fmt.Fprintln(os.Stderr, "Error reading from Stdin:", err)
		os.Exit(-1)
	} else {
		mailmd5 := md5hex(email)
		text := strings.TrimSpace(string(input))
		fmt.Fprintf(os.Stdout, `<img src='%s/%s?s=13&amp;d=retro' width='13' height='13' alt='Gravatar' /> %s`,
			avatarServer, mailmd5, text)
	}
}

func md5hex(str string) string  {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func printVersion(w io.Writer) {
	fmt.Fprintf(w, "cgit-email-avatar %s"+
		"\nAvailable at https://github.com/ttys3/cgit-email-avatar \n\n"+
		"Copyright © 2020 荒野無燈 <https://ttys3.net>\n"+
		"Distributed under the Simplified BSD License\n\n", appVersion)
}

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

	// cleanup email, the email maybe `<foo@outlook.com>` format
	email = strings.ToLower(strings.TrimSpace(email))
	if email[0] == '<' {
		email = email[1:]
	}
	if email[len(email)-1] == '>' {
		email = email[0:len(email)-1]
	}

	//fmt.Fprintf(os.Stderr, "got email: %s, page: %s\n", email, page)
	//got email: foo@outlook.com, page: refs
	//got email: foo@outlook.com, page: log

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

		/**
		div#cgit span.libravatar img.onhover {
		        display: none;
		        border: 1px solid gray;
		        padding: 0px;
		        -webkit-border-radius: 4px;
		        -moz-border-radius: 4px;
		        border-radius: 4px;
		        width: 128px;
		        height: 128px;
		}

		div#cgit span.libravatar img.inline {
		        -webkit-border-radius: 3px;
		        -moz-border-radius: 3px;
		        border-radius: 3px;
		        width: 13px;
		        height: 13px;
		        margin-right: 0.2em;
		        opacity: 0.6;
		}

		div#cgit span.libravatar:hover > img.onhover {
		        display: block;
		        position: absolute;
		        margin-left: 1.5em;
		        background-color: #eeeeee;
		        box-shadow: 2px 2px 7px rgba(100,100,100,0.75);
		}
		<span class="libravatar">
			<img class="inline" src="https://seccdn.libravatar.org/avatar/a8907e723e36cf8f37cd45636b1f0af4?s=13&amp;d=retro">
			<img class="onhover" src="https://seccdn.libravatar.org/avatar/a8907e723e36cf8f37cd45636b1f0af4?s=128&amp;d=retro">
		</span>
		*/
		fmt.Fprintf(os.Stdout, `<span class="libravatar">
<img class="inline" src='%s/%s?s=13&amp;d=retro' alt='small avatar' />
<img class="onhover" src='%s/%s?s=128&amp;d=retro' alt='large avatar' />
</span> %s`,
			avatarServer, mailmd5,
			avatarServer, mailmd5,
			text)
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

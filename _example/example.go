package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	"github.com/mattn/go-xmpp"
	"log"
	"os"
	"strings"
)

var server = flag.String("server", "talk.google.com:443", "server")
var username = flag.String("username", "", "username")
var password = flag.String("password", "", "password")
var status = flag.String("status", "xa", "status")
var statusMessage = flag.String("status-msg", "I for one welcome our new codebot overlords.", "status message")
var notls = flag.Bool("notls", false, "No TLS")
var debug = flag.Bool("debug", false, "debug output")
var session = flag.Bool("session", false, "use server session")

func serverName(host string) string {
	return strings.Split(host, ":")[0]
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: example [options]\n")
		flag.PrintDefaults()
		os.Exit(2)
	}
	flag.Parse()
	if *username == "" || *password == "" {
		if *debug && *username == "" && *password == "" {
			fmt.Fprintf(os.Stderr, "no username or password were given; attempting ANONYMOUS auth\n")
		} else if *username != "" || *password != "" {
			flag.Usage()
		}
	}

	if !*notls {
		xmpp.DefaultConfig = tls.Config{
			ServerName:         serverName(*server),
			InsecureSkipVerify: false,
		}
	}

	var talk *xmpp.Client
	var err error
	options := xmpp.Options{Host: *server,
		User:          *username,
		Password:      *password,
		NoTLS:         *notls,
		Debug:         *debug,
		Session:       *session,
		Status:        *status,
		StatusMessage: *statusMessage,
	}

	talk, err = options.NewClient()

	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			chat, err := talk.Recv()
			if err != nil {
				log.Fatal(err)
			}
			switch v := chat.(type) {
			case xmpp.Chat:
				fmt.Println(v.Remote, v.Text)
			case xmpp.Presence:
				fmt.Println(v.From, v.Show)
			case []xmpp.AdhocCommand:
				for _, c := range v {
						fmt.Printf("- JID: %s\n  Name: %s\n  Node: %s\n", c.JID, c.Name, c.Node)				
				}
			case xmpp.Uptime:
				fmt.Printf("Uptime: %s\n", v.Uptime)
			case xmpp.DiscoResult:
				fmt.Println("Features:")
				for _, f := range v.Features {
					fmt.Printf("- %s\n", f)
				}

				fmt.Println("Identities:")
				for _, i := range v.Identities {
					fmt.Printf("- Name: %s\n  Type: %s\n  Category: %s\n", i.Name, i.Type, i.Category)
				}
			}
		}
	}()

	for {
		in := bufio.NewReader(os.Stdin)
		line, err := in.ReadString('\n')
		if err != nil {
			continue
		}
		line = strings.TrimRight(line, "\n")

		switch line {
		case "discoComs":
			// Perform a disco request
			//talk.DiscoverNode(xmpp.XMPPNS_DISCO_COMMANDS)
			talk.AdhocGetCommands()
		case "discoItems":
			// Perform a disco request
			talk.DiscoverItems()
		case "uptime":
			talk.AdhocExecuteCommand("uptime")
		}
		
		// tokens := strings.SplitN(line, " ", 2)
		// if len(tokens) == 2 {
		// 	talk.Send(xmpp.Chat{Remote: tokens[0], Type: "chat", Text: tokens[1]})
		// }
	}
}

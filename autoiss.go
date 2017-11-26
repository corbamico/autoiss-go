package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/andrewstuart/goq"
)

var (
	errParseHTML = errors.New("HTML not match ishadowx.com configuration")
)

// serverConfig show configuration information in ss.ishadowx.com
// * Port is for connect to server ,but  as string "Port : 443" in HTML page
//   ccs selector cannot deal with strings.Split, get it into PortHTML first
// * Method showing as string "Method:aes-256-cfb" in HTML page, get it into
//   MethodHTML first, and then call serverConfig.transform()
type serverConfig struct {
	Address    string `goquery:"h4:nth-child(1) > span[id]"`
	Port       int
	PortHTML   string `goquery:"h4:nth-child(2)"`
	Password   string `goquery:"h4:nth-child(3) > span[id]"`
	Method     string
	MethodHTML string `goquery:"h4:nth-child(4)"`
}

func (s *serverConfig) isValid() bool {
	b := s.Address != "" &&
		s.PortHTML != "" &&
		s.Password != "" &&
		s.MethodHTML != ""
	return b
}

func (s *serverConfig) transform() {
	//note: this is : in chinese ,not in english
	//note: 20171126, ':' changed  to  in english
	s.Port, _ = strconv.Atoi(strings.Split(s.PortHTML, ":")[1])
	s.Method = strings.Split(s.MethodHTML, ":")[1]
}

func main() {
	var localPort int
	var url string
	var debug bool
	var indexNumber int

	log.SetOutput(os.Stdout)

	flag.IntVar(&localPort, "l", 1080, "local socks5 proxy port")
	flag.StringVar(&url, "s", "go.ishadowx.net", "server address")
	flag.BoolVar(&debug, "d", false, "print debug message")
	flag.IntVar(&indexNumber, "n", 0, "which shadowsocks server to use(default 0)\n\t0:first one,-1:last one\n\twe use last one if index exceeds max number")

	flag.Parse()

	server, err := getServerConfig("https://"+url, indexNumber)

	if err != nil {
		log.Fatal("[autoiss-go] Failed to get shadowsocks server:", err)
	}
	runSS(server, localPort, debug)
}

func runSS(s serverConfig, localPort int, debug bool) {
	cmdStr := fmt.Sprintf("-s %s -p %d -k %s -l %d -m %s",
		s.Address,
		s.Port,
		s.Password,
		localPort,
		s.Method)
	if debug {
		cmdStr = "-d " + cmdStr
	}
	log.Println("[autoiss-go] shadowsocks-local " + cmdStr)
	cmd := exec.Command("shadowsocks-local", strings.Fields(cmdStr)...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	//err := cmd.Run()

	if err != nil {
		log.Fatal("[autoiss-go] shadowsocks-local failed, error as :", err)
	}
	time.Sleep(1 * time.Second)
	log.Println("[autoiss-go] shadowsocks-local running...")
	cmd.Wait()
}

func getServerConfig(url string, index int) (serverConfig, error) {
	var server serverConfig

	doc, err := goquery.NewDocument(url)
	if err != nil {
		return server, err
	}

	//Notes 20171126
	//      class="portfolio-items isotope"
	//      but I cannot find exact dot syntac present this class
	//s := doc.Find(".portfolio-items .isotope .portfolio-item")
	s := doc.Find(".row .portfolio-item")

	//return err, if we got empty server info
	if len(s.Nodes) == 0 {
		return server, errParseHTML
	}

	//we use first one if index < -len(s.Nodes)
	//we use last one if index > len(s.Nodes)-1
	index = min(index, len(s.Nodes)-1)
	index = max(index, -len(s.Nodes))

	err = goq.UnmarshalSelection(s.Eq(index), &server)
	if !server.isValid() {
		err = errParseHTML
	} else {
		server.transform()
	}
	return server, err
}

func min(a int, b int) (min int) {
	if a <= b {
		min = a
	} else {
		min = b
	}
	return
}
func max(a int, b int) (max int) {
	if a <= b {
		max = b
	} else {
		max = a
	}
	return
}

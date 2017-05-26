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
	s.Port, _ = strconv.Atoi(strings.Split(s.PortHTML, "：")[1]) //note: this is : in chinese ,not in english
	s.Method = strings.Split(s.MethodHTML, ":")[1]
}

func main() {
	var localPort int
	var url string

	log.SetOutput(os.Stdout)

	flag.IntVar(&localPort, "l", 1080, "local socks5 proxy port")
	flag.StringVar(&url, "s", "ss.ishadowx.com", "server address")

	flag.Parse()

	server, err := getServerConfig("http://" + url)

	if err != nil {
		log.Fatal("[autoiss-go] Failed to get shadowsocks server:", err)
	}
	runSS(server, localPort)
}

func runSS(s serverConfig, localPort int) {
	cmdStr := fmt.Sprintf("-s %s -p %d -k %s -l %d -m %s",
		s.Address,
		s.Port,
		s.Password,
		localPort,
		s.Method)

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

func getServerConfig(url string) (serverConfig, error) {
	var server serverConfig

	doc, err := goquery.NewDocument(url)
	if err != nil {
		return server, err
	}
	p := doc.Find(".portfolio-items .portfolio-item").First()

	err = goq.UnmarshalSelection(p, &server)
	if !server.isValid() {
		err = errParseHTML
	} else {
		server.transform()
	}
	return server, err
}

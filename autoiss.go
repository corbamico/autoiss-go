package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"time"

	"github.com/PuerkitoBio/goquery"
)

type serverConfig struct {
	serverAddress  string
	serverPort     int
	serverPassword string
	serverMethod   string
}

func main() {
	var localPort = 1080
	var iShadowServer = "ss.ishadowx.com"

	log.SetOutput(os.Stdout)

	flag.IntVar(&localPort, "l", 1080, "local socks5 proxy port")
	flag.StringVar(&iShadowServer, "s", "ss.ishadowx.com", "server address")

	flag.Parse()

	server := getServerConfig(iShadowServer)

	if server.serverPassword != "" {
		runSS(server, localPort)
	}

	log.Fatal("[autoiss-go] Failed to get shadowsocks server info.")
}

func runSS(s serverConfig, localPort int) {
	cmdStr := fmt.Sprintf("-s %s -p %d -k %s -l %d -m %s",
		s.serverAddress,
		s.serverPort,
		s.serverPassword,
		localPort,
		s.serverMethod)

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

func getServerConfig(url string) serverConfig {
	var server serverConfig
	req, err := http.NewRequest("GET", "http://"+url, nil)
	if err != nil {
		log.Fatal("创建请求时出错错误", err)
	}
	//req.Header.Set("User-Agent", "Mozilla/5.0 (Arch Linux kernel 4.6.5) AppleWebKit/537.36 (KHTML, like Gecko) Maxthon/4.0 Chrome/39.0.2146.0 Safari/537.36")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal("发起请求时出错错误", err)
	}
	p, err := goquery.NewDocumentFromResponse(res)
	if err != nil {
		log.Fatal("解析html时出现错误", err)
	}

	divs := p.Find(".portfolio-items .portfolio-item")
	for i := range divs.Nodes {
		div := divs.Eq(i)
		h4 := div.Find("h4")
		serverPort, _ := strconv.Atoi(strings.Split(h4.Eq(1).Text(), "：")[1])
		serverAddress := h4.Eq(0).Find("span[id]").First().Text()

		server = serverConfig{
			serverAddress:  serverAddress,
			serverPort:     serverPort,
			serverPassword: h4.Eq(2).Find("span[id]").First().Text(),
			serverMethod:   strings.Split(h4.Eq(3).Text(), ":")[1],
		}
		//return first one
		return server

	}
	return server
}

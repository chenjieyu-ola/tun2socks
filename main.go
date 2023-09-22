package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"go.uber.org/automaxprocs/maxprocs"
	"gopkg.in/yaml.v3"

	_ "github.com/xjasonlyu/tun2socks/v2/dns"
	"github.com/xjasonlyu/tun2socks/v2/engine"
	"github.com/xjasonlyu/tun2socks/v2/internal/version"
	"github.com/xjasonlyu/tun2socks/v2/log"
)

var (
	key = new(engine.Key)

	configFile  string
	versionFlag bool
)

func init() {
	flag.IntVar(&key.Mark, "fwmark", 0, "Set firewall MARK (Linux only)")
	flag.IntVar(&key.MTU, "mtu", 0, "Set device maximum transmission unit (MTU)")
	flag.DurationVar(&key.UDPTimeout, "udp-timeout", 0, "Set timeout for each UDP session")
	flag.StringVar(&configFile, "config", "", "YAML format configuration file")
	flag.StringVar(&key.Device, "device", "", "Use this device [driver://]name")
	flag.StringVar(&key.Interface, "interface", "", "Use network INTERFACE (Linux/MacOS only)")
	flag.StringVar(&key.LogLevel, "loglevel", "info", "Log level [debug|info|warning|error|silent]")
	flag.StringVar(&key.Proxy, "proxy", "", "Use this proxy [protocol://]host[:port]")
	flag.StringVar(&key.RestAPI, "restapi", "", "HTTP statistic server listen address")
	flag.StringVar(&key.TCPSendBufferSize, "tcp-sndbuf", "", "Set TCP send buffer size for netstack")
	flag.StringVar(&key.TCPReceiveBufferSize, "tcp-rcvbuf", "", "Set TCP receive buffer size for netstack")
	flag.BoolVar(&key.TCPModerateReceiveBuffer, "tcp-auto-tuning", false, "Enable TCP receive buffer auto-tuning")
	flag.StringVar(&key.MulticastGroups, "multicast-groups", "", "Set multicast groups, separated by commas")
	flag.StringVar(&key.TUNPreUp, "tun-pre-up", "", "Execute a command before TUN device setup")
	flag.StringVar(&key.TUNPostUp, "tun-post-up", "", "Execute a command after TUN device setup")
	flag.BoolVar(&versionFlag, "version", false, "Show version and then quit")
	flag.Parse()
}

func main() {
	maxprocs.Set(maxprocs.Logger(func(string, ...any) {}))

	if versionFlag {
		fmt.Println(version.String())
		fmt.Println(version.BuildString())
		os.Exit(0)
	}

	if configFile != "" {
		data, err := os.ReadFile(configFile)
		if err != nil {
			log.Fatalf("Failed to read config file '%s': %v", configFile, err)
		}
		if err = yaml.Unmarshal(data, key); err != nil {
			log.Fatalf("Failed to unmarshal config file '%s': %v", configFile, err)
		}
	}

	engine.Insert(key)

	engine.Start()
	defer engine.Stop()

	err := exec.Command("ifconfig", key.Device, "198.18.0.1", "198.18.0.1", "up").Run()
	if err == nil {
		exec.Command("route", "add", "-net", "1.0.0.0/8", "198.18.0.1").Run()
		exec.Command("route", "add", "-net", "2.0.0.0/7", "198.18.0.1").Run()
		exec.Command("route", "add", "-net", "4.0.0.0/6", "198.18.0.1").Run()
		exec.Command("route", "add", "-net", "8.0.0.0/5", "198.18.0.1").Run()
		exec.Command("route", "add", "-net", "16.0.0.0/4", "198.18.0.1").Run()
		exec.Command("route", "add", "-net", "32.0.0.0/3", "198.18.0.1").Run()
		exec.Command("route", "add", "-net", "64.0.0.0/2", "198.18.0.1").Run()
		exec.Command("route", "add", "-net", "128.0.0.0/1", "198.18.0.1").Run()
		exec.Command("route", "add", "-net", "32.0.0.0/3", "198.18.0.1").Run()
		exec.Command("route", "add", "-net", "198.18.0.0/15", "198.18.0.1").Run()
	} else {
		fmt.Println(err)
		return
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
}

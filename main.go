package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Port int `yaml:"port"`
	} `yaml:"server"`
	CIDR struct {
		File string `yaml:"file"`
	} `yaml:"cidr"`
}

type IPChecker struct {
	networks []*net.IPNet
}

func NewIPChecker(cidrFile string) (*IPChecker, error) {
	file, err := os.Open(cidrFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open CIDR file: %v", err)
	}
	defer file.Close()

	var networks []*net.IPNet
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		_, network, err := net.ParseCIDR(scanner.Text())
		if err != nil {
			return nil, fmt.Errorf("failed to parse CIDR %s: %v", scanner.Text(), err)
		}
		networks = append(networks, network)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading CIDR file: %v", err)
	}

	return &IPChecker{networks: networks}, nil
}

func (ic *IPChecker) Contains(ipStr string) bool {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}

	for _, network := range ic.networks {
		if network.Contains(ip) {
			return true
		}
	}
	return false
}

func main() {
	// Load configuration
	config := &Config{}
	configFile, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	err = yaml.Unmarshal(configFile, config)
	if err != nil {
		log.Fatalf("Error parsing config file: %v", err)
	}

	// Initialize IP checker
	checker, err := NewIPChecker(config.CIDR.File)
	if err != nil {
		log.Fatalf("Error initializing IP checker: %v", err)
	}

	// Setup Gin router
	r := gin.Default()

	r.GET("/check/:ip", func(c *gin.Context) {
		ip := c.Param("ip")
		result := checker.Contains(ip)
		c.JSON(200, gin.H{
			"ip":      ip,
			"matches": result,
		})
	})

	// Start server
	addr := fmt.Sprintf(":%d", config.Server.Port)
	log.Printf("Starting server on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

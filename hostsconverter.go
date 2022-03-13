package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"sort"
	"strings"
)

type Hosts []string

func isValidIpV4(addr string) bool {
	return net.ParseIP(addr) != nil
}

func isSpecialRecord(addr string, name string) bool {
	// ip := []string{
	// 	"127.0.0.1",
	// 	"255.255.255.255",
	// 	"::1",
	// 	"FE80::1",
	// 	"FF00::0",
	// 	"FF02::1",
	// 	"FF02::2",
	// 	"FF02::3",
	// 	"0.0.0.0",
	// }

	nm := []string{
		"LOCALHOST",
		"LOCALHOST.LOCALDOMAIN",
		"LOCAL",
		"BROADCASTHOST",
		"IP6-LOCALHOST",
		"IP6-LOOPBACK",
		"IP6-LOCALNET",
		"IP6-MCASTPREFIX",
		"IP6-ALLNODES",
		"IP6-ALLROUTERS",
		"IP6-ALLHOSTS",
	}

	for _, i := range nm {
		if i == name {
			return true
		}
	}
	return false
}

func LoadHostsFromFile(filename string) Hosts {
	log.Printf("Reading file %s\n", filename)
	var res Hosts
	f, err := os.Open(filename)
	if err != nil {
		return res
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "#") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		ip := strings.ToUpper(fields[0])
		var name string
		name = strings.ToUpper(fields[1])

		if strings.Contains(name, "#") {
			// после имени через # идет коментарий
			log.Printf("[!] name with comment: %s\n", line)

			ss := strings.Split(name, "#")
			name = ss[0]
		}

		if strings.Contains(name, ":") {
			// для случаев типа 199.7.52.190:80
			log.Printf("[!] ip with port: %s\n", line)
			ss := strings.Split(name, ":")
			name = ss[0]
		}

		if ip == name {
			// если ip и имя одинаковые
			log.Printf("[!] ip==name: %s\n", line)
			continue
		}

		if isValidIpV4(name) {
			// если в имение вписан другой ip
			log.Printf("[!] name is ip: %s\n", line)
			continue
		}

		if isSpecialRecord(ip, name) {
			log.Printf("[!] special record: %s\n", line)
			continue
		}

		if !strings.Contains(name, ".") {
			log.Printf("[!] name without dot: %s\n", line)
			continue
		}
		res = append(res, name)
	}

	return res
}

func LoadHosts(files []string) Hosts {
	var res Hosts
	log.Println("Loading hosts...")

	for _, filename := range files {
		hosts := LoadHostsFromFile(filename)
		res = append(res, hosts...)
	}

	return res
}

func (hosts Hosts) SaveHosts(filename string) {
	log.Printf("Saving hosts (%d records)...", len(hosts))
	f, err := os.Create(filename)
	defer f.Close()
	if err != nil {
		return
	}
	writer := bufio.NewWriter(f)
	for _, h := range hosts {
		fmt.Fprintf(writer, "127.0.0.1 %s\n", strings.ToLower(h))
	}
	writer.Flush()
}

func (hosts Hosts) FilterHosts(filters Hosts) Hosts {
	var res Hosts
	log.Printf("Filtering hosts (%d records)...\n", len(hosts))
	for _, h := range hosts {
		if !filters.HostExists(h) {
			res = append(res, h)
		}
	}
	return res
}

func (hosts Hosts) AllowHosts(hostsAllow Hosts) Hosts {
	var res Hosts
	log.Printf("Forcing allow hosts (%d records)...\n", len(hostsAllow))
	for _, h := range hosts {
		if !hostsAllow.HostExists(h) && !IsHostNumeric(h) {
			res = append(res, h)
		}
		if IsHostNumeric(h) {
			log.Printf("Why IP in list %s\n", h)
		}
	}
	return res
}

func isNumeric(c rune) bool {
	return '0' <= c && c <= '9'
}

func IsHostNumeric(host string) bool {
	s := strings.ReplaceAll(host, ".", "")
	for _, c := range s {
		if !isNumeric(c) {
			return false
		}
	}
	return true
}

func (hosts Hosts) HostExists(hostToCheck string) bool {
	for _, h := range hosts {
		if strings.ToUpper(h) == strings.ToUpper(hostToCheck) {
			return true
		}
	}
	return false
}

func (hosts Hosts) Sort() Hosts {
	tmp := hosts
	sort.Strings(tmp)
	return tmp
}

func (hosts Hosts) RemoveDuplicates() Hosts {
	//var newHosts Hosts
	newHosts := (Hosts)(make([]string, 0, len(hosts)))
	hostsMap := make(map[string]bool)
	log.Printf("Removing duplicates (%d records)...\n", len(hosts))
	for _, h := range hosts {
		H := strings.ToUpper(h)
		if _, ok := hostsMap[H]; !ok {
			hostsMap[H] = true
			newHosts = append(newHosts, h)
		}
	}
	return newHosts
}

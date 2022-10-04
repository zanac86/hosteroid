package main

import (
	"fmt"
	"log"
)

var (
	Hosts_Always = []string{
		"localhost",
		"localhost.localdomain",
		"broadcasthost",
		"local",
	}
)

// SystemHostFilename return path for save results
func SystemHostFilename() string {
	/*
		Windows
		"c:\\windows\\system32\\drivers\\etc\\hosts"
	*/
	return "_hosts_"
}

func main() {
	log.Println("Hosteroid")

	var exts []string
	for i := 0; i < 1000; i++ {
		exts = append(exts, fmt.Sprintf(".%d", i))
	}

	files := ListFiles(".", exts, false)

	hosts := LoadHosts(files)

	hostsAllow := LoadHostsFromFile("hosts.allow")

	hosts2 := hosts.FilterHosts(Hosts_Always).RemoveDuplicates().AllowHosts(hostsAllow).Sort()

	hosts2 = append(hosts2, Hosts_Always...)

	hosts2.SaveHosts(SystemHostFilename())

	log.Println("Have a nice day!")
}

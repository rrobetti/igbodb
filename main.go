package main

import "fmt"
import "igbodb/leader-writter"
import "igbodb/load-manager"
import "igbodb/read-replica"

func main() {
	fmt.Println("From main function..")
	leaderwritter.Print()
	loadmanager.Print()
	readreplica.Print()
}

// Simple TCP port scanner with multiple goroutines which scans local network(255.255.255.0 subnet mask hardcoded) for an opened port by specified number. 
// It was made as a pre-employement test.
package main
import (
    "log"
    "net"
	"fmt"
	"flag"
    "strconv"
	"time"
)
func GetOutboundIP() net.IP { // There should be a better way to do that, but I haven't found it out just yet.
    conn, err := net.Dial("udp", "8.8.8.8:80")
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()
    localAddr := conn.LocalAddr().(*net.UDPAddr)
    return localAddr.IP
}
func ScanRoutine(myip net.IP, port int, i int){
	ip := make(net.IP, len(myip))
	copy(ip, myip)
	ip[3] = byte(i); // Assuming 255.255.255.0 subnet mask.
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%v:%d", ip, port), 9 * time.Second)
	if err == nil {
		conn.Close()
		log.Printf("port %d is opened at %v",port, ip)
	} else {
		log.Printf("port %d is closed at %v",port, ip)
	}
}
func main() {
	flag.Parse()
	myip := GetOutboundIP()
	for _, arg := range flag.Args() { // Added just so we can scan multiple ports at the same time, not sure if it's really necessary.
		port, err := strconv.Atoi(arg)
	    if( err != nil || port < 1 || port > 65535 ){
			log.Printf("wrong port: %q", arg)
	        continue
	    }
		for i := 1; i < 255; i++{ // Assuming 255.255.255.0 subnet mask
			go ScanRoutine(myip, port, i)
		}
	}
	time.Sleep(10 * time.Second) // Not the best way to let the program finish it's goroutines, but the goal was to finish scanning in 10 secs so let it be 10 sec waiting.
}
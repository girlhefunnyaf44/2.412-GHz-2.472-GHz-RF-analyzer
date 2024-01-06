package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"github.com/fatih/color"
)
/* I used some light obfuscation techniques hence the function names, so this program is hard to read.

 At the time of writing this I was really interested in cyber security, CTF and reverse engineering. This reflects my coding style, because of how obfuscated I made it.

 The program tries to confuse you and make you scratch your head, but I added documentation so it's easier to read.*/


func rainbow(title string) (int, error) {  // This function uses winapi to change the title of the window. 
	handle, err := syscall.LoadLibrary("Kernel32.dll")
	if err != nil {
		return 0, err
	}
	defer syscall.FreeLibrary(handle)
	proc, err := syscall.GetProcAddress(handle, "SetConsoleTitleW")
	if err != nil {
		return 0, err
	}
	r, _, err := syscall.Syscall(proc, 1, uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(title))), 0, 0)
	return int(r), err

}

func purple(stdin string, ports, stdout chan int) {  // This function checks if a port is open
	for quux := range ports {
		domain := fmt.Sprintf("%s:%d", stdin, quux)
		socket, err := net.Dial("tcp", domain)
		if err != nil {
			stdout <- 0
			continue
		}
		socket.Close()
		stdout <- quux
	}
}

func black() {  // This function creates goroutines that run the purple function to check if the port is open, then relays the information via a channel. It scans port 1-1024
	var stdin string
	fmt.Print("Enter a naked Domain: ")
	fmt.Scanln(&stdin)
	ports := make(chan int, 100)
	stdout := make(chan int)
	var baz []int

	for i := 0; i < cap(ports); i++ {
		go purple(stdin, ports, stdout)
	}
	go func() {
		for i := 1; i <= 1024; i++ {
			ports <- i
		}
		close(ports)
	}()
	for i := 0; i < 1024; i++ {
		port := <-stdout
		if port != 0 {
			baz = append(baz, port)
		}
	}
	close(stdout)
	sort.Ints(baz)
	for _, port := range baz {
		fmt.Printf("%d open\n", port)
	}
}

func yellow(input string) string { // This function just scrapes tasklist.
	bananaArr := []string{"Image Name", "Session Name", ",", "AnyDesk.exe", "Services", "Console", "Session#", "Mem Usage", "PID"}
	for _, s := range bananaArr {
		input = strings.Replace(input, s, "", -1)
	}
	return input
}
func light_blue(input string) string { //regex pattern that looks for ext in double quotes.
	pattern := `"[^"]\b[^"]*"`
	re := regexp.MustCompile(pattern)
	return re.ReplaceAllString(input, "")
}
func light_green(input string) string { // This function searches for matches containing double quotes and replaces them
	re := regexp.MustCompile(`"`)
	return re.ReplaceAllString(input, "")
}

func light_pink(magenta string) (string, string, string, string) { // The program Anydesk had multiple instances of itself so searching for the name "anydesk" isn't enough. you have to sort through the Process ID's
	pidArray := strings.Fields(magenta)
	if len(pidArray) >= 4 {
		pid1, pid2, pid3, pid4 := pidArray[0], pidArray[1], pidArray[2], pidArray[3]
		return pid1, pid2, pid3, pid4
	} else {
		fmt.Println("Sucks to be you.")
		time.Sleep(3 * time.Second)
		os.Exit(0)
	}
	return "", "", "", "" // Scoping issue, returns a bunch of empty strings.
}

func red() (string, string, string, string) { // Takes 0 parameters and returns 4 strings
	errors := func() {  // Lambda for Lazy error handling
		fmt.Println("Anydesk is not currently running.")
		time.Sleep(4 * time.Second)
		os.Exit(0)
	}
	cmd := exec.Command("tasklist", "/fo", "csv", "/fi", "IMAGENAME eq Anydesk.exe") // Tasklist windows command that looks for "Anydesk.exe"
	results, err := cmd.Output() 
	cmd_output := string(results) // Converts the cmd byte slice to string
	if err != nil || strings.Count(cmd_output, "AnyDesk") <= 3 { // Checks if anydesk is running
		errors() // idiosyncrasy
	}
	result := yellow(cmd_output)
	
/* This calls the function yellow and takes the cmd output as a parameter. 
The Yellow function just scrapes the output to make it easier to find Anydesk. */

	
	orange := func(blue string) string { 
		pattern := `"[^"]*K[^"]*"`  // This checks for a string that contains the char K and is surrounded by quotes.
		re := regexp.MustCompile(pattern)
		return re.ReplaceAllString(blue, "")
	}(result) // This entire function takes the parsed cmd output and runs a regedit command to parse it even more.

	light_red := light_blue(orange)  /* These two lines are just functions that run regex commands to further narrow the results, since there was a lot. */
	white := light_green(light_red) 
	Mot, Hai, Ba, Bon := light_pink(white)  // The version of Anydesk I created this for should have 4 PID's. light_pink scans each PID with a netstat command to find out which one contains the subjects IP.
	return Mot, Hai, Ba, Bon // Vietnamese metasyntatic variable names
}
func dark_red(pid string) { // This runs a netstat command on each PID to find the IP address
	cmd := exec.Command("cmd.exe", "/C", "netstat", "-p", "TCP", "-n", "-a", "-o", "|", "findstr", pid) 
	output, err := cmd.Output()
	indigo := string(output)
	dark_grey := func() {
		switch err {
		case nil:
			fmt.Printf("Process scanned with PID %s successfully\n", pid)
			indigo = string(output)
		}
	}
	dark_grey()
	pattern := fmt.Sprintf(`(?m)^.*%s.*$`, "SYN_SENT") // More regex commands to parse then followed up by function calls
	re := regexp.MustCompile(pattern)
	matches := re.FindAllString(indigo, -1)
	filtered := strings.Join(matches, "\n")
	fmt.Println(filtered)
	fmt.Println("The ignore the 1st number the 2nd number is the IP Address.")
	time.Sleep(1 * time.Minute)
	Mot, Hai, Ba, Bon := red()
	dark_red(Mot)
	dark_red(Hai)
	dark_red(Ba)
	dark_red(Bon)

}

func main() {
	func() { // This lambda function uses a random seed to generate a number 1-10 then uses a syscall to change the programs name. "SetWindowTextA function"
		rand.Seed(time.Now().UnixNano()) 
		number := rand.Intn(10)
		ascii := fmt.Sprintf("%d", number)
		fmt.Println(ascii)
		rainbow(ascii) 
	}()
	amg := "" +
		"        ⠀ ⠀⠀⠀⠀⠀⠀⠀  ⠀⠀⠀⣠⣤⣤⣤⣤⣤⣶⣦⣤⣄⡀⠀⠀⠀⠀⠀⠀⠀⠀\n" +
		"         ⠀ ⠀⠀ ⠀ ⠀⠀⠀⢀⣴⣿⡿⠛⠉⠙⠛⠛⠛⠛⠻⢿⣿⣷⣤⡀⠀⠀⠀⠀⠀\n" +
		"          ⠀⠀⠀ ⠀⠀⠀⠀⠀⣼⣿⠋⠀⠀⠀⠀⠀⠀⠀⢀⣀⣀⠈⢻⣿⣿⡄⠀⠀⠀⠀\n" +
		"⠀        ⠀ ⠀⠀⠀⠀ ⠀⣸⣿⡏⠀⠀⠀⣠⣶⣾⣿⣿⣿⠿⠿⠿⢿⣿⣿⣿⣄⠀⠀⠀\n" +
		"⠀         ⠀⠀⠀ ⠀⠀⠀⣿⣿⠁⠀⠀⢰⣿⣿⣯⠁⠀⠀⠀⠀⠀⠀⠀⠈⠙⢿⣷⡄⠀\n" +
		"⠀         ⠀⣀⣤⣴⣶⣶⣿⡟⠀⠀⠀⢸⣿⣿⣿⣆⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣿⣷⠀\n" +
		"⠀        ⢰⣿⡟⠋⠉⣹⣿⡇⠀⠀⠀⠘⣿⣿⣿⣿⣷⣦⣤⣤⣤⣶⣶⣶⣶⣿⣿⣿⠀\n" +
		"        ⠀⢸⣿⡇⠀⠀⣿⣿⡇⠀⠀⠀⠀⠹⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡿⠃⠀\n" +
		"⠀        ⣸⣿⡇⠀⠀⣿⣿⡇⠀⠀⠀⠀⠀⠉⠻⠿⣿⣿⣿⣿⡿⠿⠿⠛⢻⣿⡇⠀⠀\n" +
		"⠀        ⣿⣿⠁⠀⠀⣿⣿⡇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢸⣿⣧⠀⠀\n" +
		"⠀        ⣿⣿⠀⠀⠀⣿⣿⡇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢸⣿⣿⠀⠀\n" +
		"⠀        ⣿⣿⠀⠀⠀⣿⣿⡇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢸⣿⣿⠀⠀\n" +
		"        ⠀⢿⣿⡆⠀⠀⣿⣿⡇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣿⣿⠃⠀⠀\n" +
		"⠀       ⠀⠛⢿⣿⣿⣿⣿⣇⠀⠀⠀⠀⠀⣰⣿⣿⣷⣶⣶⣶⣶⠶⠀⢠⣿⣿⠀⠀⠀\n" +
		"⠀⠀⠀        ⠀⠀⠀⠀⣿⣿⠀⠀⠀⠀⠀⣿⣿⡇⠀⣽⣿⡏⠁⠀⠀⢸⣿⡇⠀⠀⠀\n" +
		"        ⠀⠀⠀⠀⠀⠀⠀⣿⣿⠀⠀⠀⠀⠀⣿⣿⡇⠀⢹⣿⡆⠀⠀⠀⣸⣿⠇⠀⠀⠀\n" +
		"        ⠀⠀⠀⠀⠀⠀⠀⢿⣿⣦⣄⣀⣠⣴⣿⣿⠁⠀⠈⠻⣿⣿⣿⣿⡿⠏⠀⠀⠀⠀\n" +
		"        ⠀⠀⠀⠀⠀⠀⠀⠈⠛⠻⠿⠿⠿⠿⠋⠁⠀⠀⠀⠀⠀⠀⠀⠀"
	Color := color.New(color.FgRed).Add(color.Bold)
	Color.Println(amg)
	Colour := color.New(color.FgHiCyan)
	fmt.Println("")
	color.HiBlack("•——————————————————————————•°•✿•°•——————————————————————————•")
	Colour.Println("   Type 1 for Anydesk Resolver, Type 2 for TCP Port Scanner")
	color.HiBlack("•——————————————————————————•°•✿•°•——————————————————————————•")
	fmt.Print("Enter Here:")
	scammer := bufio.NewScanner(os.Stdin) // buffered i/o
	scammer.Scan()
	input, _ := strconv.ParseInt(scammer.Text(), 10, 64) 
	switch input {
	case 1:
		Mot, Hai, Ba, Bon := red()
		dark_red(Mot)
		dark_red(Hai)
		dark_red(Ba)
		dark_red(Bon)
	case 2:
		black()
	default:
		fmt.Println("Invalid Input.")
		time.Sleep(3 * time.Second)
		os.Exit(1)
	}
}

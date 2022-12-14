package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	f, _ := os.OpenFile("./log.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	defer f.Close()
	t, _ := os.OpenFile("./total.txt", os.O_WRONLY|os.O_CREATE, 0644)
	defer t.Close()
	onLeave := false

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Please enter the logger's purpose: ")
	text, _ := reader.ReadString('\n')
	text = strings.TrimSuffix(text, "\n")
	fmt.Println("Great. Your", text, "logger will start soon.")
	_, err := f.Write([]byte("This is your " + text + " logger. \n"))
	if err != nil {
		fmt.Println(err)
	}
	temp := time.Now()
	var totalDuration time.Duration

	for {
		if !onLeave {
			fmt.Println("Welcome back!")
			totalDuration += time.Since(temp)
			resStr := "Recorded time: " + time.Now().String() + " Time logged: " +
				time.Since(temp).String() + " Total Duration: " + totalDuration.String() + "\n"
			f.Write([]byte(resStr))
			os.Truncate("./total.txt", 0)
			fmt.Println("Your total logged time for " + text + " is " + totalDuration.String())
		} else {
			temp = time.Now()
		}
		fmt.Println("Please enter anything :)")
		_, _ = reader.ReadString('\n')
		fmt.Println("Entered at ", time.Now())
		onLeave = !onLeave
	}
}

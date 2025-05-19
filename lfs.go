package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"strings"
)

var TEAL string = "\033[38;5;6m"
var SKYBLUE string = "\033[38;5;74m"
var RESET string = "\033[0m"

func shellinf() (string, string, string) {
	/* UN means username */
	un, _ := user.Current()
	/* WD is present working directory */
	wd, _ := os.Getwd()
	/* HN is hostname */
	hn, _ := os.Hostname()
	return un.Username, wd, hn
}

func main() {
	/* shell loop */
	shloop()
}

func shloop() {
	readInput := bufio.NewReader(os.Stdin)
	fmt.Printf("Starting My Intentionally Less Friendly Shell\n")

	var status int
	for status != 1 {
		un, wd, hn := shellinf()
		fmt.Printf(TEAL+"%v"+SKYBLUE+"@%v>"+RESET+" %v >> ", un, hn, wd)

		input := CleanInput(readInput)
		pipes := ParseInput(input)
		if pipes == true {
			status = PipeHandling(input)
		} else {
			status = HandlingInput(input)
		}
	}
}

/* this function removes the newline from the end of the input */
func CleanInput(reader *bufio.Reader) string {
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Error getting input to clean it: %v \n", err)
	}
	input = strings.TrimSuffix(input, "\n")

	return input
}

func ParseInput(input string) bool {
	pipe := strings.Contains(input, "|")
	if pipe == true {
		return true
	}
	return false
}

func twopiper(twodPipe [][]string) {
	pipeEx1 := exec.Command(twodPipe[0][0], twodPipe[0][1:]...)
	pipeEx2 := exec.Command(twodPipe[1][0], twodPipe[1][1:]...)

	pipeEx2.Stdin, _ = pipeEx1.StdoutPipe()
	pipeEx2.Stdout = os.Stdout
	pipeEx2.Stderr = os.Stderr

	err := pipeEx2.Start()
	_ = pipeEx1.Run()
	pipeEx2.Wait()
	if err != nil {
		fmt.Printf("Error piping: %v\n", pipeEx2)
	}

}

func threepiper(twodPipe [][]string) {
	pipex1 := exec.Command(twodPipe[0][0], twodPipe[0][1:]...)
	pipex2 := exec.Command(twodPipe[1][0], twodPipe[1][1:]...)
	pipex3 := exec.Command(twodPipe[2][0], twodPipe[2][1:]...)

	pipex2.Stdin, _ = pipex1.StdoutPipe()
	pipex3.Stdin, _ = pipex2.StdoutPipe()
	pipex3.Stdout = os.Stdout
	pipex3.Stderr = os.Stderr

	err := pipex3.Start()
	pipex2.Start()
	pipex1.Run()
	pipex2.Wait()
	pipex3.Wait()
	if err != nil {
		fmt.Printf("Error piping: %v \n", err)
	}
}

func fourpiper(twodPipe [][]string) {
	pipex1 := exec.Command(twodPipe[0][0], twodPipe[0][1:]...)
	pipex2 := exec.Command(twodPipe[1][0], twodPipe[1][1:]...)
	pipex3 := exec.Command(twodPipe[2][0], twodPipe[2][1:]...)
	pipex4 := exec.Command(twodPipe[3][0], twodPipe[3][1:]...)

	pipex2.Stdin, _ = pipex1.StdoutPipe()
	pipex3.Stdin, _ = pipex2.StdoutPipe()
	pipex4.Stdin, _ = pipex3.StdoutPipe()
	pipex4.Stdout = os.Stdout
	pipex4.Stderr = os.Stderr

	err := pipex4.Start()
	pipex2.Start()
	pipex3.Start()
	pipex1.Run()
	pipex2.Wait()
	pipex3.Wait()
	pipex4.Wait()
	if err != nil {
		fmt.Printf("Error piping: %v \n", err)
	}

}

func PipeHandling(input string) int {
	pipeSan := strings.Split(input, "|")
	for i := range pipeSan {
		pipeSan[i] = strings.TrimSpace(pipeSan[i])
	}
	var twodPipe [][]string
	for _, pipeSan := range pipeSan {
		/* we iterate over pipeSan and the body of pipeSan gets appended to twodPipe
		 * and inside that body we split the array with whitespace */
		twodPipe = append(twodPipe, strings.Split(pipeSan, " "))
	}

	twodlen := len(twodPipe)

	switch twodlen {
	case 2:
		twopiper(twodPipe)
	case 3:
		threepiper(twodPipe)
	case 4:
		fourpiper(twodPipe)
	default:
		lawls, _ := exec.Command("sh", "-c", input).Output()
		fmt.Println(lawls)

	}

	return 0
}

func HandlingInput(input string) int {
	/* Sanitized input aka clean */
	SanInput := strings.Split(input, " ")

	if SanInput[0] == "cd" {

		user, err := user.Current()
		LoneCD := strings.Join(SanInput[:], " ")
		LoneCD = strings.Trim(LoneCD, " ")
		if LoneCD == "cd" {
			if err == nil {
				os.Chdir(user.HomeDir)
				return 0
			} else {
				fmt.Printf("Couldnt get user, defaulting to /home\n")
				os.Chdir("/home")
			}
		} else {
			SanInput[1] = strings.Replace(SanInput[1], "~", user.HomeDir, 1)
			os.Chdir(SanInput[1])
			return 0
		}

	} else if SanInput[0] == "exit" {
		return 1

	} else if SanInput[0] == "pwd" {
		pwd, err := os.Getwd()
		if err != nil {
			fmt.Printf("Error getting PWD: %v\n", err)
			return 0
		}
		fmt.Printf("%v\n", pwd)
	} else {
		/*
		 * Comodoro is the command to execute
		 * the reason its called comodoro, is because
		 * the name started as command, then commandtoexec,
		 * then commandeer
		 * then comodoro
		 */
		Comodoro := exec.Command(SanInput[0], SanInput[1:]...)
		Comodoro.Stdin = os.Stdin
		Comodoro.Stdout = os.Stdout
		Comodoro.Stderr = os.Stderr

		err := Comodoro.Start()
		Comodoro.Wait()
		if err != nil {
			fmt.Printf("Error at starting process: %v\n", err)
		}
		return 0
	}
	return 0
}

package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"strings"
)

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
		fmt.Printf("%v@%v > %v >> ", un, hn, wd)

		input := CleanInput(readInput)
		fmt.Printf("input %v\n", input)
		//check :=
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
		fmt.Printf("Contains at least 1 pipe\n")
		return true
	} else {
		fmt.Printf("No pipes here\n")
		return false
	}
	return false
}

func PipeHandling(input string) int {
	pipeSan := strings.Split(input, "|")
	//fmt.Printf("%v\n", pipeSan)
	for i := range pipeSan {
		pipeSan[i] = strings.TrimSpace(pipeSan[i])
	}
	//	pipelen := len(pipeSan)
	//	argshit := string(pipeSan[0])
	var twodPipe [][]string
	for _, pipeSan := range pipeSan {
		/* we iterate over pipeSan and the body of pipeSan gets appended to twodPipe
		 * and inside that body we split the array with whitespace */
		twodPipe = append(twodPipe, strings.Split(pipeSan, " "))
	}

	fmt.Printf("%v\n", twodPipe[0][0])
	twodlen := len(twodPipe)
	i := 0

	for i < twodlen {
		r := 0
		y := 0
		x := len(twodPipe[y][r]) - len(twodPipe[y][r])
		z := x + 1
		l := 0
		c1 := twodPipe[y][x]
		c2 := twodPipe[z][l]
		fmt.Println(c1, c2)

		piping1 := exec.Command(twodPipe[y][r], twodPipe[y][1:]...)
		piping2 := exec.Command(twodPipe[z][l], twodPipe[z][1:]...)

		piping2.Stdin, _ = piping1.StdoutPipe()
		piping2.Stdout = os.Stdout
		piping2.Stderr = os.Stderr

		if i-twodlen == 1 {

			piping2.Start()
			piping1.Run()
			piping2.Wait()
			break
		}
		y++
		i = 12
	}

	/*	for i := range twodPipe {
			fmt.Println(twodPipe[i])
		}

		chucha := exec.Command()
			chucha.Stdin = os.Stdin
			chucha.Stdout = os.Stdout
			chucha.Start()
			chucha.Wait()*/
	return 0
}

func HandlingInput(input string) int {
	/* Sanitized input aka clean */
	SanInput := strings.Split(input, " ")

	fmt.Printf("saninput type: %T\n", SanInput)

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

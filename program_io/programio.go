package programio

import (
	"bufio"
	"fmt"
	"os"
)

func InitStdoutExit(doneChan chan struct{}) {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		text := scanner.Text()
		if text == "exit" {
			fmt.Println("text was entered")
			doneChan <- struct{}{}
			return
		}
		fmt.Println(`enter "exit" to exit program`)
	}
}

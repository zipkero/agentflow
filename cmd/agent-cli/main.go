package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	fmt.Print("입력: ")

	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		fmt.Fprintln(os.Stderr, "입력 읽기 실패")
		os.Exit(1)
	}

	input := strings.TrimSpace(scanner.Text())
	if input == "" {
		fmt.Fprintln(os.Stderr, "입력이 비어있습니다")
		os.Exit(1)
	}

	// TODO(Task 1-5-2): runtime.Run(ctx, input) 호출로 교체
	fmt.Printf("수신된 입력: %q\n", input)
}

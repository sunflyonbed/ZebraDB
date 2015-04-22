package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"redis/resp"
)

var (
	info = flag.String("i", "", "")
)

func main() {
	flag.Parse()

	data := strings.Split(*info, " ")
	var cmd []string
	cmd = append(cmd, "RPUSH")
	cmd = append(cmd, "dbq")
	cmd = append(cmd, string(resp.Format(data)))
	fmt.Fprint(os.Stdout, string(resp.Format(cmd)))
}

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"cbgo"
)

func main() {
	var ipc cbgo.CbIpc

	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s <command>\n", os.Args[0])
		os.Exit(-1)
	}

	if err := ipc.Open(); err != nil {
		fmt.Println("can't open IPC socket", err)
		os.Exit(-1)
	}

	if subcmd, err := json.Marshal(map[string]string{"command": os.Args[1]}); err != nil {
		fmt.Println("can't marshal command '%s': %s\n", os.Args[1], err)
		os.Exit(-1)
	} else {
		cmd := "custom_event " + strings.ReplaceAll(string(subcmd), `"`, `\"`)
		fmt.Printf("cmd='%s'\n", cmd)
		ipc.SendCmd(cmd)
	}
}

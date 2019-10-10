package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type UsersSt struct {
	Browsers []string `json:"browsers"`
	Email    string   `json:"email"`
	Name     string   `json:"name"`
}

func FastSearch(out io.Writer) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	seenBrowsers := make(map[string]struct{})
	foundUsers := ""
	var i int = -1

	var isAndroid, isMSIE bool
	var line []byte
	var user UsersSt

	for scanner.Scan() {
		i++
		line = scanner.Bytes()
		if bytes.Contains(line, []byte("Android")) || bytes.Contains(line, []byte("MSIE")) {

			err := json.Unmarshal(line, &user)
			if err != nil {
				panic(err)
			}

			isAndroid = false
			isMSIE = false

			for _, browserRaw := range user.Browsers {
				if ok := strings.Contains(browserRaw, "Android"); ok {
					isAndroid = true
					seenBrowsers[browserRaw] = struct{}{}
				}
				if ok := strings.Contains(browserRaw, "MSIE"); ok {
					isMSIE = true
					seenBrowsers[browserRaw] = struct{}{}
				}
			}

			if !(isAndroid && isMSIE) {
				continue
			}

			foundUsers += "[" + strconv.Itoa(i) + "] " + user.Name + " <" + strings.Replace(user.Email, "@", " [at] ", 1) + ">\n"
		}
	}

	fmt.Fprintln(out, "found users:\n"+foundUsers)
	fmt.Fprintln(out, "Total unique browsers", len(seenBrowsers))
}

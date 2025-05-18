package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

const (
	MAX_CONNECTIONS = 1000
	MAX_CONNECTION_TIMEOUT = 10 // seconds
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage
	//
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	defer l.Close()
	// we may have multiple connections
	
	// for i := 0; i < MAX_CONNECTIONS; i++ {
	for {
		// TODO: how to restrict max munber of connections?
		// TODO: how to set up timeout for connections?
		// TODO: for long running connections, how to set up?
		go handleConnection(l)
	}
}

func handleConnection(l net.Listener) {
	conn, err := l.Accept()
	// defer conn.Close()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
	buf := make([]byte, 1024)
	// each connection may have multiple requests
	for {
		_, err = conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				continue
			} 
			fmt.Println("Error reading: ", err.Error())
		}
		// parse and execute command 
		if err := parseAndExecuteCommand(conn, buf); err != nil {
			fmt.Println("Error parsing command: ", err.Error())
			continue
		}
		// clean up buf each time 
		clear(buf)
	}
}

func parseAndExecuteCommand(conn net.Conn, buf []byte) error {
	if len(buf) == 0 {
		return nil
	}
	// TODO(LOW): change it could run multiple commands in one statement
	// convert byte array to string
	pos := 0
	var tokens []string
	for pos < len(buf) {
		if buf[pos] == '\x00' {
			break
		}
		switch buf[pos] {
		case '+':
			// simple string
			pos++
			token := parseSimpleString(buf, &pos)
			tokens = append(tokens, token...)
			fmt.Println("[DEBUG] simplestring=[%s]", tokens)
		// case '-':
		// 	// simple error
		// 	if err := parseAndExecuteSimpleError(conn, buf, &pos); err != nil {
		// 		return err
		// 	}
		// case ':':
		// 	// integer
		// 	if err := parseAndExecuteInteger(conn, buf, &pos); err != nil {
		// 		return err
		// 	}
		case '$':
			// bulk strings
			pos++
			tokens = parseBulkString(buf, &pos)
			fmt.Println("[DEBUG] bulkstring=[%s]", tokens)
		case '*':
			// array
			pos++
			// fmt.Println("[DEBUG] receive array. buf=%s pos=%d", buf, pos)
			token :=  parseArray(buf, &pos)
			tokens = append(tokens, token...)
			fmt.Println("[DEBUG] array=[%s]", tokens)
		// case '_':
		// 	// null
		// 	if err := parseAndExecuteNull(conn, buf, &pos); err != nil {
		// 		return err
		// 	}
		// case '#':
		// 	// boolean
		// 	if err := parseAndExecuteBoolean(conn, buf, &pos); err != nil {
		// 		return err
		// 	}
		// case ',':
		// 	// double
		// 	if err := parseAndExecuteDouble(conn, buf, &pos); err != nil {
		// 		return err
		// 	}
		// case '(':
		// 	// big numbers
		// 	if err := parseAndExecuteBigNumbers(conn, buf, &pos); err != nil {
		// 		return err
		// 	}
		// case '!':
		// 	// bulk errors
		// 	if err := parseAndExecuteBulkErrors(conn, buf, &pos); err != nil {
		// 		return err
		// 	}
		// case '=':
		// 	// verbatim strings
		// 	if err := parseAndExecuteVerbatimStrings(conn, buf, &pos); err != nil {
		// 		return err
		// 	}
		// case '%':
		// 	// maps
		// 	if err := parseAndExecuteMaps(conn, buf, &pos); err != nil {
		// 		return err
		// 	}
		// case '|':
		// 	// attributes
		// 	if err := parseAndExecuteAttributes(conn, buf, &pos); err != nil {
		// 		return err
		// 	}
		// case '~':
		// 	// sets
		// 	if err := parseAndExecuteSets(conn, buf, &pos); err != nil {
		// 		return err
		// 	}
		// case '>':
		// 	// pushes
		// 	if err := parseAndExecutePushes(conn, buf, &pos); err != nil {
		// 		return err
		// 	}

		default:
			// unknown command
			fmt.Println("[DEBUG] current pos=%d", pos)
			return fmt.Errorf("unknown command: %v", buf[pos])
		}
	}

	fmt.Println("[DEBUG] token=[%s]", tokens)
	execute(conn, tokens)
	return nil
}

func execute(conn net.Conn, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("empty request")
	}
	cmd := args[0]
	switch strings.ToLower(cmd) {
	case "ping":
		conn.Write([]byte("+PONG\r\n")) 
	case "echo":
		if len(args) < 2 {
			return nil
		}
		conn.Write([]byte(buildBulkString(args[1])))
	case "set":
	    var err error
		if len(args) < 3 {
			return nil
		}
		k, v := args[1], args[2]
		ttl := int64(-1)
		if len(args) > 3 {
			opt_k, opt_v := args[3], args[4]
			if strings.ToLower(opt_k) == "px" {
				t, _ := strconv.Atoi(opt_v)
				ttl = int64(t)
			}
		}
		err = executeCommandSetWithExpiry(k, v, int64(ttl)); 
		if err != nil {
			return fmt.Errorf("fail to set value. err=%v", err)
		}
		conn.Write([]byte(buildSimpleString("OK")))
	case "get":
		if len(args) < 2 {
			return nil
		}
		k := args[1]
		v, err := executeCommandGet(k)
		if err != nil {
			return fmt.Errorf("fail to get value. err=%v", err)
		}
		fmt.Println("[DEBUG] get value=%s", v)
		conn.Write([]byte(buildBulkString(v)))
	default:
		return fmt.Errorf("unknown command")
	}
	return nil
}





	






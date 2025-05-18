package main

import (
	"net"
	"strconv"
	"fmt"
)

// assume receive byte string e.g.: +PING\r\n
func parseSimpleString(buf []byte, pos *int) []string {
	tokens := []string{}
	s := *pos
	for *pos < len(buf) - 1 {
		if buf[*pos] == '\r' && buf[*pos + 1] == '\n' {
			break
		}
		*pos++
	}
	
	tokens = append(tokens, string(buf[s:*pos]))
	return tokens
}

func parseAndExecuteSimpleError(conn net.Conn, buf []byte, pos *int) error {
	// Parse the simple error from the buffer and send it to the connection
	// Implementation goes here
	return nil
}

func parseAndExecuteInteger(conn net.Conn, buf []byte, pos *int) error {

	// Parse the integer from the buffer and send it to the connection
	// Implementation goes here
	return nil
}

func parseBulkString(buf []byte, pos *int) []string {
	// Parse the bulk string from the buffer and send it to the connection
	// Implementation goes here
	tokens := []string{}
	l := getLength(buf, pos)
	s := *pos
	// fmt.Println("[DEBUG] parseBulkString start=%d", *pos)
	// fmt.Println("[DEBUG] slice=[%x]", buf[s:s + l])
	tokens = append(tokens, string(buf[s:s + l]))
	*pos += (l + 2)
	// fmt.Println("[DEBUG] parseBulkString end=%d", *pos)
	// fmt.Println("[DEBUG] bulk string len=[%d] val=[%s]", l, tokens)
	return tokens
}

// *2\r\n$4\r\nECHO\r\n$4\r\npear\r\n
        
func parseArray(buf []byte, pos *int) []string {
	// Parse the array from the buffer and send it to the connection
	// Implementation goes here
	tokens := []string{}
	// get length of the array
	l := getLength(buf, pos)
	if l < 1 {
		return tokens
	}

	// get token of each element in the array
	for i := 0; i < l; i++ {
		switch buf[*pos] {
		case '$':
			*pos++
			tmp := parseBulkString(buf, pos)
			// fmt.Println("[DEBUG] get bulkstring=%s", tmp)
			tokens = append(tokens, tmp...)
		}
	}
	// debug 
	// fmt.Println("[DEBUG] array=%s", tokens)
	return tokens
}

func parseAndExecuteNull(conn net.Conn, buf []byte, pos int) error {
	// Parse the null from the buffer and send it to the connection
	// Implementation goes here
	return nil
}

func parseAndExecuteDouble(conn net.Conn, buf []byte, pos int) error {
	// Parse the double from the buffer and send it to the connection
	// Implementation goes here
	return nil
}

func parseAndExecuteBigNumbers(conn net.Conn, buf []byte, pos int) error {
	// Parse the big numbers from the buffer and send it to the connection
	// Implementation goes here
	return nil
}

func parseAndExecuteBulkErrors(conn net.Conn, buf []byte, pos int) error {
	// Parse the bulk errors from the buffer and send it to the connection
	// Implementation goes here
	return nil
}

func parseAndExecuteBoolean(conn net.Conn, buf []byte, pos int) error {
	// Parse the boolean from the buffer and send it to the connection
	// Implementation goes here
	return nil
}

func parseAndExecuteVerbatimStrings(conn net.Conn, buf []byte, pos int) error {
	return nil
}

func parseAndExecuteMaps(conn net.Conn, buf []byte, pos int) error {
	// Parse the maps from the buffer and send it to the connection
	// Implementation goes here
	return nil
}

func parseAndExecuteAttributes(conn net.Conn, buf []byte, pos int) error {
	// Parse the attributes from the buffer and send it to the connection
	// Implementation goes here
	return nil
}

func parseAndExecuteSets(conn net.Conn, buf []byte, pos int) error {
	// Parse the sets from the buffer and send it to the connection
	// Implementation goes here
	return nil
}

func parseAndExecutePushes(conn net.Conn, buf []byte, pos int) error {
	// Parse the pushes from the buffer and send it to the connection
	// Implementation goes here
	return nil
}


// HELPER FUNCTION
func getLength(buf []byte, pos *int) int {
	s := *pos
	// fmt.Println("[DEBUG] getLength start=%d", *pos)
	for *pos < len(buf) - 1 {
		if buf[*pos] == '\r' && buf[*pos + 1] == '\n' {
			break
		}
		*pos++
	}
	l, _ := strconv.Atoi(string(buf[s:*pos]))
	*pos += 2
	// fmt.Println("[DEBUG] getLength end=%d", *pos)

	return l
}

func buildBulkString(s string) string {
	if (s == "-1") {
		return "$-1\r\n"
	}
	return fmt.Sprintf("$%d\r\n%s\r\n", len(s), s)
}

func buildSimpleString(s string) string {
	return fmt.Sprintf("+%s\r\n", s)
}

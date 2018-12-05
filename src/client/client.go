//server
package main

import (
        "bufio"
        "flag"
        "fmt"
        "net"
        "os"
)

func main() {
        var portFlag string
        flag.StringVar(&portFlag, "port", "9090", "Server port. Defalt: 9090")
        flag.StringVar(&portFlag, "p", "9090", "Server port. Defalt: 9090")
        var hostFlag string
        flag.StringVar(&hostFlag, "host", "127.0.0.1", "Server host. Defalt: 127.0.0.1")
        flag.StringVar(&hostFlag, "h", "127.0.0.1", "Server host. Defalt: 127.0.0.1")
        flag.Parse()
        fmt.Println(portFlag)
        conn, err := net.Dial("tcp", hostFlag+":"+portFlag)
        if err != nil {
                fmt.Println(err)
                return
        }
        defer conn.Close()
        serverWriter := bufio.NewWriter(conn)
        serverReader := bufio.NewReader(conn)
        comandReader := bufio.NewReader(os.Stdin)
        for {
                fmt.Print("> ")
                command, err := comandReader.ReadString('\n')
                if err != nil {
                        fmt.Println(err)
                }
                if command != "" {
                        n, err := serverWriter.WriteString(command)
                        if n == 0 || err != nil {
                                fmt.Println(err)
                                return
                        }
                        serverWriter.Flush()
                }
                reply, err := serverReader.ReadString('\n')
                if err != nil {
                        fmt.Println(err)
                }
                fmt.Print(reply)
        }
}
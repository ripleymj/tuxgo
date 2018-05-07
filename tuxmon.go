package main

import "bufio"
import "fmt"
import "io"
import "os"
import "os/exec"
import "strings"
import "time"

func main() {

    env := os.Environ()
    tuxConfig := "/home/psadm2/psft/pt/8.56/appserv/APPDOM/PSTUXCFG"
    env = append(env, fmt.Sprintf("TUXCONFIG=%s", tuxConfig))

    for {
        tuxCmd := exec.Command("tmadmin", "-r");

        tuxCmd.Env = env

        tuxIn, _ := tuxCmd.StdinPipe()
        tuxOut, _ := tuxCmd.StdoutPipe()
        oBuf := bufio.NewReader(tuxOut)
    
        tuxCmd.Start()

        tuxIn.Write([]byte("verbose on\n"))
        tuxIn.Write([]byte("page off\n"))

        tuxIn.Write([]byte("psr\n"))
        tuxIn.Write([]byte("pq\n"))
        tuxIn.Write([]byte("pclt\n"))
        tuxIn.Write([]byte("quit\n"))

        moreMessages := true

        MessageRead:
        for moreMessages {
            message,_ := oBuf.ReadString('\n')
            message = strings.TrimLeft(message, " >")
    
            if (strings.Index(message, "Group ID:") == 0) {
                fmt.Println("psr message")
            } else if (strings.Index(message, "Prog Name:") == 0) {
                fmt.Println("pq message")
            } else if (strings.Index(message, "LMID:") == 0) {
                fmt.Println("pclt message")
            }
    
            appendLine := true
            for appendLine {
                line, err := oBuf.ReadString('\n')
                line = strings.TrimLeft(line," >")
                if(err != nil ) {
                    if(err != io.EOF) {
                        fmt.Printf("Error: %s\n", err)
                    }
                    moreMessages = false
                    break MessageRead
                }
    
                if(line == "\n") {
                    appendLine = false
                } else {
                   message += line 
                }
            }
            fmt.Print(message)
            fmt.Println("----")
        }
        time.Sleep(1 * time.Second)
        tuxIn.Close()
        tuxCmd.Wait()
    }
} 

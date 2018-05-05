package main

import "fmt"
import "io"
import "bufio"
import "os"
import "os/exec"
import "strings"
import "time"

func main() {

    tuxCmd := exec.Command("tmadmin", "-r");

    env := os.Environ()
    tuxConfig := "/home/psadm2/psft/pt/8.56/appserv/APPDOM/PSTUXCFG"
    env = append(env, fmt.Sprintf("TUXCONFIG=%s", tuxConfig))
    tuxCmd.Env = env
    tuxIn, _ := tuxCmd.StdinPipe()
    tuxOut, _ := tuxCmd.StdoutPipe()
    oBuf := bufio.NewReader(tuxOut)
    tuxCmd.Start()
    
    go tuxWrite(tuxIn)
    doTuxLoop := true

    for doTuxLoop {
        appendMessage := true

        message,_ := oBuf.ReadString('\n')
        message = strings.TrimLeft(message, " >")

        if (strings.Index(message, "Group ID:") == 0) {
            fmt.Println("psr message")
        } else if (strings.Index(message, "Prog Name:") == 0) {
            fmt.Println("pq message")
        } else if (strings.Index(message, "LMID:") == 0) {
            fmt.Println("pclt message")
        }

        for appendMessage {
            line, err := oBuf.ReadString('\n')
            line = strings.TrimLeft(line," >")
            if(err != nil) {
                fmt.Println(err)
                break
            }

            if(line == "\n") {
                appendMessage = false
            } else {
                //fmt.Print(line)
               message += line 
            }
        }
        fmt.Print(message)
        fmt.Println("----")
        //fmt.Println("Done with loop")
        //doTuxLoop = false
    }

    tuxCmd.Wait()
} 

func tuxWrite(pipe io.WriteCloser) {
    pipe.Write([]byte("verbose on\n"))
    pipe.Write([]byte("page off\n"))

    for i := 0; i < 3; i++ {
        pipe.Write([]byte("psr\n"))
        pipe.Write([]byte("pq\n"))
        pipe.Write([]byte("pclt\n"))
        time.Sleep(5 * time.Second)
        //fmt.Println(s)
    }

    pipe.Write([]byte("quit\n"))
    pipe.Close()
}

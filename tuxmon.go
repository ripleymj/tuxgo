package main

import "fmt"
import "io"
import "bufio"
import "os"
import "os/exec"
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
    
    doTuxLoop := true

    for doTuxLoop {
        doCmdLoop := true

        go tuxWrite(tuxIn)

        for doCmdLoop {
            line, err := oBuf.ReadString('\n')

            if(err != nil) {
                fmt.Println(err)
                break
            }

            if(line == "\n") {
                fmt.Println("----")
            } else {
                fmt.Print(line)
            }
        }

        fmt.Println("Done with loop")
        doTuxLoop = false
    }

    tuxCmd.Wait()
} 

func tuxWrite(pipe io.WriteCloser) {
    pipe.Write([]byte("verbose on\n"))

    for i := 0; i < 3; i++ {
        pipe.Write([]byte("psr\n"))
        pipe.Write([]byte("pq\n"))
        pipe.Write([]byte("pclt\n"))
        time.Sleep(5 * time.Second)
        //fmt.Println(s)
    }

    pipe.Close()
}

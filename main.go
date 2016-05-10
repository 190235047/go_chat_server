package main
 
import (
    "flag"
    "fmt"
    "log"
    "os"
    "runtime"
    "mynet"
)
 
var (
    Port           = flag.String("i", ":12345", "IP port to listen on")
    logFileName    = flag.String("log", "cServer.log", "Log file name")
    //configFileName = flag.String("configfile", "config.ini", "General configuration file")
)
/*
var (
    configFile = flag.String("configfile", "config.ini", "General configuration file")
)
*/ 
func main() {
    runtime.GOMAXPROCS(runtime.NumCPU())
    flag.Parse()
 
    //set logfile Stdout
    logFile, logErr := os.OpenFile(*logFileName, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
    if logErr != nil {
        fmt.Println("Fail to find", *logFile, "cServer start Failed")
        os.Exit(1)
    }
    log.SetOutput(logFile)
    log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
    //set logfile Stdout End
 
    //start listen
    listenErr := mynet.StartListen(*Port)
    if listenErr != nil {
        log.Fatalf("Server abort! Cause:%v \n", listenErr)
    }
}

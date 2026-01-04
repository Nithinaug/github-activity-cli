package main
import (
	"fmt"
	"os"
)
func main(){
	if len(os.Args) < 2{
		fmt.Println("Github-Activity:<username> ")
		os.Exit(1)
	}
	usname := os.Args[1]
	proc,err := Process(usname)
	if err!=nil{
          fmt.Print("Error : ",err)
		  os.Exit(1)
	}
		printSummary(proc)
}
package main

import (
  //"os"
  "flag"
  "io/ioutil"
  "os"
  "bufio"
  "strings" 
  "fmt"
)

var (

  outFile  = flag.String("outFile", "", "File to output new ducky script")
  inFile = flag.String("inFile", "", "File to turn into a ducky script")
  verbose  = flag.Bool("verbose", false, "verbosly send messages to the console")
  OS  = flag.String("OS", "windows", "the operating system the script is supposed to run on: windows, mac, linux")

  duckyHeader = [4]string{
    "REM Autogenerated Script via Ducky_Maker",
    "REM Ducky_Maker Author: github.com/ahhh/Ducky_Maker",
    "REM Script Author: X",
    "DELAY 500",
  }

  winHeader = [6]string{
    "GUI r",
    "DELAY 500",
    "STRING cmd",
    "DELAY 500",
    "ENTER",
    "DELAY 500",
  }

  linHeader = [4]string{
    "CTRL-ALT T",
    "DELAY 500",
    "ENTER",
    "DELAY 500",
  }
  macHeader = [6]string{
    "GUI SPACE",
    "DELAY 500",
    "STRING terminal",
    "DELAY 500",   
    "ENTER",
    "DELAY 500",
  }
  
  duckyFooter = [2]string{
    "DELAY 300",
    "ENTER",
  }
)

func paramCheck() bool {
  canRun := true
  // Make sure outFile
  if (*outFile != "") {
    if *verbose == true {
      fmt.Println("outFile has values set to " + string(*outFile))
    }
  } else {
    fmt.Println("No outFile provided!")
    canRun = false
  }
  // Make sure inFile is set
  if (*inFile != "") {
    if *verbose == true {
      fmt.Println("inFile has values set to " + string(*inFile))
    }
  } else {
    fmt.Println("No inFile provided!")
    canRun = false
  }
  // Make sure OS is set
  if ((*OS != "windows") || (*OS != "linux") || (*OS != "mac")) {
    if *verbose == true {
      fmt.Println("OS has values set to " + string(*OS))
    }
  } else {
    fmt.Println("No valid OS provided!")
    canRun = false
  }

  if !canRun {
    fmt.Println("Missing mandatory paramaters. use -h for the help menu.")
    return false
  } else {
    return true
  }

}


func main() {

  flag.Parse()
  inFileFlag := flag.Lookup("inFile")
  outFileFlag := flag.Lookup("outFile")
  OSFlag := flag.Lookup("OS")

  shouldRun := paramCheck()
  if !shouldRun {
    if *verbose == true {
      fmt.Println("Failed mandatory flag checks, plz set required flags!!")
    }
  }

  fileData := ReadFile(inFileFlag.Value.String())
  new_script := Format_lines(fileData, OSFlag.Value.String())
  new_file := strings.Join(new_script[:],"\n")
  err := WriteFile(new_file, outFileFlag.Value.String())
  if err != nil {
    fmt.Println(err.Error())
  }
}


func Format_lines(work []string, operatingSys string) []string{
  //fmt.Printf("dat work: %s\n", work)
  var newHotness []string
  // Copy over Ducky header
  for _, line := range duckyHeader {
    newHotness = append(newHotness, string(line))
  }
  // Copy over OS specific prompt opening
  if operatingSys == "windows" {
    for _, line := range winHeader {
      newHotness = append(newHotness, string(line))
    }
  } else if (operatingSys == "linux"){
    for _, line := range linHeader {
      newHotness = append(newHotness, string(line))
    }
  } else if (operatingSys == "mac") {
    for _, line := range macHeader {
      newHotness = append(newHotness, string(line))
    }
  }
  // Copy over script content
  for _, line := range work {
    //fmt.Printf("sentance %d: %s\n", i, string(line))
    newHotness = append(newHotness,"STRING " + string(line))
    newHotness = append(newHotness, "ENTER")
  }
  // Copy over Ducky footer
  for _, line := range duckyFooter {
    newHotness = append(newHotness, string(line))
  }
  return newHotness

}

// readLines reads a whole file into memory
// and returns a slice of its lines.
func ReadFile(path string) []string {
  file, err := os.Open(path)
  if err != nil {
    fmt.Println(err.Error())
    return nil
  }
  defer file.Close()

  var lines []string
  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    lines = append(lines, scanner.Text())
  }
  return lines
}


func WriteFile(text, outFile string) (error) {
  //Check if log file exists to create or append
  if Exists(outFile) {
    //Write the lines to the file
    f, err := os.OpenFile(outFile,
       os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
      fmt.Println(err.Error())
      return err
    }
    defer f.Close()
    if _, err := f.WriteString(text); err != nil {
      fmt.Println(err.Error())
      return err
    }  
  } else {
    //Create file and write first lines
    err := ioutil.WriteFile(outFile, []byte(text), 0700)
    if err != nil {
  	  fmt.Println(err.Error())
      return err
    }
  }
  return nil
}


func Exists(path string) bool {
  // Run stat on a file
  _, err := os.Stat(path)
  // If it runs fine the file exists
  if err == nil {
    return true
  }
  // If stat fails then the file does not exist
  if *verbose == true { fmt.Println("Creating new file as " + err.Error()) }
  return false
}


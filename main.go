package goPSRemoting

import (
	"bytes"
	"errors"
	"log"
	"os/exec"
	"runtime"
	"strings"
)

func runCommand(args ...string) (string, error) {
	cmd := exec.Command(args[0], args[1:]...)
	var out bytes.Buffer
	var err bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &err
	cmd.Run()
	// convert err to an error type if there is an error returned
	var e error
	if err.String() != "" {
		e = errors.New(err.String())
	}
	return strings.TrimRight(out.String(), "\r\n"), e
}
func RunPowershellCommand(username string, password string, server string, command string, usessl string, usessh string, authentication string) (string, error) {
	var pscommand string
	if runtime.GOOS == "windows" {
		log.Printf("Windows host detected.")
		pscommand = "powershell.exe -ExecutionPolicy ByPass"
	} else {
		pscommand = "pwsh"
	}
	var winRMPre string
	if usessh == "1" {
		winRMPre = "$s = New-PSSession -HostName " + server + " -Username " + username + " -SSHTransport"
	} else {
		if authentication == "negociate" {
			log.Printf("Authentication detected: Negociate")
			winRMPre = "$SecurePassword = '" + password + "' | ConvertTo-SecureString -AsPlainText -Force; $cred = New-Object System.Management.Automation.PSCredential -ArgumentList '" + username + "', $SecurePassword; $s = New-PSSession -Authentication Negotiate -ComputerName " + server + " -Credential $cred"
			log.Printf("WinRMPre command: %v",winRMPre)
		}
	}
	var winRMPost string
	if runtime.GOOS == "windows" {
		winRMPost = "; Invoke-Command -Session $s -Scriptblock { " + command + " }; Remove-PSSession $s"
	} else {
		winRMPost = "; Invoke-Command -Session $s -Scriptblock { powershell '" + command + "' }; Remove-PSSession $s"
	}
	
	log.Printf("WinRMPost command: %v",winRMPost)
	var winRMCommand string
	if usessl == "1" {
		winRMCommand = winRMPre + " -UseSSL" + winRMPost
	} else {
		winRMCommand = winRMPre + winRMPost
	}
	
	log.Printf("WinRMCommand: %v",winRMCommand)
	
	out, err := runCommand(pscommand, "-command", winRMCommand)
	return out, err
}

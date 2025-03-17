package cli

import (
	"os"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/zorg113/golang_dipl/atibruteforce/internal/controller/httpapi/handlers"
	"github.com/zorg113/golang_dipl/atibruteforce/model/entity"
	"github.com/zorg113/golang_dipl/atibruteforce/model/service"
)

var suggestions = []prompt.Suggest{
	{Text: "blacklist add [ip_address] [mask]", Description: "Add a blacklisted IP address"},
	{Text: "blacklist remove [ip_address] [mask]", Description: "Remove a blacklisted IP address"},
	{Text: "blacklist get", Description: "Get all blacklisted IP addresses"},
	{Text: "whitelist add [ip_address] [mask]", Description: "Add a whitelisted IP address"},
	{Text: "whitelist remove [ip_address] [mask]", Description: "Remove a whitelisted IP address"},
	{Text: "whitelist get", Description: "Get all whitelisted IP addresses"},
	{Text: "bucket remove [login] [ip_address]", Description: "Remove login and ip adress from bucket"},
	{Text: "help", Description: "show list commands"},
	{Text: "exit", Description: "Exit the program"},
}

type CommandLineInterface struct {
	serviceAuth      *service.Authorization
	serviceBlackList *service.BlackList
	serviceWhiteList *service.WhiteList
}

func NewCommandLineInterface(serviceAuth *service.Authorization, serviceBlackList *service.BlackList, serviceWhiteList *service.WhiteList) *CommandLineInterface {
	return &CommandLineInterface{
		serviceAuth:      serviceAuth,
		serviceBlackList: serviceBlackList,
		serviceWhiteList: serviceWhiteList,
	}
}

func (cli *CommandLineInterface) Run(ch chan os.Signal) {
	executer := prompt.Executor(func(s string) {
		s = strings.TrimSpace(s)
		setCommand := strings.Split(s, " ")
		switch setCommand[0] {
		case "blacklist":
			cli.blackListCommand(setCommand)
		case "whitelist":
			cli.whiteListCommand(setCommand)
		case "bucket":
			cli.bucketCommand(setCommand)
		case "help":
			println("Available commands:")
			for _, suggestion := range suggestions {
				println(suggestion.Text + " - " + suggestion.Description)
			}
		case "exit":
			ch <- os.Interrupt
			return
		default:
			println("Invalid command")
		}
	})
	completer := prompt.Completer(func(in prompt.Document) []prompt.Suggest {
		w := in.GetWordBeforeCursor()
		if w == "" {
			return []prompt.Suggest{}
		}
		return prompt.FilterHasPrefix(suggestions, w, true)
	})
	defer func() {
		if a := recover(); a != nil {
			println("Command line interface not available")
		}
	}()
	prompt.New(executer, completer).Run()
}
func checkLenCommand(command []string, msg string) bool {
	ret := len(command) != 4
	if !ret {
		println("Usage: " + msg)
	}
	return ret
}

func (cli *CommandLineInterface) blackListCommand(command []string) {
	switch command[1] {
	case "add":
		if !checkLenCommand(command, "blacklist add [ip_address] [mask]") {
			return
		}
		cli.addIpToBlackList(
			entity.IpNetwork{
				Ip:   command[2],
				Mask: command[3],
			})
	case "remove":
		if !checkLenCommand(command, "blacklist remove [ip_address] [mask]") {
			return
		}
		cli.deleteIpFromBlackList(entity.IpNetwork{
			Ip:   command[2],
			Mask: command[3],
		})
	case "get":
		cli.getIpsFromBlackList()
	default:
		println("Invalid command")
	}
}

func (cli *CommandLineInterface) whiteListCommand(command []string) {
	switch command[1] {
	case "add":
		if !checkLenCommand(command, "whitelist add [ip_address] [mask]") {
			return
		}
		cli.addIpToWhiteList(entity.IpNetwork{
			Ip:   command[2],
			Mask: command[3],
		})
	case "remove":
		if !checkLenCommand(command, "whitelist remove [ip_address] [mask]") {
			return
		}
		cli.deleteIpFromWhiteList(entity.IpNetwork{
			Ip:   command[2],
			Mask: command[3],
		})
	case "get":
		cli.getIpsFromWhiteList()
	default:
		println("Invalid command")
	}
}

func (cli *CommandLineInterface) bucketCommand(command []string) {
	switch command[1] {
	case "reset":
		if !checkLenCommand(command, "bucket reset [login] [ip_address]") {
			return
		}
		cli.resetBucket(entity.Request{
			Login: command[2],
			Ip:    command[3],
		})
	default:
		println("Invalid command")
	}
}

func (cli *CommandLineInterface) addIpToBlackList(ipNet entity.IpNetwork) {
	isValidated := handlers.ValidateIP(ipNet)
	if !isValidated {
		println("Invalid IP or mask")
		return
	}
	err := cli.serviceBlackList.AddIP(ipNet)
	if err != nil {
		println(err.Error())
		return
	}
	println("IP address added to black list")
}

func (cli *CommandLineInterface) deleteIpFromBlackList(ipNet entity.IpNetwork) {
	isValidated := handlers.ValidateIP(ipNet)
	if !isValidated {
		println("Invalid IP or mask")
		return
	}
	err := cli.serviceBlackList.DeleteIP(ipNet)
	if err != nil {
		println(err.Error())
		return
	}
	println("IP address removed from black list")
}

func (cli *CommandLineInterface) getIpsFromBlackList() {
	ips, err := cli.serviceBlackList.GetIPs()
	if err != nil {
		println(err.Error())
		return
	}
	for _, ip := range ips {
		println("ip:%s mask:%s", ip.Ip, ip.Mask)
	}
}

func (cli *CommandLineInterface) addIpToWhiteList(ipNet entity.IpNetwork) {
	isValidated := handlers.ValidateIP(ipNet)
	if !isValidated {
		println("Invalid IP or mask")
		return
	}
	err := cli.serviceWhiteList.AddIP(ipNet)
	if err != nil {
		println(err.Error())
		return
	}
	println("IP address added to whitelist")
}

func (cli *CommandLineInterface) deleteIpFromWhiteList(ipNet entity.IpNetwork) {
	isValidated := handlers.ValidateIP(ipNet)
	if !isValidated {
		println("Invalid IP or mask")
		return
	}
	err := cli.serviceWhiteList.DeleteIP(ipNet)
	if err != nil {
		println(err.Error())
		return
	}
	println("IP address removed from whitelist")
}

func (cli *CommandLineInterface) getIpsFromWhiteList() {
	ips, err := cli.serviceWhiteList.GetIPs()
	if err != nil {
		println(err.Error())
		return
	}
	for _, ip := range ips {
		println("ip:%s mask:%s", ip.Ip, ip.Mask)
	}
}

func (cli *CommandLineInterface) resetBucket(req entity.Request) {
	isValidated := handlers.ValidateRequest(req)
	if !isValidated {
		println("Invalid request")
		return
	}
	isReset := cli.serviceAuth.ResetIpBucket(req.Ip)
	if !isReset {
		println("ip address not find")
		return
	}
	println("Bucket reset for IP address: " + req.Ip)
	isReset = cli.serviceAuth.ResetLoginInBucket(req.Login)
	if !isReset {
		println("login not find")
		return
	}
	println("Bucket reset for login: " + req.Login)
}

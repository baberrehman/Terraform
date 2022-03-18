package terraform

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/printer"
	jsonParser "github.com/hashicorp/hcl/json/parser"
	"golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

// Helper functions for Terraform integration

var basePath string // = constants.TempDir + "chef" //"/home/baber"

func run_terraform_script(projectID string, chefServerUrl string, IP string, userName string, environmentName ...string) bool {

	cmd := exec.Command("knife", "bootstrap", IP, "--ssh-user", userName)

	println("Executing command: ", "knife", "bootstrap", IP, "--ssh-user", userName)

	stdout, err3 := cmd.StdoutPipe()
	if err3 != nil {
		println(err3)
		return false
	}
	stderr, err4 := cmd.StderrPipe()
	if err4 != nil {
		println(err4)
		return false
	}

	//cmd.Stdout = os.Stdout
	//cmd.Stderr = os.Stderr
	multi := io.MultiReader(stdout, stderr)
	inp := bufio.NewScanner(multi)
	err := cmd.Start()
	if err != nil {
		log.Println("Error while creating client\n", err)
		return false
	}
	fmt.Println("** Generating logs **")
	for inp.Scan() {

		line := inp.Text()
		println(line, projectID)
	}
	exitStatus := cmd.Wait()
	if exiterr, ok := exitStatus.(*exec.ExitError); ok {
		println("**** Exit Status is:", exiterr, "****")
		return false
	}

	println("Terraform Orchestration Successful!")
	return true
}

func terraformInit() {
	cmd := exec.Command("terraform", "init")
	runCMDLocal(cmd)
}

func terraformFormat() {
	cmd := exec.Command("terraform", "fmt")
	runCMDLocal(cmd)
}

func terraformValidate() {
	cmd := exec.Command("terraform", "validate")
	runCMDLocal(cmd)
}

func terraformApply() {
	cmd := exec.Command("terraform", "apply")
	runCMDLocal(cmd)
}

func isFolderExist(path string) bool {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			println("File does not exist at: " + path)
			println(err.Error())
			return false
		}
	}
	println("File exist at: " + path)
	return true
}

func runCMDLocal(cmd *exec.Cmd) bool {
	stdout, err1 := cmd.StdoutPipe()
	if err1 != nil {
		println(err1.Error())
		return false
	}
	stderr, err2 := cmd.StderrPipe()
	if err2 != nil {
		println(err2.Error())
		return false
	}

	multi := io.MultiReader(stdout, stderr)
	inp := bufio.NewScanner(multi)
	err3 := cmd.Start()
	if err3 != nil {
		log.Println("Error starting command\n", err3)

		return false
	}
	fmt.Println("**** Generating logs ****")
	for inp.Scan() {

		line := inp.Text()
		println(line)
	}
	exitStatus := cmd.Wait()
	if exiterr, ok := exitStatus.(*exec.ExitError); ok {
		println("**** Exit Status is:", exiterr, "****")
		return false
	}
	return true
}

func chefServerUrlToDirectoryName(chefServerUrl string) string {
	chefServerUrl = strings.Replace(chefServerUrl, ":", "", -1)
	chefServerUrl = strings.Replace(chefServerUrl, "/", "", -1)
	chefServerUrl = strings.Replace(chefServerUrl, ".", "_", -1)
	chefServerUrl = strings.ToLower(chefServerUrl)
	return chefServerUrl
}

func runCMDRemote(cmd, hostname string, config *ssh.ClientConfig) (string, error) {
	conn, err := ssh.Dial("tcp", hostname+":22", config)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	session, err := conn.NewSession()
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer session.Close()
	var stdoutBuf bytes.Buffer
	session.Stdout = &stdoutBuf
	session.Stderr = &stdoutBuf
	session.Run(cmd)
	return stdoutBuf.String(), nil
}

/*
hcl_to_json converts HCL to JSON
*/
func hcl_to_json(script string) string {
	cmd := exec.Command("json2hcl", "-reverse", script)
	runCMDLocal(cmd)
	return ""
}

/*
ToJSON converts terraform script to JSON
file arguments contains the path to terraform script
*/

func ToJSON(file string) (error, string) {
	// TODO: This needs to change to read from another source
	// file arguments contains the path to terraform script
	input, err := ioutil.ReadFile(file)
	if err != nil {
		return fmt.Errorf("unable to read from stdin: %s", err), ""
	}

	var v interface{}
	err = hcl.Unmarshal(input, &v)
	if err != nil {
		return fmt.Errorf("unable to parse HCL: %s", err), ""
	}

	json, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("unable to marshal json: %s", err), ""
	}

	output := string(json)

	//fmt.Println(output)
	return nil, output
}

func ToHCL(input string) (error, string) {
	// TODO: This needs to change to read from another source
	//input, err := ioutil.ReadAll(os.Stdin)
	//if err != nil {
	//	return fmt.Errorf("unable to read from stdin: %s", err)
	//}

	ast, err := jsonParser.Parse([]byte(input))
	if err != nil {
		return fmt.Errorf("unable to parse JSON: %s", err), ""
	}

	err = printer.Fprint(os.Stdout, ast)
	if err != nil {
		return fmt.Errorf("unable to print HCL: %s", err), ""
	}

	return nil, ""
}

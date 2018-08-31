package cmd

import (
	"bufio"
	"io"
	"os"
	"os/exec"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/spf13/cobra"
	"github.com/youyo/credentor"
)

func init() {}

func NewCmdRoot() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "credentor",
		Short: "",
		Run: func(cmd *cobra.Command, args []string) {
			credentials, err := getCredentials()
			if err != nil {
				cmd.SetOutput(os.Stderr)
				cmd.Println(err)
				os.Exit(1)
			}

			if len(args) != 0 {
				os.Setenv("AWS_ACCESS_KEY_ID", credentials.AccessKeyID)
				os.Setenv("AWS_SECRET_ACCESS_KEY", credentials.SecretAccessKey)
				os.Setenv("AWS_SESSION_TOKEN", credentials.SessionToken)
				executeCommand(cmd, args)
			} else {
				creds := format(credentials)
				for _, v := range creds {
					cmd.Println(v)
				}
			}

		},
	}

	cobra.OnInitialize(initConfig)
	return cmd
}

func getCredentials() (*credentials.Value, error) {
	c := credentor.NewConfig()
	if err := c.ExtractRoleInfo(); err != nil {
		return nil, err
	}

	v, err := credentor.GetCredentials(c.GetRoleArn(), c.ExportSessionOptions())
	return v, err
}

func format(creds *credentials.Value) []string {
	return []string{
		"export AWS_ACCESS_KEY_ID=" + creds.AccessKeyID,
		"export AWS_SECRET_ACCESS_KEY=" + creds.SecretAccessKey,
		"export AWS_SESSION_TOKEN=" + creds.SessionToken,
	}
}

func executeCommand(cmd *cobra.Command, args []string) {
	var command *exec.Cmd
	if len(args) == 1 {
		command = exec.Command(args[0])
	} else if len(args) >= 2 {
		command = exec.Command(args[0], args[1:]...)
	}

	stdout, err := command.StdoutPipe()
	if err != nil {
		cmd.SetOutput(os.Stderr)
		cmd.Println(err)
		os.Exit(1)
	}

	stderr, err := command.StderrPipe()
	if err != nil {
		cmd.SetOutput(os.Stderr)
		cmd.Println(err)
		os.Exit(1)
	}

	go writer(cmd, stdout)
	go writer(cmd, stderr)

	command.Start()
	command.Wait()
}

func writer(cmd *cobra.Command, input io.ReadCloser) {
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		cmd.Println(scanner.Text())
	}
}

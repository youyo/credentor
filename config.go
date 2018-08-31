package credentor

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	ini "gopkg.in/ini.v1"
)

type Config struct {
	EnvironmentVariables EnvironmentVariable
	FilePaths            FilePath
	Role                 RoleInfo
}

type EnvironmentVariable struct {
	AwsProfile    string
	AwsConfigFile string
}

type FilePath struct {
	AwsConfigFile string
}

type RoleInfo struct {
	RoleArn       string
	SourceProfile string
	MfaSerial     string
	ExternalID    string
}

func NewConfig() *Config {
	c := new(Config)
	c.readEnvironmentVariables()
	c.setFilePaths()
	return c
}

func (c *Config) readEnvironmentVariables() {
	e := &c.EnvironmentVariables
	e.AwsProfile = os.Getenv("AWS_PROFILE")
	e.AwsConfigFile = os.Getenv("AWS_CONFIG_FILE")
}

func (c *Config) setFilePaths() {
	c.FilePaths.AwsConfigFile = filepath.Join(os.Getenv("HOME"), ".aws/config")
	if c.EnvironmentVariables.AwsConfigFile != "" {
		c.FilePaths.AwsConfigFile = c.EnvironmentVariables.AwsConfigFile
	}
}

func (c *Config) ExtractRoleInfo() error {
	// validate
	if c.EnvironmentVariables.AwsProfile == "" {
		return errors.New("AWS_PROFILE is not fulfilled. AWS_PROFILE is required.")
	}

	// load AwsConfigFile
	iniData, err := ini.Load(c.FilePaths.AwsConfigFile)
	if err != nil {
		return err
	}

	// extract RoleInfo
	sectionName := "profile " + c.EnvironmentVariables.AwsProfile
	r := iniData.Section(sectionName)
	c.Role.RoleArn = r.Key("role_arn").String()
	c.Role.SourceProfile = r.Key("source_profile").String()
	c.Role.MfaSerial = r.Key("mfa_serial").String()
	c.Role.ExternalID = r.Key("external_id").String()

	return nil
}

func (c *Config) ExportSessionOptions() session.Options {
	return session.Options{
		Profile:                 c.EnvironmentVariables.AwsProfile,
		SharedConfigState:       session.SharedConfigEnable,
		AssumeRoleTokenProvider: stscreds.StdinTokenProvider,
	}
}

func (c *Config) GetRoleArn() string {
	return c.Role.RoleArn
}

package credentor

import (
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
)

func GetCredentials(arn string, options session.Options) (*credentials.Value, error) {
	sess := session.Must(session.NewSessionWithOptions(options))
	creds := stscreds.NewCredentials(sess, arn, func(p *stscreds.AssumeRoleProvider) {
		p.RoleSessionName = "credentor-session-" + time.Now().Format("20060102150405")
		p.Duration = time.Duration(60) * time.Minute
	})
	credentialsValue, err := creds.Get()
	if err != nil {
		return nil, err
	}

	return &credentialsValue, nil
}

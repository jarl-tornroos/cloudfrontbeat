package awsfacade

import (
	"os/user"
	"path/filepath"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/ec2rolecreds"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
)

// Auth is a common authenticator for AWS services
type Auth struct {
	Config          *aws.Config
	Session         *session.Session
	SessionWithConf *session.Session
}

// NewAuth tries to create the authentication in the following order
// Environment variables, IAM roles and credentials file in located in ~/.aws/credentials
func NewAuth(region string, environment string) *Auth {
	usr, _ := user.Current()
	credentialsFile := filepath.Join(usr.HomeDir, ".aws/credentials")

	sess := session.Must(session.NewSession())

	// Try first environment variables, then with IAM roles and lastly with user configuration
	creds := credentials.NewChainCredentials(
		[]credentials.Provider{
			&credentials.EnvProvider{},
			&ec2rolecreds.EC2RoleProvider{
				Client: ec2metadata.New(sess),
			},
			&credentials.SharedCredentialsProvider{
				Filename: credentialsFile,
				Profile:  environment,
			},
		})

	config := &aws.Config{
		Region:      aws.String(region),
		Credentials: creds,
	}

	return &Auth{
		Config:          config,
		Session:         sess,
		SessionWithConf: session.Must(session.NewSession(config)),
	}
}

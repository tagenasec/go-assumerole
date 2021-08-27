package assumerole

import (
	"context"

	"github.com/apex/log"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

type AssumeRole struct {
	cfg    aws.Config
	stsSvc *sts.Client
}

func NewAssumeRoleFromProfileName(profileName string) (*AssumeRole, error) {
	log.WithField("profile", profileName).Info("Configuring using profile name")
	awsConfig, err := config.LoadDefaultConfig(context.TODO(),
		config.WithSharedConfigProfile(profileName))
	if err != nil {
		log.WithField("profile", profileName).WithError(err).Error("Unable to load AWS config from profile name")
		return nil, err
	}
	awsConfig.Region = "us-east-1"
	return &AssumeRole{
		cfg:    awsConfig,
		stsSvc: sts.NewFromConfig(awsConfig),
	}, nil
}

package assumerole

import (
	"os"

	"github.com/apex/log"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

func SetEnvironmentForAssumedRole(assumedRole *sts.AssumeRoleOutput) error {
	os.Unsetenv("AWS_PROFILE")
	err := os.Setenv("AWS_ACCESS_KEY_ID", *assumedRole.Credentials.AccessKeyId)
	if err != nil {
		log.WithError(err).Error("Unable to set environment")
		return err
	}
	os.Setenv("AWS_SECRET_ACCESS_KEY", *assumedRole.Credentials.SecretAccessKey)
	if err != nil {
		log.WithError(err).Error("Unable to set environment")
		return err
	}
	os.Setenv("AWS_SESSION_TOKEN", *assumedRole.Credentials.SessionToken)
	if err != nil {
		log.WithError(err).Error("Unable to set environment")
		return err
	}
	return nil
}

func (self *AssumeRole) AssumeRoleAndOpenShell(roleArn string) error {
	credentials, err := self.AssumeRoleArn(roleArn)
	if err != nil {
		log.WithError(err).Error("Unable to assume role")
		return err
	}
	err = SetEnvironmentForAssumedRole(credentials)
	if err != nil {
		log.WithError(err).Error("Unable to set environment")
		return err
	}
	shell := os.Getenv("SHELL")
	err = doExec(shell, []string{}, map[string]string{})
	if err != nil {
		log.WithField("shell", shell).WithError(err).Error("Unable to exec SHELL")
		return err
	}
	return nil
}

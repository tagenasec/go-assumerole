package assumerole

import (
	"context"
	"fmt"
	"strings"

	"github.com/apex/log"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

func (self *AssumeRole) AssumeRoleArn(targetRoleArn string) (*sts.AssumeRoleOutput, error) {
	log.Info("Attempting to assume role")
	callerIdentity, err := self.stsSvc.GetCallerIdentity(context.TODO(), nil)
	if err != nil {
		log.WithError(err).Error("Unable to get caller identity")
		return nil, err
	}
	userName := strings.Split(*callerIdentity.Arn, "/")[1]
	iamSvc := iam.NewFromConfig(self.cfg)
	mfaDevicesList, err := iamSvc.ListMFADevices(context.TODO(), &iam.ListMFADevicesInput{
		UserName: &userName,
	})
	if err != nil {
		log.WithError(err).Error("Unable to fetch mfa devices")
		return nil, err
	}
	if len(mfaDevicesList.MFADevices) == 0 {
		log.Error("User has no mfa devices")
		return nil, fmt.Errorf("User has no mfa devices")
	}
	if len(mfaDevicesList.MFADevices) > 1 {
		log.Error("User has more than one MFA device")
		return nil, fmt.Errorf("User has more than one MFA device")
	}
	serialNumber := mfaDevicesList.MFADevices[0].SerialNumber
	token, err := stscreds.StdinTokenProvider()
	if err != nil {
		log.WithError(err).Error("Unable to get Token from user")
		return nil, err
	}
	roleSessionName := fmt.Sprintf("cli-by-%s", userName)
	credentials, err := self.stsSvc.AssumeRole(context.TODO(), &sts.AssumeRoleInput{
		RoleArn:         &targetRoleArn,
		RoleSessionName: &roleSessionName,
		SerialNumber:    serialNumber,
		TokenCode:       &token,
	})
	if err != nil {
		log.WithError(err).Error("Unable to assume role")
		return nil, err
	}
	log.Info("Successfully assumed role")
	return credentials, nil
}

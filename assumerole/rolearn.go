package assumerole

import (
	"context"
	"fmt"

	"github.com/apex/log"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

func RoleArnFromAccountAndRoleName(accountId string, roleName string) string {
	return fmt.Sprintf("arn:aws:iam::%s:role/%s", accountId, roleName)
}

func (self *AssumeRole) RoleArnFromRoleNameInThisAccount(roleName string) (string, error) {
	callerIdentity, err := self.stsSvc.GetCallerIdentity(context.TODO(), &sts.GetCallerIdentityInput{})
	if err != nil {
		log.WithError(err).Error("Unable to get caller identity")
	}
	return RoleArnFromAccountAndRoleName(*callerIdentity.Account, roleName), nil
}

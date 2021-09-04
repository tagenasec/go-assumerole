package main

import (
	"os"
	"strings"

	"github.com/apex/log"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/tagenasec/go-assumerole/assumerole"
)

func main() {
	log.SetLevel(log.DebugLevel)

	role := ""
	profile := ""
	rootCmd := &cobra.Command{
		Use:   "assumerole [options]",
		Short: "Open a new shell with aws environment variables for an assumed role",
		RunE: func(cmd *cobra.Command, args []string) error {
			if role == "" {
				return errors.Errorf("Must specify -r (or --role)")
			}
			if profile == "" {
				profile = "default"
			}
			assumeRole, err := assumerole.NewAssumeRoleFromProfileName(profile)
			if err != nil {
				log.WithError(err).Error("Unable to create assume role object")
				return err
			}
			if !strings.HasPrefix(role, "arn:") {
				roleArn, err := assumeRole.RoleArnFromRoleNameInThisAccount(role)
				if err != nil {
					log.WithError(err).Error("Unable to get role arn")
					return err
				}
				role = roleArn
			}
			err = assumeRole.AssumeRoleAndOpenShell(role)
			if err != nil {
				log.WithError(err).Error("Assume operation failed")
				return err
			}
			return nil
		},
	}
	rootCmd.PersistentFlags().StringVarP(&role, "role", "r", "", "role name to assume to")
	rootCmd.PersistentFlags().StringVarP(&profile, "profile", "p", os.Getenv("AWS_PROFILE"), "aws profile to assume from")

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}

	os.Exit(0)
}

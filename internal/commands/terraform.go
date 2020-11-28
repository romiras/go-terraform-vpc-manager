package commands

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/hashicorp/terraform-exec/tfinstall"
	"github.com/romiras/go-terraform-vpc-manager/internal/helpers"
	registry "github.com/romiras/go-terraform-vpc-manager/internal/registries"
)

const DefaultExecPath = "/usr/bin/terraform"

type ExecFinder struct{}

func (ef *ExecFinder) ExecPath(ctx context.Context) (string, error) {
	return registry.Reg.ExecPath, nil
}

func varOpt(instanceTagNamePrefix string) *tfexec.VarOption {
	return tfexec.Var("instance_tag_name=" + instanceTagNamePrefix)
}

// See https://github.com/hashicorp/terraform/issues/12917#issuecomment-437861756
func targetOption(instanceFullTagName string) *tfexec.TargetOption {
	// -target=aws_instance.docker-nginx-demo-instance-0
	resource := "aws_instance" + "." + instanceFullTagName
	return tfexec.Target(resource)
}

func tfPrepare() *tfexec.Terraform {
	_, err := os.Stat(registry.Reg.WorkingDir)
	if os.IsNotExist(err) {
		helpers.AbortOnError(err)
	}

	execPath, err := tfinstall.Find(context.Background(), &ExecFinder{})
	helpers.AbortOnError(err)

	tf, err := tfexec.NewTerraform(registry.Reg.WorkingDir, execPath)
	helpers.AbortOnError(err)

	// err = tf.Init(context.Background(), tfexec.Upgrade(true), tfexec.LockTimeout("60s"))
	// if err != nil {
	// 	panic(err)
	// }

	return tf
}

func CreateVPC(instanceName string) error {
	helpers.DebugMsg("CreateVPC: " + instanceName)

	tf := tfPrepare()
	if instanceName != "" {
		return tf.Apply(context.Background(), varOpt(instanceName))
	}
	return tf.Apply(context.Background())
}

func DestroyVPC(instanceName string) error {
	helpers.DebugMsg("DestroyVPC: " + instanceName)

	tf := tfPrepare()
	if instanceName != "" {
		return tf.Destroy(context.Background(), targetOption(instanceName))
	}
	return tf.Destroy(context.Background())
}

/*
this program should be a go binary with CLI flags
1. create
2. destroy
when creating there should be a sub-flag - instance_name
when destroying should be anther sub-flag - all, instance_name
you may use the cobra go CLI package

-create -instance_name name
-destroy -instance_name name
*/

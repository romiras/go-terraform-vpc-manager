package commands

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/hashicorp/terraform-exec/tfinstall"
)

const DefaultExecPath = "/usr/bin/terraform"
const workingDir = "/home/romiras/Devel/devops/terraform-vpc-docker-helloworld"

type ExecFinder struct{}

func (ef *ExecFinder) ExecPath(ctx context.Context) (string, error) {
	return DefaultExecPath, nil
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

func tfPrepare(workingDir string) *tfexec.Terraform {
	execPath, err := tfinstall.Find(context.Background(), &ExecFinder{})
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	tf, err := tfexec.NewTerraform(workingDir, execPath)
	if err != nil {
		panic(err)
	}

	// err = tf.Init(context.Background(), tfexec.Upgrade(true), tfexec.LockTimeout("60s"))
	// if err != nil {
	// 	panic(err)
	// }

	return tf
}

func CreateInstance(name string) error {
	fmt.Println("CreateInstance: " + name)

	tf := tfPrepare(workingDir)
	return tf.Apply(context.Background(), varOpt(name))
}

func DestroyInstance(name string) error {
	fmt.Println("DestroyInstance: " + name)

	tf := tfPrepare(workingDir)
	return tf.Destroy(context.Background(), targetOption(name))
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

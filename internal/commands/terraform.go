package commands

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/romiras/go-terraform-vpc-manager/internal/helpers"
	registry "github.com/romiras/go-terraform-vpc-manager/internal/registries"
)

// See https://github.com/hashicorp/terraform/issues/12917#issuecomment-437861756
func targetOption(instanceFullTagName string) *tfexec.TargetOption {
	// -target=aws_instance.docker-nginx-demo-instance-0
	resource := "aws_instance" + "." + instanceFullTagName
	return tfexec.Target(resource)
}

func tfPrepare() (*tfexec.Terraform, error) {
	_, err := os.Stat(registry.Reg.WorkingDir)
	if os.IsNotExist(err) {
		helpers.AbortOnError(err)
	}

	tf, err := tfexec.NewTerraform(registry.Reg.WorkingDir, registry.Reg.ExecPath)
	if err != nil {
		return nil, err
	}

	return tf, nil
}

func CreateVPC(instanceName string) error {
	helpers.DebugMsg("CreateVPC: " + instanceName)

	tf, err := tfPrepare()
	if err != nil {
		return err
	}

	initOptions := []tfexec.InitOption{
		// tfexec.Upgrade(true),
		// tfexec.LockTimeout("60s"),
	}

	err = tf.Init(context.TODO(), initOptions...)
	if err != nil {
		return err
	}

	applyOptions := make([]tfexec.ApplyOption, 0)
	if instanceName != "" {
		applyOptions = append(applyOptions, tfexec.Var("instance_tag_name="+instanceName))
	}

	return tf.Apply(context.TODO(), applyOptions...)
}

func DestroyVPC(instanceName string) error {
	helpers.DebugMsg("DestroyVPC: " + instanceName)

	tf, err := tfPrepare()
	if err != nil {
		return err
	}

	destroyOptions := make([]tfexec.DestroyOption, 0)
	if instanceName != "" {
		destroyOptions = append(destroyOptions, targetOption(instanceName))
	}

	return tf.Destroy(context.TODO(), destroyOptions...)
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

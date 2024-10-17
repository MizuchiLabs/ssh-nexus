package provider

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pocketbase/pocketbase/models"
)

// AWSProvider represents a provider implementation for AWS
type AWSProvider struct {
	Config *models.Record
}

func NewAWSProvider(config *models.Record) *AWSProvider {
	return &AWSProvider{Config: config}
}

func (p *AWSProvider) Sync() ([]ProviderMachine, error) {
	sess, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(
			p.Config.GetString("username"),
			p.Config.GetString("password"),
			p.Config.GetString("token"),
		),
	})
	if err != nil {
		return nil, err
	}

	svc := ec2.New(sess)
	instances, err := svc.DescribeInstances(nil)
	if err != nil {
		return nil, err
	}

	var machines []ProviderMachine
	for _, reservation := range instances.Reservations {
		for _, instance := range reservation.Instances {
			if instance.PublicIpAddress != nil && instance.InstanceId != nil {
				machine := ProviderMachine{
					ID:      *instance.InstanceId,
					Host:    *instance.PublicIpAddress,
					Name:    *instance.InstanceId,
					Running: *instance.State.Name == "running",
				}
				machines = append(machines, machine)
			}
		}
	}

	return machines, nil
}

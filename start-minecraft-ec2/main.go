package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

func startInstances(cfg aws.Config, region string) error {
	log.Println("#" + region)

	client := ec2.NewFromConfig(cfg, func(o *ec2.Options) {
		o.Region = region
	})
	dii := &ec2.DescribeInstancesInput{
		Filters: []types.Filter{
			{
				Name:   aws.String("tag:Service"),
				Values: []string{*aws.String("Minecraft")},
			},
			{
				Name:   aws.String("tag:AutoStart"),
				Values: []string{*aws.String("True")},
			},
		},
	}

	dio, err := client.DescribeInstances(context.TODO(), dii)
	if err != nil {
		return err
	}

	ids := make([]string, 0)
	for _, r := range dio.Reservations {
		for _, i := range r.Instances {
			log.Printf("* %s", *i.InstanceId)
			ids = append(ids, *i.InstanceId)
		}
	}

	if len(ids) == 0 {
		return nil
	}

	sii := &ec2.StartInstancesInput{
		InstanceIds: ids,
	}

	_, err = client.StartInstances(context.TODO(), sii)
	if err != nil {
		return err
	}

	return nil
}

func handleRequest(ctx context.Context) (string, error) {
	log.Println("Start Minecraft EC2 instances.")
	regions := []string{"ap-northeast-1", "ap-northeast-3"}
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Println(err)
		return "", err
	}

	for _, r := range regions {
		err = startInstances(cfg, r)
		if err != nil {
			log.Println(err)
			return "", err
		}
	}

	return "", nil
}

func main() {
	lambda.Start(handleRequest)
}

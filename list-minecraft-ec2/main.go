package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

func listInstances(cfg aws.Config, region string) error {
	fmt.Println("#" + region)

	client := ec2.NewFromConfig(cfg, func(o *ec2.Options) {
		o.Region = region
	})
	input := &ec2.DescribeInstancesInput{
		Filters: []types.Filter{
			{
				Name:   aws.String("tag:Service"),
				Values: []string{*aws.String("Minecraft")},
			},
		},
	}

	output, err := client.DescribeInstances(context.TODO(), input)
	if err != nil {
		return err
	}

	for _, r := range output.Reservations {
		for _, i := range r.Instances {
			hostName := getHostNameTag(i.Tags)
			fmt.Printf("%s\t%s\t%s\n", *i.InstanceId, i.State.Name, hostName)
		}
	}

	return nil
}

func getHostNameTag(tags []types.Tag) string {
	for _, t := range tags {
		if *t.Key == "HostName" {
			return *t.Value
		}
	}

	return ""
}

func handleRequest(ctx context.Context) (string, error) {
	regions := []string{"ap-northeast-1", "ap-northeast-3"}
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	for _, r := range regions {
		err = listInstances(cfg, r)
		if err != nil {
			fmt.Println(err)
			return "", err
		}
	}

	return "", nil
}

func main() {
	lambda.Start(handleRequest)
}

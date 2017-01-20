package aws

import (
	"errors"
	"fmt"
	"reflect"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/wallix/awless/rdf"
)

type propertyTransform struct {
	name      string
	transform transformFn
}

type transformFn func(i interface{}) (interface{}, error)

var ErrTagNotFound = errors.New("aws tag key not found")
var ErrFieldNotFound = errors.New("aws struct field not found")

var extractValueFn = func(i interface{}) (interface{}, error) {
	iv := reflect.ValueOf(i)
	if iv.Kind() == reflect.Ptr {
		return iv.Elem().Interface(), nil
	}
	return nil, fmt.Errorf("aws type unknown: %T", i)
}

var extractFieldFn = func(field string) transformFn {
	return func(i interface{}) (interface{}, error) {
		value := reflect.ValueOf(i)
		struc := value.Elem()

		structField := struc.FieldByName(field)

		if !structField.IsValid() {
			return nil, ErrFieldNotFound
		}

		return extractValueFn(structField.Interface())
	}
}

var extractTagFn = func(key string) transformFn {
	return func(i interface{}) (interface{}, error) {
		tags, ok := i.([]*ec2.Tag)
		if !ok {
			return nil, fmt.Errorf("aws model: unexpected type %T", i)
		}
		for _, t := range tags {
			if key == awssdk.StringValue(t.Key) {
				return awssdk.StringValue(t.Value), nil
			}
		}

		return nil, ErrTagNotFound
	}
}

var awsResourcesDef = map[rdf.ResourceType]map[string]*propertyTransform{
	rdf.Instance: {
		"Id":         {name: "InstanceId", transform: extractValueFn},
		"Name":       {name: "Tags", transform: extractTagFn("Name")},
		"Type":       {name: "InstanceType", transform: extractValueFn},
		"SubnetId":   {name: "SubnetId", transform: extractValueFn},
		"VpcId":      {name: "VpcId", transform: extractValueFn},
		"PublicIp":   {name: "PublicIpAddress", transform: extractValueFn},
		"PrivateIp":  {name: "PrivateIpAddress", transform: extractValueFn},
		"ImageId":    {name: "ImageId", transform: extractValueFn},
		"LaunchTime": {name: "LaunchTime", transform: extractValueFn},
		"State":      {name: "State", transform: extractFieldFn("Name")},
		"KeyName":    {name: "KeyName", transform: extractValueFn},
	},
	rdf.Vpc: {
		"Id":        {name: "VpcId", transform: extractValueFn},
		"Name":      {name: "Tags", transform: extractTagFn("Name")},
		"IsDefault": {name: "IsDefault", transform: extractValueFn},
		"State":     {name: "State", transform: extractValueFn},
		"CidrBlock": {name: "CidrBlock", transform: extractValueFn},
	},
	rdf.Subnet: {
		"Id":                  {name: "SubnetId", transform: extractValueFn},
		"Name":                {name: "Tags", transform: extractTagFn("Name")},
		"VpcId":               {name: "VpcId", transform: extractValueFn},
		"MapPublicIpOnLaunch": {name: "MapPublicIpOnLaunch", transform: extractValueFn},
		"State":               {name: "State", transform: extractValueFn},
		"CidrBlock":           {name: "CidrBlock", transform: extractValueFn},
	},
	rdf.SecurityGroup: {
		"Id":      {name: "GroupId", transform: extractValueFn},
		"Name":    {name: "GroupName", transform: extractValueFn},
		"OwnerId": {name: "OwnerId", transform: extractValueFn},
		"VpcId":   {name: "VpcId", transform: extractValueFn},
	},
	rdf.User: {
		"Id":                   {name: "UserId", transform: extractValueFn},
		"Name":                 {name: "UserName", transform: extractValueFn},
		"Arn":                  {name: "Arn", transform: extractValueFn},
		"Path":                 {name: "Path", transform: extractValueFn},
		"PasswordLastUsedDate": {name: "PasswordLastUsed", transform: extractValueFn},
	},
	rdf.Role: {
		"Id":         {name: "RoleId", transform: extractValueFn},
		"Name":       {name: "RoleName", transform: extractValueFn},
		"Arn":        {name: "Arn", transform: extractValueFn},
		"CreateDate": {name: "CreateDate", transform: extractValueFn},
		"Path":       {name: "Path", transform: extractValueFn},
	},
	rdf.Group: {
		"Id":         {name: "GroupId", transform: extractValueFn},
		"Name":       {name: "GroupName", transform: extractValueFn},
		"Arn":        {name: "Arn", transform: extractValueFn},
		"CreateDate": {name: "CreateDate", transform: extractValueFn},
		"Path":       {name: "Path", transform: extractValueFn},
	},
	rdf.Policy: {
		"Id":           {name: "PolicyId", transform: extractValueFn},
		"Name":         {name: "PolicyName", transform: extractValueFn},
		"Arn":          {name: "Arn", transform: extractValueFn},
		"CreateDate":   {name: "CreateDate", transform: extractValueFn},
		"UpdateDate":   {name: "UpdateDate", transform: extractValueFn},
		"Description":  {name: "Description", transform: extractValueFn},
		"IsAttachable": {name: "IsAttachable", transform: extractValueFn},
		"Path":         {name: "Path", transform: extractValueFn},
	},
}

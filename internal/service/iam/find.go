package iam

import (
	"context"
	"regexp"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func FindGroupAttachedPolicy(conn *iam.IAM, groupARN string, policyARN string) (*iam.Policy, error) {
	input := &iam.ListGroupPoliciesInput{
		GroupArn: aws.String(groupARN),
	}

	output, err := conn.ListGroupPolicies(input)

	if err != nil {
		return nil, err
	}

	if output == nil {
		return nil, tfresource.NewEmptyResultError(input)
	}

	for _, policy := range output.Policies {
		if policy == nil {
			continue
		}

		if aws.StringValue(policy.PolicyArn) == policyARN {
			return policy, nil
		}
	}

	return nil, &retry.NotFoundError{}
}

func FindUserAttachedGlobalPolicy(conn *iam.IAM, userName, policyARN string) (*iam.Policy, error) {
	input := &iam.ListUserGlobalPoliciesInput{
		UserName: aws.String(userName),
	}

	output, err := conn.ListUserGlobalPolicies(input)

	if err != nil {
		return nil, err
	}

	if output == nil {
		return nil, tfresource.NewEmptyResultError(input)
	}

	for _, policy := range output.Policies {
		if policy == nil {
			continue
		}

		if aws.StringValue(policy.PolicyArn) == policyARN {
			return policy, nil
		}
	}

	return nil, &retry.NotFoundError{}
}

func FindUserAttachedProjectPolicy(conn *iam.IAM, userName, policyARN, projectName string) (*iam.Policy, error) {
	input := &iam.ListUserProjectPoliciesInput{
		UserName:    aws.String(userName),
		ProjectName: aws.String(projectName),
	}

	output, err := conn.ListUserProjectPolicies(input)

	if err != nil {
		return nil, err
	}

	if output == nil {
		return nil, tfresource.NewEmptyResultError(input)
	}

	for _, policy := range output.Policies {
		if policy == nil {
			continue
		}

		if aws.StringValue(policy.PolicyArn) == policyARN {
			return policy, nil
		}
	}

	return nil, &retry.NotFoundError{}
}

// FindUserAttachedPolicy returns the AttachedPolicy corresponding to the specified user and policy ARN.
func FindUserAttachedPolicy(conn *iam.IAM, userName string, policyARN string) (*iam.AttachedPolicy, error) {
	input := &iam.ListAttachedUserPoliciesInput{
		UserName: aws.String(userName),
	}

	var result *iam.AttachedPolicy

	err := conn.ListAttachedUserPoliciesPages(input, func(page *iam.ListAttachedUserPoliciesOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, attachedPolicy := range page.AttachedPolicies {
			if attachedPolicy == nil {
				continue
			}

			if aws.StringValue(attachedPolicy.PolicyArn) == policyARN {
				result = attachedPolicy
				return false
			}
		}

		return !lastPage
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

// FindPolicies returns the FindPolicies corresponding to the specified ARN, name, and/or path-prefix.
func FindPolicies(conn *iam.IAM, arn, name, pathPrefix string) ([]*iam.Policy, error) {
	input := &iam.ListPoliciesInput{}

	if pathPrefix != "" {
		input.PathPrefix = aws.String(pathPrefix)
	}

	var results []*iam.Policy

	err := conn.ListPoliciesPages(input, func(page *iam.ListPoliciesOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, p := range page.Policies {
			if p == nil {
				continue
			}

			if arn != "" && arn != aws.StringValue(p.PolicyArn) {
				continue
			}

			if name != "" && name != aws.StringValue(p.PolicyName) {
				continue
			}

			results = append(results, p)
		}

		return !lastPage
	})

	return results, err
}

func FindPolicyByArn(conn *iam.IAM, arn string) (*iam.Policy, error) {
	input := &iam.GetPolicyInput{
		PolicyArn: aws.String(arn),
	}

	output, err := conn.GetPolicy(input)

	if tfawserr.ErrCodeEquals(err, PolicyNotFoundCode) {
		return nil, &retry.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if err != nil {
		return nil, err
	}

	if output == nil || output.Policy == nil {
		return nil, tfresource.NewEmptyResultError(input)
	}

	return output.Policy, nil
}

func FindGroups(conn *iam.IAM, name, groupType string) ([]*iam.Group, error) {
	input := &iam.ListGroupsInput{}

	if groupType != "" {
		input.Type = aws.String(groupType)
	}

	var results []*iam.Group

	err := conn.ListGroupsPages(input, func(page *iam.ListGroupsOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, g := range page.Groups {
			if g == nil {
				continue
			}

			if name != "" && name != aws.StringValue(g.GroupName) {
				continue
			}

			results = append(results, g)
		}

		return !lastPage
	})

	return results, err
}

func FindGroupByArn(conn *iam.IAM, arn string) (*iam.Group, []*iam.User, error) {
	input := &iam.GetGroupInput{
		GroupArn: aws.String(arn),
	}

	output, err := conn.GetGroup(input)

	if tfawserr.ErrCodeEquals(err, GroupNotFoundCode) {
		return nil, nil, &retry.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if err != nil {
		return nil, nil, err
	}

	if output == nil || output.Group == nil {
		return nil, nil, tfresource.NewEmptyResultError(input)
	}

	return output.Group, output.Users, nil
}

func FindUserGlobalGroups(conn *iam.IAM, userName string) ([]*iam.Group, error) {
	input := &iam.ListUserGlobalGroupsInput{
		UserName: aws.String(userName),
	}

	output, err := conn.ListUserGlobalGroups(input)

	if err != nil {
		return nil, err
	}

	if output == nil {
		return nil, tfresource.NewEmptyResultError(input)
	}

	return output.Groups, nil
}

func FindUserProjectGroups(conn *iam.IAM, userName, projectName string) ([]*iam.Group, error) {
	input := &iam.ListUserProjectGroupsInput{
		UserName:    aws.String(userName),
		ProjectName: aws.String(projectName),
	}

	output, err := conn.ListUserProjectGroups(input)

	if err != nil {
		return nil, err
	}

	if output == nil {
		return nil, tfresource.NewEmptyResultError(input)
	}

	return output.Groups, nil
}

func FindUsers(conn *iam.IAM, nameRegex, pathPrefix string) ([]*iam.User, error) {
	input := &iam.ListUsersInput{}

	if pathPrefix != "" {
		input.PathPrefix = aws.String(pathPrefix)
	}

	var results []*iam.User

	err := conn.ListUsersPages(input, func(page *iam.ListUsersOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, user := range page.Users {
			if user == nil {
				continue
			}

			if nameRegex != "" && !regexp.MustCompile(nameRegex).MatchString(aws.StringValue(user.UserName)) {
				continue
			}

			results = append(results, user)
		}

		return !lastPage
	})

	return results, err
}

func FindUserByName(conn *iam.IAM, name string) (*iam.User, error) {
	input := &iam.GetUserInput{
		UserName: aws.String(name),
	}

	output, err := conn.GetUser(input)

	if tfawserr.ErrCodeEquals(err, UserNotFoundCode) {
		return nil, &retry.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if err != nil {
		return nil, err
	}

	if output == nil || output.User == nil {
		return nil, tfresource.NewEmptyResultError(input)
	}

	return output.User, nil
}

func FindProjectByName(conn *iam.IAM, name string) (*iam.Project, error) {
	input := &iam.GetProjectInput{
		ProjectName: aws.String(name),
	}

	output, err := conn.GetProject(input)

	if tfawserr.ErrCodeEquals(err, ProjectNotFoundCode) {
		return nil, &retry.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if err != nil {
		return nil, err
	}

	if output == nil || output.Project == nil {
		return nil, tfresource.NewEmptyResultError(input)
	}

	project := output.Project

	if aws.StringValue(project.State) == iam.ProjectStateTypeDeleted {
		return nil, &retry.NotFoundError{}
	}

	return project, nil
}

func FindRoleByName(conn *iam.IAM, name string) (*iam.Role, error) {
	input := &iam.GetRoleInput{
		RoleName: aws.String(name),
	}

	output, err := conn.GetRole(input)

	if tfawserr.ErrCodeEquals(err, iam.ErrCodeNoSuchEntityException) {
		return nil, &resource.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if err != nil {
		return nil, err
	}

	if output == nil || output.Role == nil {
		return nil, tfresource.NewEmptyResultError(input)
	}

	return output.Role, nil
}

func FindVirtualMFADevice(conn *iam.IAM, serialNum string) (*iam.VirtualMFADevice, error) {
	input := &iam.ListVirtualMFADevicesInput{}

	output, err := conn.ListVirtualMFADevices(input)

	if err != nil {
		return nil, err
	}

	if len(output.VirtualMFADevices) == 0 || output.VirtualMFADevices[0] == nil {
		return nil, tfresource.NewEmptyResultError(output)
	}

	var device *iam.VirtualMFADevice

	for _, dvs := range output.VirtualMFADevices {
		if aws.StringValue(dvs.SerialNumber) == serialNum {
			device = dvs
			break
		}
	}

	if device == nil {
		return nil, tfresource.NewEmptyResultError(device)
	}

	return device, nil
}

func FindServiceSpecificCredential(conn *iam.IAM, serviceName, userName, credID string) (*iam.ServiceSpecificCredentialMetadata, error) {
	input := &iam.ListServiceSpecificCredentialsInput{
		ServiceName: aws.String(serviceName),
		UserName:    aws.String(userName),
	}

	output, err := conn.ListServiceSpecificCredentials(input)

	if tfawserr.ErrCodeEquals(err, iam.ErrCodeNoSuchEntityException) {
		return nil, &resource.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if err != nil {
		return nil, err
	}

	if len(output.ServiceSpecificCredentials) == 0 || output.ServiceSpecificCredentials[0] == nil {
		return nil, tfresource.NewEmptyResultError(output)
	}

	var cred *iam.ServiceSpecificCredentialMetadata

	for _, crd := range output.ServiceSpecificCredentials {
		if aws.StringValue(crd.ServiceName) == serviceName &&
			aws.StringValue(crd.UserName) == userName &&
			aws.StringValue(crd.ServiceSpecificCredentialId) == credID {
			cred = crd
			break
		}
	}

	if cred == nil {
		return nil, tfresource.NewEmptyResultError(cred)
	}

	return cred, nil
}

func FindSigningCertificate(conn *iam.IAM, userName, certId string) (*iam.SigningCertificate, error) {
	input := &iam.ListSigningCertificatesInput{
		UserName: aws.String(userName),
	}

	output, err := conn.ListSigningCertificates(input)

	if tfawserr.ErrCodeEquals(err, iam.ErrCodeNoSuchEntityException) {
		return nil, &resource.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if err != nil {
		return nil, err
	}

	if len(output.Certificates) == 0 || output.Certificates[0] == nil {
		return nil, tfresource.NewEmptyResultError(output)
	}

	var cert *iam.SigningCertificate

	for _, crt := range output.Certificates {
		if aws.StringValue(crt.UserName) == userName &&
			aws.StringValue(crt.CertificateId) == certId {
			cert = crt
			break
		}
	}

	if cert == nil {
		return nil, tfresource.NewEmptyResultError(cert)
	}

	return cert, nil
}

func FindSAMLProviderByARN(ctx context.Context, conn *iam.IAM, arn string) (*iam.GetSAMLProviderOutput, error) {
	input := &iam.GetSAMLProviderInput{
		SAMLProviderArn: aws.String(arn),
	}

	output, err := conn.GetSAMLProviderWithContext(ctx, input)

	if tfawserr.ErrCodeEquals(err, iam.ErrCodeNoSuchEntityException) {
		return nil, &resource.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if err != nil {
		return nil, err
	}

	if output == nil {
		return nil, tfresource.NewEmptyResultError(input)
	}

	return output, nil
}

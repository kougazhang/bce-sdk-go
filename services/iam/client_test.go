/*
 * Copyright 2021 Baidu, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
 * except in compliance with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the
 * License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions
 * and limitations under the License.
 */

package iam

import (
	"encoding/json"
	"fmt"
	"github.com/kougazhang/bce-sdk-go/services/iam/api"
	"github.com/kougazhang/bce-sdk-go/util/log"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

// For security reason, ak/sk should not hard write here.
type Conf struct {
	AK string
	SK string
}

var IAM_CLIENT *Client

func init() {
	_, f, _, _ := runtime.Caller(0)
	for i := 0; i < 7; i++ {
		f = filepath.Dir(f)
	}
	conf := filepath.Join(f, "config.json")
	fp, err := os.Open(conf)
	if err != nil {
		fmt.Printf("config json file of ak/sk not given: %+v\n", conf)
		os.Exit(1)
	}
	decoder := json.NewDecoder(fp)
	confObj := &Conf{}
	decoder.Decode(confObj)
	IAM_CLIENT, _ = NewClient(confObj.AK, confObj.SK)
	log.SetLogLevel(log.DEBUG)
}

// ExpectEqual is the helper function for test each case
func ExpectEqual(alert func(format string, args ...interface{}),
	expected interface{}, actual interface{}) bool {
	expectedValue, actualValue := reflect.ValueOf(expected), reflect.ValueOf(actual)
	equal := false
	switch {
	case expected == nil && actual == nil:
		return true
	case expected != nil && actual == nil:
		equal = expectedValue.IsNil()
	case expected == nil && actual != nil:
		equal = actualValue.IsNil()
	default:
		if actualType := reflect.TypeOf(actual); actualType != nil {
			if expectedValue.IsValid() && expectedValue.Type().ConvertibleTo(actualType) {
				equal = reflect.DeepEqual(expectedValue.Convert(actualType).Interface(), actual)
			}
		}
	}
	if !equal {
		_, file, line, _ := runtime.Caller(1)
		alert("%s:%d: missmatch, expect %v but %v", file, line, expected, actual)
		return false
	}
	return true
}

func TestCreateListUpdateDeleteUser(t *testing.T) {
	name := "test-user-sdk-go"
	args := &api.CreateUserArgs{
		Name:        name,
		Description: "description",
	}
	res, err := IAM_CLIENT.CreateUser(args)
	ExpectEqual(t.Errorf, err, nil)
	jsonRes, _ := json.Marshal(res)
	t.Logf(string(jsonRes))
	ExpectEqual(t.Errorf, res.Name, args.Name)
	ExpectEqual(t.Errorf, res.Description, args.Description)

	users, err := IAM_CLIENT.ListUser()
	ExpectEqual(t.Errorf, err, nil)
	if users == nil || len(users.Users) == 0 {
		t.Errorf("list user return no result")
	}
	jsonRes, _ = json.Marshal(users)
	t.Logf(string(jsonRes))

	updateArgs := &api.UpdateUserArgs{
		Description: "updated description",
	}
	updated, err := IAM_CLIENT.UpdateUser(name, updateArgs)
	ExpectEqual(t.Errorf, err, nil)
	jsonRes, _ = json.Marshal(updated)
	t.Logf(string(jsonRes))
	ExpectEqual(t.Errorf, updated.Name, name)
	ExpectEqual(t.Errorf, updated.Description, updateArgs.Description)

	err = IAM_CLIENT.DeleteUser(name)
	ExpectEqual(t.Errorf, err, nil)
}

func TestUpdateGetDeleteUserLoginProfile(t *testing.T) {
	name := "test-user-sdk-go-login-profile"
	args := &api.CreateUserArgs{
		Name: name,
	}
	_, err := IAM_CLIENT.CreateUser(args)
	ExpectEqual(t.Errorf, err, nil)

	updateArgs := &api.UpdateUserLoginProfileArgs{
		Password:        "1@3Qwe4f",
		EnabledLoginMfa: true,
		LoginMfaType:    "PHONE",
	}
	updateRes, err := IAM_CLIENT.UpdateUserLoginProfile(name, updateArgs)
	ExpectEqual(t.Errorf, err, nil)
	jsonRes, _ := json.Marshal(updateRes)
	t.Logf(string(jsonRes))
	ExpectEqual(t.Errorf, updateRes.EnabledLoginMfa, true)
	ExpectEqual(t.Errorf, updateRes.LoginMfaType, "PHONE")

	getRes, err := IAM_CLIENT.GetUserLoginProfile(name)
	ExpectEqual(t.Errorf, err, nil)
	jsonRes, _ = json.Marshal(getRes)
	t.Logf(string(jsonRes))
	ExpectEqual(t.Errorf, updateRes.EnabledLoginMfa, true)
	ExpectEqual(t.Errorf, updateRes.LoginMfaType, "PHONE")

	err = IAM_CLIENT.DeleteUser(name)
	ExpectEqual(t.Errorf, err, nil)
}

func TestCreateGroupUpdateGetListDeleteGroup(t *testing.T) {
	name := "test_sdk_go_group"
	args := &api.CreateUserArgs{
		Name: name,
	}
	_, err := IAM_CLIENT.CreateUser(args)
	ExpectEqual(t.Errorf, err, nil)

	groupArgs := &api.CreateGroupArgs{
		Name:        name,
		Description: "description",
	}
	group, err := IAM_CLIENT.CreateGroup(groupArgs)
	ExpectEqual(t.Errorf, err, nil)
	ExpectEqual(t.Errorf, name, group.Name)
	ExpectEqual(t.Errorf, groupArgs.Description, group.Description)

	updateGroupArgs := &api.UpdateGroupArgs{
		Description: "updated group",
	}
	updated, err := IAM_CLIENT.UpdateGroup(name, updateGroupArgs)
	ExpectEqual(t.Errorf, err, nil)
	ExpectEqual(t.Errorf, name, updated.Name)
	ExpectEqual(t.Errorf, updateGroupArgs.Description, updated.Description)

	getRes, err := IAM_CLIENT.GetGroup(name)
	ExpectEqual(t.Errorf, err, nil)
	ExpectEqual(t.Errorf, name, getRes.Name)
	ExpectEqual(t.Errorf, updateGroupArgs.Description, getRes.Description)

	listRes, err := IAM_CLIENT.ListGroup()
	ExpectEqual(t.Errorf, err, nil)
	if listRes == nil || len(listRes.Groups) == 0 {
		t.Errorf("list group return no result")
	}

	err = IAM_CLIENT.DeleteUser(name)
	ExpectEqual(t.Errorf, err, nil)
	err = IAM_CLIENT.DeleteGroup(name)
	ExpectEqual(t.Errorf, err, nil)
}

func TestAddDeleteUserFromGroup(t *testing.T) {
	name := "test_sdk_go_group"
	args := &api.CreateUserArgs{
		Name: name,
	}
	user, err := IAM_CLIENT.CreateUser(args)
	ExpectEqual(t.Errorf, err, nil)

	groupArgs := &api.CreateGroupArgs{
		Name:        name,
		Description: "description",
	}
	group, err := IAM_CLIENT.CreateGroup(groupArgs)
	ExpectEqual(t.Errorf, err, nil)

	err = IAM_CLIENT.AddUserToGroup(name, name)
	ExpectEqual(t.Errorf, err, nil)

	usersRes, err := IAM_CLIENT.ListUsersInGroup(name)
	ExpectEqual(t.Errorf, err, nil)
	if usersRes == nil || len(usersRes.Users) != 1 {
		t.Errorf("list group result not 1")
	}
	ExpectEqual(t.Errorf, 1, len(usersRes.Users))
	ExpectEqual(t.Errorf, user.Id, usersRes.Users[0].Id)
	ExpectEqual(t.Errorf, user.Name, usersRes.Users[0].Name)

	groupsRes, err := IAM_CLIENT.ListGroupsForUser(name)
	ExpectEqual(t.Errorf, err, nil)
	if groupsRes == nil || len(groupsRes.Groups) != 1 {
		t.Errorf("list user result not 1")
	}
	ExpectEqual(t.Errorf, 1, len(groupsRes.Groups))
	ExpectEqual(t.Errorf, group.Id, groupsRes.Groups[0].Id)
	ExpectEqual(t.Errorf, group.Name, groupsRes.Groups[0].Name)

	err = IAM_CLIENT.DeleteUserFromGroup(name, name)
	ExpectEqual(t.Errorf, err, nil)

	usersRes, err = IAM_CLIENT.ListUsersInGroup(name)
	ExpectEqual(t.Errorf, err, nil)
	ExpectEqual(t.Errorf, 0, len(usersRes.Users))

	groupsRes, err = IAM_CLIENT.ListGroupsForUser(name)
	ExpectEqual(t.Errorf, err, nil)
	ExpectEqual(t.Errorf, 0, len(groupsRes.Groups))

	err = IAM_CLIENT.DeleteUser(name)
	ExpectEqual(t.Errorf, err, nil)

	err = IAM_CLIENT.DeleteGroup(name)
	ExpectEqual(t.Errorf, err, nil)
}

func getPolicyDocument() string {
	aclEntry := api.AclEntry{
		Service:    "bos",
		Region:     "bj",
		Permission: []string{"ListBucket"},
		Resource:   []string{"*"},
		Effect:     "Allow",
	}
	acl := &api.Acl{
		AccessControlList: []api.AclEntry{aclEntry},
	}
	document, _ := json.Marshal(acl)
	return string(document)
}

func TestCreateGetListDeletePolicy(t *testing.T) {
	name := "test_sdk_go_policy"
	args := &api.CreatePolicyArgs{
		Name:        name,
		Description: "description",
		Document:    getPolicyDocument(),
	}

	res, err := IAM_CLIENT.CreatePolicy(args)
	ExpectEqual(t.Errorf, err, nil)
	ExpectEqual(t.Errorf, name, res.Name)
	ExpectEqual(t.Errorf, args.Description, res.Description)

	getRes, err := IAM_CLIENT.GetPolicy(name, "")
	ExpectEqual(t.Errorf, err, nil)
	ExpectEqual(t.Errorf, name, getRes.Name)
	ExpectEqual(t.Errorf, args.Description, getRes.Description)

	listRes, err := IAM_CLIENT.ListPolicy(name, "")
	ExpectEqual(t.Errorf, err, nil)
	if listRes == nil || len(listRes.Policies) == 0 {
		t.Errorf("list policy result is empty")
	}

	err = IAM_CLIENT.DeletePolicy(name)
	ExpectEqual(t.Errorf, err, nil)
}

func TestUserAttachDetachPolicy(t *testing.T) {
	userName := "test_sdk_go_policy"
	args := &api.CreateUserArgs{
		Name: userName,
	}
	_, err := IAM_CLIENT.CreateUser(args)
	ExpectEqual(t.Errorf, err, nil)

	policyName := "test_sdk_go_policy"
	policyArgs := &api.CreatePolicyArgs{
		Name:        policyName,
		Description: "description",
		Document:    getPolicyDocument(),
	}

	_, err = IAM_CLIENT.CreatePolicy(policyArgs)
	ExpectEqual(t.Errorf, err, nil)

	attachArgs := &api.AttachPolicyToUserArgs{
		UserName:   userName,
		PolicyName: policyName,
	}
	err = IAM_CLIENT.AttachPolicyToUser(attachArgs)
	ExpectEqual(t.Errorf, err, nil)

	policies, err := IAM_CLIENT.ListUserAttachedPolicies(userName)
	ExpectEqual(t.Errorf, err, nil)
	if policies == nil || len(policies.Policies) != 1 {
		t.Errorf("list policy result is not 1")
	}
	policy := policies.Policies[0]
	ExpectEqual(t.Errorf, policyName, policy.Name)

	detachArgs := &api.DetachPolicyFromUserArgs{
		UserName:   userName,
		PolicyName: policyName,
	}
	err = IAM_CLIENT.DetachPolicyFromUser(detachArgs)
	ExpectEqual(t.Errorf, err, nil)

	err = IAM_CLIENT.DeletePolicy(policyName)
	ExpectEqual(t.Errorf, err, nil)
	err = IAM_CLIENT.DeleteUser(userName)
	ExpectEqual(t.Errorf, err, nil)
}

func TestGroupAttachDetachPolicy(t *testing.T) {
	groupName := "test_sdk_go_policy"
	args := &api.CreateGroupArgs{
		Name: groupName,
	}
	_, err := IAM_CLIENT.CreateGroup(args)
	ExpectEqual(t.Errorf, err, nil)

	policyName := "test_sdk_go_policy"
	policyArgs := &api.CreatePolicyArgs{
		Name:        policyName,
		Description: "description",
		Document:    getPolicyDocument(),
	}

	_, err = IAM_CLIENT.CreatePolicy(policyArgs)
	ExpectEqual(t.Errorf, err, nil)

	attachArgs := &api.AttachPolicyToGroupArgs{
		GroupName:  groupName,
		PolicyName: policyName,
	}
	err = IAM_CLIENT.AttachPolicyToGroup(attachArgs)
	ExpectEqual(t.Errorf, err, nil)

	policies, err := IAM_CLIENT.ListGroupAttachedPolicies(groupName)
	ExpectEqual(t.Errorf, err, nil)
	if policies == nil || len(policies.Policies) != 1 {
		t.Errorf("list policy result is not 1")
	}
	policy := policies.Policies[0]
	ExpectEqual(t.Errorf, policyName, policy.Name)

	detachArgs := &api.DetachPolicyFromGroupArgs{
		GroupName:  groupName,
		PolicyName: policyName,
	}
	err = IAM_CLIENT.DetachPolicyFromGroup(detachArgs)
	ExpectEqual(t.Errorf, err, nil)

	err = IAM_CLIENT.DeletePolicy(policyName)
	ExpectEqual(t.Errorf, err, nil)
	err = IAM_CLIENT.DeleteGroup(groupName)
	ExpectEqual(t.Errorf, err, nil)
}

package controllers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/vedicsoft/vamps-core/commons"
)

const SPLIT_SYMBOL string = "."
const ALL_SYMBOL string = "*"

type Statement struct {
	Id         string   `json:"id"`
	Effect     string   `json:"effect"`
	Actions    []string `json:"actions"`
	Resources  []string `json:"resources"`
	Conditions []string `json:"conditions"`
}

type Policy struct {
	Id         string      `json:"id"`
	Statements []Statement `json:"statements"`
}

func (policy *Policy) evaluate(requestedAction, requestedResource string) bool {
	decision := false
	for _, statement := range policy.Statements {
		if assertAction(statement.Actions, requestedAction) {
			if assertResource(statement.Resources, requestedResource) {
				if statement.Effect == "denied" {
					return false
				} else {
					decision = true
				}
			}
		}
	}
	return decision
}

func assertAction(policyActions []string, requestedAction string) bool {
	k := strings.Split(requestedAction, SPLIT_SYMBOL)
	i := len(k)
	for _, action := range policyActions {
		checkLength := i
		p := strings.Split(action, SPLIT_SYMBOL)
		if len(p) < i && p[len(p)-1] == ALL_SYMBOL {
			checkLength = len(p)
		}
		var matches int
		for j := 0; j < checkLength; j++ {
			if k[j] != p[j] && p[j] != ALL_SYMBOL {
				break
			} else if p[j] == ALL_SYMBOL || k[j] == p[j] {
				matches++
				continue
			}
		}
		if matches > 0 && matches == checkLength {
			fmt.Printf("requested action: %s matched with policy action: %s \n", requestedAction, action)
			return true
		}
	}
	return false
}

func assertResource(policyResources []string, requestedResource string) bool {
	k := strings.Split(requestedResource, SPLIT_SYMBOL)
	i := len(k)
	for _, action := range policyResources {
		checkLength := i
		p := strings.Split(action, SPLIT_SYMBOL)
		if len(p) < i && p[len(p)-1] == ALL_SYMBOL {
			checkLength = len(p)
		}
		var matches int
		for j := 0; j < checkLength; j++ {
			if k[j] != p[j] && p[j] != ALL_SYMBOL {
				break
			} else if p[j] == ALL_SYMBOL || k[j] == p[j] {
				matches++
				continue
			}
		}
		if matches > 0 && matches == checkLength {
			fmt.Printf("requested resource: %s matched with policy resource: %s \n", requestedResource, action)
			return true
		}
	}
	return false
}

func (p *Policy) IsValid() bool {
	return false
}

func getUserPolicies(tenantID int, username string) []Policy {
	return nil
}

func isAuthorized2(tenantID int, username string, r *http.Request) bool {
	resourcePrefix := commons.ServerConfigurations.Prefix
	requestedResource := resourcePrefix + SPLIT_SYMBOL + r.URL.Path
	requestedAction := strings.ToLower(requestedResource + SPLIT_SYMBOL + r.Method)

	userPolicies := getUserPolicies(tenantID, username)
	isAuthorized := false
	for _, userPolicy := range userPolicies {
		if userPolicy.evaluate(requestedAction, requestedResource) {
			isAuthorized = true
		} else {
			isAuthorized = false
			break
		}
	}
	return isAuthorized
}

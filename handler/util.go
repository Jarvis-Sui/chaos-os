package handler

import (
	"fmt"
	"net/url"
)

func checkQueryParams(params url.Values, allowedParams []string) error {
	set := make(map[string]bool)
	for _, v := range allowedParams {
		set[v] = true
	}

	for k := range params {
		if _, ok := set[k]; !ok {
			return fmt.Errorf("parameter %s not allowed. allowed params: %v", k, allowedParams)
		}
	}
	return nil
}

func checkRequiredQueryParams(params url.Values, requiredParams []string) error {
	for _, v := range requiredParams {
		if _, ok := params[v]; !ok {
			return fmt.Errorf("parameter %s not found. required params: %v", v, requiredParams)
		}
	}
	return nil
}

func errorMsg(err error) map[string]string {
	return map[string]string{"msg": err.Error()}
}

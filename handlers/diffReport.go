package handlers

import "hash/fnv"

func Hash(str string) uint {
	alg := fnv.New32a()
	alg.Write([]byte(str))
	return uint(alg.Sum32())
}

func diffReport(data map[string]interface{}) map[string]interface{} {
	//var newVulns map[string]interface{} = make(map[string]interface{})
	//var remVulns map[string]interface{} = make(map[string]interface{})
	return map[string]interface{}{}
}

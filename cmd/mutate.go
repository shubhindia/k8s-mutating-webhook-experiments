package cmd

import (
	"encoding/json"
	"fmt"

	v1beta1 "k8s.io/api/admission/v1beta1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Mutate mutates
func Mutate(body []byte, verbose bool) ([]byte, error) {

	// unmarshal request into AdmissionReview struct
	admissionReview := v1beta1.AdmissionReview{}
	if err := json.Unmarshal(body, &admissionReview); err != nil {
		return nil, fmt.Errorf("unmarshaling request failed with %s", err)
	}

	admissionRequest := admissionReview.Request
	response := v1beta1.AdmissionResponse{}

	var pod *v1.Pod
	var operations []PatchOperation

	// unmarshal the object into a pod (since we are doing operations on the pod)
	if err := json.Unmarshal(admissionRequest.Object.Raw, &pod); err != nil {
		return nil, fmt.Errorf("unable unmarshal pod json object %v", err)
	}

	// set required response options
	// allow the request
	response.Allowed = true
	// copy over the UID since this is unique and had to be same as request
	response.UID = admissionRequest.UID
	// define patchType as JSONPath
	pT := v1beta1.PatchTypeJSONPatch
	response.PatchType = &pT

	var containers []v1.Container
	// add existing containers to new list since we want to be there after patch
	containers = append(containers, pod.Spec.Containers...)

	// add new container to the list
	newContainer := v1.Container{
		Name:    "test-sidecar",
		Image:   "busybox:stable",
		Command: []string{"sleep", "infinity"},
	}
	containers = append(containers, newContainer)

	operations = append(operations, ReplacePatchOperation("/spec/containers", containers))

	patch, err := json.Marshal(operations)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal patch operations %s", err)
	}

	// finally, set the patch and result
	response.Patch = patch
	response.Result = &metav1.Status{
		Status: "Success",
	}

	admissionReview.Response = &response

	// convert everything back to json to return as response
	responseBody, err := json.Marshal(admissionReview)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal response %s", err)
	}

	return responseBody, nil
}

type PatchOperation struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	From  string      `json:"from"`
	Value interface{} `json:"value,omitempty"`
}

func ReplacePatchOperation(path string, value interface{}) PatchOperation {
	return PatchOperation{
		Op:    replaceOperation,
		Path:  path,
		Value: value,
	}
}

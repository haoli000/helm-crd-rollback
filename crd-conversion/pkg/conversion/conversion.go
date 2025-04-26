package conversion

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// ConversionReview represents the conversion webhook request/response
type ConversionReview struct {
	apiextensionsv1.ConversionReview
}

// HandleConversion handles the conversion webhook requests
func HandleConversion(w http.ResponseWriter, r *http.Request) {
	var body []byte
	if r.Body != nil {
		if data, err := ioutil.ReadAll(r.Body); err == nil {
			body = data
		}
	}

	if len(body) == 0 {
		log.Println("Empty body")
		http.Error(w, "empty body", http.StatusBadRequest)
		return
	}

	// Decode the request
	var review ConversionReview
	if err := json.Unmarshal(body, &review); err != nil {
		log.Printf("Failed to decode request: %v", err)
		http.Error(w, "could not decode request", http.StatusBadRequest)
		return
	}

	// Set response UID from request
	review.Response = &apiextensionsv1.ConversionResponse{
		UID:              review.Request.UID,
		ConvertedObjects: make([]runtime.RawExtension, len(review.Request.Objects)),
		Result:           metav1.Status{Status: "Success"},
	}

	// Process each object based on desired API version
	for i, obj := range review.Request.Objects {
		var converted runtime.RawExtension

		if review.Request.DesiredAPIVersion == "sample.operator.javaoperatorsdk.io/v1alpha2" {
			converted = convertToV1Alpha2(obj)
		} else if review.Request.DesiredAPIVersion == "sample.operator.javaoperatorsdk.io/v1alpha1" {
			converted = convertToV1Alpha1(obj)
		} else {
			review.Response.Result = metav1.Status{
				Status:  "Failure",
				Message: "Unsupported conversion",
			}
			break
		}

		review.Response.ConvertedObjects[i] = converted
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(review)
}

// convertToV1Alpha2 converts an object to v1alpha2 version
func convertToV1Alpha2(obj runtime.RawExtension) runtime.RawExtension {
	var objMap map[string]interface{}
	json.Unmarshal(obj.Raw, &objMap)

	// Set the correct apiVersion
	objMap["apiVersion"] = "sample.operator.javaoperatorsdk.io/v1alpha2"

	// Add the addedField to status if it exists
	if status, ok := objMap["status"].(map[string]interface{}); ok {
		status["addedField"] = "XXXXXX"
	}

	raw, _ := json.Marshal(objMap)
	return runtime.RawExtension{Raw: raw}
}

// convertToV1Alpha1 converts an object to v1alpha1 version
func convertToV1Alpha1(obj runtime.RawExtension) runtime.RawExtension {
	var objMap map[string]interface{}
	json.Unmarshal(obj.Raw, &objMap)

	// Set the correct apiVersion
	objMap["apiVersion"] = "sample.operator.javaoperatorsdk.io/v1alpha1"

	// Remove the addedField from status if it exists
	if status, ok := objMap["status"].(map[string]interface{}); ok {
		delete(status, "addedField")
	}

	raw, _ := json.Marshal(objMap)
	return runtime.RawExtension{Raw: raw}
}

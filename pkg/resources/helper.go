//
// Copyright 2020 IBM Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package resources

import (
	"math/rand"
	"reflect"
	"time"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// cannot set to const due to k8s struct needing pointers to primitive types

var TrueVar = true
var FalseVar = false

var DefaultSecretMode int32 = 420
var Seconds60 int64 = 60

// Important product values needed for annotations
const LicensingProductName = "IBM Cloud Platform Common Services"
const LicensingProductID = "068a62892a1e4db39641342e592daa25"
const LicensingProductMetric = "FREE"
const LicensingProductVersion = "3.4.0"

const randStringCharset string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const randStringCharsetLength = len(randStringCharset)

type ResourceObject interface {
	metav1.Object
	runtime.Object
}

func RandString(length int) string {
	randFunc := rand.New(rand.NewSource(time.Now().UnixNano())) //#nosec
	outputStringByte := make([]byte, length)
	for i := 0; i < length; i++ {
		outputStringByte[i] = randStringCharset[randFunc.Intn(randStringCharsetLength)]
	}
	return string(outputStringByte)
}

func Contains(s []corev1.LocalObjectReference, e corev1.LocalObjectReference) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func AnnotationsForPod() map[string]string {
	return map[string]string{"productName": LicensingProductName,
		"productID": LicensingProductID, "productVersion": LicensingProductVersion, "productMetric": LicensingProductMetric,
		"clusterhealth.ibm.com/dependencies": "metering"}
}

func WatchForResources(log logr.Logger, o runtime.Object, c controller.Controller, watchTypes []ResourceObject) error {
	for _, restype := range watchTypes {
		log.Info("Watching", "restype", reflect.TypeOf(restype).String())
		err := c.Watch(&source.Kind{Type: restype}, &handler.EnqueueRequestForOwner{
			IsController: true,
			OwnerType:    o,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

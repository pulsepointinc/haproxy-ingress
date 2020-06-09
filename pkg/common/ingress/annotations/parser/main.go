/*
Copyright 2015 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package parser

import (
	"k8s.io/api/core/v1"
	"strconv"

	"github.com/jcmoraisjr/haproxy-ingress/pkg/common/ingress/errors"
	extensions "k8s.io/api/extensions/v1beta1"
)

// IngressAnnotation has a method to parse annotations located in Ingress
type IngressAnnotation interface {
	Parse(ing *extensions.Ingress) (interface{}, error)
}

type metaAnnotations map[string]string

func (a metaAnnotations) parseBool(name string) (bool, error) {
	val, ok := a[name]
	if ok {
		b, err := strconv.ParseBool(val)
		if err != nil {
			return false, errors.NewInvalidAnnotationContent(name, val)
		}
		return b, nil
	}
	return false, errors.ErrMissingAnnotations
}

func (a metaAnnotations) parseString(name string) (string, error) {
	val, ok := a[name]
	if ok {
		return val, nil
	}
	return "", errors.ErrMissingAnnotations
}

func (a metaAnnotations) parseInt(name string) (int, error) {
	val, ok := a[name]
	if ok {
		i, err := strconv.Atoi(val)
		if err != nil {
			return 0, errors.NewInvalidAnnotationContent(name, val)
		}
		return i, nil
	}
	return 0, errors.ErrMissingAnnotations
}

func checkAnnotation(name string, ing *extensions.Ingress) error {
	if ing == nil || len(ing.GetAnnotations()) == 0 {
		return errors.ErrMissingAnnotations
	}
	if name == "" {
		return errors.ErrInvalidAnnotationName
	}

	return nil
}

// GetBoolAnnotation extracts a boolean from an Ingress annotation
func GetBoolAnnotation(name string, ing *extensions.Ingress) (bool, error) {
	err := checkAnnotation(name, ing)
	if err != nil {
		return false, err
	}
	return metaAnnotations(ing.GetAnnotations()).parseBool(name)
}

// GetStringAnnotation extracts a string from an Ingress annotation
func GetStringAnnotation(name string, ing *extensions.Ingress) (string, error) {
	err := checkAnnotation(name, ing)
	if err != nil {
		return "", err
	}
	return metaAnnotations(ing.GetAnnotations()).parseString(name)
}

// GetIntAnnotation extracts an int from an Ingress annotation
func GetIntAnnotation(name string, ing *extensions.Ingress) (int, error) {
	err := checkAnnotation(name, ing)
	if err != nil {
		return 0, err
	}
	return metaAnnotations(ing.GetAnnotations()).parseInt(name)
}

// GetIntNodeAnnotation extracts an int from an Ingress annotation
func GetIntNodeAnnotation(name string, node *v1.Node) (int, error) {
	if node == nil || len(node.GetAnnotations()) == 0 {
		return 0, errors.ErrMissingAnnotations
	}
	if name == "" {
		return 0, errors.ErrInvalidAnnotationName
	}
	return metaAnnotations(node.GetAnnotations()).parseInt(name)
}

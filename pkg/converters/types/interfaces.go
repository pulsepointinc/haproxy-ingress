/*
Copyright 2019 The HAProxy Ingress Controller Authors.

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

package types

import (
	"net"
	"time"

	api "k8s.io/api/core/v1"
	discoveryv1 "k8s.io/api/discovery/v1"
	networking "k8s.io/api/networking/v1"
	gatewayv1alpha1 "sigs.k8s.io/gateway-api/apis/v1alpha1"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"

	hatypes "github.com/jcmoraisjr/haproxy-ingress/pkg/haproxy/types"
)

// Cache ...
type Cache interface {
	ExternalNameLookup(externalName string) ([]net.IP, error)
	GetIngress(ingressName string) (*networking.Ingress, error)
	GetIngressList() ([]*networking.Ingress, error)
	GetIngressClass(className string) (*networking.IngressClass, error)
	GetGatewayA1(gatewayName string) (*gatewayv1alpha1.Gateway, error)
	GetGatewayA1List() ([]*gatewayv1alpha1.Gateway, error)
	GetHTTPRouteA1List(namespace string, match map[string]string) ([]*gatewayv1alpha1.HTTPRoute, error)
	GetGatewayMap() (map[string]*gatewayv1alpha2.Gateway, error)
	GetHTTPRouteList() ([]*gatewayv1alpha2.HTTPRoute, error)
	GetService(defaultNamespace, serviceName string) (*api.Service, error)
	GetEndpoints(service *api.Service) (*api.Endpoints, error)
	GetEndpointSlices(service *api.Service) ([]*discoveryv1.EndpointSlice, error)
	GetConfigMap(configMapName string) (*api.ConfigMap, error)
	GetNamespace(name string) (*api.Namespace, error)
	GetTerminatingPods(service *api.Service, track []TrackingRef) ([]*api.Pod, error)
	GetPod(podName string) (*api.Pod, error)
	GetPodNamespace() string
	GetTLSSecretPath(defaultNamespace, secretName string, track []TrackingRef) (CrtFile, error)
	GetCASecretPath(defaultNamespace, secretName string, track []TrackingRef) (ca, crl File, err error)
	GetDHSecretPath(defaultNamespace, secretName string) (File, error)
	GetPasswdSecretContent(defaultNamespace, secretName string, track []TrackingRef) ([]byte, error)
	SwapChangedObjects() *ChangedObjects
	GetNodeByName(nodeName string) (*api.Node, error)
}

// ChangedObjects ...
type ChangedObjects struct {
	//
	GlobalConfigMapDataCur, GlobalConfigMapDataNew map[string]string
	//
	TCPConfigMapDataCur, TCPConfigMapDataNew map[string]string
	//
	IngressesDel, IngressesUpd, IngressesAdd []*networking.Ingress
	//
	IngressClassesDel, IngressClassesUpd, IngressClassesAdd []*networking.IngressClass
	//
	GatewaysA1Del, GatewaysA1Upd, GatewaysA1Add []*gatewayv1alpha1.Gateway
	//
	GatewayClassesA1Del, GatewayClassesA1Upd, GatewayClassesA1Add []*gatewayv1alpha1.GatewayClass
	//
	HTTPRoutesA1Del, HTTPRoutesA1Upd, HTTPRoutesA1Add []*gatewayv1alpha1.HTTPRoute
	//
	GatewaysDel, GatewaysUpd, GatewaysAdd []*gatewayv1alpha2.Gateway
	//
	GatewayClassesDel, GatewayClassesUpd, GatewayClassesAdd []*gatewayv1alpha2.GatewayClass
	//
	HTTPRoutesDel, HTTPRoutesUpd, HTTPRoutesAdd []*gatewayv1alpha2.HTTPRoute
	//
	EndpointsNew []*api.Endpoints
	//
	EndpointSlicesUpd []*discoveryv1.EndpointSlice
	//
	ServicesDel, ServicesUpd, ServicesAdd []*api.Service
	//
	SecretsDel, SecretsUpd, SecretsAdd []*api.Secret
	//
	ConfigMapsDel, ConfigMapsUpd, ConfigMapsAdd []*api.ConfigMap
	//
	PodsNew []*api.Pod
	//
	NeedFullSync bool
	//
	NodesUpd []*api.Node
	//
	Objects []string
	Links   TrackingLinks
}

// ResourceType ...
type ResourceType string

// ...
const (
	ResourceIngress      ResourceType = "Ingress"
	ResourceIngressClass ResourceType = "IngressClass"

	ResourceGatewayA1      ResourceType = "GatewayA1"
	ResourceGatewayClassA1 ResourceType = "GatewayClassA1"
	ResourceHTTPRouteA1    ResourceType = "HTTPRouteA1"

	ResourceGateway      ResourceType = "Gateway"
	ResourceGatewayClass ResourceType = "GatewayClass"
	ResourceHTTPRoute    ResourceType = "HTTPRoute"

	ResourceConfigMap ResourceType = "ConfigMap"
	ResourceService   ResourceType = "Service"
	ResourceEndpoints ResourceType = "Endpoints"
	ResourceSecret    ResourceType = "Secret"
	ResourcePod       ResourceType = "Pod"
	ResourceNode      ResourceType = "Node"

	ResourceHATCPService ResourceType = "HATCPService"
	ResourceHAHostname   ResourceType = "HAHostname"
	ResourceHABackend    ResourceType = "HABackend"
	ResourceHAUserlist   ResourceType = "HAUserlist"

	ResourceAcmeData ResourceType = "AcmeData"
)

// TrackingRef ...
type TrackingRef struct {
	Context    ResourceType
	UniqueName string
}

// TrackingLinks ...
type TrackingLinks map[ResourceType][]string

// Tracker ...
type Tracker interface {
	TrackNames(leftContext ResourceType, leftName string, rightContext ResourceType, rightName string)
	TrackRefName(left []TrackingRef, rightContext ResourceType, rightName string)
	TrackRefs(left, right TrackingRef)
	QueryLinks(input TrackingLinks, removeMatches bool) TrackingLinks
	ClearLinks()
}

// AnnotationReader ...
type AnnotationReader interface {
	ReadAnnotations(backend *hatypes.Backend, services []*api.Service, pathLinks []hatypes.PathLink)
}

// File ...
type File struct {
	Filename string
	SHA1Hash string
}

// CrtFile ...
type CrtFile struct {
	Filename   string
	SHA1Hash   string
	CommonName string
	NotAfter   time.Time
}

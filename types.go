// Copyright 2016 Google Inc. All Rights Reserved.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

type InstanceGroup struct {
	Autoscaler       Autoscaler `yaml:"autoscaler"`
	InstanceTemplate string     `yaml:"instanceTemplate"`
	TargetPools      []string   `yaml:"targetPools"`
}

type Autoscaler struct {
	AutoscalingPolicy AutoscalingPolicy `yaml:"autoscalingPolicy"`
}

type AutoscalingPolicy struct {
	CPUUtilization CPUUtilization `yaml:"cpuUtilization"`
	MaxNumReplicas int64          `yaml:"maxNumReplicas"`
	MinNumReplicas int64          `yaml:"minNumReplicas"`
}

type CPUUtilization struct {
	UtilizationTarget float64 `yaml:"utilizationTarget"`
}

type InstanceTemplate struct {
	Properties Properties `yaml:"properties"`
}

type Properties struct {
	Metadata PropertiesMetadata `yaml:"metadata"`
}

type PropertiesMetadata struct {
	Items []Item `yaml:"items"`
}

type Item struct {
	Key   string `yaml:"key"`
	Value string `yaml:"value"`
}

type CloudInit struct {
	WriteFiles []Path `yaml:"write_files"`
}

type Path struct {
	Path    string `yaml:"path"`
	Content string `yaml:"content"`
}

type Deployment struct {
	ApiVersion string         `json:"apiVersion,omitempty" yaml:"apiVersion,omitempty"`
	Kind       string         `json:"kind,omitempty" yaml:"kind,omitempty"`
	Metadata   Metadata       `json:"metadata" yaml:"metadata"`
	Spec       DeploymentSpec `json:"spec" yaml:"spec"`
}

type Service struct {
	ApiVersion string      `json:"apiVersion,omitempty" yaml:"apiVersion,omitempty"`
	Kind       string      `json:"kind,omitempty" yaml:"kind,omitempty"`
	Metadata   Metadata    `json:"metadata" yaml:"metadata"`
	Spec       ServiceSpec `json:"spec" yaml:"spec"`
}

type Metadata struct {
	Name            string            `json:"name,omitempty" yaml:"name,omitempty"`
	GenerateName    string            `json:"generateName,omitempty" yaml:"generateName,omitempty"`
	ResourceVersion string            `json:"resourceVersion,omitempty" yaml:"resourceVersion,omitempty"`
	Labels          map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
	Annotations     map[string]string `json:"annotations,omitempty" yaml:"annotations,omitempty"`
	Uid             string            `json:"uid,omitempty" yaml:"uid,omitempty"`
}

type DeploymentSpec struct {
	Replicas int64       `json:"replicas"`
	Template PodTemplate `json:"template"`
}

type ServiceSpec struct {
	ClusterIP string            `json:"clusterIP,omitempty" yaml:"clusterIP,omitempty"`
	Type      string            `json:"type"`
	Ports     []Port            `json:"ports"`
	Selector  map[string]string `yaml:"selector"`
}

type Port struct {
	Name       string `json:"name" yaml:"name"`
	Protocol   string `json:"protocol" yaml:"protocol"`
	Port       int64  `json:"port" yaml:"port"`
	TargetPort int64  `json:"targetPort" yaml:"targetPort"`
}

type PodTemplate struct {
	Metadata Metadata `json:"metadata" yaml:"metadata"`
	Spec     PodSpec  `json:"spec" yaml:"spec"`
}

type PodSpec struct {
	Containers []Container `json:"containers" yaml:"containers"`
}

type Container struct {
	Image string `json:"image" yaml:"image"`
	Name  string `json:"name" yaml:"name"`
}

type HorizontalPodAutoscaler struct {
	ApiVersion string                      `json:"apiVersion,omitempty" yaml:"apiVersion,omitempty"`
	Kind       string                      `json:"kind,omitempty" yaml:"kind,omitempty"`
	Metadata   Metadata                    `json:"metadata" yaml:"metadata"`
	Spec       HorizontalPodAutoscalerSpec `json:"spec" yaml:"spec"`
}

type HorizontalPodAutoscalerSpec struct {
	MinReplicas    int64                `json:"minReplicas" yaml:"minReplicas"`
	MaxReplicas    int64                `json:"maxReplicas" yaml:"maxReplicas"`
	CPUUtilization CPUTargetUtilization `json:"cpuUtilization" yaml:"cpuUtilization"`
	ScaleRef       SubresourceReference `json:"scaleRef" yaml:"scaleRef"`
}

type CPUTargetUtilization struct {
	TargetPercentage int64 `json:"targetPercentage" yaml:"targetPercentage"`
}

type SubresourceReference struct {
	Kind        string `json:"kind" yaml:"kind"`
	Name        string `json:"name" yaml:"name"`
	Subresource string `json:"subresource" yaml:"subresource"`
}

type ForwardingRule struct {
	ID string `json:"id" yaml:"id"`
}

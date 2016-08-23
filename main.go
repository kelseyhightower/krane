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

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"gopkg.in/yaml.v2"
)

var (
	instanceGroup string
)

func main() {
	flag.StringVar(&instanceGroup, "instance-group", "", "Instance group")
	flag.Parse()

	// Get the instance group details.
	cmd := exec.Command("gcloud", "compute", "instance-groups", "managed", "describe", instanceGroup)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	var ig InstanceGroup
	err = yaml.Unmarshal(out.Bytes(), &ig)
	if err != nil {
		log.Fatal(err)
	}

	// Extract the pod manifest from the instance template.
	out.Reset()
	cmd = exec.Command("gcloud", "compute", "instance-templates", "describe", ig.InstanceTemplate)
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	var it InstanceTemplate
	err = yaml.Unmarshal(out.Bytes(), &it)
	if err != nil {
		log.Fatal(err)
	}

	// We assume instances are created with the Google Container Image and
	// a cloud-init cloud config.
	var ci CloudInit
	for _, i := range it.Properties.Metadata.Items {
		if i.Key == "user-data" {
			err := yaml.Unmarshal([]byte(i.Value), &ci)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	for _, pod := range ci.WriteFiles {
		var pt PodTemplate
		err := yaml.Unmarshal([]byte(pod.Content), &pt)

		labels := make(map[string]string)
		labels["name"] = pt.Metadata.Name
		labels["app"] = pt.Metadata.Name

		deployment := Deployment{
			ApiVersion: "extensions/v1beta1",
			Kind:       "Deployment",
			Metadata:   Metadata{Name: pt.Metadata.Name},
			Spec: DeploymentSpec{
				Replicas: 1,
				Template: PodTemplate{
					Metadata: Metadata{
						Labels: labels,
					},
					Spec: pt.Spec,
				},
			},
		}

		// Print a Kubernetes Deployment config to stdout.
		deploymentData, err := yaml.Marshal(&deployment)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(deploymentData))
		fmt.Println("---")

		// Create HPA
		targetPercentage := int64(ig.Autoscaler.AutoscalingPolicy.CPUUtilization.UtilizationTarget * 100)
		if targetPercentage == 0 {
			targetPercentage = 1
		}

		hpa := HorizontalPodAutoscaler{
			ApiVersion: "extensions/v1beta1",
			Kind:       "HorizontalPodAutoscaler",
			Metadata:   Metadata{Name: pt.Metadata.Name},
			Spec: HorizontalPodAutoscalerSpec{
				MinReplicas: ig.Autoscaler.AutoscalingPolicy.MinNumReplicas,
				MaxReplicas: ig.Autoscaler.AutoscalingPolicy.MaxNumReplicas,
				CPUUtilization: CPUTargetUtilization{
					TargetPercentage: targetPercentage,
				},
				ScaleRef: SubresourceReference{
					Kind:        "Deployment",
					Name:        pt.Metadata.Name,
					Subresource: "scale",
				},
			},
		}

		// Print a Kubernetes HorizontalPodAutoscaler config to stdout.
		hpaData, err := yaml.Marshal(&hpa)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(hpaData))

		// Determine if the instance group belongs to a target pool that is
		// referenced by a Google Cloud Platform network load balancer.
		// If so, create Kubernetes Service with type LoadBalancer.
		filter := fmt.Sprintf("target=%s", strings.SplitAfterN(ig.TargetPools[0], "/", 9)[8])
		cmd = exec.Command("gcloud", "compute", "forwarding-rules", "list", "--filter", filter, "--format", "yaml")

		out.Reset()
		cmd.Stdout = &out
		err = cmd.Run()
		if err != nil {
			log.Fatal(err)
		}

		var fr ForwardingRule
		err = yaml.Unmarshal(out.Bytes(), &fr)
		if err != nil {
			log.Fatal(err)
		}

		if fr.ID != "" {
			fmt.Println("---")
			service := Service{
				ApiVersion: "v1",
				Kind:       "Service",
				Metadata:   Metadata{Name: pt.Metadata.Name},
				Spec: ServiceSpec{
					Type: "LoadBalancer",
					Ports: []Port{
						Port{
							Name:       "http",
							Protocol:   "TCP",
							Port:       80,
							TargetPort: 80,
						},
					},
					Selector: map[string]string{
						"app": pt.Metadata.Name,
					},
				},
			}

			// Print a Kubernetes Service config to stdout.
			svcData, err := yaml.Marshal(&service)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(string(svcData))
		}
	}
}

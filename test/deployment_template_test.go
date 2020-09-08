package test

import (
	"testing"
	appsv1 "k8s.io/api/apps/v1"  
	"github.com/gruntwork-io/terratest/modules/helm"
)

func TestDeploymentTemplateRendersContainerImage(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := "../"

	// Setup the args. For this test, we will set the following input values:
	// - crccheck/hello-world
	options := &helm.Options{
		SetValues: map[string]string{"image.repository": "crccheck/hello-world"},
	}

	// Run RenderTemplate to render the template and capture the output.
	output := helm.RenderTemplate(t, options, helmChartPath, "deployment", []string{"templates/deployment.yaml"})

	// Now we use kubernetes/client-go library to render the template output into the Deployment struct. This will
	// ensure the Deployment resource is rendered correctly.
	var deployment appsv1.Deployment
	helm.UnmarshalK8SYaml(t, output, &deployment)

	// Finally, we verify the pod spec is set to the expected container image value
	expectedContainerImage := "crccheck/hello-world:latest"
	podContainers := deployment.Spec.Template.Spec.Containers
	if podContainers[0].Image != expectedContainerImage {
		t.Fatalf("Rendered container image (%s) is not expected (%s)", podContainers[0].Image, expectedContainerImage)
	}
}
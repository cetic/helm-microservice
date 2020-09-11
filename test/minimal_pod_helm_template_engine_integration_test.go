package test

import (
	"fmt"
	"strings"
	"testing"
	"time"
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/random"
	http_helper "github.com/gruntwork-io/terratest/modules/http-helper"
)

func TestPodDeploysContainerImageHelmTemplateEngine(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := "../"

	// Setup the kubectl config and context. Here we choose to use the defaults, which is:
	// - HOME/.kube/config for the kubectl config file
	// - Current context of the kubectl config file
	// We also specify that we are working in the default namespace (required to get the Pod)
	kubectlOptions := k8s.NewKubectlOptions("", "", "default")

	// Setup the args. For this test, we will set the following input values:
	// - image=nginx:1.15.8
	// - fullnameOverride=minimal-pod-RANDOM_STRING
	// We use a fullnameOverride so we can find the Pod later during verification
	podName := fmt.Sprintf("minimal-pod-%s", strings.ToLower(random.UniqueId()))
	options := &helm.Options{
		SetValues: map[string]string{"fullnameOverride": podName, "nameOverride": podName},
	}

	// Run RenderTemplate to render the template and capture the output.
	output := helm.RenderTemplate(t, options, helmChartPath, "minimal-pod", []string{})

	// Make sure to delete the resources at the end of the test
	defer k8s.KubectlDeleteFromString(t, kubectlOptions, output)

	// Now use kubectl to apply the rendered template
	k8s.KubectlApplyFromString(t, kubectlOptions, output)

	// Now that the chart is deployed, verify the deployment. This function will open a tunnel to the Pod and hit the
	// nginx container endpoint.
	verifyNginxPod(t, kubectlOptions, podName)
}

// verifyNginxPod will open a tunnel to the Pod and hit the endpoint to verify the nginx welcome page is shown.
func verifyNginxPod(t *testing.T, kubectlOptions *k8s.KubectlOptions, podName string) {
	// Wait for the pod to come up. It takes some time for the Pod to start, so retry a few times.
	retries := 15
	sleep := 5 * time.Second
	k8s.WaitUntilPodAvailable(t, kubectlOptions, podName, retries, sleep)

	// We will first open a tunnel to the pod, making sure to close it at the end of the test.
	tunnel := k8s.NewTunnel(kubectlOptions, k8s.ResourceTypePod, podName, 0, 80)
	defer tunnel.Close()
	tunnel.ForwardPort(t)

	// ... and now that we have the tunnel, we will verify that we get back a 200 OK with the nginx welcome page.
	// It takes some time for the Pod to start, so retry a few times.
	endpoint := fmt.Sprintf("http://%s", tunnel.Endpoint())
	http_helper.HttpGetWithRetryWithCustomValidation(
		t,
		endpoint,
		nil,
		retries,
		sleep,
		func(statusCode int, body string) bool {
			return statusCode == 200 && strings.Contains(body, "Hello World")
		},
	)
}

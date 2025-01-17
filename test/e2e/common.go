package e2e

import (
	"fmt"
	"os/exec"
	"time"

	"golang.org/x/net/context"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"

	"github.com/vmware-tanzu/velero/pkg/builder"
)

// EnsureClusterExists returns whether or not a kubernetes cluster exists for tests to be run on.
func EnsureClusterExists(ctx context.Context) error {
	return exec.CommandContext(ctx, "kubectl", "cluster-info").Run()
}

// CreateNamespace creates a kubernetes namespace
func CreateNamespace(ctx context.Context, client *kubernetes.Clientset, namespace string) error {
	ns := builder.ForNamespace(namespace).Result()
	_, err := client.CoreV1().Namespaces().Create(ctx, ns, metav1.CreateOptions{})
	if apierrors.IsAlreadyExists(err) {
		return nil
	}
	return err
}

// WaitForNamespaceDeletion Waits for namespace to be deleted.
func WaitForNamespaceDeletion(interval, timeout time.Duration, client *kubernetes.Clientset, ns string) error {
	err := wait.PollImmediate(interval, timeout, func() (bool, error) {
		_, err := client.CoreV1().Namespaces().Get(context.TODO(), ns, metav1.GetOptions{})
		if err != nil {
			if apierrors.IsNotFound(err) {
				return true, nil
			}
			return false, err
		}
		fmt.Printf("Namespace %s is still being deleted...\n", ns)
		return false, nil
	})
	return err
}

package test

import (
	"flag"
	"github.com/appscode/go/crypto/rand"
	"k8s.io/client-go/util/homedir"
	"testing"
	"time"

	"e2e_test/framework"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/reporters"
	. "github.com/onsi/gomega"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
)

var (
	StorageClass = "linode-block-storage"
	useExisting  = false
	dlt          = false
	ClusterName  = rand.WithUniqSuffix("csi-linode")
)

func init() {
	flag.StringVar(&framework.Image, "image", framework.Image, "registry/repository:tag")
	flag.StringVar(&framework.ApiToken, "api-token", os.Getenv("LINODE_API_TOKEN"), "linode api token")
	flag.BoolVar(&dlt, "delete", dlt, "Delete cluster after test")
	flag.BoolVar(&useExisting, "use-existing", useExisting, "Use existing kubernetes cluster")
	flag.StringVar(&framework.KubeConfigFile, "kubeconfig", filepath.Join(homedir.HomeDir(), ".kube/config"), "To use existing cluster provide kubeconfig file")
}

const (
	TIMEOUT = 20 * time.Minute
)

var (
	root *framework.Framework
)

func TestE2e(t *testing.T) {
	RegisterFailHandler(Fail)
	SetDefaultEventuallyTimeout(TIMEOUT)

	junitReporter := reporters.NewJUnitReporter("junit.xml")
	RunSpecsWithDefaultAndCustomReporters(t, "e2e Suite", []Reporter{junitReporter})
}

var _ = BeforeSuite(func() {

	if !useExisting {
		err := framework.CreateCluster(ClusterName)
		Expect(err).NotTo(HaveOccurred())
		dir, err := os.Getwd()
		Expect(err).NotTo(HaveOccurred())
		framework.KubeConfigFile = filepath.Join(dir, ClusterName+".conf")
	}

	By("Using kubeconfig from " + framework.KubeConfigFile)
	config, err := clientcmd.BuildConfigFromFlags("", framework.KubeConfigFile)
	Expect(err).NotTo(HaveOccurred())

	// Clients
	kubeClient := kubernetes.NewForConfigOrDie(config)

	// Framework
	root = framework.New(config, kubeClient, StorageClass)

	By("Using Namespace " + root.Namespace())
	err = root.CreateNamespace()
	Expect(err).NotTo(HaveOccurred())

	By("Applying Manifest ")
	err = root.ApplyManifest()
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	if dlt || !useExisting {
		err := framework.DeleteCluster(ClusterName)
		Expect(err).NotTo(HaveOccurred())
	}
})

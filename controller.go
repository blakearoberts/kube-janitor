package main

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	batchV1 "k8s.io/client-go/informers/batch/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"

	log "github.com/sirupsen/logrus"
)

// Controller is the controller implementation for Job resources
type Controller struct {
	informer batchV1.JobInformer
}

// New returns a new Job controller
func New(kubeClient kubernetes.Interface) (*Controller, error) {
	informer := informers.NewSharedInformerFactory(kubeClient, 0).Batch().V1().Jobs()

	c := &Controller{
		informer: informer,
	}

	informer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    c.add,
		UpdateFunc: c.update,
		DeleteFunc: c.delete,
	})

	return c, nil
}

// Run starts the controller. Run is blocking and therefore expected to be
// called from a dedicated goroutine
func (c *Controller) Run() error {
	stopper := make(chan struct{})
	defer close(stopper)
	c.informer.Informer().Run(stopper)
	return nil
}

func (c *Controller) add(obj interface{}) {
	log.Infof("job added: %s", obj.(metav1.Object).GetName())
}

func (c *Controller) update(old interface{}, new interface{}) {
	obj := new.(metav1.Object)

	job, err := c.informer.Lister().Jobs(obj.GetNamespace()).Get(obj.GetName())
	if err != nil {
		log.Warningf("failed to get job during update event: %v", err)
		return
	}
	log.Infof("job '%s' updated with status: %v", obj.GetName(), job.Status)
}

func (c *Controller) delete(obj interface{}) {
	log.Infof("job deleted: %s", obj.(metav1.Object).GetName())
}

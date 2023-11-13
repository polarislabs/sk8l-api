package main

import (
	"context"
	"fmt"
	"log"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	// "k8s.io/apimachinery/pkg/util/yaml"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	//
	// Uncomment to load all auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth"
	//
	// Or uncomment to load specific auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth/azure"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
)

type K8sClient struct {
	namespace string
	ClientSet *kubernetes.Clientset
}

// A ClientOption is used to configure a Client.
type ClientOption func(*K8sClient)

func WithNamespace(namespace string) ClientOption {
	return func(kc *K8sClient) {
		kc.namespace = namespace
	}
}

func NewK8sClient(options ...ClientOption) *K8sClient {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	k8sClient := &K8sClient{
		ClientSet: clientset,
	}

	for _, option := range options {
		option(k8sClient)
	}

	return k8sClient
}

func (kc *K8sClient) GetCronjobs() *batchv1.CronJobList {
	ctx := context.TODO()

	cronJobs, err := kc.ClientSet.BatchV1().CronJobs(kc.namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	log.Printf("There are %d cronjobs in the cluster\n", len(cronJobs.Items))

	return cronJobs
}

func (kc *K8sClient) GetCronjob(cronjobNamespace, cronjobName string) *batchv1.CronJob {
	ctx := context.TODO()
	cronJob, err := kc.ClientSet.BatchV1().CronJobs(cronjobNamespace).Get(ctx, cronjobName, metav1.GetOptions{})

	if errors.IsNotFound(err) {
		log.Printf("Cronjob %s found in default namespace\n", cronjobName)
	} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		log.Printf("Error getting cronjob %v\n", statusError.ErrStatus.Message)
	} else if err != nil {
		panic(err.Error())
	} else {
		log.Printf("Found example-xxxxx cronjob in default namespace\n")
	}

	return cronJob
}

func (kc *K8sClient) GetJobPodsForJob(job *batchv1.Job) *corev1.PodList {
	ctx := context.TODO()

	fmt.Println("searching for", job.Namespace, job.Name)
	jobPods, err := kc.ClientSet.CoreV1().Pods(job.Namespace).List(ctx, metav1.ListOptions{
		LabelSelector: fmt.Sprintf("job-name=%s", job.Name),
	})

	if errors.IsNotFound(err) {
		log.Printf("Pods found in default namespace\n")
	} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		log.Printf("Error getting pod %v\n", statusError.ErrStatus.Message)
	} else if err != nil {
		panic(err.Error())
	} else {
		log.Printf("There are %d jobPods in the cluster\n", len(jobPods.Items))
	}

	return jobPods
}

func (kc *K8sClient) GetPod(jobNamespace, podName string) *corev1.Pod {
	ctx := context.TODO()
	pod, err := kc.ClientSet.CoreV1().Pods(jobNamespace).Get(ctx, podName, metav1.GetOptions{})

	// Examples for error handling:
	// - Use helper functions e.g. errors.IsNotFound()
	// - And/or cast to StatusError and use its properties like e.g. ErrStatus.Message
	if errors.IsNotFound(err) {
		log.Printf("Pod example-xxxxx not found in default namespace\n")
	} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		log.Printf("Error getting pod %v\n", statusError.ErrStatus.Message)
	} else if err != nil {
		panic(err.Error())
	} else {
		log.Printf("Found example-xxxxx pod in default namespace\n")
	}

	return pod
}

func (kc *K8sClient) GetJob(jobNamespace, jobName string) *batchv1.Job {
	ctx := context.TODO()
	job, err := kc.ClientSet.BatchV1().Jobs(jobNamespace).Get(ctx, jobName, metav1.GetOptions{})
	if err != nil {
		panic(err.Error())
	}

	return job
}

func (kc *K8sClient) GetAllJobs() *batchv1.JobList {
	ctx := context.TODO()

	// get pods in all the namespaces by omitting namespace
	// Or specify namespace to get pods in particular namespace
	// Limit: 10, // need to fix this - last duration / current duration get messed up
	jobs, err := kc.ClientSet.BatchV1().Jobs(kc.namespace).List(ctx, metav1.ListOptions{})

	if err != nil {
		panic(err.Error())
	}
	log.Printf("There are %d jobs in the cluster\n", len(jobs.Items))

	return jobs
}
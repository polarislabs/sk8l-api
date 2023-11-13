package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/danroux/sk8l/protos"
	gyaml "github.com/ghodss/yaml"
	"google.golang.org/grpc/health/grpc_health_v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
)

type Sk8lServer struct {
	grpc_health_v1.UnimplementedHealthServer
	protos.UnimplementedCronjobServer
	K8sClient *K8sClient
}

func (h Sk8lServer) Check(ctx context.Context, req *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	log.Default().Println("serving health")
	return &grpc_health_v1.HealthCheckResponse{Status: grpc_health_v1.HealthCheckResponse_SERVING}, nil
}

func (h Sk8lServer) Watch(req *grpc_health_v1.HealthCheckRequest, stream grpc_health_v1.Health_WatchServer) error {
	response := &grpc_health_v1.HealthCheckResponse{Status: grpc_health_v1.HealthCheckResponse_SERVING}

	if err := stream.Send(response); err != nil {
		return err
	}

	return nil
}

func (s Sk8lServer) GetCronjobs(context.Context, *protos.CronjobsRequest) (*protos.CronjobsResponse, error) {
	cronJobList := s.K8sClient.GetCronjobs()
	// cronJobList := getMocks()
	// mocked := getMocks().Items
	// cronJobList.Items = append(cronJobList.Items, mocked...)

	n := len(cronJobList.Items)
	cronjobs := make([]*protos.CronjobResponse, 0, n)

	for _, cronjobItem := range cronJobList.Items {
		cronjob := s.cronJobResponse(cronjobItem)
		cronjobs = append(cronjobs, cronjob)
	}

	y := &protos.CronjobsResponse{
		Cronjobs: cronjobs,
	}

	return y, nil
}

func (s Sk8lServer) GetCronjob(ctx context.Context, in *protos.CronjobRequest) (*protos.CronjobResponse, error) {
	log.Printf("Received: %v, %v", in.CronjobName, in.CronjobNamespace)

	cronjobName := in.CronjobName
	cronjobNamespace := in.CronjobNamespace

	cronjob := s.K8sClient.GetCronjob(cronjobNamespace, cronjobName)
	cronjobPodsResponse := s.cronJobResponse(*cronjob)

	return cronjobPodsResponse, nil
}

func (s Sk8lServer) GetCronjobPods(ctx context.Context, in *protos.CronjobPodsRequest) (*protos.CronjobPodsResponse, error) {
	cronjobName := in.CronjobName
	cronjobNamespace := in.CronjobNamespace

	cronjob := s.K8sClient.GetCronjob(cronjobNamespace, cronjobName)
	cronjobResponse := s.cronJobResponse(*cronjob)
	lightweightCronjobPodsResponse := &protos.CronjobResponse{
		Name:      cronjob.Name,
		Namespace: cronjob.Namespace,
		Jobs:      cronjobResponse.Jobs,
	}

	cronjobPodsResponse := &protos.CronjobPodsResponse{
		Pods:    cronjobResponse.JobsPods,
		Cronjob: lightweightCronjobPodsResponse,
	}

	return cronjobPodsResponse, nil
}

func (s Sk8lServer) GetCronjobYAML(ctx context.Context, in *protos.CronjobRequest) (*protos.CronjobYAMLResponse, error) {
	cronjob := s.K8sClient.GetCronjob(in.CronjobNamespace, in.CronjobName)
	prettyJson, _ := json.MarshalIndent(cronjob, "", "  ")

	y, _ := gyaml.JSONToYAML(prettyJson)

	x := &protos.CronjobYAMLResponse{
		Cronjob: string(y),
	}

	return x, nil
}

func (s Sk8lServer) GetJobYAML(ctx context.Context, in *protos.JobRequest) (*protos.JobYAMLResponse, error) {
	job := s.K8sClient.GetJob(in.JobNamespace, in.JobName)
	prettyJson, _ := json.MarshalIndent(job, "", "  ")

	y, _ := gyaml.JSONToYAML(prettyJson)

	x := &protos.JobYAMLResponse{
		Job: string(y),
	}

	return x, nil
}

func (s Sk8lServer) GetPodYAML(ctx context.Context, in *protos.PodRequest) (*protos.PodYAMLResponse, error) {
	pod := s.K8sClient.GetPod(in.PodNamespace, in.PodName)
	prettyJson, _ := json.MarshalIndent(pod, "", "  ")

	y, _ := gyaml.JSONToYAML(prettyJson)

	x := &protos.PodYAMLResponse{
		Pod: string(y),
	}

	return x, nil
}

// Revisit this. JobConditions are not being used yet anywhere. PodResponse.TerminationReasons.TerminationDetails -> ContainerStateTerminated
func jobFailed(job batchv1.Job, jobPodsResponses []*protos.PodResponse) (bool, *batchv1.JobCondition, []*batchv1.JobCondition) {
	var jobFailed bool
	var failureCondition *batchv1.JobCondition

	n := len(job.Status.Conditions)
	jobConditions := make([]*batchv1.JobCondition, 0, n)
	for _, jobCondition := range job.Status.Conditions {
		log.Println("\n\njobCondition", jobCondition)
		log.Println("checking Type", jobCondition.Type == batchv1.JobFailed, jobCondition.Type, "\n\n")

		if jobFailed != true {
			if jobCondition.Type == batchv1.JobFailed {
				jobFailed = true
				failureCondition = &jobCondition
			}
		}
		jobConditions = append(jobConditions, &jobCondition)
	}

	for _, pr := range jobPodsResponses {
		if pr.Failed {
			jobFailed = true
		}
	}

	return jobFailed, failureCondition, jobConditions
}

func (s Sk8lServer) buildJobResponse(batchJob batchv1.Job) *protos.JobResponse {
	jobPodsForJob := s.K8sClient.GetJobPodsForJob(&batchJob)
	jobPodsResponses := buildJobPodsResponses(jobPodsForJob)

	jobFailed, failureCondition, jobConditions := jobFailed(batchJob, jobPodsResponses)

	duration := toDuration(batchJob, jobFailed, failureCondition)
	durationInS := toDurationInS(batchJob, jobFailed, failureCondition)
	completionTimeInS := toCompletionTimeInS(batchJob)

	terminationReasons := make([]*protos.TerminationReason, 0)

	for _, podResponse := range jobPodsResponses {
		terminationReasons = append(terminationReasons, podResponse.TerminationReasons...)
	}

	jobResponse := &protos.JobResponse{
		Name:              batchJob.Name,
		Namespace:         batchJob.Namespace,
		Uuid:              string(batchJob.UID),
		CreationTimestamp: batchJob.GetCreationTimestamp().UTC().Format(time.RFC3339),
		Generation:        batchJob.Generation,
		Duration:          duration.String(),
		DurationInS:       durationInS,
		Spec:              &batchJob.Spec,
		Status: &protos.JobStatus{
			StartTime:         batchJob.Status.StartTime,
			StartTimeInS:      batchJob.Status.StartTime.Unix(),
			CompletionTime:    batchJob.Status.CompletionTime,
			CompletionTimeInS: completionTimeInS,
			Active:            &batchJob.Status.Active,
			Failed:            &batchJob.Status.Failed,
			Ready:             batchJob.Status.Ready,
			Succeeded:         &batchJob.Status.Succeeded,
			Conditions:        jobConditions,
		},
		Succeeded:          jobSucceded(batchJob),
		Failed:             jobFailed,
		FailureCondition:   failureCondition,
		Pods:               jobPodsResponses,
		TerminationReasons: terminationReasons,
	}

	return jobResponse
}

func (s Sk8lServer) allAndRunningJobsAnPods(jobs *batchv1.JobList, jobUID types.UID) ([]*protos.JobResponse, []*protos.PodResponse, []*protos.JobResponse, []*protos.PodResponse) {
	jn := len(jobs.Items)
	allJobsForCronJob := make([]*protos.JobResponse, 0, jn)
	allJobPodsForCronjob := make([]*protos.PodResponse, 0)
	runningJobs := make([]*protos.JobResponse, 0, jn)
	runningPods := make([]*protos.PodResponse, 0)

	// go through all jobs and get the ones that match the jobUID(owner)
	for _, batchJob := range jobs.Items {
		for _, owr := range batchJob.ObjectMeta.OwnerReferences {
			if owr.UID == jobUID {
				jobResponse := s.buildJobResponse(batchJob)
				allJobsForCronJob = append(allJobsForCronJob, jobResponse)
				allJobPodsForCronjob = append(allJobPodsForCronjob, jobResponse.Pods...)

				if *jobResponse.Status.Active > 0 {
					runningJobs = append(runningJobs, jobResponse)
				}

				for _, pod := range jobResponse.Pods {
					if pod.Phase == string(corev1.PodRunning) {
						runningPods = append(runningPods, pod)
					}
				}
			}
		}
	}

	log.Printf("There are %d jobs for cronjob %s in the cluster\n", len(jobs.Items), jobUID)

	return allJobsForCronJob, allJobPodsForCronjob, runningJobs, runningPods
}

func (s Sk8lServer) cronJobResponse(cronJob batchv1.CronJob) *protos.CronjobResponse {
	allJobsForNamespace := s.K8sClient.GetAllJobs()

	allJobsForCronJob, jobPodsForCronJob, runningJobs, runningJobPods := s.allAndRunningJobsAnPods(allJobsForNamespace, cronJob.UID)
	lastDuration := getLastDuration(allJobsForCronJob)
	currentDuration := getCurrentDuration(runningJobs)

	commands := buildCronJobCommand(cronJob)
	lastSuccessfulTime, lastScheduleTime := buildLastTimes(cronJob)

	cjr := &protos.CronjobResponse{
		Name:               cronJob.Name,
		Namespace:          cronJob.Namespace,
		Uid:                string(cronJob.UID),
		ContainerCommands:  commands,
		Definition:         cronJob.Spec.Schedule,
		CreationTimestamp:  cronJob.GetCreationTimestamp().UTC().Format(time.RFC3339),
		LastSuccessfulTime: lastSuccessfulTime,
		LastScheduleTime:   lastScheduleTime,
		Active:             len(cronJob.Status.Active) > 0,
		LastDuration:       lastDuration,
		CurrentDuration:    currentDuration,
		Jobs:               allJobsForCronJob,
		RunningJobs:        runningJobs,
		RunningJobsPods:    runningJobPods,
		JobsPods:           jobPodsForCronJob,
		Spec:               &cronJob.Spec,
	}

	return cjr
}

func terminatedAndFailedContainers(pod *corev1.Pod) (*protos.TerminatedContainers, *protos.TerminatedContainers) {
	terminatedEphContainers := make([]*protos.ContainerResponse, 0)
	terminatedInitContainers := make([]*protos.ContainerResponse, 0)
	terminatedContainers := make([]*protos.ContainerResponse, 0)
	failedEphContainers := make([]*protos.ContainerResponse, 0)
	failedInitContainers := make([]*protos.ContainerResponse, 0)
	failedContainers := make([]*protos.ContainerResponse, 0)
	terminatedReasons := make([]*protos.TerminationReason, 0)

	for _, container := range pod.Status.EphemeralContainerStatuses {
		// ephStates = append(ephStates, container.State)
		// if container.State.Waiting != nil && container.State.Waiting.Reason == "Error" {
		//      failedEphContainers = append(failedEphContainers, &container)
		// }
		container := container
		l := []*corev1.PodCondition{}
		for _, pc := range pod.Status.Conditions {
			pc := pc
			l = append(l, &pc)
		}

		if container.State.Terminated != nil {
			cr := &protos.ContainerResponse{
				Status:     &container,
				Phase:      string(pod.Status.Phase),
				Conditions: l,
			}
			terminatedEphContainers = append(terminatedEphContainers, cr)
			if container.State.Terminated.Reason == "Error" {
				cr.TerminatedReason = &protos.TerminationReason{
					TerminationDetails: container.State.Terminated,
					ContainerName:      cr.Status.Name,
				}
				terminatedReasons = append(terminatedReasons, cr.TerminatedReason)
				failedEphContainers = append(failedEphContainers, cr)
			}
		}
	}

	for _, container := range pod.Status.InitContainerStatuses {
		container := container
		l := []*corev1.PodCondition{}
		for _, pc := range pod.Status.Conditions {
			pc := pc
			l = append(l, &pc)
		}

		if container.State.Terminated != nil {
			cr := &protos.ContainerResponse{
				Status:     &container,
				Phase:      string(pod.Status.Phase),
				Conditions: l,
			}

			terminatedInitContainers = append(terminatedInitContainers, cr)
			if container.State.Terminated.Reason == "Error" {
				cr.TerminatedReason = &protos.TerminationReason{
					TerminationDetails: container.State.Terminated,
					ContainerName:      cr.Status.Name,
				}
				terminatedReasons = append(terminatedReasons, cr.TerminatedReason)
				failedInitContainers = append(failedInitContainers, cr)
			}
		}
	}

	for _, container := range pod.Status.ContainerStatuses {
		container := container
		l := []*corev1.PodCondition{}
		for _, pc := range pod.Status.Conditions {
			pc := pc
			l = append(l, &pc)
		}

		cr := &protos.ContainerResponse{
			Status:     &container,
			Phase:      string(pod.Status.Phase),
			Conditions: l,
		}

		if container.State.Terminated != nil {
			terminatedContainers = append(terminatedContainers, cr)
			if container.State.Terminated.Reason == "Error" {
				cr.TerminatedReason = &protos.TerminationReason{
					TerminationDetails: container.State.Terminated,
					ContainerName:      cr.Status.Name,
				}
				terminatedReasons = append(terminatedReasons, cr.TerminatedReason)
				failedContainers = append(failedContainers, cr)
			}
		}
	}

	terminatedContainersResponse := &protos.TerminatedContainers{
		InitContainers:      terminatedInitContainers,
		EphemeralContainers: terminatedEphContainers,
		Containers:          terminatedContainers,
	}

	failedContainersResponse := &protos.TerminatedContainers{
		InitContainers:      failedInitContainers,
		EphemeralContainers: failedEphContainers,
		Containers:          failedContainers,
		TerminationReasons:  terminatedReasons,
	}

	return terminatedContainersResponse, failedContainersResponse
}

func buildJobPodsResponses(gJobPods *corev1.PodList) []*protos.PodResponse {
	n := len(gJobPods.Items)
	jobPodsResponses := make([]*protos.PodResponse, 0, n)

	// ephStates := make([]corev1.ContainerState, 0)
	for _, pod := range gJobPods.Items {
		pod := pod
		// jobPodsForJob.Items[0].Status.ContainerStatuses
		// jobPodsForJob.Items[0].Status.InitContainerStatuses
		terminatedContainers, failedContainers := terminatedAndFailedContainers(&pod)

		failed := len(failedContainers.TerminationReasons) > 0

		jobResponse := &protos.PodResponse{
			Metadata:             &pod.ObjectMeta,
			Spec:                 &pod.Spec,
			Status:               &pod.Status,
			TerminatedContainers: terminatedContainers,
			FailedContainers:     failedContainers,
			Failed:               failed,
			Phase:                string(pod.Status.Phase),
			TerminationReasons:   failedContainers.TerminationReasons,
		}
		jobPodsResponses = append(jobPodsResponses, jobResponse)
	}

	return jobPodsResponses
}

func jobSucceded(job batchv1.Job) bool {
	// The completion time is only set when the job finishes successfully.
	return job.Status.CompletionTime != nil
}

func buildLastTimes(cronJob batchv1.CronJob) (string, string) {
	var lastSuccessfulTime string
	var lastScheduleTime string
	if cronJob.Status.LastSuccessfulTime != nil {
		lastSuccessfulTime = cronJob.Status.LastSuccessfulTime.UTC().Format(time.RFC3339)
	}

	if cronJob.Status.LastScheduleTime != nil {
		lastScheduleTime = cronJob.Status.LastScheduleTime.UTC().Format(time.RFC3339)
	}

	return lastSuccessfulTime, lastScheduleTime
}

func buildCronJobCommand(cronJob batchv1.CronJob) map[string]*protos.ContainerCommands {
	// cronJob.Spec.JobTemplate.Spec.Template.Spec.Containers
	// cronJob.Spec.JobTemplate.Spec.Template.Spec.InitContainers.Image
	// cronJob.Spec.JobTemplate.Spec.Template.Spec.InitContainers.Command
	// spec:
	//   backoffLimit:6
	//   commentmpletionMode: NonIndexed
	//   completions: 1
	//   parallelism: 1
	commands := make(map[string]*protos.ContainerCommands)
	n := len(cronJob.Spec.JobTemplate.Spec.Template.Spec.InitContainers)
	initContainersCommands := make([]string, 0, n)
	var command bytes.Buffer
	for _, container := range cronJob.Spec.JobTemplate.Spec.Template.Spec.InitContainers {
		for _, ccmd := range container.Command {
			command.WriteString(fmt.Sprintf("%s ", ccmd))
		}
		initContainersCommands = append(initContainersCommands, command.String())
		command.Reset()
	}

	containersCommands := make([]string, 0, n)
	for _, container := range cronJob.Spec.JobTemplate.Spec.Template.Spec.Containers {
		for _, ccmd := range container.Command {
			command.WriteString(fmt.Sprintf("%s ", ccmd))
		}
		containersCommands = append(containersCommands, command.String())
		command.Reset()
	}

	ephemeralContainersinersCommands := make([]string, 0, n)
	for _, container := range cronJob.Spec.JobTemplate.Spec.Template.Spec.EphemeralContainers {
		for _, ccmd := range container.Command {
			command.WriteString(fmt.Sprintf("%s ", ccmd))
		}
		ephemeralContainersinersCommands = append(ephemeralContainersinersCommands, command.String())
		command.Reset()
	}

	commands["InitContainers"] = &protos.ContainerCommands{
		Commands: initContainersCommands,
	}
	commands["Containers"] = &protos.ContainerCommands{
		Commands: containersCommands,
	}
	commands["EphemeralContainers"] = &protos.ContainerCommands{
		Commands: ephemeralContainersinersCommands,
	}

	return commands
}

func getCurrentDuration(runningJobsForCronJob []*protos.JobResponse) int64 {
	var lastDuration int64
	if len(runningJobsForCronJob) > 0 {
		last := runningJobsForCronJob[len(runningJobsForCronJob)-1]
		if last.DurationInS != 0 {
			lastDuration = last.DurationInS
		}
	}

	return lastDuration
}

func getLastDuration(allJobsForCronJob []*protos.JobResponse) int64 {
	var lastDuration int64
	if len(allJobsForCronJob) > 0 {
		var i int
		if len(allJobsForCronJob) > 2 {
			i = 2
		} else {
			i = 1
		}
		last := allJobsForCronJob[len(allJobsForCronJob)-i]
		lastDuration = last.DurationInS
	}

	return lastDuration
}

func toDuration(job batchv1.Job, jobFailed bool, failureCondition *batchv1.JobCondition) time.Duration {
	var d time.Duration

	status := job.Status
	if jobFailed && failureCondition != nil {
		d = failureCondition.LastTransitionTime.Sub(status.StartTime.Time)
		return d
	}

	if status.StartTime == nil {
		return d
	}

	switch {
	case status.CompletionTime == nil:
		d = time.Since(status.StartTime.Time)
	default:
		d = status.CompletionTime.Sub(status.StartTime.Time)
	}

	return d
	// return duration.HumanDuration(d)
}

func toDurationInS(job batchv1.Job, jobFailed bool, failureCondition *batchv1.JobCondition) int64 {
	d := toDuration(job, jobFailed, failureCondition)

	return int64(d.Seconds())
}

func toCompletionTimeInS(job batchv1.Job) int64 {
	if job.Status.CompletionTime != nil {
		return job.Status.CompletionTime.Unix()
	}

	return int64(0)
}

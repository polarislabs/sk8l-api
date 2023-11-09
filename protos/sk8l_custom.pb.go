// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v4.24.3
// source: sk8l_custom.proto

package protos

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	v1 "k8s.io/api/batch/v1"
	v11 "k8s.io/apimachinery/pkg/apis/meta/v1"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// JobStatus represents the current state of a Job.
type JobStatus struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The latest available observations of an object's current state. When a Job
	// fails, one of the conditions will have type "Failed" and status true. When
	// a Job is suspended, one of the conditions will have type "Suspended" and
	// status true; when the Job is resumed, the status of this condition will
	// become false. When a Job is completed, one of the conditions will have
	// type "Complete" and status true.
	// More info: https://kubernetes.io/docs/concepts/workloads/controllers/jobs-run-to-completion/
	// +optional
	// +patchMergeKey=type
	// +patchStrategy=merge
	// +listType=atomic
	Conditions []*v1.JobCondition `protobuf:"bytes,1,rep,name=conditions,proto3" json:"conditions,omitempty"`
	// Represents time when the job controller started processing a job. When a
	// Job is created in the suspended state, this field is not set until the
	// first time it is resumed. This field is reset every time a Job is resumed
	// from suspension. It is represented in RFC3339 form and is in UTC.
	// +optional
	StartTime *v11.Time `protobuf:"bytes,2,opt,name=startTime,proto3,oneof" json:"startTime,omitempty"`
	// Represents time when the job was completed. It is not guaranteed to
	// be set in happens-before order across separate operations.
	// It is represented in RFC3339 form and is in UTC.
	// The completion time is only set when the job finishes successfully.
	// +optional
	CompletionTime *v11.Time `protobuf:"bytes,3,opt,name=completionTime,proto3,oneof" json:"completionTime,omitempty"`
	// The number of pending and running pods.
	// +optional
	Active *int32 `protobuf:"varint,4,opt,name=active,proto3,oneof" json:"active,omitempty"`
	// The number of pods which reached phase Succeeded.
	// +optional
	Succeeded *int32 `protobuf:"varint,5,opt,name=succeeded,proto3,oneof" json:"succeeded,omitempty"`
	// The number of pods which reached phase Failed.
	// +optional
	Failed *int32 `protobuf:"varint,6,opt,name=failed,proto3,oneof" json:"failed,omitempty"`
	// The number of pods which are terminating (in phase Pending or Running
	// and have a deletionTimestamp).
	//
	// This field is alpha-level. The job controller populates the field when
	// the feature gate JobPodReplacementPolicy is enabled (disabled by default).
	// +optional
	Terminating *int32 `protobuf:"varint,11,opt,name=terminating,proto3,oneof" json:"terminating,omitempty"`
	// completedIndexes holds the completed indexes when .spec.completionMode =
	// "Indexed" in a text format. The indexes are represented as decimal integers
	// separated by commas. The numbers are listed in increasing order. Three or
	// more consecutive numbers are compressed and represented by the first and
	// last element of the series, separated by a hyphen.
	// For example, if the completed indexes are 1, 3, 4, 5 and 7, they are
	// represented as "1,3-5,7".
	// +optional
	CompletedIndexes *string `protobuf:"bytes,7,opt,name=completedIndexes,proto3,oneof" json:"completedIndexes,omitempty"`
	// FailedIndexes holds the failed indexes when backoffLimitPerIndex=true.
	// The indexes are represented in the text format analogous as for the
	// `completedIndexes` field, ie. they are kept as decimal integers
	// separated by commas. The numbers are listed in increasing order. Three or
	// more consecutive numbers are compressed and represented by the first and
	// last element of the series, separated by a hyphen.
	// For example, if the failed indexes are 1, 3, 4, 5 and 7, they are
	// represented as "1,3-5,7".
	// This field is alpha-level. It can be used when the `JobBackoffLimitPerIndex`
	// feature gate is enabled (disabled by default).
	// +optional
	FailedIndexes *string `protobuf:"bytes,10,opt,name=failedIndexes,proto3,oneof" json:"failedIndexes,omitempty"`
	// uncountedTerminatedPods holds the UIDs of Pods that have terminated but
	// the job controller hasn't yet accounted for in the status counters.
	//
	// The job controller creates pods with a finalizer. When a pod terminates
	// (succeeded or failed), the controller does three steps to account for it
	// in the job status:
	//
	// 1. Add the pod UID to the arrays in this field.
	// 2. Remove the pod finalizer.
	// 3. Remove the pod UID from the arrays while increasing the corresponding
	//     counter.
	//
	// Old jobs might not be tracked using this field, in which case the field
	// remains null.
	// +optional
	UncountedTerminatedPods *v1.UncountedTerminatedPods `protobuf:"bytes,8,opt,name=uncountedTerminatedPods,proto3,oneof" json:"uncountedTerminatedPods,omitempty"`
	// The number of pods which have a Ready condition.
	//
	// This field is beta-level. The job controller populates the field when
	// the feature gate JobReadyPods is enabled (enabled by default).
	// +optional
	Ready *int32 `protobuf:"varint,9,opt,name=ready,proto3,oneof" json:"ready,omitempty"`
	//////// sk8l custom
	StartTimeInS      int64 `protobuf:"varint,880,opt,name=startTimeInS,json=start_time_in_s,proto3" json:"startTimeInS,omitempty"`
	CompletionTimeInS int64 `protobuf:"varint,881,opt,name=completionTimeInS,json=completion_time_in_s,proto3" json:"completionTimeInS,omitempty"`
}

func (x *JobStatus) Reset() {
	*x = JobStatus{}
	if protoimpl.UnsafeEnabled {
		mi := &file_sk8l_custom_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *JobStatus) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*JobStatus) ProtoMessage() {}

func (x *JobStatus) ProtoReflect() protoreflect.Message {
	mi := &file_sk8l_custom_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use JobStatus.ProtoReflect.Descriptor instead.
func (*JobStatus) Descriptor() ([]byte, []int) {
	return file_sk8l_custom_proto_rawDescGZIP(), []int{0}
}

func (x *JobStatus) GetConditions() []*v1.JobCondition {
	if x != nil {
		return x.Conditions
	}
	return nil
}

func (x *JobStatus) GetStartTime() *v11.Time {
	if x != nil {
		return x.StartTime
	}
	return nil
}

func (x *JobStatus) GetCompletionTime() *v11.Time {
	if x != nil {
		return x.CompletionTime
	}
	return nil
}

func (x *JobStatus) GetActive() int32 {
	if x != nil && x.Active != nil {
		return *x.Active
	}
	return 0
}

func (x *JobStatus) GetSucceeded() int32 {
	if x != nil && x.Succeeded != nil {
		return *x.Succeeded
	}
	return 0
}

func (x *JobStatus) GetFailed() int32 {
	if x != nil && x.Failed != nil {
		return *x.Failed
	}
	return 0
}

func (x *JobStatus) GetTerminating() int32 {
	if x != nil && x.Terminating != nil {
		return *x.Terminating
	}
	return 0
}

func (x *JobStatus) GetCompletedIndexes() string {
	if x != nil && x.CompletedIndexes != nil {
		return *x.CompletedIndexes
	}
	return ""
}

func (x *JobStatus) GetFailedIndexes() string {
	if x != nil && x.FailedIndexes != nil {
		return *x.FailedIndexes
	}
	return ""
}

func (x *JobStatus) GetUncountedTerminatedPods() *v1.UncountedTerminatedPods {
	if x != nil {
		return x.UncountedTerminatedPods
	}
	return nil
}

func (x *JobStatus) GetReady() int32 {
	if x != nil && x.Ready != nil {
		return *x.Ready
	}
	return 0
}

func (x *JobStatus) GetStartTimeInS() int64 {
	if x != nil {
		return x.StartTimeInS
	}
	return 0
}

func (x *JobStatus) GetCompletionTimeInS() int64 {
	if x != nil {
		return x.CompletionTimeInS
	}
	return 0
}

var File_sk8l_custom_proto protoreflect.FileDescriptor

var file_sk8l_custom_proto_rawDesc = []byte{
	0x0a, 0x11, 0x73, 0x6b, 0x38, 0x6c, 0x5f, 0x63, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x0b, 0x73, 0x6b, 0x38, 0x6c, 0x5f, 0x63, 0x75, 0x73, 0x74, 0x6f, 0x6d,
	0x1a, 0x34, 0x6b, 0x38, 0x73, 0x2e, 0x69, 0x6f, 0x2f, 0x61, 0x70, 0x69, 0x6d, 0x61, 0x63, 0x68,
	0x69, 0x6e, 0x65, 0x72, 0x79, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x61, 0x70, 0x69, 0x73, 0x2f, 0x6d,
	0x65, 0x74, 0x61, 0x2f, 0x76, 0x31, 0x2f, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x65, 0x64,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x23, 0x6b, 0x38, 0x73, 0x2e, 0x69, 0x6f, 0x2f, 0x61,
	0x70, 0x69, 0x2f, 0x62, 0x61, 0x74, 0x63, 0x68, 0x2f, 0x76, 0x31, 0x2f, 0x67, 0x65, 0x6e, 0x65,
	0x72, 0x61, 0x74, 0x65, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xda, 0x06, 0x0a, 0x09,
	0x4a, 0x6f, 0x62, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x41, 0x0a, 0x0a, 0x63, 0x6f, 0x6e,
	0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x21, 0x2e,
	0x6b, 0x38, 0x73, 0x2e, 0x69, 0x6f, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x62, 0x61, 0x74, 0x63, 0x68,
	0x2e, 0x76, 0x31, 0x2e, 0x4a, 0x6f, 0x62, 0x43, 0x6f, 0x6e, 0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e,
	0x52, 0x0a, 0x63, 0x6f, 0x6e, 0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x4d, 0x0a, 0x09,
	0x73, 0x74, 0x61, 0x72, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x2a, 0x2e, 0x6b, 0x38, 0x73, 0x2e, 0x69, 0x6f, 0x2e, 0x61, 0x70, 0x69, 0x6d, 0x61, 0x63, 0x68,
	0x69, 0x6e, 0x65, 0x72, 0x79, 0x2e, 0x70, 0x6b, 0x67, 0x2e, 0x61, 0x70, 0x69, 0x73, 0x2e, 0x6d,
	0x65, 0x74, 0x61, 0x2e, 0x76, 0x31, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x48, 0x00, 0x52, 0x09, 0x73,
	0x74, 0x61, 0x72, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x88, 0x01, 0x01, 0x12, 0x57, 0x0a, 0x0e, 0x63,
	0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x2a, 0x2e, 0x6b, 0x38, 0x73, 0x2e, 0x69, 0x6f, 0x2e, 0x61, 0x70, 0x69,
	0x6d, 0x61, 0x63, 0x68, 0x69, 0x6e, 0x65, 0x72, 0x79, 0x2e, 0x70, 0x6b, 0x67, 0x2e, 0x61, 0x70,
	0x69, 0x73, 0x2e, 0x6d, 0x65, 0x74, 0x61, 0x2e, 0x76, 0x31, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x48,
	0x01, 0x52, 0x0e, 0x63, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x69, 0x6d,
	0x65, 0x88, 0x01, 0x01, 0x12, 0x1b, 0x0a, 0x06, 0x61, 0x63, 0x74, 0x69, 0x76, 0x65, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x05, 0x48, 0x02, 0x52, 0x06, 0x61, 0x63, 0x74, 0x69, 0x76, 0x65, 0x88, 0x01,
	0x01, 0x12, 0x21, 0x0a, 0x09, 0x73, 0x75, 0x63, 0x63, 0x65, 0x65, 0x64, 0x65, 0x64, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x05, 0x48, 0x03, 0x52, 0x09, 0x73, 0x75, 0x63, 0x63, 0x65, 0x65, 0x64, 0x65,
	0x64, 0x88, 0x01, 0x01, 0x12, 0x1b, 0x0a, 0x06, 0x66, 0x61, 0x69, 0x6c, 0x65, 0x64, 0x18, 0x06,
	0x20, 0x01, 0x28, 0x05, 0x48, 0x04, 0x52, 0x06, 0x66, 0x61, 0x69, 0x6c, 0x65, 0x64, 0x88, 0x01,
	0x01, 0x12, 0x25, 0x0a, 0x0b, 0x74, 0x65, 0x72, 0x6d, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6e, 0x67,
	0x18, 0x0b, 0x20, 0x01, 0x28, 0x05, 0x48, 0x05, 0x52, 0x0b, 0x74, 0x65, 0x72, 0x6d, 0x69, 0x6e,
	0x61, 0x74, 0x69, 0x6e, 0x67, 0x88, 0x01, 0x01, 0x12, 0x2f, 0x0a, 0x10, 0x63, 0x6f, 0x6d, 0x70,
	0x6c, 0x65, 0x74, 0x65, 0x64, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x65, 0x73, 0x18, 0x07, 0x20, 0x01,
	0x28, 0x09, 0x48, 0x06, 0x52, 0x10, 0x63, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x49,
	0x6e, 0x64, 0x65, 0x78, 0x65, 0x73, 0x88, 0x01, 0x01, 0x12, 0x29, 0x0a, 0x0d, 0x66, 0x61, 0x69,
	0x6c, 0x65, 0x64, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x65, 0x73, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x09,
	0x48, 0x07, 0x52, 0x0d, 0x66, 0x61, 0x69, 0x6c, 0x65, 0x64, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x65,
	0x73, 0x88, 0x01, 0x01, 0x12, 0x6b, 0x0a, 0x17, 0x75, 0x6e, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x65,
	0x64, 0x54, 0x65, 0x72, 0x6d, 0x69, 0x6e, 0x61, 0x74, 0x65, 0x64, 0x50, 0x6f, 0x64, 0x73, 0x18,
	0x08, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2c, 0x2e, 0x6b, 0x38, 0x73, 0x2e, 0x69, 0x6f, 0x2e, 0x61,
	0x70, 0x69, 0x2e, 0x62, 0x61, 0x74, 0x63, 0x68, 0x2e, 0x76, 0x31, 0x2e, 0x55, 0x6e, 0x63, 0x6f,
	0x75, 0x6e, 0x74, 0x65, 0x64, 0x54, 0x65, 0x72, 0x6d, 0x69, 0x6e, 0x61, 0x74, 0x65, 0x64, 0x50,
	0x6f, 0x64, 0x73, 0x48, 0x08, 0x52, 0x17, 0x75, 0x6e, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x65, 0x64,
	0x54, 0x65, 0x72, 0x6d, 0x69, 0x6e, 0x61, 0x74, 0x65, 0x64, 0x50, 0x6f, 0x64, 0x73, 0x88, 0x01,
	0x01, 0x12, 0x19, 0x0a, 0x05, 0x72, 0x65, 0x61, 0x64, 0x79, 0x18, 0x09, 0x20, 0x01, 0x28, 0x05,
	0x48, 0x09, 0x52, 0x05, 0x72, 0x65, 0x61, 0x64, 0x79, 0x88, 0x01, 0x01, 0x12, 0x26, 0x0a, 0x0c,
	0x73, 0x74, 0x61, 0x72, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x49, 0x6e, 0x53, 0x18, 0xf0, 0x06, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x0f, 0x73, 0x74, 0x61, 0x72, 0x74, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x5f,
	0x69, 0x6e, 0x5f, 0x73, 0x12, 0x30, 0x0a, 0x11, 0x63, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x69,
	0x6f, 0x6e, 0x54, 0x69, 0x6d, 0x65, 0x49, 0x6e, 0x53, 0x18, 0xf1, 0x06, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x14, 0x63, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x74, 0x69, 0x6d,
	0x65, 0x5f, 0x69, 0x6e, 0x5f, 0x73, 0x42, 0x0c, 0x0a, 0x0a, 0x5f, 0x73, 0x74, 0x61, 0x72, 0x74,
	0x54, 0x69, 0x6d, 0x65, 0x42, 0x11, 0x0a, 0x0f, 0x5f, 0x63, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74,
	0x69, 0x6f, 0x6e, 0x54, 0x69, 0x6d, 0x65, 0x42, 0x09, 0x0a, 0x07, 0x5f, 0x61, 0x63, 0x74, 0x69,
	0x76, 0x65, 0x42, 0x0c, 0x0a, 0x0a, 0x5f, 0x73, 0x75, 0x63, 0x63, 0x65, 0x65, 0x64, 0x65, 0x64,
	0x42, 0x09, 0x0a, 0x07, 0x5f, 0x66, 0x61, 0x69, 0x6c, 0x65, 0x64, 0x42, 0x0e, 0x0a, 0x0c, 0x5f,
	0x74, 0x65, 0x72, 0x6d, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6e, 0x67, 0x42, 0x13, 0x0a, 0x11, 0x5f,
	0x63, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x65, 0x73,
	0x42, 0x10, 0x0a, 0x0e, 0x5f, 0x66, 0x61, 0x69, 0x6c, 0x65, 0x64, 0x49, 0x6e, 0x64, 0x65, 0x78,
	0x65, 0x73, 0x42, 0x1a, 0x0a, 0x18, 0x5f, 0x75, 0x6e, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x65, 0x64,
	0x54, 0x65, 0x72, 0x6d, 0x69, 0x6e, 0x61, 0x74, 0x65, 0x64, 0x50, 0x6f, 0x64, 0x73, 0x42, 0x08,
	0x0a, 0x06, 0x5f, 0x72, 0x65, 0x61, 0x64, 0x79, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_sk8l_custom_proto_rawDescOnce sync.Once
	file_sk8l_custom_proto_rawDescData = file_sk8l_custom_proto_rawDesc
)

func file_sk8l_custom_proto_rawDescGZIP() []byte {
	file_sk8l_custom_proto_rawDescOnce.Do(func() {
		file_sk8l_custom_proto_rawDescData = protoimpl.X.CompressGZIP(file_sk8l_custom_proto_rawDescData)
	})
	return file_sk8l_custom_proto_rawDescData
}

var file_sk8l_custom_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_sk8l_custom_proto_goTypes = []interface{}{
	(*JobStatus)(nil),                  // 0: sk8l_custom.JobStatus
	(*v1.JobCondition)(nil),            // 1: k8s.io.api.batch.v1.JobCondition
	(*v11.Time)(nil),                   // 2: k8s.io.apimachinery.pkg.apis.meta.v1.Time
	(*v1.UncountedTerminatedPods)(nil), // 3: k8s.io.api.batch.v1.UncountedTerminatedPods
}
var file_sk8l_custom_proto_depIdxs = []int32{
	1, // 0: sk8l_custom.JobStatus.conditions:type_name -> k8s.io.api.batch.v1.JobCondition
	2, // 1: sk8l_custom.JobStatus.startTime:type_name -> k8s.io.apimachinery.pkg.apis.meta.v1.Time
	2, // 2: sk8l_custom.JobStatus.completionTime:type_name -> k8s.io.apimachinery.pkg.apis.meta.v1.Time
	3, // 3: sk8l_custom.JobStatus.uncountedTerminatedPods:type_name -> k8s.io.api.batch.v1.UncountedTerminatedPods
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_sk8l_custom_proto_init() }
func file_sk8l_custom_proto_init() {
	if File_sk8l_custom_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_sk8l_custom_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*JobStatus); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	file_sk8l_custom_proto_msgTypes[0].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_sk8l_custom_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_sk8l_custom_proto_goTypes,
		DependencyIndexes: file_sk8l_custom_proto_depIdxs,
		MessageInfos:      file_sk8l_custom_proto_msgTypes,
	}.Build()
	File_sk8l_custom_proto = out.File
	file_sk8l_custom_proto_rawDesc = nil
	file_sk8l_custom_proto_goTypes = nil
	file_sk8l_custom_proto_depIdxs = nil
}

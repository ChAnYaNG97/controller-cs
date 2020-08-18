/*
Copyright 2020 yangchen.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	csv1 "controller-cs/api/v1"
)

// AliyunCKReconciler reconciles a AliyunCK object
type AliyunCKReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=cs.cloudplus.io,resources=aliyuncks,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=cs.cloudplus.io,resources=aliyuncks/status,verbs=get;update;patch

func (r *AliyunCKReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("aliyunck", req.NamespacedName)

	// your logic here
	// 1. Get，获得AliyunCK对象
	ack := &csv1.AliyunCK{}
	if err := r.Get(ctx, req.NamespacedName, ack); err != nil {
		return reconcile.Result{}, client.IgnoreNotFound(err)

	}

		switch ack.Status.Phase {
	case csv1.Phase_None:
		ack.Status.Phase = csv1.Phase_Create
		if err := r.Update(ctx, ack); err != nil {
			log.Error(err, "update status from none to initing error")
			return ctrl.Result{}, err
		}
	case csv1.Phase_Create:
		// 判断是不是exist
		// 这里先简单通过是否有cluster_id来判断
		if ack.Status.ClusterId == "" {
			// 不存在，就创建
		} else {
			// 存在去查状态
			// status = getStatus();
			status := "creating"
			if status == "running" {
				// 更新状态
			} else {
				// 不变
			}
		}



	case csv1.Phase_Run:

	case csv1.Phase_Delete:


	}

	//
	//// 2. 删除
	//myFinalizerName := "finalizer.cloudplus.io"
	//
	//
	//if ack.ObjectMeta.DeletionTimestamp.IsZero() {
	//	// deleteTimestamp为空，说明没有进入删除状态
	//	// 加上finalizer
	//	ack.ObjectMeta.Finalizers = append(ack.ObjectMeta.Finalizers, myFinalizerName)
	//	if err := r.Update(ctx, ack); err != nil {
	//		r.Log.Error(err, "add finalizer error", ack.Name)
	//	}
	//} else {
	//	// deleteTimestamp不为空，进入删除状态
	//	if containsString(ack.Finalizers, myFinalizerName) {
	//		ack.ObjectMeta.Finalizers = removeString(ack.ObjectMeta.Finalizers, myFinalizerName)
	//
	//		if err := r.Update(ctx, ack); err != nil {
	//			return ctrl.Result{}, nil
	//		}
	//	}
	//
	//	return ctrl.Result{}, nil
	//}
	//// 3. 创建/更新
	//
	//
	//// 4. 更新状态？
	//
	//// 5. 记录结果

	return ctrl.Result{}, nil
}

func (r *AliyunCKReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&csv1.AliyunCK{}).
		Complete(r)
}

func containsString(strs []string, s string) bool {
	for _, str := range strs {
		if str == s {
			return true
		}
	}

	return false

}

func removeString(strs []string, s string) (result []string) {
	for _, str := range strs {
		if str != s {
			result = append(result, str)
		}
	}

	return
}

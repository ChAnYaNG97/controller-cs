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

	csv1 "controller-cs/api/v1"
	"controller-cs/driver"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"time"
)

const (
	FinalizerName = "finalizer.cloudplus.io"
)

// AliyunCKReconciler reconciles a AliyunCK object
type AliyunCKReconciler struct {
	client.Client
	Log         logr.Logger
	Scheme      *runtime.Scheme
	CloudClient driver.ClientInterface
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
			// 不存在，创建
			req := &driver.CreateRequest{
				Name:         ack.Name,
				Region:       ack.Spec.Region,
				InstanceType: ack.Spec.InstanceType,
			}
			resp, err := r.CloudClient.CreateCluster(req)
			if err != nil {
				log.Error(err, "create cluster error")
				return ctrl.Result{}, err
			}

			ack.Status.ClusterId = resp.ClusterId
			ack.Status.VPCId = resp.VPCId
			ack.Status.Phase = csv1.Phase_Create

			if err := r.Update(ctx, ack); err != nil {
				log.Error(err, "update status cid and vid error")
				return ctrl.Result{}, err
			}
		} else {
			// 存在去查状态
			// status = getStatus();
			req := &driver.CommonRequest{
				ClusterId: ack.Status.ClusterId,
			}
			resp, err := r.CloudClient.GetClusterStatus(req)
			if err != nil {
				return ctrl.Result{}, err
			}
			status := resp.Status
			if status == "running" {
				// 更新状态
				ack.Status.Phase = csv1.Phase_Run
				if err := r.Update(ctx, ack); err != nil {
					log.Error(err, "update status cid and vid error")
					return ctrl.Result{}, err
				}

			}
			return ctrl.Result{
				RequeueAfter: 1 * time.Minute,
			}, nil
		}
	case csv1.Phase_Run:
		// 判断是否有deleteTimestamp
		if !ack.DeletionTimestamp.IsZero() {
			// 要删除的
			if containsString(ack.ObjectMeta.Finalizers, FinalizerName) {
				//存在finalizer，删除finalizer
				//先把删除命令下了
				req := &driver.CommonRequest{
					ClusterId: ack.Status.ClusterId,
				}
				err := r.CloudClient.DeleteCluster(req)
				if err != nil {
					log.Error(err, "driver delete cluster error")
					return ctrl.Result{}, err
				}

				ack.Status.Phase = csv1.Phase_Delete

				if err := r.Update(ctx, ack); err != nil {
					log.Error(err, "update status error")
					return ctrl.Result{}, err
				}
			}
		} else {
			// deleteTimestamp 不存在
			if !containsString(ack.ObjectMeta.Finalizers, FinalizerName) {
				// 没有finalizer，加上finalizer
				ack.ObjectMeta.Finalizers = append(ack.ObjectMeta.Finalizers, FinalizerName)
				if err := r.Update(ctx, ack); err != nil {
					log.Error(err, "add finalizer error")
				}
			}
		}

	case csv1.Phase_Delete:
		// 判断阿里云集群的实际状态
		req := &driver.CommonRequest{
			ClusterId: ack.Status.ClusterId,
		}
		resp, err := r.CloudClient.IsClusterExist(req)
		if err != nil {
			log.Error(err, "driver get status error")
			return ctrl.Result{}, err
		}
		if !resp.Exist {
			// 阿里云上已经被删除了，就删除finalizer进入gc状态
			ack.ObjectMeta.Finalizers = removeString(ack.ObjectMeta.Finalizers, FinalizerName)
			if err := r.Update(ctx, ack); err != nil {
				log.Error(err, "delete finalizer error")
				return ctrl.Result{}, err
			}

		}

		return ctrl.Result{
			RequeueAfter: 1 * time.Minute,
		}, nil
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

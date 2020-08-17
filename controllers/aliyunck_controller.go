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
	_ = context.Background()
	_ = r.Log.WithValues("aliyunck", req.NamespacedName)

	// your logic here

	return ctrl.Result{}, nil
}

func (r *AliyunCKReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&csv1.AliyunCK{}).
		Complete(r)
}

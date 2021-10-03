/*
Copyright 2021.

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

	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	testv1 "leo/api/v1"
)

// LeoPodReconciler reconciles a LeoPod object
type LeoPodReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=test.leo.io,resources=leopods,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=test.leo.io,resources=leopods/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=test.leo.io,resources=leopods/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the LeoPod object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *LeoPodReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// your logic here
	leopod := testv1.LeoPod{}
	if err := r.Client.Get(ctx, req.NamespacedName, &leopod); err != nil {
		klog.Errorln("fail to get my kind resource: ", err)
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	deployment := &apps.Deployment{}
	err := r.Client.Get(ctx, client.ObjectKey{
		Name:      leopod.Name,
		Namespace: leopod.Namespace,
	}, deployment)
	// //need to delete deployment
	// if leopod.ObjectMeta.DeletionTimestamp.IsZero() {
	// 	err := r.Client.Delete(ctx, deployment)
	// 	if err != nil {
	// 		klog.Errorln("delete deployment:", deployment.Name, err)
	// 		return ctrl.Result{}, err
	// 	}
	// 	klog.Infoln("delete deployment:", deployment.Name)
	// 	return ctrl.Result{}, nil
	// }
	klog.Warningln("get deployment info response:", err)
	if apierrors.IsNotFound(err) {
		// create my deployment
		deployment = r.buildDeployment(leopod)
		if err := r.Client.Create(ctx, deployment); err != nil {
			klog.Errorln("fail to create delployment:", err)
			return ctrl.Result{}, err
		}
		klog.Infof("create deloyment %s resource for leo\n", deployment.Name)
		return ctrl.Result{}, nil
	}
	if err != nil {
		klog.Errorln("fail to get delpoyments for leo pod :", err)
		return ctrl.Result{}, err
	}
	klog.Info("existing Deployment resource already exists for leo pod, checking replica count")
	expectReplicas := leopod.Spec.Replicas
	if deployment.Spec.Replicas != &expectReplicas {
		klog.Info("updating replica count", "old_count", *deployment.Spec.Replicas, "new_count", expectReplicas)
		deployment.Spec.Replicas = &expectReplicas
		if err := r.Client.Update(ctx, deployment); err != nil {
			klog.Error(err, "failed to Deployment update replica count")
			return ctrl.Result{}, err
		}
		klog.Warningf("Scaled deployment %+v to %+v replicas", deployment.Spec.Replicas, expectReplicas)
		return ctrl.Result{}, nil
	}
	klog.Infoln("everything is ok,do nothing")
	leopod.Status.ReadyReplicas = *deployment.Spec.Replicas
	if r.Client.Status().Update(ctx, &leopod); err != nil {
		klog.Errorln("fail to update MyKind Status :", err)
		return ctrl.Result{}, err
	}
	klog.Info("resource status synced")
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *LeoPodReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&testv1.LeoPod{}).
		Complete(r)
}
func (r *LeoPodReconciler) buildDeployment(leopod testv1.LeoPod) *apps.Deployment {
	d := apps.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      leopod.Name,
			Namespace: leopod.ObjectMeta.Namespace,
		},
		Spec: apps.DeploymentSpec{
			Replicas: &leopod.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"leopods.test.leo.io": leopod.Name,
				},
			},
			Template: core.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"leopods.test.leo.io": leopod.Name,
					},
				},
				Spec: core.PodSpec{
					Containers: []core.Container{
						{
							Name:            leopod.Name,
							Image:           leopod.Spec.Image,
							ImagePullPolicy: core.PullIfNotPresent,
						},
					},
				},
			},
		},
	}
	return &d
}

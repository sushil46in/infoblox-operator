/*
Copyright 2021 Sushil Saini.

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
	k8serr "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	infoboxv1alpha1 "infoblox-operator/api/v1alpha1"
)

// DnsrecordReconciler reconciles a Dnsrecord object
type DnsrecordReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=infobox.infobloxdns.io,resources=dnsrecords,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=infobox.infobloxdns.io,resources=dnsrecords/status,verbs=get;update;patch

func (r *DnsrecordReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("dnsrecord", req.NamespacedName)
	vg := &infoboxv1alpha1.Dnsrecord{}

	if err := r.Client.Get(ctx, req.NamespacedName, vg); err != nil {
		if !k8serr.IsNotFound(err) {
			log.Error(err, "unable to fetch Dnsrecords")
		}
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	myFinalizerName := "dns-record"

	// examine DeletionTimestamp to determine if object is under deletion
	if vg.ObjectMeta.DeletionTimestamp.IsZero() {
		// The object is not being deleted, so if it does not have our finalizer,
		// then lets add the finalizer and update the object. This is equivalent
		// registering our finalizer.
		if !containsString(vg.ObjectMeta.Finalizers, myFinalizerName) {
			vg.ObjectMeta.Finalizers = append(vg.ObjectMeta.Finalizers, myFinalizerName)
			if err := r.Update(context.Background(), vg); err != nil {
				return ctrl.Result{}, err
			}
			sendmsg(vg.Spec.Hostname, vg.Spec.Ipaddress)
		} else {
			log.Info("DNS record has been updated")
		}
	} else {
		// The object is being deleted
		if containsString(vg.ObjectMeta.Finalizers, myFinalizerName) {
			// our finalizer is present, so lets handle any external dependency
			sendmsg("dd", "sdf")

			// remove our finalizer from the list and update it.
			vg.ObjectMeta.Finalizers = removeString(vg.ObjectMeta.Finalizers, myFinalizerName)
			if err := r.Update(context.Background(), vg); err != nil {
				return ctrl.Result{}, err
			}
		}

		// Stop reconciliation as the item is being deleted
		return ctrl.Result{}, nil
	}

	return ctrl.Result{}, nil
}

func (r *DnsrecordReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&infoboxv1alpha1.Dnsrecord{}).
		Complete(r)
}

func containsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

func removeString(slice []string, s string) (result []string) {
	for _, item := range slice {
		if item == s {
			continue
		}
		result = append(result, item)
	}
	return
}

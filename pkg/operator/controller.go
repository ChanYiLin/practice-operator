/*
Copyright 2018 Jack Lin (jacklin@laslab.cs.nthu.edu)

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


package main


import (
    "fmt"
    "reflect"
    "time"

    opkit "github.com/rook/operator-kit"
    studentv1 "practice-operator/pkg/apis/student/v1"
    nthuclientset "practice-operator/pkg/client/clientset/versioned/typed/nthu/v1"
    "k8s.io/client-go/tools/cache"


    appsv1 "k8s.io/api/apps/v1"
    corev1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

    "k8s.io/apimachinery/pkg/runtime/schema"
    "k8s.io/apimachinery/pkg/util/clock"
    "k8s.io/client-go/tools/cache"
)


var msg = map[studentv1.StudentLifeState]string{
    studentv1.Health:   "我非常健康!!",
    studentv1.Sick:     "我生病了!!",
    studentv1.Dead:     "我掛了!!",
}

// Controller represents a controller object for student custom resources
type Controller struct {
    context         *opkit.Context
    clientset       nthuclientset.NthuV1Interface
    clock           clock.Clock
}


func NewController (context *opkit.Context, clientset nthuclientset.NthuV1Interface) *Controller {
    return &Controller{
        context:    context,
        clientset:  clientset,
        clock:      clock.RealClock{},
    }
}

// Watch watches for instances of student custom resources and acts on them
func (c *Controller) StartWatch(namespace string, stopCh chan struct{}) error {

    resourceHandlers := cache.ResourceEventHandlerFuncs{
        AddFunc:    c.onAdd,
        UpdateFunc: c.onUpdate,
        DeleteFunc: c.onDelete,
    }

    fmt.Printf("start watching resources in namespace %s", namespace)
    
    restClient := c.clientset.RESTClient()
    watcher := opkit.NewWatcher(studentv1.StudentResource, namespace, resourceHandlers, restClient)
    
    go watcher.Watch(&studentv1.Student{}, stopCh)
    
    return nil
}



func (c *Controller) onAdd(obj interface{}) {
    student := obj.(*studentv1.Student).DeepCopy()
    
    fmt.Printf("%s resource onAdd.", student.Name)
    //fmt.Printf("Added Sample '%s' with Hello=%s\n", s.Name, s.Spec.Hello)

    deployment, _ := c.context.Clientset.AppsV1().Deployments(student.Namespace).Create(newDeployment(student))
    c.updateStatus(student, deployment)
}


func (c *Controller) onUpdate(oldObj, newObj interface{}) {
    oldStudent := oldObj.(*studentv1.Student).DeepCopy()
    newStudent := newObj.(*studentv1.Student).DeepCopy()
    fmt.Printf("%s resource onUpdate.", Student.Name)
    //fmt.Printf("Updated sample '%s' from %s to %s\n", newSample.Name, oldSample.Spec.Hello, newSample.Spec.Hello)

    deployment, _ := c.context.Clientset.AppsV1().Deployments(student.Namespace).Update(newDeployment(student))
    c.updateStatus(newStudent, deployment)
}

func (c *Controller) onDelete(obj interface{}) {
    student := obj.(*studentv1.Student).DeepCopy()

    fmt.Printf("%s resource onDelete.", student.Name)

    c.context.Clientset.AppsV1().Deployments(student.Namespace).Delete(student.Name, &metav1.DeleteOptions{})
    //fmt.Printf("Deleted sample '%s' with Hello=%s\n", s.Name, s.Spec.Hello)
}

func newDeployment(student *studentv1.Student) *appsv1.Deployment {
    labels := map[string]string{
        "app":        "nginx",
        "controller": student.Name,
    }

    return &appsv1.Deployment{
        ObjectMeta: metav1.ObjectMeta{
            Name:      student.Spec.TaskName,
            Namespace: student.Namespace,
            OwnerReferences: []metav1.OwnerReference{
                *metav1.NewControllerRef(student, schema.GroupVersionKind{
                    Group:   studentv1.SchemeGroupVersion.Group,
                    Version: studentv1.SchemeGroupVersion.Version,
                    Kind:    "Student",
                }),
            },
        },
        Spec: appsv1.DeploymentSpec{
            Replicas: student.Spec.Threads,
            Selector: &metav1.LabelSelector{
                MatchLabels: labels,
            },
            Template: corev1.PodTemplateSpec{
                ObjectMeta: metav1.ObjectMeta{
                    Labels: labels,
                },
                Spec: corev1.PodSpec{
                    Containers: []corev1.Container{
                        {
                            Name:  "nginx",
                            Image: "nginx:latest",
                        },
                    },
                },
            },
        },
    }
}

func (c *Controller) updateStatus(student *studentv1.Student, deployment *appsv1.Deployment) error {
    studentCopy := student.DeepCopy()
    studentCopy.Status.AvailableThreads = *student.Spec.Threads

    r := studentCopy.Status.AvailableThreads
    switch {
        case r <= 3:
            studentCopy.Status.LiverState = studentv1.Health
        case r <= 9 && r > 3:
            studentCopy.Status.LiverState = studentv1.Sick
        case r > 9:
            studentCopy.Status.LiverState = studentv1.Dead
    }

    studentCopy.Status.Message = msg[studentCopy.Status.LifeState]
    studentCopy.Status.LastLiveTime = metav1.NewTime(c.clock.Now())
    if _, err := c.clientset.Students(namespace).Update(studentCopy); err != nil {
        return fmt.Errorf("failed to update student %s status: %+v", student.Namespace, err)
    }
    return nil
}



















/*
Copyright 2016 Jack Lin (jacklin@lsalab.cs.nthu.edu.tw)

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
    "os"
    "os/signal"
    "syscall"
    "time"

    opkit "github.com/rook/operator-kit"
    
    studentv1 "practice-operator/pkg/apis/student/v1"
    nthuclientset "practice-operator/pkg/client/clientset/versioned/typed/nthu/v1"
    
    "k8s.io/api/core/v1"
    apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/rest"
)



func main() {

    fmt.Println("Getting kubernetes context")
    
    context, clientset, err := createContext()
    if err != nil {
        fmt.Printf("failed to create context. %+v\n", err)
        os.Exit(1)
    }

    // Create and wait for CRD resources
    fmt.Println("Registering the student resource")
    
    resources := []opkit.CustomResource{studentv1.StudentResource} // define student resources

    err = opkit.CreateCustomResources(*context, resources)
    if err != nil {
        fmt.Printf("failed to create custom resource. %+v\n", err)
        os.Exit(1)
    }

    // create signals to stop watching the resources
    signalChan := make(chan os.Signal, 1)
    stopChan := make(chan struct{})
    signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

    // start watching the student resource
    fmt.Println("Watching the student resource")
    
    //TODO
    controller := NewController(context, clientset) //create the contollerfrom controller.go
    controller.StartWatch(v1.NamespaceAll, stopChan) 

    for {
        select {
        case <- signalChan:
            fmt.Println("shutdown signal received, exiting...")
            close(stopChan)
            return
        }
    }
    
}


func createContext() (*opkit.Context, nthuclientset.NthuV1Interface, error) {
    config, err := rest.InClusterConfig()
    if err != nil {
        return nil, nil, fmt.Errorf("failed to get k8s config. %+v", err)
    }

    clientset, err := kubernetes.NewForConfig(config)
    if err != nil {
        return nil, nil, fmt.Errorf("failed to get k8s client. %+v", err)
    }

    apiExtClientset, err := apiextensionsclient.NewForConfig(config)
    if err != nil {
        return nil, nil, fmt.Errorf("failed to create k8s API extension clientset. %+v", err)
    }

    nthuClientset, err := nthuclientset.NewForConfig(config)
    if err != nil {
        return nil, nil, fmt.Errorf("failed to create nthu clientset. %+v", err)
    }

    context := &opkit.Context{
        Clientset:             clientset,
        APIExtensionClientset: apiExtClientset,
        Interval:              500 * time.Millisecond,
        Timeout:               60 * time.Second,
    }
    return context, nthuClientset, nil   
}





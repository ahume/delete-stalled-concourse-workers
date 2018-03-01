# delete-stalled-concourse-worker

In certain situations its possible for workers to [stay in a stalled state](https://github.com/concourse/concourse/issues/1497). This application can be run to prune those workers.

## GKE use case

At Brandwatch we run Concourse on Kubernetes (GKE) using preemptbile VMs. When VMs are preempted we are left with stalled workers, and so run this as a cronjob every minute to clean them up.
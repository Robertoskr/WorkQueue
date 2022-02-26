# WorkQueue

Most simple work queue implementation that you can encounter, not suitable for production environment, its only for learning purposes !!!!!


If you have a docker image of a golang server with predefined endpoints, you can use them, only thing that you need to modify are the proxies in the main function.

/queue_node.go ==> entry point, it holds the work queue and that our server workers are going to consume,
/docker-api ==> example of worker, 

(important) all the workers should be equal, otherwise may be some 404 not found errors!

client.go //client testing the workqueue

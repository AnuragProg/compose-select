



# COMPOSE-SELECT CLI TOOL


## PROBLEM:
When I was developing microservices platform, I had to create new compose file consisting of the services that my application depends on and then run the container.
My platform consists of various services so writing various compose files was repetitive and once an error is there or modification, then all the subsequent files needs that change :-)
Classic "Spend more time in automating the task than manually doing it", but wanted to do something other than that application.

## SOLUTION:
Make user create a compose file consisting of all services in it defining dependencies using "depends_on", then use the compose-select executable to mention the file, then the service that you want to run locally, thereafter enter the output file name. Voila, you have your compose file for that particular service ready.
Whenever there is an update, just make the change in central compose file and then repeat the above process. You have your file updated again.


## ISSUES:
For databases like services, you have to expose the ports manually so that there is no clash of ports and changes be made according to the developers choice.


```
E.g Walkthrough

Enter filename:
compose.yaml

Enter service name:
auth-app

Enter output filename:
compose-auth.yaml

```

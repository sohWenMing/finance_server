# Finance Project: A Project to put all concepts that have been learnt together

This project is to try to integrate different concepts such as:
* Creating and running servers
* Deployment, both in development phases using docker containers and in production
* The usage of Makefiles to streamline development operations
* git, and eventually the use of github actions to streamline deployment to production environments
* deployment to AWS, using EC2 instance for the server, and RDS for the postgres database instance
* Auth using JWT credentials - using the golang JWT library
* unit testing

# The Testing Environment
While it might be a bit overkill, because this project was also used to explore the usage of Docker, all testing should be done with containerization of all the different components of this project which should include
 * server 
 * database 
 * front end (if in the future there is a front end server involved)

# Execution of Commands
The project should be built in such a way that all execution of commands for testing, building and deployment should be done through the use of a Makefile
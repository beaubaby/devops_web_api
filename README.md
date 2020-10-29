# DevOps Instruction - Web application

For detailed explanation on how things work with kotlin spring boot build instruction on pipeline. I've used `bash` script in the main focus.

## Infrastructure as a Code on AWS Cloud Platform.

1.  Create a Github repository for the code they have written (including actions based build workflow)

`git clone https://github.com/beaubaby/devops_web_api.git`

2.  Install terraform and author the following on any cloud provider

 * Additional non-default VPC with internet gateway and route table 
 * Private and Public Subnets
 * SSH Key Setup
 * Virtual Compute instances that run a web service
 * Virtual Compute instances that run a database for leveraging an external source such as RDS
 * Render a simple website that shows information being either pulled out of the data layer or from some 3rd party API
 * Logging enabled to a central place

Alternatively,
 - All of the above, but running on EKS, GKE, or AKS.
 
#### High-Level Infrastructure diagram #### 


 
#### Deployment Instruction ####

#### Usage

To test connection with own database

1.Containerize source code. Run test before containerize.


    ./go lint
    ./go test
    ./go dependency_check
    ./go static_check
    ./go build
    ./go containerize


2.Initial own database
   - Start DB local for dev application on local

    ./go startDb


3.Check your microservice to connect the database on local with set configuration profile

`-Dspring.profiles.active=local`

4. Run 

4.After test passed, can stop your database. Run

    ./go stopDb

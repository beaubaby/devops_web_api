## LOAN_ELIGIBILITY_SERVICE (SKELETON)

This repo contains resources and database configuration scripts for building microservices on AWS EKS cluster in Topup Coreplatfrom

It's separate 3 paths 

- Containerize microservice packages including application configuration and source code.
- Initialize own database, user, password and grant privilege to it including keep the password in [AWS secret manager](https://docs.aws.amazon.com/secretsmanager/latest/userguide/intro.html).
- Apply microservice connect to [AWS RDS](https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/Welcome.html) cluster each environment.

###### Note: Refer RDS cluster in [RDS Cluster repository](https://gitlab.tools.buk0.com/core-platform/rds-cluster)

#### Usage

To test connection with own database

- Test on local

1.Containerize microservice source code. Run test before containerize.


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

4.After test passed, can stop your database. Run

    ./go stopDb


###### Note:

if you need to update some database config. you can update file `/script/init/init-postgres-db.sql`

### Reference GOCD pipeline

loan-eligibility-service : https://gocd-server.tools.buk0.com/go/tab/pipeline/history/loan-eligibility-service
RDS Cluster: https://gocd-server.tools.buk0.com/go/tab/pipeline/history/rds-cluster
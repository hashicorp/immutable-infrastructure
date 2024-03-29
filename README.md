# Immutable-Infrastructure

**BACKGROUND INFORMATION**

**Before Virtualization and Cloud Computing:**
* Server infrastructure was enpasulated in physical servers that were:
  * expensive
  * time consuming to create 
  * setup took days or weeks because of the need to order new hardware, configuration the machine, and other manual processes.
  
* Mutable Infrastructure = older way of confguring infrastructure that was leveraged prior to the emergence of virtualization and cloud computing that has allowed for the shift to immutable infrastructure


**Shift to Immutable Infrastructure**

With Immutable Infrastructure, we leverage virtualization and cloud computing to allow for new deployment workflows in which new servers can be provisioned quickly and programmatically 
  * **virtualization** = when a simulated computing environment is created rather than a physical environment
  * **cloud computing** = using network of remote servers hosted via internet to store, manage + process data
  
  **Benefits of Cloud Computing**
  
 1. No Configuration Drift or Snowflake Serveres
    * config. drift = changes to software or hardware are mad ein data center environment and not systematically tracked
    * snowflake servers = a server that requires special configuration beyond the means of automated deployment (servers become difficult for replication because each snowflake server has a unique fragile part of overall infrastructure)
 2. Consistent staging environments and easy horizontal scaling
 3. Simple rollback and recovery processes

**REPOSITORY CONTENTS**

images:
* directory storing packer image template (image.pkr.hcl)

instances:
* main.tf : declaration of terraform resources, aws_instance resource, and output public ip address
* terraform.tfvars : region declaration
* variables.tf : terraform vars. w/ associated descriptions

scripts:
* setup.sh

Additional files:
* tf-packer : stores private key
* tf-packer.pub : contains ssh key


**NOTE:** This repository will be associated with an Instruqt lab tracK. The track will utilize this repo to show the process of deploying a machine image, updating the image, and showing it redeploy

**NOTE:** This is just a copy from the Companion code repository. Using the same template to get a feeling of building an image with packer and deploying it using terraform. There is a lot of different ways to build an image and use different provisioners for web-apps: I found this the most easiest way instead of using a .service file.. We just use a shell script and issue a go-get with a github link to your web-app.

**Resources:**
1. https://www.digitalocean.com/community/tutorials/what-is-immutable-infrastructure
2. https://www.sumologic.com/insight/mutable-immutable-infrastructure/
3. https://www.hashicorp.com/resources/what-is-mutable-vs-immutable-infrastructure
4. https://www.hashicorp.com/resources/modern-infrastructure-automation-with-packer-terraform-and-consul

**Web-App:** 
1. Make sure to set ssh-rsa keys: tf-packer and tf-packer.pub 
2. Also when ssh-ing to server the directory path to run the web-app is: 
3. /go/pkg/mod/github.com/sabinlehaci/go-web-app@v0.0.0-20220706130610-b856cb378344$ 


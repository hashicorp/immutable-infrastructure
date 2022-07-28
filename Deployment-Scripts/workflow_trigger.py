import os
import requests
import sys

TOKEN= str(sys.argv[0])
OWNER= str(sys.argv[1])
REPO= str(sys.argv[2])
Workflow_Name= str(sys.argv[3])


print( "the token value is")
def trigger_workflow(Workflow_Name):

      headers = {
        "Accept": "application/vnd.github.v3+json",
        "Authorization": f"token {TOKEN}",
      }

      data = {
        "event_type": Workflow_Name
      }


      responsevalue=requests.post(f"https://api.github.com/repos/{OWNER}/{REPO}/dispatches",json=data,headers=headers)
      print("The response message is ",responsevalue.content)    


trigger_workflow(Workflow_Name)

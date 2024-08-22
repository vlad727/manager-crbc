#!/bin/bash
## script allow to build, tag, push and restart new image 

TAG=$1
echo "Starting......"

echo "Prune images"
docker system prune --all --force

echo "============================================================================================================"
# build image
docker build -t manager-crbc . 
   if [ "$?" -ne 0 ]
   then
      echo "Nope!"
      exit 1
   else
      echo "Success build image..."
   fi
# tag image 
echo "============================================================================================================"
docker tag manager-crbc registry.apps.k8s.ose-prod.solution.sbt/vlku/manager-crbc:$TAG
   if [ "$?" -ne 0 ]
   then
      echo "Nope!"
      exit 1
   else
      echo "Success tag image..."
   fi
echo "============================================================================================================"
docker push  registry.apps.k8s.ose-prod.solution.sbt/vlku/manager-crbc:$TAG
   if [ "$?" -ne 0 ]
   then
      echo "Nope!"
      exit 1
   else
      echo "Success push image..."
   fi
echo "============================================================================================================"

# helm upgrade 
helm upgrade crbc chart/ --set image.tag="$TAG"
   if [ "$?" -ne 0 ]
   then
      echo "Failed to helm upgrade"
      exit 1
   else
      echo "Success push image..."
   fi
echo "============================================================================================================"
# wait until pods start 
sleep 5

# get logs from new pod
NEWPODNAME=$(k get pods -n manager-crbc  | awk 'NR!=1 {print $1}')

# get logs for new pod and container
k logs $NEWPODNAME -c manager-crbc   -n manager-crbc -f
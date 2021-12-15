#!/bin/bash

for sdir in io subnets testrtn ./ utils internal ; do 
  pushd $sdir
	for fl in *go; do 
		echo $sdir/$fl
    	sed -s s"/jpaulm/tyoung3/" $fl > /tmp/$fl && mv /tmp/$fl $fl
    done
  popd
done   

done

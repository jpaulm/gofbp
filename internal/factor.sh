#!/bin/bash

pfx="github.com/"
from="$pfx$1"; to="$pfx$2"
log=/tmp/factor.log
	
case $# in 
	2)
	echo FACTOR.SH: `date`  >> $log 
	for sdir in ../ ../io ../subnets ../testrtn  ../utils ../internal ; do 
  		pushd $sdir >/dev/null
			for fl in *go; do 
				echo $sdir/$fl >> $log
    			sed -s s"!$from!$to!" $fl > /tmp/$fl && mv /tmp/$fl $fl && count=$(($count+1))
    		done
  		popd >/dev/null
	done   
	echo $count hits
	 ;;
	 *)echo USAGE:  $0 FROM TO; exit 1;;
esac 	 

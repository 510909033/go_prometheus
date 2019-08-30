#!/bin/sh
function p() {
	# $1 msg  $2 level
	echo $1
}
function pull() {
	git stash save 'auto'
	if [[ $? != 0 ]];then
		echo "ERROR:git save "
		exit 1
	fi	
	echo "git stash save SUCCESS" 

	git status -s
	if [[ $? != 0 ]];then
		echo "ERROR:git status -s "
		exit 1
	fi	


	git pull origin $branch
	if [[ $? != 0 ]];then
		echo "ERROR:git pull origin $branch "
		exit 1
	fi	
	echo "git pull origin $branch SUCCESS"

	git stash pop
	if [[ $? != 0 ]];then
		echo "ERROR:git stash pop"
		exit 1
	fi	
	echo "git stash pop SUCCESS"
	echo "pull method OVER"
}

function push() {

}



branch=`git branch|grep '*'|awk '{print $2}'`
echo "branch=$branch"

echo "over"
exit 0
  

# get process id
pro_id=`ps -ef | grep ./scp294 | grep -v grep | awk '{print $2}'`
# stop process
if [ "$pro_id" != "" ];
then
	echo "SCP-294 is stopped, process id: $pro_id"
	kill -9 $pro_id
else
echo "SCP-294 was not started!"
fi
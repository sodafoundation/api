#!bin/bash

CHKLOG=/var/log/scan/hkconfig.log
LINUXLOG=/var/log/scan/linuxinstall.log
ISCSILOG=/var/log/scan/iscsi.log
ISCITARLOG=/var/log/scan/iscitar.log
LOGINLOG=/var/log/scan/login.log
SESSIONLOG=/var/log/scan/session.log
DISKLOG=/var/log/scan/disk.log
TARGETIPLOG=/var/log/scan/targetip.log
TARGETLOG=/var/log/scan/target.log

echo >$CHKLOG
echo >$LINUXLOG
echo >$ISCSILOG
echo >$ISCITARLOG
echo >$LOGINLOG
echo >$SESSIONLOG
echo >$DISKLOG
echo >$TARGETIPLOG
echo >$TARGETLOG

declare target

#install chkconfig
echo "Begin to Scan Volume"
sudo apt install sysv-rc-conf>$CHKLOG
sudo apt-get install chkconfig>$LINUXLOG

#echo out install log to stdio
echo "install result:"
echo "apt-get install chkconfig result:"
cat $LINUXLOG
echo "apt install sysv-rc-conf result:"
cat $CHKLOG

#check if sysv-rc-conf install is finish
count=1
while ( [ $? -ne 0 ] || [ `grep -c "newly installed" $CHKLOG` -ne 1 ] )
   do
   echo "sysv_rc_conf  install not finish"
   count=$((count+1))
   if [ $count -eq 6 ];then
       exit 1
   fi
   done
echo "sysv_rc_conf install finished"

#force link chkconfig
cp -f /usr/sbin/sysv-rc-conf /usr/sbin/chkconfig
echo `$?`
if [ $? -eq 0  ];then
  echo "link chkconfig sucsee!"
else
  echo "link chkconfig fail!"
fi

#check /usr/sbin/chkconfig is exist
echo "detail of /usr/sbin/chkconfig"
echo `ls /usr/sbin/chkconfig`

#GET TARGETIP 
echo >ip.log
echo >addr.log
ifconfig |grep -w "inet"|grep addr >>addr.log
cat addr.log|while read line
do
   array=(${line// / })
   for var in ${array[@]}
   do
      echo $var | grep "addr">>ip.log
   done
done
#echo ip out
echo "get ip log:"
cat ip.log

##Echo IP to log
cat ip.log|while read line
do
  TARGETIP=`echo $line | cut -d \: -f 2`
  TARGETIP=$TARGETIP:3260
  echo $TARGETIP>$TARGETIPLOG
done

echo "targerip log:"
cat $TARGETIPLOG

##get ip from log
cat $TARGETIPLOG |while read line
do
   ip=`echo $line`
   #find target
  chkconfig iscsi on
  chkconfig iscsi --list > $ISCSILOG
  echo >$TARGETLOG
  iscsiadm -m discovery -t sendtargets -p $ip>$ISCITARLOG

   #Check TARGETLOG
   echo "TARGET LOG SHOW:"
   cat $TARGETLOG
  #login target
cat $ISCITARLOG |while read line
  do
     a=$line
     echo $a
      target=`echo $a | cut -d \, -f 2` #iqn.2017-10.io.opensds:d3a3059d-7e31-4093-8c44-391528e748b0
      echo $target>$TARGETLOG
  #login

        echo "Login Out all session before login"
        iscsiadm -m node -U all

                #check if all session have been Login OUT

        echo `iscsiadm -m session`

        echo "login target:$target -p $ip"
        #if $target or ip ="",exit
        iscsiadm -m node –T $target -p $ip -l >$LOGINLOG

        echo "check login session:"
        echo `iscsiadm -m session`
   done

echo `grep -c "successful" $LOGINLOG`

if [ `grep -c "successful" $LOGINLOG` -ne 0 ];then
      echo "login target note success!"
   elif [ `grep -c "already present" $LOGINLOG` -eq 1 ];then
     echo "the not has been login in,please login out first!"
   else
     echo "login target note fail!"
fi
#view login session
iscsiadm -m session >$SESSIONLOG
if [ `awk '{print NR}' $SESSIONLOG|tail -n1` -eq 1 ];then
   echo "Have been Login in Target!"
fi

#show disk info
echo > $DISKLOG
fdisk -l|grep Disk >$DISKLOG
#if [ `grep -c "Disk /dev/sd: 2 GiB" $DISKLOG` -eq 1 ];then
#     echo "volume attachment successfully!"
#   else
#     echo "volume attachment fail!"
#fi
echo `grep "Disk /dev/sd" $DISKLOG |grep "1 GiB"`

#login out from the target 
iscsiadm -m node –T $target -p $ip -u
echo >$SESSIONLOG
iscsiadm -m session >>$SESSIONLOG
if [ -s $SESSIONLOG];then
     echo "LOGIN OUT FAIL"
   else
     echo "LOGIN OUT SUCCESS"
fi
echo "Scan volume end!"

##remove target
 echo "remove target:"
 # iscsiadm -m node -o delete -T $target -p $ip
 echo "Tagget Ip Log:"
 cat $TARGETIPLOG   #127.0.0.1:3260
echo "Target log:"
 cat $TARGETLOG  #1 iqn.2017-10.io.opensds:18f9ba11-543a-4bf2-b15e-10de92aba274

 iqn=`cat $TARGETLOG`
 tariqn=`echo $iqn | cut -d \  -f 2`
  echo $tariqn #iqn.2017-10.io.opensds:18f9ba11-543a-4bf2-b15e-10de92aba274

 tarip=`cat $TARGETIPLOG`
 iscsiadm -m node -o delete -T $iqn -p $tarip

 iscsiadm -m node > node.log
 if [[ ! -s node.log ]];then
   echo "remove target node fai"
 else
   echo "remove target node success"
 fi
done

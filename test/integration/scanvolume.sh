#!bin/bash

CHKLOG=/var/log/scan/hkconfig.log
#TARGETIP=10.10.3.156:3260
ISCSILOG=/var/log/scan/iscsi.log
TARGETLOG=/var/log/scan/target.log
LOGINLOG=/var/log/scan/login.log
SESSIONLOG=/var/log/scan/session.log
DISKLOG=/var/log/scan/disk.log
TARGETIPLOG=/var/log/scan/targetip.log
DISKRSLOG=/var/log/scan/diskrs.log

declare target

#install chkconfig

sudo apt install sysv-rc-conf>$CHKLOG

#check if install is finish

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

#link chkconfig
cp /usr/sbin/sysv-rc-conf /usr/sbin/chkconfig
if [ $? -eq 0  ];then
  echo "link chkconfig sucsee!"
else
  echo "link chkconfig fail!"
fi

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

##Echo IP to log
cat ip.log|while read line
do
  TARGETIP=`echo $line | cut -d \: -f 2`
  TARGETIP=$TARGETIP:3260
  echo $TARGETIP>$TARGETIPLOG
done

##get ip from log
cat $TARGETIPLOG |while read line
do
   ip=`echo $line`
   #find target
  chkconfig iscsi on
  chkconfig iscsi --list > $ISCSILOG
  iscsiadm -m discovery -t sendtargets -p $ip>$TARGETLOG

  #login target
cat $TARGETLOG |while read line
  do
     a=$line
     echo $a
      target=`echo $a | cut -d \, -f 2` #iqn.2017-10.io.opensds:d3a3059d-7e31-4093-8c44-391528e748b0
      echo $target
  #login
    iscsiadm -m node –T $ip，$target -l >$LOGINLOG
   done

if [ `grep -c "successful" $LOGINLOG` -eq 1 ];then
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
fdisk -l|grep Disk >$DISKLOG
#if [ `grep -c "Disk /dev/sd: 2 GiB" $DISKLOG` -eq 1 ];then
#     echo "volume attachment successfully!"
#   else
#     echo "volume attachment fail!"
#fi
grep "Disk /dev/sd" $DISKLOG |grep "2 GiB">$DISKRSLOG


#login out from the target 
iscsiadm -m node –T $target -p $ip -u
echo >$SESSIONLOG
iscsiadm -m session >>$SESSIONLOG
if [ -s $SESSIONLOG];then
     echo "LOGIN OUT FAIL"
   else
     echo "LOGIN OUT SUCCESS"
fi
done

#!bin/bash

CHKLOG=/var/log/scan/hkconfig.log
ISCSILOG=/var/log/scan/iscsi.log
TARGETLOG=/var/log/scan/target.log
LOGINLOG=/var/log/scan/login.log
SESSIONLOG=/var/log/scan/session.log
DISKLOG=/var/log/scan/disk.log
TARGETIPLOG=/var/log/scan/targetip.log

touch $CHKLOG
touch $ISCSILOG
touch $TARGETLOG
touch $LOGINLOG
touch $SESSIONLOG
touch $DISKLOG
touch $TARGETIPLOG

declare target
echo "Begin to Scan Attachment Volume"
#install chkconfig
sudo apt install sysv-rc-conf>$CHKLOG
echo "print install chkconfig result:"
cat $CHKLOG

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
echo "locate the sysv-rc-conf:"
echo `locate sysv-rc-conf`
cp /usr/sbin/sysv-rc-conf /usr/sbin/chkconfig
if [ $? -eq 0  ];then
  echo "link chkconfig sucsee!"
else
  echo "link chkconfig fail!"
fi


#GET TARGETIP
touch ip.log
touch addr.log
echo >ip.log
echo >addr.log
echo `ifconfig |grep -w "inet"|grep addr`
ifconfig |grep -w "inet"|grep addr >>addr.log
cat addr.log|while read line
do
   array=(${line// / })
   for var in ${array[@]}
   do
      echo `$var | grep "addr"`
      echo $var | grep "addr">>ip.log
   done
done
##Echo IP to log
echo "Get Ip:"
cat ip.log
echo "Get AddrIp:"
cat addr.log

echo "Get Target Ip:"
cat ip.log|while read line
do
  TARGETIP=`echo $line | cut -d \: -f 2`
  TARGETIP=$TARGETIP:3260
  echo $TARGETIP>$TARGETIPLOG
  cat $TARGETIPLOG
done

##get ip from log
echo "Get TARGETLOG:"
cat $TARGETIPLOG |while read line
do
   ip=`echo $line`
   #find target
  chkconfig iscsi on
  chkconfig iscsi --list > $ISCSILOG
  echo "GET ISCILOG:"
  cat $ISCSILOG
  iscsiadm -m discovery -t sendtargets -p $ip>$TARGETLOG
  
  cat $TARGETLOG
  #login target
 cat $TARGETLOG |while read line
  do
     a=$line
     echo $a
      target=`echo $a | cut -d \, -f 2` #iqn.2017-10.io.opensds:d3a3059d-7e31-4093-8c44-391528e748b0
      echo $target
  #login
    iscsiadm -m node –T $ip，$target -l >$LOGINLOG
    cat $LOGINLOG
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
echo `grep "Disk /dev/sd" $DISKLOG |grep "2 GiB"`

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



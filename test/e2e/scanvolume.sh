#!bin/bash

CHKLOG=hkconfig.log
LINUXLOG=linuxinstall.log
ISCSILOG=iscsi.log
ISCITARLOG=iscitar.log
LOGINLOG=login.log
SESSIONLOG=session.log
DISKLOG=disk.log
TARGETIPLOG=targetip.log
TARGETLOG=target.log

echo >$CHKLOG
echo >$LINUXLOG
echo >$ISCSILOG
echo >$ISCITARLOG
echo >$LOGINLOG
echo >$SESSIONLOG
echo >$DISKLOG
echo >$TARGETIPLOG
echo >$TARGETLOG


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
  echo $line,${#line}
  #if len(line) !=0,get it and break
  if [ ${#line} -ne 0 ];then
     TARGETIP=`echo $line | cut -d \: -f 2`
     echo $TARGETIP:3260 >$TARGETIPLOG
     break
  fi
done

echo "targerip log:"
cat $TARGETIPLOG

#show isci status 
chkconfig iscsi on
chkconfig iscsi --list > $ISCSILOG

#get the target
ip=`cat $TARGETIPLOG`
iscsiadm -m discovery -t sendtargets -p $ip>$ISCITARLOG

#Check TARGETLOG
echo "TARGET LOG SHOW:"
cat $ISCITARLOG

#login target
 echo "begin to login target..."
 iscitar=`cat $ISCITARLOG`
  isciqn=`echo $iscitar | cut -d \, -f 2` #1 iqn.2017-10.io.opensds:d3a3059d-7e31-4093-8c44-391528e748b0
 targetiqn=`echo $iscitar | cut -d \  -f 2`
 echo "target ipn:"
 echo $targetiqn
 echo $targetiqn > $TARGETLOG
##print log to CI
 echo "TARGET LOG(1 IPN):"
 cat $TARGETLOG
##Login out all node
  iscsiadm -m node -U all
##Login
  iscsiadm -m node –T $targetiqn -p $ip -l >$LOGINLOG
##Print Login Log to CI
  echo "Login Log:"
  cat $LOGINLOG
##Check if exist 'successful' in $LOGINLOG
echo `grep -c "successful" $LOGINLOG`
if [ `grep -c "successful" $LOGINLOG` -ne 0 ];then
      echo "login target note success!"
   else
     echo "login target note fail!"
fi
#view login session
echo `iscsiadm -m session`
iscsiadm -m session >$SESSIONLOG
echo "Print Session Log after Login:"
cat $SESSIONLOG
##Check the number of Row
if [ `awk '{print NR}' $SESSIONLOG|tail -n1` -eq 1 ];then
   echo "Have been Login in Target!"
fi
#show disk info
fdisk -l|grep Disk >$DISKLOG
cat $DISKLOG
echo `grep "Disk /dev/sd" $DISKLOG |grep "1 GiB"`

#login out from the target 
iscsiadm -m node –T $targetiqn -p $ip -u
echo >$SESSIONLOG
iscsiadm -m session >>$SESSIONLOG
echo "Print Session Log after Login Out:"
cat $SESSIONLOG
##remove target
 echo "start to remove target:"
 # iscsiadm -m node -o delete -T $target -p $ip
 echo "Tagget Ip Log:"
 cat $TARGETIPLOG   #127.0.0.1:3260
echo "Target log:"
 echo $targetiqn  #iqn.2017-10.io.opensds:18f9ba11-543a-4bf2-b15e-10de92aba274

# iqn=`cat $TARGETLOG`
# tariqn=`echo $iqn | cut -d \  -f 2`
#  echo $tariqn #iqn.2017-10.io.opensds:18f9ba11-543a-4bf2-b15e-10de92aba274

 tarip=`cat $TARGETIPLOG`
 iscsiadm -m node -o delete -T $targetiqn -p $tarip

 iscsiadm -m node > node.log
 if [[ ! -s node.log ]];then
   echo "remove target node fai"
 else
   echo "remove target node success"
 fi
echo "Scan volume end!"


import { Component, OnInit, Input } from '@angular/core';

@Component({
  selector: 'app-suspension-frame',
  templateUrl: './suspension-frame.component.html',
  styleUrls: [

  ]
})
export class SuspensionFrameComponent implements OnInit {
    data=[];
    policyName:string;
  @Input()
    set policy(policy: any) {
        let extra = policy[1];
        this.policyName = policy[0];
        if(this.policyName === "QoS"){
            let maxIpos ="MaxIOPS = " + extra[":provisionPolicy"].ioConnectivityLoS.maxIOPS + " IOPS/TB";
            this.data.push(maxIpos);
            let maxBWS = "MaxBWS = " + extra[":provisionPolicy"].ioConnectivityLoS.maxBWS + " MBPS/TB";
            this.data.push(maxBWS);
        }else if(this.policyName === "Replication"){
            let type ="Type = " + extra[":replicationPolicy"].dataProtectionLoS.replicaTypes;
            this.data.push(type);
            let mode = "Mode = " + extra[":replicationPolicy"].replicaInfos.replicaUpdateMode;
            this.data.push(mode);
            let Period = "Period = " + extra[":replicationPolicy"].replicaInfos.replicationPeriod +" Minutes";
            this.data.push(Period);
            let Bandwidth = "Bandwidth = " + extra[":replicationPolicy"].replicaInfos.replicationBandwidth +" MBPS/TB";
            this.data.push(Bandwidth);
            let RGO = "RGO = " + extra[":replicationPolicy"].dataProtectionLoS.recoveryGeographicObject;
            this.data.push(RGO);
            let RTO = "RTO = " + extra[":replicationPolicy"].dataProtectionLoS.recoveryTimeObjective;
            this.data.push(RTO);
            let RPO = "RPO = " + extra[":replicationPolicy"].dataProtectionLoS.recoveryPointObjective;
            this.data.push(RPO);
            let Consistency = "Consistency = " + extra[":replicationPolicy"].replicaInfos.consistencyEnabled;
            this.data.push(Consistency);
        }else{
            let schedule ="Schedule = " + extra[":snapshotPolicy"].schedule.occurrence;
            this.data.push(schedule);
            let execution = "Execution Time = " + extra[":snapshotPolicy"].schedule.datetime.split("T")[1] ;
            this.data.push(execution);
            let Retention  = "Retention  = " + (extra[":snapshotPolicy"].retention["number"] ? extra[":snapshotPolicy"].retention["number"]: (extra[":snapshotPolicy"].retention.duration+" Days"));
            this.data.push(Retention );
        }
    };
  constructor() { }

  ngOnInit() {
    // console.log(this.policy);
  }
}

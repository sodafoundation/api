import {Component, Input, OnInit} from '@angular/core';
import { I18NService } from 'app/shared/api';
import { VolumeService ,ReplicationService} from './../../volume.service';
import { ConfirmationService,ConfirmDialogModule} from '../../../../components/common/api';

@Component({
  selector: 'app-replication-list',
  templateUrl: './replication-list.component.html',
  providers: [ConfirmationService],
  styleUrls: [

  ]
})
export class ReplicationListComponent implements OnInit {

    @Input() volumeId;
    volume={
        name:""
    };
    replication ={
        name:"replication",
        replicationPeriod:30,
        id:""
    };
    arrowEnable = true;
    showReplication:boolean=false;
    //[0]:status:enble;[1]:status:disabled;[2]:status:failover;[3]:status:other
    swichMode = [
        {
            disableBtnDisplay:true,
            disableBtnDisabled:false,
            enableBtnDisplay:false,
            enableBtnDisabled:true,
            failoverBtnDisabled:false,
        },
        {
            disableBtnDisplay:false,
            disableBtnDisabled:true,
            enableBtnDisplay:true,
            enableBtnDisabled:false,
            failoverBtnDisabled:false,
        },
        {
            disableBtnDisplay:false,
            disableBtnDisabled:true,
            enableBtnDisplay:true,
            enableBtnDisabled:false,
            failoverBtnDisabled:true,
        },
        {
            disableBtnDisplay:true,
            disableBtnDisabled:true,
            enableBtnDisplay:false,
            enableBtnDisabled:true,
            failoverBtnDisabled:true,
        }
    ];
    operationStatus:any;
    constructor(
       public I18N:I18NService,
       private VolumeService:VolumeService,
       private replicationService:ReplicationService,
       private confirmationService:ConfirmationService
              ) {  }

    ngOnInit() {
        this.operationStatus = this.swichMode[3];
        this.getAllReplicationsDetail();
    }
    getAllReplicationsDetail(){
        this.replicationService.getAllReplicationsDetail().subscribe((resRep)=>{
            let replications = resRep.json();
            replications.forEach(element => {
                if(element.primaryVolumeId == this.volumeId){
                    this.getVolumeById(this.volumeId);
                    this.replication = element;
                    //ReplicationStatus
                    switch(this.replication['replicationStatus']){
                        case "enabled":
                            this.operationStatus = this.swichMode[0];
                            this.arrowEnable = true;
                            break;
                        case "disabled":
                            this.operationStatus = this.swichMode[1];
                            this.arrowEnable = false;
                            break;
                        case "failed_over":
                            this.operationStatus = this.swichMode[2];
                            this.arrowEnable = false;
                            break;
                        default:
                            this.operationStatus = this.swichMode[3];
                            this.arrowEnable = false;
                    }
                    this.showReplication = true;
                }
                if(element.secondaryVolumeId == this.volumeId){
                    this.replication = element;
                    //ReplicationStatus
                    switch(this.replication['replicationStatus']){
                        case "enabled":
                            this.operationStatus = this.swichMode[0];
                            this.arrowEnable = true;
                            break;
                        case "disabled":
                            this.operationStatus = this.swichMode[1];
                            this.arrowEnable = false;
                            break;
                        case "failed_over":
                            this.operationStatus = this.swichMode[2];
                            this.arrowEnable = false;
                            break;
                        default:
                            this.operationStatus = this.swichMode[3];
                            this.arrowEnable = false;
                    }
                    this.getVolumeById(element.primaryVolumeId);
                    this.showReplication = true;
                }
            });
        });
    }
    getVolumeById(volumeId){
        this.VolumeService.getVolumeById(volumeId).subscribe((res) => {
            this.volume = res.json();
        });
    }
    getReplicationByVolumeId = function(volumeId){
        let param = {
            "key": "PrimaryVolumeId",
            "value":volumeId
        }
        this.replicationService.getReplicationDetailByVolumeId(param).subscribe((res) => {
            var data = res.json();
            if(data.length !== 0){
                this.replication = data[0];
                this.showReplication = true;
            }else{
                this.showReplication = false;
            }
        });
    }
    disableReplication(){
        let msg = "<div>Are you sure you want to disable the Replication?</div><h3>[ "+ this.replication.name +" ]</h3>";
        let header ="Disable Replication";
        let acceptLabel = "Disable";
        let warming = false;
        this.confirmDialog([msg,header,acceptLabel,warming,"disable"])
    }
    enableReplication(){
        let msg = "<div>Are you sure you want to enable the Replication?</div><h3>[ "+ this.replication.name +" ]</h3>";
        let header ="Enable Replication";
        let acceptLabel = "Enable";
        let warming = false;
        this.confirmDialog([msg,header,acceptLabel,warming,"enable"])
    }
    failoverReplication(){
        let msg = "<div>Are you sure you want to failover the Replication?</div><h3>[ "+ this.replication.name +" ]</h3>";
        let header ="Failover Replication";
        let acceptLabel = "Failover";
        let warming = true;
        this.confirmDialog([msg,header,acceptLabel,warming,"failover"])
    }
    deleteReplication(){
        let msg = "<div>Are you sure you want to delete the Replication?</div><h3>[ "+ this.replication.name +" ]</h3>";
        let header ="Delete Replication";
        let acceptLabel = "Delete";
        let warming = true;
        this.confirmDialog([msg,header,acceptLabel,warming,"delete"])
    }
    confirmDialog([msg,header,acceptLabel,warming=true,func]){
        this.confirmationService.confirm({
            message: msg,
            header: header,
            acceptLabel: acceptLabel,
            isWarning: warming,
            accept: ()=>{
                try {
                    switch(func){
                        case "disable":
                            this.replicationService.disableReplication(this.replication.id).subscribe((res)=>{
                                this.getAllReplicationsDetail();
                            });
                            break;
                        case "delete":
                            this.replicationService.deleteReplication(this.replication.id).subscribe((res)=>{
                                this.getAllReplicationsDetail();
                            });
                            break;
                        case "failover":
                            this.replicationService.failoverReplication(this.replication.id).subscribe((res)=>{
                                this.getAllReplicationsDetail();
                            });
                            break;
                        case "enable":
                            this.replicationService.enableReplication(this.replication.id).subscribe((res)=>{
                                this.getAllReplicationsDetail();
                            });
                            break;
                    }
                }
                catch (e) {
                    console.log(e);
                }
                finally {
                    this.getAllReplicationsDetail();
                }
            },
            reject:()=>{}
        })
    }
}

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
    showReplication:boolean=false;
    constructor(
       public I18N:I18NService,
       private VolumeService:VolumeService,
       private replicationService:ReplicationService,
       private confirmationService:ConfirmationService
              ) {  }

    ngOnInit() {
      this.getVolumeById(this.volumeId);
      this.getReplicationByVolumeId(this.volumeId)
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
                            this.replicationService.disableReplication(this.replication.id).subscribe((res)=>{});
                            break;
                        case "delete":
                            this.replicationService.deleteReplication(this.replication.id).subscribe((res)=>{
                                this.getReplicationByVolumeId(this.volumeId);
                            });
                            break;
                        case "failover":
                            this.replicationService.failoverReplication(this.replication.id).subscribe((res)=>{
                                this.getReplicationByVolumeId(this.volumeId);
                            });
                            break;
                    }
                }
                catch (e) {
                    console.log(e);
                }
                finally {
                    this.getReplicationByVolumeId(this.volumeId);
                }
            },
            reject:()=>{}
        })
    }
}

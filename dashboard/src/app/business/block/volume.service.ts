import { Injectable } from '@angular/core';
import { I18NService, HttpService, ParamStorService } from '../../shared/api';
import { Observable } from 'rxjs';

@Injectable()
export class VolumeService {
  constructor(
    private http: HttpService,
    private paramStor: ParamStorService
  ) { }

  url = 'v1beta/{project_id}/block/volumes';

  //Create volume
  createVolume(param) {
    return this.http.post(this.url, param);
  }

  //Update volume
  modifyVolume(id,param) {
    let modifyUrl = this.url + '/' + id
    return this.http.put(modifyUrl, param);
  }

  //Delete volume
  deleteVolume(id): Observable<any> {
    let deleteUrl = this.url + '/' + id
    return this.http.delete(deleteUrl);
  }

  //Search all volumes
  getVolumes(): Observable<any> {
    return this.http.get(this.url);
  }

  //Search volume
  getVolumeById(id): Observable<any> {
    let url = this.url + '/' + id;
    return this.http.get(url);
  }
  //Search volume by groupId
  getVolumeByGroupId(id): Observable<any> {
    let url = this.url + '?GroupId=' + id;
    return this.http.get(url);
  }

  //Create volumesGroup
  createVolumesGroup(param) {
    return this.http.post(this.url, param);
  }

  //Delete volumesGroup
  deleteVolumesGroup(id): Observable<any> {
    let deleteUrl = this.url + '/' + id
    return this.http.delete(deleteUrl);
  }

  //Search volumesGroups
  getVolumesGroups(): Observable<any> {
    return this.http.get(this.url);
  }
  expandVolume(id,param):Observable<any> {
      let expandVolumeUrl = 'v1beta/{project_id}/block/volumes' + '/' + id + "/resize"
      return this.http.post(expandVolumeUrl,param);
  }
}

@Injectable()
export class SnapshotService {
  constructor(
    private http: HttpService,
    private paramStor: ParamStorService
  ) { }

  url = 'v1beta/{project_id}/block/snapshots';

  //Create snapshot
  createSnapshot(param) {
    return this.http.post(this.url, param);
  }

  //Delete snapshot
  deleteSnapshot(id){
    let url = this.url + "/" + id;
    return this.http.delete(url);
  }

  //Search snapshot
  getSnapshots(filter?){
    let url = this.url;
    if(filter){
      url = this.url + "?" + filter.key + "=" + filter.value;
    }
    return this.http.get(url);
  }

  //Update snapshot
  modifySnapshot(id,param){
    let url = this.url + "/" + id;
    return this.http.put(url,param);
  }
}
@Injectable()
export class ReplicationService {
    constructor(
        private http: HttpService,
        private paramStor: ParamStorService
    ) { }

    project_id = this.paramStor.CURRENT_TENANT().split("|")[1];
    replicationUrl = 'v1beta/{project_id}/block/replications';
    //create replication
    createReplication(param){
        let url = this.replicationUrl;
        return this.http.post(url,param);
    }
    getReplicationDetailByVolumeId(filter?){
        let url = this.replicationUrl+"/detail";
        if(filter){
            url = url + "?" + filter.key + "=" + filter.value;
        }
        return this.http.get(url);
    }
    disableReplication(param){
        let url = this.replicationUrl+"/"+param+"/disable";
        return this.http.post(url,param);
    }
    enableReplication(param){
      let url = this.replicationUrl+"/"+param+"/enable";
      return this.http.post(url,param);
  }
    failoverReplication(id){
        let url = this.replicationUrl+"/"+id+"/failover";
        let param = {
            "allowAttachedVolume": true,
            "secondaryBackendId": "default"
        }
        return this.http.post(url,param);
    }
    deleteReplication(param){
        let url = this.replicationUrl+"/"+param;
        return this.http.delete(url);
    }
    //get all replications
    getAllReplicationsDetail(){
      let url = this.replicationUrl+"/detail";
      return this.http.get(url);
  }
}
@Injectable()
export class VolumeGroupService {
    constructor(
        private http: HttpService,
        private paramStor: ParamStorService
    ) { }

    project_id = this.paramStor.CURRENT_TENANT().split("|")[1];
    volumeGroupUrl = 'v1beta/{project_id}/block/volumeGroup';
    //create volume group
    createVolumeGroup(param){
        let url = this.volumeGroupUrl;
        return this.http.post(url,param);
    }
    //get volume group
    getVolumeGroups(): Observable<any> {
        return this.http.get(this.volumeGroupUrl);
    }
    //delete volume group
    deleteVolumeGroup(groupId): Observable<any> {
      let url = this.volumeGroupUrl+"/" + groupId
      return this.http.delete(url);
    }
    //modify volume group
    modifyVolumeGroup(groupId,param): Observable<any> {
      let url = this.volumeGroupUrl+"/" + groupId
      return this.http.put(url,param);
    }
    //get volume group by id
    getVolumeGroupById(groupId): Observable<any> {
      let url = this.volumeGroupUrl+"/"+groupId;
      return this.http.get(url);
    }
    //add or remove volumes 
    addOrRemovevolumes(groupId,param): Observable<any> {
      let url = this.volumeGroupUrl+"/"+groupId;
      return this.http.put(url,param);
    }
}

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

  //创建 volume
  createVolume(param) {
    return this.http.post(this.url, param);
  }

  //修改 volume
  modifyVolume(id,param) {
    let modifyUrl = this.url + '/' + id
    return this.http.put(modifyUrl, param);
  }

  //删除 volume
  deleteVolume(id): Observable<any> {
    let deleteUrl = this.url + '/' + id
    return this.http.delete(deleteUrl);
  }

  //查询 volumes
  getVolumes(): Observable<any> {
    return this.http.get(this.url);
  }

  //查询指定Id volume
  getVolumeById(id): Observable<any> {
    let url = this.url + '/' + id;
    return this.http.get(url);
  }

  //创建 volumesGroup
  createVolumesGroup(param) {
    return this.http.post(this.url, param);
  }

  //删除 volumesGroup
  deleteVolumesGroup(id): Observable<any> {
    let deleteUrl = this.url + '/' + id
    return this.http.delete(deleteUrl);
  }

  //查询 volumesGroups
  getVolumesGroups(): Observable<any> {
    return this.http.get(this.url);
  }
  expandVolume(id,param):Observable<any> {
      let expandVolumeUrl = 'v1beta/{project_id}/volumes' + '/' + id + "/action"
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

  //创建 snapshot
  createSnapshot(param) {
    return this.http.post(this.url, param);
  }

  //删除 snapshot
  deleteSnapshot(id){
    let url = this.url + "/" + id;
    return this.http.delete(url);
  }

  //查询 snapshot
  getSnapshots(filter?){
    let url = this.url;
    if(filter){
      url = this.url + "?" + filter.key + "=" + filter.value;
    }
    console.log(url);
    return this.http.get(url);
  }

  //修改 snapshot
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
    replicationUrl = 'v1beta/'+ this.project_id +'/block/replications';
    //create replication
    createReplication(param){
        let url = this.replicationUrl;
        return this.http.post(url,param);
    }
}
@Injectable()
export class VolumeGroupService {
    constructor(
        private http: HttpService,
        private paramStor: ParamStorService
    ) { }

    project_id = this.paramStor.CURRENT_TENANT().split("|")[1];
    volumeGroupUrl = 'v1beta/'+ this.project_id +'/block/volumeGroup';
    //create volume group
    createVolumeGroup(param){
        let url = this.volumeGroupUrl;
        return this.http.post(url,param);
    }
    //查询 volumes
    getVolumeGroups(): Observable<any> {
        return this.http.get(this.volumeGroupUrl);
    }
}

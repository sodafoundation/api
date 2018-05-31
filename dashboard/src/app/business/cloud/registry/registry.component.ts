import { Component, OnInit } from '@angular/core';

@Component({
  selector: 'app-registry',
  templateUrl: './registry.component.html',
  styleUrls: ['./registry.component.scss']
})
export class RegistryComponent implements OnInit {

  registryServiceDisplay = false;
  label;
  cloudServices = [];
  serviceTypeOptions = [];
  regionOptions = [];
  selectedType = '';
  selectedRegion = '';
  

  constructor() { }

  ngOnInit() {
    //界面文本
    this.label = {
      name: 'Service Name',
      type: 'Service Type',
      region: 'Region',
      accessKey: 'AWS Access Key',
      secretKey: 'Secret Key'
    };

    //type下拉框的项
    this.serviceTypeOptions = [
      { label: 'AWS S3', value: { id: 0, name: 'New York', code: 'NY' } },
      { label: 'Microsoft Azure Blob Storage', value: { id: 1, name: 'New York', code: 'NY' } },
      { label: 'Huawei OBS', value: { id: 2, name: 'Rome', code: 'RM' } }
    ];

    //region下拉框的项
    this.regionOptions = [
      { label: 'CN North', value: { id: 0, name: 'New York', code: 'NY' } },
      { label: 'CN South', value: { id: 1, name: 'New York', code: 'NY' } },
      { label: 'Huawei OBS', value: { id: 2, name: 'Rome', code: 'RM' } }
    ];

    //查询回来已有的云服务
    this.cloudServices = [
      {
        name:'service_for_analytics',
        region:'EU(Paris)',
        type:'AWS S3'
      },
      {
        name:'service_for_finance',
        region:'CN North',
        type:'Huawei OBS'
      },
      {
        name:'service_for_media',
        region:'North Europe',
        type:'Microsoft Azure Blob Storage'
      }
    ];





  }

  showRegistryService() {
    this.registryServiceDisplay = true;
  }

  registryCloud() {
    //http注册新的云
    alert("registry");
  }

  getServiceType() {
    //http获取service type
  }

  getRegions() {
    //http获取regions
  }

  getCloudServices() {
    //http获取cloudServices
  }
}

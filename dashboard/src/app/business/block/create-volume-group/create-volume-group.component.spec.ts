import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { CreateVolumeGroupComponent } from './create-volume-group.component';

describe('CreateVolumeGroupComponent', () => {
  let component: CreateVolumeGroupComponent;
  let fixture: ComponentFixture<CreateVolumeGroupComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ CreateVolumeGroupComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(CreateVolumeGroupComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});

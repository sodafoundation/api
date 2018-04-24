import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { CreateVolumeComponent } from './create-volume.component';

describe('CreateVolumeComponent', () => {
  let component: CreateVolumeComponent;
  let fixture: ComponentFixture<CreateVolumeComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ CreateVolumeComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(CreateVolumeComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});

import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ReplicationGroupComponent } from './replication-group.component';

describe('ReplicationGroupComponent', () => {
  let component: ReplicationGroupComponent;
  let fixture: ComponentFixture<ReplicationGroupComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ReplicationGroupComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ReplicationGroupComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});

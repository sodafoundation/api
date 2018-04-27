import { Component, OnInit, Input } from '@angular/core';
import { FormControl, FormGroup, FormBuilder, Validators, ValidatorFn, AbstractControl } from '@angular/forms';
import { SnapshotService } from './../../volume.service';

@Component({
  selector: 'app-snapshot-list',
  templateUrl: './snapshot-list.component.html',
  styleUrls: ['./snapshot-list.component.scss']
})
export class SnapshotListComponent implements OnInit {

  @Input() volumeId;
  label;
  selectedSnapshot;
  selectedSnapshots = [];
  snapshortfilter;
  snapshots;
  snapshotfilter;
  modifyDialogDisplay = false;
  createSnapshotDisplay = false;
  snapshotFormGroup;
  modifyFormGroup;

  constructor(
    private SnapshotService: SnapshotService,
    private fb: FormBuilder
  ) {
    this.snapshotFormGroup = this.fb.group({
      "name": ["", Validators.required],
      "description": ["", Validators.maxLength(200)]
    });
    this.modifyFormGroup = this.fb.group({
      "name": ['', Validators.required]
    });
  }

  ngOnInit() {
    this.label = {
      name: 'name',
      volume: 'Volume',
      description: 'Description'
    }
    this.getSnapshots(
      {
        key: 'volumeId',
        value: this.volumeId
      }
    );
  }

  showCreateSnapshot() {
    this.createSnapshotDisplay = true;
  }

  createSnapshot() {
    let param = {
      name: this.snapshotFormGroup.value.name,
      volumeId: this.volumeId,
      description: this.snapshotFormGroup.value.description
    }
    this.SnapshotService.createSnapshot(param).subscribe((res) => {
      this.getSnapshots(
        {
          key: 'volumeId',
          value: this.volumeId
        }
      );
      this.createSnapshotDisplay = false;
    });
  }

  batchDeleteSnapshort() {
    if (this.selectedSnapshots) {
      this.selectedSnapshots.forEach(snapshot => {
        // this.deleteSnapshot(snapshot.id);
        console.log(this.selectedSnapshots);
      });
    }
  }

  getSnapshots(filter?) {
    this.SnapshotService.getSnapshots(filter).subscribe((res) => {
      this.snapshots = res.json();
    });
  }

  deleteSnapshot(id) {
    this.SnapshotService.deleteSnapshot(id).subscribe((res) => {
      // this.snapshots = res.json();
      console.log('delete success');
    });
  }


  showModifyDialog(id) {
    this.modifyDialogDisplay = true;
  }

}

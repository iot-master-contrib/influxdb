import { Component, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';
import { NzMessageService } from 'ng-zorro-antd/message';
import { RequestService } from '../request.service';
@Component({
  selector: 'app-setting',
  templateUrl: './setting.component.html',
  styleUrls: ['./setting.component.scss'],
})
export class SettingComponent implements OnInit {
  group!: FormGroup;
  loading = false;
  query = {};
  dbData = [];
  constructor(
    private fb: FormBuilder,
    private router: Router,
    private route: ActivatedRoute,
    private rs: RequestService,
    private msg: NzMessageService
  ) {
    this.load();
  }
  switchValue = false;

  ngOnInit(): void {
    // this.rs.get(`config`).subscribe(res => {
    //   //let data = res.data;
    //   this.build(res.data)
    // })

    this.build();
  }

  load() {
    this.rs.get(`/app/influx/api/config/influxdb`).subscribe((res) => {
      this.dbData = res.data; 
      this.group.patchValue({
        bucket: res.data.Bucket,
        org: res.data.Org,
        url: res.data.Url,
        token: res.data.Token
      });
    });
  }

  build(obj?: any) {
    obj = obj || {};
    this.group = this.fb.group({
      bucket: [obj.bucket || '', []],
      org: [obj.org || '', []],
      url: [obj.url || '', []],
      token: [obj.token || '', []],
    });
  }

  submit() {
    if (this.group.valid) {
      this.group.patchValue({ LogLevel: Number(this.group.value.LogLevel) });
      this.rs
        .post(`/app/influx/api/config/influxdb`, this.group.value)
        .subscribe((res) => {
          this.msg.success('保存成功');
        });

      return;
    } else {
      Object.values(this.group.controls).forEach((control) => {
        if (control.invalid) {
          control.markAsDirty();
          control.updateValueAndValidity({ onlySelf: true });
        }
      });
    }
  }
}

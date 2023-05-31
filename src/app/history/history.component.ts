import { Component } from '@angular/core';
import { RequestService } from 'src/app/request.service';
import {
  FormGroup,
  UntypedFormBuilder,
  UntypedFormControl,
  UntypedFormGroup,
  Validators,
} from '@angular/forms';
import { DatePipe } from '@angular/common';
import { NzMessageService } from 'ng-zorro-antd/message';
import * as dayjs from 'dayjs';
@Component({
  selector: 'app-history',
  templateUrl: './history.component.html',
  styleUrls: ['./history.component.scss'],
  providers: [DatePipe],
})
export class HistoryComponent {
  chart:any;
  validateForm!: FormGroup;
  href!: string;
  option: any = {};
  query: any = {};
  searchData: any[] = [];
  searchTotal = 1;
  proData: any[] = [];
  proTotal = 1;
  devData: any[] = [];
  devTotal = 1;
  properties: any[] = [];
  selectPropertyObj: { label: string } = { label: '' };
  isVisible!: boolean
  loading = true;
  isFirst = true;
  constructor(
    private readonly datePipe: DatePipe,
    private rs: RequestService,
    private fb: UntypedFormBuilder,
    private msg: NzMessageService
  ) {
    let base = +new Date(1988, 9, 3);
    let oneDay = 24 * 3600 * 1000;
    let data: any = [[base, Math.random() * 300]];
    for (let i = 1; i < 20000; i++) {
      let now = new Date((base += oneDay));
      data.push([
        +now,
        Math.round((Math.random() - 0.5) * 20 + data[i - 1][1]),
      ]);
    }
    for (let i = 0; i < data.length; i++) data[i][1] += 300;

    this.option = {
      tooltip: {
        trigger: 'axis',
        position: function (pt: any) {
          return [pt[0], '10%'];
        },
      },
      title: {
        left: 'center',
        text: '[能耗]变化曲线',
      },
      toolbox: {
        feature: {
          dataZoom: {
            yAxisIndex: 'none',
          },
          restore: {},
          saveAsImage: {},
        },
      },
      xAxis: {
        type: 'time',
        boundaryGap: false,
      },
      yAxis: {
        type: 'value',
        boundaryGap: [0, '100%'],
      },
      dataZoom: [
        {
          type: 'inside',
          start: 70,
          end: 80,
        },
        {
          start: 70,
          end: 80,
        },
      ],
      series: [
        {
          name: 'Fake Data',
          type: 'line',
          smooth: true,
          symbol: 'none',
          areaStyle: {},
          data: data,
        },
      ],
    };
    this.load();
  }


  chartInit(ec:any){ this.chart=ec}


  search() {

    this.query = {};
    if (this.validateForm.valid) {
      let value = this.validateForm.value;
      this.query = {
        start: dayjs(value.strEnd[0]).toISOString(),
        end: dayjs(value.strEnd[1]).toISOString(),
        window: value.window ? value.window + value.winTp : '1h',
        fn: 'last',
      };
      this.rs
        .get(
          `/app/influxdb/api/query/${value.pid}/${value.id}/${value.name}`,
          this.query
        )
        .subscribe((res) => {
          this.searchData = res.data;
          this.searchTotal = res.total;
          //图表渲染
         // this.chart.setOption(this.option);
        })
        .add(() => {
          this.loading = false;
        });
      return;
    } else {
      Object.values(this.validateForm.controls).forEach((control) => {
        if (control.invalid) {
          control.markAsDirty();
          control.updateValueAndValidity({ onlySelf: true });
        }
      });
    }
  }
  change(mes: any) {
    // this.resetForm()
    let vl = this.validateForm;
    vl.patchValue({ id: '' });
    vl.patchValue({ name: '' });
    let tmp = this.proData.filter((item) => {
      if (item.id === mes) return item;
    });
    this.properties = tmp[0].properties;
    this.query.filter = { product_id: mes };
    this.rs
      .post('/app/classify/api/device/search', this.query)
      .subscribe((res) => {
        const resData = res && Array.isArray(res.data) ? res.data : [];
        this.devData = resData;
        this.devTotal = res.total;
        if (this.isFirst && resData.length) {
          this.validateForm.patchValue({
            id: resData[0].id,
            name: this.properties[0].name
          });
          this.search();
          this.isFirst = false;
        }
      })
      .add(() => {
        this.loading = false;
      });
  }
  load() {
    this.query = {};
    this.rs
      .post('/api/product/search', this.query)
      .subscribe((res) => {
        const resData = res && Array.isArray(res.data) ? res.data : [];
        resData.filter((item: any) => {
          item.checked = false;
        });
        this.proData = resData;
        this.proTotal = res.total;
        if (resData.length) {
          this.validateForm.patchValue({ pid: resData[0].id || '' });
        }
      })
      .add(() => {
        this.loading = false;
      });

  }
  edit(mes: any) {
    this.isVisible = true;

  }

  resetForm(): void {
    //this.validateForm.reset();
    let vl = this.validateForm;
    vl.patchValue({ strEnd: [] });
    vl.patchValue({ window: null });
    //vl.patchValue({ fn: '' });
  }

  ngOnInit(): void {
    this.validateForm = this.fb.group({
      pid: ['', [Validators.required]],
      id: ['', [Validators.required]],
      name: ['', [Validators.required]],
      strEnd: [
        [
          this.datePipe.transform(
            new Date(new Date().getTime() - 7 * 24 * 3600 * 1000),
            'yyyy-MM-dd'
          ),
          this.datePipe.transform(new Date(), 'yyyy-MM-dd'),
        ],
      ],
      window: [1],
      winTp: ['h'],
      fn: ['last'],
    });
  }
  checked = false;
  indeterminate = false;
  setOfCheckedId = new Set<number>();
  cancel() {
    this.msg.info('取消删除!');
  }
  updateCheckedSet(id: number, checked: boolean): void {
    if (checked) {
      this.setOfCheckedId.add(id);
    } else {
      this.setOfCheckedId.delete(id);
    }
  }

  handlePropertyChange(value: string) {
    const obj = this.properties.find((item) => item.name === value);
    this.selectPropertyObj = obj;
  }
}

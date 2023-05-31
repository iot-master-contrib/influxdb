import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { NZ_I18N } from 'ng-zorro-antd/i18n';
import { zh_CN } from 'ng-zorro-antd/i18n';
import { APP_BASE_HREF, registerLocaleData } from '@angular/common';
import zh from '@angular/common/locales/zh';
import { FormsModule , ReactiveFormsModule} from '@angular/forms';
import { HttpClientModule } from '@angular/common/http';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { SettingComponent } from './setting/setting.component';
import { HistoryComponent } from './history/history.component';
import { NzGridModule } from 'ng-zorro-antd/grid';
import { NzDatePickerModule } from 'ng-zorro-antd/date-picker';
import { NzSelectModule } from 'ng-zorro-antd/select';
import { NzFormModule } from 'ng-zorro-antd/form';
import { NzInputNumberModule } from 'ng-zorro-antd/input-number';
import { NzButtonModule } from 'ng-zorro-antd/button';
import { NzSpaceModule } from 'ng-zorro-antd/space';
import { NzLayoutModule } from 'ng-zorro-antd/layout';
import { NgxEchartsModule } from 'ngx-echarts';
import { NzMessageModule } from 'ng-zorro-antd/message';
import { NzCardModule } from 'ng-zorro-antd/card';
import { NzSwitchModule } from 'ng-zorro-antd/switch';
registerLocaleData(zh);

@NgModule({
  declarations: [
    AppComponent,
    SettingComponent,
    HistoryComponent
  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    NzSelectModule,
    NzSpaceModule ,
    NzSwitchModule ,
    NzCardModule ,
    NzLayoutModule,
    NzFormModule,
    NzMessageModule,
    NgxEchartsModule.forRoot({
      echarts: () => import('echarts')
    }), 
    NzInputNumberModule,
    NzButtonModule,
    ReactiveFormsModule,
    NzGridModule,
    NzDatePickerModule,
    FormsModule,
    HttpClientModule,
    BrowserAnimationsModule
  ],
  providers: [
    { provide: APP_BASE_HREF, useValue: '/app/influxdb/' },
    { provide: NZ_I18N, useValue: zh_CN }
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }

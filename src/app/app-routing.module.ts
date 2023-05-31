import { SettingComponent } from './setting/setting.component';
import { HistoryComponent } from './history/history.component';
import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { PageNotFoundComponent } from './base/page-not-found/page-not-found.component';
const routes: Routes = [
  { path: '', pathMatch: 'full', redirectTo: 'history' },
  { path: 'setting', component: SettingComponent },
  { path: 'history', component: HistoryComponent },
  { path: '**', component: PageNotFoundComponent },
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule],
})
export class AppRoutingModule {}

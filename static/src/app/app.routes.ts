// app.routes.ts
import { Routes } from '@angular/router';
import { HomeComponent } from './home.component';
import { ProfileComponent } from './profile.component';
import { LoginComponent } from './login.component';
import { RegisterComponent } from './register.component';

export const routes: Routes = [
  { path: 'login', component: LoginComponent },
  { path: 'register', component: RegisterComponent },
  { path: 'home', component: HomeComponent },          
  { path: 'profile', component: ProfileComponent } 
];

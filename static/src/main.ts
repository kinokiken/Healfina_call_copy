import { bootstrapApplication } from '@angular/platform-browser';
import { ApplicationConfig } from '@angular/core';
import { provideRouter } from '@angular/router';
import { routes } from './app/app.routes';  // Ваши маршруты
import { AppComponent } from './app/app.component';
import { provideHttpClient, withInterceptorsFromDi } from '@angular/common/http';

const appConfig: ApplicationConfig = {
  providers: [
    provideRouter(routes),      // Подключаем маршруты
    provideHttpClient(),        // Подключаем HttpClient
  ]
};

bootstrapApplication(AppComponent, appConfig).catch((err) => console.error(err));

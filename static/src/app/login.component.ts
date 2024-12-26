import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { HttpClient } from '@angular/common/http';
import { Router } from '@angular/router'; 
import { RouterModule } from '@angular/router';

@Component({
  selector: 'app-login',
  standalone: true,
  imports: [CommonModule, FormsModule, RouterModule],
  template: `
    <h2>Авторизация</h2>
    <div>
      <label>User ID:
        <input [(ngModel)]="userID" type="text"/>
      </label>
    </div>
    <div>
      <label>Password:
        <input [(ngModel)]="password" type="password"/>
      </label>
    </div>
    <button (click)="login()">Войти</button>
    <button (click)="register()">Пройти регистрацию</button>
  `
})
export class LoginComponent {
  userID = '';
  password = '';

  constructor(private http: HttpClient, private router: Router) { }

  login() {
    if (!this.userID || !this.password) {
      alert('Введите user_id и пароль!');
      return;
    }

    // Отправим POST /api/login
    const body = { user_id: this.userID, password: this.password };
    this.http.post('/api/login', body).subscribe({
      next: (resp) => {
        console.log('Логин успешен:', resp);
        this.router.navigate(['/home']);
      },
      error: (err) => {
        console.error('Ошибка логина:', err);
        alert('Неверные данные или ошибка сервера');
      }
    });
  }

  register() {
    this.router.navigate(['/register']); 
  }
}
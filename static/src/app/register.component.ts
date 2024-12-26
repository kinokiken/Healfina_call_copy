import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { HttpClient } from '@angular/common/http';
import { Router } from '@angular/router'; 

@Component({
  selector: 'app-register',
  standalone: true,
  imports: [CommonModule, FormsModule],
  template: `
    <h2>Регистрация</h2>
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
    <button (click)="register()">Зарегистрироваться</button>
    <button (click)="login()">Уже есть аккаунт</button>
  `
})
export class RegisterComponent {
  userID = '';
  password = '';

  constructor(private http: HttpClient, private router: Router) { }

  register() {
    if (!this.userID || !this.password) {
      alert('Введите user_id и пароль!');
      return;
    }

    // POST /api/register { "user_id":"...", "password":"..." }
    const body = { user_id: this.userID, password: this.password };
    this.http.post('/api/register', body).subscribe({
      next: (resp) => {
        console.log('Регистрация успешна:', resp);
        alert('Пароль установлен, теперь можете логиниться');
        this.router.navigate(['/login']);
      },
      error: (err) => {
        console.error('Ошибка регистрации:', err);
        alert('Не удалось зарегистрироваться!');
      }
    });
  }

  login() {
    this.router.navigate(['/login']); 
  }

}
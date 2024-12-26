// src/app/home.component.ts (standalone пример)
import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterLink } from '@angular/router';
import { FontAwesomeModule } from '@fortawesome/angular-fontawesome';
import { AudioService } from './services/audio.service';
import { faPhone, faBellSlash, faToggleOn, faToggleOff, faUser } from '@fortawesome/free-solid-svg-icons';

@Component({
  selector: 'app-home',
  standalone: true,
  imports: [CommonModule, RouterLink, FontAwesomeModule],
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.css'],
})
export class HomeComponent {
    darkMode: boolean = false;
    isSessionActive: boolean = false;

    faPhone = faPhone;
    faUser = faUser;
    faBellSlash = faBellSlash;
    faToggleOn = faToggleOn;
    faToggleOff = faToggleOff;

    constructor(private audioService: AudioService) { 
        console.log('AppComponent: Конструктор вызван');
    }

    toggleDarkMode() {
        this.darkMode = !this.darkMode;
        console.log('toggleDarkMode() вызван. Новое состояние:', this.darkMode);
        if (this.darkMode) {
          document.body.classList.add('dark-mode');
        } else {
          document.body.classList.remove('dark-mode');
        }
        fetch('/set_dark_mode', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({ darkMode: this.darkMode })
        })
        .then(response => response.json())
        .then(data => {
          console.log('Dark mode updated:', data);
        })
        .catch(error => console.error('Ошибка:', error));
    }

    toggleCall() {
        console.log('toggleCall() вызван. Текущее состояние:', this.isSessionActive);
        if (this.isSessionActive) {
          this.audioService.stopSession();
          console.log('Завершение звонка...');
        } else {
          this.audioService.startSession()
            .then(() => {
              console.log('Начало звонка...');
            })
            .catch(err => {
              console.error("Ошибка при запуске сессии:", err);
              alert("Не удалось начать сессию. Проверьте разрешения или подключение к серверу.");
            });
        }
        this.isSessionActive = !this.isSessionActive;
        console.log('Новое состояние звонка:', this.isSessionActive);
      }
}

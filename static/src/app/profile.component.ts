import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';                 // ← Добавили
import { HttpClient } from '@angular/common/http';
import { FontAwesomeModule } from '@fortawesome/angular-fontawesome';
import { faTrash, faPen } from '@fortawesome/free-solid-svg-icons';
import { Router } from '@angular/router';
import { NgbModule } from '@ng-bootstrap/ng-bootstrap';

@Component({
  selector: 'app-profile',
  standalone: true,
  imports: [CommonModule, FormsModule, FontAwesomeModule, NgbModule],     
  templateUrl: './profile.component.html',
  styleUrls: ['./profile.component.css'],
})
export class ProfileComponent implements OnInit {
  user: any = null;

  faTrash = faTrash;
  faPen = faPen;
  isTableView: boolean = true; 
  records = [
    { id: '1', time: new Date(2023, 5, 24, 14, 30), record: 'Запись 1' },
    { id: '2', time: new Date(2023, 4, 13, 11, 15), record: 'Запись 2' },
    { id: '3', time: new Date(2023, 6, 5, 9, 0), record: 'Запись 3' },
  ];
  originalRecords: any[] = []

  sortByDateAsc = true; 
  sortByTimeAsc = true; 
  searchQuery: string = '';
  filterPeriod: string = 'all';

  filterRecordsByPeriod(period: string) {
    const now = new Date();
    let startDate: Date;
  
    if (period === 'all') {
      return [...this.originalRecords];  // Возвращаем все записи
    }
  
    switch (period) {
      case 'day':
        startDate = new Date();
        startDate.setDate(now.getDate() - 1); // За последний день
        break;
      case 'week':
        startDate = new Date();
        startDate.setDate(now.getDate() - 7); // За последнюю неделю
        break;
      case 'month':
        startDate = new Date();
        startDate.setMonth(now.getMonth() - 1); // За последний месяц
        break;
      default:
        return [...this.originalRecords]; // Если период не распознан, возвращаем все записи
    }
  
    // Возвращаем только те записи, которые соответствуют выбранному периоду
    return this.originalRecords.filter(rec => new Date(rec.time) >= startDate);
  }

  // Применение фильтра по периоду и поиску
  applyFilter() {
    let filteredRecords = this.filterRecordsByPeriod(this.filterPeriod);
  
    if (this.searchQuery.trim() !== '') {
      filteredRecords = filteredRecords.filter(rec =>
        rec.record.toLowerCase().includes(this.searchQuery.toLowerCase())
      );
    }
  
    this.records = filteredRecords;  // Обновляем только отображаемые записи
  }

  // Метод, который вызывается при смене периода
  filterByPeriod(period: string) {
    this.filterPeriod = period;
    this.applyFilter();  // Применяем фильтрацию с новым периодом
  }

  sortByDate() {
    this.sortByDateAsc = !this.sortByDateAsc; 

    this.records.sort((a, b) => {
      return this.sortByDateAsc
        ? a.time > b.time ? 1 : a.time < b.time ? -1 : 0
        : b.time > a.time ? 1 : b.time < a.time ? -1 : 0;
    });
  }

  // sortByTime() {
  //   this.sortByTimeAsc = !this.sortByTimeAsc; 

  //   this.records.sort((a, b) => {
  //     return this.sortByTimeAsc
  //       ? a.time.getTime() > b.time.getTime() ? 1 : a.time.getTime() < b.time.getTime() ? -1 : 0
  //       : b.time.getTime() > a.time.getTime() ? 1 : b.time.getTime() < a.time.getTime() ? -1 : 0;
  //   });
  // }

  switchToTableView() {
    this.isTableView = true;
  }

  switchToCardView() {
    this.isTableView = false;
  }

  newRecordText: string = '';

  constructor(private http: HttpClient, private router: Router) {  
    console.log("ProfileComponent: Конструктор вызван");
  }

  ngOnInit() {
    console.log("ProfileComponent: ngOnInit вызван");
    this.loadProfile();
    this.loadRecords();
  }

  goBack() {
    this.router.navigate(['/home']);  
  }

  loadProfile() {
    this.http.get('/profile').subscribe({
      next: (data) => {
        this.user = data;
      },
      error: (err) => {
        console.error('Ошибка при загрузке профиля:', err);
      }
    });
  }

  loadRecords() {
    this.http.get<any[]>('/records').subscribe({
      next: (data) => {
        this.originalRecords = data || [];  // Сохраняем оригинальные записи
        this.records = [...this.originalRecords];
        console.log('Загружены записи:', this.records);
      },
      error: (err) => {
        console.error('Ошибка при загрузке записей:', err);
      }
    });
  }

  deleteRecord(recordID: string) {
    const confirmed = window.confirm('Вы уверены, что хотите удалить эту запись?');
    if (confirmed) {
      this.http.delete(`/records/delete/${recordID}`).subscribe({
        next: (response) => {
          console.log('Запись успешно удалена:', response);
          this.records = this.records.filter(rec => rec.id !== recordID);
        },
        error: (err) => {
          console.error('Ошибка при удалении записи:', err);
          alert('Не удалось удалить запись. Попробуйте позже.');
        }
      });
    }
  }

  editRecord(rec: any) {
    console.log('Редактировать:', rec);
  
    const newText = prompt('Введите новый текст записи:', rec.record);
    if (newText === null) {

      return;
    }
    const trimmed = newText.trim();
    if (!trimmed) {
      alert('Текст не может быть пустым.');
      return;
    }
  
    this.http.put(`/records/update/${rec.id}`, { updatedData: trimmed }).subscribe({
      next: (response) => {
        console.log('Запись обновлена успешно', response);

        rec.record = trimmed;
      },
      error: (err) => {
        console.error('Ошибка при обновлении записи:', err);
        alert('Не удалось обновить запись. Попробуйте позже.');
      }
    });
  }

  addRecord() {
    if (!this.newRecordText.trim()) {
      alert('Введите текст записи!');
      return;
    }
    const body = { Record: this.newRecordText };
    console.log(body)

    this.http.put('/records/add', body).subscribe({
      next: (resp) => {
        console.log('Запись добавлена:', resp);
        this.newRecordText = '';
        this.loadRecords();
      },
      error: (err) => {
        console.error('Ошибка при добавлении записи:', err);
        alert('Не удалось добавить запись!');
      }
    });
  }

  searchRecords() {
    if (this.searchQuery.trim() === '') {
      this.loadRecords();  // Если строка поиска пустая, загружаем все записи
    } else {
      this.http.get<any[]>(`/records/search?searchQuery=${this.searchQuery}`).subscribe({
        next: (data) => {
          this.records = data || [];
        },
        error: (err) => {
          console.error('Ошибка при поиске записей:', err);
        }
      });
    }
}
}

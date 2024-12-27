document.addEventListener("DOMContentLoaded", () => {
    // Добавляем боковую панель
    const sidebar = `
        <nav class="bg-dark text-white p-3" id="sidebar" style="width: 250px; height: 100vh; position: fixed;">
            <h4>Меню</h4>
            <ul class="nav flex-column">
                <li class="nav-item"><a class="nav-link text-white" href="index.html">Главная</a></li>
                <li class="nav-item"><a class="nav-link text-white" href="scan.html">Сканирование</a></li>
                <li class="nav-item"><a class="nav-link text-white" href="reports.html">Отчёты</a></li>
                <li class="nav-item"><a class="nav-link text-white" href="hosts_groups.html">Хосты</a></li>
            </ul>
        </nav>`;
    document.body.insertAdjacentHTML('afterbegin', sidebar);

    // Добавляем подвал
    // const footer = `
    //     <footer class="text-center py-3 bg-dark text-white" style="margin-left: 250px;">
    //         <p>&copy; 2024 Система управления сканированием. Все права защищены.</p>
    //     </footer>`;
    // document.body.insertAdjacentHTML('beforeend', footer);
});

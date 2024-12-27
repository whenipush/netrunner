import { NmapAPI } from "./api.js";

async function loadTasks() {
    try {
        const tasks = await NmapAPI.getAllTasks(); // Получаем данные задач из API
        const taskList = document.getElementById("taskList");
        taskList.innerHTML = ""; // Очищаем текущий список

        tasks.forEach((task) => {
            // Создаём элемент списка
            const li = document.createElement("li");
            li.className = "list-group-item d-flex justify-content-between align-items-center";
            li.innerHTML = `
                <div>
                    <strong>${task.number_task}</strong> - ${task.name} <br>
                    <small>Статус: ${task.status}, Прогресс: ${task.percent}%</small>
                </div>
                <button class="btn btn-primary btn-sm view-task-details" data-id="${task.ID}">Детали</button>
            `;

            // Добавляем элемент в список
            taskList.appendChild(li);
        });

        // Добавляем обработчики для кнопок "Детали"
        document.querySelectorAll(".view-task-details").forEach((button) => {
            button.addEventListener("click", (event) => {
                const taskId = event.target.getAttribute("data-id");
                showTaskDetails(taskId, tasks);
            });
        });
    } catch (error) {
        console.error("Ошибка загрузки задач:", error.message);
    }
}

// Показ деталей задачи
function showTaskDetails(taskId, tasks) {
    const task = tasks.find((t) => t.ID === parseInt(taskId));
    if (task) {
        const details = `
            <h5>${task.number_task}: ${task.name}</h5>
            <p><strong>Статус:</strong> ${task.status}</p>
            <p><strong>Прогресс:</strong> ${task.percent}%</p>
            <p><strong>Скрипт:</strong> ${task.script}</p>
            <p><strong>Хосты:</strong></p>
            <ul>
                ${task.hosts.map((host) => `<li>${host.ip}</li>`).join("")}
            </ul>
        `;
        document.getElementById("taskDetails").innerHTML = details;
    } else {
        console.error("Задача не найдена");
    }
}

// Инициализация
document.addEventListener("DOMContentLoaded", () => {
    loadTasks();
});

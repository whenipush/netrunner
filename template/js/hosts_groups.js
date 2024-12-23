import { HostAPI, GroupAPI } from "./api.js";

// Обновление списка хостов
async function loadHosts() {
    try {
        const hosts = await HostAPI.getAllHosts();
        const hostsList = document.getElementById("hostsList");
        hostsList.innerHTML = ""; // Очищаем текущий список

        hosts.forEach((host) => {
            const li = document.createElement("li");
            li.className = "list-group-item d-flex justify-content-between align-items-center";
            li.textContent = host.ip;

            // Кнопка удаления
            const deleteButton = document.createElement("button");
            deleteButton.className = "btn btn-danger btn-sm";
            deleteButton.textContent = "Удалить";
            deleteButton.onclick = async () => {
                try {
                    await HostAPI.deleteHost(host.id);
                    loadHosts(); // Перезагрузка списка
                } catch (error) {
                    console.error("Ошибка при удалении хоста:", error.message);
                }
            };

            li.appendChild(deleteButton);
            hostsList.appendChild(li);
        });
    } catch (error) {
        console.error("Ошибка загрузки хостов:", error.message);
    }
}

// Обновление списка групп
async function loadGroups() {
    try {
        const groups = await GroupAPI.getAllGroups();
        const groupsList = document.getElementById("groupsList");
        groupsList.innerHTML = ""; // Очищаем текущий список

        groups.forEach((group) => {
            const li = document.createElement("li");
            li.className = "list-group-item d-flex justify-content-between align-items-center";
            li.textContent = group.name;

            // Кнопка удаления
            const deleteButton = document.createElement("button");
            deleteButton.className = "btn btn-danger btn-sm";
            deleteButton.textContent = "Удалить";
            deleteButton.onclick = async () => {
                try {
                    await GroupAPI.deleteGroup(group.id);
                    loadGroups(); // Перезагрузка списка
                } catch (error) {
                    console.error("Ошибка при удалении группы:", error.message);
                }
            };

            li.appendChild(deleteButton);
            groupsList.appendChild(li);
        });
    } catch (error) {
        console.error("Ошибка загрузки групп:", error.message);
    }
}

// Обработка формы добавления хоста
document.getElementById("addHostForm").addEventListener("submit", async (event) => {
    event.preventDefault();
    const hostIP = document.getElementById("hostIP").value;

    try {
        await HostAPI.createHost({ ip: hostIP });
        document.getElementById("hostIP").value = ""; // Очистка поля
        loadHosts(); // Перезагрузка списка
    } catch (error) {
        console.error("Ошибка при добавлении хоста:", error.message);
    }
});

// Обработка формы добавления группы
document.getElementById("createGroupForm").addEventListener("submit", async (event) => {
    event.preventDefault();
    const groupName = document.getElementById("groupName").value;

    try {
        await GroupAPI.createGroup({ name: groupName });
        document.getElementById("groupName").value = ""; // Очистка поля
        loadGroups(); // Перезагрузка списка
    } catch (error) {
        console.error("Ошибка при добавлении группы:", error.message);
    }
});

// Инициализация
document.addEventListener("DOMContentLoaded", () => {
    loadHosts();
    loadGroups();
});
